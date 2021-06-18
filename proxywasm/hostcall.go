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

package proxywasm

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
)

// GetVMConfiguration is used for retrieving configurations ginve in the "vm_config.configuration" field.
// This hostcall is only available during types.RootContext.OnVMStart call.
// "size" argument specifies homw many bytes you want to read. Set it to "vmConfigurationSize" given in OnVMStart.
func GetVMConfiguration(size int) ([]byte, error) {
	return getBuffer(internal.BufferTypeVMConfiguration, 0, size)
}

// GetPluginConfiguration is used for retrieving configurations ginve in the "config.configuration" field.
// This hostcall is only available during types.RootContext.OnPluginStart call.
// "size" argument specifies homw many bytes you want to read. Set it to "pluginConfigurationSize" given in OnVMStart.
func GetPluginConfiguration(size int) ([]byte, error) {
	return getBuffer(internal.BufferTypePluginConfiguration, 0, size)
}

// SetTickPeriodMilliSeconds sets the tick interval of types.RootContext.OnTick calls.
// Only available for types.RootContext.
func SetTickPeriodMilliSeconds(millSec uint32) error {
	return internal.StatusToError(internal.ProxySetTickPeriodMilliseconds(millSec))
}

// RegisterSharedQueue registers the shared queue on this root context.
// "Register" means that OnQueueReady is called for this root context whenever a new item is enqueued on that queueID.
// Only available for types.RootContext. The returned ququeID can be used for Enqueue/DequeueSharedQueue.
// Note that "name" must be unique across all Wasm VMs which share a same "vm_id".
// That means you can use "vm_id" can be used for separating shared queue namespace.
func RegisterSharedQueue(name string) (ququeID uint32, err error) {
	var queueID uint32
	ptr := internal.StringBytePtr(name)
	st := internal.ProxyRegisterSharedQueue(ptr, len(name), &queueID)
	return queueID, internal.StatusToError(st)
}

// ResolveSharedQueue acquires the queueID for the given vm_id and queue name.
// The returned ququeID can be used for Enqueue/DequeueSharedQueue.
func ResolveSharedQueue(vmID, queueName string) (ququeID uint32, err error) {
	var ret uint32
	st := internal.ProxyResolveSharedQueue(internal.StringBytePtr(vmID),
		len(vmID), internal.StringBytePtr(queueName), len(queueName), &ret)
	return ret, internal.StatusToError(st)
}

// EnqueueSharedQueue enqueues an data to the shared queue of the given queueID.
// In order to get queue id for a target queue, use "ResolveSharedQueue" first.
func EnqueueSharedQueue(queueID uint32, data []byte) error {
	return internal.StatusToError(internal.ProxyEnqueueSharedQueue(queueID, &data[0], len(data)))
}

// DequeueSharedQueue dequeues an data from the shared queue of the given queueID.
// In order to get queue id for a target queue, use "ResolveSharedQueue" first.
func DequeueSharedQueue(queueID uint32) ([]byte, error) {
	var raw *byte
	var size int
	st := internal.ProxyDequeueSharedQueue(queueID, &raw, &size)
	if st != internal.StatusOK {
		return nil, internal.StatusToError(st)
	}
	return internal.RawBytePtrToByteSlice(raw, size), nil
}

// Done must be callsed when OnPluginDone returnes false indicating that the plugin is in pending state
// right before deletion by hots. Only available for types.RootContext.
func Done() {
	internal.ProxyDone()
}

// DispatchHttpCall is for dipatching http calls to a remote cluster. This can be used by all contexts
// including Tcp and Root contexts. "cluster" arg specifies the remote cluster the host will send
// the request against with "headers", "body", "trailers" arguments. If the host successfully made the request
// and recevived the response from the remote cluster, then "callBack" function is called.
// During callBack is called, "GetHttpCallResponseHeaders", "GetHttpCallResponseBody", "GetHttpCallResponseTrailers"
// calls are available for accessing the response information.
func DispatchHttpCall(
	cluster string,
	headers [][2]string,
	body []byte,
	trailers [][2]string,
	timeoutMillisecond uint32,
	callBack func(numHeaders, bodySize, numTrailers int),
) (calloutID uint32, err error) {
	shs := internal.SerializeMap(headers)
	hp := &shs[0]
	hl := len(shs)

	sts := internal.SerializeMap(trailers)
	tp := &sts[0]
	tl := len(sts)

	var bodyPtr *byte
	if len(body) > 0 {
		bodyPtr = &body[0]
	}

	u := internal.StringBytePtr(cluster)
	switch st := internal.ProxyHttpCall(u, len(cluster),
		hp, hl, bodyPtr, len(body), tp, tl, timeoutMillisecond, &calloutID); st {
	case internal.StatusOK:
		internal.RegisterHttpCallout(calloutID, callBack)
		return calloutID, nil
	default:
		return 0, internal.StatusToError(st)
	}
}

