//go:build !proxywasm_timing

package proxytest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type noopPlugin struct {
	types.DefaultVMContext
	tcp bool
}

// NewPluginContext implements the same method on types.DefaultVMContext.
func (p *noopPlugin) NewPluginContext(uint32) types.PluginContext {
	return &noopPluginContext{tcp: p.tcp}
}

type noopPluginContext struct {
	types.DefaultPluginContext
	tcp bool
}

// NewHttpContext implements the same method on types.DefaultPluginContext.
func (p *noopPluginContext) NewHttpContext(uint32) types.HttpContext {
	if !p.tcp {
		return &noopHttpContext{}
	}
	return nil
}

// NewTcpContext implements the same method on types.DefaultPluginContext.
func (p *noopPluginContext) NewTcpContext(uint32) types.TcpContext {
	if p.tcp {
		return &noopTcpContext{}
	}
	return nil
}

type noopHttpContext struct {
	types.DefaultHttpContext
}

type noopTcpContext struct {
	types.DefaultTcpContext
}

// Execute lifecycle methods, there should be no logs for the no-op plugin.
func TestTimingOff(t *testing.T) {
	t.Run("http", func(t *testing.T) {
		host, reset := NewHostEmulator(NewEmulatorOption().WithVMContext(&noopPlugin{}))
		defer reset()

		require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())
		id := host.InitializeHttpContext()

		require.Equal(t, types.ActionContinue, host.CallOnRequestHeaders(id, nil, false))
		require.Equal(t, types.ActionContinue, host.CallOnRequestBody(id, nil, false))
		require.Equal(t, types.ActionContinue, host.CallOnRequestTrailers(id, nil))
		require.Equal(t, types.ActionContinue, host.CallOnResponseHeaders(id, nil, false))
		require.Equal(t, types.ActionContinue, host.CallOnResponseBody(id, nil, false))
		require.Equal(t, types.ActionContinue, host.CallOnResponseTrailers(id, nil))
		host.CompleteHttpContext(id)

		require.Empty(t, host.GetDebugLogs())
	})

	t.Run("tcp", func(t *testing.T) {
		host, reset := NewHostEmulator(NewEmulatorOption().WithVMContext(&noopPlugin{tcp: true}))
		defer reset()

		require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())
		id, action := host.InitializeConnection()
		require.Equal(t, types.ActionContinue, action)

		require.Equal(t, types.ActionContinue, host.CallOnDownstreamData(id, nil))
		require.Equal(t, types.ActionContinue, host.CallOnUpstreamData(id, nil))
		host.CompleteConnection(id)

		require.Empty(t, host.GetDebugLogs())
	})
}
