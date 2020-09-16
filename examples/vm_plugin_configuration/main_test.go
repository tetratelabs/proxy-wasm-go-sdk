// Copyright 2020 Tetrate
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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestContext_OnConfigure(t *testing.T) {
	pluginConfigData := `{"name": "tinygo plugin configuration"}`
	ctx := context{}
	host, done := proxytest.NewRootFilterHost(ctx, []byte(pluginConfigData), nil)
	defer done() // release the host emulation lock so that other test cases can insert their own host emulation

	host.ConfigurePlugin() // invoke OnConfigure

	logs := host.GetLogs(types.LogLevelInfo)
	msg := logs[len(logs)-1]
	assert.True(t, strings.Contains(msg, pluginConfigData))
}

func TestContext_OnVMStart(t *testing.T) {
	vmConfigData := `{"name": "tinygo vm configuration"}`
	ctx := context{}
	host, done := proxytest.NewRootFilterHost(ctx, nil, []byte(vmConfigData))
	defer done() // release the host emulation lock so that other test cases can insert their own host emulation

	host.StartVM() // invoke OnConfigure

	logs := host.GetLogs(types.LogLevelInfo)
	msg := logs[len(logs)-1]
	assert.True(t, strings.Contains(msg, vmConfigData))
}
