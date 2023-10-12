package properties

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
)

func TestGetXdsClusterName(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(xdsClusterName, []byte("outbound|80||httpbin.org"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetXdsClusterName()
	require.NoError(t, err)
	require.Equal(t, "outbound|80||httpbin.org", result)
}

func TestGetXdsRouteName(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(xdsRouteName, []byte("routename"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetXdsRouteName()
	require.NoError(t, err)
	require.Equal(t, "routename", result)
}
func TestGetXdsListenerFilterChainName(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(xdsListenerFilterChainName, []byte("mychain"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetXdsListenerFilterChainName()
	require.NoError(t, err)
	require.Equal(t, "mychain", result)
}
