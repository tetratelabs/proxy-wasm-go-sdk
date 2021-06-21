package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHelloWorld_OnTick(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnPluginStart.
	require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())
	require.Equal(t, tickMilliseconds, host.GetTickPeriod())

	// Call OnTick.
	host.Tick()

	// Check Envoy logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, "OnTick called")
}

func TestHelloWorld_OnPluginStart(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnPluginStart.
	require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())

	// Check Envoy logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, "proxy_on_vm_start from Go!")
	require.Equal(t, tickMilliseconds, host.GetTickPeriod())
}
