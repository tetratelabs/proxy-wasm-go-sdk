//go:build proxywasm_timing

package proxytest

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type timedPlugin struct {
	types.DefaultVMContext
	tcp bool
}

// NewPluginContext implements the same method on types.DefaultVMContext.
func (p *timedPlugin) NewPluginContext(uint32) types.PluginContext {
	time.Sleep(1 * time.Millisecond)
	return &timedPluginContext{tcp: p.tcp}
}

type timedPluginContext struct {
	types.DefaultPluginContext
	tcp bool
}

// OnPluginStart implements the same method on types.PluginContext.
func (p *timedPluginContext) OnPluginStart(int) types.OnPluginStartStatus {
	time.Sleep(1 * time.Millisecond)
	return true
}

// NewHttpContext implements the same method on types.DefaultPluginContext.
func (p *timedPluginContext) NewHttpContext(uint32) types.HttpContext {
	time.Sleep(1 * time.Millisecond)
	if !p.tcp {
		return &timedHttpContext{}
	}
	return nil
}

// NewTcpContext implements the same method on types.DefaultPluginContext.
func (p *timedPluginContext) NewTcpContext(uint32) types.TcpContext {
	time.Sleep(1 * time.Millisecond)
	if p.tcp {
		return &timedTcpContext{}
	}
	return nil
}

type timedHttpContext struct {
}

// OnHttpRequestHeaders implements the same method on types.HttpContext.
func (c *timedHttpContext) OnHttpRequestHeaders(int, bool) types.Action {
	time.Sleep(1 * time.Millisecond)
	return types.ActionContinue
}

// OnHttpRequestBody implements the same method on types.HttpContext.
func (c *timedHttpContext) OnHttpRequestBody(int, bool) types.Action {
	time.Sleep(1 * time.Millisecond)
	return types.ActionContinue
}

// OnHttpRequestTrailers implements the same method on types.HttpContext.
func (c *timedHttpContext) OnHttpRequestTrailers(int) types.Action {
	time.Sleep(1 * time.Millisecond)
	return types.ActionContinue
}

// OnHttpResponseHeaders implements the same method on types.HttpContext.
func (c *timedHttpContext) OnHttpResponseHeaders(int, bool) types.Action {
	time.Sleep(1 * time.Millisecond)
	return types.ActionContinue
}

// OnHttpResponseBody implements the same method on types.HttpContext.
func (c *timedHttpContext) OnHttpResponseBody(int, bool) types.Action {
	time.Sleep(1 * time.Millisecond)
	return types.ActionContinue
}

// OnHttpResponseTrailers implements the same method on types.HttpContext.
func (c *timedHttpContext) OnHttpResponseTrailers(int) types.Action {
	time.Sleep(1 * time.Millisecond)
	return types.ActionContinue
}

// OnHttpStreamDone implements the same method on types.HttpContext.
func (c *timedHttpContext) OnHttpStreamDone() {
	time.Sleep(1 * time.Millisecond)
}

type timedTcpContext struct {
}

// OnNewConnection implements the same method on types.TcpContext.
func (t timedTcpContext) OnNewConnection() types.Action {
	time.Sleep(1 * time.Millisecond)
	return types.ActionContinue
}

// OnDownstreamData implements the same method on types.TcpContext.
func (t timedTcpContext) OnDownstreamData(int, bool) types.Action {
	time.Sleep(1 * time.Millisecond)
	return types.ActionContinue
}

// OnDownstreamClose implements the same method on types.TcpContext.
func (t timedTcpContext) OnDownstreamClose(types.PeerType) {
	time.Sleep(1 * time.Millisecond)
}

// OnUpstreamData implements the same method on types.TcpContext.
func (t timedTcpContext) OnUpstreamData(int, bool) types.Action {
	time.Sleep(1 * time.Millisecond)
	return types.ActionContinue
}

// OnUpstreamClose implements the same method on types.TcpContext.
func (t timedTcpContext) OnUpstreamClose(types.PeerType) {
	time.Sleep(1 * time.Millisecond)
}

// OnStreamDone implements the same method on types.TcpContext.
func (t timedTcpContext) OnStreamDone() {
	time.Sleep(1 * time.Millisecond)
}

// Execute lifecycle methods, there should be logs for the no-op plugin.
func TestTimingOn(t *testing.T) {
	t.Run("http", func(t *testing.T) {
		host, reset := NewHostEmulator(NewEmulatorOption().WithVMContext(&timedPlugin{}))
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
		host, reset := NewHostEmulator(NewEmulatorOption().WithVMContext(&timedPlugin{tcp: true}))
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
			// 0 since we add a small sleep in each function, while bugs like forgetting "defer" might cause it. So do
			// a special check for it.
			require.NotContains(t, l, "took 0s")
			return
		}
	}
	require.Failf(t, "unexpected log", "expected log %s, got %v", msg, logs)
}
