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

package proxytest

import (
	"log"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type (
	httpHostEmulator struct {
		httpStreams map[uint32]*httpStreamState
	}
	httpStreamState struct {
		requestHeaders, responseHeaders   [][2]string
		requestTrailers, responseTrailers [][2]string
		requestBody, responseBody         []byte

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

// impl internal.ProxyWasmHost: delegated from hostEmulator
func (h *httpHostEmulator) httpHostEmulatorProxyGetBufferBytes(bt internal.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) internal.Status {
	active := internal.VMStateGetActiveContextID()
	stream := h.httpStreams[active]
	var buf []byte
	switch bt {
	case internal.BufferTypeHttpRequestBody:
		buf = stream.requestBody
	case internal.BufferTypeHttpResponseBody:
		buf = stream.requestBody
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	if len(buf) == 0 {
		return internal.StatusNotFound
	} else if start >= len(buf) {
		log.Printf("start index out of range: %d (start) >= %d ", start, len(buf))
		return internal.StatusBadArgument
	}

	*returnBufferData = &buf[start]
	if maxSize > len(buf)-start {
		*returnBufferSize = len(buf) - start
	} else {
		*returnBufferSize = maxSize
	}
	return internal.StatusOK
}

func (h *httpHostEmulator) httpHostEmulatorProxySetBufferBytes(bt internal.BufferType, start int, maxSize int,
	bufferData *byte, bufferSize int) internal.Status {
	active := internal.VMStateGetActiveContextID()
	stream := h.httpStreams[active]
	var targetBuf *[]byte
	switch bt {
	case internal.BufferTypeHttpRequestBody:
		targetBuf = &stream.requestBody
	case internal.BufferTypeHttpResponseBody:
		targetBuf = &stream.responseBody
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	body := internal.RawBytePtrToByteSlice(bufferData, bufferSize)
	if start == 0 {
		if maxSize == 0 {
			// Prepend
			*targetBuf = append(body, *targetBuf...)
			return internal.StatusOK
		} else if maxSize >= len(*targetBuf) {
			// Replace
			*targetBuf = body
			return internal.StatusOK
		} else {
			return internal.StatusBadArgument
		}
	} else if start >= len(*targetBuf) {
		// Append.
		*targetBuf = append(*targetBuf, body...)
		return internal.StatusOK
	} else {
		return internal.StatusBadArgument
	}
}

// impl internal.ProxyWasmHost: delegated from hostEmulator
func (h *httpHostEmulator) httpHostEmulatorProxyGetHeaderMapValue(mapType internal.MapType, keyData *byte,
	keySize int, returnValueData **byte, returnValueSize *int) internal.Status {
	key := internal.RawBytePtrToString(keyData, keySize)
	active := internal.VMStateGetActiveContextID()
	stream := h.httpStreams[active]

	var headers [][2]string
	switch mapType {
	case internal.MapTypeHttpRequestHeaders:
		headers = stream.requestHeaders
	case internal.MapTypeHttpResponseHeaders:
		headers = stream.responseHeaders
	case internal.MapTypeHttpRequestTrailers:
		headers = stream.requestTrailers
	case internal.MapTypeHttpResponseTrailers:
		headers = stream.responseTrailers
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	for _, h := range headers {
		if h[0] == key {
			value := []byte(h[1])
			*returnValueData = &value[0]
			*returnValueSize = len(value)
			return internal.StatusOK
		}
	}

	return internal.StatusNotFound
}

// impl internal.ProxyWasmHost
func (h *httpHostEmulator) ProxyAddHeaderMapValue(mapType internal.MapType, keyData *byte,
	keySize int, valueData *byte, valueSize int) internal.Status {

	key := internal.RawBytePtrToString(keyData, keySize)
	value := internal.RawBytePtrToString(valueData, valueSize)
	active := internal.VMStateGetActiveContextID()
	stream := h.httpStreams[active]

	switch mapType {
	case internal.MapTypeHttpRequestHeaders:
		stream.requestHeaders = addMapValue(stream.requestHeaders, key, value)
	case internal.MapTypeHttpResponseHeaders:
		stream.responseHeaders = addMapValue(stream.responseHeaders, key, value)
	case internal.MapTypeHttpRequestTrailers:
		stream.requestTrailers = addMapValue(stream.requestTrailers, key, value)
	case internal.MapTypeHttpResponseTrailers:
		stream.responseTrailers = addMapValue(stream.responseTrailers, key, value)
	default:
		panic("unimplemented")
	}

	return internal.StatusOK
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

// impl internal.ProxyWasmHost
func (h *httpHostEmulator) ProxyReplaceHeaderMapValue(mapType internal.MapType, keyData *byte,
	keySize int, valueData *byte, valueSize int) internal.Status {
	key := internal.RawBytePtrToString(keyData, keySize)
	value := internal.RawBytePtrToString(valueData, valueSize)
	active := internal.VMStateGetActiveContextID()
	stream := h.httpStreams[active]

	switch mapType {
	case internal.MapTypeHttpRequestHeaders:
		stream.requestHeaders = replaceMapValue(stream.requestHeaders, key, value)
	case internal.MapTypeHttpResponseHeaders:
		stream.responseHeaders = replaceMapValue(stream.responseHeaders, key, value)
	case internal.MapTypeHttpRequestTrailers:
		stream.requestTrailers = replaceMapValue(stream.requestTrailers, key, value)
	case internal.MapTypeHttpResponseTrailers:
		stream.responseTrailers = replaceMapValue(stream.responseTrailers, key, value)
	default:
		panic("unimplemented")
	}
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
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

// impl internal.ProxyWasmHost
func (h *httpHostEmulator) ProxyRemoveHeaderMapValue(mapType internal.MapType, keyData *byte, keySize int) internal.Status {
	key := internal.RawBytePtrToString(keyData, keySize)
	active := internal.VMStateGetActiveContextID()
	stream := h.httpStreams[active]

	switch mapType {
	case internal.MapTypeHttpRequestHeaders:
		stream.requestHeaders = removeHeaderMapValue(stream.requestHeaders, key)
	case internal.MapTypeHttpResponseHeaders:
		stream.responseHeaders = removeHeaderMapValue(stream.responseHeaders, key)
	case internal.MapTypeHttpRequestTrailers:
		stream.requestTrailers = removeHeaderMapValue(stream.requestTrailers, key)
	case internal.MapTypeHttpResponseTrailers:
		stream.responseTrailers = removeHeaderMapValue(stream.responseTrailers, key)
	default:
		panic("unimplemented")
	}
	return internal.StatusOK
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

// impl internal.ProxyWasmHost: delegated from hostEmulator
func (h *httpHostEmulator) httpHostEmulatorProxyGetHeaderMapPairs(mapType internal.MapType, returnValueData **byte,
	returnValueSize *int) internal.Status {
	active := internal.VMStateGetActiveContextID()
	stream := h.httpStreams[active]

	var m []byte
	switch mapType {
	case internal.MapTypeHttpRequestHeaders:
		m = internal.SerializeMap(stream.requestHeaders)
	case internal.MapTypeHttpResponseHeaders:
		m = internal.SerializeMap(stream.responseHeaders)
	case internal.MapTypeHttpRequestTrailers:
		m = internal.SerializeMap(stream.requestTrailers)
	case internal.MapTypeHttpResponseTrailers:
		m = internal.SerializeMap(stream.responseTrailers)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	if len(m) == 0 {
		// The host might reutrn OK without setting the data pointer,
		// if there's nothing to pass to Wasm VM.
		*returnValueData = nil
		*returnValueSize = 0
		return internal.StatusOK
	}

	*returnValueData = &m[0]
	*returnValueSize = len(m)
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (h *httpHostEmulator) ProxySetHeaderMapPairs(mapType internal.MapType, mapData *byte, mapSize int) internal.Status {
	m := internal.DeserializeMap(internal.RawBytePtrToByteSlice(mapData, mapSize))
	active := internal.VMStateGetActiveContextID()
	stream := h.httpStreams[active]

	switch mapType {
	case internal.MapTypeHttpRequestHeaders:
		stream.requestHeaders = m
	case internal.MapTypeHttpResponseHeaders:
		stream.responseHeaders = m
	case internal.MapTypeHttpRequestTrailers:
		stream.requestTrailers = m
	case internal.MapTypeHttpResponseTrailers:
		stream.responseTrailers = m
	default:
		panic("unimplemented")
	}
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (h *httpHostEmulator) ProxyContinueStream(internal.StreamType) internal.Status {
	active := internal.VMStateGetActiveContextID()
	stream := h.httpStreams[active]
	stream.action = types.ActionContinue
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (h *httpHostEmulator) ProxySendLocalResponse(statusCode uint32,
	statusCodeDetailData *byte, statusCodeDetailsSize int, bodyData *byte, bodySize int,
	headersData *byte, headersSize int, grpcStatus int32) internal.Status {
	active := internal.VMStateGetActiveContextID()
	stream := h.httpStreams[active]
	stream.sentLocalResponse = &LocalHttpResponse{
		StatusCode:       statusCode,
		StatusCodeDetail: internal.RawBytePtrToString(statusCodeDetailData, statusCodeDetailsSize),
		Data:             internal.RawBytePtrToByteSlice(bodyData, bodySize),
		Headers:          internal.DeserializeMap(internal.RawBytePtrToByteSlice(headersData, headersSize)),
		GRPCStatus:       grpcStatus,
	}
	return internal.StatusOK
}

// impl HostEmulator
func (h *httpHostEmulator) InitializeHttpContext() (contextID uint32) {
	contextID = getNextContextID()
	internal.ProxyOnContextCreate(contextID, PluginContextID)
	h.httpStreams[contextID] = &httpStreamState{action: types.ActionContinue}
	return
}

// impl HostEmulator
func (h *httpHostEmulator) CallOnRequestHeaders(contextID uint32, headers [][2]string, endOfStream bool) types.Action {
	cs, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.requestHeaders = headers
	cs.action = internal.ProxyOnRequestHeaders(contextID,
		len(headers), endOfStream)
	return cs.action
}

// impl HostEmulator
func (h *httpHostEmulator) CallOnResponseHeaders(contextID uint32, headers [][2]string, endOfStream bool) types.Action {
	cs, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.responseHeaders = headers
	cs.action = internal.ProxyOnResponseHeaders(contextID, len(headers), endOfStream)
	return cs.action
}

// impl HostEmulator
func (h *httpHostEmulator) CallOnRequestTrailers(contextID uint32, trailers [][2]string) types.Action {
	cs, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.requestTrailers = trailers
	cs.action = internal.ProxyOnRequestTrailers(contextID, len(trailers))
	return cs.action
}

// impl HostEmulator
func (h *httpHostEmulator) CallOnResponseTrailers(contextID uint32, trailers [][2]string) types.Action {
	cs, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.responseTrailers = trailers
	cs.action = internal.ProxyOnResponseTrailers(contextID, len(trailers))
	return cs.action
}

// impl HostEmulator
func (h *httpHostEmulator) CallOnRequestBody(contextID uint32, body []byte, endOfStream bool) types.Action {
	cs, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.requestBody = append(cs.requestBody, body...)
	cs.action = internal.ProxyOnRequestBody(contextID,
		len(body), endOfStream)
	return cs.action
}

// impl HostEmulator
func (h *httpHostEmulator) CallOnResponseBody(contextID uint32, body []byte, endOfStream bool) types.Action {
	cs, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.responseBody = body
	cs.action = internal.ProxyOnResponseBody(contextID,
		len(body), endOfStream)
	return cs.action
}

// impl HostEmulator
func (h *httpHostEmulator) CompleteHttpContext(contextID uint32) {
	internal.ProxyOnLog(contextID)
	internal.ProxyOnDelete(contextID)
}

// impl HostEmulator
func (h *httpHostEmulator) GetCurrentHttpStreamAction(contextID uint32) types.Action {
	stream, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}
	return stream.action
}

// impl HostEmulator
func (h *httpHostEmulator) GetCurrentRequestHeaders(contextID uint32) [][2]string {
	stream, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}
	return stream.requestHeaders
}

// impl HostEmulator
func (h *httpHostEmulator) GetCurrentRequestBody(contextID uint32) []byte {
	stream, ok := h.httpStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}
	return stream.requestBody
}

// impl HostEmulator
func (h *httpHostEmulator) GetSentLocalResponse(contextID uint32) *LocalHttpResponse {
	return h.httpStreams[contextID].sentLocalResponse
}
