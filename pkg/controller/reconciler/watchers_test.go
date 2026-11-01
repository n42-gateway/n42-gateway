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

package reconciler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"github.com/n42-gateway/n42-gateway/pkg/controller/config"
)

func TestControllerPodUpdatePredicate(t *testing.T) {
	now := &metav1.Time{Time: time.Now()}
	pod := func(ip string, deletion *metav1.Time) *api.Pod {
		return &api.Pod{
			ObjectMeta: metav1.ObjectMeta{Namespace: "ingress", Name: "p2", DeletionTimestamp: deletion},
			Status:     api.PodStatus{PodIP: ip},
		}
	}
	testCases := map[string]struct {
		old, new *api.Pod
		exp      bool
	}{
		"no-change":         {old: pod("10.0.0.1", nil), new: pod("10.0.0.1", nil), exp: false},
		"pod-ip-changed":    {old: pod("", nil), new: pod("10.0.0.1", nil), exp: true},
		"start-terminating": {old: pod("10.0.0.1", nil), new: pod("10.0.0.1", now), exp: true},
	}
	w := &watchers{cfg: &config.Config{ControllerPod: types.NamespacedName{Namespace: "ingress", Name: "p1"}}}
	var podHandler *hdlr
	for _, h := range w.handlersCore() {
		if _, ok := h.typ.(*api.Pod); ok {
			podHandler = h
		}
	}
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.exp, podHandler.pr[0].Update(event.UpdateEvent{ObjectOld: test.old, ObjectNew: test.new}))
		})
	}
}
