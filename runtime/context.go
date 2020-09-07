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
	"strconv"

	"github.com/mathetake/proxy-wasm-go/runtime/hostcall"
	"github.com/mathetake/proxy-wasm-go/runtime/types"
)

type Context interface {
	DispatchHttpCall(upstream string, headers [][2]string,
		body string, trailers [][2]string, timeoutMillisecond uint32) (calloutID uint32, status types.Status)
	OnHttpCallResponse(calloutID uint32, numHeaders, bodySize, numTrailers int)
	GetHttpCallResponseHeaders() ([][2]string, types.Status)
	GetHttpCallResponseBody(start, maxSize int) ([]byte, types.Status)
	GetHttpCallResponseTrailers() ([][2]string, types.Status)
	GetCurrentTime() int64
	OnDone() bool
	Done()
}

type RootContext interface {
	Context
	OnVMStart(vmConfigurationSize int) bool
	OnConfigure(pluginConfigurationSize int) bool
	GetPluginConfiguration(dataSize int) ([]byte, types.Status)
	SetTickPeriod(period uint32) types.Status
	OnQueueReady(queueID uint32)
	OnTick()
	OnLog()
}

type StreamContext interface {
	Context
	OnNewConnection() types.Action
	OnDownstreamData(dataSize int, endOfStream bool) types.Action
	GetDownStreamData(start, maxSize int) ([]byte, types.Status)
	OnDownStreamClose(peerType types.PeerType)

	OnUpstreamData(dataSize int, endOfStream bool) types.Action
	GetUpstreamData(start, maxSize int) ([]byte, types.Status)
	OnUpstreamStreamClose(peerType types.PeerType)
	OnLog()
}

type HttpContext interface {
	Context

	// request
	OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action
	GetHttpRequestHeaders() ([][2]string, types.Status)
	SetHttpRequestHeaders(headers [][2]string) types.Status
	GetHttpRequestHeader(key string) (string, types.Status)
	RemoveHttpRequestHeader(key string) types.Status
	SetHttpRequestHeader(key, value string) types.Status
	AddHttpRequestHeader(key, value string) types.Status

	OnHttpRequestBody(bodySize int, endOfStream bool) types.Action
	GetHttpRequestBody(start, maxSize int) ([]byte, types.Status)

	OnHttpRequestTrailers(numTrailers int) types.Action
	GetHttpRequestTrailers() ([][2]string, types.Status)
	SetHttpRequestTrailers(headers [][2]string) types.Status
	GetHttpRequestTrailer(key string) (string, types.Status)
	RemoveHttpRequestTrailer(key string) types.Status
	SetHttpRequestTrailer(key, value string) types.Status
	AddHttpRequestTrailer(key, value string) types.Status

	ResumeHttpRequest() types.Status

	// response
	OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action
	GetHttpResponseHeaders() ([][2]string, types.Status)
	SetHttpResponseHeaders(headers [][2]string) types.Status
	GetHttpResponseHeader(key string) (string, types.Status)
	RemoveHttpResponseHeader(key string) types.Status
	SetHttpResponseHeader(key, value string) types.Status
	AddHttpResponseHeader(key, value string) types.Status

	OnHttpResponseBody(bodySize int, endOfStream bool) types.Action
	GetHttpResponseBody(start, maxSize int) ([]byte, types.Status)

	OnHttpResponseTrailers(numTrailers int) types.Action
	GetHttpResponseTrailers() ([][2]string, types.Status)
	SetHttpResponseTrailers(headers [][2]string) types.Status
	GetHttpResponseTrailer(key string) (string, types.Status)
	RemoveHttpResponseTrailer(key string) types.Status
	SetHttpResponseTrailer(key, value string) types.Status
	AddHttpResponseTrailer(key, value string) types.Status

	ResumeHttpResponse() types.Status

	SendHttpResponse(statusCode uint32, headers [][2]string, body string) types.Status
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
	hostcall.ProxyGetCurrentTimeNanoseconds(&t)
	return t
}

// impl Context
func (d *DefaultContext) DispatchHttpCall(upstream string,
	headers [][2]string, body string, trailers [][2]string, timeoutMillisecond uint32) (uint32, types.Status) {
	return dispatchHttpCall(upstream, headers, body, trailers, timeoutMillisecond)
}

