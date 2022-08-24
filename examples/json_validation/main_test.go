// The framework emulates the expected behavior of Envoyproxy, and you can test your extensions without running Envoy and with
// the standard Go CLI. To run tests, simply run
// go test ./...

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestOnHTTPRequestHeaders(t *testing.T) {
	type testCase struct {
		contentType    string
		expectedAction types.Action
	}

	vmTest(t, func(t *testing.T, vm types.VMContext) {
		for name, tCase := range map[string]testCase{
			"fails due to unsupported content type": {
				contentType:    "text/html",
				expectedAction: types.ActionPause,
			},
			"success for JSON": {
				contentType:    "application/json",
				expectedAction: types.ActionContinue,
			},
		} {
			t.Run(name, func(t *testing.T) {
				opt := proxytest.NewEmulatorOption().WithVMContext(vm)
				host, reset := proxytest.NewHostEmulator(opt)
				defer reset()

				require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())

				id := host.InitializeHttpContext()

				hs := [][2]string{{"content-type", tCase.contentType}}

				action := host.CallOnRequestHeaders(id, hs, false)
				assert.Equal(t, tCase.expectedAction, action)
			})
		}
	})
}

func TestOnHTTPRequestBody(t *testing.T) {
	type testCase struct {
		body           string
		expectedAction types.Action
	}

	vmTest(t, func(t *testing.T, vm types.VMContext) {

		for name, tCase := range map[string]testCase{
			"pauses due to invalid payload": {
				body:           "invalid_payload",
				expectedAction: types.ActionPause,
			},
			"pauses due to unknown keys": {
				body:           `{"unknown_key":"unknown_value"}`,
				expectedAction: types.ActionPause,
			},
			"success": {
				body:           "{\"my_key\":\"my_value\"}",
				expectedAction: types.ActionContinue,
			},
		} {
			t.Run(name, func(t *testing.T) {
				opt := proxytest.
					NewEmulatorOption().
					WithPluginConfiguration([]byte(`{"requiredKeys": ["my_key"]}`)).
					WithVMContext(vm)
				host, reset := proxytest.NewHostEmulator(opt)
				defer reset()

				require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())

				id := host.InitializeHttpContext()

				action := host.CallOnRequestBody(id, []byte(tCase.body), true)
				assert.Equal(t, tCase.expectedAction, action)
			})
		}
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
