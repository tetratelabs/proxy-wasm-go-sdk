package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHttpAuthRandom_OnHttpRequestHeaders(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewHttpContext(newContext)
	host := proxytest.NewHostEmulator(opt)
	defer host.Done()

	contextID := host.HttpFilterInitContext()
	host.HttpFilterPutRequestHeaders(contextID, types.Headers{{"key", "value"}}) // OnHttpRequestHeaders called

	attrs := host.GetCalloutAttributesFromContext(contextID)
	require.Equal(t, len(attrs), 1) // verify DispatchHttpCall is called

	require.Equal(t, "httpbin", attrs[0].Upstream)
	require.Equal(t, types.ActionPause,
		host.HttpFilterGetCurrentStreamAction(contextID)) // check if the current action is pause

	logs := host.GetLogs(types.LogLevelInfo)
	require.GreaterOrEqual(t, len(logs), 2)

	assert.Equal(t, "http call dispatched to "+clusterName, logs[len(logs)-1])
	assert.Equal(t, "request header: key: value", logs[len(logs)-2])
}

func TestHttpAuthRandom_OnHttpCallResponse(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewHttpContext(newContext)
	host := proxytest.NewHostEmulator(opt)
	defer host.Done()

	// http://httpbin.org/uuid
	headers := [][2]string{
		{"HTTP/1.1", "200 OK"}, {"Date:", "Thu, 17 Sep 2020 02:47:07 GMT"},
		{"Content-Type", "application/json"}, {"Content-Length", "53"},
		{"Connection", "keep-alive"}, {"Server", "gunicorn/19.9.0"},
		{"Access-Control-Allow-Origin", "*"}, {"Access-Control-Allow-Credentials", "true"},
	}

	// access granted body
	contextID := host.HttpFilterInitContext()
	host.HttpFilterPutRequestHeaders(contextID, nil) // OnHttpRequestHeaders called

	body := []byte(`{"uuid": "7b10a67a-1c67-4199-835b-cbefcd4a63d4"}`)
	attrs := host.GetCalloutAttributesFromContext(contextID)
	require.Equal(t, len(attrs), 1) // verify DispatchHttpCall is called

	host.PutCalloutResponse(attrs[0].CalloutID, headers, nil, body)
	assert.Nil(t, host.HttpFilterGetSentLocalResponse(contextID))

	logs := host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 1)
	assert.Equal(t, "access granted", logs[len(logs)-1])

	// access denied body
	contextID = host.HttpFilterInitContext()
	host.HttpFilterPutRequestHeaders(contextID, nil) // OnHttpRequestHeaders called

	body = []byte(`{"uuid": "aaaaaaaa-1c67-4199-835b-cbefcd4a63d4"}`)
	attrs = host.GetCalloutAttributesFromContext(contextID)
	require.Equal(t, len(attrs), 1) // verify DispatchHttpCall is called

	host.PutCalloutResponse(attrs[0].CalloutID, headers, nil, body)
	localResponse := host.HttpFilterGetSentLocalResponse(contextID) // check local responses
	assert.NotNil(t, localResponse)
	logs = host.GetLogs(types.LogLevelInfo)
	assert.Equal(t, "access forbidden", logs[len(logs)-1])

	assert.Equal(t, uint32(403), localResponse.StatusCode)
	assert.Equal(t, []byte("access forbidden"), localResponse.Data)
	require.Len(t, localResponse.Headers, 1)
	assert.Equal(t, "powered-by", localResponse.Headers[0][0])
	assert.Equal(t, "proxy-wasm-go-sdk!!", localResponse.Headers[0][1])
}
