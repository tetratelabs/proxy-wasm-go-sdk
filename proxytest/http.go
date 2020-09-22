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
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/rawhostcall"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type HttpFilterHost struct {
	*baseHost

	newContext func(contextID uint32) proxywasm.HttpContext
	contexts   map[uint32]*httpContextState
}

type httpContextState struct {
	context proxywasm.HttpContext
	requestHeaders, responseHeaders,
	requestTrailers, responseTrailers [][2]string
	requestBody, responseBody []byte

	action            types.Action
	sentLocalResponse *LocalHttpResponse
}

type LocalHttpResponse struct {
	StatusCode       uint32
	StatusCodeDetail string
	Data             []byte
	Headers          [][2]string
	GRPCStatus       int32
}

func NewHttpFilterHost(f func(contextID uint32) proxywasm.HttpContext) (*HttpFilterHost, func()) {
	host := &HttpFilterHost{
		newContext: f,
		contexts:   map[uint32]*httpContextState{},
	}

	host.baseHost = newBaseHost(func(contextID uint32, numHeaders, bodySize, numTrailers int) {
		ctx, ok := host.contexts[contextID]
		if !ok {
			log.Fatalf("invalid context id for callback: %d", contextID)
		}

		ctx.context.OnHttpCallResponse(numHeaders, bodySize, numTrailers)
	})
	hostMux.Lock()
	rawhostcall.RegisterMockWASMHost(host)
	return host, func() {
		hostMux.Unlock()
	}
}

func (h *HttpFilterHost) InitContext() uint32 {
	contextID := uint32(len(h.contexts)) + 1
	ctx := h.newContext(contextID)

	h.contexts[contextID] = &httpContextState{
		context: ctx,
		action:  types.ActionContinue,
	}
	return contextID
}

func (h *HttpFilterHost) GetAction(contextID uint32) types.Action {
	cs, ok := h.contexts[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}
	return cs.action
}

