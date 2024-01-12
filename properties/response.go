package properties

import (
	"sort"
	"strings"
)

// This file hosts helper functions to retrieve response-related properties as described in:
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#response-attributes

var (
	responseCode           = []string{"response", "code"}
	responseCodeDetails    = []string{"response", "code_details"}
	responseFlags          = []string{"response", "flags"}
	responseGrpcStatusCode = []string{"response", "grpc_status"}
	responseHeaders        = []string{"response", "headers"}
	responseTrailers       = []string{"response", "trailers"}
	responseSize           = []string{"response", "size"}
	responseTotalSize      = []string{"response", "total_size"}
)

// GetResponseCode returns the response HTTP status code.
func GetResponseCode() (uint64, error) {
	return getPropertyUint64(responseCode)
}

// GetResponseCodeDetails returns the internal response code details (subject to change).
func GetResponseCodeDetails() (string, error) {
	return getPropertyString(responseCodeDetails)
}

// GetResponseFlags returns additional details about the response beyond the standard
// response code encoded as a bit-vector.
//
// https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage#config-access-log-format-response-flags
func GetResponseFlags() (uint64, error) {
	return getPropertyUint64(responseFlags)
}

// ResponseFlag enum
const (
	FailedLocalHealthCheck uint64 = 1 << iota
	NoHealthyUpstream
	UpstreamRequestTimeout
	LocalReset
	UpstreamRemoteReset
	UpstreamConnectionFailure
	UpstreamConnectionTermination
	UpstreamOverflow
	NoRouteFound
	DelayInjected
	FaultInjected
	RateLimited
	UnauthorizedExternalService
	RateLimitServiceError
	DownstreamConnectionTermination
	UpstreamRetryLimitExceeded
	StreamIdleTimeout
	InvalidEnvoyRequestHeaders
	DownstreamProtocolError
	UpstreamMaxStreamDurationReached
	ResponseFromCacheFilter
	NoFilterConfigFound
	DurationTimeout
	UpstreamProtocolError
	NoClusterFound
	OverloadManager
	DnsResolutionFailed
	DropOverLoad
	// LastFlag should always be equal to the last flag.
	LastFlag = DropOverLoad
)

var flagNames = map[uint64]string{
	FailedLocalHealthCheck:           "FailedLocalHealthCheck",
	NoHealthyUpstream:                "NoHealthyUpstream",
	UpstreamRequestTimeout:           "UpstreamRequestTimeout",
	LocalReset:                       "LocalReset",
	UpstreamRemoteReset:              "UpstreamRemoteReset",
	UpstreamConnectionFailure:        "UpstreamConnectionFailure",
	UpstreamConnectionTermination:    "UpstreamConnectionTermination",
	UpstreamOverflow:                 "UpstreamOverflow",
	NoRouteFound:                     "NoRouteFound",
	DelayInjected:                    "DelayInjected",
	FaultInjected:                    "FaultInjected",
	RateLimited:                      "RateLimited",
	UnauthorizedExternalService:      "UnauthorizedExternalService",
	RateLimitServiceError:            "RateLimitServiceError",
	DownstreamConnectionTermination:  "DownstreamConnectionTermination",
	UpstreamRetryLimitExceeded:       "UpstreamRetryLimitExceeded",
	StreamIdleTimeout:                "StreamIdleTimeout",
	InvalidEnvoyRequestHeaders:       "InvalidEnvoyRequestHeaders",
	DownstreamProtocolError:          "DownstreamProtocolError",
	UpstreamMaxStreamDurationReached: "UpstreamMaxStreamDurationReached",
	ResponseFromCacheFilter:          "ResponseFromCacheFilter",
	NoFilterConfigFound:              "NoFilterConfigFound",
	DurationTimeout:                  "DurationTimeout",
	UpstreamProtocolError:            "UpstreamProtocolError",
	NoClusterFound:                   "NoClusterFound",
	OverloadManager:                  "OverloadManager",
	DnsResolutionFailed:              "DnsResolutionFailed",
	DropOverLoad:                     "DropOverLoad",
}

// getResponseFlagsShort returns a comma-separated string of flag names set in the flags parameter.
func GetResponseFlagsShort() (string, error) {
	var names []string

	flags, err := GetResponseFlags()
	if err != nil {
		return "", err
	}

	for flag, name := range flagNames {
		if flags&flag != 0 {
			names = append(names, name)
		}
	}

	sort.Strings(names)
	return strings.Join(names, ","), nil
}

// GetResponseGrpcStatusCode returns the response gRPC status code.
func GetResponseGrpcStatusCode() (uint64, error) {
	return getPropertyUint64(responseGrpcStatusCode)
}

// GetResponseHeaders returns all response headers indexed by the lower-cased header name.
func GetResponseHeaders() (map[string]string, error) {
	return getPropertyStringMap(responseHeaders)
}

// GetResponseTrailers returns all response trailers indexed by the lower-cased trailer name.
func GetResponseTrailers() (map[string]string, error) {
	return getPropertyStringMap(responseTrailers)
}

// GetResponseSize returns the size of the response body.
func GetResponseSize() (uint64, error) {
	return getPropertyUint64(responseSize)
}

// GetResponseTotalSize returns the total size of the response including the approximate
// uncompressed size of the headers and the trailers.
func GetResponseTotalSize() (uint64, error) {
	return getPropertyUint64(responseTotalSize)
}
