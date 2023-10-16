package properties

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
	result, err := getPropertyUint64(responseCode)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// GetResponseCodeDetails returns the internal response code details (subject to change).
func GetResponseCodeDetails() (string, error) {
	result, err := getPropertyString(responseCodeDetails)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetResponseFlags returns additional details about the response beyond the standard
// response code encoded as a bit-vector.
//
// https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage#config-access-log-format-response-flags
func GetResponseFlags() (uint64, error) {
	result, err := getPropertyUint64(responseFlags)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// GetResponseGrpcStatusCode returns the response gRPC status code.
func GetResponseGrpcStatusCode() (uint64, error) {
	result, err := getPropertyUint64(responseGrpcStatusCode)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// GetResponseHeaders returns all response headers indexed by the lower-cased header name.
func GetResponseHeaders() (map[string]string, error) {
	result, err := getPropertyStringMap(responseHeaders)
	if err != nil {
		return map[string]string{}, err
	}
	return result, nil
}

// GetResponseTrailers returns all response trailers indexed by the lower-cased trailer name.
func GetResponseTrailers() (map[string]string, error) {
	result, err := getPropertyStringMap(responseTrailers)
	if err != nil {
		return map[string]string{}, err
	}
	return result, nil
}

// GetResponseSize returns the size of the response body.
func GetResponseSize() (uint64, error) {
	result, err := getPropertyUint64(responseSize)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// GetResponseTotalSize returns the total size of the response including the approximate
// uncompressed size of the headers and the trailers.
func GetResponseTotalSize() (uint64, error) {
	result, err := getPropertyUint64(responseTotalSize)
	if err != nil {
		return 0, err
	}
	return result, nil
}
