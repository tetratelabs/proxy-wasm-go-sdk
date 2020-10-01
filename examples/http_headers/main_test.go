package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHttpHeaders_OnHttpRequestHeaders(t *testing.T) {
	host, done := proxytest.NewHttpFilterHost(newContext)
	defer done()
	id := host.InitContext()

	hs := [][2]string{{"key1", "value1"}, {"key2", "value2"}}
	host.PutRequestHeaders(id, hs) // call OnHttpRequestHeaders

	logs := host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 1)

	assert.Equal(t, "request header --> key2: value2", logs[len(logs)-1])
	assert.Equal(t, "request header --> key1: value1", logs[len(logs)-2])
}

func TestHttpHeaders_OnHttpResponseHeaders(t *testing.T) {
	host, done := proxytest.NewHttpFilterHost(newContext)
	defer done()
	id := host.InitContext()

	hs := [][2]string{{"key1", "value1"}, {"key2", "value2"}}
	host.PutResponseHeaders(id, hs) // call OnHttpResponseHeaders

	logs := host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 1)

	assert.Equal(t, "response header <-- key2: value2", logs[len(logs)-1])
	assert.Equal(t, "response header <-- key1: value1", logs[len(logs)-2])
}
