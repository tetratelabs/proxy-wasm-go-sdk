package properties

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
)

func TestGetNodeMetaAnnotations(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]string
		expect map[string]string
	}{
		{
			name: "Non-Empty Annotations",
			input: map[string]string{
				"inject.istio.io/templates":   "gateway",
				"istio.io/rev":                "default",
				"kubernetes.io/config.seen":   "2023-10-13T10:39:01.174733724Z",
				"kubernetes.io/config.source": "api",
				"prometheus.io/path":          "/stats/prometheus",
				"prometheus.io/port":          "15020",
				"prometheus.io/scrape":        "true",
				"proxy.istio.io/overrides":    `{"containers":[{"name":"istio-proxy","ports":[{"name":"http-envoy-prom","containerPort":15090,"protocol":"TCP"}],"resources":{"limits":{"cpu":"2","memory":"1Gi"},"requests":{"cpu":"100m","memory":"128Mi"}},"volumeMounts":[{"name":"kube-api-access-6mm2z","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"Always","securityContext":{"capabilities":{"drop":["ALL"]},"privileged":false,"runAsUser":1337,"runAsGroup":1337,"runAsNonRoot":true,"readOnlyRootFilesystem":true,"allowPrivilegeEscalation":false}}]} sidecar.istio.io/componentLogLevel:wasm:debug sidecar.istio.io/inject:true sidecar.istio.io/status:{"initContainers":null,"containers":["istio-proxy"],"volumes":["workload-socket","credential-socket","workload-certs","istio-envoy","istio-data","istio-podinfo","istio-token","istiod-ca-cert"],"imagePullSecrets":null,"revision":"default"}]`,
			},
			expect: map[string]string{
				"inject.istio.io/templates":   "gateway",
				"istio.io/rev":                "default",
				"kubernetes.io/config.seen":   "2023-10-13T10:39:01.174733724Z",
				"kubernetes.io/config.source": "api",
				"prometheus.io/path":          "/stats/prometheus",
				"prometheus.io/port":          "15020",
				"prometheus.io/scrape":        "true",
				"proxy.istio.io/overrides":    `{"containers":[{"name":"istio-proxy","ports":[{"name":"http-envoy-prom","containerPort":15090,"protocol":"TCP"}],"resources":{"limits":{"cpu":"2","memory":"1Gi"},"requests":{"cpu":"100m","memory":"128Mi"}},"volumeMounts":[{"name":"kube-api-access-6mm2z","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"Always","securityContext":{"capabilities":{"drop":["ALL"]},"privileged":false,"runAsUser":1337,"runAsGroup":1337,"runAsNonRoot":true,"readOnlyRootFilesystem":true,"allowPrivilegeEscalation":false}}]} sidecar.istio.io/componentLogLevel:wasm:debug sidecar.istio.io/inject:true sidecar.istio.io/status:{"initContainers":null,"containers":["istio-proxy"],"volumes":["workload-socket","credential-socket","workload-certs","istio-envoy","istio-data","istio-podinfo","istio-token","istiod-ca-cert"],"imagePullSecrets":null,"revision":"default"}]`,
			},
		},
		{
			name:   "Empty Annotations",
			input:  map[string]string{},
			expect: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaAnnotations, serializeStringMap(tt.input))
			_, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			result, err := GetNodeMetaAnnotations()
			require.NoError(t, err)
			require.Equal(t, tt.expect, result)
		})
	}
}

func TestGetNodeMetaAppContainers(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaAppContainers, []byte("metadata"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaAppContainers()
	require.NoError(t, err)
	require.Equal(t, "metadata", result)
}

func TestGetNodeMetaClusterId(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaClusterId, []byte("Kubernetes"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaClusterId()
	require.NoError(t, err)
	require.Equal(t, "Kubernetes", result)
}

func TestGetNodeMetaEnvoyPrometheusPort(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaEnvoyPrometheusPort, serializeFloat64(15090))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaEnvoyPrometheusPort()
	require.NoError(t, err)
	require.Equal(t, float64(15090), result)
}

func TestGetNodeMetaEnvoyStatusPort(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaEnvoyStatusPort, serializeFloat64(15021))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaEnvoyStatusPort()
	require.NoError(t, err)
	require.Equal(t, float64(15021), result)
}

