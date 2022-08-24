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

func TestData(t *testing.T) {
	vmTest(t, func(t *testing.T, vm types.VMContext) {
		opt := proxytest.NewEmulatorOption().WithVMContext(vm)
		host, reset := proxytest.NewHostEmulator(opt)
		defer reset()

		// Call OnVMStart -> set initial value.
		require.Equal(t, types.OnVMStartStatusOK, host.StartVM())
		// Initialize http context.
		contextID := host.InitializeHttpContext()
		// Call OnHttpRequestHeaders.
		action := host.CallOnRequestHeaders(contextID, nil, false)
		require.Equal(t, types.ActionContinue, action)

		// Check Envoy logs.
		logs := host.GetInfoLogs()
		require.Contains(t, logs, "shared value: 10000001")

		// Call OnHttpRequestHeaders again.
		action = host.CallOnRequestHeaders(contextID, nil, false)
		require.Equal(t, types.ActionContinue, action)
		action = host.CallOnRequestHeaders(contextID, nil, false)
		require.Equal(t, types.ActionContinue, action)

		// Check Envoy logs.
		logs = host.GetInfoLogs()
		require.Contains(t, logs, "shared value: 10000003")
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
