package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHttpHeaders_OnHttpRequestHeaders(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewHttpContext(newContext)
	host := proxytest.NewHostEmulator(opt)
	defer host.Done()
	id := host.HttpFilterInitContext()

	hs := types.Headers{{"key1", "value1"}, {"key2", "value2"}}
	host.HttpFilterPutRequestHeaders(id, hs) // call OnHttpRequestHeaders

	host.HttpFilterCompleteHttpStream(id)

	logs := host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 1)

	assert.Equal(t, fmt.Sprintf("%d finished", id), logs[len(logs)-1])
	assert.Equal(t, "request header --> key2: value2", logs[len(logs)-2])
	assert.Equal(t, "request header --> key1: value1", logs[len(logs)-3])
}

func TestHttpHeaders_OnHttpResponseHeaders(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewHttpContext(newContext)
	host := proxytest.NewHostEmulator(opt)
	defer host.Done()
	id := host.HttpFilterInitContext()

	hs := types.Headers{{"key1", "value1"}, {"key2", "value2"}}
	host.HttpFilterPutResponseHeaders(id, hs) // call OnHttpRequestHeaders
	host.HttpFilterCompleteHttpStream(id)     // call OnHttpStreamDone

	logs := host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 1)

	assert.Equal(t, fmt.Sprintf("%d finished", id), logs[len(logs)-1])
	assert.Equal(t, "response header <-- key2: value2", logs[len(logs)-2])
	assert.Equal(t, "response header <-- key1: value1", logs[len(logs)-3])
}
