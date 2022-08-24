// The framework emulates the expected behavior of Envoyproxy, and you can test your extensions without running Envoy and with
// the standard Go CLI. To run tests, simply run
// go test ./...

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHttpContext_OnHttpRequestHeaders(t *testing.T) {
	vmTest(t, func(t *testing.T, vm types.VMContext) {
		opt := proxytest.NewEmulatorOption().WithVMContext(vm)
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
	})
}

// vmTest executes f twice, once with a types.VMContext that executes plugin code directly
// in the host, and again by executing the plugin code within the compiled main.wasm binary.
// Execution with main.wasm will be skipped if the file cannot be found.
func vmTest(t *testing.T, f func(*testing.T, types.VMContext)) {
	t.Helper()

	t.Run("go", func(t *testing.T) {
		f(t, &vmContext{})
	})

	t.Run("wasm", func(t *testing.T) {
		wasm, err := os.ReadFile("main.wasm")
		if err != nil {
			t.Skip("wasm not found")
		}
		v, err := proxytest.NewWasmVMContext(wasm)
		require.NoError(t, err)
		defer v.Close()
		f(t, v)
	})
}
