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
)

func main() {
	proxywasm.SetNewRootContext(func(uint32) proxywasm.RootContext { return context{} })
}

type context struct{ proxywasm.DefaultContext }

// override
func (ctx context) OnVMStart(vmConfigurationSize int) bool {
	data, err := proxywasm.HostCallGetVMConfiguration(vmConfigurationSize)
	if err != nil {
		proxywasm.LogCritical("error reading vm configuration", err.Error())
	}

	proxywasm.LogInfo("vm config: \n", string(data))
	return true
}

func (ctx context) OnConfigure(pluginConfigurationSize int) bool {
	data, err := proxywasm.HostCallGetPluginConfiguration(pluginConfigurationSize)
	if err != nil {
		proxywasm.LogCritical("error reading plugin configuration", err.Error())
	}

	proxywasm.LogInfo("plugin config: \n", string(data))
	return true
}
