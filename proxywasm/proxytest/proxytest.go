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
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

// HostEmulator implements the host side of proxy-wasm. Methods on HostEmulator will either invoke
// methods in the plugin or return state about the host itself, reflecting any changes from previous
// plugin invocations.
type HostEmulator interface {
	// StartVM executes types.VMContext.OnVMStart in the plugin.
	StartVM() types.OnVMStartStatus
	// StartPlugin executes types.PluginContext.OnPluginStart in the plugin.
	StartPlugin() types.OnPluginStartStatus
	// FinishVM executes types.PluginContext.OnPluginDone in the plugin.
	FinishVM() bool
	// GetCalloutAttributesFromContext returns the current HTTP callout attributes for the given HTTP context in the
	// host.
	GetCalloutAttributesFromContext(contextID uint32) []HttpCalloutAttribute
	// CallOnHttpCallResponse executes the callback for the HTTP call with ID calloutID in the plugin.
	CallOnHttpCallResponse(calloutID uint32, headers [][2]string, trailers [][2]string, body []byte)
	// GetCounterMetric returns the value for the counter in the host.
	GetCounterMetric(name string) (uint64, error)
	// GetGaugeMetric returns the value for the gauge in the host.
	GetGaugeMetric(name string) (uint64, error)
	// GetHistogramMetric returns the value for the histogram in the host.
	GetHistogramMetric(name string) (uint64, error)
	// GetTraceLogs returns the trace logs that have been collected in the host.
	GetTraceLogs() []string
	// GetDebugLogs returns the debug logs that have been collected in the host.
	GetDebugLogs() []string
	// GetInfoLogs returns the info logs that have been collected in the host.
	GetInfoLogs() []string
	// GetWarnLogs returns the warn logs that have been collected in the host.
	GetWarnLogs() []string
	// GetErrorLogs returns the error logs that have been collected in the host.
	GetErrorLogs() []string
	// GetCriticalLogs returns the critical logs that have been collected in the host.
	GetCriticalLogs() []string
	// GetTickPeriod returns the current tick period in the host.
	GetTickPeriod() uint32
	// Tick executes types.PluginContext.OnTick in the plugin.
	Tick()
	// GetQueueSize gets the current size of the queue in the host.
	GetQueueSize(queueID uint32) int
	// RegisterForeignFunction registers the foreign function in the host.
	RegisterForeignFunction(name string, f func([]byte) []byte)

	// InitializeConnection executes types.TcpContext.OnNewConnection in the plugin.
	InitializeConnection() (contextID uint32, action types.Action)
	// CallOnUpstreamData executes types.TcpContext.OnUpstreamData in the plugin.
	CallOnUpstreamData(contextID uint32, data []byte) types.Action
	// CallOnDownstreamData executes types.TcpContext.OnDownstreamData in the plugin.
	CallOnDownstreamData(contextID uint32, data []byte) types.Action
	// CloseUpstreamConnection executes types.TcpContext.OnUpstreamClose in the plugin.
	CloseUpstreamConnection(contextID uint32)
	// CloseDownstreamConnection executes types.TcpContext.OnDownstreamClose in the plugin.
	CloseDownstreamConnection(contextID uint32)
	// CompleteConnection executes types.TcpContext.OnStreamDone in the plugin.
	CompleteConnection(contextID uint32)

	// InitializeHttpContext executes types.PluginContext.NewHttpContext in the plugin.
	InitializeHttpContext() (contextID uint32)
	// CallOnResponseHeaders executes types.HttpContext.OnHttpResponseHeaders in the plugin.
	// The number of headers and endOfStream are passed to the plugin and the content of headers are visible in
	// the plugin for methods like proxywasm.GetHttpResponseHeaders.
	CallOnResponseHeaders(contextID uint32, headers [][2]string, endOfStream bool) types.Action
	// CallOnResponseBody executes types.HttpContext.OnHttpResponseBody in the plugin.
	// The number of bytes and endOfStream are passed to the plugin and the content of bytes are visible in the plugin
	// for methods like proxywasm.GetHttpResponseBody.
	CallOnResponseBody(contextID uint32, body []byte, endOfStream bool) types.Action
	// CallOnResponseTrailers executes types.HttpContext.OnHttpResponseTrailers in the plugin.
	// The number of trailers and endOfStream are passed to the plugin and the content of trailers are visible in the
	// plugin for methods like proxywasm.GetHttpResponseTrailers.
	CallOnResponseTrailers(contextID uint32, trailers [][2]string) types.Action
	// CallOnRequestHeaders executes types.HttpContext.OnHttpRequestHeaders in the plugin.
	// The number of headers and endOfStream are passed to the plugin and the content of headers are visible in the
	// plugin for methods like proxywasm.GetHttpRequestHeaders.
	CallOnRequestHeaders(contextID uint32, headers [][2]string, endOfStream bool) types.Action
	// CallOnRequestTrailers executes types.HttpContext.OnHttpRequestTrailers in the plugin.
	// The number of trailers and endOfStream are passed to the plugin and the content of trailers are visible in the
	// plugin for methods like proxywasm.GetHttpRequestTrailers.
	CallOnRequestTrailers(contextID uint32, trailers [][2]string) types.Action
	// CallOnRequestBody executes types.HttpContext.OnHttpRequestBody in the plugin.
	// The number of bytes and endOfStream are passed to the plugin and the content of bytes are visible in the plugin
	// for methods like proxywasm.GetHttpRequestBody.
	CallOnRequestBody(contextID uint32, body []byte, endOfStream bool) types.Action
	// CompleteHttpContext executes types.HttpContext.OnHttpStreamDone in the plugin.
	CompleteHttpContext(contextID uint32)
	// GetCurrentHttpStreamAction returns the current action for the HTTP stream with ID contextID in the host.
	// This is the return value of a previous lifecycle call in the plugin.
	GetCurrentHttpStreamAction(contextID uint32) types.Action
	// GetCurrentRequestHeaders returns the current request headers for the HTTP stream with ID contextID in the host.
	// This will reflect any mutations made by the plugin such as with proxywasm.AddHttpRequestHeader.
	GetCurrentRequestHeaders(contextID uint32) [][2]string
	// GetCurrentResponseHeaders returns the current response headers for the HTTP stream with ID contextID in the host.
	// This will reflect any mutations made by the plugin such as with proxywasm.AddHttpResponseHeader.
	GetCurrentResponseHeaders(contextID uint32) [][2]string
	// GetCurrentRequestBody returns the current request body for the HTTP stream with ID contextID in the host.
	// This will reflect any mutations made by th eplugin such as with proxywasm.AppendHttpRequestBody.
	GetCurrentRequestBody(contextID uint32) []byte
	// GetCurrentResponseBody returns the current response body for the HTTP stream with ID contextID in the host.
	// This will reflect any mutations made by the plugin such as with proxywasm.AppendHttpResponseBody.
	GetCurrentResponseBody(contextID uint32) []byte
	// GetSentLocalResponse returns the local response that has been sent for the HTTP stream with ID contextID in the
	// host. This contains the arguments passed to proxywasm.SendHttpResponse in the plugin. If
	// proxywasm.SendHttpResponse hasn't been invoked by the plugin, this will return nil.
	GetSentLocalResponse(contextID uint32) *LocalHttpResponse
	// GetProperty returns property data from the host, for a given path.
	GetProperty(path []string) ([]byte, error)
	// SetProperty sets property data on the host, for a given path.
	SetProperty(path []string, data []byte) error
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
	properties         map[string][]byte
}

