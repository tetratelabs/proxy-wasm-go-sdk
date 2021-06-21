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

type timerContext struct {
	types.DefaultPluginContext
	onTick bool
}

func (ctx *timerContext) OnTick() {
	ctx.onTick = true
}

func Test_onTick(t *testing.T) {
	var id uint32 = 100
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	currentState = &state{pluginContexts: map[uint32]*pluginContextState{id: {context: &timerContext{}}}}
	ctx, ok := currentState.pluginContexts[id].context.(*timerContext)
	require.True(t, ok)
	proxyOnTick(id)
	require.True(t, ctx.onTick)
}
