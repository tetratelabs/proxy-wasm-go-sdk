package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestRootContext_OnTick(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)

	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// Call OnVMStart.
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())
	require.Equal(t, tickMilliseconds, host.GetTickPeriod())

	// Register foreign function named "compress".
	host.RegisterForeignFunction("compress", func(b []byte) []byte { return b })

	for i := 1; i < 10; i++ {
		host.Tick()
		// Check Envoy logs.
		logs := host.GetLogs(types.LogLevelInfo)
		require.Contains(t, logs, fmt.Sprintf("foreign function (compress) called: %d, result: %s", i, "68656c6c6f20776f726c6421"))
	}
}
