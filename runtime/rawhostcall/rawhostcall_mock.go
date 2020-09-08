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

package rawhostcall

import "github.com/mathetake/proxy-wasm-go/runtime/types"

var currentHost ProxyWASMHost

func RegisterMockWASMHost(host ProxyWASMHost) {
	currentHost = host
}

type ProxyWASMHost interface {
	ProxyLog(logLevel types.LogLevel, messageData *byte, messageSize int) types.Status
	ProxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int)
	ProxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int)
	ProxySendLocalResponse(statusCode uint32, statusCodeDetailData *byte, statusCodeDetailsSize int, bodyData *byte, bodySize int, headersData *byte, headersSize int, grpcStatus int32) types.Status
	ProxyGetSharedData(keyData *byte, keySize int, returnValueData **byte, returnValueSize *byte, returnCas *uint32) types.Status
	ProxySetSharedData(keyData *byte, keySize int, valueData *byte, valueSize int, cas uint32) types.Status
	ProxyRegisterSharedQueue(nameData *byte, nameSize int, returnID *uint32) types.Status
	ProxyResolveSharedQueue(vmIDData *byte, vmIDSize int, nameData *byte, nameSize int, returnID *uint32) types.Status
	ProxyDequeueSharedQueue(queueID uint32, returnValueData **byte, returnValueSize *int) types.Status
	ProxyEnqueueSharedQueue(queueID uint32, valueData *byte, valueSize int) types.Status
	ProxyGetHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) types.Status
	ProxyAddHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status
	ProxyReplaceHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status
	ProxyContinueStream(streamType types.StreamType) types.Status
	ProxyCloseStream(streamType types.StreamType) types.Status
	ProxyRemoveHeaderMapValue(mapType types.MapType, keyData *byte, keySize int) types.Status
	ProxyGetHeaderMapPairs(mapType types.MapType, returnValueData **byte, returnValueSize *int) types.Status
	ProxySetHeaderMapPairs(mapType types.MapType, mapData *byte, mapSize int) types.Status
	ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) types.Status
	ProxyHttpCall(upstreamData *byte, upstreamSize int, headerData *byte, headerSize int, bodyData *byte, bodySize int, trailersData *byte, trailersSize int, timeout uint32, calloutIDPtr *uint32) types.Status
	ProxySetTickPeriodMilliseconds(period uint32) types.Status
	ProxyGetCurrentTimeNanoseconds(returnTime *int64) types.Status
	ProxySetEffectiveContext(contextID uint32) types.Status
	ProxyDone() types.Status
	ProxyDefineMetric(metricType types.MetricType, metricNameData *byte, metricNameSize int, returnMetricIDPtr *uint32) types.Status
	ProxyIncrementMetric(metricID uint32, offset int64) types.Status
	ProxyRecordMetric(metricID uint32, value uint64) types.Status
	ProxyGetMetric(metricID uint32, returnMetricValue *uint64) types.Status
}

func ProxyLog(logLevel types.LogLevel, messageData *byte, messageSize int) types.Status {
	return currentHost.ProxyLog(logLevel, messageData, messageSize)
}

func ProxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int) {
	currentHost.ProxySetProperty(pathData, pathSize, valueData, valueSize)
}

func ProxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int) {
	currentHost.ProxyGetProperty(pathData, pathSize, returnValueData, returnValueSize)
}

func ProxySendLocalResponse(statusCode uint32, statusCodeDetailData *byte,
	statusCodeDetailsSize int, bodyData *byte, bodySize int, headersData *byte, headersSize int, grpcStatus int32) types.Status {
	return currentHost.ProxySendLocalResponse(statusCode,
		statusCodeDetailData, statusCodeDetailsSize, bodyData, bodySize, headersData, headersSize, grpcStatus)
}

func ProxyGetSharedData(keyData *byte, keySize int, returnValueData **byte, returnValueSize *byte, returnCas *uint32) types.Status {
	return currentHost.ProxyGetSharedData(keyData, keySize, returnValueData, returnValueSize, returnCas)
}

func ProxySetSharedData(keyData *byte, keySize int, valueData *byte, valueSize int, cas uint32) types.Status {
	return currentHost.ProxySetSharedData(keyData, keySize, valueData, valueSize, cas)
}

