// Copyright 2021 Tetratea
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxytest

import "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"

type EmulatorOption struct {
	pluginConfiguration, vmConfiguration []byte
	newRootContext                       func(uint32) proxywasm.RootContext
	newStreamContext                     func(uint32, uint32) proxywasm.StreamContext
	newHttpContext                       func(uint32, uint32) proxywasm.HttpContext
}

func NewEmulatorOption() *EmulatorOption {
	return &EmulatorOption{}
}

func (o *EmulatorOption) WithNewRootContext(f func(uint32) proxywasm.RootContext) *EmulatorOption {
	o.newRootContext = f
	return o
}

func (o *EmulatorOption) WithNewHttpContext(f func(uint32, uint32) proxywasm.HttpContext) *EmulatorOption {
	o.newHttpContext = f
	return o
}

func (o *EmulatorOption) WithNewStreamContext(f func(uint32, uint32) proxywasm.StreamContext) *EmulatorOption {
	o.newStreamContext = f
	return o
}

func (o *EmulatorOption) WithPluginConfiguration(data []byte) *EmulatorOption {
	o.pluginConfiguration = data
	return o
}

func (o *EmulatorOption) WithVMConfiguration(data []byte) *EmulatorOption {
	o.vmConfiguration = data
	return o
}
