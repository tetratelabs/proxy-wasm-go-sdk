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
	vmTest(t, func(t *testing.T, vm types.VMContext) {
		opt := proxytest.NewEmulatorOption().WithVMContext(vm)
		host, reset := proxytest.NewHostEmulator(opt)
		defer reset()

		// Call OnVMStart.
		require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())
		require.Equal(t, tickMilliseconds, host.GetTickPeriod())

		// Register foreign function named "compress".
		host.RegisterForeignFunction("compress", func(b []byte) []byte { return b })

		for i := 1; i < 10; i++ {
			host.Tick()
			// Check Envoy logs.
			logs := host.GetInfoLogs()
			require.Contains(t, logs, fmt.Sprintf("foreign function (compress) called: %d, result: %s", i, "68656c6c6f20776f726c6421"))
		}
	})
}

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
