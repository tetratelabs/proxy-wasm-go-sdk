package runtime

import (
	"strconv"
)

type Context interface {
	DispatchHttpCall(upstream string, headers [][2]string,
		body string, trailers [][2]string, timeoutMillisecond uint32) (calloutID uint32, status Status)
	OnHttpCallResponse(calloutID uint32, numHeaders, bodySize, numTrailers int)
	GetHttpCallResponseHeaders() ([][2]string, Status)
	GetHttpCallResponseBody(start, maxSize int) ([]byte, Status)
	GetHttpCallResponseTrailers() ([][2]string, Status)
	GetCurrentTime() int64
	OnDone() bool
	Done()
}

type RootContext interface {
	Context
	OnVMStart(vmConfigurationSize int) bool
	OnConfigure(pluginConfigurationSize int) bool
	GetConfiguration() ([]byte, Status)
	SetTickPeriod(period uint32) Status
	OnQueueReady(queueID uint32)
	OnTick()
	OnLog()
}

type StreamContext interface {
	Context
	OnNewConnection() Action
	OnDownstreamData(dataSize int, endOfStream bool) Action
	GetDownStreamData(start, maxSize int) ([]byte, Status)
	OnDownStreamClose(peerType PeerType)

	OnUpstreamData(dataSize int, endOfStream bool) Action
	GetUpstreamData(start, maxSize int) ([]byte, Status)
	OnUpstreamStreamClose(peerType PeerType)
	OnLog()
}

type HttpContext interface {
	Context

	// request
	OnHttpRequestHeaders(numHeaders int) Action
	GetHttpRequestHeaders() ([][2]string, Status)
	SetHttpRequestHeaders(headers [][2]string) Status
	GetHttpRequestHeader(key string) (string, Status)
	RemoveHttpRequestHeader(key string) Status
	SetHttpRequestHeader(key, value string) Status
	AddHttpRequestHeader(key, value string) Status

	OnHttpRequestBody(bodySize int, endOfStream bool) Action
	GetHttpRequestBody(start, maxSize int) ([]byte, Status)

	OnHttpRequestTrailers(numTrailers int) Action
	GetHttpRequestTrailers() ([][2]string, Status)
	SetHttpRequestTrailers(headers [][2]string) Status
	GetHttpRequestTrailer(key string) (string, Status)
	RemoveHttpRequestTrailer(key string) Status
	SetHttpRequestTrailer(key, value string) Status
	AddHttpRequestTrailer(key, value string) Status

	ResumeHttpRequest() Status

	// response
	OnHttpResponseHeaders(numHeaders int) Action
	GetHttpResponseHeaders() ([][2]string, Status)
	SetHttpResponseHeaders(headers [][2]string) Status
	GetHttpResponseHeader(key string) (string, Status)
	RemoveHttpResponseHeader(key string) Status
	SetHttpResponseHeader(key, value string) Status
	AddHttpResponseHeader(key, value string) Status

	OnHttpResponseBody(bodySize int, endOfStream bool) Action
	GetHttpResponseBody(start, maxSize int) ([]byte, Status)

	OnHttpResponseTrailers(numTrailers int) Action
	GetHttpResponseTrailers() ([][2]string, Status)
	SetHttpResponseTrailers(headers [][2]string) Status
	GetHttpResponseTrailer(key string) (string, Status)
	RemoveHttpResponseTrailer(key string) Status
	SetHttpResponseTrailer(key, value string) Status
	AddHttpResponseTrailer(key, value string) Status

	ResumeHttpResponse() Status

	SendHttpResponse(statusCode uint32, headers [][2]string, body string) Status
	ClearHttpRouteCache() Status
	OnLog()
}

type DefaultContext struct{}

var (
	_ Context       = &DefaultContext{}
	_ RootContext   = &DefaultContext{}
	_ StreamContext = &DefaultContext{}
	_ HttpContext   = &DefaultContext{}
)

// impl Context
func (d *DefaultContext) GetCurrentTime() int64 {
	var t int64
	proxyGetCurrentTimeNanoseconds(&t)
	return t
}

// impl Context
func (d *DefaultContext) DispatchHttpCall(upstream string,
	headers [][2]string, body string, trailers [][2]string, timeoutMillisecond uint32) (uint32, Status) {
	return dispatchHttpCall(upstream, headers, body, trailers, timeoutMillisecond)
}

