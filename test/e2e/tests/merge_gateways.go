// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

//go:build e2e
// +build e2e

package tests

import (
	"testing"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/gateway-api/conformance/utils/http"
	"sigs.k8s.io/gateway-api/conformance/utils/kubernetes"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"
)

func init() {
	ConformanceTests = append(ConformanceTests, MergeGatewaysTest)
}

var MergeGatewaysTest = suite.ConformanceTest{
	ShortName:   "MergeGateways",
	Description: "Merge gateways on to a single EnvoyProxy",
	Manifests:   []string{"testdata/merge-gateways.yaml"},
	Test: func(t *testing.T, suite *suite.ConformanceTestSuite) {
		ns1 := "gateway-conformance-infra"
		ns2 := "default"

		route1NN := types.NamespacedName{Name: "merged-gateway-route-1", Namespace: ns1}
		gw1NN := types.NamespacedName{Name: "merged-gateway-1", Namespace: ns1}
		gw1Addr := kubernetes.GatewayAndHTTPRoutesMustBeAccepted(t, suite.Client, suite.TimeoutConfig, suite.ControllerName, kubernetes.NewGatewayRef(gw1NN), route1NN)

		route2NN := types.NamespacedName{Name: "merged-gateway-route-2", Namespace: ns1}
		gw2NN := types.NamespacedName{Name: "merged-gateway-2", Namespace: ns1}
		gw2Addr := kubernetes.GatewayAndHTTPRoutesMustBeAccepted(t, suite.Client, suite.TimeoutConfig, suite.ControllerName, kubernetes.NewGatewayRef(gw2NN), route2NN)

		route3NN := types.NamespacedName{Name: "merged-gateway-route-3", Namespace: ns1}
		gw3NN := types.NamespacedName{Name: "merged-gateway-3", Namespace: ns2}
		gw3Addr := kubernetes.GatewayAndHTTPRoutesMustBeAccepted(t, suite.Client, suite.TimeoutConfig, suite.ControllerName, kubernetes.NewGatewayRef(gw3NN), route3NN)

		if gw1Addr != gw2Addr {
			t.Errorf("fail to merge gateways: inconsistent gateway address %s and %s", gw1Addr, gw2Addr)
		}
		if gw2Addr != gw3Addr {
			t.Errorf("fail to merge gateways: inconsistent gateway address %s and %s", gw2Addr, gw3Addr)
		}

		t.Run("merged three gateways with http routes", func(t *testing.T) {
			// Three gateways have the same address.
			gwAddr := gw1Addr

			http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, http.ExpectedResponse{
				Request:   http.Request{Path: "/merge1"},
				Response:  http.Response{StatusCode: 200},
				Namespace: ns1,
				Backend:   "infra-backend-v1",
			})

			http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, http.ExpectedResponse{
				Request:   http.Request{Path: "/merge2"},
				Response:  http.Response{StatusCode: 200},
				Namespace: ns1,
				Backend:   "infra-backend-v2",
			})

			http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, http.ExpectedResponse{
				Request:   http.Request{Path: "/merge3"},
				Response:  http.Response{StatusCode: 200},
				Namespace: ns2,
				Backend:   "infra-backend-v3",
			})
		})
	},
}
