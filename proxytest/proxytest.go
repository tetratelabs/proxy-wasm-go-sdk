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
	"fmt"
	"log"
	"sync"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/rawhostcall"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type HostEmulator interface {
	Done()

	// Root
	StartVM() types.OnVMStartStatus
	StartPlugin() types.OnPluginStartStatus
	FinishVM() bool
	GetCalloutAttributesFromContext(contextID uint32) []HttpCalloutAttribute
	CallOnHttpCallResponse(contextID uint32, headers types.Headers, trailers types.Trailers, body []byte)
	GetCounterMetric(name string) (uint64, error)
	GetGaugeMetric(name string) (uint64, error)
	GetHistogramMetric(name string) (uint64, error)
	GetLogs(level types.LogLevel) []string
	GetTickPeriod() uint32
	Tick()
	GetQueueSize(queueID uint32) int
	RegisterForeignFunction(name string, f func([]byte) []byte)

	// network
	InitializeConnection() (contextID uint32, action types.Action)
	CallOnUpstreamData(contextID uint32, data []byte) types.Action
	CallOnDownstreamData(contextID uint32, data []byte) types.Action
	CloseUpstreamConnection(contextID uint32)
	CloseDownstreamConnection(contextID uint32)
	CompleteConnection(contextID uint32)

	// http
	InitializeHttpContext() (contextID uint32)
	CallOnResponseHeaders(contextID uint32, headers types.Headers, endOfStream bool) types.Action
	CallOnResponseBody(contextID uint32, body []byte, endOfStream bool) types.Action
	CallOnResponseTrailers(contextID uint32, trailers types.Trailers) types.Action
	CallOnRequestHeaders(contextID uint32, headers types.Headers, endOfStream bool) types.Action
	CallOnRequestTrailers(contextID uint32, trailers types.Trailers) types.Action
	CallOnRequestBody(contextID uint32, body []byte, endOfStream bool) types.Action
	CompleteHttpContext(contextID uint32)
	GetCurrentHttpStreamAction(contextID uint32) types.Action
	GetCurrentRequestHeaders(contextID uint32) types.Headers
	GetSentLocalResponse(contextID uint32) *LocalHttpResponse
	CallOnLogForAccessLogger(requestHeaders, responseHeaders types.Headers)
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
		panic(fmt.Sprintf("buffer type %d is not supported by proxytest frame work yet", bt))
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
