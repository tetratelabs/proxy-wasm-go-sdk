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
	"reflect"
	"unsafe"
)

// thin wrappers of raw host calls

func setMap(mapType MapType, headers [][2]string) Status {
	shs := serializeMap(headers)
	hp := &shs[0]
	hl := len(shs)
	return proxySetHeaderMapPairs(mapType, hp, hl)
}

// TODO: not tested yet
func getMapValue(mapType MapType, key string) (string, Status) {
	var rvs int
	var raw *byte
	if st := proxyGetHeaderMapValue(mapType, stringToBytePtr(key), len(key), &raw, &rvs); st != StatusOk {
		return "", st
	}

	ret := *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(raw)),
		Len:  uintptr(rvs),
		Cap:  uintptr(rvs),
	}))
	return ret, StatusOk
}

// TODO: not tested yet
func removeMapValue(mapType MapType, key string) Status {
	return proxyRemoveHeaderMapValue(mapType, stringToBytePtr(key), len(key))
}

// TODO: not tested yet
func setMapValue(mapType MapType, key, value string) Status {
	return proxyReplaceHeaderMapValue(mapType, stringToBytePtr(key), len(key), stringToBytePtr(value), len(value))
}

// TODO: not tested yet
func addMapValue(mapType MapType, key, value string) Status {
	return proxyAddHeaderMapValue(mapType, stringToBytePtr(key), len(key), stringToBytePtr(value), len(value))
}

func getMap(mapType MapType) ([][2]string, Status) {
	var rvs int
	var raw *byte

	st := proxyGetHeaderMapPairs(mapType, &raw, &rvs)
	if st != StatusOk {
		return nil, st
	}

	bs := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(raw)),
		Len:  uintptr(rvs),
		Cap:  uintptr(rvs),
	}))
	return deserializeMap(bs), StatusOk
}

func getBuffer(bufType BufferType, start, maxSize int) ([]byte, Status) {
	var retData *byte
	var retSize int
	switch st := proxyGetBufferBytes(bufType, start, maxSize, &retData, &retSize); st {
	case StatusOk:
		// is this correct handling...?
		if retData == nil {
			return nil, StatusNotFound
		}
		return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
			Data: uintptr(unsafe.Pointer(retData)),
			Len:  uintptr(retSize),
			Cap:  uintptr(retSize),
		})), st
	default:
		return nil, st
	}
}

// TODO: not tested yet
func getConfiguration() ([]byte, Status) {
	var retData *byte
	var retSize int
	switch st := proxyGetConfiguration(&retData, &retSize); st {
	case StatusOk:
		// is this correct handling...?
		if retData == nil {
			return nil, StatusNotFound
		}
		return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
			Data: uintptr(unsafe.Pointer(retData)),
			Len:  uintptr(retSize),
			Cap:  uintptr(retSize),
		})), st
	default:
		return nil, st
	}
}

func sendHttpResponse(statusCode uint32, headers [][2]string, body string) Status {
	shs := serializeMap(headers)
	hp := &shs[0]
	hl := len(shs)
	return proxySendLocalResponse(statusCode, nil, 0,
		stringToBytePtr(body), len(body), hp, hl, -1,
	)
}

func setEffectiveContext(contextID uint32) Status {
	return proxySetEffectiveContext(contextID)
}

func dispatchHttpCall(upstream string,
	headers [][2]string, body string, trailers [][2]string, timeoutMillisecond uint32) (uint32, Status) {
	shs := serializeMap(headers)
	hp := &shs[0]
	hl := len(shs)

	sts := serializeMap(trailers)
	tp := &sts[0]
	tl := len(sts)

	var calloutID uint32

	u := []byte(upstream)
	switch retStatus := proxyHttpCall(&u[0], len(u),
		hp, hl, stringToBytePtr(body), len(body), tp, tl, timeoutMillisecond, &calloutID); retStatus {
	case StatusOk:
		currentState.registerCallout(calloutID)
		return calloutID, StatusOk
	default:
		return 0, retStatus
	}
}

func setTickPeriodMilliSeconds(millSec uint32) Status {
	return proxySetTickPeriodMilliseconds(millSec)
}

func stringToBytePtr(in string) *byte {
	var ret *byte
	if len(in) > 0 {
		b := []byte(in) // TODO: zero alloc
		ret = &b[0]
	}
	return ret
}
