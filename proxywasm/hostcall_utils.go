package proxywasm

import (
	"math"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal/rawhostcall"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func setMap(mapType types.MapType, headers [][2]string) types.Status {
	shs := internal.SerializeMap(headers)
	hp := &shs[0]
	hl := len(shs)
	return rawhostcall.ProxySetHeaderMapPairs(mapType, hp, hl)
}

func getMapValue(mapType types.MapType, key string) (string, types.Status) {
	var rvs int
	var raw *byte
	if st := rawhostcall.ProxyGetHeaderMapValue(mapType, internal.StringBytePtr(key), len(key), &raw, &rvs); st != types.StatusOK {
		return "", st
	}

	ret := internal.RawBytePtrToString(raw, rvs)
	return ret, types.StatusOK
}

func removeMapValue(mapType types.MapType, key string) types.Status {
	return rawhostcall.ProxyRemoveHeaderMapValue(mapType, internal.StringBytePtr(key), len(key))
}

func setMapValue(mapType types.MapType, key, value string) types.Status {
	return rawhostcall.ProxyReplaceHeaderMapValue(mapType, internal.StringBytePtr(key), len(key), internal.StringBytePtr(value), len(value))
}

func addMapValue(mapType types.MapType, key, value string) types.Status {
	return rawhostcall.ProxyAddHeaderMapValue(mapType, internal.StringBytePtr(key), len(key), internal.StringBytePtr(value), len(value))
}

func getMap(mapType types.MapType) ([][2]string, types.Status) {
	var rvs int
	var raw *byte

	st := rawhostcall.ProxyGetHeaderMapPairs(mapType, &raw, &rvs)
	if st != types.StatusOK {
		return nil, st
	}

	bs := internal.RawBytePtrToByteSlice(raw, rvs)
	return internal.DeserializeMap(bs), types.StatusOK
}

func getBuffer(bufType types.BufferType, start, maxSize int) ([]byte, types.Status) {
	var retData *byte
	var retSize int
	switch st := rawhostcall.ProxyGetBufferBytes(bufType, start, maxSize, &retData, &retSize); st {
	case types.StatusOK:
		// is this correct handling...?
		if retData == nil {
			return nil, types.StatusNotFound
		}
		return internal.RawBytePtrToByteSlice(retData, retSize), st
	default:
		return nil, st
	}
}

func appendToBuffer(bufType types.BufferType, buffer []byte) error {
	var bufferData *byte
	if len(buffer) != 0 {
		bufferData = &buffer[0]
	}
	return types.StatusToError(rawhostcall.ProxySetBufferBytes(bufType, math.MaxInt32, 0, bufferData, len(buffer)))
}

func prependToBuffer(bufType types.BufferType, buffer []byte) error {
	var bufferData *byte
	if len(buffer) != 0 {
		bufferData = &buffer[0]
	}
	return types.StatusToError(rawhostcall.ProxySetBufferBytes(bufType, 0, 0, bufferData, len(buffer)))
}

func replaceBuffer(bufType types.BufferType, buffer []byte) error {
	var bufferData *byte
	if len(buffer) != 0 {
		bufferData = &buffer[0]
	}
	return types.StatusToError(rawhostcall.ProxySetBufferBytes(bufType, 0, math.MaxInt32, bufferData, len(buffer)))
}
