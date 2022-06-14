// These tests are supposed to run with `proxytest` build tag, and this way we can leverage the testing framework in "proxytest" package.
// The framework emulates the expected behavior of Envoyproxy, and you can test your extensions without running Envoy and with
// the standard Go CLI. To run tests, simply run
// go test -tags=proxytest ./...

//go:build proxytest

package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHttpHeaders_OnHttpRequestHeaders(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Initialize http context.
	id := host.InitializeHttpContext()

	// Call OnHttpResponseHeaders.
	hs := [][2]string{{"key1", "value1"}, {"key2", "value2"}}
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
	logs := host.GetInfoLogs()
	require.Contains(t, logs, fmt.Sprintf("%d finished", id))
	require.Contains(t, logs, "request header --> key2: value2")
	require.Contains(t, logs, "request header --> key1: value1")
}

func TestHttpHeaders_OnHttpResponseHeaders(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Initialize http context.
	id := host.InitializeHttpContext()

	// Call OnHttpResponseHeaders.
	hs := [][2]string{{"key1", "value1"}, {"key2", "value2"}}
	action := host.CallOnResponseHeaders(id, hs, false)
	require.Equal(t, types.ActionContinue, action)

	// Call OnHttpStreamDone.
	host.CompleteHttpContext(id)

	resHeaders := host.GetCurrentResponseHeaders(id)
	require.Equal(t, "key1", resHeaders[0][0])
	require.Equal(t, "value1", resHeaders[0][1])
	require.Equal(t, "key2", resHeaders[1][0])
	require.Equal(t, "value2", resHeaders[1][1])

	// Check Envoy logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, fmt.Sprintf("%d finished", id))
	require.Contains(t, logs, "response header <-- key2: value2")
	require.Contains(t, logs, "response header <-- key1: value1")
}
