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
	"strconv"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetNewRootContext(func(uint32) proxywasm.RootContext { return &metrics{} })
	proxywasm.SetNewHttpContext(func(uint32) proxywasm.HttpContext { return &metrics{} })
}

var counter proxywasm.MetricCounter

const metricsName = "proxy_wasm_go.request_counter"

type metrics struct{ proxywasm.DefaultContext }

// override
func (ctx *metrics) OnVMStart(int) bool {
	ct, err := proxywasm.DefineCounterMetric(metricsName)
	if err != nil {
		proxywasm.LogCritical("error defining metrics: ", err.Error())
	}
	counter = ct
	return true
}

// override
func (ctx *metrics) OnHttpRequestHeaders(int, bool) types.Action {
	prev, err := counter.Get()
	if err != nil {
		proxywasm.LogCritical("error retrieving previous metric: ", err.Error())
	}

	proxywasm.LogInfo("previous value of ", metricsName, ": ", strconv.Itoa(int(prev)))

	if err := counter.Increment(1); err != nil {
		proxywasm.LogCritical("error incrementing metrics", err.Error())
	}
	proxywasm.LogInfo("incremented")
	return types.ActionContinue
}
