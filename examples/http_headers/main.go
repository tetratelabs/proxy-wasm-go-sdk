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
	"strconv"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetNewHttpContext(newContext)
}

type httpHeaders struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	proxywasm.DefaultContext
	contextID uint32
}

func newContext(contextID uint32) proxywasm.HttpContext {
	return &httpHeaders{contextID: contextID}
}

// override
func (ctx *httpHeaders) OnHttpRequestHeaders(int, bool) types.Action {
	hs, err := proxywasm.HostCallGetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCritical("failed to get request headers: ", err.Error())
	}

	for _, h := range hs {
		proxywasm.LogInfo("request header: ", h[0], ": ", h[1])
	}
	return types.ActionContinue
}

// override
func (ctx *httpHeaders) OnHttpResponseHeaders(int, bool) types.Action {
	hs, err := proxywasm.HostCallGetHttpResponseHeaders()
	if err != nil {
		proxywasm.LogCritical("failed to get request headers: ", err.Error())
	}

	for _, h := range hs {
		proxywasm.LogInfo("response header: ", h[0], ": ", h[1])
	}
	return types.ActionContinue
}

// override
func (ctx *httpHeaders) OnLog() {
	proxywasm.LogInfo(strconv.FormatUint(uint64(ctx.contextID), 10), " finished")
}
