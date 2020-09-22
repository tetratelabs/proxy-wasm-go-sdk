// Copyright 2020 Tetrate
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
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func newStreamContext(uint32) proxywasm.StreamContext {
	return context{}
}

func TestNetwork_OnNewConnection(t *testing.T) {
	host, done := proxytest.NewNetworkFilterHost(newStreamContext)
	defer done() // release the host emulation lock so that other test cases can insert their own host emulation

	_ = host.InitConnection() // OnNewConnection is called

	logs := host.GetLogs(types.LogLevelInfo) // retrieve logs emitted to Envoy
	assert.Equal(t, logs[0], "new connection!")
}

func TestNetwork_OnDownstreamClose(t *testing.T) {
	host, done := proxytest.NewNetworkFilterHost(newStreamContext)
	defer done() // release the host emulation lock so that other test cases can insert their own host emulation

	contextID := host.InitConnection()        // OnNewConnection is called
	host.CloseDownstreamConnection(contextID) // OnDownstreamClose is called

	logs := host.GetLogs(types.LogLevelInfo) // retrieve logs emitted to Envoy
	require.Len(t, logs, 2)
	assert.Equal(t, logs[1], "downstream connection close!")
}

func TestNetwork_OnDownstreamData(t *testing.T) {
	host, done := proxytest.NewNetworkFilterHost(newStreamContext)
	defer done() // release the host emulation lock so that other test cases can insert their own host emulation

	contextID := host.InitConnection() // OnNewConnection is called

	msg := "this is downstream data"
	data := []byte(msg)
	host.PutDownstreamData(contextID, data) // OnDownstreamData is called

	logs := host.GetLogs(types.LogLevelInfo) // retrieve logs emitted to Envoy
	assert.Equal(t, "downstream data received: "+msg, logs[len(logs)-1])
}

func TestNetwork_OnUpstreamData(t *testing.T) {
	host, done := proxytest.NewNetworkFilterHost(newStreamContext)
	defer done() // release the host emulation lock so that other test cases can insert their own host emulation

	contextID := host.InitConnection() // OnNewConnection is called

	msg := "this is upstream data"
	data := []byte(msg)
	host.PutUpstreamData(contextID, data) // OnUpstreamData is called

	logs := host.GetLogs(types.LogLevelInfo) // retrieve logs emitted to Envoy
	assert.Equal(t, "upstream data received: "+msg, logs[len(logs)-1])
}

func TestNetwork_counter(t *testing.T) {
	host, done := proxytest.NewNetworkFilterHost(newStreamContext)
	defer done() // release the host emulation lock so that other test cases can insert their own host emulation

	context{}.OnVMStart(0) // init metric

	contextID := host.InitConnection()
	host.CompleteConnection(contextID) // call OnDone on contextID -> increment the connection counter

	logs := host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 0)

	assert.Equal(t, "connection complete!", logs[len(logs)-1])
	actual, err := counter.Get()
	require.NoError(t, err)
	assert.Equal(t, uint64(1), actual)
}
