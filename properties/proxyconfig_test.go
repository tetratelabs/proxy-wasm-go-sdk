package properties

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
)

func TestGetNodeMetaProxyConfigBinaryPath(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaProxyConfigBinaryPath, []byte("/usr/local/bin/envoy"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaProxyConfigBinaryPath()
	require.NoError(t, err)
	require.Equal(t, "/usr/local/bin/envoy", result)
}

func TestGetNodeMetaProxyConfigConcurrency(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaProxyConfigConcurrency, serializeFloat64(4))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaProxyConfigConcurrency()
	require.NoError(t, err)
	require.Equal(t, float64(4), result)
}

func TestGetNodeMetaProxyConfigConfigPath(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeMetaProxyConfigConfigPath, []byte("./etc/istio/proxy"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeMetaProxyConfigConfigPath()
	require.NoError(t, err)
	require.Equal(t, "./etc/istio/proxy", result)
}

func TestGetNodeProxyConfigControlPlaneAuthPolicy(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeProxyConfigControlPlaneAuthPolicy, []byte("MUTUAL_TLS"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeProxyConfigControlPlaneAuthPolicy()
	require.NoError(t, err)
	require.Equal(t, "MUTUAL_TLS", result)
}

func TestGetNodeProxyConfigDiscoveryAddress(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeProxyConfigDiscoveryAddress, []byte("istiod.istio-system.svc:15012"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeProxyConfigDiscoveryAddress()
	require.NoError(t, err)
	require.Equal(t, "istiod.istio-system.svc:15012", result)
}

func TestGetNodeProxyConfigDrainDuration(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeProxyConfigDrainDuration, []byte("45s"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeProxyConfigDrainDuration()
	require.NoError(t, err)
	require.Equal(t, "45s", result)
}

func TestGetNodeProxyConfigExtraStatTags(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeProxyConfigExtraStatTags, serializeStringSlice([]string{"tag1", "tag2", "tag3"}))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeProxyConfigExtraStatTags()
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"tag1", "tag2", "tag3"}, result)
}

func TestGetNodeProxyConfigHoldApplicationUntilProxyStarts(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected bool
	}{
		{
			name:     "Test with HoldApplicationUntilProxyStarts true",
			input:    true,
			expected: true,
		},
		{
			name:     "Test with HoldApplicationUntilProxyStarts false",
			input:    false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := proxytest.NewEmulatorOption().WithProperty(nodeProxyConfigHoldApplicationUntilProxyStarts, serializeBool(tt.input))
			_, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			result, err := GetNodeProxyConfigHoldApplicationUntilProxyStarts()
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestGetNodeProxyConfigProxyAdminPort(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeProxyConfigProxyAdminPort, serializeFloat64(15000))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeProxyConfigProxyAdminPort()
	require.NoError(t, err)
	require.Equal(t, float64(15000), result)
}

func TestGetNodeProxyConfigProxyStatsMatcher(t *testing.T) {
	tests := []struct {
		name             string
		emulatorOption   *proxytest.EmulatorOption
		expectedMatcher  IstioProxyStatsMatcher
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name: "full matcher",
			emulatorOption: proxytest.NewEmulatorOption().
				WithProperty(nodeProxyConfigProxyStatsMatcherInclusionPrefixes, serializeStringSlice([]string{"prefix1", "prefix2"})).
				WithProperty(nodeProxyConfigProxyStatsMatcherInclusionRegexps, serializeStringSlice([]string{"regexp1", "regexp2"})).
				WithProperty(nodeProxyConfigProxyStatsMatcherInclusionSuffixes, serializeStringSlice([]string{"suffix1", "suffix2"})),
			expectedMatcher: IstioProxyStatsMatcher{
				InclusionPrefixes: []string{"prefix1", "prefix2"},
				InclusionRegexps:  []string{"regexp1", "regexp2"},
				InclusionSuffixes: []string{"suffix1", "suffix2"},
			},
			expectError: false,
		},
		{
			name: "only prefixes",
			emulatorOption: proxytest.NewEmulatorOption().
				WithProperty(nodeProxyConfigProxyStatsMatcherInclusionPrefixes, serializeStringSlice([]string{"prefix1", "prefix2"})),
			expectedMatcher: IstioProxyStatsMatcher{
				InclusionPrefixes: []string{"prefix1", "prefix2"},
			},
			expectError: false,
		},
		{
			name: "only regexps",
			emulatorOption: proxytest.NewEmulatorOption().
				WithProperty(nodeProxyConfigProxyStatsMatcherInclusionRegexps, serializeStringSlice([]string{"regexp1", "regexp2"})),
			expectedMatcher: IstioProxyStatsMatcher{
				InclusionRegexps: []string{"regexp1", "regexp2"},
			},
			expectError: false,
		},
		{
			name: "only suffixes",
			emulatorOption: proxytest.NewEmulatorOption().
				WithProperty(nodeProxyConfigProxyStatsMatcherInclusionSuffixes, serializeStringSlice([]string{"suffix1", "suffix2"})),
			expectedMatcher: IstioProxyStatsMatcher{
				InclusionSuffixes: []string{"suffix1", "suffix2"},
			},
			expectError: false,
		},
		{
			name:             "no properties",
			emulatorOption:   proxytest.NewEmulatorOption(),
			expectedMatcher:  IstioProxyStatsMatcher{},
			expectError:      true,
			expectedErrorMsg: "failed to fetch any components of IstioProxyStatsMatcher",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, reset := proxytest.NewHostEmulator(tt.emulatorOption)
			defer reset()

			matcher, err := GetNodeProxyConfigProxyStatsMatcher()

			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, tt.expectedErrorMsg, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedMatcher, matcher)
			}
		})
	}
}

