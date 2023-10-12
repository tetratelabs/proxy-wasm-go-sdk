package properties

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
)

func TestGetRequestPath(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(requestPath, []byte("/headers"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetRequestPath()
	require.NoError(t, err)
	require.Equal(t, "/headers", result)
}

func TestGetRequestUrlPath(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(requestUrlPath, []byte("/headers"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetRequestUrlPath()
	require.NoError(t, err)
	require.Equal(t, "/headers", result)
}

func TestGetRequestHost(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(requestHost, []byte("wasm.httpbin.org"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetRequestHost()
	require.NoError(t, err)
	require.Equal(t, "wasm.httpbin.org", result)
}

func TestGetRequestScheme(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "Request Scheme: HTTP",
			input:  "http",
			expect: "http",
		},
		{
			name:   "Request Scheme: HTTPS",
			input:  "https",
			expect: "https",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := proxytest.NewEmulatorOption().WithProperty(requestScheme, []byte(tt.input))
			_, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			result, err := GetRequestScheme()
			require.NoError(t, err)
			require.Equal(t, tt.expect, result)
		})
	}
}

func TestGetRequestMethod(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(requestMethod, []byte("GET"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetRequestMethod()
	require.NoError(t, err)
	require.Equal(t, "GET", result)
}

func TestGetRequestHeaders(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]string
		output map[string]string
	}{
		{
			name: "Test with populated map",
			input: map[string]string{
				":authority":                  "wasm.httpbin.org",
				":method":                     "GET",
				":path":                       "/headers",
				":scheme":                     "http",
				"accept":                      "*/*",
				"user-agent":                  "curl/7.81.0",
				"x-b3-sampled":                "1",
				"x-envoy-decorator-operation": "httpbin.org:80/*",
				"x-envoy-internal":            "true",
				"x-envoy-peer-metadata":       "ChQKDkFQUF9DT05UQUlORVJTE",
				"x-envoy-peer-metadata-id":    "router~10.244.0.13~istio-ingress-67cddc6d57-kk2cr.istio-ingress~istio-ingress.svc.cluster.local",
				"x-forwarded-for":             "10.244.0.1",
				"x-forwarded-proto":           "http",
				"x-request-id":                "7490e0f7-87f0-4c81-92aa-8ea3d5896189",
			},
			output: map[string]string{
				":authority":                  "wasm.httpbin.org",
				":method":                     "GET",
				":path":                       "/headers",
				":scheme":                     "http",
				"accept":                      "*/*",
				"user-agent":                  "curl/7.81.0",
				"x-b3-sampled":                "1",
				"x-envoy-decorator-operation": "httpbin.org:80/*",
				"x-envoy-internal":            "true",
				"x-envoy-peer-metadata":       "ChQKDkFQUF9DT05UQUlORVJTE",
				"x-envoy-peer-metadata-id":    "router~10.244.0.13~istio-ingress-67cddc6d57-kk2cr.istio-ingress~istio-ingress.svc.cluster.local",
				"x-forwarded-for":             "10.244.0.1",
				"x-forwarded-proto":           "http",
				"x-request-id":                "7490e0f7-87f0-4c81-92aa-8ea3d5896189",
			},
		},
		{
			name:   "Test with empty map",
			input:  map[string]string{},
			output: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := proxytest.NewEmulatorOption().WithProperty(requestHeaders, serializeStringMap(tt.input))
			_, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			result, err := GetRequestHeaders()
			require.NoError(t, err)
			require.Equal(t, tt.output, result)
		})
	}
}

func TestGetRequestReferer(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(requestReferer, []byte("https://site.com/page"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetRequestReferer()
	require.NoError(t, err)
	require.Equal(t, "https://site.com/page", result)
}

func TestGetRequestUserAgent(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(requestUserAgent, []byte("curl/7.81.0"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetRequestUserAgent()
	require.NoError(t, err)
	require.Equal(t, "curl/7.81.0", result)
}

func TestGetRequestTime(t *testing.T) {
	now := time.Now().UTC()
	opt := proxytest.NewEmulatorOption().WithProperty(requestTime, serializeTimestamp(now))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetRequestTime()
	require.NoError(t, err)
	require.Equal(t, now, result)
}

func TestGetRequestId(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(requestId, []byte("7490e0f7-87f0-4c81-92aa-8ea3d5896189"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetRequestId()
	require.NoError(t, err)
	require.Equal(t, "7490e0f7-87f0-4c81-92aa-8ea3d5896189", result)
}

func TestGetRequestProtocol(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(requestProtocol, []byte("HTTP/1.1"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetRequestProtocol()
	require.NoError(t, err)
	require.Equal(t, "HTTP/1.1", result)
}

func TestGetRequestQuery(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(requestQuery, []byte("?page=1&limit=10"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetRequestQuery()
	require.NoError(t, err)
	require.Equal(t, "?page=1&limit=10", result)
}

func TestGetRequestDuration(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(requestDuration, serializeUint64(1000))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetRequestDuration()
	require.NoError(t, err)
	require.Equal(t, uint64(1000), result)
}

func TestGetRequestSize(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(requestSize, serializeUint64(256))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetRequestSize()
	require.NoError(t, err)
	require.Equal(t, uint64(256), result)
}

func TestGetRequestTotalSize(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(requestTotalSize, serializeUint64(1024))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetRequestTotalSize()
	require.NoError(t, err)
	require.Equal(t, uint64(1024), result)
}
