package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestMetric(t *testing.T) {
	host := proxytest.NewHostEmulator(nil, nil,
		newRootContext, nil, newHttpContext,
	)
	defer host.Done() // release the host emulation lock so that other test cases can insert their own host emulation

	host.StartVM() // call OnVMStart: define metric

	contextID := host.HttpFilterInitContext()
	exp := uint64(3)
	for i := uint64(0); i < exp; i++ {
		host.HttpFilterPutRequestHeaders(contextID, nil)
	}

	logs := host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 0)

	assert.Equal(t, "incremented", logs[len(logs)-1])

	value, err := counter.Get()
	require.NoError(t, err)
	assert.Equal(t, uint64(3), value)
}
