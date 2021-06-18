// Copyright 2020-2021 Tetrate
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

// TODO: Auto generate this file from rawhostcall.go.
package rawhostcall

import (
	"sync"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

var (
	currentHost ProxyWasmHost
	mutex       = &sync.Mutex{}
)

func RegisterMockWasmHost(host ProxyWasmHost) (release func()) {
	mutex.Lock()
	currentHost = host
	return func() {
		mutex.Unlock()
	}
}

type ProxyWasmHost interface {
	ProxyLog(logLevel uint32, messageData *byte, messageSize int) uint32
	ProxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int) uint32
	ProxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int) uint32
	ProxySendLocalResponse(statusCode uint32, statusCodeDetailData *byte, statusCodeDetailsSize int, bodyData *byte, bodySize int, headersData *byte, headersSize int, grpcStatus int32) uint32
	ProxyGetSharedData(keyData *byte, keySize int, returnValueData **byte, returnValueSize *int, returnCas *uint32) uint32
	ProxySetSharedData(keyData *byte, keySize int, valueData *byte, valueSize int, cas uint32) uint32
	ProxyRegisterSharedQueue(nameData *byte, nameSize int, returnID *uint32) uint32
	ProxyResolveSharedQueue(vmIDData *byte, vmIDSize int, nameData *byte, nameSize int, returnID *uint32) uint32
	ProxyDequeueSharedQueue(queueID uint32, returnValueData **byte, returnValueSize *int) uint32
	ProxyEnqueueSharedQueue(queueID uint32, valueData *byte, valueSize int) uint32
	ProxyGetHeaderMapValue(mapType uint32, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) uint32
	ProxyAddHeaderMapValue(mapType uint32, keyData *byte, keySize int, valueData *byte, valueSize int) uint32
	ProxyReplaceHeaderMapValue(mapType uint32, keyData *byte, keySize int, valueData *byte, valueSize int) uint32
	ProxyContinueStream(streamType uint32) uint32
	ProxyCloseStream(streamType uint32) uint32
	ProxyRemoveHeaderMapValue(mapType uint32, keyData *byte, keySize int) uint32
	ProxyGetHeaderMapPairs(mapType uint32, returnValueData **byte, returnValueSize *int) uint32
	ProxySetHeaderMapPairs(mapType uint32, mapData *byte, mapSize int) uint32
	ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) uint32
	ProxySetBufferBytes(bt types.BufferType, start int, maxSize int, bufferData *byte, bufferSize int) uint32
	ProxyHttpCall(upstreamData *byte, upstreamSize int, headerData *byte, headerSize int, bodyData *byte, bodySize int, trailersData *byte, trailersSize int, timeout uint32, calloutIDPtr *uint32) uint32
	ProxyCallForeignFunction(funcNamePtr *byte, funcNameSize int, paramPtr *byte, paramSize int, returnData **byte, returnSize *int) uint32
	ProxySetTickPeriodMilliseconds(period uint32) uint32
	ProxySetEffectiveContext(contextID uint32) uint32
	ProxyDone() uint32
	ProxyDefineMetric(metricType types.MetricType, metricNameData *byte, metricNameSize int, returnMetricIDPtr *uint32) uint32
	ProxyIncrementMetric(metricID uint32, offset int64) uint32
	ProxyRecordMetric(metricID uint32, value uint64) uint32
	ProxyGetMetric(metricID uint32, returnMetricValue *uint64) uint32
}

type DefaultProxyWAMSHost struct{}

var _ ProxyWasmHost = DefaultProxyWAMSHost{}

