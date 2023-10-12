package properties

import "fmt"

// Identifies the direction of the traffic relative to the local Envoy
//
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#enum-config-core-v3-trafficdirection
type EnvoyTrafficDirection int

const (
	// ⁣Default option is unspecified
	Unspecified EnvoyTrafficDirection = iota
	// The transport is used for incoming traffic
	Inbound
	// ⁣The transport is used for outgoing traffic
	Outbound
)

func (t EnvoyTrafficDirection) String() string {
	switch t {
	case Unspecified:
		return "UNSPECIFIED"
	case Inbound:
		return "INBOUND"
	case Outbound:
		return "OUTBOUND"
	}
	return "UNSPECIFIED"
}

// Identifies location of where either Envoy runs or where upstream hosts run
//
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#config-core-v3-locality
type EnvoyLocality struct {
	Region  string
	Zone    string
	Subzone string
}

// Version and identification for an Envoy extension
//
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#config-core-v3-extension
type EnvoyExtension struct {
	Name     string
	Category string
	TypeUrls []string
}

// Metadata provides additional inputs to filters based on matched listeners,
// filter chains, routes and endpoints. It is structured as a map, usually from
// filter name (in reverse DNS format) to metadata specific to the filter. Metadata
// key-values for a filter are merged as connection and request handling occurs,
// with later values for the same key overriding earlier values
//
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#config-core-v3-metadata
type IstioFilterMetadata struct {
	Config   string
	Services []IstioService
}

type IstioService struct {
	Host      string
	Name      string
	Namespace string
}

// Proxy stats name matchers for stats creation. Note this is in addition to the minimum Envoy stats that
// Istio generates by default
//
// https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#ProxyConfig-ProxyStatsMatcher
type IstioProxyStatsMatcher struct {
	InclusionPrefixes []string
	InclusionRegexps  []string
	InclusionSuffixes []string
}

// TrafficInterceptionMode indicates how traffic to/from the workload is captured and sent to Envoy. This
// should not be confused with the CaptureMode in the API that indicates how the user wants traffic to be
// intercepted for the listener. TrafficInterceptionMode is always derived from the Proxy metadata.
//
// https://pkg.go.dev/istio.io/istio/pilot/pkg/model#TrafficInterceptionMode
type IstioTrafficInterceptionMode int

const (
	// InterceptionNone indicates that the workload is not using IPtables for traffic interception
	None IstioTrafficInterceptionMode = iota
	// InterceptionTproxy implies traffic intercepted by IPtables with TPROXY mode
	Tproxy
	// InterceptionRedirect implies traffic intercepted by IPtables with REDIRECT mode. This is our default mode
	Redirect
)

// String converts the IstioTrafficInterceptionMode enum value to its corresponding string representation.
// It returns "NONE" for None, "TPROXY" for Tproxy, and "REDIRECT" for Redirect.
// If the enum value doesn't match any of the predefined values, it defaults to "REDIRECT".
func (t IstioTrafficInterceptionMode) String() string {
	switch t {
	case None:
		return "NONE"
	case Tproxy:
		return "TPROXY"
	case Redirect:
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
		return None, nil
	case "TPROXY":
		return Tproxy, nil
	case "REDIRECT":
		return Redirect, nil
	default:
		return Redirect, fmt.Errorf("invalid IstioTrafficInterceptionMode: %s", s)
	}
}
