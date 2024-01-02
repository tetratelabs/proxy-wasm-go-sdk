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
	"fmt"
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &setBodyContext{}
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	return types.OnPluginStartStatusOK
}

type setBodyContext struct {
	// Embed the default root http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	totalRequestBodyReadSize int
	receivedChunks           int
}

// Override types.DefaultHttpContext.
func (ctx *setBodyContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	proxywasm.LogInfof("OnHttpRequestBody called. BodySize: %d, totalRequestBodyReadSize: %d, endOfStream: %v", bodySize, ctx.totalRequestBodyReadSize, endOfStream)

	// If some data has been received, we read it.
	// Reading the body chunk by chunk, bodySize is the size of the current chunk, not the total size of the body.
	chunkSize := bodySize - ctx.totalRequestBodyReadSize
	if chunkSize > 0 {
		ctx.receivedChunks++
		chunk, err := proxywasm.GetHttpRequestBody(ctx.totalRequestBodyReadSize, chunkSize)
		if err != nil {
			proxywasm.LogCriticalf("failed to get request body: %v", err)
			return types.ActionContinue
		}
		proxywasm.LogInfof("read chunk size: %d", len(chunk))
		if len(chunk) != chunkSize {
			proxywasm.LogErrorf("read data does not match the expected size: %d != %d", len(chunk), chunkSize)
		}
		ctx.totalRequestBodyReadSize += len(chunk)
		if strings.Contains(string(chunk), "pattern") {
			patternFound := fmt.Sprintf("pattern found in chunk: %d", ctx.receivedChunks)
			proxywasm.LogInfo(patternFound)
			if err := proxywasm.SendHttpResponse(403, [][2]string{
				{"powered-by", "proxy-wasm-go-sdk"},
			}, []byte(patternFound), -1); err != nil {
				proxywasm.LogCriticalf("failed to send local response: %v", err)
				proxywasm.ResumeHttpRequest()
			} else {
				proxywasm.LogInfo("local 403 response sent")
			}
			return types.ActionPause
		}
	}

	if !endOfStream {
		// Wait until we see the entire body before sending the request upstream.
		return types.ActionPause
	}
	// When endOfStream is true, we have received the entire body. We expect the total size is equal to the sum of the sizes of the chunks.
	if ctx.totalRequestBodyReadSize != bodySize {
		proxywasm.LogErrorf("read data does not match the expected total size: %d != %d", ctx.totalRequestBodyReadSize, bodySize)
	}
	proxywasm.LogInfof("pattern not found")
	return types.ActionContinue
}
