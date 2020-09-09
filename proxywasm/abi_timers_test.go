package proxywasm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type timerContext struct {
	DefaultContext
	onTick bool
}

func (ctx *timerContext) OnTick() {
	ctx.onTick = true
}

func Test_onTick(t *testing.T) {
	var id uint32 = 100
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	currentState = &state{rootContexts: map[uint32]RootContext{id: &timerContext{}}}
	ctx, ok := currentState.rootContexts[id].(*timerContext)
	require.True(t, ok)
	proxyOnTick(id)
	assert.True(t, ctx.onTick)
}
