package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHttpBody_OnHttpRequestBody(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newContext)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// Create http context.
	id := host.InitializeHttpContext()

	// Call OnRequestBody.
	action := host.CallOnRequestBody(id, []byte(`{ "initial": "request body" }`), false)
	require.Equal(t, types.ActionContinue, action)

	logs := host.GetLogs(types.LogLevelInfo)

	// Check Envoy logs.
	assert.Contains(t, logs, "on http request body finished")
	assert.Contains(t, logs, `initial request body: { "initial": "request body" }`)
	assert.Contains(t, logs, "body size: 29")
}