// NewHostEmulator returns a new HostEmulator that can be used to test a plugin. Plugin tests will
// often involve calling methods on HostEmulator to invoke methods in the plugin while checking
// the state within the host after plugin execution.
func NewHostEmulator(opt *EmulatorOption) (host HostEmulator, reset func()) {
	root := newRootHostEmulator(opt.pluginConfiguration, opt.vmConfiguration)
	network := newNetworkHostEmulator()
	http := newHttpHostEmulator()
	emulator := &hostEmulator{
		root,
		network,
		http,
		0,
		make(map[string][]byte),
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

func cloneWithLowerCaseMapKeys(m [][2]string) [][2]string {
	r := make([][2]string, len(m))
	for i, entry := range m {
		r[i] = [2]string{strings.ToLower(entry[0]), entry[1]}
	}
	return r
}

func deserializeRawBytePtrToMap(aw *byte, size int) [][2]string {
	m := internal.DeserializeMap(internal.RawBytePtrToByteSlice(aw, size))
	for _, entry := range m {
		entry[0] = strings.ToLower(entry[0])
	}
	return m
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
func (h *hostEmulator) ProxySetProperty(pathPtr *byte, pathSize int, dataPtr *byte, dataSize int) internal.Status {
	path := internal.RawBytePtrToString(pathPtr, pathSize)
	data := internal.RawBytePtrToByteSlice(dataPtr, dataSize)
	h.properties[path] = data
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (h *hostEmulator) ProxyGetProperty(pathPtr *byte, pathSize int, dataPtrPtr **byte, dataSizePtr *int) internal.Status {
	path := internal.RawBytePtrToString(pathPtr, pathSize)
	if _, ok := h.properties[path]; !ok {
		return internal.StatusNotFound
	}
	data := h.properties[path]
	*dataPtrPtr = &data[0]
	dataSize := len(data)
	*dataSizePtr = dataSize
	return internal.StatusOK
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
