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
	"fmt"

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
	return &metricPluginContext{}
}

// metricPluginContext implements types.PluginContext.
type metricPluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

// NewHttpContext implements types.PluginContext.
func (ctx *metricPluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &metricHttpContext{}
}

// metricHttpContext implements types.HttpContext.
type metricHttpContext struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
}

const (
	customHeaderKey         = "my-custom-header"
	customHeaderValueTagKey = "value"
)

// counters is a map from custom header value to a counter metric.
// Note that Proxy-Wasm plugins are single threaded, so no need to use a lock.
var counters = map[string]proxywasm.MetricCounter{}

// OnHttpRequestHeaders implements types.HttpContext.
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
