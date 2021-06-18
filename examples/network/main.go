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
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

var (
	connectionCounterName = "proxy_wasm_go.connection_counter"
	counter               proxywasm.MetricCounter
)

func main() {
	proxywasm.SetNewRootContextFn(newRootContext)
}

func newRootContext(contextID uint32) types.RootContext {
	return &rootContext{}
}

type rootContext struct {
	// Embed the default root context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultRootContext
}

// Override DefaultRootContext.
func (ctx *rootContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	counter = proxywasm.DefineCounterMetric(connectionCounterName)
	return types.OnVMStartStatusOK
}

// Override DefaultRootContext.
func (ctx *rootContext) NewTcpContext(contextID uint32) types.TcpContext {
	return &networkContext{}
}

type networkContext struct {
	// Embed the default tcp context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultTcpContext
}

// Override DefaultTcpContext.
func (ctx *networkContext) OnNewConnection() types.Action {
	proxywasm.LogInfo("new connection!")
	return types.ActionContinue
}

// Override DefaultTcpContext.
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

// Override DefaultTcpContext.
func (ctx *networkContext) OnDownstreamClose(types.PeerType) {
	proxywasm.LogInfo("downstream connection close!")
	return
}

// Override DefaultTcpContext.
func (ctx *networkContext) OnUpstreamData(dataSize int, endOfStream bool) types.Action {
	if dataSize == 0 {
		return types.ActionContinue
	}

	ret, err := proxywasm.GetProperty([]string{"upstream", "address"})
	if err != nil {
		proxywasm.LogCriticalf("failed to get upstream data: %v", err)
		return types.ActionContinue
	}

	proxywasm.LogInfof("remote address: %s", string(ret))

	data, err := proxywasm.GetUpstreamData(0, dataSize)
	if err != nil && err != types.ErrorStatusNotFound {
		proxywasm.LogCritical(err.Error())
	}

	proxywasm.LogInfof("<<<<<< upstream data received <<<<<<\n%s", string(data))
	return types.ActionContinue
}

// Override DefaultTcpContext.
func (ctx *networkContext) OnStreamDone() {
	counter.Increment(1)
	proxywasm.LogInfo("connection complete!")
}