// GetHttpCallResponseHeaders is used for retrieving http response headers
// returned by a remote cluster in reponse to the DispatchHttpCall.
// Only available during "callback" function passed to DispatchHttpCall.
func GetHttpCallResponseHeaders() ([][2]string, error) {
	return getMap(internal.MapTypeHttpCallResponseHeaders)
}

// GetHttpCallResponseBody is used for retrieving http response body
// returned by a remote cluster in reponse to the DispatchHttpCall.
// Only available during "callback" function passed to DispatchHttpCall.
func GetHttpCallResponseBody(start, maxSize int) ([]byte, error) {
	return getBuffer(internal.BufferTypeHttpCallResponseBody, start, maxSize)
}

// GetHttpCallResponseTrailers is used for retrieving http response trailers
// returned by a remote cluster in reponse to the DispatchHttpCall.
// Only available during "callback" function passed to DispatchHttpCall.
func GetHttpCallResponseTrailers() ([][2]string, error) {
	return getMap(internal.MapTypeHttpCallResponseTrailers)
}

// GetDownstreamData can be used for retrieving tcp downstream data in buffered in the host.
// Returned bytes begining from "start" to "start" +"maxSize" in the buffer.
// Only available during types.TcpContext.OnDownstreamData.
func GetDownstreamData(start, maxSize int) ([]byte, error) {
	return getBuffer(internal.BufferTypeDownstreamData, start, maxSize)
}

// AppendDownstreamData appends the given bytes to the downstream tcp data in buffered in the host.
// Only available during types.TcpContext.OnDownstreamData.
func AppendDownstreamData(data []byte) error {
	return appendToBuffer(internal.BufferTypeDownstreamData, data)
}

// PrependDownstreamData prepends the given bytes to the downstream tcp data in buffered in the host.
// Only available during types.TcpContext.OnDownstreamData.
func PrependDownstreamData(data []byte) error {
	return prependToBuffer(internal.BufferTypeDownstreamData, data)
}

// ReplaceDownstreamData replaces the downstream tcp data in buffered in the host
// with the given bytes. Only available during types.TcpContext.OnDownstreamData.
func ReplaceDownstreamData(data []byte) error {
	return replaceBuffer(internal.BufferTypeDownstreamData, data)
}

// GetDownstreamData can be used for retrieving upstream tcp data in buffered in the host.
// Returned bytes begining from "start" to "start" +"maxSize" in the buffer.
// Only available during types.TcpContext.OnUpstreamData.
func GetUpstreamData(start, maxSize int) ([]byte, error) {
	return getBuffer(internal.BufferTypeUpstreamData, start, maxSize)
}

// AppendUpstreamData appends the given bytes to the upstream tcp data in buffered in the host.
// Only available during types.TcpContext.OnUpstreamData.
func AppendUpstreamData(data []byte) error {
	return appendToBuffer(internal.BufferTypeUpstreamData, data)
}

// PrependUpstreamData prepends the given bytes to the upstream tcp data in buffered in the host.
// Only available during types.TcpContext.OnUpstreamData.
func PrependUpstreamData(data []byte) error {
	return prependToBuffer(internal.BufferTypeUpstreamData, data)
}

// ReplaceUpstreamData replaces the upstream tcp data in buffered in the host
// with the given bytes. Only available during types.TcpContext.OnUpstreamData.
func ReplaceUpstreamData(data []byte) error {
	return replaceBuffer(internal.BufferTypeUpstreamData, data)
}

