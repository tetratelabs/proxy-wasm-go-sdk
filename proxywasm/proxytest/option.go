// Copyright 2021 Tetrate
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

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

// EmulatorOption is an option that can be passed to NewHostEmulator.
type EmulatorOption struct {
	pluginConfiguration []byte
	vmConfiguration     []byte
	vmContext           types.VMContext
	properties          map[string][]byte
}

// NewEmulatorOption creates a new EmulatorOption.
func NewEmulatorOption() *EmulatorOption {
	return &EmulatorOption{vmContext: &types.DefaultVMContext{}}
}

// WithVMContext sets the VMContext.
func (o *EmulatorOption) WithVMContext(context types.VMContext) *EmulatorOption {
	o.vmContext = context
	return o
}

// WithPluginConfiguration sets the plugin configuration.
func (o *EmulatorOption) WithPluginConfiguration(data []byte) *EmulatorOption {
	o.pluginConfiguration = data
	return o
}

// WithVMConfiguration sets the VM configuration.
func (o *EmulatorOption) WithVMConfiguration(data []byte) *EmulatorOption {
	o.vmConfiguration = data
	return o
}

// WithProperty sets a property. If the property already exists, it will be overwritten.
func (o *EmulatorOption) WithProperty(path []string, value []byte) *EmulatorOption {
	if o.properties == nil {
		o.properties = map[string][]byte{}
	}
	o.properties[string(internal.SerializePropertyPath(path))] = value
	return o
}
