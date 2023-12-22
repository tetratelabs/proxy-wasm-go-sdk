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

func TestSetBodyContext_OnHttpRequestHeaders(t *testing.T) {
	vmTest(t, func(t *testing.T, vm types.VMContext) {
		opt := proxytest.NewEmulatorOption().WithVMContext(vm)
		host, reset := proxytest.NewHostEmulator(opt)
		defer reset()

		t.Run("remove content length", func(t *testing.T) {
			// Create http context.
			id := host.InitializeHttpContext()

			// Call OnRequestHeaders.
			action := host.CallOnRequestHeaders(id, [][2]string{
				{"content-length", "10"},
				{"buffer-operation", "replace"},
			}, false)

			// Must be continued.
			require.Equal(t, types.ActionContinue, action)

			// Check the final request headers
			headers := host.GetCurrentRequestHeaders(id)
			require.Equal(t,
				[][2]string{{"buffer-operation", "replace"}},
				headers,
				"content-length header must be removed.")
		})

		t.Run("400 response", func(t *testing.T) {
			// Create http context.
			id := host.InitializeHttpContext()

			// Call OnRequestHeaders without "content-length"
			action := host.CallOnRequestHeaders(id, nil, false)

			// Must be paused.
			require.Equal(t, types.ActionPause, action)

			// Check the local response.
			localResponse := host.GetSentLocalResponse(id)
			require.NotNil(t, localResponse)
			require.Equal(t, uint32(400), localResponse.StatusCode)
			require.Equal(t, "content must be provided", string(localResponse.Data))
		})
	})
}

func TestSetBodyContext_OnHttpRequestBody(t *testing.T) {
	vmTest(t, func(t *testing.T, vm types.VMContext) {
		opt := proxytest.NewEmulatorOption().WithVMContext(vm)
		host, reset := proxytest.NewHostEmulator(opt)
		defer reset()

		t.Run("pause until EOS", func(t *testing.T) {
			// Create http context.
			id := host.InitializeHttpContext()

			// Call OnRequestBody.
			action := host.CallOnRequestBody(id, []byte("aaaa"), false /* end of stream */)

			// Must be paused
			require.Equal(t, types.ActionPause, action)
		})

		t.Run("append", func(t *testing.T) {
			// Create http context.
			id := host.InitializeHttpContext()

			// Call OnRequestHeaders.
			action := host.CallOnRequestHeaders(id, [][2]string{
				{"content-length", "10"},
				{"buffer-operation", "append"},
			}, false)

			// Must be continued.
			require.Equal(t, types.ActionContinue, action)

			// Call OnRequestBody.
			action = host.CallOnRequestBody(id, []byte(`[original body]`), true)
			require.Equal(t, types.ActionContinue, action)

			// Check Envoy logs.
			logs := host.GetInfoLogs()
			require.Contains(t, logs, `original request body: [original body]`)

			// Check the final request body is the replaced one.
			require.Equal(t, "[original body][this is appended body]", string(host.GetCurrentRequestBody(id)))
		})

		t.Run("prepend", func(t *testing.T) {
			// Create http context.
			id := host.InitializeHttpContext()

			// Call OnRequestHeaders.
			action := host.CallOnRequestHeaders(id, [][2]string{
				{"content-length", "10"},
				{"buffer-operation", "prepend"},
			}, false)

			// Must be continued.
			require.Equal(t, types.ActionContinue, action)

			// Call OnRequestBody.
			action = host.CallOnRequestBody(id, []byte(`[original body]`), true)
			require.Equal(t, types.ActionContinue, action)

			// Check Envoy logs.
			logs := host.GetInfoLogs()
			require.Contains(t, logs, `original request body: [original body]`)

			// Check the final request body is the replaced one.
			require.Equal(t, "[this is prepended body][original body]", string(host.GetCurrentRequestBody(id)))
		})

		t.Run("replace", func(t *testing.T) {
			// Create http context.
			id := host.InitializeHttpContext()

			// Call OnRequestHeaders.
			action := host.CallOnRequestHeaders(id, [][2]string{
				{"content-length", "10"},
				{"buffer-operation", "replace"},
			}, false)

			// Must be continued.
			require.Equal(t, types.ActionContinue, action)

			// Call OnRequestBody.
			action = host.CallOnRequestBody(id, []byte(`[original body]`), true)
			require.Equal(t, types.ActionContinue, action)

			// Check Envoy logs.
			logs := host.GetInfoLogs()
			require.Contains(t, logs, `original request body: [original body]`)

			// Check the final request body is the replaced one.
			require.Equal(t, "[this is replaced body]", string(host.GetCurrentRequestBody(id)))
		})
	})
}

