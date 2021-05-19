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
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestQueue(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// Call OnVMStart -> register the queue, and set tick period.
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	require.Contains(t, logs, fmt.Sprintf("queue registered, name: %s, id: %d", queueName, queueID))

	// Initialize http context.
	contextID := host.InitializeHttpContext()

	// Call OnRequestHeaders
	action := host.CallOnRequestHeaders(contextID, nil, false)
	require.Equal(t, types.ActionContinue, action)

	// Check Envoy logs.
	logs = host.GetLogs(types.LogLevelInfo)
	require.Contains(t, logs, "dequeued data: hello")
	require.Contains(t, logs, "dequeued data: world")
	require.Contains(t, logs, "dequeued data: hello")
	require.Contains(t, logs, "dequeued data: proxy-wasm")
}
