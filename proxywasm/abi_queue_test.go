package proxywasm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

type queueContext struct {
	DefaultContext
	onQueueReady bool
}

func (ctx *queueContext) OnQueueReady(uint32) {
	ctx.onQueueReady = true
}

func Test_queueReady(t *testing.T) {
	var id uint32 = 100
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	currentState = &state{rootContexts: map[uint32]RootContext{id: &queueContext{}}}
	ctx, ok := currentState.rootContexts[id].(*queueContext)
	require.True(t, ok)
	proxyOnQueueReady(id, 10)
	assert.True(t, ctx.onQueueReady)
}
