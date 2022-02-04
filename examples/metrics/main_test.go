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

func TestMetric(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnVMStart.
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())

	// Initialize http context.
	headers := [][2]string{{"my-custom-header", "foo"}}
	contextID := host.InitializeHttpContext()
	exp := uint64(3)
	for i := uint64(0); i < exp; i++ {
		// Call OnRequestHeaders
		action := host.CallOnRequestHeaders(contextID, headers, false)
		require.Equal(t, types.ActionContinue, action)
	}

	// Check metrics.
	value, err := host.GetCounterMetric("custom_header_value_counts_value=foo_reporter=wasmgosdk")
	require.NoError(t, err)
	require.Equal(t, uint64(3), value)
}
