package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestMetric(t *testing.T) {

	ctx := metric{}
	host, done := proxytest.NewRootFilterHost(ctx, nil, nil)
	defer done() // release the host emulation lock so that other test cases can insert their own host emulation

	host.StartVM() // define metric

	exp := uint64(3)
	for i := uint64(0); i < exp; i++ {
		ctx.OnHttpRequestHeaders(0, false) // increment
	}

	logs := host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 0)

	assert.Equal(t, "incremented", logs[len(logs)-1])

	value, err := counter.Get()
	require.NoError(t, err)
	assert.Equal(t, uint64(3), value)
}
