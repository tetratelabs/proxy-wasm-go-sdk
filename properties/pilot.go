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
	result, err := getPropertyStringMap(nodeMetaAnnotations)
	if err != nil {
		return make(map[string]string), err
	}
	return result, nil
}

// GetNodeMetaAppContainers returns the app containers of the node
func GetNodeMetaAppContainers() (string, error) {
	result, err := getPropertyString(nodeMetaAppContainers)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetNodeMetaClusterId returns the cluster ID of the node, which defines the
// cluster the node belongs to
func GetNodeMetaClusterId() (string, error) {
	result, err := getPropertyString(nodeMetaClusterId)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetNodeMetaEnvoyPrometheusPort returns the Envoy Prometheus port of the node
func GetNodeMetaEnvoyPrometheusPort() (float64, error) {
	result, err := getPropertyFloat64(nodeMetaEnvoyPrometheusPort)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// GetNodeMetaEnvoyStatusPort returns the Envoy status port of the node
func GetNodeMetaEnvoyStatusPort() (float64, error) {
	result, err := getPropertyFloat64(nodeMetaEnvoyStatusPort)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// GetNodeMetaInstanceIps returns the instance IPs of the node
func GetNodeMetaInstanceIps() (string, error) {
	result, err := getPropertyString(nodeMetaInstanceIps)
	if err != nil {
		return "", err
	}
	return result, nil
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
		return Redirect, err
	}
	mode, err := ParseIstioTrafficInterceptionMode(result)
	if err != nil {
		return Redirect, err
	}
	return mode, nil
}

// GetNodeMetaIstioProxySha returns the Istio proxy SHA of the node
func GetNodeMetaIstioProxySha() (string, error) {
	result, err := getPropertyString(nodeMetaIstioProxySha)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetNodeMetaIstioVersion returns the Istio version of the node
func GetNodeMetaIstioVersion() (string, error) {
	result, err := getPropertyString(nodeMetaIstioVersion)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetNodeMetaLabels returns the labels of the node
func GetNodeMetaLabels() (map[string]string, error) {
	result, err := getPropertyStringMap(nodeMetaLabels)
	if err != nil {
		return make(map[string]string), err
	}
	return result, nil
}

// GetNodeMetaMeshId returns the mesh ID of the node
func GetNodeMetaMeshId() (string, error) {
	result, err := getPropertyString(nodeMetaMeshId)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetNodeMetaName returns the name of the node
func GetNodeMetaName() (string, error) {
	result, err := getPropertyString(nodeMetaName)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetNodeMetaNamespace returns the namespace of the node
func GetNodeMetaNamespace() (string, error) {
	result, err := getPropertyString(nodeMetaNamespace)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetNodeMetaNodeName returns the node name of the node
func GetNodeMetaNodeName() (string, error) {
	result, err := getPropertyString(nodeMetaNodeName)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetNodeMetaOwner returns the owner of the node (opaque string). Typically, this is the
// owning controller of of the workload instance (ex: k8s deployment for a k8s pod)
func GetNodeMetaOwner() (string, error) {
	result, err := getPropertyString(nodeMetaOwner)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetNodeMetaPilotSan returns the pilot SAN (subject alternate names) of the node's xDS server
func GetNodeMetaPilotSan() ([]string, error) {
	result, err := getPropertyStringSlice(nodeMetaPilotSan)
	if err != nil {
		return make([]string, 0), err
	}
	return result, nil
}

// GetNodeMetaPodPorts returns the pod ports of the node. This is used to lookup named ports
func GetNodeMetaPodPorts() (string, error) {
	result, err := getPropertyString(nodeMetaPodPorts)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetNodeMetaServiceAccount returns the service account of the node
func GetNodeMetaServiceAccount() (string, error) {
	result, err := getPropertyString(nodeMetaServiceAccount)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetNodeMetaWorkloadName returns the workload name of the node
func GetNodeMetaWorkloadName() (string, error) {
	result, err := getPropertyString(nodeMetaWorkloadName)
	if err != nil {
		return "", err
	}
	return result, nil
}
