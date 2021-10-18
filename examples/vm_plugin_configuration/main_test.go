// These tests are supposed to run with `proxytest` build tag, and this way we can leverage the testing framework in "proxytest" package.
// The framework emulates the expected behavior of Envoyproxy, and you can test your extensions without running Envoy and with
// the standard Go CLI. To run tests, simply run
// go test -tags=proxytest ./...

//go:build proxytest

package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestContext_OnPluginStart(t *testing.T) {
	// Setup configurations.
	pluginConfigData := `tinygo plugin configuration`
	opt := proxytest.NewEmulatorOption().
		WithPluginConfiguration([]byte(pluginConfigData)).
		WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnPluginStart.
	require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())

	// Check Envoy logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, "plugin config: "+pluginConfigData)
}

func TestContext_OnVMStart(t *testing.T) {
	// Setup configurations.
	vmConfigData := `tinygo vm configuration`
	opt := proxytest.NewEmulatorOption().
		WithVMConfiguration([]byte(vmConfigData)).
		WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnVMStart.
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())

	// Check Envoy logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, "vm config: "+vmConfigData)
}
