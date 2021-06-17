package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestMetric(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnVMStart.
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())

	// Initialize http context.
	contextID := host.InitializeHttpContext()
	exp := uint64(3)
	for i := uint64(0); i < exp; i++ {
		// Call OnRequestHeaders
		action := host.CallOnRequestHeaders(contextID, nil, false)
		require.Equal(t, types.ActionContinue, action)
	}

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	require.Contains(t, logs, "incremented")

	// Check metrics.
	value, err := host.GetCounterMetric(metricsName)
	require.NoError(t, err)
	require.Equal(t, uint64(3), value)
}
