package properties

import (
	"fmt"
)

// This file hosts helper functions to retrieve node-metadata-related properties as described in:
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#wasm-attributes
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#envoy-v3-api-msg-config-core-v3-node
//
// The ProxyConfig variables are described in:
// https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#ProxyConfig

var (
	nodeMetaProxyConfigBinaryPath                     = []string{"node", "metadata", "PROXY_CONFIG", "binaryPath"}
	nodeMetaProxyConfigConcurrency                    = []string{"node", "metadata", "PROXY_CONFIG", "concurrency"}
	nodeMetaProxyConfigConfigPath                     = []string{"node", "metadata", "PROXY_CONFIG", "configPath"}
	nodeProxyConfigControlPlaneAuthPolicy             = []string{"node", "metadata", "PROXY_CONFIG", "controlPlaneAuthPolicy"}
	nodeProxyConfigDiscoveryAddress                   = []string{"node", "metadata", "PROXY_CONFIG", "discoveryAddress"}
	nodeProxyConfigDrainDuration                      = []string{"node", "metadata", "PROXY_CONFIG", "drainDuration"}
	nodeProxyConfigExtraStatTags                      = []string{"node", "metadata", "PROXY_CONFIG", "extraStatTags"}
	nodeProxyConfigHoldApplicationUntilProxyStarts    = []string{"node", "metadata", "PROXY_CONFIG", "holdApplicationUntilProxyStarts"}
	nodeProxyConfigProxyAdminPort                     = []string{"node", "metadata", "PROXY_CONFIG", "proxyAdminPort"}
	nodeProxyConfigProxyStatsMatcherInclusionPrefixes = []string{"node", "metadata", "PROXY_CONFIG", "proxyStatsMatcher", "inclusionPrefixes"}
	nodeProxyConfigProxyStatsMatcherInclusionRegexps  = []string{"node", "metadata", "PROXY_CONFIG", "proxyStatsMatcher", "inclusionRegexps"}
	nodeProxyConfigProxyStatsMatcherInclusionSuffixes = []string{"node", "metadata", "PROXY_CONFIG", "proxyStatsMatcher", "inclusionSuffixes"}
	nodeProxyConfigServiceCluster                     = []string{"node", "metadata", "PROXY_CONFIG", "serviceCluster"}
	nodeProxyConfigStatNameLength                     = []string{"node", "metadata", "PROXY_CONFIG", "statNameLength"}
	nodeProxyConfigStatusPort                         = []string{"node", "metadata", "PROXY_CONFIG", "statusPort"}
	nodeProxyConfigTerminationDrainDuration           = []string{"node", "metadata", "PROXY_CONFIG", "terminationDrainDuration"}
	nodeProxyConfigTracingDatadogAddress              = []string{"node", "metadata", "PROXY_CONFIG", "tracing", "datadog", "address"}
	nodeProxyConfigTracingOpenCensusAgentAddress      = []string{"node", "metadata", "PROXY_CONFIG", "tracing", "opencensusagent", "address"}
	nodeProxyConfigTracingZipkinAddress               = []string{"node", "metadata", "PROXY_CONFIG", "tracing", "zipkin", "address"}
)

// GetNodeMetaProxyConfigBinaryPath returns the path to the proxy binary
func GetNodeMetaProxyConfigBinaryPath() (string, error) {
	return getPropertyString(nodeMetaProxyConfigBinaryPath)
}

// GetNodeMetaProxyConfigConcurrency returns the concurrency configuration of the proxy which
// is the number of worker threads to run. If unset, this will be automatically determined based
// on CPU requests/limits. If set to 0, all cores on the machine will be used. Default is 2 worker
// threads
func GetNodeMetaProxyConfigConcurrency() (float64, error) {
	return getPropertyFloat64(nodeMetaProxyConfigConcurrency)
}

// GetNodeMetaProxyConfigConfigPath returns the path to the proxy configuration, Proxy agent
// generates the actual configuration and stores it in this directory
func GetNodeMetaProxyConfigConfigPath() (string, error) {
	return getPropertyString(nodeMetaProxyConfigConfigPath)
}

// GetNodeProxyConfigControlPlaneAuthPolicy returns the control plane authentication policy of
// the proxy. The authenticationPolicy defines how the proxy is authenticated when it connects
// to the control plane. Default is set to MUTUAL_TLS
func GetNodeProxyConfigControlPlaneAuthPolicy() (string, error) {
	return getPropertyString(nodeProxyConfigControlPlaneAuthPolicy)
}

// GetNodeProxyConfigDiscoveryAddress returns the discovery address of the proxy. The discovery
// service exposes xDS over an mTLS connection. The inject configuration may override this value
func GetNodeProxyConfigDiscoveryAddress() (string, error) {
	return getPropertyString(nodeProxyConfigDiscoveryAddress)
}

// GetNodeProxyConfigDrainDuration returns the drain duration of the proxy, the time in seconds
// that Envoy will drain connections during a hot restart. MUST be >=1s (e.g., 1s/1m/1h).
// Default drain duration is 45s
func GetNodeProxyConfigDrainDuration() (string, error) {
	return getPropertyString(nodeProxyConfigDrainDuration)
}

// GetNodeProxyConfigExtraStatTags returns the extra stat tags of the proxy to extract from the
// in-proxy Istio telemetry. These extra tags can be added by configuring the telemetry extension.
// Each additional tag needs to be present in this list. Extra tags emitted by the telemetry
// extensions must be listed here so that they can be processed and exposed as Prometheus metrics
func GetNodeProxyConfigExtraStatTags() ([]string, error) {
	return getPropertyStringSlice(nodeProxyConfigExtraStatTags)
}

