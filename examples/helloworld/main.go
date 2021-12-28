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
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &helloWorld{}
}

type helloWorld struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
	contextID uint32
}

// Override types.DefaultPluginContext.
func (ctx *helloWorld) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	rand.Seed(time.Now().UnixNano())

	proxywasm.LogInfo("OnPluginStart from Go!")
	if err := proxywasm.SetTickPeriodMilliSeconds(tickMilliseconds); err != nil {
		proxywasm.LogCriticalf("failed to set tick period: %v", err)
	}

	return types.OnPluginStartStatusOK
}

// Override types.DefaultPluginContext.
func (ctx *helloWorld) OnTick() {
	t := time.Now().UnixNano()
	proxywasm.LogInfof("It's %d: random value: %d", t, rand.Uint64())
	proxywasm.LogInfof("OnTick called")
}

// Override types.DefaultPluginContext.
func (*helloWorld) NewHttpContext(uint32) types.HttpContext { return &types.DefaultHttpContext{} }
