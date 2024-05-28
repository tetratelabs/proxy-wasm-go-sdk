// Copyright 2020-2024 Tetrate
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
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

// vmContext implements types.VMContext.
type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// NewPluginContext implements types.VMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{counter: proxywasm.DefineCounterMetric("proxy_wasm_go.connection_counter")}
}

// pluginContext implements types.PluginContext.
type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
	counter proxywasm.MetricCounter
}

// NewTcpContext implements types.PluginContext.
func (ctx *pluginContext) NewTcpContext(contextID uint32) types.TcpContext {
	return &networkContext{counter: ctx.counter}
}

// networkContext implements types.TcpContext.
type networkContext struct {
	// Embed the default tcp context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultTcpContext
	counter proxywasm.MetricCounter
}

// OnNewConnection implements types.TcpContext.
func (ctx *networkContext) OnNewConnection() types.Action {
	proxywasm.LogInfo("new connection!")
	return types.ActionContinue
}

// OnDownstreamData implements types.TcpContext.
func (ctx *networkContext) OnDownstreamData(dataSize int, endOfStream bool) types.Action {
	if dataSize == 0 {
		return types.ActionContinue
	}

	data, err := proxywasm.GetDownstreamData(0, dataSize)
	if err != nil && err != types.ErrorStatusNotFound {
		proxywasm.LogCriticalf("failed to get downstream data: %v", err)
		return types.ActionContinue
	}

	proxywasm.LogInfof(">>>>>> downstream data received >>>>>>\n%s", string(data))
	return types.ActionContinue
}

// OnDownstreamClose implements types.TcpContext.
func (ctx *networkContext) OnDownstreamClose(types.PeerType) {
	proxywasm.LogInfo("downstream connection close!")
}

// OnUpstreamData implements types.TcpContext.
func (ctx *networkContext) OnUpstreamData(dataSize int, endOfStream bool) types.Action {
	if dataSize == 0 {
		return types.ActionContinue
	}

	// Get the remote ip address of the upstream cluster.
	address, err := proxywasm.GetProperty([]string{"upstream", "address"})
	if err != nil {
		proxywasm.LogWarnf("failed to get upstream remote address: %v", err)
	}

	proxywasm.LogInfof("remote address: %s", string(address))

	// Get the upstream cluster's metadata in the cluster configuration.
	metadataKeyValues, err := proxywasm.GetPropertyMap([]string{"cluster_metadata", "filter_metadata", "location"})
	if err != nil {
		proxywasm.LogWarnf("failed to get upstream location metadata: %v", err)
	}

	for _, metadata := range metadataKeyValues {
		key, value := metadata[0], metadata[1]
		proxywasm.LogInfof("upstream cluster metadata location[%s]=%s", string(key), string(value))
	}

	data, err := proxywasm.GetUpstreamData(0, dataSize)
	if err != nil && err != types.ErrorStatusNotFound {
		proxywasm.LogCritical(err.Error())
	}

	proxywasm.LogInfof("<<<<<< upstream data received <<<<<<\n%s", string(data))
	return types.ActionContinue
}

// OnStreamDone implements types.TcpContext.
func (ctx *networkContext) OnStreamDone() {
	ctx.counter.Increment(1)
	proxywasm.LogInfo("connection complete!")
}
