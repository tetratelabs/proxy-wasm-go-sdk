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

type testOnContextCreateRootContext struct {
	types.DefaultRootContext
	cnt int
}

func (ctx *testOnContextCreateRootContext) NewStreamContext(contextID uint32) types.StreamContext {
	if contextID == 100 {
		ctx.cnt += 100
		return &types.DefaultStreamContext{}
	}
	return nil
}

func (ctx *testOnContextCreateRootContext) NewHttpContext(contextID uint32) types.HttpContext {
	if contextID == 1000 {
		ctx.cnt += 1000
		return &types.DefaultHttpContext{}
	}
	return nil
}

func Test_proxyOnContextCreateHttpContext(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	var rootPtr *testOnContextCreateRootContext
	currentState = &state{
		rootContexts: map[uint32]*rootContextState{},
		httpStreams:  map[uint32]types.HttpContext{},
		streams:      map[uint32]types.StreamContext{},
		newRootContext: func(contextID uint32) types.RootContext {
			return &testOnContextCreateRootContext{}
		},
		contextIDToRootID: map[uint32]uint32{},
	}

	SetNewRootContextFn(func(contextID uint32) types.RootContext {
		rootPtr = &testOnContextCreateRootContext{cnt: 1}
		return rootPtr
	})

	proxyOnContextCreate(1, 0)
	require.Equal(t, 1, rootPtr.cnt)

	proxyOnContextCreate(100, 1)
	require.Equal(t, 101, rootPtr.cnt)

	proxyOnContextCreate(1000, 1)
	require.Equal(t, 1101, rootPtr.cnt)
}

type lifecycleContext struct {
	types.DefaultRootContext
	types.DefaultHttpContext
	types.DefaultStreamContext
	onDoneCalled, onLogCalled bool
}

func (ctx *lifecycleContext) OnPluginDone() bool {
	ctx.onDoneCalled = true
	return true
}

func (ctx *lifecycleContext) OnStreamDone() {
	ctx.onDoneCalled = true
}

func (ctx *lifecycleContext) OnHttpStreamDone() {
	ctx.onDoneCalled = true
}

func (ctx *lifecycleContext) OnLog() {
	ctx.onLogCalled = true
}

func Test_onDone(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	currentState = &state{
		rootContexts: map[uint32]*rootContextState{},
		httpStreams:  map[uint32]types.HttpContext{},
		streams:      map[uint32]types.StreamContext{},
	}

	var id uint32 = 1
	ctx := &lifecycleContext{}
	currentState.httpStreams[id] = ctx
	proxyOnDone(id)
	require.True(t, ctx.onDoneCalled)
	require.Equal(t, id, currentState.activeContextID)

	id = 2
	ctx = &lifecycleContext{}
	currentState.streams[id] = ctx
	proxyOnDone(id)
	require.True(t, ctx.onDoneCalled)
	require.Equal(t, id, currentState.activeContextID)

	id = 3
	ctx = &lifecycleContext{}
	currentState.rootContexts[id] = &rootContextState{context: ctx}
	proxyOnDone(id)
	require.True(t, ctx.onDoneCalled)
	require.Equal(t, id, currentState.activeContextID)
}

func Test_onLog(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	currentState = &state{
		rootContexts: map[uint32]*rootContextState{},
		httpStreams:  map[uint32]types.HttpContext{},
		streams:      map[uint32]types.StreamContext{},
	}

	var id uint32 = 1
	ctx := &lifecycleContext{}
	currentState.httpStreams[id] = ctx
	proxyOnLog(id)
	require.True(t, ctx.onLogCalled)
	require.Equal(t, id, currentState.activeContextID)

	id = 2
	ctx = &lifecycleContext{}
	currentState.streams[id] = ctx
	proxyOnLog(id)
	require.True(t, ctx.onLogCalled)
	require.Equal(t, id, currentState.activeContextID)

	id = 3
	ctx = &lifecycleContext{}
	currentState.rootContexts[id] = &rootContextState{context: ctx}
	proxyOnLog(id)
	require.True(t, ctx.onLogCalled)
	require.Equal(t, id, currentState.activeContextID)
}
