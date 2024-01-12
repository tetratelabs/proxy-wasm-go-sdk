package properties

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
)

func TestGetResponseCode(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(responseCode, serializeUint64(200))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetResponseCode()
	require.NoError(t, err)
	require.Equal(t, uint64(200), result)
}

func TestGetResponseCodeDetails(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(responseCodeDetails, []byte("Not Found"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetResponseCodeDetails()
	require.NoError(t, err)
	require.Equal(t, "Not Found", result)
}

func TestGetResponseFlags(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(responseFlags, serializeUint64(123))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetResponseFlags()
	require.NoError(t, err)
	require.Equal(t, uint64(123), result)
}

func TestGetResponseFlagsShort(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(responseFlags, serializeUint64(123))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetResponseFlagsShort()
	require.NoError(t, err)
	require.Equal(t, "FailedLocalHealthCheck,LocalReset,NoHealthyUpstream,UpstreamConnectionFailure,UpstreamConnectionTermination,UpstreamRemoteReset", result)
}

func TestGetResponseGrpcStatusCode(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(responseGrpcStatusCode, serializeUint64(200))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetResponseGrpcStatusCode()
	require.NoError(t, err)
	require.Equal(t, uint64(200), result)
}

func TestGetResponseHeaders(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]string
		expect map[string]string
	}{
		{
			name: "With Headers",
			input: map[string]string{
				":status":                          "200",
				"access-control-allow-credentials": "true",
				"access-control-allow-origin":      "*",
				"connection":                       "keep-alive",
				"content-length":                   "1383",
				"content-type":                     "application/json",
				"date":                             "Fri, 13 Oct 2023 11:38:01 GMT",
				"server":                           "gunicorn/19.9.0",
				"x-envoy-upstream-service-time":    "199",
			},
			expect: map[string]string{
				":status":                          "200",
				"access-control-allow-credentials": "true",
				"access-control-allow-origin":      "*",
				"connection":                       "keep-alive",
				"content-length":                   "1383",
				"content-type":                     "application/json",
				"date":                             "Fri, 13 Oct 2023 11:38:01 GMT",
				"server":                           "gunicorn/19.9.0",
				"x-envoy-upstream-service-time":    "199",
			},
		},
		{
			name:   "Empty Headers",
			input:  map[string]string{},
			expect: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := proxytest.NewEmulatorOption().WithProperty(responseHeaders, serializeStringMap(tt.input))
			_, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			result, err := GetResponseHeaders()
			require.NoError(t, err)
			require.Equal(t, tt.expect, result)
		})
	}
}

func TestGetResponseTrailers(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]string
		expect map[string]string
	}{
		{
			name: "With Trailers",
			input: map[string]string{
				"Expires:": "Wed, 21 Oct 2015 07:28:00 GMT",
			},
			expect: map[string]string{
				"Expires:": "Wed, 21 Oct 2015 07:28:00 GMT",
			},
		},
		{
			name:   "Empty Trailers",
			input:  map[string]string{},
			expect: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := proxytest.NewEmulatorOption().WithProperty(responseTrailers, serializeStringMap(tt.input))
			_, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			result, err := GetResponseTrailers()
			require.NoError(t, err)
			require.Equal(t, tt.expect, result)
		})
	}
}

func TestGetResponseSize(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(responseSize, serializeUint64(512))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetResponseSize()
	require.NoError(t, err)
	require.Equal(t, uint64(512), result)
}

func TestGetResponseTotalSize(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(responseTotalSize, serializeUint64(2048))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetResponseTotalSize()
	require.NoError(t, err)
	require.Equal(t, uint64(2048), result)
}