// ContinueTcpStream continues interating on the tcp connection
// after types.Action.Pause was returned by types.TcpContext.
// Only available for types.TcpContext.
func ContinueTcpStream() error {
	// Note that internal.ProxyContinueStream is not implemented in Envoy,
	// so we intentionally choose to pass StreamTypeDownstream here while
	// the name itself is not indiciating "continue downstream".
	return internal.StatusToError(internal.ProxyContinueStream(internal.StreamTypeDownstream))
}

// CloseDownstream closes the downstream tcp connection for this Tcp context.
// Only available for types.TcpContext.
func CloseDownstream() error {
	return internal.StatusToError(internal.ProxyCloseStream(internal.StreamTypeDownstream))
}

// CloseUpstream closes the upstream tcp connection for this Tcp context.
// Only available for types.TcpContext.
func CloseUpstream() error {
	return internal.StatusToError(internal.ProxyCloseStream(internal.StreamTypeUpstream))
}

// GetHttpRequestHeaders is used for retrieving http request headers.
// Only available during types.HttpContext.OnHttpRequestHeaders and
// types.HttpContext.OnHttpStreamDone.
func GetHttpRequestHeaders() ([][2]string, error) {
	return getMap(internal.MapTypeHttpRequestHeaders)
}

func ReplaceHttpRequestHeaders(headers [][2]string) error {
	return setMap(internal.MapTypeHttpRequestHeaders, headers)
}

func GetHttpRequestHeader(key string) (string, error) {
	return getMapValue(internal.MapTypeHttpRequestHeaders, key)
}

func RemoveHttpRequestHeader(key string) error {
	return removeMapValue(internal.MapTypeHttpRequestHeaders, key)
}

func SetHttpRequestHeader(key, value string) error {
	return setMapValue(internal.MapTypeHttpRequestHeaders, key, value)
}

func AddHttpRequestHeader(key, value string) error {
	return addMapValue(internal.MapTypeHttpRequestHeaders, key, value)
}

func GetHttpRequestBody(start, maxSize int) ([]byte, error) {
	return getBuffer(internal.BufferTypeHttpRequestBody, start, maxSize)
}

func AppendHttpRequestBody(data []byte) error {
	return appendToBuffer(internal.BufferTypeHttpRequestBody, data)
}

func PrependHttpRequestBody(data []byte) error {
	return prependToBuffer(internal.BufferTypeHttpRequestBody, data)
}

func ReplaceHttpRequestBody(data []byte) error {
	return replaceBuffer(internal.BufferTypeHttpRequestBody, data)
}

func GetHttpRequestTrailers() ([][2]string, error) {
	return getMap(internal.MapTypeHttpRequestTrailers)
}

func ReplaceHttpRequestTrailers(trailers [][2]string) error {
	return setMap(internal.MapTypeHttpRequestTrailers, trailers)
}

func GetHttpRequestTrailer(key string) (string, error) {
	return getMapValue(internal.MapTypeHttpRequestTrailers, key)
}

func RemoveHttpRequestTrailer(key string) error {
	return removeMapValue(internal.MapTypeHttpRequestTrailers, key)
}

func SetHttpRequestTrailer(key, value string) error {
	return setMapValue(internal.MapTypeHttpRequestTrailers, key, value)
}

func AddHttpRequestTrailer(key, value string) error {
	return addMapValue(internal.MapTypeHttpRequestTrailers, key, value)
}

func ResumeHttpRequest() error {
	return internal.StatusToError(internal.ProxyContinueStream(internal.StreamTypeRequest))
}

// GetHttpResponseHeaders is used for retrieving http response headers.
// Only available during types.HttpContext.OnHttpResponseHeaders and
// types.HttpContext.OnHttpStreamDone.
func GetHttpResponseHeaders() ([][2]string, error) {
	return getMap(internal.MapTypeHttpResponseHeaders)
}

func ReplaceHttpResponseHeaders(headers [][2]string) error {
	return setMap(internal.MapTypeHttpResponseHeaders, headers)
}

func GetHttpResponseHeader(key string) (string, error) {
	return getMapValue(internal.MapTypeHttpResponseHeaders, key)
}

func RemoveHttpResponseHeader(key string) error {
	return removeMapValue(internal.MapTypeHttpResponseHeaders, key)
}

func SetHttpResponseHeader(key, value string) error {
	return setMapValue(internal.MapTypeHttpResponseHeaders, key, value)
}

