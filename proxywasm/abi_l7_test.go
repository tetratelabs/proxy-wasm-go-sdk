// Copyright 2021 Tetratea
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxywasm

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/rawhostcall"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type l7Context struct {
	DefaultHttpContext
	onHttpRequestHeaders,
	onHttpRequestBody,
	onHttpRequestTrailers,
	onHttpResponseHeaders,
	onHttpResponseBody,
	onHttpResponseTrailers,
	onHttpCallResponse bool
}

func (ctx *l7Context) OnHttpRequestHeaders(int, bool) types.Action {
	ctx.onHttpRequestHeaders = true
	return types.ActionContinue
}

func (ctx *l7Context) OnHttpRequestBody(int, bool) types.Action {
	ctx.onHttpRequestBody = true
	return types.ActionContinue
}

func (ctx *l7Context) OnHttpRequestTrailers(int) types.Action {
	ctx.onHttpRequestTrailers = true
	return types.ActionContinue
}

func (ctx *l7Context) OnHttpResponseHeaders(int, bool) types.Action {
	ctx.onHttpResponseHeaders = true
	return types.ActionContinue
}

func (ctx *l7Context) OnHttpResponseBody(int, bool) types.Action {
	ctx.onHttpResponseBody = true
	return types.ActionContinue
}

func (ctx *l7Context) OnHttpResponseTrailers(int) types.Action {
	ctx.onHttpResponseTrailers = true
	return types.ActionContinue
}

func (ctx *l7Context) OnHttpCallResponse(int, int, int) {
	ctx.onHttpCallResponse = true
}

func Test_l7(t *testing.T) {
	var cID uint32 = 100
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	currentState = &state{httpStreams: map[uint32]HttpContext{cID: &l7Context{}}}
	ctx, ok := currentState.httpStreams[cID].(*l7Context)
	require.True(t, ok)

	proxyOnRequestHeaders(cID, 0, false)
	require.True(t, ctx.onHttpRequestHeaders)
	proxyOnRequestBody(cID, 0, false)
	require.True(t, ctx.onHttpRequestBody)
	proxyOnRequestTrailers(cID, 0)
	require.True(t, ctx.onHttpRequestTrailers)
	proxyOnResponseHeaders(cID, 0, false)
	require.True(t, ctx.onHttpResponseHeaders)
	proxyOnResponseBody(cID, 0, false)
	require.True(t, ctx.onHttpResponseBody)
	proxyOnResponseTrailers(cID, 0)
	require.True(t, ctx.onHttpResponseTrailers)
}

func Test_proxyOnHttpCallResponse(t *testing.T) {
	hostMutex.Lock()
	defer hostMutex.Unlock()
	rawhostcall.RegisterMockWASMHost(rawhostcall.DefaultProxyWAMSHost{})

	var (
		rootContextID uint32 = 1
		callOutID     uint32 = 10
	)

	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	ctx := &l7Context{}
	currentState = &state{
		rootContexts: map[uint32]*rootContextState{rootContextID: {
			httpCallbacks: map[uint32]*httpCallbackAttribute{callOutID: {callback: ctx.OnHttpCallResponse}},
		}},
	}

	proxyOnHttpCallResponse(rootContextID, callOutID, 0, 0, 0)
	_, ok := currentState.rootContexts[rootContextID].httpCallbacks[callOutID]
	require.False(t, ok)
	require.True(t, ctx.onHttpCallResponse)

	ctx = &l7Context{}
	currentState = &state{
		rootContexts: map[uint32]*rootContextState{rootContextID: {
			httpCallbacks: map[uint32]*httpCallbackAttribute{callOutID: {callback: ctx.OnHttpCallResponse}},
		}},
	}

	proxyOnHttpCallResponse(rootContextID, callOutID, 0, 0, 0)
	_, ok = currentState.rootContexts[rootContextID].httpCallbacks[callOutID]
	require.False(t, ok)
	require.True(t, ctx.onHttpCallResponse)
}
