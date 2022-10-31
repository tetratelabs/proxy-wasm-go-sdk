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
	return &properties{contextID: contextID}
}

type properties struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID uint32
}

var propertyPrefix = []string{
	"route_metadata",
	"filter_metadata",
	"envoy.filters.http.wasm",
}

// Override types.DefaultHttpContext.
func (ctx *properties) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	auth, err := proxywasm.GetProperty(append(propertyPrefix, "auth"))
	if err != nil {
		if err == types.ErrorStatusNotFound {
			proxywasm.LogInfo("no auth header for route")
			return types.ActionContinue
		}
		proxywasm.LogCriticalf("failed to read properties: %v", err)
	}
	proxywasm.LogInfof("auth header is \"%s\"", auth)

	hs, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
	}

	// Verify authentication header exists
	authHeader := false
	for _, h := range hs {
		if h[0] == string(auth) {
			authHeader = true
			break
		}
	}

	// Reject requests without authentication header
	if !authHeader {
		proxywasm.SendHttpResponse(401, nil, nil, 16)
		return types.ActionPause
	}

	return types.ActionContinue
}

// Override types.DefaultHttpContext.
func (ctx *properties) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}