// impl Context
func (d *DefaultContext) GetHttpCallResponseHeaders() ([][2]string, Status) {
	return getMap(MapTypeHttpCallResponseHeaders)
}

// impl Context
func (d *DefaultContext) GetHttpCallResponseBody(start, maxSize int) ([]byte, Status) {
	return getBuffer(BufferTypeHttpCallResponseBody, start, maxSize)
}

// impl Context
func (d *DefaultContext) GetHttpCallResponseTrailers() ([][2]string, Status) {
	return getMap(MapTypeHttpCallResponseTrailers)
}

// impl Context
func (d *DefaultContext) OnHttpCallResponse(calloutID uint32, numHeaders, bodySize, numTrailers int) {
}

// impl Context
func (d *DefaultContext) OnDone() bool {
	return true
}

// impl Context
func (d *DefaultContext) Done() {
	switch st := proxyDone(); st {
	case StatusOk:
		return
	default:
		panic("unexpected status: " + strconv.FormatUint(uint64(st), 10))
	}
}

// impl HttpContext, StreamContext, RootContext
func (d *DefaultContext) OnLog() {}

// impl RootContext
func (d *DefaultContext) OnVMStart(_ int) bool {
	return true
}

// impl RootContext
func (d *DefaultContext) OnConfigure(_ int) bool {
	return true
}

// impl RootContext
func (d *DefaultContext) GetConfiguration() ([]byte, Status) {
	return getConfiguration()
}

// impl RootContext
func (d *DefaultContext) SetTickPeriod(milliSec uint32) Status {
	return setTickPeriodMilliSeconds(milliSec)
}

// impl RootContext
func (d *DefaultContext) OnTick() {}

// impl RootContext
func (d *DefaultContext) OnQueueReady(_ uint32) {}

// impl StreamContext
func (d *DefaultContext) OnNewConnection() Action {
	return ActionContinue
}

// impl StreamContext
func (d *DefaultContext) OnDownstreamData(dataSize int, endOfStream bool) Action {
	return ActionContinue
}

// impl StreamContext
func (d *DefaultContext) GetDownStreamData(start, maxSize int) ([]byte, Status) {
	return getBuffer(BufferTypeDownstreamData, start, maxSize)
}

// impl StreamContext
func (d *DefaultContext) OnDownStreamClose(_ PeerType) {}

// impl StreamContext
func (d *DefaultContext) OnUpstreamData(_ int, _ bool) Action {
	return ActionContinue
}

// impl StreamContext
func (d *DefaultContext) GetUpstreamData(start, maxSize int) ([]byte, Status) {
	return getBuffer(BufferTypeUpstreamData, start, maxSize)
}

// impl StreamContext
func (d *DefaultContext) OnUpstreamStreamClose(_ PeerType) {}

