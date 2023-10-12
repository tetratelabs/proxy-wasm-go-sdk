package properties

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
)

func TestGetPropertyBool(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty([]string{"someBoolPath"}, serializeBool(true))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := getPropertyBool([]string{"someBoolPath"})
	require.NoError(t, err)
	require.Equal(t, true, result)
}

func TestGetPropertyByteSliceMap(t *testing.T) {
	input := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
	}

	opt := proxytest.NewEmulatorOption().WithProperty([]string{"someByteSliceMapPath"}, serializeByteSliceMap(input))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := getPropertyByteSliceMap([]string{"someByteSliceMapPath"})
	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestGetPropertyByteSliceSlice(t *testing.T) {
	input := [][]byte{
		[]byte("value1"),
		[]byte("value2"),
	}

	opt := proxytest.NewEmulatorOption().WithProperty([]string{"someByteSliceSlicePath"}, serializeByteSliceSlice(input))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := getPropertyByteSliceSlice([]string{"someByteSliceSlicePath"})
	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestGetPropertyFloat64(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty([]string{"someFloat64Path"}, serializeFloat64(3.14))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := getPropertyFloat64([]string{"someFloat64Path"})
	require.NoError(t, err)
	require.Equal(t, 3.14, result)
}

func TestGetPropertyString(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty([]string{"someStringPath"}, []byte("testString"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := getPropertyString([]string{"someStringPath"})
	require.NoError(t, err)
	require.Equal(t, "testString", result)
}

func TestGetPropertyStringMap(t *testing.T) {
	input := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	opt := proxytest.NewEmulatorOption().WithProperty([]string{"someStringMapPath"}, serializeStringMap(input))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := getPropertyStringMap([]string{"someStringMapPath"})
	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestGetPropertyStringSlice(t *testing.T) {
	input := []string{"value1", "value2"}

	opt := proxytest.NewEmulatorOption().WithProperty([]string{"someStringSlicePath"}, serializeStringSlice(input))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := getPropertyStringSlice([]string{"someStringSlicePath"})
	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestGetPropertyTimestamp(t *testing.T) {
	now := time.Now().UTC()

	opt := proxytest.NewEmulatorOption().WithProperty([]string{"someTimestampPath"}, serializeTimestamp(now))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := getPropertyTimestamp([]string{"someTimestampPath"})
	require.NoError(t, err)
	require.Equal(t, now, result)
}

func TestGetPropertyUint64(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty([]string{"someUint64Path"}, serializeUint64(12345))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := getPropertyUint64([]string{"someUint64Path"})
	require.NoError(t, err)
	require.Equal(t, uint64(12345), result)
}
