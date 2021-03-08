// Copyright 2021 Tetratea
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxytest

import (
	"log"
	"sync"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/rawhostcall"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type HostEmulator interface {
	Done()

	// Root
	StartVM()
	StartPlugin()
	FinishVM()

	GetCalloutAttributesFromContext(contextID uint32) []HttpCalloutAttribute
	PutCalloutResponse(contextID uint32, headers, trailers [][2]string, body []byte)

	GetLogs(level types.LogLevel) []string
	GetTickPeriod() uint32
	Tick()
	GetQueueSize(queueID uint32) int

	// network
	NetworkFilterInitConnection() (contextID uint32)
	NetworkFilterPutUpstreamData(contextID uint32, data []byte)
	NetworkFilterPutDownstreamData(contextID uint32, data []byte)
	NetworkFilterCloseUpstreamConnection(contextID uint32)
	NetworkFilterCloseDownstreamConnection(contextID uint32)
	NetworkFilterCompleteConnection(contextID uint32)

	// http
	HttpFilterInitContext() (contextID uint32)
	HttpFilterPutRequestHeaders(contextID uint32, headers [][2]string)
	HttpFilterGetRequestHeaders(contextID uint32) (headers [][2]string)
	HttpFilterPutRequestHeadersEndOfStream(contextID uint32, headers [][2]string, endOfStream bool)
	HttpFilterPutResponseHeaders(contextID uint32, headers [][2]string)
	HttpFilterGetResponseHeaders(contextID uint32) (headers [][2]string)
	HttpFilterPutResponseHeadersEndOfStream(contextID uint32, headers [][2]string, endOfStream bool)
	HttpFilterPutRequestTrailers(contextID uint32, headers [][2]string)
	HttpFilterPutResponseTrailers(contextID uint32, headers [][2]string)
	HttpFilterPutRequestBody(contextID uint32, body []byte)
	HttpFilterPutRequestBodyEndOfStream(contextID uint32, body []byte, endOfStream bool)
	HttpFilterGetRequestBody(contextID uint32) []byte
	HttpFilterPutResponseBody(contextID uint32, body []byte)
	HttpFilterPutResponseBodyEndOfStream(contextID uint32, body []byte, endOfStream bool)
	HttpFilterGetResponseBody(contextID uint32) []byte
	HttpFilterCompleteHttpStream(contextID uint32)
	HttpFilterGetCurrentStreamAction(contextID uint32) types.Action
	HttpFilterGetSentLocalResponse(contextID uint32) *LocalHttpResponse
	CallOnLogForAccessLogger(requestHeaders, responseHeaders [][2]string)
}

const (
	RootContextID uint32 = 1 // TODO: support multiple rootContext
)

var (
	hostMux       = sync.Mutex{}
	nextContextID = RootContextID + 1
)

type hostEmulator struct {
	*rootHostEmulator
	*networkHostEmulator
	*httpHostEmulator

	effectiveContextID uint32
}

func NewHostEmulator(opt *EmulatorOption) HostEmulator {
	root := newRootHostEmulator(opt.pluginConfiguration, opt.vmConfiguration)
	network := newNetworkHostEmulator()
	http := newHttpHostEmulator()
	emulator := &hostEmulator{
		root,
		network,
		http,
		0,
	}

	hostMux.Lock() // acquire the lock of host emulation
	rawhostcall.RegisterMockWASMHost(emulator)

	// set up state
	proxywasm.SetNewRootContext(opt.newRootContext)
	proxywasm.SetNewStreamContext(opt.newStreamContext)
	proxywasm.SetNewHttpContext(opt.newHttpContext)

	// create root context: TODO: support multiple root contexts
	proxywasm.ProxyOnContextCreate(RootContextID, 0)

	return emulator
}

func getNextContextID() (ret uint32) {
	ret = nextContextID
	nextContextID++
	return
}

// impl HostEmulator
func (*hostEmulator) Done() {
	defer hostMux.Unlock()
	defer proxywasm.VMStateReset()
}

// impl rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) types.Status {
	switch bt {
	case types.BufferTypePluginConfiguration, types.BufferTypeVMConfiguration, types.BufferTypeHttpCallResponseBody:
		return h.rootHostEmulatorProxyGetBufferBytes(bt, start, maxSize, returnBufferData, returnBufferSize)
	case types.BufferTypeDownstreamData, types.BufferTypeUpstreamData:
		return h.networkHostEmulatorProxyGetBufferBytes(bt, start, maxSize, returnBufferData, returnBufferSize)
	case types.BufferTypeHttpRequestBody, types.BufferTypeHttpResponseBody:
		return h.httpHostEmulatorProxyGetBufferBytes(bt, start, maxSize, returnBufferData, returnBufferSize)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}
}

func (h *hostEmulator) ProxySetBufferBytes(bt types.BufferType, start int, maxSize int, bufferData *byte, bufferSize int) types.Status {
	switch bt {
	case types.BufferTypeHttpRequestBody, types.BufferTypeHttpResponseBody:
		return h.httpHostEmulatorProxySetBufferBytes(bt, start, maxSize, bufferData, bufferSize)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}
}

// impl rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxyGetHeaderMapValue(mapType types.MapType, keyData *byte,
	keySize int, returnValueData **byte, returnValueSize *int) types.Status {
	switch mapType {
	case types.MapTypeHttpRequestHeaders, types.MapTypeHttpResponseHeaders,
		types.MapTypeHttpRequestTrailers, types.MapTypeHttpResponseTrailers:
		return h.httpHostEmulatorProxyGetHeaderMapValue(mapType, keyData,
			keySize, returnValueData, returnValueSize)
	case types.MapTypeHttpCallResponseHeaders, types.MapTypeHttpCallResponseTrailers:
		return h.rootHostEmulatorProxyGetMapValue(mapType, keyData,
			keySize, returnValueData, returnValueSize)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}
}

// impl rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxyGetHeaderMapPairs(mapType types.MapType, returnValueData **byte,
	returnValueSize *int) types.Status {
	switch mapType {
	case types.MapTypeHttpRequestHeaders, types.MapTypeHttpResponseHeaders,
		types.MapTypeHttpRequestTrailers, types.MapTypeHttpResponseTrailers:
		return h.httpHostEmulatorProxyGetHeaderMapPairs(mapType, returnValueData, returnValueSize)
	case types.MapTypeHttpCallResponseHeaders, types.MapTypeHttpCallResponseTrailers:
		return h.rootHostEmulatorProxyGetHeaderMapPairs(mapType, returnValueData, returnValueSize)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}
}

// impl rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxySetEffectiveContext(contextID uint32) types.Status {
	h.effectiveContextID = contextID
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxySetProperty(*byte, int, *byte, int) types.Status {
	panic("unimplemented")
}

// impl rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxyGetProperty(*byte, int, **byte, *int) types.Status {
	log.Printf("ProxyGetProperty not implemented in the host emulator yet")
	return 0
}

// impl rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxyResolveSharedQueue(vmIDData *byte, vmIDSize int, nameData *byte, nameSize int, returnID *uint32) types.Status {
	log.Printf("ProxyResolveSharedQueue not implemented in the host emulator yet")
	return 0
}

// impl rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxyCloseStream(streamType types.StreamType) types.Status {
	log.Printf("ProxyCloseStream not implemented in the host emulator yet")
	return 0
}

// impl rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxyDone() types.Status {
	log.Printf("ProxyDone not implemented in the host emulator yet")
	return 0
}
