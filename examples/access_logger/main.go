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
)

func main() {
	proxywasm.SetNewRootContext(newAccessLogger)
}

type accessLogger struct {
	// You'd better embed the default root context
	// so that you don't need to reimplement all the methods by yourself.
	proxywasm.DefaultRootContext
}

func newAccessLogger(contextID uint32) proxywasm.RootContext {
	return &accessLogger{}
}

// Override DefaultRootContext.
func (ctx *accessLogger) OnLog() {
	hdr, err := proxywasm.GetHttpRequestHeader(":path")
	if err != nil {
		proxywasm.LogCritical(err.Error())
		return
	}

	proxywasm.LogInfof("OnLog: :path = %s", hdr)
}
