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
		body string, trailers [][2]string, timeoutMillisecond uint32) (calloutID uint32, status error)
	OnHttpCallResponse(calloutID uint32, numHeaders, bodySize, numTrailers int)
	GetHttpCallResponseHeaders() ([][2]string, error)
	GetHttpCallResponseBody(start, maxSize int) ([]byte, error)
	GetHttpCallResponseTrailers() ([][2]string, error)
	GetCurrentTime() int64
	OnDone() bool
	Done()
}

type RootContext interface {
	Context
	OnVMStart(vmConfigurationSize int) bool
	OnConfigure(pluginConfigurationSize int) bool
	GetPluginConfiguration(dataSize int) ([]byte, error)
	SetTickPeriod(period uint32) error
	OnQueueReady(queueID uint32)
	OnTick()
	OnLog()
}

type StreamContext interface {
	Context
	OnNewConnection() types.Action
	OnDownstreamData(dataSize int, endOfStream bool) types.Action
	GetDownStreamData(start, maxSize int) ([]byte, error)
	OnDownStreamClose(peerType types.PeerType)

	OnUpstreamData(dataSize int, endOfStream bool) types.Action
	GetUpstreamData(start, maxSize int) ([]byte, error)
	OnUpstreamStreamClose(peerType types.PeerType)
	OnLog()
}

type HttpContext interface {
	Context

	// request
	OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action
	GetHttpRequestHeaders() ([][2]string, error)
	SetHttpRequestHeaders(headers [][2]string) error
	GetHttpRequestHeader(key string) (string, error)
	RemoveHttpRequestHeader(key string) error
	SetHttpRequestHeader(key, value string) error
	AddHttpRequestHeader(key, value string) error

	OnHttpRequestBody(bodySize int, endOfStream bool) types.Action
	GetHttpRequestBody(start, maxSize int) ([]byte, error)

	OnHttpRequestTrailers(numTrailers int) types.Action
	GetHttpRequestTrailers() ([][2]string, error)
	SetHttpRequestTrailers(headers [][2]string) error
	GetHttpRequestTrailer(key string) (string, error)
	RemoveHttpRequestTrailer(key string) error
	SetHttpRequestTrailer(key, value string) error
	AddHttpRequestTrailer(key, value string) error

	ResumeHttpRequest() error

	// response
	OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action
	GetHttpResponseHeaders() ([][2]string, error)
	SetHttpResponseHeaders(headers [][2]string) error
	GetHttpResponseHeader(key string) (string, error)
	RemoveHttpResponseHeader(key string) error
	SetHttpResponseHeader(key, value string) error
	AddHttpResponseHeader(key, value string) error

	OnHttpResponseBody(bodySize int, endOfStream bool) types.Action
	GetHttpResponseBody(start, maxSize int) ([]byte, error)

	OnHttpResponseTrailers(numTrailers int) types.Action
	GetHttpResponseTrailers() ([][2]string, error)
	SetHttpResponseTrailers(headers [][2]string) error
	GetHttpResponseTrailer(key string) (string, error)
	RemoveHttpResponseTrailer(key string) error
	SetHttpResponseTrailer(key, value string) error
	AddHttpResponseTrailer(key, value string) error

	ResumeHttpResponse() error

	SendHttpResponse(statusCode uint32, headers [][2]string, body string) error
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
	headers [][2]string, body string, trailers [][2]string, timeoutMillisecond uint32) (uint32, error) {
	ret, st := dispatchHttpCall(upstream, headers, body, trailers, timeoutMillisecond)
	return ret, types.StatusToError(st)
}

// impl Context
func (d *DefaultContext) GetHttpCallResponseHeaders() ([][2]string, error) {
	ret, st := getMap(types.MapTypeHttpCallResponseHeaders)
	return ret, types.StatusToError(st)
}

// impl Context
func (d *DefaultContext) GetHttpCallResponseBody(start, maxSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypeHttpCallResponseBody, start, maxSize)
	return ret, types.StatusToError(st)
}

// impl Context
func (d *DefaultContext) GetHttpCallResponseTrailers() ([][2]string, error) {
	ret, st := getMap(types.MapTypeHttpCallResponseTrailers)
	return ret, types.StatusToError(st)
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
		panic("unexpected status on proxy_done: " + strconv.FormatUint(uint64(st), 10))
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
func (d *DefaultContext) GetPluginConfiguration(dataSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypePluginConfiguration, 0, dataSize)
	return ret, types.StatusToError(st)
}

// impl RootContext
func (d *DefaultContext) SetTickPeriod(milliSec uint32) error {
	return types.StatusToError(setTickPeriodMilliSeconds(milliSec))
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
func (d *DefaultContext) GetDownStreamData(start, maxSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypeDownstreamData, start, maxSize)
	return ret, types.StatusToError(st)
}

// impl StreamContext
func (d *DefaultContext) OnDownStreamClose(_ types.PeerType) {}

// impl StreamContext
func (d *DefaultContext) OnUpstreamData(_ int, _ bool) types.Action {
	return types.ActionContinue
}

// impl StreamContext
func (d *DefaultContext) GetUpstreamData(start, maxSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypeUpstreamData, start, maxSize)
	return ret, types.StatusToError(st)
}

// impl StreamContext
func (d *DefaultContext) OnUpstreamStreamClose(_ types.PeerType) {}