func AddHttpResponseHeader(key, value string) error {
	return addMapValue(internal.MapTypeHttpResponseHeaders, key, value)
}

func GetHttpResponseBody(start, maxSize int) ([]byte, error) {
	return getBuffer(internal.BufferTypeHttpResponseBody, start, maxSize)
}

func AppendHttpResponseBody(data []byte) error {
	return appendToBuffer(internal.BufferTypeHttpResponseBody, data)
}

func PrependHttpResponseBody(data []byte) error {
	return prependToBuffer(internal.BufferTypeHttpResponseBody, data)
}

func ReplaceHttpResponseBody(data []byte) error {
	return replaceBuffer(internal.BufferTypeHttpResponseBody, data)
}

func GetHttpResponseTrailers() ([][2]string, error) {
	return getMap(internal.MapTypeHttpResponseTrailers)
}

func ReplaceHttpResponseTrailers(trailers [][2]string) error {
	return setMap(internal.MapTypeHttpResponseTrailers, trailers)
}

func GetHttpResponseTrailer(key string) (string, error) {
	return getMapValue(internal.MapTypeHttpResponseTrailers, key)
}

func RemoveHttpResponseTrailer(key string) error {
	return removeMapValue(internal.MapTypeHttpResponseTrailers, key)
}

func SetHttpResponseTrailer(key, value string) error {
	return setMapValue(internal.MapTypeHttpResponseTrailers, key, value)
}

func AddHttpResponseTrailer(key, value string) error {
	return addMapValue(internal.MapTypeHttpResponseTrailers, key, value)
}

func ResumeHttpResponse() error {
	return internal.StatusToError(internal.ProxyContinueStream(internal.StreamTypeResponse))
}

// SendHttpResponse sends a http response to the downstream with given information (headers, statuscode, body).
// This call cannot be used outside types.HttpContext otherwise an error should be returned.
// Also please note that this cannot be used after types.HttpContext.OnHttpResponseHeaders returns Continue
// since in that case, the response headers may have already arrived at the downstream and there is no way
// to override the already sent headers.
func SendHttpResponse(statusCode uint32, headers [][2]string, body []byte) error {
	shs := internal.SerializeMap(headers)
	var bp *byte
	if len(body) > 0 {
		bp = &body[0]
	}
	hp := &shs[0]
	hl := len(shs)
	return internal.StatusToError(
		internal.ProxySendLocalResponse(
			statusCode, nil, 0,
			bp, len(body), hp, hl, -1,
		),
	)
}

func GetSharedData(key string) (value []byte, cas uint32, err error) {
	var raw *byte
	var size int

	st := internal.ProxyGetSharedData(internal.StringBytePtr(key), len(key), &raw, &size, &cas)
	if st != internal.StatusOK {
		return nil, 0, internal.StatusToError(st)
	}
	return internal.RawBytePtrToByteSlice(raw, size), cas, nil
}

func SetSharedData(key string, data []byte, cas uint32) error {
	st := internal.ProxySetSharedData(internal.StringBytePtr(key),
		len(key), &data[0], len(data), cas)
	return internal.StatusToError(st)
}

func GetProperty(path []string) ([]byte, error) {
	var ret *byte
	var retSize int
	raw := internal.SerializePropertyPath(path)

	err := internal.StatusToError(internal.ProxyGetProperty(&raw[0], len(raw), &ret, &retSize))
	if err != nil {
		return nil, err
	}

	return internal.RawBytePtrToByteSlice(ret, retSize), nil

}

func SetProperty(path []string, data []byte) error {
	raw := internal.SerializePropertyPath(path)
	return internal.StatusToError(internal.ProxySetProperty(
		&raw[0], len(path), &data[0], len(data),
	))
}

// CallForeignFunction calls a foreign function of given funcName defined by host implementations.
// Foreign functions are host specific functions so please refer to the doc of your host implementation for detail.
func CallForeignFunction(funcName string, param []byte) (ret []byte, err error) {
	f := internal.StringBytePtr(funcName)

	var returnData *byte
	var returnSize int

	switch st := internal.ProxyCallForeignFunction(f, len(funcName), &param[0], len(param), &returnData, &returnSize); st {
	case internal.StatusOK:
		return internal.RawBytePtrToByteSlice(returnData, returnSize), nil
	default:
		return nil, internal.StatusToError(st)
	}
}
