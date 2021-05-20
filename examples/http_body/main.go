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

func main() {
	proxywasm.SetNewRootContext(newContext)
}
func newContext(uint32) proxywasm.RootContext { return &rootContext{} }

type rootContext struct {
	// You'd better embed the default root context
	// so that you don't need to reimplement all the methods by yourself.
	proxywasm.DefaultRootContext
	shouldEchoBody bool
}

// Override DefaultRootContext.
func (ctx *rootContext) NewHttpContext(contextID uint32) proxywasm.HttpContext {
	if ctx.shouldEchoBody {
		return &echoBodyContext{}
	}
	return &setBodyContext{}
}

// Override DefaultRootContext.
func (ctx *rootContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	data, err := proxywasm.GetPluginConfiguration(pluginConfigurationSize)
	if err != nil {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
	}
	ctx.shouldEchoBody = string(data) == "echo"
	return types.OnPluginStartStatusOK
}

type setBodyContext struct {
	// You'd better embed the default root context
	// so that you don't need to reimplement all the methods by yourself.
	proxywasm.DefaultHttpContext
}

// Override DefaultHttpContext.
func (ctx *setBodyContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	proxywasm.LogInfof("body size: %d", bodySize)
	if bodySize != 0 {
		initialBody, err := proxywasm.GetHttpRequestBody(0, bodySize)
		if err != nil {
			proxywasm.LogErrorf("failed to get request body: %v", err)
			return types.ActionContinue
		}
		proxywasm.LogInfof("initial request body: %s", string(initialBody))

		b := []byte(`{ "another": "body" }`)

		err = proxywasm.SetHttpRequestBody(b)
		if err != nil {
			proxywasm.LogErrorf("failed to set request body: %v", err)
			return types.ActionContinue
		}

		proxywasm.LogInfof("on http request body finished")
	}

	return types.ActionContinue
}

type echoBodyContext struct {
	// You'd better embed the default root context
	// so that you don't need to reimplement all the methods by yourself.
	proxywasm.DefaultHttpContext
}

func (ctx *echoBodyContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	var reject bool
	clen, err := proxywasm.GetHttpRequestHeader("content-length")
	if err != nil {
		reject = true
	} else if l, err := strconv.Atoi(clen); err != nil || l < 1 {
		reject = true
	}

	if reject {
		if err := proxywasm.SendHttpResponse(400, nil, []byte("content must be provided")); err != nil {
			panic(err)
		}
		return types.ActionPause
	}
	return types.ActionContinue
}

// Override DefaultHttpContext.
func (ctx *echoBodyContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	if bodySize != 0 {
		body, _ := proxywasm.GetHttpRequestBody(0, bodySize)
		if err := proxywasm.SendHttpResponse(200, nil, body); err != nil {
			panic(err)
		}
		return types.ActionPause
	}
	return types.ActionPause
}
