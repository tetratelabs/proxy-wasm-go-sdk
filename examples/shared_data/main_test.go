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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestData(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	defer host.Done() // release the host emulation lock so that other test cases can insert their own host emulation

	host.StartVM() // set initial value
	contextID := host.HttpFilterInitContext()
	host.HttpFilterPutRequestHeaders(contextID, nil) // OnHttpRequestHeaders is called

	logs := host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 0)

	assert.Equal(t, "shared value: 1", logs[len(logs)-1])
	host.HttpFilterPutRequestHeaders(contextID, nil) // OnHttpRequestHeaders is called
	host.HttpFilterPutRequestHeaders(contextID, nil) // OnHttpRequestHeaders is called

	logs = host.GetLogs(types.LogLevelInfo)
	assert.Equal(t, "shared value: 3", logs[len(logs)-1])
}
