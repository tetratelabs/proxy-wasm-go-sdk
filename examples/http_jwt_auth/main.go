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
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const (
	// The secret key used to sign the JWT token.
	secretKey = "secret"
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

	// Validate header format.
	slice := strings.Fields(authorization)
	if len(slice) != 2 || slice[0] != "Bearer" {
		if err := proxywasm.SendHttpResponse(400, nil, []byte("invalid authorization header"), -1); err != nil {
			panic(err)
		}
		return types.ActionPause
	}

	// Verify token.
	ok, err := verifyToken(slice[1])
	if err != nil {
		proxywasm.LogWarnf("failed to verify token with an error: %s", err)
	}
	if !ok {
		if err := proxywasm.SendHttpResponse(401, nil, []byte("invalid token"), -1); err != nil {
			panic(err)
		}
		return types.ActionPause
	}

	proxywasm.LogInfof("request authorized!")

	return types.ActionContinue
}

// verifyToken checks if the JWT token is valid.
func verifyToken(token string) (bool, error) {
	sepIdx := strings.LastIndex(token, ".")
	if sepIdx == -1 {
		return false, nil
	}
	unsignedToken := token[:sepIdx]
	signature, err := base64.RawURLEncoding.DecodeString(token[sepIdx+1:])
	if err != nil {
		return false, err
	}
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(unsignedToken))
	expectedSignature := mac.Sum(nil)
	return hmac.Equal(signature, expectedSignature), nil
}