// impl HttpContext
func (d *DefaultContext) OnHttpRequestHeaders(_ int) Action {
	return ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpRequestHeaders() ([][2]string, Status) {
	return getMap(MapTypeHttpRequestHeaders)
}

// impl HttpContext
func (d *DefaultContext) SetHttpRequestHeaders(headers [][2]string) Status {
	return setMap(MapTypeHttpRequestHeaders, headers)
}

// impl HttpContext
func (d *DefaultContext) GetHttpRequestHeader(key string) (string, Status) {
	return getMapValue(MapTypeHttpRequestHeaders, key)
}

// impl HttpContext
func (d *DefaultContext) RemoveHttpRequestHeader(key string) Status {
	return removeMapValue(MapTypeHttpRequestHeaders, key)
}

// impl HttpContext
func (d *DefaultContext) SetHttpRequestHeader(key, value string) Status {
	return setMapValue(MapTypeHttpRequestHeaders, key, value)
}

// impl HttpContext
func (d *DefaultContext) AddHttpRequestHeader(key, value string) Status {
	return addMapValue(MapTypeHttpRequestHeaders, key, value)
}

// impl HttpContext
func (d *DefaultContext) OnHttpRequestBody(_ int, _ bool) Action {
	return ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpRequestBody(start, maxSize int) ([]byte, Status) {
	return getBuffer(BufferTypeHttpRequestBody, start, maxSize)
}

// impl HttpContext
func (d *DefaultContext) OnHttpRequestTrailers(numTrailers int) Action {
	return ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpRequestTrailers() ([][2]string, Status) {
	return getMap(MapTypeHttpRequestTrailers)
}

// impl HttpContext
func (d *DefaultContext) SetHttpRequestTrailers(headers [][2]string) Status {
	return setMap(MapTypeHttpRequestTrailers, headers)
}

// impl HttpContext
func (d *DefaultContext) GetHttpRequestTrailer(key string) (string, Status) {
	return getMapValue(MapTypeHttpRequestTrailers, key)
}

// impl HttpContext
func (d *DefaultContext) RemoveHttpRequestTrailer(key string) Status {
	return removeMapValue(MapTypeHttpRequestTrailers, key)
}

// impl HttpContext
func (d *DefaultContext) SetHttpRequestTrailer(key, value string) Status {
	return setMapValue(MapTypeHttpRequestTrailers, key, value)
}

// impl HttpContext
func (d *DefaultContext) AddHttpRequestTrailer(key, value string) Status {
	return addMapValue(MapTypeHttpRequestTrailers, key, value)
}

// impl HttpContext
func (d *DefaultContext) ResumeHttpRequest() Status {
	return proxyContinueRequest()
}

// impl HttpContext
func (d *DefaultContext) OnHttpResponseHeaders(_ int) Action {
	return ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpResponseHeaders() ([][2]string, Status) {
	return getMap(MapTypeHttpResponseHeaders)
}

// impl HttpContext
func (d *DefaultContext) SetHttpResponseHeaders(headers [][2]string) Status {
	return setMap(MapTypeHttpResponseHeaders, headers)
}

// impl HttpContext
func (d *DefaultContext) GetHttpResponseHeader(key string) (string, Status) {
	return getMapValue(MapTypeHttpResponseHeaders, key)
}

// impl HttpContext
func (d *DefaultContext) RemoveHttpResponseHeader(key string) Status {
	return removeMapValue(MapTypeHttpResponseHeaders, key)
}

// impl HttpContext
func (d *DefaultContext) SetHttpResponseHeader(key, value string) Status {
	return setMapValue(MapTypeHttpResponseHeaders, key, value)
}

// impl HttpContext
func (d *DefaultContext) AddHttpResponseHeader(key, value string) Status {
	return addMapValue(MapTypeHttpResponseHeaders, key, value)
}

// impl HttpContext
func (d *DefaultContext) OnHttpResponseBody(size int, endOfStream bool) Action {
	return ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpResponseBody(start, maxSize int) ([]byte, Status) {
	return getBuffer(BufferTypeHttpResponseBody, start, maxSize)
}

// impl HttpContext
func (d *DefaultContext) OnHttpResponseTrailers(numTrailers int) Action {
	return ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpResponseTrailers() ([][2]string, Status) {
	return getMap(MapTypeHttpResponseTrailers)

}

// impl HttpContext
func (d *DefaultContext) SetHttpResponseTrailers(headers [][2]string) Status {
	return setMap(MapTypeHttpResponseTrailers, headers)
}

// impl HttpContext
func (d *DefaultContext) GetHttpResponseTrailer(key string) (string, Status) {
	return getMapValue(MapTypeHttpResponseTrailers, key)
}

// impl HttpContext
func (d *DefaultContext) RemoveHttpResponseTrailer(key string) Status {
	return removeMapValue(MapTypeHttpResponseTrailers, key)
}

// impl HttpContext
func (d *DefaultContext) SetHttpResponseTrailer(key, value string) Status {
	return setMapValue(MapTypeHttpResponseTrailers, key, value)
}

// impl HttpContext
func (d *DefaultContext) AddHttpResponseTrailer(key, value string) Status {
	return addMapValue(MapTypeHttpResponseTrailers, key, value)
}

// impl HttpContext
func (d *DefaultContext) ResumeHttpResponse() Status {
	return proxyContinueResponse()
}

// impl HttpContext
func (d *DefaultContext) SendHttpResponse(statusCode uint32, headers [][2]string, body string) Status {
	return sendHttpResponse(statusCode, headers, body)
}

// impl HttpContext
func (d *DefaultContext) ClearHttpRouteCache() Status {
	return proxyClearRouteCache()
}
