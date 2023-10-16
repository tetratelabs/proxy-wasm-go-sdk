package properties

// This file hosts helper functions to retrieve node-metadata-related properties as described in:
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#wasm-attributes
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#envoy-v3-api-msg-config-core-v3-node
//
// The pilot environment variables are described in:
// https://pkg.go.dev/istio.io/istio/pilot/pkg/model

var (
	nodeMetaAnnotations         = []string{"node", "metadata", "ANNOTATIONS"}
	nodeMetaAppContainers       = []string{"node", "metadata", "APP_CONTAINERS"}
	nodeMetaClusterId           = []string{"node", "metadata", "CLUSTER_ID"}
	nodeMetaEnvoyPrometheusPort = []string{"node", "metadata", "ENVOY_PROMETHEUS_PORT"}
	nodeMetaEnvoyStatusPort     = []string{"node", "metadata", "ENVOY_STATUS_PORT"}
	nodeMetaInstanceIps         = []string{"node", "metadata", "INSTANCE_IPS"}
	nodeMetaInterceptionMode    = []string{"node", "metadata", "INTERCEPTION_MODE"}
	nodeMetaIstioProxySha       = []string{"node", "metadata", "ISTIO_PROXY_SHA"}
	nodeMetaIstioVersion        = []string{"node", "metadata", "ISTIO_VERSION"}
	nodeMetaLabels              = []string{"node", "metadata", "LABELS"}
	nodeMetaMeshId              = []string{"node", "metadata", "MESH_ID"}
	nodeMetaName                = []string{"node", "metadata", "NAME"}
	nodeMetaNamespace           = []string{"node", "metadata", "NAMESPACE"}
	nodeMetaNodeName            = []string{"node", "metadata", "NODE_NAME"}
	nodeMetaOwner               = []string{"node", "metadata", "OWNER"}
	nodeMetaPilotSan            = []string{"node", "metadata", "PILOT_SAN"}
	nodeMetaPodPorts            = []string{"node", "metadata", "POD_PORTS"}
	nodeMetaServiceAccount      = []string{"node", "metadata", "SERVICE_ACCOUNT"}
	nodeMetaWorkloadName        = []string{"node", "metadata", "WORKLOAD_NAME"}
)

// GetNodeMetaAnnotations returns the node annotations
func GetNodeMetaAnnotations() (map[string]string, error) {
	return getPropertyStringMap(nodeMetaAnnotations)
}

// GetNodeMetaAppContainers returns the app containers of the node
func GetNodeMetaAppContainers() (string, error) {
	return getPropertyString(nodeMetaAppContainers)
}

// GetNodeMetaClusterId returns the cluster ID of the node, which defines the
// cluster the node belongs to
func GetNodeMetaClusterId() (string, error) {
	return getPropertyString(nodeMetaClusterId)
}

// GetNodeMetaEnvoyPrometheusPort returns the Envoy Prometheus port of the node
func GetNodeMetaEnvoyPrometheusPort() (float64, error) {
	return getPropertyFloat64(nodeMetaEnvoyPrometheusPort)
}

// GetNodeMetaEnvoyStatusPort returns the Envoy status port of the node
func GetNodeMetaEnvoyStatusPort() (float64, error) {
	return getPropertyFloat64(nodeMetaEnvoyStatusPort)
}

// GetNodeMetaInstanceIps returns the instance IPs of the node
func GetNodeMetaInstanceIps() (string, error) {
	return getPropertyString(nodeMetaInstanceIps)
}

// GetNodeMetaInterceptionMode returns the interception mode of the node
//
// Possible values:
//
//	REDIRECT	: REDIRECT mode uses iptables REDIRECT to NAT and redirect to Envoy. This mode
//							loses source IP addresses during redirection
//	TPROXY		: TPROXY mode uses iptables TPROXY to redirect to Envoy. This mode preserves both
//							the source and destination IP addresses and ports, so that they can be used for
//							advanced filtering and manipulation. This mode also configures the sidecar to
//							run with the CAP_NET_ADMIN capability, which is required to use TPROXY
//	NONE			: NONE mode does not configure redirect to Envoy at all. This is an advanced
//							configuration that typically requires changes to user applications.
func GetNodeMetaInterceptionMode() (IstioTrafficInterceptionMode, error) {
	result, err := getPropertyString(nodeMetaInterceptionMode)
	if err != nil {
		return IstioTrafficInterceptionModeRedirect, err
	}
	mode, err := ParseIstioTrafficInterceptionMode(result)
	if err != nil {
		return IstioTrafficInterceptionModeRedirect, err
	}
	return mode, nil
}

// GetNodeMetaIstioProxySha returns the Istio proxy SHA of the node
func GetNodeMetaIstioProxySha() (string, error) {
	return getPropertyString(nodeMetaIstioProxySha)
}

// GetNodeMetaIstioVersion returns the Istio version of the node
func GetNodeMetaIstioVersion() (string, error) {
	return getPropertyString(nodeMetaIstioVersion)
}

// GetNodeMetaLabels returns the labels of the node
func GetNodeMetaLabels() (map[string]string, error) {
	return getPropertyStringMap(nodeMetaLabels)
}

// GetNodeMetaMeshId returns the mesh ID of the node
func GetNodeMetaMeshId() (string, error) {
	return getPropertyString(nodeMetaMeshId)
}

// GetNodeMetaName returns the name of the node
func GetNodeMetaName() (string, error) {
	return getPropertyString(nodeMetaName)
}

// GetNodeMetaNamespace returns the namespace of the node
func GetNodeMetaNamespace() (string, error) {
	return getPropertyString(nodeMetaNamespace)
}

// GetNodeMetaNodeName returns the node name of the node
func GetNodeMetaNodeName() (string, error) {
	return getPropertyString(nodeMetaNodeName)
}

// GetNodeMetaOwner returns the owner of the node (opaque string). Typically, this is the
// owning controller of of the workload instance (ex: k8s deployment for a k8s pod)
func GetNodeMetaOwner() (string, error) {
	return getPropertyString(nodeMetaOwner)
}

// GetNodeMetaPilotSan returns the pilot SAN (subject alternate names) of the node's xDS server
func GetNodeMetaPilotSan() ([]string, error) {
	return getPropertyStringSlice(nodeMetaPilotSan)
}

// GetNodeMetaPodPorts returns the pod ports of the node. This is used to lookup named ports
func GetNodeMetaPodPorts() (string, error) {
	return getPropertyString(nodeMetaPodPorts)
}

// GetNodeMetaServiceAccount returns the service account of the node
func GetNodeMetaServiceAccount() (string, error) {
	return getPropertyString(nodeMetaServiceAccount)
}

// GetNodeMetaWorkloadName returns the workload name of the node
func GetNodeMetaWorkloadName() (string, error) {
	return getPropertyString(nodeMetaWorkloadName)
}
