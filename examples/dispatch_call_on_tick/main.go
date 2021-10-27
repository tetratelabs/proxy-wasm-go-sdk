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
	"crypto/rand"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const tickMilliseconds uint32 = 100

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
	return &pluginContext{contextID: contextID}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
	contextID uint32
	callBack  func(numHeaders, bodySize, numTrailers int)
	cnt       int
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	if err := proxywasm.SetTickPeriodMilliSeconds(tickMilliseconds); err != nil {
		proxywasm.LogCriticalf("failed to set tick period: %v", err)
		return types.OnPluginStartStatusFailed
	}
	proxywasm.LogInfof("set tick period milliseconds: %d", tickMilliseconds)
	ctx.callBack = func(numHeaders, bodySize, numTrailers int) {
		ctx.cnt++
		proxywasm.LogInfof("called %d for contextID=%d", ctx.cnt, ctx.contextID)
		headers, err := proxywasm.GetHttpCallResponseHeaders()
		if err != nil && err != types.ErrorStatusNotFound {
			panic(err)
		}
		for _, h := range headers {
			proxywasm.LogInfof("response header for the dispatched call: %s: %s", h[0], h[1])
		}
		headers, err = proxywasm.GetHttpCallResponseTrailers()
		if err != nil && err != types.ErrorStatusNotFound {
			panic(err)
		}
		for _, h := range headers {
			proxywasm.LogInfof("response trailer for the dispatched call: %s: %s", h[0], h[1])
		}
	}
	return types.OnPluginStartStatusOK
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) OnTick() {
	headers := [][2]string{
		{":method", "GET"}, {":authority", "some_authority"}, {"accept", "*/*"},
	}
	// Pick random value to select the request path.
	buf := make([]byte, 1)
	_, _ = rand.Read(buf)
	if buf[0]%2 == 0 {
		headers = append(headers, [2]string{":path", "/ok"})
	} else {
		headers = append(headers, [2]string{":path", "/fail"})
	}
	if _, err := proxywasm.DispatchHttpCall("web_service", headers, nil, nil, 5000, ctx.callBack); err != nil {
		proxywasm.LogCriticalf("dispatch httpcall failed: %v", err)
	}
}
