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

const (
	bufferOperationAppend  = "append"
	bufferOperationPrepend = "prepend"
	bufferOperationReplace = "replace"
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
	shouldEchoBody bool
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	if ctx.shouldEchoBody {
		return &echoBodyContext{}
	}
	return &setBodyContext{}
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	data, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
	}
	ctx.shouldEchoBody = string(data) == "echo"
	return types.OnPluginStartStatusOK
}

type setBodyContext struct {
	// Embed the default root http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	totalRequestBodySize int
	bufferOperation      string
}

// Override types.DefaultHttpContext.
func (ctx *setBodyContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	if _, err := proxywasm.GetHttpRequestHeader("content-length"); err != nil {
		if err := proxywasm.SendHttpResponse(400, nil, []byte("content must be provided"), -1); err != nil {
			panic(err)
		}
		return types.ActionPause
	}

	// Remove Content-Length in order to prevent severs from crashing if we set different body from downstream.
	if err := proxywasm.RemoveHttpRequestHeader("content-length"); err != nil {
		panic(err)
	}

	// Get "Buffer-Operation" header value.
	op, err := proxywasm.GetHttpRequestHeader("buffer-operation")
	if err != nil || (op != bufferOperationAppend &&
		op != bufferOperationPrepend &&
		op != bufferOperationReplace) {
		// Fallback to replace
		op = bufferOperationReplace
	}
	ctx.bufferOperation = op
	return types.ActionContinue
}

// Override types.DefaultHttpContext.
func (ctx *setBodyContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	ctx.totalRequestBodySize += bodySize
	if !endOfStream {
		// Wait until we see the entire body to replace.
		return types.ActionPause
	}

	originalBody, err := proxywasm.GetHttpRequestBody(0, ctx.totalRequestBodySize)
	if err != nil {
		proxywasm.LogErrorf("failed to get request body: %v", err)
		return types.ActionContinue
	}
	proxywasm.LogInfof("original request body: %s", string(originalBody))

	switch ctx.bufferOperation {
	case bufferOperationAppend:
		err = proxywasm.AppendHttpRequestBody([]byte(`[this is appended body]`))
	case bufferOperationPrepend:
		err = proxywasm.PrependHttpRequestBody([]byte(`[this is prepended body]`))
	case bufferOperationReplace:
		err = proxywasm.ReplaceHttpRequestBody([]byte(`[this is replaced body]`))
	}
	if err != nil {
		proxywasm.LogErrorf("failed to %s request body: %v", ctx.bufferOperation, err)
		return types.ActionContinue
	}
	return types.ActionContinue
}

type echoBodyContext struct {
	// mbed the default plugin context
	// so that you don't need to reimplement all the methods by yourself.
	types.DefaultHttpContext
	totalRequestBodySize int
}

// Override types.DefaultHttpContext.
func (ctx *echoBodyContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	ctx.totalRequestBodySize += bodySize
	if !endOfStream {
		// Wait until we see the entire body to replace.
		return types.ActionPause
	}

	// Send the request body as the response body.
	body, _ := proxywasm.GetHttpRequestBody(0, ctx.totalRequestBodySize)
	if err := proxywasm.SendHttpResponse(200, nil, body, -1); err != nil {
		panic(err)
	}
	return types.ActionPause
}
