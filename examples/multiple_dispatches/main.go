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
	"strconv"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const clusterName = "httpbin"

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
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

// Override types.DefaultPluginContext.
func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{contextID: contextID}
}

type httpContext struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	// contextID is the unique identifier assigned to each httpContext.
	contextID uint32
	// pendingDispatchedRequest is the number of pending dispatched requests.
	pendingDispatchedRequest int
}

const totalDispatchNum = 10

// Override types.DefaultHttpContext.
func (ctx *httpContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	// On each request response, we dispatch the http calls `totalDispatchNum` times.
	// Note: DispatchHttpCall is asynchronously processed, so each loop is non-blocking.
	for i := 0; i < totalDispatchNum; i++ {
		if _, err := proxywasm.DispatchHttpCall(clusterName, [][2]string{
			{":path", "/"},
			{":method", "GET"},
			{":authority", ""}},
			nil, nil, 50000, ctx.dispatchCallback); err != nil {
			panic(err)
		}
		// Now we have made a dispatched request, so we record it.
		ctx.pendingDispatchedRequest++
	}
	return types.ActionPause
}

// dispatchCallback is the callback function called in response to the response arrival from the dispatched request.
func (ctx *httpContext) dispatchCallback(numHeaders, bodySize, numTrailers int) {
	// Decrement the pending request counter.
	ctx.pendingDispatchedRequest--
	if ctx.pendingDispatchedRequest == 0 {
		// This case, all the dispatched request was processed.
		// Adds a response header to the original response.
		proxywasm.AddHttpResponseHeader("total-dispatched", strconv.Itoa(totalDispatchNum))
		// And then contniue the original reponse.
		proxywasm.ResumeHttpResponse()
		proxywasm.LogInfof("response resumed after processed %d dispatched request", totalDispatchNum)
	} else {
		proxywasm.LogInfof("pending dispatched requests: %d", ctx.pendingDispatchedRequest)
	}
}
