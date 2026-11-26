/*
Copyright 2024 The HAProxy Ingress Controller Authors.

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

package gateway

import (
	"testing"

	api "k8s.io/api/core/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	hatypes "github.com/n42-gateway/n42-gateway/pkg/haproxy/types"
)

type recordingAnnotationReader struct {
	pathLinks []*hatypes.PathLink
}

func (r *recordingAnnotationReader) ReadAnnotations(_ *hatypes.Backend, _ []*api.Service, pathLinks []*hatypes.PathLink) {
	r.pathLinks = append(r.pathLinks, pathLinks...)
}

func TestTLSPassthroughAppliesServiceAnnotations(t *testing.T) {
	passthrough := gatewayv1.TLSModePassthrough
	c := setup(t)
	_ = c.createGatewayClass("haproxy")
	g := c.createGateway1("default/web", "l1:8443")
	route := c.createTLSRoute1("default/web", "web", "echoserver:8443")
	c.createService1("default/echoserver", "8443", "172.17.0.11")
	g.Spec.Listeners[0].TLS = &gatewayv1.ListenerTLSConfig{Mode: &passthrough}
	route.Spec.Hostnames = append(route.Spec.Hostnames, "domain.local")

	reader := &recordingAnnotationReader{}
	conv := c.createConverter()
	conv.ann = reader
	_ = conv.Sync(true)

	if len(reader.pathLinks) != 1 {
		t.Fatalf("annotation reader got %d path links, want 1", len(reader.pathLinks))
	}
}
