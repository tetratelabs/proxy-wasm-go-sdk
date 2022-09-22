//go:build proxywasm_timing

package proxytest

import (
	"fmt"
	"regexp"
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

// Execute lifecycle methods, there should be logs for the no-op plugin.
func TestTimingOn(t *testing.T) {
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

		host.GetDebugLogs()
		requireLogged(t, "proxyOnContextCreate", host.GetDebugLogs())
		requireLogged(t, "proxyOnConfigure", host.GetDebugLogs())
		requireLogged(t, "proxyOnContextCreate", host.GetDebugLogs())
		requireLogged(t, "proxyOnRequestHeaders", host.GetDebugLogs())
		requireLogged(t, "proxyOnRequestBody", host.GetDebugLogs())
		requireLogged(t, "proxyOnRequestTrailers", host.GetDebugLogs())
		requireLogged(t, "proxyOnResponseHeaders", host.GetDebugLogs())
		requireLogged(t, "proxyOnResponseBody", host.GetDebugLogs())
		requireLogged(t, "proxyOnResponseTrailers", host.GetDebugLogs())
		requireLogged(t, "proxyOnLog", host.GetDebugLogs())
		requireLogged(t, "proxyOnDelete", host.GetDebugLogs())
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

		requireLogged(t, "proxyOnContextCreate", host.GetDebugLogs())
		requireLogged(t, "proxyOnConfigure", host.GetDebugLogs())
		requireLogged(t, "proxyOnContextCreate", host.GetDebugLogs())
		requireLogged(t, "proxyOnNewConnection", host.GetDebugLogs())
		requireLogged(t, "proxyOnDownstreamData", host.GetDebugLogs())
		requireLogged(t, "proxyOnUpstreamData", host.GetDebugLogs())
		requireLogged(t, "proxyOnLog", host.GetDebugLogs())
		requireLogged(t, "proxyOnDelete", host.GetDebugLogs())
	})
}

func requireLogged(t *testing.T, msg string, logs []string) {
	t.Helper()
	re := regexp.MustCompile(fmt.Sprintf(`%s took \d+`, msg))
	for _, l := range logs {
		if re.MatchString(l) {
			// While we can't make reliable assertions on the actual time took, it is inconceivable to have a time of
			// 0, though bugs like forgetting "defer" might cause it. So do a special check for it.
			require.NotContains(t, l, "took 0s")
			return
		}
	}
	require.Failf(t, "unexpected log", "expected log %s, got %v", msg, logs)
}
