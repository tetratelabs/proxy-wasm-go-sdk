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

// +build proxytest

package hostcall

import "github.com/mathetake/proxy-wasm-go/runtime/types"

func ProxyLog(logLevel types.LogLevel, messageData *byte, messageSize int) types.Status {
	return 0
}

func ProxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int) {}

func ProxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int) {}

func ProxySendLocalResponse(statusCode uint32, statusCodeDetailData *byte, statusCodeDetailsSize int, bodyData *byte, bodySize int, headersData *byte, headersSize int, grpcStatus int32) types.Status {
	return 0
}
func ProxyGetSharedData(keyData *byte, keySize int, returnValueData **byte, returnValueSize *byte, returnCas *uint32) types.Status {
	return 0
}
func ProxySetSharedData(keyData *byte, keySize int, valueData *byte, valueSize int, cas uint32) types.Status {
	return 0
}
func ProxyRegisterSharedQueue(nameData *byte, nameSize uint, returnID *uint32) types.Status {
	return 0
}
func ProxyResolveSharedQueue(vmIDData *byte, vmIDSize uint, nameData *byte, nameSize uint, returnID *uint32) types.Status {
	return 0
}
func ProxyDequeueSharedQueue(queueID uint32, returnValueData **byte, returnValueSize *byte) types.Status {
	return 0
}
func ProxyEnqueueSharedQueue(queueID uint32, valueData *byte, valueSize uint) types.Status {
	return 0
}
func ProxyGetHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) types.Status {
	return 0
}
func ProxyAddHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status {
	return 0
}
func ProxyReplaceHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status {
	return 0
}
func ProxyContinueStream(streamType types.StreamType) types.Status {
	return 0
}
func ProxyCloseStream(streamType types.StreamType) types.Status {
	return 0
}
func ProxyRemoveHeaderMapValue(mapType types.MapType, keyData *byte, keySize int) types.Status {
	return 0
}
func ProxyGetHeaderMapPairs(mapType types.MapType, returnValueData **byte, returnValueSize *int) types.Status {
	return 0
}
func ProxySetHeaderMapPairs(mapType types.MapType, mapData *byte, mapSize int) types.Status {
	return 0
}
func ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) types.Status {
	return 0
}
func ProxyHttpCall(upstreamData *byte, upstreamSize int, headerData *byte, headerSize int, bodyData *byte, bodySize int, trailersData *byte, trailersSize int, timeout uint32, calloutIDPtr *uint32) types.Status {
	return 0
}
func ProxySetTickPeriodMilliseconds(period uint32) types.Status {
	return 0
}
func ProxyGetCurrentTimeNanoseconds(returnTime *int64) types.Status {
	return 0
}
func ProxySetEffectiveContext(contextID uint32) types.Status {
	return 0
}
func ProxyDone() types.Status {
	return 0
}
