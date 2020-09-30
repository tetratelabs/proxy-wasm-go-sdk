// Copyright 2020 Tetrate
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

func main() {
	proxywasm.SetNewRootContext(func(uint32) proxywasm.RootContext { return queue{} })
	proxywasm.SetNewHttpContext(func(uint32) proxywasm.HttpContext { return queue{} })
}

type queue struct{ proxywasm.DefaultContext }

const (
	queueName               = "proxy_wasm_go.queue"
	tickMilliseconds uint32 = 100
)

var queueID uint32

// override
func (ctx queue) OnVMStart(int) bool {
	qID, err := proxywasm.HostCallRegisterSharedQueue(queueName)
	if err != nil {
		panic(err.Error())
	}
	queueID = qID
	proxywasm.LogInfof("queue registered, name: %s, id: %d", queueName, qID)

	if err := proxywasm.HostCallSetTickPeriodMilliSeconds(tickMilliseconds); err != nil {
		proxywasm.LogCriticalf("failed to set tick period: %v", err)
	}
	proxywasm.LogInfof("set tick period milliseconds: %d", tickMilliseconds)
	return true
}

// override
func (ctx queue) OnHttpRequestHeaders(int, bool) types.Action {
	for _, msg := range []string{"hello", "world", "hello", "proxy-wasm"} {
		if err := proxywasm.HostCallEnqueueSharedQueue(queueID, []byte(msg)); err != nil {
			proxywasm.LogCriticalf("error queueing: %v", err)
		}
	}
	return types.ActionContinue
}

// override
func (ctx queue) OnTick() {
	data, err := proxywasm.HostCallDequeueSharedQueue(queueID)
	switch err {
	case types.ErrorStatusEmpty:
		return
	case nil:
		proxywasm.LogInfof("dequeued data: %s", string(data))
	default:
		proxywasm.LogCriticalf("error retrieving data from queue %d: %v", queueID, err)
	}
}
