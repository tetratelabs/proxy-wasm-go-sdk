package proxywasm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func Test_proxyOnContextCreate(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	var cnt int
	currentState = &state{
		rootContexts:   map[uint32]RootContext{},
		httpContexts:   map[uint32]HttpContext{},
		streamContexts: map[uint32]StreamContext{},
	}

	SetNewRootContext(func(contextID uint32) RootContext {
		cnt++
		return nil
	})

	proxyOnContextCreate(100, 0)
	require.Equal(t, 1, cnt)
	SetNewHttpContext(func(contextID uint32) HttpContext {
		cnt += 100
		return nil
	})
	proxyOnContextCreate(100, 100)
	require.Equal(t, 101, cnt)
	currentState.newHttpContext = nil

	SetNewStreamContext(func(contextID uint32) StreamContext {
		cnt += 1000
		return nil
	})
	proxyOnContextCreate(100, 100)
	require.Equal(t, 1101, cnt)
}

func Test_proxyOnDelete(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	currentState = &state{
		rootContexts:   map[uint32]RootContext{},
		httpContexts:   map[uint32]HttpContext{},
		streamContexts: map[uint32]StreamContext{},
	}

	var id uint32 = 100
	var ctx = &DefaultContext{}
	currentState.streamContexts[id] = ctx
	proxyOnDelete(id)
	assert.Nil(t, currentState.streamContexts[id])

	currentState.httpContexts[id] = ctx
	proxyOnDelete(id)
	assert.Nil(t, currentState.httpContexts[id])

	currentState.rootContexts[id] = ctx
	proxyOnDelete(id)
	assert.Nil(t, currentState.rootContexts[id])
}

type lifecycleContext struct {
	DefaultContext
	onDone, onLog bool
}

func (ctx *lifecycleContext) OnLog() {
	ctx.onLog = true
}
func (ctx *lifecycleContext) OnDone() bool {
	ctx.onDone = true
	return true
}

func Test_onDone(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	currentState = &state{
		rootContexts:   map[uint32]RootContext{},
		httpContexts:   map[uint32]HttpContext{},
		streamContexts: map[uint32]StreamContext{},
	}

	var id uint32 = 1
	ctx := &lifecycleContext{}
	currentState.rootContexts[id] = ctx
	proxyOnDone(id)
	assert.True(t, ctx.onDone)
	assert.Equal(t, id, currentState.activeContextID)

	id = 2
	ctx = &lifecycleContext{}
	currentState.httpContexts[id] = ctx
	proxyOnDone(id)
	assert.True(t, ctx.onDone)
	assert.Equal(t, id, currentState.activeContextID)

	id = 3
	ctx = &lifecycleContext{}
	currentState.rootContexts[id] = ctx
	proxyOnDone(id)
	assert.True(t, ctx.onDone)
	assert.Equal(t, id, currentState.activeContextID)
}

func Test_onLog(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	currentState = &state{
		rootContexts:   map[uint32]RootContext{},
		httpContexts:   map[uint32]HttpContext{},
		streamContexts: map[uint32]StreamContext{},
	}

	var id uint32 = 1
	ctx := &lifecycleContext{}
	currentState.rootContexts[id] = ctx
	proxyOnLog(id)
	assert.True(t, ctx.onLog)
	assert.Equal(t, id, currentState.activeContextID)

	id = 2
	ctx = &lifecycleContext{}
	currentState.httpContexts[id] = ctx
	proxyOnLog(id)
	assert.True(t, ctx.onLog)
	assert.Equal(t, id, currentState.activeContextID)

	id = 3
	ctx = &lifecycleContext{}
	currentState.rootContexts[id] = ctx
	proxyOnLog(id)
	assert.True(t, ctx.onLog)
	assert.Equal(t, id, currentState.activeContextID)
}
