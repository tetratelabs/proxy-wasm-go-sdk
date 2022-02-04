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
	"fmt"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &metricPluginContext{}
}

type metricPluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

// Override types.DefaultPluginContext.
func (ctx *metricPluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &metricHttpContext{}
}

type metricHttpContext struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
}

const (
	customHeaderKey         = "my-custom-header"
	customHeaderValueTagKey = "value"
)

var counters = map[string]proxywasm.MetricCounter{}

// Override types.DefaultHttpContext.
func (ctx *metricHttpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	customHeaderValue, err := proxywasm.GetHttpRequestHeader(customHeaderKey)
	if err == nil {
		counter, ok := counters[customHeaderValue]
		if !ok {
			// This metric is processed as: custom_header_value_counts{value="foo",reporter="wasmgosdk"} n.
			// The extraction rule is defined in envoy.yaml as a bootstrap configuration.
			// See https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/metrics/v3/stats.proto#config-metrics-v3-statsconfig.
			fqn := fmt.Sprintf("custom_header_value_counts_%s=%s_reporter=wasmgosdk", customHeaderValueTagKey, customHeaderValue)
			counter = proxywasm.DefineCounterMetric(fqn)
			counters[customHeaderValue] = counter
		}
		counter.Increment(1)
	}
	return types.ActionContinue
}
