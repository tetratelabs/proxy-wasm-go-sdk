// Copyright 2021 Tetrate
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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"runtime"
	"strings"
	"time"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const (
	// The secret key used to sign the JWT token.
	secretKey = "secret"
	// tick period in milliseconds.
	tickMilliseconds uint32 = 100
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
	return &pluginContext{contextID: contextID}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
	contextID uint32
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{contextID: contextID}
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	// Set tick period. OnTick is called every tick period.
	if err := proxywasm.SetTickPeriodMilliSeconds(tickMilliseconds); err != nil {
		proxywasm.LogCriticalf("failed to set tick period: %v", err)
	}

	return types.OnPluginStartStatusOK
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) OnTick() {
	// These lines are used for observing the memory stability during the e2e test.
	// TODO(musaprg): split this into another package and import it.
	// 				  These lines are only used for testing, which should be separated from example code.
	t := time.Now().UnixNano()
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	proxywasm.LogInfof("[memstat][contextID=%d][unixnanotime=%d] heap size: in-use / reserved = %d / %d bytes", ctx.contextID, t, mem.HeapInuse, mem.HeapSys)
}

type httpContext struct {
	// Embed the default plugin context
	// so that you don't need to reimplement all the methods by yourself.
	types.DefaultHttpContext
	contextID uint32
}

// Override types.DefaultHttpContext.
func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	authorization, err := proxywasm.GetHttpRequestHeader("Authorization")
	if err != nil {
		if err := proxywasm.SendHttpResponse(400, nil, []byte("authorization header must be provided"), -1); err != nil {
			panic(err)
		}
		return types.ActionPause
	}

	proxywasm.LogInfof("authorization token: %s", authorization)

	// Validate format and verify token.
	slice := strings.Fields(authorization)
	if len(slice) != 2 || slice[0] != "Bearer" || !verifyToken(slice[1]) {
		if err := proxywasm.SendHttpResponse(401, nil, []byte("invalid token"), -1); err != nil {
			panic(err)
		}
		return types.ActionPause
	}

	proxywasm.LogInfof("request authorized!")

	return types.ActionContinue
}

// Override types.DefaultHttpContext.
func (ctx *httpContext) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}

// verifyToken checks if the JWT token is valid.
func verifyToken(token string) bool {
	slice := strings.Split(token, ".")
	if len(slice) != 3 {
		return false
	}
	unsignedToken := strings.Join(slice[:2], ".")
	signature, err := base64.RawURLEncoding.DecodeString(slice[2])
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(unsignedToken))
	expectedSignature := mac.Sum(nil)
	return hmac.Equal(signature, expectedSignature)
}
