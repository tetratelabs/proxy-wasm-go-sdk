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
	proxywasm.SetNewRootContextFn(newRootContext)
}

type rootContext struct {
	// Embed the default root context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultRootContext
}

func newRootContext(uint32) types.RootContext { return &rootContext{} }

// Override DefaultRootContext.
func (*rootContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpHeaders{contextID: contextID}
}

type httpHeaders struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID uint32
}

// Override DefaultHttpContext.
func (ctx *httpHeaders) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	err := proxywasm.SetHttpRequestHeader("test", "best")
	if err != nil {
		proxywasm.LogCritical("failed to set request header: test")
	}

	hs, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
	}

	for _, h := range hs {
		proxywasm.LogInfof("request header --> %s: %s", h[0], h[1])
	}
	return types.ActionContinue
}

// Override DefaultHttpContext.
func (ctx *httpHeaders) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	hs, err := proxywasm.GetHttpResponseHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get response headers: %v", err)
	}

	for _, h := range hs {
		proxywasm.LogInfof("response header <-- %s: %s", h[0], h[1])
	}
	return types.ActionContinue
}

// Override DefaultHttpContext.
func (ctx *httpHeaders) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}