// impl HttpContext
func (d *DefaultContext) OnHttpRequestHeaders(_ int, _ bool) types.Action {
	return types.ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpRequestHeaders() ([][2]string, error) {
	ret, st := getMap(types.MapTypeHttpRequestHeaders)
	return ret, types.StatusToError(st)
}

// impl HttpContext
func (d *DefaultContext) SetHttpRequestHeaders(headers [][2]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpRequestHeaders, headers))
}

// impl HttpContext
func (d *DefaultContext) GetHttpRequestHeader(key string) (string, error) {
	ret, st := getMapValue(types.MapTypeHttpRequestHeaders, key)
	return ret, types.StatusToError(st)
}

// impl HttpContext
func (d *DefaultContext) RemoveHttpRequestHeader(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpRequestHeaders, key))
}

// impl HttpContext
func (d *DefaultContext) SetHttpRequestHeader(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpRequestHeaders, key, value))
}

// impl HttpContext
func (d *DefaultContext) AddHttpRequestHeader(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpRequestHeaders, key, value))
}

// impl HttpContext
func (d *DefaultContext) OnHttpRequestBody(_ int, _ bool) types.Action {
	return types.ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpRequestBody(start, maxSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypeHttpRequestBody, start, maxSize)
	return ret, types.StatusToError(st)
}

// impl HttpContext
func (d *DefaultContext) OnHttpRequestTrailers(numTrailers int) types.Action {
	return types.ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpRequestTrailers() ([][2]string, error) {
	ret, st := getMap(types.MapTypeHttpRequestTrailers)
	return ret, types.StatusToError(st)
}

// impl HttpContext
func (d *DefaultContext) SetHttpRequestTrailers(headers [][2]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpRequestTrailers, headers))
}

// impl HttpContext
func (d *DefaultContext) GetHttpRequestTrailer(key string) (string, error) {
	ret, st := getMapValue(types.MapTypeHttpRequestTrailers, key)
	return ret, types.StatusToError(st)
}

// impl HttpContext
func (d *DefaultContext) RemoveHttpRequestTrailer(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpRequestTrailers, key))
}

// impl HttpContext
func (d *DefaultContext) SetHttpRequestTrailer(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpRequestTrailers, key, value))
}

// impl HttpContext
func (d *DefaultContext) AddHttpRequestTrailer(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpRequestTrailers, key, value))
}

// impl HttpContext
func (d *DefaultContext) ResumeHttpRequest() error {
	return types.StatusToError(hostcall.ProxyContinueStream(types.StreamTypeRequest))
}

// impl HttpContext
func (d *DefaultContext) OnHttpResponseHeaders(_ int, _ bool) types.Action {
	return types.ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpResponseHeaders() ([][2]string, error) {
	ret, st := getMap(types.MapTypeHttpResponseHeaders)
	return ret, types.StatusToError(st)
}

// impl HttpContext
func (d *DefaultContext) SetHttpResponseHeaders(headers [][2]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpResponseHeaders, headers))
}

// impl HttpContext
func (d *DefaultContext) GetHttpResponseHeader(key string) (string, error) {
	ret, st := getMapValue(types.MapTypeHttpResponseHeaders, key)
	return ret, types.StatusToError(st)
}

// impl HttpContext
func (d *DefaultContext) RemoveHttpResponseHeader(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpResponseHeaders, key))
}

// impl HttpContext
func (d *DefaultContext) SetHttpResponseHeader(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpResponseHeaders, key, value))
}

// impl HttpContext
func (d *DefaultContext) AddHttpResponseHeader(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpResponseHeaders, key, value))
}

// impl HttpContext
func (d *DefaultContext) OnHttpResponseBody(size int, endOfStream bool) types.Action {
	return types.ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpResponseBody(start, maxSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypeHttpResponseBody, start, maxSize)
	return ret, types.StatusToError(st)
}

// impl HttpContext
func (d *DefaultContext) OnHttpResponseTrailers(numTrailers int) types.Action {
	return types.ActionContinue
}

// impl HttpContext
func (d *DefaultContext) GetHttpResponseTrailers() ([][2]string, error) {
	ret, st := getMap(types.MapTypeHttpResponseTrailers)
	return ret, types.StatusToError(st)
}

// impl HttpContext
func (d *DefaultContext) SetHttpResponseTrailers(headers [][2]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpResponseTrailers, headers))
}

// impl HttpContext
func (d *DefaultContext) GetHttpResponseTrailer(key string) (string, error) {
	ret, st := getMapValue(types.MapTypeHttpResponseTrailers, key)
	return ret, types.StatusToError(st)
}

// impl HttpContext
func (d *DefaultContext) RemoveHttpResponseTrailer(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpResponseTrailers, key))
}

// impl HttpContext
func (d *DefaultContext) SetHttpResponseTrailer(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpResponseTrailers, key, value))
}

// impl HttpContext
func (d *DefaultContext) AddHttpResponseTrailer(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpResponseTrailers, key, value))
}

// impl HttpContext
func (d *DefaultContext) ResumeHttpResponse() error {
	return types.StatusToError(hostcall.ProxyContinueStream(types.StreamTypeResponse))
}

// impl HttpContext
func (d *DefaultContext) SendHttpResponse(statusCode uint32, headers [][2]string, body string) error {
	return types.StatusToError(sendHttpResponse(statusCode, headers, body))
}
