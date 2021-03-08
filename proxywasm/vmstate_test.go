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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var currentStateMux sync.Mutex

func TestSetNewRootContext(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	var cnt int
	f := func(uint32) RootContext {
		cnt++
		return nil
	}
	SetNewRootContext(f)
	currentState.newRootContext(0)
	assert.Equal(t, 1, cnt)
}

func TestSetNewHttpContext(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	var cnt int
	f := func(uint32, uint32) HttpContext {
		cnt++
		return nil
	}
	SetNewHttpContext(f)
	currentState.newHttpContext(0, 0)
	assert.Equal(t, 1, cnt)
}

func TestSetNewStreamContext(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	var cnt int
	f := func(uint32, uint32) StreamContext {
		cnt++
		return nil
	}
	SetNewStreamContext(f)
	currentState.newStreamContext(0, 0)
	assert.Equal(t, 1, cnt)
}

func TestState_createRootContext(t *testing.T) {
	t.Run("newRootContext exists", func(t *testing.T) {
		type rc struct{ DefaultRootContext }
		s := &state{
			rootContexts:      map[uint32]*rootContextState{},
			newRootContext:    func(contextID uint32) RootContext { return &rc{} },
			contextIDToRootID: map[uint32]uint32{},
		}

		var cid uint32 = 100
		s.createRootContext(cid)
		assert.NotNil(t, s.rootContexts[cid])
	})

	t.Run("non exists", func(t *testing.T) {
		s := &state{rootContexts: map[uint32]*rootContextState{}, contextIDToRootID: map[uint32]uint32{}}
		var cid uint32 = 100
		s.createRootContext(cid)
		c, ok := s.rootContexts[cid]
		require.True(t, ok)
		_, ok = c.context.(*DefaultRootContext)
		assert.True(t, ok)
	})
}

func TestState_createStreamContext(t *testing.T) {
	type sc struct{ DefaultStreamContext }

	var (
		cid uint32 = 100
		rid uint32 = 10
	)
	s := &state{
		rootContexts:      map[uint32]*rootContextState{rid: nil},
		streams:           map[uint32]StreamContext{},
		newStreamContext:  func(rootContextID, contextID uint32) StreamContext { return &sc{} },
		contextIDToRootID: map[uint32]uint32{},
	}

	s.createStreamContext(cid, rid)
	c, ok := s.streams[cid]
	require.True(t, ok)
	_, ok = c.(*sc)
	assert.True(t, ok)
}

func TestState_createHttpContext(t *testing.T) {
	type hc struct{ DefaultHttpContext }

	var (
		cid uint32 = 100
		rid uint32 = 10
	)
	s := &state{
		rootContexts:      map[uint32]*rootContextState{rid: nil},
		httpStreams:       map[uint32]HttpContext{},
		newHttpContext:    func(rootContextID, contextID uint32) HttpContext { return &hc{} },
		contextIDToRootID: map[uint32]uint32{},
	}

	s.createHttpContext(cid, rid)
	c, ok := s.httpStreams[cid]
	require.True(t, ok)
	_, ok = c.(*hc)
	assert.True(t, ok)

}
