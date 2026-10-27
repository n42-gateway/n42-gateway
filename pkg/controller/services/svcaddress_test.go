/*
Copyright 2026 The HAProxy Ingress Controller Authors.

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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/n42-gateway/n42-gateway/pkg/controller/config"
	convtypes "github.com/n42-gateway/n42-gateway/pkg/converters/types"
	"github.com/n42-gateway/n42-gateway/pkg/utils"
)

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
			objects:   []client.Object{gwClass, testGateway("gw1", "n42")},
			expGwAddr: map[string][]gatewayv1.GatewayStatusAddress{"gw1": nil},
		},
		"full sync updates unchanged gateway": {
			objects:   []client.Object{gwClass, testGateway("gw1", "n42")},
			fullSync:  true,
			expGwAddr: map[string][]gatewayv1.GatewayStatusAddress{"gw1": gwAddress(lb)},
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
				cfg.PublishAddressIPs = nil
				s.curr = nil
			}
			if test.fullSync {
				s.requireFullSync()
			}

			err := s.checkChanged(ctx, utils.NewTimer(nil), &convtypes.ChangedObjects{Links: test.links})
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

func testGateway(name, class string) *gatewayv1.Gateway {
	return &gatewayv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: name},
		Spec:       gatewayv1.GatewaySpec{GatewayClassName: gatewayv1.ObjectName(class)},
	}
}
