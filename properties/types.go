package properties

import "fmt"

// EnvoyTrafficDirection identifies the direction of the traffic relative to the local Envoy.
//
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#enum-config-core-v3-trafficdirection
type EnvoyTrafficDirection int

const (
	// EnvoyTrafficDirectionUnspecified means that the direction is not specified.
	EnvoyTrafficDirectionUnspecified EnvoyTrafficDirection = iota
	// EnvoyTrafficDirectionInbound means that the transport is used for incoming traffic.
	EnvoyTrafficDirectionInbound
	// EnvoyTrafficDirectionOutbound means that the transport is used for outgoing traffic.
	EnvoyTrafficDirectionOutbound
)

// String converts the EnvoyTrafficDirection enum value to its corresponding string representation.
// It returns "UNSPECIFIED" for Unspecified, "INBOUND" for Inbound, and "OUTBOUND" for Outbound.
// If the enum value doesn't match any of the predefined values, it defaults to "UNSPECIFIED".
func (t EnvoyTrafficDirection) String() string {
	switch t {
	case EnvoyTrafficDirectionUnspecified:
		return "UNSPECIFIED"
	case EnvoyTrafficDirectionInbound:
		return "INBOUND"
	case EnvoyTrafficDirectionOutbound:
		return "OUTBOUND"
	}
	return "UNSPECIFIED"
}

// EnvoyLocality identifies location of where either Envoy runs or where upstream hosts run.
//
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#config-core-v3-locality
type EnvoyLocality struct {
	Region  string
	Zone    string
	Subzone string
}

// EnvoyExtension holds version and identification for an Envoy extension.
//
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#config-core-v3-extension
type EnvoyExtension struct {
	Name     string
	Category string
	TypeUrls []string
}

// IstioFilterMetadata provides additional inputs to filters based on matched listeners,
// filter chains, routes and endpoints. It is structured as a map, usually from
// filter name (in reverse DNS format) to metadata specific to the filter. Metadata
// key-values for a filter are merged as connection and request handling occurs,
// with later values for the same key overriding earlier values.
//
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#config-core-v3-metadata
type IstioFilterMetadata struct {
	Config   string
	Services []IstioService
}

// IstioService holds information of the host, name and namespace of an Istio Service.
type IstioService struct {
	Host      string
	Name      string
	Namespace string
}

// IstioProxyStatsMatcher holds proxy stats name matches for stats creation. Note this is in addition to the minimum Envoy stats that
// Istio generates by default.
//
// https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#ProxyConfig-ProxyStatsMatcher
type IstioProxyStatsMatcher struct {
	InclusionPrefixes []string
	InclusionRegexps  []string
	InclusionSuffixes []string
}

// IstioTrafficInterceptionMode indicates how traffic to/from the workload is captured and sent to Envoy. This
// should not be confused with the CaptureMode in the API that indicates how the user wants traffic to be
// intercepted for the listener. IstioTrafficInterceptionMode is always derived from the Proxy metadata.
//
// https://pkg.go.dev/istio.io/istio/pilot/pkg/model#TrafficInterceptionMode
type IstioTrafficInterceptionMode int

const (
	// IstioTrafficInterceptionModeNone indicates that the workload is not using IPtables for traffic interception.
	IstioTrafficInterceptionModeNone IstioTrafficInterceptionMode = iota
	// IstioTrafficInterceptionModeTproxy implies traffic intercepted by IPtables with TPROXY mode.
	IstioTrafficInterceptionModeTproxy
	// IstioTrafficInterceptionModeRedirect implies traffic intercepted by IPtables with REDIRECT mode. This is our default mode.
	IstioTrafficInterceptionModeRedirect
)

// String converts the IstioTrafficInterceptionMode enum value to its corresponding string representation.
// It returns "NONE" for None, "TPROXY" for Tproxy, and "REDIRECT" for Redirect.
// If the enum value doesn't match any of the predefined values, it defaults to "REDIRECT".
func (t IstioTrafficInterceptionMode) String() string {
	switch t {
	case IstioTrafficInterceptionModeNone:
		return "NONE"
	case IstioTrafficInterceptionModeTproxy:
		return "TPROXY"
	case IstioTrafficInterceptionModeRedirect:
		return "REDIRECT"
	}
	return "REDIRECT"
}

// ParseIstioTrafficInterceptionMode converts a string representation of IstioTrafficInterceptionMode to
// its corresponding enum value. It returns None for "NONE", Tproxy for "TPROXY", and Redirect for "REDIRECT".
// If the provided string doesn't match any of the predefined values, it returns an error and the default
// value Redirect.
func ParseIstioTrafficInterceptionMode(s string) (IstioTrafficInterceptionMode, error) {
	switch s {
	case "NONE":
		return IstioTrafficInterceptionModeNone, nil
	case "TPROXY":
		return IstioTrafficInterceptionModeTproxy, nil
	case "REDIRECT":
		return IstioTrafficInterceptionModeRedirect, nil
	default:
		return IstioTrafficInterceptionModeRedirect, fmt.Errorf("invalid IstioTrafficInterceptionMode: %s", s)
	}
}
