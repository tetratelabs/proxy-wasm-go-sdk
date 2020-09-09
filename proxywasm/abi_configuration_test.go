package proxywasm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type configurationContext struct {
	DefaultContext
	onVMStartCalled, onConfigureCalled bool
}

func (c *configurationContext) OnVMStart(vmConfigurationSize int) bool {
	c.onVMStartCalled = true
	return true
}

func (c *configurationContext) OnConfigure(pluginConfigurationSize int) bool {
	c.onConfigureCalled = true
	return true
}

func Test_proxyOnVMStart(t *testing.T) {
	var rID uint32 = 100
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	currentState = &state{rootContexts: map[uint32]RootContext{rID: &configurationContext{}}}

	proxyOnVMStart(rID, 0)
	ctx, ok := currentState.rootContexts[rID].(*configurationContext)
	require.True(t, ok)
	assert.True(t, ctx.onVMStartCalled)
	assert.Equal(t, rID, currentState.activeContextID)

	proxyOnConfigure(rID, 0)
	assert.True(t, ctx.onConfigureCalled)
	assert.Equal(t, rID, currentState.activeContextID)
}