func ProxyRegisterSharedQueue(nameData *byte, nameSize int, returnID *uint32) types.Status {
	return currentHost.ProxyRegisterSharedQueue(nameData, nameSize, returnID)
}

func ProxyResolveSharedQueue(vmIDData *byte, vmIDSize int, nameData *byte, nameSize int, returnID *uint32) types.Status {
	return currentHost.ProxyResolveSharedQueue(vmIDData, vmIDSize, nameData, nameSize, returnID)
}

func ProxyDequeueSharedQueue(queueID uint32, returnValueData **byte, returnValueSize *int) types.Status {
	return currentHost.ProxyDequeueSharedQueue(queueID, returnValueData, returnValueSize)
}

func ProxyEnqueueSharedQueue(queueID uint32, valueData *byte, valueSize int) types.Status {
	return currentHost.ProxyEnqueueSharedQueue(queueID, valueData, valueSize)
}

func ProxyGetHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) types.Status {
	return currentHost.ProxyGetHeaderMapValue(mapType, keyData, keySize, returnValueData, returnValueSize)
}

func ProxyAddHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status {
	return currentHost.ProxyAddHeaderMapValue(mapType, keyData, keySize, valueData, valueSize)
}

func ProxyReplaceHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status {
	return currentHost.ProxyReplaceHeaderMapValue(mapType, keyData, keySize, valueData, valueSize)
}

func ProxyContinueStream(streamType types.StreamType) types.Status {
	return currentHost.ProxyContinueStream(streamType)
}

func ProxyCloseStream(streamType types.StreamType) types.Status {
	return currentHost.ProxyCloseStream(streamType)
}
func ProxyRemoveHeaderMapValue(mapType types.MapType, keyData *byte, keySize int) types.Status {
	return currentHost.ProxyRemoveHeaderMapValue(mapType, keyData, keySize)
}

func ProxyGetHeaderMapPairs(mapType types.MapType, returnValueData **byte, returnValueSize *int) types.Status {
	return currentHost.ProxyGetHeaderMapPairs(mapType, returnValueData, returnValueSize)
}

func ProxySetHeaderMapPairs(mapType types.MapType, mapData *byte, mapSize int) types.Status {
	return currentHost.ProxySetHeaderMapPairs(mapType, mapData, mapSize)
}

func ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) types.Status {
	return currentHost.ProxyGetBufferBytes(bt, start, maxSize, returnBufferData, returnBufferSize)
}

func ProxyHttpCall(upstreamData *byte, upstreamSize int, headerData *byte, headerSize int, bodyData *byte,
	bodySize int, trailersData *byte, trailersSize int, timeout uint32, calloutIDPtr *uint32) types.Status {
	return currentHost.ProxyHttpCall(upstreamData, upstreamSize,
		headerData, headerSize, bodyData, bodySize, trailersData, trailersSize, timeout, calloutIDPtr)
}

func ProxySetTickPeriodMilliseconds(period uint32) types.Status {
	return currentHost.ProxySetTickPeriodMilliseconds(period)
}

func ProxyGetCurrentTimeNanoseconds(returnTime *int64) types.Status {
	return currentHost.ProxyGetCurrentTimeNanoseconds(returnTime)
}

func ProxySetEffectiveContext(contextID uint32) types.Status {
	return currentHost.ProxySetEffectiveContext(contextID)
}

func ProxyDone() types.Status {
	return currentHost.ProxyDone()
}

func ProxyDefineMetric(metricType types.MetricType,
	metricNameData *byte, metricNameSize int, returnMetricIDPtr *uint32) types.Status {
	return currentHost.ProxyDefineMetric(metricType, metricNameData, metricNameSize, returnMetricIDPtr)
}

func ProxyIncrementMetric(metricID uint32, offset int64) types.Status {
	return currentHost.ProxyIncrementMetric(metricID, offset)
}

func ProxyRecordMetric(metricID uint32, value uint64) types.Status {
	return currentHost.ProxyRecordMetric(metricID, value)
}

func ProxyGetMetric(metricID uint32, returnMetricValue *uint64) types.Status {
	return currentHost.ProxyGetMetric(metricID, returnMetricValue)
}
