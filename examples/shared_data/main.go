// Copyright 2020-2024 Tetrate
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
	sharedDataKey = "shared_data_key"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type (
	// vmContext implements types.VMContext.
	vmContext struct{}
	// pluginContext implements types.PluginContext.
	pluginContext struct {
		// Embed the default plugin context here,
		// so that we don't need to reimplement all the methods.
		types.DefaultPluginContext
	}
	// httpContext implements types.HttpContext.
	httpContext struct {
		// Embed the default http context here,
		// so that we don't need to reimplement all the methods.
		types.DefaultHttpContext
	}
)

// OnVMStart implements types.VMContext.
func (*vmContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	initialValueBuf := make([]byte, 0) // Empty data to indicate that the data is not initialized.
	if err := proxywasm.SetSharedData(sharedDataKey, initialValueBuf, 0); err != nil {
		proxywasm.LogWarnf("error setting shared data on OnVMStart: %v", err)
	}
	return types.OnVMStartStatusOK
}

// NewPluginContext implements types.VMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

// NewHttpContext implements types.PluginContext.
func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{}
}

// OnHttpRequestHeaders implements types.HttpContext.
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

// incrementData increments the shared data value by 1.
func (ctx *httpContext) incrementData() (uint64, error) {
	data, cas, err := proxywasm.GetSharedData(sharedDataKey)
	if err != nil {
		proxywasm.LogWarnf("error getting shared data on OnHttpRequestHeaders: %v", err)
		return 0, err
	}

	var nextValue uint64
	if len(data) > 0 {
		nextValue = binary.LittleEndian.Uint64(data) + 1
	} else {
		nextValue = 1
	}

	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, nextValue)
	if err := proxywasm.SetSharedData(sharedDataKey, buf, cas); err != nil {
		proxywasm.LogWarnf("error setting shared data on OnHttpRequestHeaders: %v", err)
		return 0, err
	}
	return nextValue, err
}
