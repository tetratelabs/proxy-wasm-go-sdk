// Copyright 2020-2021 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestContext_OnPluginStart(t *testing.T) {
	// Setup configurations.
	pluginConfigData := `tinygo plugin configuration`
	opt := proxytest.NewEmulatorOption().
		WithPluginConfiguration([]byte(pluginConfigData)).
		WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnPluginStart.
	require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())

	// Check Envoy logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, "plugin config: "+pluginConfigData)
}

func TestContext_OnVMStart(t *testing.T) {
	// Setup configurations.
	vmConfigData := `tinygo vm configuration`
	opt := proxytest.NewEmulatorOption().
		WithVMConfiguration([]byte(vmConfigData)).
		WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnVMStart.
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())

	// Check Envoy logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, "vm config: "+vmConfigData)
}
