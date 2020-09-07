// Copyright 2020 Takeshi Yoneda(@mathetake)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runtime

import (
	"github.com/mathetake/proxy-wasm-go/runtime/hostcall"
	"github.com/mathetake/proxy-wasm-go/runtime/types"
)

// thin wrappers of raw host calls

func setMap(mapType types.MapType, headers [][2]string) types.Status {
	shs := serializeMap(headers)
	hp := &shs[0]
	hl := len(shs)
	return hostcall.ProxySetHeaderMapPairs(mapType, hp, hl)
}

// TODO: not tested yet
func getMapValue(mapType types.MapType, key string) (string, types.Status) {
	var rvs int
	var raw *byte
	if st := hostcall.ProxyGetHeaderMapValue(mapType, stringToBytePtr(key), len(key), &raw, &rvs); st != types.StatusOk {
		return "", st
	}

	ret := rawBytePtrToString(raw, rvs)
	return ret, types.StatusOk
}

// TODO: not tested yet
func removeMapValue(mapType types.MapType, key string) types.Status {
	return hostcall.ProxyRemoveHeaderMapValue(mapType, stringToBytePtr(key), len(key))
}

// TODO: not tested yet
func setMapValue(mapType types.MapType, key, value string) types.Status {
	return hostcall.ProxyReplaceHeaderMapValue(mapType, stringToBytePtr(key), len(key), stringToBytePtr(value), len(value))
}

// TODO: not tested yet
func addMapValue(mapType types.MapType, key, value string) types.Status {
	return hostcall.ProxyAddHeaderMapValue(mapType, stringToBytePtr(key), len(key), stringToBytePtr(value), len(value))
}

func getMap(mapType types.MapType) ([][2]string, types.Status) {
	var rvs int
	var raw *byte

	st := hostcall.ProxyGetHeaderMapPairs(mapType, &raw, &rvs)
	if st != types.StatusOk {
		return nil, st
	}

	bs := rawBytePtrToByteSlice(raw, rvs)
	return deserializeMap(bs), types.StatusOk
}

func getBuffer(bufType types.BufferType, start, maxSize int) ([]byte, types.Status) {
	var retData *byte
	var retSize int
	switch st := hostcall.ProxyGetBufferBytes(bufType, start, maxSize, &retData, &retSize); st {
	case types.StatusOk:
		// is this correct handling...?
		if retData == nil {
			return nil, types.StatusNotFound
		}
		return rawBytePtrToByteSlice(retData, retSize), st
	default:
		return nil, st
	}
}

func sendHttpResponse(statusCode uint32, headers [][2]string, body string) types.Status {
	shs := serializeMap(headers)
	hp := &shs[0]
	hl := len(shs)
	return hostcall.ProxySendLocalResponse(statusCode, nil, 0,
		stringToBytePtr(body), len(body), hp, hl, -1,
	)
}

func setEffectiveContext(contextID uint32) types.Status {
	return hostcall.ProxySetEffectiveContext(contextID)
}

func dispatchHttpCall(upstream string,
	headers [][2]string, body string, trailers [][2]string, timeoutMillisecond uint32) (uint32, types.Status) {
	shs := serializeMap(headers)
	hp := &shs[0]
	hl := len(shs)

	sts := serializeMap(trailers)
	tp := &sts[0]
	tl := len(sts)

	var calloutID uint32

	u := []byte(upstream)
	switch retStatus := hostcall.ProxyHttpCall(&u[0], len(u),
		hp, hl, stringToBytePtr(body), len(body), tp, tl, timeoutMillisecond, &calloutID); retStatus {
	case types.StatusOk:
		currentState.registerCallout(calloutID)
		return calloutID, types.StatusOk
	default:
		return 0, retStatus
	}
}

func setTickPeriodMilliSeconds(millSec uint32) types.Status {
	return hostcall.ProxySetTickPeriodMilliseconds(millSec)
}

func stringToBytePtr(in string) *byte {
	var ret *byte
	if len(in) > 0 {
		b := []byte(in) // TODO: zero alloc
		ret = &b[0]
	}
	return ret
}
