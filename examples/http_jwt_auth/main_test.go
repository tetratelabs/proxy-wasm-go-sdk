// These tests are supposed to run with `proxytest` build tag, and this way we can leverage the testing framework in "proxytest" package.
// The framework emulates the expected behavior of Envoyproxy, and you can test your extensions without running Envoy and with
// the standard Go CLI. To run tests, simply run
// go test -tags=proxytest ./...

//+build proxytest

package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const (
	// sampleJwtToken is a sample JWT token which is already signed with secret key "secret".
	sampleJwtToken = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.t-IDcSemACt8x4iTMCda8Yhe3iZaWbvV5XKSTbuAn0M`
)

func TestHttpJwtAuth_OnHttpRequestHeaders(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	t.Run("success authorization", func(t *testing.T) {
		// Create http context.
		id := host.InitializeHttpContext()

		// Call OnRequestHeaders.
		action := host.CallOnRequestHeaders(id, [][2]string{
			{"Authorization", fmt.Sprintf("Bearer %s", sampleJwtToken)},
		}, false)

		// Must be continued.
		require.Equal(t, types.ActionContinue, action)

		// Call OnHttpStreamDone
		host.CompleteHttpContext(id)

		// Check Envoy logs.
		logs := host.GetInfoLogs()
		require.Contains(t, logs, fmt.Sprintf("%d finished", id))
	})

	t.Run("400 response", func(t *testing.T) {
		// Create http context.
		id := host.InitializeHttpContext()

		// Call OnRequestHeaders.
		action := host.CallOnRequestHeaders(id, nil, false)

		// Must be paused.
		require.Equal(t, types.ActionPause, action)

		// Check the local response.
		localResponse := host.GetSentLocalResponse(id)
		require.NotNil(t, localResponse)
		require.Equal(t, uint32(400), localResponse.StatusCode)
		require.Equal(t, "authorization header must be provided", string(localResponse.Data))

		// Call OnHttpStreamDone
		host.CompleteHttpContext(id)

		// Check Envoy logs.
		logs := host.GetInfoLogs()
		require.Contains(t, logs, fmt.Sprintf("%d finished", id))
	})

	t.Run("401 response for invalid token", func(t *testing.T) {
		// Create http context.
		id := host.InitializeHttpContext()

		// Call OnRequestHeaders.
		action := host.CallOnRequestHeaders(id, [][2]string{
			{"Authorization", "invalidtoken"},
		}, false)

		// Must be paused.
		require.Equal(t, types.ActionPause, action)

		// Check the local response.
		localResponse := host.GetSentLocalResponse(id)
		require.NotNil(t, localResponse)
		require.Equal(t, uint32(401), localResponse.StatusCode)
		require.Equal(t, "invalid token", string(localResponse.Data))

		// Call OnHttpStreamDone
		host.CompleteHttpContext(id)

		// Check Envoy logs.
		logs := host.GetInfoLogs()
		require.Contains(t, logs, fmt.Sprintf("%d finished", id))
	})
}
