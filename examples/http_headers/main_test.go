package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHttpHeaders_OnHttpRequestHeaders(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// Initialize http context.
	id := host.InitializeHttpContext()

	// Call OnHttpResponseHeaders.
	hs := types.Headers{{"key1", "value1"}, {"key2", "value2"}}
	action := host.CallOnRequestHeaders(id,
		hs, false)
	require.Equal(t, types.ActionContinue, action)

	// Check headers.
	resultHeaders := host.GetCurrentRequestHeaders(id)
	var found bool
	for _, val := range resultHeaders {
		if val[0] == "test" {
			require.Equal(t, "best", val[1])
			found = true
		}
	}
	require.True(t, found)

	// Call OnHttpStreamDone.
	host.CompleteHttpContext(id)

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	require.Contains(t, logs, fmt.Sprintf("%d finished", id))
	require.Contains(t, logs, "request header --> key2: value2")
	require.Contains(t, logs, "request header --> key1: value1")
}

func TestHttpHeaders_OnHttpResponseHeaders(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// Initialize http context.
	id := host.InitializeHttpContext()

	// Call OnHttpResponseHeaders.
	hs := types.Headers{{"key1", "value1"}, {"key2", "value2"}}
	action := host.CallOnResponseHeaders(id, hs, false)
	require.Equal(t, types.ActionContinue, action)

	// Call OnHttpStreamDone.
	host.CompleteHttpContext(id)

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	require.Contains(t, logs, fmt.Sprintf("%d finished", id))
	require.Contains(t, logs, "response header <-- key2: value2")
	require.Contains(t, logs, "response header <-- key1: value1")
}
