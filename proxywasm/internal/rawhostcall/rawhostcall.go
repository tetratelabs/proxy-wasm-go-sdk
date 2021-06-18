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

// +build !proxytest

package rawhostcall

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

//export proxy_log
func ProxyLog(logLevel uint32, messageData *byte, messageSize int) uint32

//export proxy_send_local_response
func ProxySendLocalResponse(statusCode uint32, statusCodeDetailData *byte, statusCodeDetailsSize int,
	bodyData *byte, bodySize int, headersData *byte, headersSize int, grpcStatus int32) uint32

//export proxy_get_shared_data
func ProxyGetSharedData(keyData *byte, keySize int, returnValueData **byte, returnValueSize *int, returnCas *uint32) uint32

//export proxy_set_shared_data
func ProxySetSharedData(keyData *byte, keySize int, valueData *byte, valueSize int, cas uint32) uint32

//export proxy_register_shared_queue
func ProxyRegisterSharedQueue(nameData *byte, nameSize int, returnID *uint32) uint32

//export proxy_resolve_shared_queue
func ProxyResolveSharedQueue(vmIDData *byte, vmIDSize int, nameData *byte, nameSize int, returnID *uint32) uint32

//export proxy_dequeue_shared_queue
func ProxyDequeueSharedQueue(queueID uint32, returnValueData **byte, returnValueSize *int) uint32

//export proxy_enqueue_shared_queue
func ProxyEnqueueSharedQueue(queueID uint32, valueData *byte, valueSize int) uint32

//export proxy_get_header_map_value
func ProxyGetHeaderMapValue(mapType uint32, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) uint32

//export proxy_add_header_map_value
func ProxyAddHeaderMapValue(mapType uint32, keyData *byte, keySize int, valueData *byte, valueSize int) uint32

//export proxy_replace_header_map_value
func ProxyReplaceHeaderMapValue(mapType uint32, keyData *byte, keySize int, valueData *byte, valueSize int) uint32

//export proxy_remove_header_map_value
func ProxyRemoveHeaderMapValue(mapType uint32, keyData *byte, keySize int) uint32

//export proxy_get_header_map_pairs
func ProxyGetHeaderMapPairs(mapType uint32, returnValueData **byte, returnValueSize *int) uint32

//export proxy_set_header_map_pairs
func ProxySetHeaderMapPairs(mapType uint32, mapData *byte, mapSize int) uint32

//export proxy_get_buffer_bytes
func ProxyGetBufferBytes(bufferType uint32, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) uint32

//export proxy_set_buffer_bytes
func ProxySetBufferBytes(bufferType uint32, start int, maxSize int, bufferData *byte, bufferSize int) uint32

//export proxy_continue_stream
func ProxyContinueStream(streamType uint32) uint32

//export proxy_close_stream
func ProxyCloseStream(streamType uint32) uint32

//export proxy_http_call
func ProxyHttpCall(upstreamData *byte, upstreamSize int, headerData *byte, headerSize int,
	bodyData *byte, bodySize int, trailersData *byte, trailersSize int, timeout uint32, calloutIDPtr *uint32,
) uint32

//export proxy_call_foreign_function
func ProxyCallForeignFunction(funcNamePtr *byte, funcNameSize int, paramPtr *byte, paramSize int, returnData **byte, returnSize *int) uint32

//export proxy_set_tick_period_milliseconds
func ProxySetTickPeriodMilliseconds(period uint32) uint32

//export proxy_set_effective_context
func ProxySetEffectiveContext(contextID uint32) uint32

//export proxy_done
func ProxyDone() uint32

//export proxy_define_metric
func ProxyDefineMetric(metricType types.MetricType, metricNameData *byte, metricNameSize int, returnMetricIDPtr *uint32) uint32

//export proxy_increment_metric
func ProxyIncrementMetric(metricID uint32, offset int64) uint32

//export proxy_record_metric
func ProxyRecordMetric(metricID uint32, value uint64) uint32

//export proxy_get_metric
func ProxyGetMetric(metricID uint32, returnMetricValue *uint64) uint32

//export proxy_get_property
func ProxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int) uint32

//export proxy_set_property
func ProxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int) uint32
