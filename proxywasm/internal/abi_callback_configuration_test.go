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

type testConfigurationVMContext struct {
	types.DefaultPluginContext
	onVMStartCalled bool
}

func (c *testConfigurationVMContext) OnVMStart(int) types.OnVMStartStatus {
	c.onVMStartCalled = true
	return true
}

func (c *testConfigurationVMContext) NewPluginContext(uint32) types.PluginContext {
	return &testConfigurationPluginContext{}
}

type testConfigurationPluginContext struct {
	types.DefaultPluginContext
	onPluginStartCalled bool
}

func (c *testConfigurationPluginContext) OnPluginStart(int) types.OnPluginStartStatus {
	c.onPluginStartCalled = true
	return true
}

func Test_pluginInitialization(t *testing.T) {
	currentStateMux.Lock()
	defer currentStateMux.Unlock()

	vmContext := &testConfigurationVMContext{}
	currentState = &state{
		vmContext:         vmContext,
		pluginContexts:    map[uint32]*pluginContextState{},
		contextIDToRootID: map[uint32]uint32{},
	}

	// Check ABI version. There is no return value so just make sure it doesn't panic.
	proxyABIVersion()

	// Call OnVMStart.
	proxyOnVMStart(0, 0)
	require.True(t, vmContext.onVMStartCalled)
	require.Equal(t, uint32(0), currentState.activeContextID)

	// Allocate memory
	require.NotNil(t, proxyOnMemoryAllocate(100))

	// Create plugin context.
	pluginContextID := uint32(100)
	proxyOnContextCreate(pluginContextID, 0)
	require.Contains(t, currentState.pluginContexts, pluginContextID)
	pluginContext := currentState.pluginContexts[pluginContextID].context.(*testConfigurationPluginContext)

	// Call OnPluginStart.
	proxyOnConfigure(pluginContextID, 0)
	require.True(t, pluginContext.onPluginStartCalled)
	require.Equal(t, pluginContextID, currentState.activeContextID)
}
