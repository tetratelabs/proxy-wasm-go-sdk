package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHttpAuthRandom_OnHttpRequestHeaders(t *testing.T) {
	host, done := proxytest.NewHttpFilterHost(newContext)
	defer done()

	id := host.InitContext()
	host.PutRequestHeaders(id, [][2]string{{"key", "value"}}) // OnHttpRequestHeaders called

	require.True(t, host.IsDispatchCalled(id))                     // check if http call is dispatched
	require.Equal(t, types.ActionPause, host.GetCurrentAction(id)) // check if the current action is pause

	logs := host.GetLogs(types.LogLevelInfo)
	require.GreaterOrEqual(t, len(logs), 2)

	assert.Equal(t, "http call dispatched to "+clusterName, logs[len(logs)-1])
	assert.Equal(t, "request header: key: value", logs[len(logs)-2])
}

func TestHttpAuthRandom_OnHttpCallResponse(t *testing.T) {
	host, done := proxytest.NewHttpFilterHost(newContext)
	defer done()

	// http://httpbin.org/uuid
	headers := [][2]string{
		{"HTTP/1.1", "200 OK"}, {"Date:", "Thu, 17 Sep 2020 02:47:07 GMT"},
		{"Content-Type", "application/json"}, {"Content-Length", "53"},
		{"Connection", "keep-alive"}, {"Server", "gunicorn/19.9.0"},
		{"Access-Control-Allow-Origin", "*"}, {"Access-Control-Allow-Credentials", "true"},
	}

	// access granted body
	id := host.InitContext()
	body := []byte(`{"uuid": "7b10a67a-1c67-4199-835b-cbefcd4a63d4"}`)
	host.PutCalloutResponse(id, headers, nil, body)
	assert.Nil(t, host.GetSentLocalResponse(id))

	logs := host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 1)
	assert.Equal(t, "access granted", logs[len(logs)-1])

	// access denied body
	id = host.InitContext()
	body = []byte(`{"uuid": "aaaaaaaa-1c67-4199-835b-cbefcd4a63d4"}`)
	host.PutCalloutResponse(id, headers, nil, body)
	localResponse := host.GetSentLocalResponse(id) // check local responses
	assert.NotNil(t, localResponse)
	logs = host.GetLogs(types.LogLevelInfo)
	assert.Equal(t, "access forbidden", logs[len(logs)-1])

	assert.Equal(t, uint32(403), localResponse.StatusCode)
	assert.Equal(t, []byte("access forbidden"), localResponse.Data)
	require.Len(t, localResponse.Headers, 1)
	assert.Equal(t, "powered-by", localResponse.Headers[0][0])
	assert.Equal(t, "proxy-wasm-go-sdk!!", localResponse.Headers[0][1])
}
