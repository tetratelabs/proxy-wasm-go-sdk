package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHelloWorld_OnTick(t *testing.T) {
	ctx := newHelloWorld(100)
	host := proxytest.NewHostEmulator(nil, nil, newHelloWorld, nil, nil)
	defer host.Done() // release the host emulation lock so that other test cases can insert their own host emulation
	ctx.OnTick()

	logs := host.GetLogs(types.LogLevelInfo)
	msg := logs[len(logs)-1]

	assert.True(t, strings.Contains(msg, "OnTick on"))
}

func TestHelloWorld_OnVMStart(t *testing.T) {
	host := proxytest.NewHostEmulator(nil, nil, newHelloWorld, nil, nil)
	defer host.Done() // release the host emulation lock so that other test cases can insert their own host emulation

	host.StartVM()
	logs := host.GetLogs(types.LogLevelInfo)
	msg := logs[len(logs)-1]

	assert.True(t, strings.Contains(msg, "proxy_on_vm_start from Go!"))
	assert.Equal(t, tickMilliseconds, host.GetTickPeriod())
}
