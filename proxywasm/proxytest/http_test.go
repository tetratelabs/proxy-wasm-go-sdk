package proxytest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type testPlugin struct {
	types.DefaultVMContext
	buffered bool
}

type testPluginContext struct {
	types.DefaultPluginContext
	buffered bool
}

type testHttpContext struct {
	types.DefaultHttpContext
	buffered bool
}

// NewPluginContext implements the same method on types.VMContext.
func (p *testPlugin) NewPluginContext(uint32) types.PluginContext {
	return &testPluginContext{buffered: p.buffered}
}

// NewPluginContext implements the same method on types.PluginContext.
func (p *testPluginContext) NewHttpContext(uint32) types.HttpContext {
	return &testHttpContext{buffered: p.buffered}
}

// OnHttpRequestBody implements the same method on types.HttpContext.
func (h *testHttpContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	if !endOfStream {
		if h.buffered {
			return types.ActionPause
		} else {
			return types.ActionContinue
		}
	}

	body, err := proxywasm.GetHttpRequestBody(0, bodySize)
	if err != nil {
		panic(err)
	}
	proxywasm.LogInfo(fmt.Sprintf("request body:%s", string(body)))

	return types.ActionContinue
}

// OnHttpResponseBody implements the same method on types.HttpContext.
func (h *testHttpContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {
	if !endOfStream {
		if h.buffered {
			return types.ActionPause
		} else {
			return types.ActionContinue
		}
	}

	body, err := proxywasm.GetHttpResponseBody(0, bodySize)
	if err != nil {
		panic(err)
	}
	proxywasm.LogInfo(fmt.Sprintf("response body:%s", string(body)))

	return types.ActionContinue
}

func TestBodyBuffering(t *testing.T) {
	tests := []struct {
		name     string
		buffered bool
		action   types.Action
		logged   string
	}{
		{
			name:     "buffered",
			buffered: true,
			action:   types.ActionPause,
			logged:   "11111",
		},
		{
			name:     "unbuffered",
			buffered: false,
			action:   types.ActionContinue,
			logged:   "22222",
		},
	}

	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			host, reset := NewHostEmulator(NewEmulatorOption().WithVMContext(&testPlugin{buffered: tt.buffered}))
			defer reset()

			id := host.InitializeHttpContext()

			action := host.CallOnRequestBody(id, []byte("11111"), false)
			require.Equal(t, tt.action, action)

			action = host.CallOnRequestBody(id, []byte("22222"), true)
			require.Equal(t, types.ActionContinue, action)

			action = host.CallOnResponseBody(id, []byte("11111"), false)
			require.Equal(t, tt.action, action)

			action = host.CallOnResponseBody(id, []byte("22222"), true)
			require.Equal(t, types.ActionContinue, action)

			logs := host.GetInfoLogs()
			require.Contains(t, logs, fmt.Sprintf("request body:%s", tt.logged))
			require.Contains(t, logs, fmt.Sprintf("response body:%s", tt.logged))
		})
	}
}

func TestProperties(t *testing.T) {
	t.Run("Set and get properties", func(t *testing.T) {
		host, reset := NewHostEmulator(NewEmulatorOption().WithVMContext(&testPlugin{}))
		defer reset()

		_ = host.InitializeHttpContext()

		propertyPath := []string{
			"route_metadata",
			"filter_metadata",
			"envoy.filters.http.wasm",
			"hello",
		}
		propertyData := []byte("world")

		err := host.SetProperty(propertyPath, propertyData)
		require.Equal(t, err, nil)

		data, err := host.GetProperty(propertyPath)
		require.Equal(t, err, nil)
		require.Equal(t, data, propertyData)

		_, err = host.GetProperty([]string{"non-existent path"})
		require.Equal(t, err, internal.StatusToError(internal.StatusNotFound))
	})
}