// impl Context
func (d *DefaultContext) GetHttpCallResponseHeaders() ([][2]string, types.Status) {
	return getMap(types.MapTypeHttpCallResponseHeaders)
}

// impl Context
func (d *DefaultContext) GetHttpCallResponseBody(start, maxSize int) ([]byte, types.Status) {
	return getBuffer(types.BufferTypeHttpCallResponseBody, start, maxSize)
}

// impl Context
func (d *DefaultContext) GetHttpCallResponseTrailers() ([][2]string, types.Status) {
	return getMap(types.MapTypeHttpCallResponseTrailers)
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
	switch st := hostcall.ProxyDone(); st {
	case types.StatusOk:
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
func (d *DefaultContext) GetPluginConfiguration(dataSize int) ([]byte, types.Status) {
	return getBuffer(types.BufferTypePluginConfiguration, 0, dataSize)
}

// impl RootContext
func (d *DefaultContext) SetTickPeriod(milliSec uint32) types.Status {
	return setTickPeriodMilliSeconds(milliSec)
}

// impl RootContext
func (d *DefaultContext) OnTick() {}

// impl RootContext
func (d *DefaultContext) OnQueueReady(_ uint32) {}

// impl StreamContext
func (d *DefaultContext) OnNewConnection() types.Action {
	return types.ActionContinue
}

// impl StreamContext
func (d *DefaultContext) OnDownstreamData(dataSize int, endOfStream bool) types.Action {
	return types.ActionContinue
}

// impl StreamContext
func (d *DefaultContext) GetDownStreamData(start, maxSize int) ([]byte, types.Status) {
	return getBuffer(types.BufferTypeDownstreamData, start, maxSize)
}

// impl StreamContext
func (d *DefaultContext) OnDownStreamClose(_ types.PeerType) {}

// impl StreamContext
func (d *DefaultContext) OnUpstreamData(_ int, _ bool) types.Action {
	return types.ActionContinue
}

// impl StreamContext
func (d *DefaultContext) GetUpstreamData(start, maxSize int) ([]byte, types.Status) {
	return getBuffer(types.BufferTypeUpstreamData, start, maxSize)
}

// impl StreamContext
func (d *DefaultContext) OnUpstreamStreamClose(_ types.PeerType) {}

// impl HttpContext
func (d *DefaultContext) OnHttpRequestHeaders(_ int, _ bool) types.Action {
	return types.ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpRequestHeaders() ([][2]string, types.Status) {
	return getMap(types.MapTypeHttpRequestHeaders)
}

// impl HttpContext
func (d *DefaultContext) SetHttpRequestHeaders(headers [][2]string) types.Status {
	return setMap(types.MapTypeHttpRequestHeaders, headers)
}

// impl HttpContext
func (d *DefaultContext) GetHttpRequestHeader(key string) (string, types.Status) {
	return getMapValue(types.MapTypeHttpRequestHeaders, key)
}

// impl HttpContext
func (d *DefaultContext) RemoveHttpRequestHeader(key string) types.Status {
	return removeMapValue(types.MapTypeHttpRequestHeaders, key)
}

// impl HttpContext
func (d *DefaultContext) SetHttpRequestHeader(key, value string) types.Status {
	return setMapValue(types.MapTypeHttpRequestHeaders, key, value)
}

// impl HttpContext
func (d *DefaultContext) AddHttpRequestHeader(key, value string) types.Status {
	return addMapValue(types.MapTypeHttpRequestHeaders, key, value)
}

// impl HttpContext
func (d *DefaultContext) OnHttpRequestBody(_ int, _ bool) types.Action {
	return types.ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpRequestBody(start, maxSize int) ([]byte, types.Status) {
	return getBuffer(types.BufferTypeHttpRequestBody, start, maxSize)
}

// impl HttpContext
func (d *DefaultContext) OnHttpRequestTrailers(numTrailers int) types.Action {
	return types.ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpRequestTrailers() ([][2]string, types.Status) {
	return getMap(types.MapTypeHttpRequestTrailers)
}

// impl HttpContext
func (d *DefaultContext) SetHttpRequestTrailers(headers [][2]string) types.Status {
	return setMap(types.MapTypeHttpRequestTrailers, headers)
}

// impl HttpContext
func (d *DefaultContext) GetHttpRequestTrailer(key string) (string, types.Status) {
	return getMapValue(types.MapTypeHttpRequestTrailers, key)
}

// impl HttpContext
func (d *DefaultContext) RemoveHttpRequestTrailer(key string) types.Status {
	return removeMapValue(types.MapTypeHttpRequestTrailers, key)
}

// impl HttpContext
func (d *DefaultContext) SetHttpRequestTrailer(key, value string) types.Status {
	return setMapValue(types.MapTypeHttpRequestTrailers, key, value)
}

// impl HttpContext
func (d *DefaultContext) AddHttpRequestTrailer(key, value string) types.Status {
	return addMapValue(types.MapTypeHttpRequestTrailers, key, value)
}

// impl HttpContext
func (d *DefaultContext) ResumeHttpRequest() types.Status {
	return hostcall.ProxyContinueStream(types.StreamTypeRequest)
}

// impl HttpContext
func (d *DefaultContext) OnHttpResponseHeaders(_ int, _ bool) types.Action {
	return types.ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpResponseHeaders() ([][2]string, types.Status) {
	return getMap(types.MapTypeHttpResponseHeaders)
}

// impl HttpContext
func (d *DefaultContext) SetHttpResponseHeaders(headers [][2]string) types.Status {
	return setMap(types.MapTypeHttpResponseHeaders, headers)
}

// impl HttpContext
func (d *DefaultContext) GetHttpResponseHeader(key string) (string, types.Status) {
	return getMapValue(types.MapTypeHttpResponseHeaders, key)
}

// impl HttpContext
func (d *DefaultContext) RemoveHttpResponseHeader(key string) types.Status {
	return removeMapValue(types.MapTypeHttpResponseHeaders, key)
}

// impl HttpContext
func (d *DefaultContext) SetHttpResponseHeader(key, value string) types.Status {
	return setMapValue(types.MapTypeHttpResponseHeaders, key, value)
}

// impl HttpContext
func (d *DefaultContext) AddHttpResponseHeader(key, value string) types.Status {
	return addMapValue(types.MapTypeHttpResponseHeaders, key, value)
}

// impl HttpContext
func (d *DefaultContext) OnHttpResponseBody(size int, endOfStream bool) types.Action {
	return types.ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpResponseBody(start, maxSize int) ([]byte, types.Status) {
	return getBuffer(types.BufferTypeHttpResponseBody, start, maxSize)
}

// impl HttpContext
func (d *DefaultContext) OnHttpResponseTrailers(numTrailers int) types.Action {
	return types.ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpResponseTrailers() ([][2]string, types.Status) {
	return getMap(types.MapTypeHttpResponseTrailers)

}

// impl HttpContext
func (d *DefaultContext) SetHttpResponseTrailers(headers [][2]string) types.Status {
	return setMap(types.MapTypeHttpResponseTrailers, headers)
}

// impl HttpContext
func (d *DefaultContext) GetHttpResponseTrailer(key string) (string, types.Status) {
	return getMapValue(types.MapTypeHttpResponseTrailers, key)
}

// impl HttpContext
func (d *DefaultContext) RemoveHttpResponseTrailer(key string) types.Status {
	return removeMapValue(types.MapTypeHttpResponseTrailers, key)
}

// impl HttpContext
func (d *DefaultContext) SetHttpResponseTrailer(key, value string) types.Status {
	return setMapValue(types.MapTypeHttpResponseTrailers, key, value)
}

// impl HttpContext
func (d *DefaultContext) AddHttpResponseTrailer(key, value string) types.Status {
	return addMapValue(types.MapTypeHttpResponseTrailers, key, value)
}

// impl HttpContext
func (d *DefaultContext) ResumeHttpResponse() types.Status {
	return hostcall.ProxyContinueStream(types.StreamTypeResponse)
}

// impl HttpContext
func (d *DefaultContext) SendHttpResponse(statusCode uint32, headers [][2]string, body string) types.Status {
	return sendHttpResponse(statusCode, headers, body)
}
