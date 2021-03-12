// Copyright 2020-2021 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestNetwork_OnNewConnection(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// Call OnVMStart -> initialize metric
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())

	// OnNewConnection is called.
	_, action := host.InitializeConnection()
	require.Equal(t, types.ActionContinue, action)

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	assert.Contains(t, logs, "new connection!")
}

func TestNetwork_OnDownstreamClose(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// OnNewConnection is called.
	contextID, action := host.InitializeConnection()
	require.Equal(t, types.ActionContinue, action)

	// OnDownstreamClose is called.
	host.CloseDownstreamConnection(contextID)

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	assert.Contains(t, logs, "downstream connection close!")
}

func TestNetwork_OnDownstreamData(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// OnNewConnection is called.
	contextID, action := host.InitializeConnection()
	require.Equal(t, types.ActionContinue, action)

	// OnDownstreamData is called.
	msg := "this is downstream data"
	data := []byte(msg)
	host.CallOnDownstreamData(contextID, data)

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	assert.Contains(t, logs, ">>>>>> downstream data received >>>>>>\n"+msg)
}

func TestNetwork_OnUpstreamData(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// OnNewConnection is called.
	contextID, action := host.InitializeConnection()
	require.Equal(t, types.ActionContinue, action)

	// OnUpstreamData is called.
	msg := "this is upstream data"
	data := []byte(msg)
	host.CallOnUpstreamData(contextID, data)

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	assert.Contains(t, logs, "<<<<<< upstream data received <<<<<<\n"+msg)
}

func TestNetwork_counter(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// Call OnVMStart -> initialize metric
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())

	// OnNewConnection is called.
	contextID, action := host.InitializeConnection()
	require.Equal(t, types.ActionContinue, action)

	// call OnStreamDone on contextID -> increment the connection counter.
	host.CompleteConnection(contextID)

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	assert.Contains(t, logs, "connection complete!")

	// Check counter metric.
	value, err := host.GetCounterMetric("proxy_wasm_go.connection_counter")
	require.NoError(t, err)
	assert.Equal(t, uint64(1), value)
}
