// These tests are supposed to run with `proxytest` build tag, and this way we can leverage the testing framework in "proxytest" package.
// The framework emulates the expected behavior of Envoyproxy, and you can test your extensions without running Envoy and with
// the standard Go CLI. To run tests, simply run
// go test -tags=proxytest ./...

//go:build proxytest

package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHttpContext_OnHttpRequestHeaders(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Initialize context.
	contextID := host.InitializeHttpContext()

	// Call OnHttpResponseHeaders.
	action := host.CallOnResponseHeaders(contextID,
		[][2]string{{"key", "value"}}, false)
	require.Equal(t, types.ActionPause, action)

	// Verify DispatchHttpCall is called.
	callouts := host.GetCalloutAttributesFromContext(contextID)
	require.Equal(t, len(callouts), 10)

	// At this point, none of dispatched callouts received response.
	// Therefore, the current status must be paused.
	require.Equal(t, types.ActionPause, host.GetCurrentHttpStreamAction(contextID))

	// Emulates that Envoy received all the response to the dispatched callouts.
	for _, callout := range callouts {
		host.CallOnHttpCallResponse(callout.CalloutID, nil, nil, nil)
	}

	// Check if the current action is continued.
	require.Equal(t, types.ActionContinue, host.GetCurrentHttpStreamAction(contextID))

	// Check logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, "pending dispatched requests: 9")
	require.Contains(t, logs, "pending dispatched requests: 1")
	require.Contains(t, logs, "response resumed after processed 10 dispatched request")
}
