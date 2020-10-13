// Copyright 2020 Tetrate
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

package proxytest

import (
	"log"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type (
	httpHostEmulator struct {
		httpStreams map[uint32]*httpStreamState
	}
	httpStreamState struct {
		requestHeaders, responseHeaders,
		requestTrailers, responseTrailers [][2]string
		requestBody, responseBody []byte

		action            types.Action
		sentLocalResponse *LocalHttpResponse
	}
	LocalHttpResponse struct {
		StatusCode       uint32
		StatusCodeDetail string
		Data             []byte
		Headers          [][2]string
		GRPCStatus       int32
	}
)

func newHttpHostEmulator() *httpHostEmulator {
	host := &httpHostEmulator{httpStreams: map[uint32]*httpStreamState{}}
	return host
}

// impl rawhostcall.ProxyWASMHost: delegated from hostEmulator
func (h *httpHostEmulator) httpHostEmulatorProxyGetBufferBytes(bt types.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) types.Status {
	active := proxywasm.VMStateGetActiveContextID()
	stream := h.httpStreams[active]
	var buf []byte
	switch bt {
	case types.BufferTypeHttpRequestBody:
		buf = stream.requestBody
	case types.BufferTypeHttpResponseBody:
		buf = stream.requestBody
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	if start >= len(buf) {
		log.Printf("start index out of range: %d (start) >= %d ", start, len(buf))
		return types.StatusBadArgument
	}

	*returnBufferData = &buf[start]
	if maxSize > len(buf)-start {
		*returnBufferSize = len(buf) - start
	} else {
		*returnBufferSize = maxSize
	}
	return types.StatusOK
}

func (h *httpHostEmulator) httpHostEmulatorProxySetBufferBytes(bt types.BufferType, start int, maxSize int,
	bufferData *byte, bufferSize int) types.Status {
	body := proxywasm.RawBytePtrToByteSlice(bufferData, bufferSize)
	active := proxywasm.VMStateGetActiveContextID()
	stream := h.httpStreams[active]
	switch bt {
	case types.BufferTypeHttpRequestBody:
		stream.requestBody = body
	case types.BufferTypeHttpResponseBody:
		stream.responseBody = body
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost: delegated from hostEmulator
func (h *httpHostEmulator) httpHostEmulatorProxyGetHeaderMapValue(mapType types.MapType, keyData *byte,
	keySize int, returnValueData **byte, returnValueSize *int) types.Status {
	key := proxywasm.RawBytePtrToString(keyData, keySize)
	active := proxywasm.VMStateGetActiveContextID()
	stream := h.httpStreams[active]

	var headers [][2]string
	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		headers = stream.requestHeaders
	case types.MapTypeHttpResponseHeaders:
		headers = stream.responseHeaders
	case types.MapTypeHttpRequestTrailers:
		headers = stream.requestTrailers
	case types.MapTypeHttpResponseTrailers:
		headers = stream.responseTrailers
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	for _, h := range headers {
		if h[0] == key {
			value := []byte(h[1])
			*returnValueData = &value[0]
			*returnValueSize = len(value)
			return types.StatusOK
		}
	}

	return types.StatusNotFound
}

// impl rawhostcall.ProxyWASMHost
func (h *httpHostEmulator) ProxyAddHeaderMapValue(mapType types.MapType, keyData *byte,
	keySize int, valueData *byte, valueSize int) types.Status {

	key := proxywasm.RawBytePtrToString(keyData, keySize)
	value := proxywasm.RawBytePtrToString(valueData, valueSize)
	active := proxywasm.VMStateGetActiveContextID()
	stream := h.httpStreams[active]

	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		stream.requestHeaders = addMapValue(stream.requestHeaders, key, value)
	case types.MapTypeHttpResponseHeaders:
		stream.responseHeaders = addMapValue(stream.responseHeaders, key, value)
	case types.MapTypeHttpRequestTrailers:
		stream.requestTrailers = addMapValue(stream.requestTrailers, key, value)
	case types.MapTypeHttpResponseTrailers:
		stream.responseTrailers = addMapValue(stream.responseTrailers, key, value)
	default:
		panic("unimplemented")
	}

	return types.StatusOK
}

func addMapValue(base [][2]string, key, value string) [][2]string {
	for i, h := range base {
		if h[0] == key {
			h[1] += value
			base[i] = h
			return base
		}
	}
	return append(base, [2]string{key, value})
}

// impl rawhostcall.ProxyWASMHost
func (h *httpHostEmulator) ProxyReplaceHeaderMapValue(mapType types.MapType, keyData *byte,
	keySize int, valueData *byte, valueSize int) types.Status {
	key := proxywasm.RawBytePtrToString(keyData, keySize)
	value := proxywasm.RawBytePtrToString(valueData, valueSize)
	active := proxywasm.VMStateGetActiveContextID()
	stream := h.httpStreams[active]

	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		stream.requestHeaders = replaceMapValue(stream.requestHeaders, key, value)
	case types.MapTypeHttpResponseHeaders:
		stream.responseHeaders = replaceMapValue(stream.responseHeaders, key, value)
	case types.MapTypeHttpRequestTrailers:
		stream.requestTrailers = replaceMapValue(stream.requestTrailers, key, value)
	case types.MapTypeHttpResponseTrailers:
		stream.responseTrailers = replaceMapValue(stream.responseTrailers, key, value)
	default:
		panic("unimplemented")
	}
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func replaceMapValue(base [][2]string, key, value string) [][2]string {
	for i, h := range base {
		if h[0] == key {
			h[1] = value
			base[i] = h
			return base
		}
	}
	return append(base, [2]string{key, value})
}

// impl rawhostcall.ProxyWASMHost
func (h *httpHostEmulator) ProxyRemoveHeaderMapValue(mapType types.MapType, keyData *byte, keySize int) types.Status {
	key := proxywasm.RawBytePtrToString(keyData, keySize)
	active := proxywasm.VMStateGetActiveContextID()
	stream := h.httpStreams[active]

	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		stream.requestHeaders = removeHeaderMapValue(stream.requestHeaders, key)
	case types.MapTypeHttpResponseHeaders:
		stream.responseHeaders = removeHeaderMapValue(stream.responseHeaders, key)
	case types.MapTypeHttpRequestTrailers:
		stream.requestTrailers = removeHeaderMapValue(stream.requestTrailers, key)
	case types.MapTypeHttpResponseTrailers:
		stream.responseTrailers = removeHeaderMapValue(stream.responseTrailers, key)
	default:
		panic("unimplemented")
	}
	return types.StatusOK
}

func removeHeaderMapValue(base [][2]string, key string) [][2]string {
	for i, h := range base {
		if h[0] == key {
			if len(base)-1 == i {
				return base[:i]
			} else {
				return append(base[:i], base[i+1:]...)
			}
		}
	}
	return base
}

// impl rawhostcall.ProxyWASMHost: delegated from hostEmulator
func (h *httpHostEmulator) httpHostEmulatorProxyGetHeaderMapPairs(mapType types.MapType, returnValueData **byte,
	returnValueSize *int) types.Status {
	active := proxywasm.VMStateGetActiveContextID()
	stream := h.httpStreams[active]

	var m []byte
	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		m = proxywasm.SerializeMap(stream.requestHeaders)
	case types.MapTypeHttpResponseHeaders:
		m = proxywasm.SerializeMap(stream.responseHeaders)
	case types.MapTypeHttpRequestTrailers:
		m = proxywasm.SerializeMap(stream.requestTrailers)
	case types.MapTypeHttpResponseTrailers:
		m = proxywasm.SerializeMap(stream.responseTrailers)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	*returnValueData = &m[0]
	*returnValueSize = len(m)
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func (h *httpHostEmulator) ProxySetHeaderMapPairs(mapType types.MapType, mapData *byte, mapSize int) types.Status {
	m := proxywasm.DeserializeMap(proxywasm.RawBytePtrToByteSlice(mapData, mapSize))
	active := proxywasm.VMStateGetActiveContextID()
	stream := h.httpStreams[active]

	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		stream.requestHeaders = m
	case types.MapTypeHttpResponseHeaders:
		stream.responseHeaders = m
	case types.MapTypeHttpRequestTrailers:
		stream.requestTrailers = m
	case types.MapTypeHttpResponseTrailers:
		stream.responseTrailers = m
	default:
		panic("unimplemented")
	}
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func (h *httpHostEmulator) ProxyContinueStream(types.StreamType) types.Status {
	active := proxywasm.VMStateGetActiveContextID()
	stream := h.httpStreams[active]
	stream.action = types.ActionContinue
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func (h *httpHostEmulator) ProxySendLocalResponse(statusCode uint32,
	statusCodeDetailData *byte, statusCodeDetailsSize int, bodyData *byte, bodySize int,
	headersData *byte, headersSize int, grpcStatus int32) types.Status {
	active := proxywasm.VMStateGetActiveContextID()
	stream := h.httpStreams[active]
	stream.sentLocalResponse = &LocalHttpResponse{
		StatusCode:       statusCode,
		StatusCodeDetail: proxywasm.RawBytePtrToString(statusCodeDetailData, statusCodeDetailsSize),
		Data:             proxywasm.RawBytePtrToByteSlice(bodyData, bodySize),
		Headers:          proxywasm.DeserializeMap(proxywasm.RawBytePtrToByteSlice(headersData, headersSize)),
		GRPCStatus:       grpcStatus,
	}
	return types.StatusOK
}

// impl HostEmulator
func (h *httpHostEmulator) HttpFilterInitContext() (contextID uint32) {
	contextID = getNextContextID()
	proxywasm.ProxyOnContextCreate(contextID, rootContextID)
	h.httpStreams[contextID] = &httpStreamState{action: types.ActionContinue}
	return
}

// impl HostEmulator
func (h *httpHostEmulator) HttpFilterPutRequestHeaders(contextID uint32, headers [][2]string) {
	h.HttpFilterPutRequestHeadersEndOfStream(contextID, headers, false)
}

// impl HostEmulator
func (h *httpHostEmulator) HttpFilterPutRequestHeadersEndOfStream(contextID uint32, headers [][2]string, endOfStream bool) {
	cs, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.requestHeaders = headers
	cs.action = proxywasm.ProxyOnRequestHeaders(contextID,
		len(headers), endOfStream)
}

// impl HostEmulator
func (h *httpHostEmulator) HttpFilterPutResponseHeaders(contextID uint32, headers [][2]string) {
	h.HttpFilterPutResponseHeadersEndOfStream(contextID, headers, false)
}

// impl HostEmulator
func (h *httpHostEmulator) HttpFilterPutResponseHeadersEndOfStream(contextID uint32, headers [][2]string, endOfStream bool) {
	cs, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.responseHeaders = headers

	cs.action = proxywasm.ProxyOnResponseHeaders(contextID,
		len(headers), endOfStream)
}

// impl HostEmulator
func (h *httpHostEmulator) HttpFilterPutRequestTrailers(contextID uint32, headers [][2]string) {
	cs, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.requestTrailers = headers
	cs.action = proxywasm.ProxyOnRequestTrailers(contextID, len(headers))
}

// impl HostEmulator
func (h *httpHostEmulator) HttpFilterPutResponseTrailers(contextID uint32, headers [][2]string) {
	cs, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.responseTrailers = headers
	cs.action = proxywasm.ProxyOnResponseTrailers(contextID, len(headers))
}

// impl HostEmulator
func (h *httpHostEmulator) HttpFilterPutRequestBody(contextID uint32, body []byte) {
	h.HttpFilterPutRequestBodyEndOfStream(contextID, body, false)
}

// impl HostEmulator
func (h *httpHostEmulator) HttpFilterPutRequestBodyEndOfStream(contextID uint32, body []byte, endOfStream bool) {
	cs, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.requestBody = body
	cs.action = proxywasm.ProxyOnRequestBody(contextID,
		len(body), endOfStream)
}

// impl HostEmulator
func (h *httpHostEmulator) HttpFilterGetRequestBody(contextID uint32) []byte {
	cs, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	return cs.requestBody
}

// impl HostEmulator
func (h *httpHostEmulator) HttpFilterPutResponseBody(contextID uint32, body []byte) {
	h.HttpFilterPutResponseBodyEndOfStream(contextID, body, false)
}

func (h *httpHostEmulator) HttpFilterPutResponseBodyEndOfStream(contextID uint32, body []byte, endOfStream bool) {
	cs, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.responseBody = body
	cs.action = proxywasm.ProxyOnResponseBody(contextID,
		len(body), endOfStream)
}

// impl HostEmulator
func (h *httpHostEmulator) HttpFilterGetResponseBody(contextID uint32) []byte {
	cs, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	return cs.responseBody
}

// impl HostEmulator
func (h *httpHostEmulator) HttpFilterCompleteHttpStream(contextID uint32) {
	proxywasm.ProxyOnDone(contextID)
}

// impl HostEmulator
func (h *httpHostEmulator) HttpFilterGetCurrentStreamAction(contextID uint32) types.Action {
	stream, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}
	return stream.action
}

// impl HostEmulator
func (h *httpHostEmulator) HttpFilterGetSentLocalResponse(contextID uint32) *LocalHttpResponse {
	return h.httpStreams[contextID].sentLocalResponse
}
