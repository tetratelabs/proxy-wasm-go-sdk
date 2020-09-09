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
	f := func(uint32) HttpContext {
		cnt++
		return nil
	}
	SetNewHttpContext(f)
	currentState.newHttpContext(0)
	assert.Equal(t, 1, cnt)
}

func TestSetNewStreamContext(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	var cnt int
	f := func(uint32) StreamContext {
		cnt++
		return nil
	}
	SetNewStreamContext(f)
	currentState.newStreamContext(0)
	assert.Equal(t, 1, cnt)
}

func TestState_createRootContext(t *testing.T) {
	t.Run("newRootContext exists", func(t *testing.T) {
		type rc struct{ DefaultContext }
		s := &state{
			rootContexts:   map[uint32]RootContext{},
			newRootContext: func(contextID uint32) RootContext { return &rc{} },
		}

		var cid uint32 = 100
		s.createRootContext(cid)
		assert.NotNil(t, s.rootContexts[cid])
	})

	t.Run("non exists", func(t *testing.T) {
		s := &state{rootContexts: map[uint32]RootContext{}}
		var cid uint32 = 100
		s.createRootContext(cid)
		c, ok := s.rootContexts[cid]
		require.True(t, ok)
		_, ok = c.(*DefaultContext)
		assert.True(t, ok)
	})
}

func TestState_createStreamContext(t *testing.T) {
	type sc struct{ DefaultContext }

	var (
		cid uint32 = 100
		rid uint32 = 10
	)
	s := &state{
		rootContexts:     map[uint32]RootContext{rid: nil},
		streamContexts:   map[uint32]StreamContext{},
		newStreamContext: func(contextID uint32) StreamContext { return &sc{} },
	}

	s.createStreamContext(cid, rid)
	c, ok := s.streamContexts[cid]
	require.True(t, ok)
	_, ok = c.(*sc)
	assert.True(t, ok)
}

func TestState_createHttpContext(t *testing.T) {
	type hc struct{ DefaultContext }

	var (
		cid uint32 = 100
		rid uint32 = 10
	)
	s := &state{
		rootContexts:   map[uint32]RootContext{rid: nil},
		httpContexts:   map[uint32]HttpContext{},
		newHttpContext: func(contextID uint32) HttpContext { return &hc{} },
	}

	s.createHttpContext(cid, rid)
	c, ok := s.httpContexts[cid]
	require.True(t, ok)
	_, ok = c.(*hc)
	assert.True(t, ok)

}

func TestState_registerCallout(t *testing.T) {
	var calloutID uint32 = 100
	s := &state{callOuts: map[uint32]uint32{}, activeContextID: 200}
	s.registerCallout(calloutID)
	assert.Equal(t, s.callOuts[calloutID], s.activeContextID)
}

func TestState_setActiveContextID(t *testing.T) {
	s := state{}
	var cID uint32 = 100
	s.setActiveContextID(cID)
	assert.Equal(t, s.activeContextID, cID)
}
