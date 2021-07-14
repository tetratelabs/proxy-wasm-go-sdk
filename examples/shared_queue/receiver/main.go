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
	return &receiverPluginContext{contextID: contextID}
}

type receiverPluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	contextID uint32
	types.DefaultPluginContext
	queueName string
}

// Override types.DefaultPluginContext.
func (ctx *receiverPluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	// Get Plugin configuration.
	config, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		panic(fmt.Sprintf("failed to get plugin config: %v", err))
	}

	// Treat the config as the queue name for receiving.
	ctx.queueName = string(config)

	queueID, err := proxywasm.RegisterSharedQueue(ctx.queueName)
	if err != nil {
		panic("failed register queue")
	}
	proxywasm.LogInfof("queue \"%s\" registered as queueID=%d by contextID=%d", ctx.queueName, queueID, ctx.contextID)
	return types.OnPluginStartStatusOK
}

// Override types.DefaultPluginContext.
func (ctx *receiverPluginContext) OnQueueReady(queueID uint32) {
	data, err := proxywasm.DequeueSharedQueue(queueID)
	switch err {
	case types.ErrorStatusEmpty:
		return
	case nil:
		proxywasm.LogInfof("(contextID=%d) dequeued data from %s(queueID=%d): %s", ctx.contextID, ctx.queueName, queueID, string(data))
	default:
		proxywasm.LogCriticalf("error retrieving data from queue %d: %v", queueID, err)
	}
}
