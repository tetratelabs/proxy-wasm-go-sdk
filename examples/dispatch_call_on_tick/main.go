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

const tickMilliseconds uint32 = 100

func main() {
	proxywasm.SetNewRootContext(newRootContext)
}

type rootContext struct {
	// You'd better embed the default root context
	// so that you don't need to reimplement all the methods by yourself.
	proxywasm.DefaultRootContext
	contextID uint32
}

func newRootContext(contextID uint32) proxywasm.RootContext {
	return &rootContext{contextID: contextID}
}

// Override DefaultRootContext.
func (ctx *rootContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	if err := proxywasm.SetTickPeriodMilliSeconds(tickMilliseconds); err != nil {
		proxywasm.LogCriticalf("failed to set tick period: %v", err)
		return types.OnVMStartStatusFailed
	}
	proxywasm.LogInfof("set tick period milliseconds: %d", tickMilliseconds)
	return types.OnVMStartStatusOK
}

// Override DefaultRootContext.
func (ctx *rootContext) OnTick() {
	hs := types.Headers{
		{":method", "GET"}, {":authority", "some_authority"}, {":path", "/path/to/service"}, {"accept", "*/*"},
	}
	if _, err := proxywasm.DispatchHttpCall("web_service", hs, "", nil,
		5000, callback); err != nil {
		proxywasm.LogCriticalf("dispatch httpcall failed: %v", err)
	}
}

var cnt int

func callback(numHeaders, bodySize, numTrailers int) {
	cnt++
	proxywasm.LogInfof("called! %d", cnt)
}
