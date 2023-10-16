package properties

import (
	"fmt"
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// This file hosts helper functions to retrieve wasm-related properties as described in:
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#wasm-attributes

var (
	pluginName                = []string{"plugin_name"}
	pluginRootId              = []string{"plugin_root_id"}
	pluginVmId                = []string{"plugin_vm_id"}
	clusterName               = []string{"cluster_name"}
	routeName                 = []string{"route_name"}
	listenerDirection         = []string{"listener_direction"}
	nodeId                    = []string{"node", "id"}
	nodeCluster               = []string{"node", "cluster"}
	nodeDynamicParams         = []string{"node", "dynamic_parameters", "params"}
	nodeLocalityRegion        = []string{"node", "locality", "region"}
	nodeLocalityZone          = []string{"node", "locality", "zone"}
	nodeLocalitySubzone       = []string{"node", "locality", "subzone"}
	nodeUserAgentName         = []string{"node", "user_agent_name"}
	nodeUserAgentVersion      = []string{"node", "user_agent_version"}
	nodeUserAgentBuildVersion = []string{"node", "user_agent_build_version", "metadata"}
	nodeExtensions            = []string{"node", "extensions"}
	nodeClientFeatures        = []string{"node", "client_features"}
	nodeListeningAddresses    = []string{"node", "listening_addresses"}
	clusterMetadata           = []string{"node", "cluster_metadata", "filter_metadata", "istio"}
	listenerMetadata          = []string{"node", "listener_metadata", "filter_metadata", "istio"}
	routeMetadata             = []string{"node", "route_metadata", "filter_metadata", "istio"}
	upstreamHostMetadata      = []string{"node", "upstream_host_metadata", "filter_metadata", "istio"}
)

// GetPluginName returns the plugin name.
//
// This matches <metadata.name>.<metadata.namespace> in an istio WasmPlugin CR.
func GetPluginName() (string, error) {
	return getPropertyString(pluginName)
}

// GetPluginRootId returns the plugin root id.
//
// This matches the <spec.pluginName> in the istio WasmPlugin CR.
func GetPluginRootId() (string, error) {
	return getPropertyString(pluginRootId)
}

// GetPluginVmId returns the plugin vm id.
func GetPluginVmId() (string, error) {
	return getPropertyString(pluginVmId)
}

// GetClusterName returns the upstream cluster name.
//
// Example value: "outbound|80||httpbin.org".
func GetClusterName() (string, error) {
	return getPropertyString(clusterName)
}

// GetRouteName returns the route name, only available in the response path (cfr getXdsRouteName()).
//
// This matches the <spec.http.name> in the istio VirtualService CR.
func GetRouteName() (string, error) {
	return getPropertyString(routeName)
}

// GetListenerDirection returns the listener direction.
//
// Possible values are:
//
//   - UNSPECIFIED: 0 (default option is unspecified)
//   - INBOUND: 1 (⁣the transport is used for incoming traffic)
//   - OUTBOUND: 2 (the transport is used for outgoing traffic)
func GetListenerDirection() (EnvoyTrafficDirection, error) {
	result, err := getPropertyUint64(listenerDirection)
	if err != nil {
		return EnvoyTrafficDirection(Unspecified), err
	}
	return EnvoyTrafficDirection(int(result)), nil
}

// GetNodeId returns the node id, an opaque node identifier for the Envoy node. This also
// provides the local service node name. It should be set if any of the following features
// are used: statsd, CDS, and HTTP tracing, either in this message or via --service-node
//
// Example value:
// router~10.244.0.22~istio-ingress-6d78c67d85-qsbtz.istio-ingress~istio-ingress.svc.cluster.local
func GetNodeId() (string, error) {
	return getPropertyString(nodeId)
}

// GetNodeCluster returns the node cluster, which defines the local service cluster
// name where envoy is running. Though optional, it should be set if any of the following
// features are used: statsd, health check cluster verification, runtime override directory,
// user agent addition, HTTP global rate limiting, CDS, and HTTP tracing, either in this
// message or via --service-cluster
//
// Example value: istio-ingress.istio-ingress
func GetNodeCluster() (string, error) {
	return getPropertyString(nodeCluster)
}

// GetNodeDynamicParams returns the node dynamic parameters. These may vary at
// runtime (unlike other fields in this message). For example, the xDS client may have a
// shared identifier that changes during the lifetime of the xDS client. In Envoy, this
// would be achieved by updating the dynamic context on the Server::Instance’s LocalInfo
// context provider. The shard ID dynamic parameter then appears in this field during
// future discovery requests
func GetNodeDynamicParams() (string, error) {
	return getPropertyString(nodeDynamicParams)
}

// GetNodeLocality returns the node locality.
func GetNodeLocality() (EnvoyLocality, error) {
	result := EnvoyLocality{}
	var errors []string
	var successCount int

	region, err := getPropertyString(nodeLocalityRegion)
	if err != nil {
		errors = append(errors, err.Error())
	} else {
		result.Region = region
		successCount++
	}

	zone, err := getPropertyString(nodeLocalityZone)
	if err != nil {
		errors = append(errors, err.Error())
	} else {
		result.Zone = zone
		successCount++
	}

	subzone, err := getPropertyString(nodeLocalitySubzone)
	if err != nil {
		errors = append(errors, err.Error())
	} else {
		result.Subzone = subzone
		successCount++
	}

	if successCount == 0 {
		return result, fmt.Errorf(strings.Join(errors, "; "))
	}

	return result, nil
}

// GetNodeUserAgentName returns the node user agent name.
//
// Example: “envoy” or “grpc”.
func GetNodeUserAgentName() (string, error) {
	return getPropertyString(nodeUserAgentName)
}

// GetNodeUserAgentVersion returns the node user agent version.
//
// Example “1.12.2” or “abcd1234”, or “SpecialEnvoyBuild”.
func GetNodeUserAgentVersion() (string, error) {
	return getPropertyString(nodeUserAgentVersion)
}

// GetNodeUserAgentBuildVersion returns the node user agent build version.
func GetNodeUserAgentBuildVersion() (string, error) {
	return getPropertyString(nodeUserAgentBuildVersion)
}

// GetNodeExtensions returns the node extensions.
func GetNodeExtensions() ([]EnvoyExtension, error) {
	result := make([]EnvoyExtension, 0)
	extensionsRawSlice, err := getPropertyByteSliceSlice(nodeExtensions)
	if err != nil {
		return []EnvoyExtension{}, err
	}

	for _, extensionRawSlice := range extensionsRawSlice {
		if extensionRawSlice == nil {
			continue
		}
		extensionStringSlice := deserializeProtoStringSlice(extensionRawSlice)
		extension := EnvoyExtension{}

		if len(extensionStringSlice) > 0 {
			extension.Name = string(extensionStringSlice[0])
		}
		if len(extensionStringSlice) > 1 {
			extension.Category = string(extensionStringSlice[1])
		}
		if len(extensionStringSlice) > 2 {
			extenstionTypeUrls := []string{}
			extenstionTypeUrls = append(extenstionTypeUrls, extensionStringSlice[2:]...)
			extension.TypeUrls = extenstionTypeUrls
		}

		result = append(result, extension)
	}

	return result, nil
}

// GetNodeClientFeatures returns the node client features. These are well known features
// described in the Envoy API repository for a given major version of an API. Client
// features use reverse DNS naming scheme, for example "com.acme.feature".
func GetNodeClientFeatures() ([]string, error) {
	result, err := proxywasm.GetProperty(nodeClientFeatures)
	if err != nil {
		return []string{}, err
	}
	return deserializeProtoStringSlice(result), nil
}

// GetNodeListeningAddresses returns the node listening addresses.
func GetNodeListeningAddresses() ([]string, error) {
	return getPropertyStringSlice(nodeListeningAddresses)
}

// GetClusterMetadata returns the cluster metadata.
func GetClusterMetadata() (IstioFilterMetadata, error) {
	return getIstioFilterMetadata(clusterMetadata)
}

// GetListenerMetadata returns the listener metadata.
func GetListenerMetadata() (IstioFilterMetadata, error) {
	return getIstioFilterMetadata(listenerMetadata)
}

// GetRouteMetadata returns the route metadata.
func GetRouteMetadata() (IstioFilterMetadata, error) {
	return getIstioFilterMetadata(routeMetadata)
}

// GetUpstreamHostMetadata returns the upstream host metadata.
func GetUpstreamHostMetadata() (IstioFilterMetadata, error) {
	return getIstioFilterMetadata(upstreamHostMetadata)
}
