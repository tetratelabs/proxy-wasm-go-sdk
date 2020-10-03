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
	"hash/fnv"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const clusterName = "httpbin"

func main() {
	proxywasm.SetNewHttpContext(newContext)
}

type httpAuthRandom struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	proxywasm.DefaultHttpContext
	contextID uint32
}

func newContext(rootContextID, contextID uint32) proxywasm.HttpContext {
	return &httpAuthRandom{contextID: contextID}
}

// override default
func (ctx *httpAuthRandom) OnHttpRequestHeaders(int, bool) types.Action {
	hs, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
		return types.ActionContinue
	}
	for _, h := range hs {
		proxywasm.LogInfof("request header: %s: %s", h[0], h[1])
	}

	if _, err := proxywasm.DispatchHttpCall(clusterName, hs, "", [][2]string{},
		50000, httpCallResponseCallback); err != nil {
		proxywasm.LogCriticalf("dipatch httpcall failed: %v", err)
		return types.ActionContinue
	}

	proxywasm.LogInfof("http call dispatched to %s", clusterName)
	return types.ActionPause
}

func httpCallResponseCallback(_ int, bodySize int, _ int) {
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

	msg := "access forbidden"
	proxywasm.LogInfo(msg)
	proxywasm.SendHttpResponse(403, [][2]string{
		{"powered-by", "proxy-wasm-go-sdk!!"},
	}, msg)
}
