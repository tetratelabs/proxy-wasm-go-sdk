package properties

// This file hosts helper functions to retrieve xsd-configuration-related properties as described in:
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#configuration-attributes

var (
	xdsClusterName             = []string{"xds", "cluster_name"}
	xdsClusterMetadata         = []string{"xds", "cluster_metadata", "filter_metadata", "istio"}
	xdsRouteName               = []string{"xds", "route_name"}
	xdsRouteMetadata           = []string{"xds", "route_metadata", "filter_metadata", "istio"}
	xdsUpstreamHostMetadata    = []string{"xds", "upstream_host_metadata", "filter_metadata", "istio"}
	xdsListenerFilterChainName = []string{"xds", "filter_chain_name"}
)

// GetXdsClusterName returns the upstream cluster name.
//
// Example value: "outbound|80||httpbin.org".
func GetXdsClusterName() (string, error) {
	return getPropertyString(xdsClusterName)
}

// GetXdsClusterMetadata returns the upstream cluster metadata.
func GetXdsClusterMetadata() (IstioFilterMetadata, error) {
	return getIstioFilterMetadata(xdsClusterMetadata)
}

// GetXdsRouteName returns the upstream route name (available in both
// the request response path, cfr getRouteName()). This matches the
// <spec.http.name> in an istio VirtualService CR.
func GetXdsRouteName() (string, error) {
	return getPropertyString(xdsRouteName)
}

// GetXdsRouteMetadata returns the upstream route metadata.
func GetXdsRouteMetadata() (IstioFilterMetadata, error) {
	return getIstioFilterMetadata(xdsRouteMetadata)
}

// GetXdsUpstreamHostMetadata returns the upstream host metadata.
func GetXdsUpstreamHostMetadata() (IstioFilterMetadata, error) {
	return getIstioFilterMetadata(xdsUpstreamHostMetadata)
}

// GetXdsListenerFilterChainName returns the listener filter chain name.
func GetXdsListenerFilterChainName() (string, error) {
	return getPropertyString(xdsListenerFilterChainName)
}
