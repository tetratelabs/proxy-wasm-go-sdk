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

type testOnContextCreatePluginVMContext struct{}

func (*testOnContextCreatePluginVMContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	return types.OnVMStartStatusOK
}

func (*testOnContextCreatePluginVMContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &testOnContextCreatePluginContext{cnt: 1}
}

type testOnContextCreatePluginContext struct {
	types.DefaultPluginContext
	cnt int
}

func (ctx *testOnContextCreatePluginContext) NewTcpContext(contextID uint32) types.TcpContext {
	if contextID == 100 {
		ctx.cnt += 100
		return &types.DefaultTcpContext{}
	}
	return nil
}

func (ctx *testOnContextCreatePluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	if contextID == 1000 {
		ctx.cnt += 1000
		return &types.DefaultHttpContext{}
	}
	return nil
}

func Test_proxyOnContextCreate(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	vmContext := &testOnContextCreatePluginVMContext{}
	currentState = &state{
		pluginContexts:    map[uint32]*pluginContextState{},
		httpContexts:      map[uint32]types.HttpContext{},
		tcpContexts:       map[uint32]types.TcpContext{},
		contextIDToRootID: map[uint32]uint32{},
	}

	// Set the VM context.
	SetVMContext(vmContext)

	// Create Plugin context.
	proxyOnContextCreate(1, 0)
	require.Contains(t, currentState.pluginContexts, uint32(1))
	pluginContext := currentState.pluginContexts[1].context.(*testOnContextCreatePluginContext)
	require.Equal(t, 1, pluginContext.cnt)

	// Create Http contexts.
	proxyOnContextCreate(100, 1)
	require.Equal(t, 101, pluginContext.cnt)
	proxyOnContextCreate(1000, 1)
	require.Equal(t, 1101, pluginContext.cnt)
}

type lifecycleContext struct {
	types.DefaultPluginContext
	types.DefaultHttpContext
	types.DefaultTcpContext
	onDoneCalled bool
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

func Test_onDone_or_onLog(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	currentState = &state{
		pluginContexts: map[uint32]*pluginContextState{},
		httpContexts:   map[uint32]types.HttpContext{},
		tcpContexts:    map[uint32]types.TcpContext{},
	}

	// Stream Contexts are only called on on_log, not on on_done.
	var id uint32 = 1
	ctx := &lifecycleContext{}
	currentState.httpContexts[id] = ctx
	proxyOnLog(id)
	require.True(t, ctx.onDoneCalled)
	require.Equal(t, id, currentState.activeContextID)

	id = 2
	ctx = &lifecycleContext{}
	currentState.tcpContexts[id] = ctx
	proxyOnLog(id)
	require.True(t, ctx.onDoneCalled)
	require.Equal(t, id, currentState.activeContextID)

	// Root Contexts are only called on on_done, not on on_log.
	id = 3
	ctx = &lifecycleContext{}
	currentState.pluginContexts[id] = &pluginContextState{context: ctx}
	proxyOnDone(id)
	require.True(t, ctx.onDoneCalled)
	require.Equal(t, id, currentState.activeContextID)
}
