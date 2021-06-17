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

func TestSetNewRootContextFn(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	var cnt int
	f := func(uint32) types.RootContext {
		cnt++
		return nil
	}
	SetNewRootContextFn(f)
	currentState.newRootContext(0)
	require.Equal(t, 1, cnt)
}

func TestState_createRootContext(t *testing.T) {
	t.Run("newRootContext exists", func(t *testing.T) {
		type rc struct{ types.DefaultRootContext }
		s := &state{
			rootContexts:      map[uint32]*rootContextState{},
			newRootContext:    func(contextID uint32) types.RootContext { return &rc{} },
			contextIDToRootID: map[uint32]uint32{},
		}

		var cid uint32 = 100
		s.createRootContext(cid)
		require.NotNil(t, s.rootContexts[cid])
	})

	t.Run("non exists", func(t *testing.T) {
		s := &state{rootContexts: map[uint32]*rootContextState{}, contextIDToRootID: map[uint32]uint32{}}
		var cid uint32 = 100
		s.createRootContext(cid)
		c, ok := s.rootContexts[cid]
		require.True(t, ok)
		_, ok = c.context.(*types.DefaultRootContext)
		require.True(t, ok)
	})
}

type (
	testStateRootContext   struct{ types.DefaultRootContext }
	testStateStreamContext struct {
		contextID uint32
		types.DefaultStreamContext
	}
	testStateHttpContext struct {
		contextID uint32
		types.DefaultHttpContext
	}
)

func (ctx *testStateRootContext) NewStreamContext(contextID uint32) types.StreamContext {
	return &testStateStreamContext{contextID: contextID}
}

func (ctx *testStateRootContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &testStateHttpContext{contextID: contextID}
}

func TestState_createStreamContext(t *testing.T) {
	var (
		cid uint32 = 100
		rid uint32 = 10
	)
	s := &state{
		rootContexts: map[uint32]*rootContextState{rid: nil},
		streams:      map[uint32]types.StreamContext{},
		newRootContext: func(contextID uint32) types.RootContext {
			return &testStateRootContext{}
		},
		contextIDToRootID: map[uint32]uint32{},
	}

	s.createRootContext(rid)
	s.createStreamContext(cid, rid)
	c, ok := s.streams[cid]
	require.True(t, ok)
	ctx, ok := c.(*testStateStreamContext)
	require.True(t, ok)
	require.Equal(t, cid, ctx.contextID)
}

func TestState_createHttpContext(t *testing.T) {
	var (
		cid uint32 = 100
		rid uint32 = 10
	)
	s := &state{
		rootContexts: map[uint32]*rootContextState{rid: nil},
		httpStreams:  map[uint32]types.HttpContext{},
		newRootContext: func(contextID uint32) types.RootContext {
			return &testStateRootContext{}
		},
		contextIDToRootID: map[uint32]uint32{},
	}

	s.createRootContext(rid)
	s.createHttpContext(cid, rid)
	c, ok := s.httpStreams[cid]
	require.True(t, ok)
	ctx, ok := c.(*testStateHttpContext)
	require.True(t, ok)
	require.Equal(t, cid, ctx.contextID)
}