func TestSetBodyContext_OnHttpResponseBody(t *testing.T) {
	vmTest(t, func(t *testing.T, vm types.VMContext) {
		opt := proxytest.NewEmulatorOption().WithVMContext(vm)
		host, reset := proxytest.NewHostEmulator(opt)
		defer reset()

		t.Run("append", func(t *testing.T) {
			// Create http context.
			id := host.InitializeHttpContext()

			// Call OnRequestHeaders.
			action := host.CallOnRequestHeaders(id, [][2]string{
				{"buffer-replace-at", "response"},
				{"content-length", "10"},
				{"buffer-operation", "append"},
			}, false)

			// Must be continued.
			require.Equal(t, types.ActionContinue, action)

			// Call OnResponseBody.
			action = host.CallOnResponseBody(id, []byte(`[original body]`), true)
			require.Equal(t, types.ActionContinue, action)

			// Check Envoy logs.
			logs := host.GetInfoLogs()
			require.Contains(t, logs, `original response body: [original body]`)

			// Check the final response body is the replaced one.
			require.Equal(t, "[original body][this is appended body]", string(host.GetCurrentResponseBody(id)))
		})

		t.Run("prepend", func(t *testing.T) {
			// Create http context.
			id := host.InitializeHttpContext()

			// Call OnRequestHeaders.
			action := host.CallOnRequestHeaders(id, [][2]string{
				{"buffer-replace-at", "response"},
				{"content-length", "10"},
				{"buffer-operation", "prepend"},
			}, false)

			// Must be continued.
			require.Equal(t, types.ActionContinue, action)

			// Call OnRequestBody.
			action = host.CallOnResponseBody(id, []byte(`[original body]`), true)
			require.Equal(t, types.ActionContinue, action)

			// Check Envoy logs.
			logs := host.GetInfoLogs()
			require.Contains(t, logs, `original response body: [original body]`)

			// Check the final request body is the replaced one.
			require.Equal(t, "[this is prepended body][original body]", string(host.GetCurrentResponseBody(id)))
		})

		t.Run("replace", func(t *testing.T) {
			// Create http context.
			id := host.InitializeHttpContext()

			// Call OnRequestHeaders.
			action := host.CallOnRequestHeaders(id, [][2]string{
				{"buffer-replace-at", "response"},
				{"content-length", "10"},
				{"buffer-operation", "replace"},
			}, false)

			// Must be continued.
			require.Equal(t, types.ActionContinue, action)

			// Call OnRequestBody.
			action = host.CallOnResponseBody(id, []byte(`[original body]`), true)
			require.Equal(t, types.ActionContinue, action)

			// Check Envoy logs.
			logs := host.GetInfoLogs()
			require.Contains(t, logs, `original response body: [original body]`)

			// Check the final request body is the replaced one.
			require.Equal(t, "[this is replaced body]", string(host.GetCurrentResponseBody(id)))
		})
	})
}

func TestEchoBodyContext_OnHttpRequestBody(t *testing.T) {
	vmTest(t, func(t *testing.T, vm types.VMContext) {
		opt := proxytest.NewEmulatorOption().WithVMContext(vm).
			WithPluginConfiguration([]byte("echo"))
		host, reset := proxytest.NewHostEmulator(opt)
		defer reset()

		require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())

		t.Run("pause until EOS", func(t *testing.T) {
			// Create http context.
			id := host.InitializeHttpContext()

			// Call OnRequestBody.
			action := host.CallOnRequestBody(id, []byte("aaaa"), false /* end of stream */)

			// Must be paused
			require.Equal(t, types.ActionPause, action)
		})

		t.Run("echo request", func(t *testing.T) {
			// Create http context.
			id := host.InitializeHttpContext()

			for _, frame := range []string{"frame1...", "frame2...", "frame3..."} {
				// Call OnRequestHeaders without "content-length"
				action := host.CallOnRequestBody(id, []byte(frame), false /* end of stream */)

				// Must be paused.
				require.Equal(t, types.ActionPause, action)
			}

			// End stream.
			action := host.CallOnRequestBody(id, nil, true /* end of stream */)

			// Must be paused.
			require.Equal(t, types.ActionPause, action)

			// Check the local response.
			localResponse := host.GetSentLocalResponse(id)
			require.NotNil(t, localResponse)
			require.Equal(t, uint32(200), localResponse.StatusCode)
			require.Equal(t, "frame1...frame2...frame3...", string(localResponse.Data))
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
