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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type l4Context struct {
	types.DefaultTcpContext
	onDownstreamData,
	onDownstreamClose,
	onNewConnection,
	onUpstreamData,
	onUpstreamStreamClose bool
}

func (ctx *l4Context) OnDownstreamData(int, bool) types.Action {
	ctx.onDownstreamData = true
	return types.ActionContinue
}

func (ctx *l4Context) OnDownstreamClose(types.PeerType) { ctx.onDownstreamClose = true }

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

	currentState = &state{tcpContexts: map[uint32]types.TcpContext{cID: &l4Context{}}}
	ctx, ok := currentState.tcpContexts[cID].(*l4Context)
	require.True(t, ok)

	proxyOnNewConnection(cID)
	require.True(t, ctx.onNewConnection)
	proxyOnDownstreamData(cID, 0, false)
	require.True(t, ctx.onDownstreamData)
	proxyOnDownstreamConnectionClose(cID, 0)
	require.True(t, ctx.onDownstreamClose)
	proxyOnUpstreamData(cID, 0, false)
	require.True(t, ctx.onUpstreamData)
	proxyOnUpstreamConnectionClose(cID, 0)
	require.True(t, ctx.onUpstreamStreamClose)
}
