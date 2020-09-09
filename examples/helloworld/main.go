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
	"strconv"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

func main() {
	proxywasm.SetNewRootContext(newHelloWorld)
}

type helloWorld struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	proxywasm.DefaultContext
	contextID uint32
}

func newHelloWorld(contextID uint32) proxywasm.RootContext {
	return &helloWorld{contextID: contextID}
}

// override
func (ctx *helloWorld) OnVMStart(int) bool {
	proxywasm.LogInfo("proxy_on_vm_start from Go!")
	if err := proxywasm.HostCallSetTickPeriodMilliSeconds(1000); err != nil {
		proxywasm.LogCritical("failed to set tick period: ", err.Error())
	}
	return true
}

// override
func (ctx *helloWorld) OnTick() {
	t := proxywasm.HostCallGetCurrentTime()
	proxywasm.LogInfo("OnTick on ", strconv.FormatUint(uint64(ctx.contextID), 10),
		", it's ", strconv.FormatInt(t, 10))
}
