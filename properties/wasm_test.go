package properties

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
)

func TestGetPluginName(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(pluginName, []byte("istio-ingress.print-properties"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetPluginName()
	require.NoError(t, err)
	require.Equal(t, "istio-ingress.print-properties", result)
}

func TestGetPluginRootId(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(pluginRootId, []byte("print-properties"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetPluginRootId()
	require.NoError(t, err)
	require.Equal(t, "print-properties", result)
}

func TestGetPluginVmId(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(pluginVmId, []byte("plugin-vm-id-value"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetPluginVmId()
	require.NoError(t, err)
	require.Equal(t, "plugin-vm-id-value", result)
}

func TestGetClusterName(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(clusterName, []byte("outbound|80||httpbin.org"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetClusterName()
	require.NoError(t, err)
	require.Equal(t, "outbound|80||httpbin.org", result)
}

func TestGetRouteName(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(routeName, []byte("route-name-value"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetRouteName()
	require.NoError(t, err)
	require.Equal(t, "route-name-value", result)
}

func TestGetListenerDirection(t *testing.T) {
	tests := []struct {
		name           string
		input          uint64
		expectedResult EnvoyTrafficDirection
	}{
		{
			name:           "Unspecified",
			input:          0,
			expectedResult: Unspecified,
		},
		{
			name:           "Inbound",
			input:          1,
			expectedResult: Inbound,
		},
		{
			name:           "Outound",
			input:          2,
			expectedResult: Outbound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := proxytest.NewEmulatorOption().WithProperty(listenerDirection, serializeUint64(tt.input))
			_, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			result, err := GetListenerDirection()
			require.NoError(t, err)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGetNodeId(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeId, []byte("router~10.244.0.22~istio-ingress-6d78c67d85-qsbtz.istio-ingress~istio-ingress.svc.cluster.local"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeId()
	require.NoError(t, err)
	require.Equal(t, "router~10.244.0.22~istio-ingress-6d78c67d85-qsbtz.istio-ingress~istio-ingress.svc.cluster.local", result)
}

func TestGetNodeCluster(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeCluster, []byte("istio-ingress.istio-ingress"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeCluster()
	require.NoError(t, err)
	require.Equal(t, "istio-ingress.istio-ingress", result)
}

func TestGetNodeDynamicParams(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeDynamicParams, []byte("dynamic-params-value"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeDynamicParams()
	require.NoError(t, err)
	require.Equal(t, "dynamic-params-value", result)
}

func TestGetNodeLocality(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithProperty(nodeLocalityRegion, []byte("region-value")).
		WithProperty(nodeLocalityZone, []byte("zone-value")).
		WithProperty(nodeLocalitySubzone, []byte("subzone-value"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeLocality()
	require.NoError(t, err)
	require.Equal(t, "region-value", result.Region)
	require.Equal(t, "zone-value", result.Zone)
	require.Equal(t, "subzone-value", result.Subzone)
}

func TestGetNodeUserAgentName(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeUserAgentName, []byte("envoy"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeUserAgentName()
	require.NoError(t, err)
	require.Equal(t, "envoy", result)
}

func TestGetNodeUserAgentVersion(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeUserAgentVersion, []byte("1.12.2"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeUserAgentVersion()
	require.NoError(t, err)
	require.Equal(t, "1.12.2", result)
}

func TestGetNodeUserAgentBuildVersion(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeUserAgentBuildVersion, []byte("build-version-value"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeUserAgentBuildVersion()
	require.NoError(t, err)
	require.Equal(t, "build-version-value", result)
}

func TestGetNodeClientFeatures(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeClientFeatures, serializeProtoStringSlice([]string{"feature1-data", "feature2-data"}))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeClientFeatures()
	require.NoError(t, err)
	require.Equal(t, []string{"feature1-data", "feature2-data"}, result)
}

func TestGetNodeListeningAddresses(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeListeningAddresses, serializeStringSlice([]string{"192.168.0.10", "10.0.0.20"}))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeListeningAddresses()
	require.NoError(t, err)
	require.Equal(t, []string{"192.168.0.10", "10.0.0.20"}, result)
}
