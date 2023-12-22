// The framework emulates the expected behavior of Envoyproxy, and you can test your extensions without running Envoy and with
// the standard Go CLI. To run tests, simply run
// go test ./...

package main

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func Test_OnHttpRequestBody(t *testing.T) {
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
		t.Run("pattern found", func(t *testing.T) {
			// Create http context.
			id := host.InitializeHttpContext()

			body := "This is a payload with the pattern word."

			// Call OnRequestHeaders.
			action := host.CallOnRequestHeaders(id, [][2]string{
				{"content-length", strconv.Itoa(len(body))},
			}, false)

			// Must be continued.
			require.Equal(t, types.ActionContinue, action)

			// Call OnRequestBody.
			action = host.CallOnRequestBody(id, []byte(body), true)

			// Must be paused
			require.Equal(t, types.ActionPause, action)

			// Check Envoy logs.
			logs := host.GetInfoLogs()
			require.Contains(t, logs, `pattern found in chunk: 1`)

			// Check the local response.
			localResponse := host.GetSentLocalResponse(id)
			require.NotNil(t, localResponse)
			require.Equal(t, uint32(403), localResponse.StatusCode)
		})
		t.Run("pattern found multiple chunks", func(t *testing.T) {
			// Create http context.
			id := host.InitializeHttpContext()

			chunks := []string{
				"chunk1...",
				"chunk2...",
				"chunk3...",
				"chunk4 with pattern ...",
			}
			var chunksSize int
			for _, chunk := range chunks {
				chunksSize += len(chunk)
			}

			// Call OnRequestHeaders.
			action := host.CallOnRequestHeaders(id, [][2]string{
				{"content-length", strconv.Itoa(chunksSize)},
			}, false)

			// Must be continued.
			require.Equal(t, types.ActionContinue, action)

			for _, chunk := range chunks {
				action := host.CallOnRequestBody(id, []byte(chunk), false /* end of stream */)

				// Must be paused.
				require.Equal(t, types.ActionPause, action)
			}

			// Check Envoy logs.
			logs := host.GetInfoLogs()
			require.Contains(t, logs, `pattern found in chunk: 4`)
			logs = host.GetErrorLogs()
			for _, log := range logs {
				require.NotContains(t, log, `read data does not match`)
			}

			// Check the local response.
			localResponse := host.GetSentLocalResponse(id)
			require.NotNil(t, localResponse)
			require.Equal(t, uint32(403), localResponse.StatusCode)
		})
		t.Run("pattern not found", func(t *testing.T) {
			// Create http context.
			id := host.InitializeHttpContext()

			body := "This is a generic payload."

			// Call OnRequestHeaders.
			action := host.CallOnRequestHeaders(id, [][2]string{
				{"content-length", strconv.Itoa(len(body))},
			}, false)

			// Must be continued.
			require.Equal(t, types.ActionContinue, action)

			// Call OnRequestBody.
			action = host.CallOnRequestBody(id, []byte(body), false)

			// Must be paused
			require.Equal(t, types.ActionPause, action)

			// Call OnRequestBody.
			action = host.CallOnRequestBody(id, nil, true)

			// Must be continued
			require.Equal(t, types.ActionContinue, action)

			// Check Envoy logs.
			logs := host.GetInfoLogs()
			require.Contains(t, logs, `pattern not found`)
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
