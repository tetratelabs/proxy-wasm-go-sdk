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

	for i := 1; i < 10; i++ {
		host.Tick() // call OnTick
		// Check Envoy logs.
		logs := host.GetLogs(types.LogLevelInfo)
		require.Contains(t, logs, fmt.Sprintf("CallForeignFunction result: %s", "hello world!"))
	}

}
