package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestNetwork_OnNewConnection(t *testing.T) {
	host, done := proxytest.NewNetworkFilterHost(newHelloWorld)
	defer done() // release the host emulation lock so that other test cases can insert their own host emulation

	_ = host.InitConnection() // OnNewConnection is called

	logs := host.GetLogs(types.LogLevelInfo) // retrieve logs emitted to Envoy
	assert.Equal(t, logs[0], "new connection!")
}

func TestNetwork_OnDownstreamClose(t *testing.T) {
	host, done := proxytest.NewNetworkFilterHost(newHelloWorld)
	defer done() // release the host emulation lock so that other test cases can insert their own host emulation

	contextID := host.InitConnection()        // OnNewConnection is called
	host.CloseDownstreamConnection(contextID) // OnDownstreamClose is called

	logs := host.GetLogs(types.LogLevelInfo) // retrieve logs emitted to Envoy
	require.Len(t, logs, 2)
	assert.Equal(t, logs[1], "downstream connection close!")
}

func TestNetwork_OnDownstreamData(t *testing.T) {
	host, done := proxytest.NewNetworkFilterHost(newHelloWorld)
	defer done() // release the host emulation lock so that other test cases can insert their own host emulation

	contextID := host.InitConnection() // OnNewConnection is called

	msg := "this is downstream data"
	data := []byte(msg)
	host.PutDownstreamData(contextID, data) // OnDownstreamData is called

	logs := host.GetLogs(types.LogLevelInfo) // retrieve logs emitted to Envoy
	assert.Equal(t, "downstream data received: "+msg, logs[len(logs)-1])
}

func TestNetwork_OnUpstreamData(t *testing.T) {
	host, done := proxytest.NewNetworkFilterHost(newHelloWorld)
	defer done() // release the host emulation lock so that other test cases can insert their own host emulation

	contextID := host.InitConnection() // OnNewConnection is called

	msg := "this is upstream data"
	data := []byte(msg)
	host.PutUpstreamData(contextID, data) // OnUpstreamData is called

	logs := host.GetLogs(types.LogLevelInfo) // retrieve logs emitted to Envoy
	assert.Equal(t, "upstream data received: "+msg, logs[len(logs)-1])
}
