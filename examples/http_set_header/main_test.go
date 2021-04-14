package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHttpHeaders_OnHttpRequestHeaders(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// Initialize http context.
	id := host.InitializeHttpContext()

	// Call OnHttpResponseHeaders.
	hs := types.Headers{{"key1", "value1"}, {"key2", "value2"}}
	action := host.CallOnRequestHeaders(id,
		hs, false)
	require.Equal(t, types.ActionContinue, action)

	// Call OnHttpStreamDone.
	host.CompleteHttpContext(id)

	// Check headers.
	logs := host.GetHttpRequestHeaders(id)
	
	var headerValue *string
	
	for _, val := range logs {
        if val[0] == "test" {
            headerValue = &val[1]
        }
    }
	
	assert.Equal(t, *headerValue, "best")
}