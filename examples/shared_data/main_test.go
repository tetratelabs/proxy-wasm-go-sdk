package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestData(t *testing.T) {

	ctx := data{}
	host, done := proxytest.NewRootFilterHost(ctx, nil, nil)
	defer done() // release the host emulation lock so that other test cases can insert their own host emulation

	host.StartVM() // set initial value

	ctx.OnHttpRequestHeaders(0, false) // update
	logs := host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 0)

	assert.Equal(t, "shared value: 1", logs[len(logs)-1])

	ctx.OnHttpRequestHeaders(0, false) // update
	ctx.OnHttpRequestHeaders(0, false) // update

	logs = host.GetLogs(types.LogLevelInfo)
	assert.Equal(t, "shared value: 3", logs[len(logs)-1])
}
