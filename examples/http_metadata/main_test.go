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

func TestHttpHeaders_OnHttpRequestMetadata(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Initialize http context.
	id := host.InitializeHttpContext()

	hs := [][2]string{{"key1", "value1"}, {"key2", "value2"}}
	// Call OnRequestMetadata.
	action := host.CallOnRequestMetadata(id, hs)
	require.Equal(t, types.ActionContinue, action)

	// Call OnHttpStreamDone.
	host.CompleteHttpContext(id)

	// Check Envoy logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, fmt.Sprintf("%d finished", id))
	require.Contains(t, logs, "request metadata --> 2")
}

func TestHttpHeaders_OnHttpResponseMetadata(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Initialize http context.
	id := host.InitializeHttpContext()

	hs := [][2]string{{"key1", "value1"}, {"key2", "value2"}}
	// Call OnHttpResponseMetata.
	action := host.CallOnResponseMetadata(id, hs)
	require.Equal(t, types.ActionContinue, action)

	// Call OnHttpStreamDone.
	host.CompleteHttpContext(id)

	// Check Envoy logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, fmt.Sprintf("%d finished", id))
	require.Contains(t, logs, "response metadata <-- 2")
}
