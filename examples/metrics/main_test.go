package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestMetric(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// Call OnVMStart.
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())

	// Initialize http context.
	contextID := host.InitializeHttpContext()
	exp := uint64(3)
	for i := uint64(0); i < exp; i++ {
		// Call OnRequestHeaders
		action := host.CallOnRequestHeaders(contextID, nil, false)
		assert.Equal(t, types.ActionContinue, action)
	}

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	assert.Contains(t, logs, "incremented")

	// Check metrics.
	value, err := host.GetCounterMetric(metricsName)
	require.NoError(t, err)
	assert.Equal(t, uint64(3), value)
}
