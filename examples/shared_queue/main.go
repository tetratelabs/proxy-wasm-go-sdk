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
	queueName               = "proxy_wasm_go.queue"
	tickMilliseconds uint32 = 100
)

func main() {
	proxywasm.SetNewRootContext(newRootContext)
}

type queueRootContext struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	proxywasm.DefaultRootContext
}

func newRootContext(uint32) proxywasm.RootContext {
	return &queueRootContext{}
}

var queueID uint32

// override
func (ctx *queueRootContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	qID, err := proxywasm.RegisterSharedQueue(queueName)
	if err != nil {
		panic(err.Error())
	}
	queueID = qID
	proxywasm.LogInfof("queue registered, name: %s, id: %d", queueName, qID)

	if err := proxywasm.SetTickPeriodMilliSeconds(tickMilliseconds); err != nil {
		proxywasm.LogCriticalf("failed to set tick period: %v", err)
	}
	proxywasm.LogInfof("set tick period milliseconds: %d", tickMilliseconds)
	return types.OnVMStartStatusOK
}

// override
func (ctx *queueRootContext) OnTick() {
	// TODO: use OnQueueReady after the bug fixed is available in istio 1.8.1
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

// override
func (*queueRootContext) NewHttpContext(contextID uint32) proxywasm.HttpContext {
	return &queueHttpContext{}
}

type queueHttpContext struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	proxywasm.DefaultHttpContext
}

// override
func (ctx *queueHttpContext) OnHttpRequestHeaders(int, bool) types.Action {
	for _, msg := range []string{"hello", "world", "hello", "proxy-wasm"} {
		if err := proxywasm.EnqueueSharedQueue(queueID, []byte(msg)); err != nil {
			proxywasm.LogCriticalf("error queueing: %v", err)
		}
	}
	return types.ActionContinue
}
