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
	"encoding/binary"
	"errors"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const (
	sharedDataKey                 = "shared_data_key"
	sharedDataInitialValue uint64 = 10000000
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type (
	vmContext     struct{}
	pluginContext struct {
		// Embed the default plugin context here,
		// so that we don't need to reimplement all the methods.
		types.DefaultPluginContext
	}

	httpContext struct {
		// Embed the default http context here,
		// so that we don't need to reimplement all the methods.
		types.DefaultHttpContext
	}
)

// Override types.VMContext.
func (*vmContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	initialValueBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(initialValueBuf, sharedDataInitialValue)
	if err := proxywasm.SetSharedData(sharedDataKey, initialValueBuf, 0); err != nil {
		proxywasm.LogWarnf("error setting shared data on OnVMStart: %v", err)
	}
	return types.OnVMStartStatusOK
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

// Override types.DefaultPluginContext.
func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{}
}

// Override types.DefaultHttpContext.
func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	for {
		value, err := ctx.incrementData()
		if err == nil {
			proxywasm.LogInfof("shared value: %d", value)
		} else if errors.Is(err, types.ErrorStatusCasMismatch) {
			continue
		}
		break
	}
	return types.ActionContinue
}

func (ctx *httpContext) incrementData() (uint64, error) {
	value, cas, err := proxywasm.GetSharedData(sharedDataKey)
	if err != nil {
		proxywasm.LogWarnf("error getting shared data on OnHttpRequestHeaders: %v", err)
		return 0, err
	}

	buf := make([]byte, 8)
	ret := binary.LittleEndian.Uint64(value) + 1
	binary.LittleEndian.PutUint64(buf, ret)
	if err := proxywasm.SetSharedData(sharedDataKey, buf, cas); err != nil {
		proxywasm.LogWarnf("error setting shared data on OnHttpRequestHeaders: %v", err)
		return 0, err
	}
	return ret, err
}
