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

const queueName = "http_headers"

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
	return &receiverPluginContext{}
}

type receiverPluginContext struct {
	// Embed the default root context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

// Override DefaultPluginContext.
func (ctx *receiverPluginContext) OnPluginStart(vmConfigurationSize int) types.OnPluginStartStatus {
	queueID, err := proxywasm.RegisterSharedQueue(queueName)
	if err != nil {
		panic("failed register queue")
	}
	proxywasm.LogInfof("queue \"%s\" registered as id=%d", queueName, queueID)
	return types.OnPluginStartStatusOK
}

// Override DefaultPluginContext.
func (ctx *receiverPluginContext) OnQueueReady(queueID uint32) {
	data, err := proxywasm.DequeueSharedQueue(queueID)
	switch err {
	case types.ErrorStatusEmpty:
		return
	case nil:
		proxywasm.LogInfof("dequeued data: %s", string(data))
	default:
		proxywasm.LogCriticalf("error retrieving data from queue %d: %v", queueID, err)
	}
}
