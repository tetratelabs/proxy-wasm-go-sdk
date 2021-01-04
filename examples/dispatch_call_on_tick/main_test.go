package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestRootContext_OnTick(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	defer host.Done() // release the host emulation lock so that other test cases can insert their own host emulation

	host.StartVM() // call OnVMStart

	assert.Equal(t, tickMilliseconds, host.GetTickPeriod())

	for i := 1; i < 10; i++ {
		host.Tick() // call OnTick
		attrs := host.GetCalloutAttributesFromContext(proxytest.RootContextID)
		require.Equal(t, len(attrs), i)                            // verify DispatchHttpCall is called
		host.PutCalloutResponse(attrs[0].CalloutID, nil, nil, nil) // receive callout response

		logs := host.GetLogs(types.LogLevelInfo)
		require.Greater(t, len(logs), 0)
		msg := logs[len(logs)-1]

		assert.True(t, strings.Contains(msg, fmt.Sprintf("called! %d", i)))
	}

}

func TestRootContext_OnVMStart(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	defer host.Done() // release the host emulation lock so that other test cases can insert their own host emulation

	host.StartVM() // call OnVMStart
	assert.Equal(t, tickMilliseconds, host.GetTickPeriod())
}