func (d DefaultProxyWAMSHost) ProxyLog(logLevel uint32, messageData *byte, messageSize int) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxySendLocalResponse(statusCode uint32, statusCodeDetailData *byte, statusCodeDetailsSize int, bodyData *byte, bodySize int, headersData *byte, headersSize int, grpcStatus int32) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyGetSharedData(keyData *byte, keySize int, returnValueData **byte, returnValueSize *int, returnCas *uint32) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxySetSharedData(keyData *byte, keySize int, valueData *byte, valueSize int, cas uint32) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyRegisterSharedQueue(nameData *byte, nameSize int, returnID *uint32) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyResolveSharedQueue(vmIDData *byte, vmIDSize int, nameData *byte, nameSize int, returnID *uint32) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyDequeueSharedQueue(queueID uint32, returnValueData **byte, returnValueSize *int) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyEnqueueSharedQueue(queueID uint32, valueData *byte, valueSize int) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyGetHeaderMapValue(mapType uint32, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyAddHeaderMapValue(mapType uint32, keyData *byte, keySize int, valueData *byte, valueSize int) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyReplaceHeaderMapValue(mapType uint32, keyData *byte, keySize int, valueData *byte, valueSize int) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyContinueStream(streamType uint32) uint32 { return 0 }
func (d DefaultProxyWAMSHost) ProxyCloseStream(streamType uint32) uint32    { return 0 }
func (d DefaultProxyWAMSHost) ProxyRemoveHeaderMapValue(mapType uint32, keyData *byte, keySize int) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyGetHeaderMapPairs(mapType uint32, returnValueData **byte, returnValueSize *int) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxySetHeaderMapPairs(mapType uint32, mapData *byte, mapSize int) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxySetBufferBytes(bt types.BufferType, start int, maxSize int, bufferData *byte, bufferSize int) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyHttpCall(upstreamData *byte, upstreamSize int, headerData *byte, headerSize int, bodyData *byte, bodySize int, trailersData *byte, trailersSize int, timeout uint32, calloutIDPtr *uint32) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyCallForeignFunction(funcNamePtr *byte, funcNameSize int, paramPtr *byte, paramSize int, returnData **byte, returnSize *int) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxySetTickPeriodMilliseconds(period uint32) uint32 { return 0 }
func (d DefaultProxyWAMSHost) ProxySetEffectiveContext(contextID uint32) uint32    { return 0 }
func (d DefaultProxyWAMSHost) ProxyDone() uint32                                   { return 0 }
func (d DefaultProxyWAMSHost) ProxyDefineMetric(metricType types.MetricType, metricNameData *byte, metricNameSize int, returnMetricIDPtr *uint32) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyIncrementMetric(metricID uint32, offset int64) uint32 {
	return 0
}
func (d DefaultProxyWAMSHost) ProxyRecordMetric(metricID uint32, value uint64) uint32 { return 0 }
func (d DefaultProxyWAMSHost) ProxyGetMetric(metricID uint32, returnMetricValue *uint64) uint32 {
	return 0
}

func ProxyLog(logLevel uint32, messageData *byte, messageSize int) uint32 {
	return currentHost.ProxyLog(logLevel, messageData, messageSize)
}

func ProxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int) uint32 {
	return currentHost.ProxySetProperty(pathData, pathSize, valueData, valueSize)
}

func ProxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int) uint32 {
	return currentHost.ProxyGetProperty(pathData, pathSize, returnValueData, returnValueSize)
}

func ProxySendLocalResponse(statusCode uint32, statusCodeDetailData *byte,
	statusCodeDetailsSize int, bodyData *byte, bodySize int, headersData *byte, headersSize int, grpcStatus int32) uint32 {
	return currentHost.ProxySendLocalResponse(statusCode,
		statusCodeDetailData, statusCodeDetailsSize, bodyData, bodySize, headersData, headersSize, grpcStatus)
}

func ProxyGetSharedData(keyData *byte, keySize int, returnValueData **byte, returnValueSize *int, returnCas *uint32) uint32 {
	return currentHost.ProxyGetSharedData(keyData, keySize, returnValueData, returnValueSize, returnCas)
}

func ProxySetSharedData(keyData *byte, keySize int, valueData *byte, valueSize int, cas uint32) uint32 {
	return currentHost.ProxySetSharedData(keyData, keySize, valueData, valueSize, cas)
}

func ProxyRegisterSharedQueue(nameData *byte, nameSize int, returnID *uint32) uint32 {
	return currentHost.ProxyRegisterSharedQueue(nameData, nameSize, returnID)
}

func ProxyResolveSharedQueue(vmIDData *byte, vmIDSize int, nameData *byte, nameSize int, returnID *uint32) uint32 {
	return currentHost.ProxyResolveSharedQueue(vmIDData, vmIDSize, nameData, nameSize, returnID)
}

