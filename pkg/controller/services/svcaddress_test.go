/*
Copyright 2026 The N42 Gateway Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/n42-gateway/n42-gateway/pkg/controller/config"
	convtypes "github.com/n42-gateway/n42-gateway/pkg/converters/types"
	"github.com/n42-gateway/n42-gateway/pkg/utils"
)

func TestReadCurrentLBNodeIPs(t *testing.T) {
	pod := func(name, node string) *api.Pod {
		return &api.Pod{
			ObjectMeta: metav1.ObjectMeta{Namespace: "ingress", Name: name, Labels: map[string]string{"app": "n42"}},
			Spec:       api.PodSpec{NodeName: node},
		}
	}
	node := func(name, ip string) *api.Node {
		return &api.Node{
			ObjectMeta: metav1.ObjectMeta{Name: name},
			Status:     api.NodeStatus{Addresses: []api.NodeAddress{{Type: api.NodeExternalIP, Address: ip}}},
		}
	}
	terminating := func(p *api.Pod) *api.Pod {
		p.DeletionTimestamp = &metav1.Time{Time: time.Now()}
		p.Finalizers = []string{"test/keep"}
		return p
	}
	testCases := map[string]struct {
		objs    []client.Object
		nodeErr string
		listErr bool
		expLB   []networking.IngressLoadBalancerIngress
		expErr  bool
	}{
		"all-pods-resolve": {
			objs:  []client.Object{pod("p1", "n1"), pod("p2", "n2"), node("n1", "10.0.0.1"), node("n2", "10.0.0.2")},
			expLB: []networking.IngressLoadBalancerIngress{{IP: "10.0.0.1"}, {IP: "10.0.0.2"}},
		},
		"skip-pending-pod": {
			objs:  []client.Object{pod("p1", "n1"), pod("p2", ""), node("n1", "10.0.0.1")},
			expLB: []networking.IngressLoadBalancerIngress{{IP: "10.0.0.1"}},
		},
		"skip-deleted-node": {
			objs:  []client.Object{pod("p1", "n1"), pod("p2", "gone"), node("n1", "10.0.0.1")},
			expLB: []networking.IngressLoadBalancerIngress{{IP: "10.0.0.1"}},
		},
		"skip-terminating-pod": {
			objs:  []client.Object{pod("p1", "n1"), terminating(pod("p2", "n2")), node("n1", "10.0.0.1"), node("n2", "10.0.0.2")},
			expLB: []networking.IngressLoadBalancerIngress{{IP: "10.0.0.1"}},
		},
		"node-read-failure-keeps-status": {
			objs:    []client.Object{pod("p1", "n1"), pod("p2", "n2"), node("n1", "10.0.0.1"), node("n2", "10.0.0.2")},
			nodeErr: "n2",
			expErr:  true,
		},
		"pod-list-failure-keeps-status": {
			objs:    []client.Object{pod("p1", "n1"), node("n1", "10.0.0.1")},
			listErr: true,
			expErr:  true,
		},
	}
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			cli := fake.NewClientBuilder().WithObjects(test.objs...).WithInterceptorFuncs(interceptor.Funcs{
				Get: func(ctx context.Context, c client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					if _, isNode := obj.(*api.Node); isNode && key.Name == test.nodeErr {
						return fmt.Errorf("connection refused")
					}
					return c.Get(ctx, key, obj, opts...)
				},
				List: func(ctx context.Context, c client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
					if test.listErr {
						return fmt.Errorf("connection refused")
					}
					return c.List(ctx, list, opts...)
				},
			}).Build()
			cfg := &config.Config{
				ControllerPod:         types.NamespacedName{Namespace: "ingress", Name: "p1"},
				ControllerPodSelector: labels.SelectorFromSet(labels.Set{"app": "n42"}),
			}
			s := &svcAddress{
				log:   logr.Discard(),
				cfg:   cfg,
				cli:   cli,
				cache: &c{ctx: ctx, config: cfg, client: cli},
			}
			lb, err := s.readCurrentLB(ctx)
			if test.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expLB, lb)
		})
	}
}

func TestCheckChanged(t *testing.T) {
	lb := []networking.IngressLoadBalancerIngress{{IP: "192.0.2.10"}}
	staleLB := []networking.IngressLoadBalancerIngress{{IP: "192.0.2.99"}}
	gwClass := &gatewayv1.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{Name: "n42"},
		Spec:       gatewayv1.GatewayClassSpec{ControllerName: "n42-gateway.github.io/controller"},
	}
	testCases := map[string]struct {
		objects   []client.Object
		links     convtypes.TrackingLinks
		emptyLB   bool
		fullSync  bool
		expIngLB  map[string][]networking.IngressLoadBalancerIngress
		expGwAddr map[string][]gatewayv1.GatewayStatusAddress
	}{
		"changed ingress is synced": {
			objects:  []client.Object{testIngress("ing1", "n42", nil)},
			links:    convtypes.TrackingLinks{convtypes.ResourceIngress: {"default/ing1"}},
			expIngLB: map[string][]networking.IngressLoadBalancerIngress{"ing1": lb},
		},
		"unchanged ingress is not synced": {
			objects:  []client.Object{testIngress("ing1", "n42", nil)},
			expIngLB: map[string][]networking.IngressLoadBalancerIngress{"ing1": nil},
		},
		"full sync updates unchanged ingress": {
			objects:  []client.Object{testIngress("ing1", "n42", nil)},
			fullSync: true,
			expIngLB: map[string][]networking.IngressLoadBalancerIngress{"ing1": lb},
		},
		"full sync updates changed and unchanged ingresses": {
			objects:  []client.Object{testIngress("ing1", "n42", nil), testIngress("ing2", "n42", nil)},
			links:    convtypes.TrackingLinks{convtypes.ResourceIngress: {"default/ing1"}},
			fullSync: true,
			expIngLB: map[string][]networking.IngressLoadBalancerIngress{"ing1": lb, "ing2": lb},
		},
		"full sync replaces stale address": {
			objects:  []client.Object{testIngress("ing1", "n42", staleLB)},
			fullSync: true,
			expIngLB: map[string][]networking.IngressLoadBalancerIngress{"ing1": lb},
		},
		"full sync clears stale address when published address is empty": {
			objects:  []client.Object{testIngress("ing1", "n42", staleLB)},
			emptyLB:  true,
			fullSync: true,
			expIngLB: map[string][]networking.IngressLoadBalancerIngress{"ing1": nil},
		},
		"full sync skips ingress of another class": {
			objects:  []client.Object{testIngress("ing1", "other", nil)},
			fullSync: true,
			expIngLB: map[string][]networking.IngressLoadBalancerIngress{"ing1": nil},
		},
		"unchanged gateway is not synced": {
			objects:   []client.Object{gwClass, testGateway("gw1", "n42", nil)},
			expGwAddr: map[string][]gatewayv1.GatewayStatusAddress{"gw1": nil},
		},
		"changed gateway is synced": {
			objects:   []client.Object{gwClass, testGateway("gw1", "n42", nil)},
			links:     convtypes.TrackingLinks{convtypes.ResourceGateway: {"default/gw1"}},
			expGwAddr: map[string][]gatewayv1.GatewayStatusAddress{"gw1": gwAddress(lb)},
		},
		"full sync updates unchanged gateway": {
			objects:   []client.Object{gwClass, testGateway("gw1", "n42", nil)},
			fullSync:  true,
			expGwAddr: map[string][]gatewayv1.GatewayStatusAddress{"gw1": gwAddress(lb)},
		},
		"full sync updates changed and unchanged gateways": {
			objects:   []client.Object{gwClass, testGateway("gw1", "n42", nil), testGateway("gw2", "n42", nil)},
			links:     convtypes.TrackingLinks{convtypes.ResourceGateway: {"default/gw1"}},
			fullSync:  true,
			expGwAddr: map[string][]gatewayv1.GatewayStatusAddress{"gw1": gwAddress(lb), "gw2": gwAddress(lb)},
		},
		"full sync replaces stale gateway address": {
			objects:   []client.Object{gwClass, testGateway("gw1", "n42", gwAddress(staleLB))},
			fullSync:  true,
			expGwAddr: map[string][]gatewayv1.GatewayStatusAddress{"gw1": gwAddress(lb)},
		},
		"full sync clears stale gateway address when published address is empty": {
			objects:   []client.Object{gwClass, testGateway("gw1", "n42", gwAddress(staleLB))},
			emptyLB:   true,
			fullSync:  true,
			expGwAddr: map[string][]gatewayv1.GatewayStatusAddress{"gw1": nil},
		},
		"full sync skips gateway of another class": {
			objects:   []client.Object{gwClass, testGateway("gw1", "other", nil)},
			fullSync:  true,
			expGwAddr: map[string][]gatewayv1.GatewayStatusAddress{"gw1": nil},
		},
	}
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			scheme := runtime.NewScheme()
			require.NoError(t, clientgoscheme.AddToScheme(scheme))
			require.NoError(t, gatewayv1.Install(scheme))
			cli := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(test.objects...).
				WithStatusSubresource(&networking.Ingress{}, &gatewayv1.Gateway{}).
				Build()
			cfg := &config.Config{
				ControllerName:    "n42-gateway.github.io/controller",
				HasGatewayV1:      true,
				IngressClass:      "n42",
				WatchIngress:      true,
				PublishAddressIPs: []string{"192.0.2.10"},
			}
			s := initSvcAddress(ctx, cfg, cli, &c{ctx: ctx, config: cfg, client: cli})
			s.run = true
			// a previous pass already published lb
			s.curr = lb
			if test.emptyLB {
				// no controller pod to read node IPs from
				cfg.PublishAddressIPs = nil
				cfg.ControllerPodSelector = labels.Everything()
				s.curr = nil
			}

			err := s.checkChanged(ctx, utils.NewTimer(nil), &convtypes.ChangedObjects{Links: test.links, NeedFullSync: test.fullSync})
			require.NoError(t, err)

			for ingName, expLB := range test.expIngLB {
				ing := &networking.Ingress{}
				require.NoError(t, cli.Get(ctx, types.NamespacedName{Namespace: "default", Name: ingName}, ing))
				assert.Equal(t, expLB, ing.Status.LoadBalancer.Ingress, ingName)
			}
			for gwName, expAddr := range test.expGwAddr {
				gw := &gatewayv1.Gateway{}
				require.NoError(t, cli.Get(ctx, types.NamespacedName{Namespace: "default", Name: gwName}, gw))
				assert.Equal(t, expAddr, gw.Status.Addresses, gwName)
			}
		})
	}
}

func testIngress(name, class string, lb []networking.IngressLoadBalancerIngress) *networking.Ingress {
	ing := &networking.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   "default",
			Name:        name,
			Annotations: map[string]string{"kubernetes.io/ingress.class": class},
		},
	}
	ing.Status.LoadBalancer.Ingress = lb
	return ing
}

func testGateway(name, class string, addr []gatewayv1.GatewayStatusAddress) *gatewayv1.Gateway {
	gw := &gatewayv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: name},
		Spec:       gatewayv1.GatewaySpec{GatewayClassName: gatewayv1.ObjectName(class)},
	}
	gw.Status.Addresses = addr
	return gw
}
