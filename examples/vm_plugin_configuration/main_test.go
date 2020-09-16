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
