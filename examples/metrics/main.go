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

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct{}

// Implement types.VMContext.
func (*vmContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	return types.OnVMStartStatusOK
}

// Implement types.VMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &metricPluginContext{
		counter: proxywasm.DefineCounterMetric("proxy_wasm_go.request_counter"),
	}
}

type metricPluginContext struct {
	// Embed the default root context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
	counter proxywasm.MetricCounter
}

// Override DefaultPluginContext.
func (ctx *metricPluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &metricHttpContext{counter: ctx.counter}
}

type metricHttpContext struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	counter proxywasm.MetricCounter
}

// Override DefaultHttpContext.
func (ctx *metricHttpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	prev := ctx.counter.Get()
	proxywasm.LogInfof("previous value of %s: %d", "proxy_wasm_go.request_counter", prev)

	ctx.counter.Increment(1)
	proxywasm.LogInfo("incremented")
	return types.ActionContinue
}
