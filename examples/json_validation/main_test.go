package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestOnHTTPRequestHeaders(t *testing.T) {
	type testCase struct {
		contentType    string
		expectedAction types.Action
	}

	for name, tCase := range map[string]testCase{
		"fails due to unsupported content type": {
			contentType:    "text/html",
			expectedAction: types.ActionPause,
		},
		"success for JSON": {
			contentType:    "application/json",
			expectedAction: types.ActionContinue,
		},
	} {
		t.Run(name, func(t *testing.T) {
			opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
			host, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())

			id := host.InitializeHttpContext()

			hs := [][2]string{{"content-type", tCase.contentType}}

			action := host.CallOnRequestHeaders(id, hs, false)
			assert.Equal(t, tCase.expectedAction, action)
		})
	}
}

func TestOnHTTPRequestBody(t *testing.T) {
	type testCase struct {
		body           string
		expectedAction types.Action
	}

	for name, tCase := range map[string]testCase{
		"pauses due to invalid payload": {
			body:           "invalid_payload",
			expectedAction: types.ActionPause,
		},
		"pauses due to unknown keys": {
			body:           `{"unknown_key":"unknown_value"}`,
			expectedAction: types.ActionPause,
		},
		"success": {
			body:           "{\"my_key\":\"my_value\"}",
			expectedAction: types.ActionContinue,
		},
	} {
		t.Run(name, func(t *testing.T) {
			opt := proxytest.
				NewEmulatorOption().
				WithPluginConfiguration([]byte(`{"requiredKeys": ["my_key"]}`)).
				WithVMContext(&vmContext{})
			host, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())

			id := host.InitializeHttpContext()

			action := host.CallOnRequestBody(id, []byte(tCase.body), true)
			assert.Equal(t, tCase.expectedAction, action)
		})
	}
}
