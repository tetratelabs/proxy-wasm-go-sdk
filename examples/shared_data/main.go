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
	"errors"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const sharedDataKey = "shared_data_key"

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type (
	vmContext               struct{}
	sharedDataPluginContext struct {
		// Embed the default root context here,
		// so that we don't need to reimplement all the methods.
		types.DefaultPluginContext
	}

	sharedDataHttpContext struct {
		// Embed the default http context here,
		// so that we don't need to reimplement all the methods.
		types.DefaultHttpContext
	}
)

// Implement types.VMContext.
func (*vmContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	if err := proxywasm.SetSharedData(sharedDataKey, []byte{0}, 0); err != nil {
		proxywasm.LogWarnf("error setting shared data on OnVMStart: %v", err)
	}
	return types.OnVMStartStatusOK
}

// Implement types.VMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &sharedDataPluginContext{}
}

// Override DefaultPluginContext.
func (*sharedDataPluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &sharedDataHttpContext{}
}

// Override DefaultHttpContext.
func (ctx *sharedDataHttpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	for {
		value, err := ctx.incrementData()
		if err == nil {
			proxywasm.LogInfof("shared value: %d", value[0])
		} else if errors.Is(err, types.ErrorStatusCasMismatch) {
			continue
		}
		break
	}
	return types.ActionContinue
}

func (ctx *sharedDataHttpContext) incrementData() ([]byte, error) {
	value, cas, err := proxywasm.GetSharedData(sharedDataKey)
	if err != nil {
		proxywasm.LogWarnf("error getting shared data on OnHttpRequestHeaders: %v", err)
		return value, err
	}

	value[0]++
	if err := proxywasm.SetSharedData(sharedDataKey, value, cas); err != nil {
		proxywasm.LogWarnf("error setting shared data on OnHttpRequestHeaders: %v", err)
		return value, err
	}
	return value, err
}
