package proxywasm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type configurationContext struct {
	DefaultRootContext
	onVMStartCalled, onPluginStartCalled bool
}

func (c *configurationContext) OnVMStart(int) bool {
	c.onVMStartCalled = true
	return true
}

func (c *configurationContext) OnPluginStart(int) bool {
	c.onPluginStartCalled = true
	return true
}

func Test_proxyOnVMStart(t *testing.T) {
	var rID uint32 = 100
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	currentState = &state{rootContexts: map[uint32]*rootContextState{rID: {context: &configurationContext{}}}}

	proxyOnVMStart(rID, 0)
	ctx, ok := currentState.rootContexts[rID].context.(*configurationContext)
	require.True(t, ok)
	assert.True(t, ctx.onVMStartCalled)
	assert.Equal(t, rID, currentState.activeContextID)

	proxyOnConfigure(rID, 0)
	assert.True(t, ctx.onPluginStartCalled)
	assert.Equal(t, rID, currentState.activeContextID)
}
