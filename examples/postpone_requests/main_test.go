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

func TestSetEffectiveContext(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnVMStart.
	require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())
	require.Equal(t, tickMilliseconds, host.GetTickPeriod())

	// Initialize context.
	contextID := host.InitializeHttpContext()

	// Call OnHttpRequestHeaders.
	action := host.CallOnRequestHeaders(contextID, [][2]string{}, false)
	require.Equal(t, types.ActionPause, action)

	// Call OnTick.
	host.Tick()

	action = host.GetCurrentHttpStreamAction(contextID)
	require.Equal(t, types.ActionContinue, action)
}
