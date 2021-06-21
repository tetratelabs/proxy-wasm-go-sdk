package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestPluginContext_OnTick(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
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
		require.Contains(t, logs, fmt.Sprintf("called! %d", i))
	}

}

func TestPluginContext_OnVMStart(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnVMStart.
	require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())
	require.Equal(t, tickMilliseconds, host.GetTickPeriod())
}
