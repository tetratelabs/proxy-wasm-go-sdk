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
	proxywasm.SetNewRootContextFn(newRootContext)
}

type receiverRootContext struct {
	// You'd better embed the default root context
	// so that you don't need to reimplement all the methods by yourself.
	types.DefaultRootContext
}

func newRootContext(uint32) types.RootContext {
	return &receiverRootContext{}
}

// Override DefaultRootContext.
func (ctx *receiverRootContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	queueID, err := proxywasm.RegisterSharedQueue(queueName)
	if err != nil {
		panic("failed register queue")
	}
	proxywasm.LogInfof("queue \"%s\" registered as id=%d", queueName, queueID)
	return types.OnVMStartStatusOK
}

// Override DefaultRootContext.
func (ctx *receiverRootContext) OnQueueReady(queueID uint32) {
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
