package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHelloWorld_OnTick(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newHelloWorld)
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

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
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnVMStart.
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	require.Contains(t, logs, "proxy_on_vm_start from Go!")
	require.Equal(t, tickMilliseconds, host.GetTickPeriod())
}
