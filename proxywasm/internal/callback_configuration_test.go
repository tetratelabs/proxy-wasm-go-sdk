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

type configurationContext struct {
	types.DefaultRootContext
	onVMStartCalled, onPluginStartCalled bool
}

func (c *configurationContext) OnVMStart(int) types.OnVMStartStatus {
	c.onVMStartCalled = true
	return true
}

func (c *configurationContext) OnPluginStart(int) types.OnPluginStartStatus {
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
	require.True(t, ctx.onVMStartCalled)
	require.Equal(t, rID, currentState.activeContextID)

	proxyOnConfigure(rID, 0)
	require.True(t, ctx.onPluginStartCalled)
	require.Equal(t, rID, currentState.activeContextID)
}
