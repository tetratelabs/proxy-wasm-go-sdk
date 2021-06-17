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
	"math/rand"
	"time"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const tickMilliseconds uint32 = 1000

func main() {
	proxywasm.SetNewRootContextFn(newHelloWorld)
}

type helloWorld struct {
	// You'd better embed the default root context
	// so that you don't need to reimplement all the methods by yourself.
	types.DefaultRootContext
	contextID uint32
}

func newHelloWorld(contextID uint32) types.RootContext {
	return &helloWorld{contextID: contextID}
}

// Override DefaultRootContext.
func (ctx *helloWorld) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	rand.Seed(time.Now().UnixNano())

	proxywasm.LogInfo("proxy_on_vm_start from Go!")
	if err := proxywasm.SetTickPeriodMilliSeconds(tickMilliseconds); err != nil {
		proxywasm.LogCriticalf("failed to set tick period: %v", err)
	}

	return types.OnVMStartStatusOK
}

// Override DefaultRootContext.
func (ctx *helloWorld) OnTick() {
	t := time.Now().UnixNano()
	proxywasm.LogInfof("It's %d: random value: %d", t, rand.Uint64())
	proxywasm.LogInfof("OnTick called")
}
