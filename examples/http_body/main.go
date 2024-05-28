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

// vmContext implements types.VMContext.
type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// NewPluginContext implements types.VMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

// pluginContext implements types.PluginContext.
type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
	shouldEchoBody bool
}

// NewHttpContext implements types.PluginContext.
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	if ctx.shouldEchoBody {
		return &echoBodyContext{}
	}
	return &setBodyContext{}
}

// OnPluginStart implements types.PluginContext.
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	data, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
	}
	ctx.shouldEchoBody = string(data) == "echo"
	return types.OnPluginStartStatusOK
}

// setBodyContext implements types.HttpContext.
type setBodyContext struct {
	// Embed the default root http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	modifyResponse  bool
	bufferOperation string
}

// OnHttpRequestHeaders implements types.HttpContext.
func (ctx *setBodyContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	mode, err := proxywasm.GetHttpRequestHeader("buffer-replace-at")
	if err == nil && mode == "response" {
		ctx.modifyResponse = true
	}

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

// OnHttpRequestBody implements types.HttpContext.
func (ctx *setBodyContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	if ctx.modifyResponse {
		return types.ActionContinue
	}

	if !endOfStream {
		// Wait until we see the entire body to replace.
		return types.ActionPause
	}
	// Being the body never been sent upstream so far, bodySize is the total size of the body received.
	originalBody, err := proxywasm.GetHttpRequestBody(0, bodySize)
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

// OnHttpResponseHeaders implements types.HttpContext.
func (ctx *setBodyContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	if !ctx.modifyResponse {
		return types.ActionContinue
	}

	// Remove Content-Length in order to prevent severs from crashing if we set different body.
	if err := proxywasm.RemoveHttpResponseHeader("content-length"); err != nil {
		panic(err)
	}

	return types.ActionContinue
}

// OnHttpResponseBody implements types.HttpContext.
func (ctx *setBodyContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {
	if !ctx.modifyResponse {
		return types.ActionContinue
	}

	if !endOfStream {
		// Wait until we see the entire body to replace.
		return types.ActionPause
	}

	originalBody, err := proxywasm.GetHttpResponseBody(0, bodySize)
	if err != nil {
		proxywasm.LogErrorf("failed to get response body: %v", err)
		return types.ActionContinue
	}
	proxywasm.LogInfof("original response body: %s", string(originalBody))

	switch ctx.bufferOperation {
	case bufferOperationAppend:
		err = proxywasm.AppendHttpResponseBody([]byte(`[this is appended body]`))
	case bufferOperationPrepend:
		err = proxywasm.PrependHttpResponseBody([]byte(`[this is prepended body]`))
	case bufferOperationReplace:
		err = proxywasm.ReplaceHttpResponseBody([]byte(`[this is replaced body]`))
	}
	if err != nil {
		proxywasm.LogErrorf("failed to %s response body: %v", ctx.bufferOperation, err)
		return types.ActionContinue
	}
	return types.ActionContinue
}

// echoBodyContext implements types.HttpContext.
type echoBodyContext struct {
	// Embed the default plugin context
	// so that you don't need to reimplement all the methods by yourself.
	types.DefaultHttpContext
	totalRequestBodySize int
}

// OnHttpRequestBody implements types.HttpContext.
func (ctx *echoBodyContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	ctx.totalRequestBodySize = bodySize
	if !endOfStream {
		// Wait until we see the entire body to replace.
		return types.ActionPause
	}
	// Send the request body as the response body.
	body, _ := proxywasm.GetHttpRequestBody(0, bodySize)
	if err := proxywasm.SendHttpResponse(200, nil, body, -1); err != nil {
		panic(err)
	}
	return types.ActionPause
}