func TestGetNodeMetaInstanceIps(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaInstanceIps, []byte("10.244.0.13"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaInstanceIps()
	require.NoError(t, err)
	require.Equal(t, "10.244.0.13", result)
}

func TestGetNodeMetaInterceptionMode(t *testing.T) {
	tests := []struct {
		name           string
		propertyValue  []byte
		expectedResult IstioTrafficInterceptionMode
		expectedError  error
	}{
		{
			name:           "Valid TPROXY mode",
			propertyValue:  []byte("TPROXY"),
			expectedResult: Tproxy,
			expectedError:  nil,
		},
		{
			name:           "Valid REDIRECT mode",
			propertyValue:  []byte("REDIRECT"),
			expectedResult: Redirect,
			expectedError:  nil,
		},
		{
			name:           "Valid NONE mode",
			propertyValue:  []byte("NONE"),
			expectedResult: None,
			expectedError:  nil,
		},
		{
			name:           "Invalid mode",
			propertyValue:  []byte("INVALID_MODE"),
			expectedResult: Redirect,
			expectedError:  fmt.Errorf("invalid IstioTrafficInterceptionMode: INVALID_MODE"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaInterceptionMode, test.propertyValue)
			_, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			result, err := GetNodeMetaInterceptionMode()
			if test.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, test.expectedError, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, test.expectedResult, result)
		})
	}
}

func TestGetNodeMetaIstioProxySha(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaIstioProxySha, []byte("3c27a1b0cf381ca854ccc3a2034e88c206928da2"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaIstioProxySha()
	require.NoError(t, err)
	require.Equal(t, "3c27a1b0cf381ca854ccc3a2034e88c206928da2", result)
}

func TestGetNodeMetaIstioVersion(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaIstioVersion, []byte("1.18.2-tetrate-v0"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaIstioVersion()
	require.NoError(t, err)
	require.Equal(t, "1.18.2-tetrate-v0", result)
}

func TestGetNodeMetaLabels(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]string
		expect map[string]string
	}{
		{
			name: "Non-Empty Labels",
			input: map[string]string{
				"app":                                 "istio-ingress",
				"istio":                               "ingress",
				"istio-locality":                      "region1",
				"service.istio.io/canonical-name":     "istio-ingress",
				"service.istio.io/canonical-revision": "latest",
				"sidecar.istio.io/inject":             "true",
			},
			expect: map[string]string{
				"app":                                 "istio-ingress",
				"istio":                               "ingress",
				"istio-locality":                      "region1",
				"service.istio.io/canonical-name":     "istio-ingress",
				"service.istio.io/canonical-revision": "latest",
				"sidecar.istio.io/inject":             "true",
			},
		},
		{
			name:   "Empty Labels",
			input:  map[string]string{},
			expect: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaLabels, serializeStringMap(tt.input))
			_, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			result, err := GetNodeMetaLabels()
			require.NoError(t, err)
			require.Equal(t, tt.expect, result)
		})
	}
}

func TestGetNodeMetaMeshId(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaMeshId, []byte("cluster.local"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaMeshId()
	require.NoError(t, err)
	require.Equal(t, "cluster.local", result)
}

func TestGetNodeMetaName(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaName, []byte("istio-ingress-67cddc6d57-kk2cr"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaName()
	require.NoError(t, err)
	require.Equal(t, "istio-ingress-67cddc6d57-kk2cr", result)
}

func TestGetNodeMetaNamespace(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaNamespace, []byte("istio-ingress"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaNamespace()
	require.NoError(t, err)
	require.Equal(t, "istio-ingress", result)
}

func TestGetNodeMetaNodeName(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaNodeName, []byte("istio-wasm-control-plane"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaNodeName()
	require.NoError(t, err)
	require.Equal(t, "istio-wasm-control-plane", result)
}

func TestGetNodeMetaOwner(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaOwner, []byte("kubernetes://apis/apps/v1/namespaces/istio-ingress/deployments/istio-ingress"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaOwner()
	require.NoError(t, err)
	require.Equal(t, "kubernetes://apis/apps/v1/namespaces/istio-ingress/deployments/istio-ingress", result)
}

func TestGetNodeMetaPilotSan(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		expect []string
	}{
		{
			name:   "Single SAN",
			input:  []string{"istiod.istio-system.svc"},
			expect: []string{"istiod.istio-system.svc"},
		},
		{
			name:   "Multiple SANs",
			input:  []string{"istiod.istio-system.svc", "istiod.istio-system.svc.cluster.local"},
			expect: []string{"istiod.istio-system.svc", "istiod.istio-system.svc.cluster.local"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaPilotSan, serializeStringSlice(tt.input))
			_, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			result, err := GetNodeMetaPilotSan()
			require.NoError(t, err)
			require.Equal(t, tt.expect, result)
		})
	}
}

func TestGetNodeMetaPodPorts(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaPodPorts, []byte(`[{"name":"http-envoy-prom","containerPort":15090,"protocol":"TCP"}]`))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaPodPorts()
	require.NoError(t, err)
	require.Equal(t, `[{"name":"http-envoy-prom","containerPort":15090,"protocol":"TCP"}]`, result)
}

func TestGetNodeMetaServiceAccount(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaServiceAccount, []byte("istio-ingress"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaServiceAccount()
	require.NoError(t, err)
	require.Equal(t, "istio-ingress", result)
}

func TestGetNodeMetaWorkloadName(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaWorkloadName, []byte("istio-ingress"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaWorkloadName()
	require.NoError(t, err)
	require.Equal(t, "istio-ingress", result)
}
