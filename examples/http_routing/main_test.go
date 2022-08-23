// These tests are supposed to run with `proxytest` build tag, and this way we can leverage the testing framework in "proxytest" package.
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

func TestHttpRouting_OnHttpRequestHeaders(t *testing.T) {
	vmTest(t, func(t *testing.T, vm types.VMContext) {
		t.Run("canary", func(t *testing.T) {
			opt := proxytest.NewEmulatorOption().WithVMContext(vm).WithPluginConfiguration([]byte{2})
			host, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())

			// Initialize http context.
			id := host.InitializeHttpContext()
			hs := [][2]string{{":authority", "my-host.com"}}
			// Call OnHttpResponseHeaders.
			action := host.CallOnRequestHeaders(id,
				hs, false)
			require.Equal(t, types.ActionContinue, action)
			resultHeaders := host.GetCurrentRequestHeaders(id)
			require.Len(t, resultHeaders, 1)
			require.Equal(t, ":authority", resultHeaders[0][0])
			require.Equal(t, "my-host.com-canary", resultHeaders[0][1])
		})

		t.Run("non-canary", func(t *testing.T) {
			opt := proxytest.NewEmulatorOption().WithVMContext(vm).WithPluginConfiguration([]byte{1})
			host, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())

			// Initialize http context.
			id := host.InitializeHttpContext()
			hs := [][2]string{{":authority", "my-host.com"}}
			// Call OnHttpResponseHeaders.
			action := host.CallOnRequestHeaders(id,
				hs, false)
			require.Equal(t, types.ActionContinue, action)
			resultHeaders := host.GetCurrentRequestHeaders(id)
			require.Len(t, resultHeaders, 1)
			require.Equal(t, ":authority", resultHeaders[0][0])
			require.Equal(t, "my-host.com", resultHeaders[0][1])
		})
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