func TestGetNodeProxyConfigServiceCluster(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeProxyConfigServiceCluster, []byte("service-cluster-name"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeProxyConfigServiceCluster()
	require.NoError(t, err)
	require.Equal(t, "service-cluster-name", result)
}

func TestGetNodeProxyConfigStatNameLength(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeProxyConfigStatNameLength, serializeFloat64(256))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeProxyConfigStatNameLength()
	require.NoError(t, err)
	require.Equal(t, float64(256), result)
}

func TestGetNodeProxyConfigStatusPort(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeProxyConfigStatusPort, serializeFloat64(15020))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeProxyConfigStatusPort()
	require.NoError(t, err)
	require.Equal(t, float64(15020), result)
}

func TestGetNodeProxyConfigTerminationDrainDuration(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeProxyConfigTerminationDrainDuration, []byte("5s"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeProxyConfigTerminationDrainDuration()
	require.NoError(t, err)
	require.Equal(t, "5s", result)
}

func TestGetNodeProxyConfigTracingDatadogAddress(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeProxyConfigTracingDatadogAddress, []byte("datadog-agent.sre.svc.cluster.local:8126"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeProxyConfigTracingDatadogAddress()
	require.NoError(t, err)
	require.Equal(t, "datadog-agent.sre.svc.cluster.local:8126", result)
}

func TestGetNodeProxyConfigTracingOpenCensusAgentAddress(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeProxyConfigTracingOpenCensusAgentAddress, []byte("opencensus-agent.sre.svc.cluster.local:55678"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeProxyConfigTracingOpenCensusAgentAddress()
	require.NoError(t, err)
	require.Equal(t, "opencensus-agent.sre.svc.cluster.local:55678", result)
}

func TestGetNodeProxyConfigTracingZipkinAddress(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(nodeProxyConfigTracingZipkinAddress, []byte("zipkin.sre.svc.cluster.local:9411"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetNodeProxyConfigTracingZipkinAddress()
	require.NoError(t, err)
	require.Equal(t, "zipkin.sre.svc.cluster.local:9411", result)
}
