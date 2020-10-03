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
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetNewRootContext(newRootContext)
	proxywasm.SetNewHttpContext(newHttpContext)
}

var counter proxywasm.MetricCounter

const metricsName = "proxy_wasm_go.request_counter"

type metricRootContext struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	proxywasm.DefaultRootContext
}

func newRootContext(uint32) proxywasm.RootContext {
	return &metricRootContext{}
}

// override
func (ctx *metricRootContext) OnVMStart(int) bool {
	ct, err := proxywasm.DefineCounterMetric(metricsName)
	if err != nil {
		proxywasm.LogCriticalf("error defining metrics: %v", err)
	}
	counter = ct
	return true
}

type metricHttpContext struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	proxywasm.DefaultHttpContext
}

func newHttpContext(uint32) proxywasm.HttpContext {
	return &metricHttpContext{}
}

// override
func (ctx *metricHttpContext) OnHttpRequestHeaders(int, bool) types.Action {
	prev, err := counter.Get()
	if err != nil {
		proxywasm.LogCriticalf("error retrieving previous metric: %v", err)
	}

	proxywasm.LogInfof("previous value of %s: %d", metricsName, prev)

	if err := counter.Increment(1); err != nil {
		proxywasm.LogCriticalf("error incrementing metrics %v", err)
	}
	proxywasm.LogInfo("incremented")
	return types.ActionContinue
}
