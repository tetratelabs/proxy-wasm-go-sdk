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

type rootContext struct {
	// Embed the default root context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultRootContext
	config []byte
}

func newRootContext(contextID uint32) types.RootContext {
	return &rootContext{}
}

// Override DefaultRootContext.
func (ctx *rootContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	data, err := proxywasm.GetPluginConfiguration(pluginConfigurationSize)
	if err != nil {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
	}

	ctx.config = data
	return types.OnPluginStartStatusOK
}

// Override DefaultRootContext.
func (ctx *rootContext) NewHttpContext(contextID uint32) types.HttpContext {
	ret := &httpContext{config: ctx.config}
	proxywasm.LogInfof("read plugin config from root context: %s", string(ret.config))
	return ret
}

type httpContext struct {
	types.DefaultHttpContext
	config []byte
}
