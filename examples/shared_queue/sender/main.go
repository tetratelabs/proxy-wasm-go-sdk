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
	"encoding/hex"
	"fmt"
	"hash/fnv"

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
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &senderPluginContext{contextID: contextID}
}

type senderPluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
	config    string
	contextID uint32
}

func newPluginContext(uint32) types.PluginContext {
	return &senderPluginContext{}
}

// Override types.DefaultPluginContext.
func (ctx *senderPluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	// Get Plugin configuration.
	config, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		panic(fmt.Sprintf("failed to get plugin config: %v", err))
	}
	ctx.config = string(config)
	proxywasm.LogInfof("contextID=%d is configured for %s", ctx.contextID, ctx.config)
	return types.OnPluginStartStatusOK
}

// Override types.DefaultPluginContext.
func (ctx *senderPluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	// If this PluginContext is not configured for Http, then return nil.
	if ctx.config != "http" {
		return nil
	}

	// Resolve queues.
	requestHeadersQueueID, err := proxywasm.ResolveSharedQueue(receiverVMID, "http_request_headers")
	if err != nil {
		proxywasm.LogCriticalf("error resolving queue id: %v", err)
	}

	responseHeadersQueueID, err := proxywasm.ResolveSharedQueue(receiverVMID, "http_response_headers")
	if err != nil {
		proxywasm.LogCriticalf("error resolving queue id: %v", err)
	}

	// Pass the resolved queueIDs to http contexts so they can enqueue.
	return &senderHttpContext{
		requestHeadersQueueID:  requestHeadersQueueID,
		responseHeadersQueueID: responseHeadersQueueID,
		contextID:              contextID,
	}
}

type senderHttpContext struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID, requestHeadersQueueID, responseHeadersQueueID uint32
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
			proxywasm.LogInfof("(contextID=%d) enqueued data: %s", ctx.contextID, msg)
		}
	}
	return types.ActionContinue
}

func (ctx *senderPluginContext) NewTcpContext(contextID uint32) types.TcpContext {
	// If this PluginContext is not configured for Tcp, then return nil.
	if ctx.config != "tcp" {
		return nil
	}

	// Resolve queue.
	queueID, err := proxywasm.ResolveSharedQueue(receiverVMID, "tcp_data_hashes")
	if err != nil {
		proxywasm.LogCriticalf("error resolving queue id: %v", err)
	}

	// Pass the resolved queueID to tcp contexts so they can enqueue.
	return &senderTcpContext{
		tcpHashesQueueID: queueID,
		contextID:        contextID,
	}
}

type senderTcpContext struct {
	types.DefaultTcpContext
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	tcpHashesQueueID uint32
	contextID        uint32
}

func (ctx *senderTcpContext) OnUpstreamData(dataSize int, endOfStream bool) types.Action {
	if dataSize == 0 {
		return types.ActionContinue
	}

	// Calculate the hash of the data frame.
	data, err := proxywasm.GetUpstreamData(0, dataSize)
	if err != nil && err != types.ErrorStatusNotFound {
		proxywasm.LogCritical(err.Error())
	}
	s := fnv.New128a()
	_, _ = s.Write(data)
	var buf []byte
	buf = s.Sum(buf)
	hash := hex.EncodeToString(buf)

	// Enqueue the hashed data frame.
	if err := proxywasm.EnqueueSharedQueue(ctx.tcpHashesQueueID, []byte(hash)); err != nil {
		proxywasm.LogCriticalf("error queueing: %v", err)
	} else {
		proxywasm.LogInfof("(contextID=%d) enqueued data: %s", ctx.contextID, hash)
	}
	return types.ActionContinue
}
