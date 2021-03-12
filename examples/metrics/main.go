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
	proxywasm.SetNewRootContext(newRootContext)
}

var counter proxywasm.MetricCounter

const metricsName = "proxy_wasm_go.request_counter"

type metricRootContext struct {
	// You'd better embed the default root context
	// so that you don't need to reimplement all the methods by yourself.
	proxywasm.DefaultRootContext
}

func newRootContext(uint32) proxywasm.RootContext {
	return &metricRootContext{}
}

// Override DefaultRootContext.
func (ctx *metricRootContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	counter = proxywasm.DefineCounterMetric(metricsName)
	return types.OnVMStartStatusOK
}

// Override DefaultRootContext.
func (*metricRootContext) NewHttpContext(contextID uint32) proxywasm.HttpContext {
	return &metricHttpContext{}
}

type metricHttpContext struct {
	// You'd better embed the default http context
	// so that you don't need to reimplement all the methods by yourself.
	proxywasm.DefaultHttpContext
}

// Override DefaultHttpContext.
func (ctx *metricHttpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	prev := counter.Get()
	proxywasm.LogInfof("previous value of %s: %d", metricsName, prev)

	counter.Increment(1)
	proxywasm.LogInfo("incremented")
	return types.ActionContinue
}
