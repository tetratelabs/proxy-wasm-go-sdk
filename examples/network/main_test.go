// These tests are supposed to run with `proxytest` build tag, and this way we can leverage the testing framework in "proxytest" package.
// The framework emulates the expected behavior of Envoyproxy, and you can test your extensions without running Envoy and with
// the standard Go CLI. To run tests, simply run
// go test -tags=proxytest ./...

//go:build proxytest

package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestNetwork_OnNewConnection(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Initialize plugin
	require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())

	// OnNewConnection is called.
	_, action := host.InitializeConnection()
	require.Equal(t, types.ActionContinue, action)

	// Check Envoy logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, "new connection!")
}

func TestNetwork_OnDownstreamClose(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// OnNewConnection is called.
	contextID, action := host.InitializeConnection()
	require.Equal(t, types.ActionContinue, action)

	// OnDownstreamClose is called.
	host.CloseDownstreamConnection(contextID)

	// Check Envoy logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, "downstream connection close!")
}

func TestNetwork_OnDownstreamData(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// OnNewConnection is called.
	contextID, action := host.InitializeConnection()
	require.Equal(t, types.ActionContinue, action)

	// OnDownstreamData is called.
	msg := "this is downstream data"
	data := []byte(msg)
	host.CallOnDownstreamData(contextID, data)

	// Check Envoy logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, ">>>>>> downstream data received >>>>>>\n"+msg)
}

func TestNetwork_OnUpstreamData(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// OnNewConnection is called.
	contextID, action := host.InitializeConnection()
	require.Equal(t, types.ActionContinue, action)

	// OnUpstreamData is called.
	msg := "this is upstream data"
	data := []byte(msg)
	host.CallOnUpstreamData(contextID, data)

	// Check Envoy logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, "<<<<<< upstream data received <<<<<<\n"+msg)
}

func TestNetwork_counter(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnVMStart -> initialize metric
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())

	// OnNewConnection is called.
	contextID, action := host.InitializeConnection()
	require.Equal(t, types.ActionContinue, action)

	// call OnStreamDone on contextID -> increment the connection counter.
	host.CompleteConnection(contextID)

	// Check Envoy logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, "connection complete!")

	// Check counter metric.
	value, err := host.GetCounterMetric("proxy_wasm_go.connection_counter")
	require.NoError(t, err)
	require.Equal(t, uint64(1), value)
}
