package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHttpRouting_OnHttpRequestHeaders(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	t.Run("canary", func(t *testing.T) {
		now = func() int { return 0 }
		// Initialize http context.
		id := host.InitializeHttpContext()
		hs := types.Headers{{":authority", "my-host.com"}}
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
		now = func() int { return 1 }
		// Initialize http context.
		id := host.InitializeHttpContext()
		hs := types.Headers{{":authority", "my-host.com"}}
		// Call OnHttpResponseHeaders.
		action := host.CallOnRequestHeaders(id,
			hs, false)
		require.Equal(t, types.ActionContinue, action)
		resultHeaders := host.GetCurrentRequestHeaders(id)
		require.Len(t, resultHeaders, 1)
		require.Equal(t, ":authority", resultHeaders[0][0])
		require.Equal(t, "my-host.com", resultHeaders[0][1])
	})
}
