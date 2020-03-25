package runtime

//go:export proxy_log
func proxyLog(logLevel LogLevel, messageData *byte, messageSize int) Status

//go:export proxy_get_configuration
func proxyGetConfiguration(returnBufferData **byte, returnBufferSize *int) Status

//go:export proxy_set_property
func proxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int)

//go:export proxy_get_property
func proxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int)

//go:export proxy_continue_request
func proxyContinueRequest() Status

//go:export proxy_continue_response
func proxyContinueResponse() Status

//go:export proxy_send_local_response
func proxySendLocalResponse(statusCode uint32, statusCodeDetailData *byte, statusCodeDetailsSize int,
	bodyData *byte, bodySize int, headersData *byte, headersSize int, grpcStatus int32) Status

//go:export proxy_clear_route_cache
func proxyClearRouteCache() Status

//go:export proxy_get_shared_data
func proxyGetSharedData(keyData *byte, keySize int, returnValueData **byte, returnValueSize *byte, returnCas *uint32) Status

//go:export proxy_set_shared_data
func proxySetSharedData(keyData *byte, keySize int, valueData *byte, valueSize int, cas uint32) Status

//go:export proxy_register_shared_queue
func proxyRegisterSharedQueue(nameData *byte, nameSize uint, returnID *uint32) Status

//go:export proxy_resolve_shared_queue
func proxyResolveSharedQueue(vmIDData *byte, vmIDSize uint, nameData *byte, nameSize uint, returnID *uint32) Status

//go:export proxy_dequeue_shared_queue
func proxyDequeueSharedQueue(queueID uint32, returnValueData **byte, returnValueSize *byte) Status

//go:export proxy_enqueue_shared_queue
func proxyEnqueueSharedQueue(queueID uint32, valueData *byte, valueSize uint) Status

//go:export proxy_get_header_map_value
func proxyGetHeaderMapValue(mapType MapType, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) Status

//go:export proxy_add_header_map_value
func proxyAddHeaderMapValue(mapType MapType, keyData *byte, keySize int, valueData *byte, valueSize int) Status

//go:export proxy_replace_header_map_value
func proxyReplaceHeaderMapValue(mapType MapType, keyData *byte, keySize int, valueData *byte, valueSize int) Status

//go:export proxy_remove_header_map_value
func proxyRemoveHeaderMapValue(mapType MapType, keyData *byte, keySize int) Status

//go:export proxy_get_header_map_pairs
func proxyGetHeaderMapPairs(mapType MapType, returnValueData **byte, returnValueSize *int) Status

//go:export proxy_set_header_map_pairs
func proxySetHeaderMapPairs(mapType MapType, mapData *byte, mapSize int) Status

//go:export proxy_get_buffer_bytes
func proxyGetBufferBytes(bt BufferType, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) Status

//go:export proxy_http_call
func proxyHttpCall(upstreamData *byte, upstreamSize int, headerData *byte, headerSize int,
	bodyData *byte, bodySize int, trailersData *byte, trailersSize int, timeout uint32, calloutIDPtr *uint32,
) Status

//go:export proxy_set_tick_period_milliseconds
func proxySetTickPeriodMilliseconds(period uint32) Status

//go:export proxy_get_current_time_nanoseconds
func proxyGetCurrentTimeNanoseconds(returnTime *int64) Status

//go:export proxy_set_effective_context
func proxySetEffectiveContext(contextID uint32) Status

//go:export proxy_done
func proxyDone() Status
