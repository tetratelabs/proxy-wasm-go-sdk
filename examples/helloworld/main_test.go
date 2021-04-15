package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHelloWorld_OnTick(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newHelloWorld)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// Call OnVMStart.
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())
	require.Equal(t, tickMilliseconds, host.GetTickPeriod())

	// Call OnTick.
	host.Tick()

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	require.Contains(t, logs, "OnTick called")
}

func TestHelloWorld_OnVMStart(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newHelloWorld)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// Call OnVMStart.
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	require.Contains(t, logs, "proxy_on_vm_start from Go!")
	require.Equal(t, tickMilliseconds, host.GetTickPeriod())
}
