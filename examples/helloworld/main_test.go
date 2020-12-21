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
		WithNewRootContext(newHelloWorld)
	host := proxytest.NewHostEmulator(opt)
	defer host.Done() // release the host emulation lock so that other test cases can insert their own host emulation

	host.StartVM() // call OnVMStart

	assert.Equal(t, tickMilliseconds, host.GetTickPeriod())
	host.Tick() // call OnTick

	logs := host.GetLogs(types.LogLevelInfo)
	require.Greater(t, len(logs), 0)
	msg := logs[len(logs)-1]

	assert.True(t, strings.Contains(msg, "It's"))
}

func TestHelloWorld_OnVMStart(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newHelloWorld)
	host := proxytest.NewHostEmulator(opt)
	defer host.Done() // release the host emulation lock so that other test cases can insert their own host emulation

	host.StartVM() // call OnVMStart
	logs := host.GetLogs(types.LogLevelInfo)
	msg := logs[len(logs)-1]

	assert.True(t, strings.Contains(msg, "proxy_on_vm_start from Go!"))
	assert.Equal(t, tickMilliseconds, host.GetTickPeriod())
}
