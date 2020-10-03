package proxywasm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type l4Context struct {
	DefaultStreamContext
	onDownstreamData,
	onDownStreamClose,
	onNewConnection,
	onUpstreamData,
	onUpstreamStreamClose bool
}

func (ctx *l4Context) OnDownstreamData(int, bool) types.Action {
	ctx.onDownstreamData = true
	return types.ActionContinue
}

func (ctx *l4Context) OnDownstreamClose(types.PeerType) { ctx.onDownStreamClose = true }

func (ctx *l4Context) OnNewConnection() types.Action {
	ctx.onNewConnection = true
	return types.ActionContinue
}

func (ctx *l4Context) OnUpstreamData(int, bool) types.Action {
	ctx.onUpstreamData = true
	return types.ActionContinue
}

func (ctx *l4Context) OnUpstreamClose(types.PeerType) {
	ctx.onUpstreamStreamClose = true
}

func Test_l4(t *testing.T) {
	var cID uint32 = 100
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	currentState = &state{streams: map[uint32]StreamContext{cID: &l4Context{}}}
	ctx, ok := currentState.streams[cID].(*l4Context)
	require.True(t, ok)

	proxyOnNewConnection(cID)
	assert.True(t, ctx.onNewConnection)
	proxyOnDownstreamData(cID, 0, false)
	assert.True(t, ctx.onDownstreamData)
	proxyOnDownstreamConnectionClose(cID, 0)
	assert.True(t, ctx.onDownStreamClose)
	proxyOnUpstreamData(cID, 0, false)
	assert.True(t, ctx.onUpstreamData)
	proxyOnUpstreamConnectionClose(cID, 0)
	assert.True(t, ctx.onUpstreamStreamClose)
}
