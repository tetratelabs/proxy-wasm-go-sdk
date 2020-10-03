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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestQueue(t *testing.T) {
	host := proxytest.NewHostEmulator(nil, nil,
		newRootContext, nil, newHttpContext)
	defer host.Done() // release the host emulation lock so that other test cases can insert their own host emulation

	host.StartVM() // register the queue,set tick period

	logs := host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 0)
	assert.Equal(t, logs[0], fmt.Sprintf("queue registered, name: %s, id: %d", queueName, queueID))
	assert.Equal(t, tickMilliseconds, host.GetTickPeriod())

	contextID := host.HttpFilterInitContext()
	host.HttpFilterPutRequestHeaders(contextID, nil) // call enqueue
	assert.Equal(t, 4, host.GetQueueSize(queueID))

	time.Sleep(time.Duration(tickMilliseconds*5) * time.Millisecond)

	logs = host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 5)

	assert.Equal(t, "dequeued data: hello", logs[len(logs)-4])
	assert.Equal(t, "dequeued data: world", logs[len(logs)-3])
	assert.Equal(t, "dequeued data: hello", logs[len(logs)-2])
	assert.Equal(t, "dequeued data: proxy-wasm", logs[len(logs)-1])
}
