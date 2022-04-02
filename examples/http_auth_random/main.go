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
	"hash/fnv"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const clusterName = "httpbin"

func main() {
	proxywasm.LogInfo("<---- main ---->")
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	proxywasm.LogInfo("<---- vmCx NewPluginContext ---->")
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

// Override types.DefaultPluginContext.
func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	proxywasm.LogInfo("<---- pluginCx NewHttpContext ---->")
	return &httpAuthRandom{contextID: contextID}
}

type httpAuthRandom struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID uint32
}

// Override types.DefaultHttpContext.
func (ctx *httpAuthRandom) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("<---- OnHttpRequestHeaders ---->")
	hs, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
		return types.ActionContinue
	}
	for _, h := range hs {
		proxywasm.LogInfof("request header: %s: %s", h[0], h[1])
	}

	if _, err := proxywasm.DispatchHttpCall(clusterName, hs, nil, nil,
		50000, httpCallResponseCallback); err != nil {
		proxywasm.LogCriticalf("dipatch httpcall failed: %v", err)
		return types.ActionContinue
	}

	proxywasm.LogInfof("http call dispatched to %s", clusterName)
	return types.ActionPause
}

func httpCallResponseCallback(numHeaders, bodySize, numTrailers int) {
	proxywasm.LogInfo("<---- httpCallResponseCallBack ---->")
	hs, err := proxywasm.GetHttpCallResponseHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get response body: %v", err)
		return
	}

	for _, h := range hs {
		proxywasm.LogInfof("response header from %s: %s: %s", clusterName, h[0], h[1])
	}

	b, err := proxywasm.GetHttpCallResponseBody(0, bodySize)
	if err != nil {
		proxywasm.LogCriticalf("failed to get response body: %v", err)
		proxywasm.ResumeHttpRequest()
		return
	}

	s := fnv.New32a()
	if _, err := s.Write(b); err != nil {
		proxywasm.LogCriticalf("failed to calculate hash: %v", err)
		proxywasm.ResumeHttpRequest()
		return
	}

	if s.Sum32()%2 == 0 {
		proxywasm.LogInfo("access granted")
		proxywasm.ResumeHttpRequest()
		return
	}

	body := "access forbidden"
	proxywasm.LogInfo(body)
	if err := proxywasm.SendHttpResponse(403, [][2]string{
		{"powered-by", "proxy-wasm-go-sdk!!"},
	}, []byte(body), -1); err != nil {
		proxywasm.LogErrorf("failed to send local response: %v", err)
		proxywasm.ResumeHttpRequest()
	}
}
