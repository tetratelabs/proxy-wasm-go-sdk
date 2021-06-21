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

type vmContext struct{}

// Implement types.VMContext.
func (*vmContext) OnVMStart(int) types.OnVMStartStatus {
	return types.OnVMStartStatusOK
}

// Implement types.VMContext.
func (*vmContext) NewPluginContext(uint32) types.PluginContext {
	return &senderPluginContext{}
}

type senderPluginContext struct {
	// Embed the default root context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

func newPluginContext(uint32) types.PluginContext {
	return &senderPluginContext{}
}

// Override DefaultPluginContext.
func (ctx *senderPluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	queueID, err := proxywasm.ResolveSharedQueue(receiverVMID, queueName)
	if err != nil {
		proxywasm.LogCriticalf("error resolving queue id: %v", err)
	}

	// Pass the resolved queueID to http contexts so they can enqueue.
	return &senderHttpContext{queueID: queueID}
}

type senderHttpContext struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	queueID uint32
}

// Override DefaultHttpContext.
func (ctx *senderHttpContext) OnHttpRequestHeaders(int, bool) types.Action {
	headers, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("error getting request headers: %v", err)
	}
	for _, h := range headers {
		msg := fmt.Sprintf("{\"key\": \"%s\",\"value\": \"%s\"}", h[0], h[1])
		if err := proxywasm.EnqueueSharedQueue(ctx.queueID, []byte(msg)); err != nil {
			proxywasm.LogCriticalf("error queueing: %v", err)
		} else {
			proxywasm.LogInfof("enqueued data: %s", msg)
		}
	}
	return types.ActionContinue
}
