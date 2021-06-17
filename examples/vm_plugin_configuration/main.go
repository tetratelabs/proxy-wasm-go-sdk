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

func main() {
	proxywasm.SetNewRootContextFn(newRootContext)
}

type context struct {
	// You'd better embed the default root context
	// so that you don't need to reimplement all the methods by yourself.
	types.DefaultRootContext
}

func newRootContext(contextID uint32) types.RootContext {
	return &context{}
}

// Override DefaultRootContext.
func (ctx context) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	data, err := proxywasm.GetVMConfiguration(vmConfigurationSize)
	if err != nil {
		proxywasm.LogCriticalf("error reading vm configuration: %v", err)
	}

	proxywasm.LogInfof("vm config: %s", string(data))
	return types.OnVMStartStatusOK
}

// Override DefaultRootContext.
func (ctx context) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	data, err := proxywasm.GetPluginConfiguration(pluginConfigurationSize)
	if err != nil {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
	}

	proxywasm.LogInfof("plugin config: %s", string(data))
	return types.OnPluginStartStatusOK
}
