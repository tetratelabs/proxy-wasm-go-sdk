// Copyright 2021 Tetrate
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

package internal

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type l7Context struct {
	types.DefaultHttpContext
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

	currentState = &state{httpContexts: map[uint32]types.HttpContext{cID: &l7Context{}}}
	ctx, ok := currentState.httpContexts[cID].(*l7Context)
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
	release := RegisterMockWasmHost(DefaultProxyWAMSHost{})
	defer release()

	var (
		pluginContextID uint32 = 1
		callerContextID uint32 = 100
		callOutID       uint32 = 10
	)

	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	t.Run("normal", func(t *testing.T) {
		ctx := &l7Context{}
		currentState = &state{
			pluginContexts: map[uint32]*pluginContextState{pluginContextID: {
				httpCallbacks: map[uint32]*httpCallbackAttribute{callOutID: {callback: ctx.OnHttpCallResponse, callerContextID: callerContextID}},
			}},
			contextIDToRootID: map[uint32]uint32{callerContextID: pluginContextID},
		}

		proxyOnHttpCallResponse(pluginContextID, callOutID, 0, 0, 0)
		_, ok := currentState.pluginContexts[pluginContextID].httpCallbacks[callOutID]
		require.False(t, ok)
		require.True(t, ctx.onHttpCallResponse)
	})

	t.Run("delete before callback", func(t *testing.T) {
		ctx := &l7Context{}
		currentState = &state{
			pluginContexts: map[uint32]*pluginContextState{pluginContextID: {
				httpCallbacks: map[uint32]*httpCallbackAttribute{callOutID: {callback: ctx.OnHttpCallResponse, callerContextID: callerContextID}},
			}},
			httpContexts:      map[uint32]types.HttpContext{callerContextID: nil},
			contextIDToRootID: map[uint32]uint32{callerContextID: pluginContextID},
		}

		proxyOnDelete(callerContextID)

		proxyOnHttpCallResponse(pluginContextID, callOutID, 0, 0, 0)
		_, ok := currentState.pluginContexts[pluginContextID].httpCallbacks[callOutID]
		require.False(t, ok)

		// If the caller context is deleted before callback is called, then
		// the callback shouldn't be called.
		require.False(t, ctx.onHttpCallResponse)
	})

}