func (h *HttpFilterHost) PutRequestHeaders(contextID uint32, headers [][2]string) {
	cs, ok := h.contexts[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.requestHeaders = headers
	h.currentContextID = contextID
	cs.action = cs.context.OnHttpRequestHeaders(len(headers), false) // TODO: allow for specifying end_of_stream
}

func (h *HttpFilterHost) PutResponseHeaders(contextID uint32, headers [][2]string) {
	cs, ok := h.contexts[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.responseHeaders = headers
	h.currentContextID = contextID
	cs.action = cs.context.OnHttpResponseHeaders(len(headers), false) // TODO: allow for specifying end_of_stream
}

func (h *HttpFilterHost) PutRequestTrailers(contextID uint32, headers [][2]string) {
	cs, ok := h.contexts[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.requestTrailers = headers
	h.currentContextID = contextID
	cs.action = cs.context.OnHttpRequestTrailers(len(headers))
}

func (h *HttpFilterHost) PutResponseTrailers(contextID uint32, headers [][2]string) {
	cs, ok := h.contexts[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.responseTrailers = headers
	h.currentContextID = contextID
	cs.action = cs.context.OnHttpResponseTrailers(len(headers))
}

func (h *HttpFilterHost) PutRequestBody(contextID uint32, body []byte) {
	cs, ok := h.contexts[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.requestBody = body
	h.currentContextID = contextID
	cs.action = cs.context.OnHttpRequestBody(len(body), false) // TODO: allow for specifying end_of_stream
}

func (h *HttpFilterHost) PutResponseBody(contextID uint32, body []byte) {
	cs, ok := h.contexts[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	cs.responseBody = body
	h.currentContextID = contextID
	cs.action = cs.context.OnHttpResponseBody(len(body), false) // TODO: allow for specifying end_of_stream
}

func (h *HttpFilterHost) ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) types.Status {
	ctx := h.contexts[h.currentContextID]
	var buf []byte
	switch bt {
	case types.BufferTypeHttpRequestBody:
		buf = ctx.requestBody
	case types.BufferTypeHttpResponseBody:
		buf = ctx.requestBody
	default:
		// delegate to baseHost
		return h.getBuffer(bt, start, maxSize, returnBufferData, returnBufferSize)
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

func (h *HttpFilterHost) ProxyGetHeaderMapValue(mapType types.MapType, keyData *byte,
	keySize int, returnValueData **byte, returnValueSize *int) types.Status {
	key := proxywasm.RawBytePtrToString(keyData, keySize)
	ctx := h.contexts[h.currentContextID]

	var headers [][2]string
	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		headers = ctx.requestHeaders
	case types.MapTypeHttpResponseHeaders:
		headers = ctx.responseHeaders
	case types.MapTypeHttpRequestTrailers:
		headers = ctx.requestTrailers
	case types.MapTypeHttpResponseTrailers:
		headers = ctx.responseTrailers
	default:
		return h.getMapValue(mapType, keyData, keySize, returnValueData, returnValueSize)
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

func (h *HttpFilterHost) ProxyAddHeaderMapValue(mapType types.MapType, keyData *byte,
	keySize int, valueData *byte, valueSize int) types.Status {

	key := proxywasm.RawBytePtrToString(keyData, keySize)
	value := proxywasm.RawBytePtrToString(valueData, valueSize)
	ctx := h.contexts[h.currentContextID]

	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		ctx.requestHeaders = addMapValue(ctx.requestHeaders, key, value)
	case types.MapTypeHttpResponseHeaders:
		ctx.responseHeaders = addMapValue(ctx.responseHeaders, key, value)
	case types.MapTypeHttpRequestTrailers:
		ctx.requestTrailers = addMapValue(ctx.requestTrailers, key, value)
	case types.MapTypeHttpResponseTrailers:
		ctx.responseTrailers = addMapValue(ctx.responseTrailers, key, value)
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

func (h *HttpFilterHost) ProxyReplaceHeaderMapValue(mapType types.MapType, keyData *byte,
	keySize int, valueData *byte, valueSize int) types.Status {
	key := proxywasm.RawBytePtrToString(keyData, keySize)
	value := proxywasm.RawBytePtrToString(valueData, valueSize)
	ctx := h.contexts[h.currentContextID]

	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		ctx.requestHeaders = replaceMapValue(ctx.requestHeaders, key, value)
	case types.MapTypeHttpResponseHeaders:
		ctx.responseHeaders = replaceMapValue(ctx.responseHeaders, key, value)
	case types.MapTypeHttpRequestTrailers:
		ctx.requestTrailers = replaceMapValue(ctx.requestTrailers, key, value)
	case types.MapTypeHttpResponseTrailers:
		ctx.responseTrailers = replaceMapValue(ctx.responseTrailers, key, value)
	default:
		panic("unimplemented")
	}
	return types.StatusOK
}

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

func (h *HttpFilterHost) ProxyRemoveHeaderMapValue(mapType types.MapType, keyData *byte, keySize int) types.Status {
	key := proxywasm.RawBytePtrToString(keyData, keySize)
	ctx := h.contexts[h.currentContextID]
	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		ctx.requestHeaders = removeHeaderMapValue(ctx.requestHeaders, key)
	case types.MapTypeHttpResponseHeaders:
		ctx.responseHeaders = removeHeaderMapValue(ctx.responseHeaders, key)
	case types.MapTypeHttpRequestTrailers:
		ctx.requestTrailers = removeHeaderMapValue(ctx.requestTrailers, key)
	case types.MapTypeHttpResponseTrailers:
		ctx.responseTrailers = removeHeaderMapValue(ctx.responseTrailers, key)
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

func (h *HttpFilterHost) ProxyGetHeaderMapPairs(mapType types.MapType, returnValueData **byte,
	returnValueSize *int) types.Status {
	ctx := h.contexts[h.currentContextID]

	var m []byte
	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		m = proxywasm.SerializeMap(ctx.requestHeaders)
	case types.MapTypeHttpResponseHeaders:
		m = proxywasm.SerializeMap(ctx.responseHeaders)
	case types.MapTypeHttpRequestTrailers:
		m = proxywasm.SerializeMap(ctx.requestTrailers)
	case types.MapTypeHttpResponseTrailers:
		m = proxywasm.SerializeMap(ctx.responseTrailers)
	default:
		return h.getMapPairs(mapType, returnValueData, returnValueSize)
	}

	*returnValueData = &m[0]
	*returnValueSize = len(m)
	return types.StatusOK
}

func (h *HttpFilterHost) ProxySetHeaderMapPairs(mapType types.MapType, mapData *byte, mapSize int) types.Status {
	m := proxywasm.DeserializeMap(proxywasm.RawBytePtrToByteSlice(mapData, mapSize))
	ctx := h.contexts[h.currentContextID]
	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		ctx.requestHeaders = m
	case types.MapTypeHttpResponseHeaders:
		ctx.responseHeaders = m
	case types.MapTypeHttpRequestTrailers:
		ctx.requestTrailers = m
	case types.MapTypeHttpResponseTrailers:
		ctx.responseTrailers = m
	default:
		panic("unimplemented")
	}
	return types.StatusOK
}

func (h *HttpFilterHost) ProxyContinueStream(types.StreamType) types.Status {
	ctx := h.contexts[h.currentContextID]
	ctx.action = types.ActionContinue
	return types.StatusOK
}

func (h *HttpFilterHost) GetCurrentAction(contextID uint32) types.Action {
	ctx, ok := h.contexts[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}
	return ctx.action
}

func (h *HttpFilterHost) ProxySendLocalResponse(statusCode uint32,
	statusCodeDetailData *byte, statusCodeDetailsSize int, bodyData *byte, bodySize int,
	headersData *byte, headersSize int, grpcStatus int32) types.Status {
	ctx := h.contexts[h.currentContextID]
	ctx.sentLocalResponse = &LocalHttpResponse{
		StatusCode:       statusCode,
		StatusCodeDetail: proxywasm.RawBytePtrToString(statusCodeDetailData, statusCodeDetailsSize),
		Data:             proxywasm.RawBytePtrToByteSlice(bodyData, bodySize),
		Headers:          proxywasm.DeserializeMap(proxywasm.RawBytePtrToByteSlice(headersData, headersSize)),
		GRPCStatus:       grpcStatus,
	}
	return types.StatusOK
}

func (h *HttpFilterHost) GetSentLocalResponse(contextID uint32) *LocalHttpResponse {
	return h.contexts[contextID].sentLocalResponse
}

func (h *HttpFilterHost) GetContext(contextID uint32) proxywasm.HttpContext {
	return h.contexts[contextID].context
}
