// These tests are supposed to run with `proxytest` build tag, and this way we can leverage the testing framework in "proxytest" package.
// The framework emulates the expected behavior of Envoyproxy, and you can test your extensions without running Envoy and with
// the standard Go CLI. To run tests, simply run
// go test -tags=proxytest ./...

//go:build proxytest

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestPluginContext_OnTick(t *testing.T) {
	vmTest(t, func(vm types.VMContext) {
		opt := proxytest.NewEmulatorOption().WithVMContext(vm)
		host, reset := proxytest.NewHostEmulator(opt)
		defer reset()

		// Call OnVMStart.
		require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())
		require.Equal(t, tickMilliseconds, host.GetTickPeriod())

		for i := 1; i < 10; i++ {
			host.Tick() // call OnTick
			attrs := host.GetCalloutAttributesFromContext(proxytest.PluginContextID)
			// Verify DispatchHttpCall is called
			require.Equal(t, len(attrs), i)
			// Receive callout response.
			host.CallOnHttpCallResponse(attrs[0].CalloutID, nil, nil, nil)
			// Check Envoy logs.
			logs := host.GetInfoLogs()
			require.Contains(t, logs, fmt.Sprintf("called %d for contextID=%d", i, proxytest.PluginContextID))
		}
	})
}

func TestPluginContext_OnVMStart(t *testing.T) {
	vmTest(t, func(vm types.VMContext) {
		opt := proxytest.NewEmulatorOption().WithVMContext(vm)
		host, reset := proxytest.NewHostEmulator(opt)
		defer reset()

		// Call OnVMStart.
		require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())
		require.Equal(t, tickMilliseconds, host.GetTickPeriod())
	})
}

func vmTest(t *testing.T, f func(types.VMContext)) {
	t.Helper()

	t.Run("go", func(t *testing.T) {
		f(&vmContext{})
	})

	t.Run("wasm", func(t *testing.T) {
		wasm, err := os.ReadFile("main.wasm")
		if err != nil {
			t.Skip("wasm not found")
		}
		v, err := proxytest.NewWasmVMContext(wasm)
		require.NoError(t, err)
		f(v)
	})
}
