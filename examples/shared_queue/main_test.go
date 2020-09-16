package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestQueue(t *testing.T) {
	ctx := queue{}
	host, done := proxytest.NewRootFilterHost(ctx, nil, nil)
	defer done() // release the host emulation lock so that other test cases can insert their own host emulation

	host.StartVM() // register the queue,set tick period

	logs := host.GetLogs(types.LogLevelInfo)
	assert.Equal(t, logs[0], fmt.Sprintf("queue registered, name: %s, id: %d", queueName, queueID))
	assert.Equal(t, tickMilliseconds, host.GetTickPeriod())

	ctx.OnHttpRequestHeaders(0, false) // call enqueue
	assert.Equal(t, 4, host.GetQueueSize(queueID))

	for i := 0; i < 4; i++ {
		ctx.OnTick() // dequeue
	}
	logs = host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 4)

	assert.Equal(t, "dequeued data: hello", logs[len(logs)-4])
	assert.Equal(t, "dequeued data: world", logs[len(logs)-3])
	assert.Equal(t, "dequeued data: hello", logs[len(logs)-2])
	assert.Equal(t, "dequeued data: proxy-wasm", logs[len(logs)-1])
}
