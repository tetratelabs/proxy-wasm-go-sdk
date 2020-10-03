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
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetNewRootContext(newRootContext)
	proxywasm.SetNewHttpContext(newHttpContext)
}

type (
	sharedDataRootContext struct {
		// you must embed the default context so that you need not to reimplement all the methods by yourself
		proxywasm.DefaultRootContext
	}

	sharedDataHttpContext struct {
		// you must embed the default context so that you need not to reimplement all the methods by yourself
		proxywasm.DefaultHttpContext
	}
)

func newRootContext(uint32) proxywasm.RootContext {
	return &sharedDataRootContext{}
}

func newHttpContext(uint32) proxywasm.HttpContext {
	return &sharedDataHttpContext{}
}

const sharedDataKey = "shared_data_key"

// override
func (ctx *sharedDataRootContext) OnVMStart(int) bool {
	if err := proxywasm.SetSharedData(sharedDataKey, []byte{0}, 0); err != nil {
		proxywasm.LogWarnf("error setting shared data on OnVMStart: %v", err)
	}
	return true
}

// override
func (ctx *sharedDataHttpContext) OnHttpRequestHeaders(int, bool) types.Action {
	value, cas, err := proxywasm.GetSharedData(sharedDataKey)
	if err != nil {
		proxywasm.LogWarnf("error getting shared data on OnHttpRequestHeaders: %v", err)
		return types.ActionContinue
	}

	value[0]++
	if err := proxywasm.SetSharedData(sharedDataKey, value, cas); err != nil {
		proxywasm.LogWarnf("error setting shared data on OnHttpRequestHeaders: %v", err)
		return types.ActionContinue
	}

	proxywasm.LogInfof("shared value: %d", value[0])
	return types.ActionContinue
}
