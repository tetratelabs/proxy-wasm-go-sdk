package proxywasm

import (
	"math"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func setMap(mapType internal.MapType, headers [][2]string) error {
	shs := internal.SerializeMap(headers)
	hp := &shs[0]
	hl := len(shs)
	return internal.StatusToError(internal.ProxySetHeaderMapPairs(mapType, hp, hl))
}

func getMapValue(mapType internal.MapType, key string) (string, error) {
	var rvs int
	var raw *byte
	if st := internal.ProxyGetHeaderMapValue(
		mapType, internal.StringBytePtr(key), len(key), &raw, &rvs,
	); st != internal.StatusOK {
		return "", internal.StatusToError(st)
	}

	ret := internal.RawBytePtrToString(raw, rvs)
	return ret, nil
}

func removeMapValue(mapType internal.MapType, key string) error {
	return internal.StatusToError(
		internal.ProxyRemoveHeaderMapValue(mapType, internal.StringBytePtr(key), len(key)),
	)
}

func replaceMapValue(mapType internal.MapType, key, value string) error {
	return internal.StatusToError(
		internal.ProxyReplaceHeaderMapValue(
			mapType, internal.StringBytePtr(key), len(key), internal.StringBytePtr(value), len(value),
		),
	)
}

func addMapValue(mapType internal.MapType, key, value string) error {
	return internal.StatusToError(
		internal.ProxyAddHeaderMapValue(
			mapType, internal.StringBytePtr(key), len(key), internal.StringBytePtr(value), len(value),
		),
	)
}

func getMap(mapType internal.MapType) ([][2]string, error) {
	var rvs int
	var raw *byte

	st := internal.ProxyGetHeaderMapPairs(mapType, &raw, &rvs)
	if st != internal.StatusOK {
		return nil, internal.StatusToError(st)
	}

	bs := internal.RawBytePtrToByteSlice(raw, rvs)
	return internal.DeserializeMap(bs), nil
}

func getBuffer(bufType internal.BufferType, start, maxSize int) ([]byte, error) {
	var retData *byte
	var retSize int
	switch st := internal.ProxyGetBufferBytes(bufType, start, maxSize, &retData, &retSize); st {
	case internal.StatusOK:
		if retData == nil {
			return nil, types.ErrorStatusNotFound
		}
		return internal.RawBytePtrToByteSlice(retData, retSize), nil
	default:
		return nil, internal.StatusToError(st)
	}
}

func appendToBuffer(bufType internal.BufferType, buffer []byte) error {
	var bufferData *byte
	if len(buffer) != 0 {
		bufferData = &buffer[0]
	}
	return internal.StatusToError(internal.ProxySetBufferBytes(bufType, math.MaxInt32, 0, bufferData, len(buffer)))
}

func prependToBuffer(bufType internal.BufferType, buffer []byte) error {
	var bufferData *byte
	if len(buffer) != 0 {
		bufferData = &buffer[0]
	}
	return internal.StatusToError(internal.ProxySetBufferBytes(bufType, 0, 0, bufferData, len(buffer)))
}

func replaceBuffer(bufType internal.BufferType, buffer []byte) error {
	var bufferData *byte
	if len(buffer) != 0 {
		bufferData = &buffer[0]
	}
	return internal.StatusToError(
		internal.ProxySetBufferBytes(bufType, 0, math.MaxInt32, bufferData, len(buffer)),
	)
}
