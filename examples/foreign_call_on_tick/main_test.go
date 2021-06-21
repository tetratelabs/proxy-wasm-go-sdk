package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestPluginContext_OnTick(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewPluginContext(newPluginContext)
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnVMStart.
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())
	require.Equal(t, tickMilliseconds, host.GetTickPeriod())

	// Register foreign function named "compress".
	host.RegisterForeignFunction("compress", func(b []byte) []byte { return b })

	for i := 1; i < 10; i++ {
		host.Tick()
		// Check Envoy logs.
		logs := host.GetInfoLogs()
		require.Contains(t, logs, fmt.Sprintf("foreign function (compress) called: %d, result: %s", i, "68656c6c6f20776f726c6421"))
	}
}
