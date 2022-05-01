// Copyright 2021 Tetrate
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

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type HostEmulator interface {
	// Root
	StartVM() types.OnVMStartStatus
	StartPlugin() types.OnPluginStartStatus
	FinishVM() bool
	GetCalloutAttributesFromContext(contextID uint32) []HttpCalloutAttribute
	CallOnHttpCallResponse(contextID uint32, headers [][2]string, trailers [][2]string, body []byte)
	GetCounterMetric(name string) (uint64, error)
	GetGaugeMetric(name string) (uint64, error)
	GetHistogramMetric(name string) (uint64, error)
	GetTraceLogs() []string
	GetDebugLogs() []string
	GetInfoLogs() []string
	GetWarnLogs() []string
	GetErrorLogs() []string
	GetCriticalLogs() []string
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
	CallOnResponseHeaders(contextID uint32, headers [][2]string, endOfStream bool) types.Action
	CallOnResponseBody(contextID uint32, body []byte, endOfStream bool) types.Action
	CallOnResponseTrailers(contextID uint32, trailers [][2]string) types.Action
	CallOnRequestHeaders(contextID uint32, headers [][2]string, endOfStream bool) types.Action
	CallOnRequestTrailers(contextID uint32, trailers [][2]string) types.Action
	CallOnRequestBody(contextID uint32, body []byte, endOfStream bool) types.Action
	CompleteHttpContext(contextID uint32)
	GetCurrentHttpStreamAction(contextID uint32) types.Action
	GetCurrentRequestHeaders(contextID uint32) [][2]string
	GetCurrentRequestBody(contextID uint32) []byte
	GetSentLocalResponse(contextID uint32) *LocalHttpResponse
}

const (
	PluginContextID uint32 = 1 // TODO: support multiple pluginContext
)

var nextContextID = PluginContextID + 1

type hostEmulator struct {
	*rootHostEmulator
	*networkHostEmulator
	*httpHostEmulator

	effectiveContextID uint32
}

func NewHostEmulator(opt *EmulatorOption) (host HostEmulator, reset func()) {
	root := newRootHostEmulator(opt.pluginConfiguration, opt.vmConfiguration)
	network := newNetworkHostEmulator()
	http := newHttpHostEmulator()
	emulator := &hostEmulator{
		root,
		network,
		http,
		0,
	}

	release := internal.RegisterMockWasmHost(emulator)

	// set up state
	proxywasm.SetVMContext(opt.vmContext)

	// create plugin context: TODO: support multiple plugin contexts
	internal.ProxyOnContextCreate(PluginContextID, 0)

	return emulator, func() {
		defer release()
		defer internal.VMStateReset()
	}
}

func getNextContextID() (ret uint32) {
	ret = nextContextID
	nextContextID++
	return
}

// impl internal.ProxyWasmHost
func (h *hostEmulator) ProxyGetBufferBytes(bt internal.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) internal.Status {
	switch bt {
	case internal.BufferTypePluginConfiguration, internal.BufferTypeVMConfiguration, internal.BufferTypeHttpCallResponseBody:
		return h.rootHostEmulatorProxyGetBufferBytes(bt, start, maxSize, returnBufferData, returnBufferSize)
	case internal.BufferTypeDownstreamData, internal.BufferTypeUpstreamData:
		return h.networkHostEmulatorProxyGetBufferBytes(bt, start, maxSize, returnBufferData, returnBufferSize)
	case internal.BufferTypeHttpRequestBody, internal.BufferTypeHttpResponseBody:
		return h.httpHostEmulatorProxyGetBufferBytes(bt, start, maxSize, returnBufferData, returnBufferSize)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}
}

func (h *hostEmulator) ProxySetBufferBytes(bt internal.BufferType, start int, maxSize int, bufferData *byte, bufferSize int) (ret internal.Status) {
	switch bt {
	case internal.BufferTypeHttpRequestBody, internal.BufferTypeHttpResponseBody:
		ret = h.httpHostEmulatorProxySetBufferBytes(bt, start, maxSize, bufferData, bufferSize)
	default:
		panic(fmt.Sprintf("buffer type %d is not supported by proxytest frame work yet", bt))
	}
	return
}

// impl internal.ProxyWasmHost
func (h *hostEmulator) ProxyGetHeaderMapValue(mapType internal.MapType, keyData *byte,
	keySize int, returnValueData **byte, returnValueSize *int) internal.Status {
	switch mapType {
	case internal.MapTypeHttpRequestHeaders, internal.MapTypeHttpResponseHeaders,
		internal.MapTypeHttpRequestTrailers, internal.MapTypeHttpResponseTrailers:
		return h.httpHostEmulatorProxyGetHeaderMapValue(mapType, keyData,
			keySize, returnValueData, returnValueSize)
	case internal.MapTypeHttpCallResponseHeaders, internal.MapTypeHttpCallResponseTrailers:
		return h.rootHostEmulatorProxyGetMapValue(mapType, keyData,
			keySize, returnValueData, returnValueSize)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}
}

// impl internal.ProxyWasmHost
func (h *hostEmulator) ProxyGetHeaderMapPairs(mapType internal.MapType, returnValueData **byte,
	returnValueSize *int) internal.Status {
	switch mapType {
	case internal.MapTypeHttpRequestHeaders, internal.MapTypeHttpResponseHeaders,
		internal.MapTypeHttpRequestTrailers, internal.MapTypeHttpResponseTrailers:
		return h.httpHostEmulatorProxyGetHeaderMapPairs(mapType, returnValueData, returnValueSize)
	case internal.MapTypeHttpCallResponseHeaders, internal.MapTypeHttpCallResponseTrailers:
		return h.rootHostEmulatorProxyGetHeaderMapPairs(mapType, returnValueData, returnValueSize)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}
}

// impl internal.ProxyWasmHost
func (h *hostEmulator) ProxySetEffectiveContext(contextID uint32) internal.Status {
	h.effectiveContextID = contextID
	// TODO(ikeeip): This is a workaround. Originally host uses both true context and
	// effective context every time. We should implement this behavior hostEmulator too.
	// see: https://github.com/proxy-wasm/proxy-wasm-cpp-host/blob/f38347360feaaf5b2a733f219c4d8c9660d626f0/src/exports.cc#L23
	internal.VMStateSetActiveContextID(contextID)
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (h *hostEmulator) ProxySetProperty(*byte, int, *byte, int) internal.Status {
	panic("unimplemented")
}

// impl internal.ProxyWasmHost
func (h *hostEmulator) ProxyGetProperty(*byte, int, **byte, *int) internal.Status {
	log.Printf("ProxyGetProperty not implemented in the host emulator yet")
	return 0
}

// impl internal.ProxyWasmHost
func (h *hostEmulator) ProxyResolveSharedQueue(vmIDData *byte, vmIDSize int, nameData *byte, nameSize int, returnID *uint32) internal.Status {
	log.Printf("ProxyResolveSharedQueue not implemented in the host emulator yet")
	return 0
}

// impl internal.ProxyWasmHost
func (h *hostEmulator) ProxyCloseStream(streamType internal.StreamType) internal.Status {
	log.Printf("ProxyCloseStream not implemented in the host emulator yet")
	return 0
}

// impl internal.ProxyWasmHost
func (h *hostEmulator) ProxyDone() internal.Status {
	log.Printf("ProxyDone not implemented in the host emulator yet")
	return 0
}
