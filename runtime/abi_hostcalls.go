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

//export proxy_log
func proxyLog(logLevel LogLevel, messageData *byte, messageSize int) Status

//export proxy_set_property
func proxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int)

//export proxy_get_property
func proxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int)

//export proxy_send_local_response
func proxySendLocalResponse(statusCode uint32, statusCodeDetailData *byte, statusCodeDetailsSize int,
	bodyData *byte, bodySize int, headersData *byte, headersSize int, grpcStatus int32) Status

//export proxy_get_shared_data
func proxyGetSharedData(keyData *byte, keySize int, returnValueData **byte, returnValueSize *byte, returnCas *uint32) Status

//export proxy_set_shared_data
func proxySetSharedData(keyData *byte, keySize int, valueData *byte, valueSize int, cas uint32) Status

//export proxy_register_shared_queue
func proxyRegisterSharedQueue(nameData *byte, nameSize uint, returnID *uint32) Status

//export proxy_resolve_shared_queue
func proxyResolveSharedQueue(vmIDData *byte, vmIDSize uint, nameData *byte, nameSize uint, returnID *uint32) Status

//export proxy_dequeue_shared_queue
func proxyDequeueSharedQueue(queueID uint32, returnValueData **byte, returnValueSize *byte) Status

//export proxy_enqueue_shared_queue
func proxyEnqueueSharedQueue(queueID uint32, valueData *byte, valueSize uint) Status

//export proxy_get_header_map_value
func proxyGetHeaderMapValue(mapType MapType, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) Status

//export proxy_add_header_map_value
func proxyAddHeaderMapValue(mapType MapType, keyData *byte, keySize int, valueData *byte, valueSize int) Status

//export proxy_replace_header_map_value
func proxyReplaceHeaderMapValue(mapType MapType, keyData *byte, keySize int, valueData *byte, valueSize int) Status

//export proxy_continue_stream
func proxyContinueStream(streamType StreamType) Status

//export proxy_close_stream
func proxyCloseStream(streamType StreamType) Status

//export proxy_remove_header_map_value
func proxyRemoveHeaderMapValue(mapType MapType, keyData *byte, keySize int) Status

//export proxy_get_header_map_pairs
func proxyGetHeaderMapPairs(mapType MapType, returnValueData **byte, returnValueSize *int) Status

//export proxy_set_header_map_pairs
func proxySetHeaderMapPairs(mapType MapType, mapData *byte, mapSize int) Status

//export proxy_get_buffer_bytes
func proxyGetBufferBytes(bt BufferType, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) Status

//export proxy_http_call
func proxyHttpCall(upstreamData *byte, upstreamSize int, headerData *byte, headerSize int,
	bodyData *byte, bodySize int, trailersData *byte, trailersSize int, timeout uint32, calloutIDPtr *uint32,
) Status

//export proxy_set_tick_period_milliseconds
func proxySetTickPeriodMilliseconds(period uint32) Status

//export proxy_get_current_time_nanoseconds
func proxyGetCurrentTimeNanoseconds(returnTime *int64) Status

//export proxy_set_effective_context
func proxySetEffectiveContext(contextID uint32) Status

//export proxy_done
func proxyDone() Status
