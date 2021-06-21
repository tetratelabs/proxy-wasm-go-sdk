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

type queueContext struct {
	types.DefaultPluginContext
	onQueueReady bool
}

func (ctx *queueContext) OnQueueReady(uint32) {
	ctx.onQueueReady = true
}

func Test_queueReady(t *testing.T) {
	var id uint32 = 100
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	currentState = &state{pluginContexts: map[uint32]*pluginContextState{id: {context: &queueContext{}}}}
	ctx, ok := currentState.pluginContexts[id].context.(*queueContext)
	require.True(t, ok)
	proxyOnQueueReady(id, 10)
	require.True(t, ctx.onQueueReady)
}
