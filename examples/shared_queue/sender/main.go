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

const (
	receiverVMID = "receiver"
	queueName    = "http_headers"
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
func (*vmContext) NewPluginContext(uint32) types.PluginContext {
	return &senderPluginContext{}
}

type senderPluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

func newPluginContext(uint32) types.PluginContext {
	return &senderPluginContext{}
}

// Override types.DefaultPluginContext.
func (ctx *senderPluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	requestHeadersQueueID, err := proxywasm.ResolveSharedQueue(receiverVMID, "http_request_headers")
	if err != nil {
		proxywasm.LogCriticalf("error resolving queue id: %v", err)
	}

	responseHeadersQueueID, err := proxywasm.ResolveSharedQueue(receiverVMID, "http_response_headers")
	if err != nil {
		proxywasm.LogCriticalf("error resolving queue id: %v", err)
	}

	// Pass the resolved queueID to http contexts so they can enqueue.
	return &senderHttpContext{
		requestHeadersQueueID:  requestHeadersQueueID,
		responseHeadersQueueID: responseHeadersQueueID,
	}
}

type senderHttpContext struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	requestHeadersQueueID, responseHeadersQueueID uint32
}

// Override types.DefaultHttpContext.
func (ctx *senderHttpContext) OnHttpRequestHeaders(int, bool) types.Action {
	headers, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("error getting request headers: %v", err)
	}
	for _, h := range headers {
		msg := fmt.Sprintf("{\"key\": \"%s\",\"value\": \"%s\"}", h[0], h[1])
		if err := proxywasm.EnqueueSharedQueue(ctx.requestHeadersQueueID, []byte(msg)); err != nil {
			proxywasm.LogCriticalf("error queueing: %v", err)
		} else {
			proxywasm.LogInfof("enqueued data: %s", msg)
		}
	}
	return types.ActionContinue
}

// Override types.DefaultHttpContext.
func (ctx *senderHttpContext) OnHttpResponseHeaders(int, bool) types.Action {
	headers, err := proxywasm.GetHttpResponseHeaders()
	if err != nil {
		proxywasm.LogCriticalf("error getting response headers: %v", err)
	}
	for _, h := range headers {
		msg := fmt.Sprintf("{\"key\": \"%s\",\"value\": \"%s\"}", h[0], h[1])
		if err := proxywasm.EnqueueSharedQueue(ctx.responseHeadersQueueID, []byte(msg)); err != nil {
			proxywasm.LogCriticalf("error queueing: %v", err)
		} else {
			proxywasm.LogInfof("enqueued data: %s", msg)
		}
	}
	return types.ActionContinue
}
