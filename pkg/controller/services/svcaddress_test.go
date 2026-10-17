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
	"fmt"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	api "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"

	"github.com/n42-gateway/n42-gateway/pkg/controller/config"
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