// GetNodeProxyConfigHoldApplicationUntilProxyStarts returns whether to hold the application until
// the proxy starts. A boolean flag for enabling/disabling the holdApplicationUntilProxyStarts
// behavior. This feature adds hooks to delay application startup until the pod proxy is ready to
// accept traffic, mitigating some startup race conditions. Default value is ‘false’
func GetNodeProxyConfigHoldApplicationUntilProxyStarts() (bool, error) {
	return getPropertyBool(nodeProxyConfigHoldApplicationUntilProxyStarts)
}

// GetNodeProxyConfigProxyAdminPort returns the admin port of the proxy for administrative commands.
// Default port is 15000
func GetNodeProxyConfigProxyAdminPort() (float64, error) {
	return getPropertyFloat64(nodeProxyConfigProxyAdminPort)
}

// GetNodeProxyConfigProxyStatsMatcher returns the proxy stats matcher, which defines
// configuration for reporting custom Envoy stats. To reduce memory and CPU overhead
// from Envoy stats system, Istio proxies by default create and expose only a subset of Envoy
// stats. This option is to control creation of additional Envoy stats with prefix, suffix,
// and regex expressions match on the name of the stats. This replaces the stats inclusion
// annotations (sidecar.istio.io/statsInclusionPrefixes, sidecar.istio.io/statsInclusionRegexps,
// and sidecar.istio.io/statsInclusionSuffixes)
func GetNodeProxyConfigProxyStatsMatcher() (IstioProxyStatsMatcher, error) {
	var matcher IstioProxyStatsMatcher
	var errorsCount int

	prefixes, err := getPropertyStringSlice(nodeProxyConfigProxyStatsMatcherInclusionPrefixes)
	if err != nil {
		errorsCount++
	} else {
		matcher.InclusionPrefixes = prefixes
	}

	regexps, err := getPropertyStringSlice(nodeProxyConfigProxyStatsMatcherInclusionRegexps)
	if err != nil {
		errorsCount++
	} else {
		matcher.InclusionRegexps = regexps
	}

	suffixes, err := getPropertyStringSlice(nodeProxyConfigProxyStatsMatcherInclusionSuffixes)
	if err != nil {
		errorsCount++
	} else {
		matcher.InclusionSuffixes = suffixes
	}

	if errorsCount == 3 {
		return IstioProxyStatsMatcher{}, fmt.Errorf("failed to fetch any components of IstioProxyStatsMatcher")
	}

	return matcher, nil
}

// GetNodeProxyConfigServiceCluster returns the name of the service cluster of the proxy
// that is shared by all Envoy instances. This setting corresponds to --service-cluster flag
// in Envoy. In a typical Envoy deployment, the service-cluster flag is used to identify the
// caller, for source-based routing scenarios. Since Istio does not assign a local service
// version to each Envoy instance, the name is same for all of them. However, the source/caller’s
// identity (e.g., IP address) is encoded in the --service-node flag when launching Envoy.
// When the RDS service receives API calls from Envoy, it uses the value of the service-node
// flag to compute routes that are relative to the service instances located at that IP address
func GetNodeProxyConfigServiceCluster() (string, error) {
	return getPropertyString(nodeProxyConfigServiceCluster)
}

// GetNodeProxyConfigStatNameLength returns the stat name length of the proxy, The length
// of the name field is determined by the length of a name field in a service and the set
// of labels that comprise a particular version of the service. The default value is set
// to 189 characters. Envoy’s internal metrics take up 67 characters, for a total of 256
// character name per metric. Increase the value of this field if you find that the metrics
// from Envoys are truncated
func GetNodeProxyConfigStatNameLength() (float64, error) {
	return getPropertyFloat64(nodeProxyConfigStatNameLength)
}

// GetNodeProxyConfigStatusPort returns the port on which the agent should listen
// for administrative commands such as readiness probe. Default is set to port 15020
func GetNodeProxyConfigStatusPort() (float64, error) {
	return getPropertyFloat64(nodeProxyConfigStatusPort)
}

// GetNodeProxyConfigTerminationDrainDuration returns the stat name length of the proxy,
// the amount of time allowed for connections to complete on proxy shutdown. On receiving
// SIGTERM or SIGINT, istio-agent tells the active Envoy to start draining, preventing
// any new connections and allowing existing connections to complete. It then sleeps for
// the termination_drain_duration and then kills any remaining active Envoy processes.
// If not set, a default of 5s will be applied
func GetNodeProxyConfigTerminationDrainDuration() (string, error) {
	return getPropertyString(nodeProxyConfigTerminationDrainDuration)
}

// GetNodeProxyConfigTracingDatadogAddress returns the address of the Datadog
// service (e.g. datadog-agent.sre.svc.cluster.local:8126)
func GetNodeProxyConfigTracingDatadogAddress() (string, error) {
	return getPropertyString(nodeProxyConfigTracingDatadogAddress)
}

// GetNodeProxyConfigTracingOpenCensusAgentAddress returns the gRPC address for
// the OpenCensus agent (e.g. dns://authority/host:port or unix:path)
func GetNodeProxyConfigTracingOpenCensusAgentAddress() (string, error) {
	return getPropertyString(nodeProxyConfigTracingOpenCensusAgentAddress)
}

// GetNodeProxyConfigTracingZipkinAddress returns address of the Zipkin service
// (e.g. zipkin.sre.svc.cluster.local:9411)
func GetNodeProxyConfigTracingZipkinAddress() (string, error) {
	return getPropertyString(nodeProxyConfigTracingZipkinAddress)
}
