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
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

var currentStateMux sync.Mutex

type testSetVMContext struct {
	cnt int
}

func (*testSetVMContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	return types.OnVMStartStatusOK
}

func (ctx *testSetVMContext) NewPluginContext(contextID uint32) types.PluginContext {
	ctx.cnt++
	return &types.DefaultPluginContext{}
}

func TestSetVMContext(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	vmContext := &testSetVMContext{}
	SetVMContext(vmContext)
	_ = currentState.vmContext.NewPluginContext(0)
	require.Equal(t, 1, vmContext.cnt)
}

func TestState_createPluginContext(t *testing.T) {
	s := &state{
		pluginContexts:    map[uint32]*pluginContextState{},
		contextIDToRootID: map[uint32]uint32{},
		vmContext:         &testSetVMContext{},
	}

	var cid uint32 = 100
	s.createPluginContext(cid)
	require.NotNil(t, s.pluginContexts[cid])
}

type (
	testStateVMContext     struct{}
	testStatePluginContext struct{ types.DefaultPluginContext }
	testStateTcpContext    struct {
		contextID uint32
		types.DefaultTcpContext
	}
	testStateHttpContext struct {
		contextID uint32
		types.DefaultHttpContext
	}
)

func (*testStateVMContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	return types.OnVMStartStatusOK
}

func (ctx *testStateVMContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &testStatePluginContext{}
}

func (ctx *testStatePluginContext) NewTcpContext(contextID uint32) types.TcpContext {
	return &testStateTcpContext{contextID: contextID}
}

func (ctx *testStatePluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &testStateHttpContext{contextID: contextID}
}

func TestState_createTcpContext(t *testing.T) {
	s := &state{
		pluginContexts:    map[uint32]*pluginContextState{},
		tcpContexts:       map[uint32]types.TcpContext{},
		vmContext:         &testStateVMContext{},
		contextIDToRootID: map[uint32]uint32{},
	}

	var (
		cid uint32 = 100
		rid uint32 = 10
	)
	s.createPluginContext(rid)
	s.createTcpContext(cid, rid)
	c, ok := s.tcpContexts[cid]
	require.True(t, ok)
	ctx, ok := c.(*testStateTcpContext)
	require.True(t, ok)
	require.Equal(t, cid, ctx.contextID)
}

func TestState_createHttpContext(t *testing.T) {
	s := &state{
		pluginContexts:    map[uint32]*pluginContextState{},
		httpContexts:      map[uint32]types.HttpContext{},
		vmContext:         &testStateVMContext{},
		contextIDToRootID: map[uint32]uint32{},
	}

	var (
		cid uint32 = 100
		rid uint32 = 10
	)
	s.createPluginContext(rid)
	s.createHttpContext(cid, rid)
	c, ok := s.httpContexts[cid]
	require.True(t, ok)
	ctx, ok := c.(*testStateHttpContext)
	require.True(t, ok)
	require.Equal(t, cid, ctx.contextID)
}
