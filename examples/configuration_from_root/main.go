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
	proxywasm.SetNewRootContext(newRootContext)
	proxywasm.SetNewHttpContext(newHttpContext)
}

type rootContext struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	proxywasm.DefaultRootContext

	config []byte
}

func newRootContext(contextID uint32) proxywasm.RootContext {
	return &rootContext{}
}

func (ctx *rootContext) OnPluginStart(pluginConfigurationSize int) bool {
	data, err := proxywasm.GetPluginConfiguration(pluginConfigurationSize)
	if err != nil {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
	}

	ctx.config = data
	return true
}

type httpContext struct {
	proxywasm.DefaultHttpContext

	config []byte
}

func newHttpContext(rootContextID, contextID uint32) proxywasm.HttpContext {
	ctx := &httpContext{}

	rootCtx, err := proxywasm.GetRootContextByID(rootContextID)
	if err != nil {
		proxywasm.LogErrorf("unable to get root context: %v", err)

		return ctx
	}

	exampleRootCtx, ok := rootCtx.(*rootContext)
	if !ok {
		proxywasm.LogError("could not cast root context")
	}

	ctx.config = exampleRootCtx.config

	proxywasm.LogInfof("plugin config: %s\n", string(ctx.config))
	return ctx
}
