package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHelloWorld_OnTick(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newAccessLogger)
	host := proxytest.NewHostEmulator(opt)
	defer host.Done() // release the host emulation lock so that other test cases can insert their own host emulation

	host.CallOnLogForAccessLogger([][2]string{{":path", "/this/is/path"}}, nil) // call OnLog

	logs := host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 0)
	msg := logs[len(logs)-1]

	assert.True(t, strings.Contains(msg, "/this/is/path"))
}
