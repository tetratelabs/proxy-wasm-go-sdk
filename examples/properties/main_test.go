// The framework emulates the expected behavior of Envoyproxy, and you can test your extensions without running Envoy and with
// the standard Go CLI. To run tests, simply run
// go test ./...

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestProperties_OnHttpRequestHeaders(t *testing.T) {
	vmTest(t, func(t *testing.T, vm types.VMContext) {
		opt := proxytest.NewEmulatorOption().WithVMContext(vm)
		host, reset := proxytest.NewHostEmulator(opt)
		defer reset()

		t.Run("route is unauthenticated", func(t *testing.T) {
			// Initialize http context.
			id := host.InitializeHttpContext()

			// Call OnRequestHeaders.
			action := host.CallOnRequestHeaders(id, nil, false)
			require.Equal(t, types.ActionContinue, action)

			// Call OnHttpStreamDone.
			host.CompleteHttpContext(id)

			// Check Envoy logs.
			logs := host.GetInfoLogs()
			require.Contains(t, logs, fmt.Sprintf("no auth header for route"))
			require.Contains(t, logs, fmt.Sprintf("%d finished", id))
		})

		// Set property
		path := "auth"
		data := "cookie"
		host.SetProperty(append(propertyPrefix, path), []byte(data))

		// Get property
		actualData, _ := host.GetProperty(append(propertyPrefix, path))
		require.Equal(t, string(actualData), data)

		t.Run("user is authenticated", func(t *testing.T) {
			// Initialize http context.
			id := host.InitializeHttpContext()

			// Call OnRequestHeaders.
			action := host.CallOnRequestHeaders(id, [][2]string{
				{"cookie", "value"},
			}, false)
			require.Equal(t, types.ActionContinue, action)

			// Call OnHttpStreamDone.
			host.CompleteHttpContext(id)

			// Check Envoy logs.
			logs := host.GetInfoLogs()
			require.Contains(t, logs, fmt.Sprintf("auth header is \"%s\"", data))
			require.Contains(t, logs, fmt.Sprintf("%d finished", id))
		})

		t.Run("user is unauthenticated", func(t *testing.T) {
			// Initialize http context.
			id := host.InitializeHttpContext()

			// Call OnRequestHeaders.
			action := host.CallOnRequestHeaders(id, nil, false)
			require.Equal(t, types.ActionPause, action)

			// Call OnHttpStreamDone.
			host.CompleteHttpContext(id)

			// Check the local response.
			localResponse := host.GetSentLocalResponse(id)
			require.NotNil(t, localResponse)
			require.Equal(t, uint32(401), localResponse.StatusCode)
			require.Nil(t, localResponse.Data)
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
