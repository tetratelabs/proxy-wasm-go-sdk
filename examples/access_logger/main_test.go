package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestAccessLogger_OnTick(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newAccessLogger)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// Call OnLog.
	host.CallOnLogForAccessLogger(types.Headers{{":path", "/this/is/path"}}, nil)

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	assert.Contains(t, logs, "OnLog: :path = /this/is/path")
}
