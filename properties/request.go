package properties

import "time"

// This file hosts helper functions to retrieve request-related properties as described in:
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#request-attributes

var (
	requestPath      = []string{"request", "path"}
	requestUrlPath   = []string{"request", "url_path"}
	requestHost      = []string{"request", "host"}
	requestScheme    = []string{"request", "scheme"}
	requestMethod    = []string{"request", "method"}
	requestHeaders   = []string{"request", "headers"}
	requestReferer   = []string{"request", "referer"}
	requestUserAgent = []string{"request", "useragent"}
	requestTime      = []string{"request", "time"}
	requestId        = []string{"request", "id"}
	requestProtocol  = []string{"request", "protocol"}
	requestQuery     = []string{"request", "query"}
	requestDuration  = []string{"request", "duration"}
	requestSize      = []string{"request", "size"}
	requestTotalSize = []string{"request", "total_size"}
)

// GetRequestPath return the path portion of the URL.
func GetRequestPath() (string, error) {
	return getPropertyString(requestPath)
}

// GetRequestUrlPath returns the path portion of the URL without the query string.
func GetRequestUrlPath() (string, error) {
	return getPropertyString(requestUrlPath)
}

// GetRequestHost returns the host portion of the URL.
func GetRequestHost() (string, error) {
	return getPropertyString(requestHost)
}

// GetRequestScheme returns the scheme portion of the URL e.g. “http”.
func GetRequestScheme() (string, error) {
	return getPropertyString(requestScheme)
}

// GetRequestMethod returns the request method e.g. “GET”.
func GetRequestMethod() (string, error) {
	return getPropertyString(requestMethod)
}

// GetRequestHeaders returns all request headers indexed by the lower-cased header name.
func GetRequestHeaders() (map[string]string, error) {
	return getPropertyStringMap(requestHeaders)
}

// GetRequestReferer returns the referer request header.
func GetRequestReferer() (string, error) {
	return getPropertyString(requestReferer)
}

// GetRequestUserAgent returns the user agent request header.
func GetRequestUserAgent() (string, error) {
	return getPropertyString(requestUserAgent)
}

// GetRequestTime returns the UTC time of the first byte received, approximated to nano-seconds.
func GetRequestTime() (time.Time, error) {
	result, err := getPropertyTimestamp(requestTime)
	if err != nil {
		return time.Now(), err
	}
	return result, nil
}

// GetRequestId returns the request ID corresponding to x-request-id header value.
func GetRequestId() (string, error) {
	return getPropertyString(requestId)
}

// GetRequestProtocol returns the request protocol (“HTTP/1.0”, “HTTP/1.1”, “HTTP/2”, or “HTTP/3”).
func GetRequestProtocol() (string, error) {
	return getPropertyString(requestProtocol)
}

// GetRequestQuery returns the query portion of the URL in the format of “name1=value1&name2=value2”.
func GetRequestQuery() (string, error) {
	return getPropertyString(requestQuery)
}

// GetRequestDuration returns the total duration of the request, approximated to nano-seconds.
func GetRequestDuration() (uint64, error) {
	return getPropertyUint64(requestDuration)
}

// GetRequestSize returns the size of the request body. Content length header is used if available.
func GetRequestSize() (uint64, error) {
	return getPropertyUint64(requestSize)
}

// GetRequestTotalSize returns the total size of the request including the approximate uncompressed size of the headers.
func GetRequestTotalSize() (uint64, error) {
	return getPropertyUint64(requestTotalSize)
}