func ProxyDequeueSharedQueue(queueID uint32, returnValueData **byte, returnValueSize *int) uint32 {
	return currentHost.ProxyDequeueSharedQueue(queueID, returnValueData, returnValueSize)
}

func ProxyEnqueueSharedQueue(queueID uint32, valueData *byte, valueSize int) uint32 {
	return currentHost.ProxyEnqueueSharedQueue(queueID, valueData, valueSize)
}

func ProxyGetHeaderMapValue(mapType uint32, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) uint32 {
	return currentHost.ProxyGetHeaderMapValue(mapType, keyData, keySize, returnValueData, returnValueSize)
}

func ProxyAddHeaderMapValue(mapType uint32, keyData *byte, keySize int, valueData *byte, valueSize int) uint32 {
	return currentHost.ProxyAddHeaderMapValue(mapType, keyData, keySize, valueData, valueSize)
}

func ProxyReplaceHeaderMapValue(mapType uint32, keyData *byte, keySize int, valueData *byte, valueSize int) uint32 {
	return currentHost.ProxyReplaceHeaderMapValue(mapType, keyData, keySize, valueData, valueSize)
}

func ProxyContinueStream(streamType uint32) uint32 {
	return currentHost.ProxyContinueStream(streamType)
}

func ProxyCloseStream(streamType uint32) uint32 {
	return currentHost.ProxyCloseStream(streamType)
}
func ProxyRemoveHeaderMapValue(mapType uint32, keyData *byte, keySize int) uint32 {
	return currentHost.ProxyRemoveHeaderMapValue(mapType, keyData, keySize)
}

func ProxyGetHeaderMapPairs(mapType uint32, returnValueData **byte, returnValueSize *int) uint32 {
	return currentHost.ProxyGetHeaderMapPairs(mapType, returnValueData, returnValueSize)
}

func ProxySetHeaderMapPairs(mapType uint32, mapData *byte, mapSize int) uint32 {
	return currentHost.ProxySetHeaderMapPairs(mapType, mapData, mapSize)
}

func ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) uint32 {
	return currentHost.ProxyGetBufferBytes(bt, start, maxSize, returnBufferData, returnBufferSize)
}

func ProxySetBufferBytes(bt types.BufferType, start int, maxSize int, bufferData *byte, bufferSize int) uint32 {
	return currentHost.ProxySetBufferBytes(bt, start, maxSize, bufferData, bufferSize)
}

func ProxyHttpCall(upstreamData *byte, upstreamSize int, headerData *byte, headerSize int, bodyData *byte,
	bodySize int, trailersData *byte, trailersSize int, timeout uint32, calloutIDPtr *uint32) uint32 {
	return currentHost.ProxyHttpCall(upstreamData, upstreamSize,
		headerData, headerSize, bodyData, bodySize, trailersData, trailersSize, timeout, calloutIDPtr)
}

func ProxyCallForeignFunction(funcNamePtr *byte, funcNameSize int, paramPtr *byte, paramSize int, returnData **byte, returnSize *int) uint32 {
	return currentHost.ProxyCallForeignFunction(funcNamePtr, funcNameSize, paramPtr, paramSize, returnData, returnSize)
}

func ProxySetTickPeriodMilliseconds(period uint32) uint32 {
	return currentHost.ProxySetTickPeriodMilliseconds(period)
}

func ProxySetEffectiveContext(contextID uint32) uint32 {
	return currentHost.ProxySetEffectiveContext(contextID)
}

func ProxyDone() uint32 {
	return currentHost.ProxyDone()
}

func ProxyDefineMetric(metricType types.MetricType,
	metricNameData *byte, metricNameSize int, returnMetricIDPtr *uint32) uint32 {
	return currentHost.ProxyDefineMetric(metricType, metricNameData, metricNameSize, returnMetricIDPtr)
}

func ProxyIncrementMetric(metricID uint32, offset int64) uint32 {
	return currentHost.ProxyIncrementMetric(metricID, offset)
}

func ProxyRecordMetric(metricID uint32, value uint64) uint32 {
	return currentHost.ProxyRecordMetric(metricID, value)
}

func ProxyGetMetric(metricID uint32, returnMetricValue *uint64) uint32 {
	return currentHost.ProxyGetMetric(metricID, returnMetricValue)
}
