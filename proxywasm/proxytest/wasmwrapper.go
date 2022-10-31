// Copyright 2020-2022 Tetrate
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
	"context"
	"io"
	"reflect"
	"unsafe"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type guestABI struct {
	proxyOnVMStart          api.Function
	proxyOnContextCreate    api.Function
	proxyOnConfigure        api.Function
	proxyOnDone             api.Function
	proxyOnQueueReady       api.Function
	proxyOnTick             api.Function
	proxyOnRequestHeaders   api.Function
	proxyOnRequestBody      api.Function
	proxyOnRequestTrailers  api.Function
	proxyOnResponseHeaders  api.Function
	proxyOnResponseBody     api.Function
	proxyOnResponseTrailers api.Function
	proxyOnLog              api.Function
}

// WasmVMContext is a VMContext that delegates execution to a compiled wasm binary.
type WasmVMContext interface {
	types.VMContext
	io.Closer
}

// vmContext implements WasmVMContext.
type vmContext struct {
	runtime wazero.Runtime
	abi     guestABI
	ctx     context.Context
}

// NewWasmVMContext returns a types.VMContext that delegates plugin invocations to the provided compiled wasm binary.
// proxytest can be run with a compiled wasm binary by passing this to proxytest.WithVMContext.
//
// Running proxytest with the compiled wasm binary helps to ensure that the plugin will run when actually compiled with
// TinyGo, however stack traces and other debug features will be much worse. It is recommended to run unit tests both
// with Go and with wasm. Tests will run much faster under Go for quicker development cycles, and the wasm runner can
// confirm the behavior matches when actually compiled.
//
// For example, this snippet allows determining the types.VMContext based on a test case flag.
//
//	var vm types.VMContext
//	switch runner {
//	case "go":
//		vm = &vmContext{}
//	case "wasm":
//		wasm, err := os.ReadFile("plugin.wasm")
//		if err != nil {
//			t.Skip("wasm not found")
//		}
//		v, err := proxytest.NewWasmVMContext(wasm)
//		require.NoError(t, err)
//		vm = v
//	}
//
// Note: Currently only HTTP plugins are supported.
func NewWasmVMContext(wasm []byte) (WasmVMContext, error) {
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)

	_, err := wasi_snapshot_preview1.Instantiate(ctx, r)
	if err != nil {
		return nil, err
	}

	err = exportHostABI(ctx, r)
	if err != nil {
		return nil, err
	}

	mod, err := r.InstantiateModuleFromBinary(ctx, wasm)
	if err != nil {
		return nil, err
	}

	abi := guestABI{
		proxyOnVMStart:          mod.ExportedFunction("proxy_on_vm_start"),
		proxyOnContextCreate:    mod.ExportedFunction("proxy_on_context_create"),
		proxyOnConfigure:        mod.ExportedFunction("proxy_on_configure"),
		proxyOnDone:             mod.ExportedFunction("proxy_on_done"),
		proxyOnQueueReady:       mod.ExportedFunction("proxy_on_queue_ready"),
		proxyOnTick:             mod.ExportedFunction("proxy_on_tick"),
		proxyOnRequestHeaders:   mod.ExportedFunction("proxy_on_request_headers"),
		proxyOnRequestBody:      mod.ExportedFunction("proxy_on_request_body"),
		proxyOnRequestTrailers:  mod.ExportedFunction("proxy_on_request_trailers"),
		proxyOnResponseHeaders:  mod.ExportedFunction("proxy_on_response_headers"),
		proxyOnResponseBody:     mod.ExportedFunction("proxy_on_response_body"),
		proxyOnResponseTrailers: mod.ExportedFunction("proxy_on_response_trailers"),
		proxyOnLog:              mod.ExportedFunction("proxy_on_log"),
	}

	return &vmContext{
		runtime: r,
		abi:     abi,
		ctx:     ctx,
	}, nil
}

// OnVMStart implements the same method on types.VMContext.
func (v *vmContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	rootContextID := uint64(0) // unused
	res, err := v.abi.proxyOnVMStart.Call(v.ctx, rootContextID, uint64(vmConfigurationSize))
	handleErr(err)
	return res[0] == 1
}

// NewPluginContext implements the same method on types.VMContext.
func (v *vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	_, err := v.abi.proxyOnContextCreate.Call(v.ctx, uint64(contextID), 0)
	handleErr(err)
	return &pluginContext{
		id:  uint64(contextID),
		abi: v.abi,
		ctx: withPluginContextID(v.ctx, contextID),
	}
}

// Close implements the same method on io.Closer.
func (v *vmContext) Close() error {
	return v.runtime.Close(v.ctx)
}

// pluginContext implements types.PluginContext.
type pluginContext struct {
	id  uint64
	abi guestABI
	ctx context.Context
}

// OnPluginStart implements the same method on types.PluginContext.
func (p *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	res, err := p.abi.proxyOnConfigure.Call(p.ctx, p.id, uint64(pluginConfigurationSize))
	handleErr(err)
	return res[0] == 1
}

// OnPluginDone implements the same method on types.PluginContext.
func (p *pluginContext) OnPluginDone() bool {
	res, err := p.abi.proxyOnDone.Call(p.ctx, p.id)
	handleErr(err)
	return res[0] == 1
}

// OnQueueReady implements the same method on types.PluginContext.
func (p *pluginContext) OnQueueReady(queueID uint32) {
	_, err := p.abi.proxyOnQueueReady.Call(p.ctx, p.id, uint64(queueID))
	handleErr(err)
}

// OnTick implements the same method on types.PluginContext.
func (p *pluginContext) OnTick() {
	_, err := p.abi.proxyOnTick.Call(p.ctx, p.id)
	handleErr(err)
}

// NewTcpContext implements the same method on types.PluginContext.
func (p *pluginContext) NewTcpContext(uint32) types.TcpContext {
	return nil
}

// NewHttpContext implements the same method on types.PluginContext.
func (p *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	_, err := p.abi.proxyOnContextCreate.Call(p.ctx, uint64(contextID), p.id)
	handleErr(err)
	return &httpContext{
		id:  uint64(contextID),
		abi: p.abi,
		ctx: p.ctx,
	}
}

// httpContext implements types.HttpContext.
type httpContext struct {
	id  uint64
	abi guestABI
	ctx context.Context
}

// OnHttpRequestHeaders implements the same method on types.HttpContext.
func (h *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	res, err := h.abi.proxyOnRequestHeaders.Call(h.ctx, h.id, uint64(numHeaders), wasmBool(endOfStream))
	handleErr(err)
	return types.Action(res[0])
}

// OnHttpRequestBody implements the same method on types.HttpContext.
func (h *httpContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	res, err := h.abi.proxyOnRequestBody.Call(h.ctx, h.id, uint64(bodySize), wasmBool(endOfStream))
	handleErr(err)
	return types.Action(res[0])
}

// OnHttpRequestTrailers implements the same method on types.HttpContext.
func (h *httpContext) OnHttpRequestTrailers(numTrailers int) types.Action {
	res, err := h.abi.proxyOnRequestTrailers.Call(h.ctx, h.id, uint64(numTrailers))
	handleErr(err)
	return types.Action(res[0])
}

// OnHttpResponseHeaders implements the same method on types.HttpContext.
func (h *httpContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	res, err := h.abi.proxyOnResponseHeaders.Call(h.ctx, h.id, uint64(numHeaders), wasmBool(endOfStream))
	handleErr(err)
	return types.Action(res[0])
}

// OnHttpResponseBody implements the same method on types.HttpContext.
func (h *httpContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {
	res, err := h.abi.proxyOnResponseBody.Call(h.ctx, h.id, uint64(bodySize), wasmBool(endOfStream))
	handleErr(err)
	return types.Action(res[0])
}

// OnHttpResponseTrailers implements the same method on types.HttpContext.
func (h *httpContext) OnHttpResponseTrailers(numTrailers int) types.Action {
	res, err := h.abi.proxyOnResponseTrailers.Call(h.ctx, h.id, uint64(numTrailers))
	handleErr(err)
	return types.Action(res[0])
}

// OnHttpStreamDone implements the same method on types.HttpContext.
func (h *httpContext) OnHttpStreamDone() {
	_, err := h.abi.proxyOnLog.Call(h.ctx, h.id)
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func handleMemoryStatus(ok bool) {
	if !ok {
		panic("could not access memory")
	}
}

func wasmBytePtr(ctx context.Context, mod api.Module, off uint32, size uint32) *byte {
	if size == 0 {
		return nil

	}
	buf, ok := mod.Memory().Read(ctx, off, size)
	handleMemoryStatus(ok)
	return &buf[0]
}

func copyBytesToWasm(ctx context.Context, mod api.Module, hostPtr *byte, size int, wasmPtrPtr uint32, wasmSizePtr uint32) {
	if size == 0 {
		return
	}
	var hostSlice []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&hostSlice))
	hdr.Data = uintptr(unsafe.Pointer(hostPtr))
	hdr.Cap = size
	hdr.Len = size

	alloc := mod.ExportedFunction("proxy_on_memory_allocate")
	res, err := alloc.Call(ctx, uint64(size))
	handleErr(err)
	buf, ok := mod.Memory().Read(ctx, uint32(res[0]), uint32(size))
	handleMemoryStatus(ok)

	copy(buf, hostSlice)
	ok = mod.Memory().WriteUint32Le(ctx, wasmPtrPtr, uint32(res[0]))
	handleMemoryStatus(ok)

	ok = mod.Memory().WriteUint32Le(ctx, wasmSizePtr, uint32(size))
	handleMemoryStatus(ok)
}

func wasmBool(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}

func exportHostABI(ctx context.Context, r wazero.Runtime) error {
	_, err := r.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, logLevel uint32, messageData uint32, messageSize uint32) uint32 {
			messageDataPtr := wasmBytePtr(ctx, mod, messageData, messageSize)
			return uint32(internal.ProxyLog(internal.LogLevel(logLevel), messageDataPtr, int(messageSize)))
		}).
		Export("proxy_log").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, pathData uint32, pathSize uint32, valueData uint32, valueSize uint32) uint32 {
			pathDataPtr := wasmBytePtr(ctx, mod, pathData, pathSize)
			valueDataPtr := wasmBytePtr(ctx, mod, valueData, valueSize)
			return uint32(internal.ProxySetProperty(pathDataPtr, int(pathSize), valueDataPtr, int(valueSize)))
		}).
		Export("proxy_set_property").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, pathData uint32, pathSize uint32,
			returnValueData uint32, returnValueSize uint32) uint32 {
			pathDataPtr := wasmBytePtr(ctx, mod, pathData, pathSize)
			var returnValueHostPtr *byte
			var returnValueSizePtr int
			ret := uint32(internal.ProxyGetProperty(pathDataPtr, int(pathSize), &returnValueHostPtr, &returnValueSizePtr))
			copyBytesToWasm(ctx, mod, returnValueHostPtr, returnValueSizePtr, returnValueData, returnValueSize)
			return ret
		}).
		Export("proxy_get_property").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module,
			statusCode uint32, statusCodeDetailData uint32, statusCodeDetailsSize uint32,
			bodyData uint32, bodySize uint32, headersData uint32, headersSize uint32, grpcStatus int32) uint32 {
			statusCodeDetailDataPtr := wasmBytePtr(ctx, mod, statusCodeDetailData, statusCodeDetailsSize)
			bodyDataPtr := wasmBytePtr(ctx, mod, bodyData, bodySize)
			headersDataPtr := wasmBytePtr(ctx, mod, headersData, headersSize)
			return uint32(internal.ProxySendLocalResponse(statusCode, statusCodeDetailDataPtr, int(statusCodeDetailsSize),
				bodyDataPtr, int(bodySize), headersDataPtr, int(headersSize), grpcStatus))
		}).
		Export("proxy_send_local_response").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, keyData uint32, keySize uint32,
			returnValueData uint32, returnValueSize uint32, returnCas uint32) uint32 {
			keyDataPtr := wasmBytePtr(ctx, mod, keyData, keySize)
			var returnValueHostPtr *byte
			var returnValueSizePtr int
			var returnCasPtr uint32
			ret := uint32(internal.ProxyGetSharedData(keyDataPtr, int(keySize), &returnValueHostPtr, &returnValueSizePtr, &returnCasPtr))
			copyBytesToWasm(ctx, mod, returnValueHostPtr, returnValueSizePtr, returnValueData, returnValueSize)
			handleMemoryStatus(mod.Memory().WriteUint32Le(ctx, returnCas, returnCasPtr))
			return ret
		}).
		Export("proxy_get_shared_data").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, keyData uint32, keySize uint32, valueData uint32, valueSize uint32, cas uint32) uint32 {
			keyDataPtr := wasmBytePtr(ctx, mod, keyData, keySize)
			valueDataPtr := wasmBytePtr(ctx, mod, valueData, valueSize)
			return uint32(internal.ProxySetSharedData(keyDataPtr, int(keySize), valueDataPtr, int(valueSize), cas))
		}).
		Export("proxy_set_shared_data").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, nameData uint32, nameSize uint32, returnID uint32) uint32 {
			namePtr := wasmBytePtr(ctx, mod, nameData, nameSize)
			var returnIDPtr uint32
			ret := uint32(internal.ProxyRegisterSharedQueue(namePtr, int(nameSize), &returnIDPtr))
			handleMemoryStatus(mod.Memory().WriteUint32Le(ctx, returnID, returnIDPtr))
			return ret
		}).
		Export("proxy_register_shared_queue").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, vmIDData uint32, vmIDSize uint32, nameData uint32, nameSize uint32, returnID uint32) uint32 {
			vmID := wasmBytePtr(ctx, mod, vmIDData, vmIDSize)
			namePtr := wasmBytePtr(ctx, mod, nameData, nameSize)
			var returnIDPtr uint32
			ret := uint32(internal.ProxyResolveSharedQueue(vmID, int(vmIDSize), namePtr, int(nameSize), &returnIDPtr))
			handleMemoryStatus(mod.Memory().WriteUint32Le(ctx, returnID, returnIDPtr))
			return ret
		}).
		Export("proxy_resolve_shared_queue").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, queueID uint32, returnValueData uint32, returnValueSize uint32) uint32 {
			var returnValueHostPtr *byte
			var returnValueSizePtr int
			ret := uint32(internal.ProxyDequeueSharedQueue(queueID, &returnValueHostPtr, &returnValueSizePtr))
			copyBytesToWasm(ctx, mod, returnValueHostPtr, returnValueSizePtr, returnValueData, returnValueSize)
			return ret
		}).
		Export("proxy_dequeue_shared_queue").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, queueID uint32, valueData uint32, valueSize uint32) uint32 {
			valuePtr := wasmBytePtr(ctx, mod, valueData, valueSize)
			return uint32(internal.ProxyEnqueueSharedQueue(queueID, valuePtr, int(valueSize)))
		}).
		Export("proxy_enqueue_shared_queue").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, mapType uint32, keyData uint32, keySize uint32, returnValueData uint32, returnValueSize uint32) uint32 {
			keyPtr := wasmBytePtr(ctx, mod, keyData, keySize)
			var retValDataHostPtr *byte
			var retValSizePtr int
			ret := uint32(internal.ProxyGetHeaderMapValue(internal.MapType(mapType), keyPtr, int(keySize), &retValDataHostPtr, &retValSizePtr))
			copyBytesToWasm(ctx, mod, retValDataHostPtr, retValSizePtr, returnValueData, returnValueSize)
			return ret
		}).
		Export("proxy_get_header_map_value").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, mapType uint32, keyData uint32, keySize uint32, valueData uint32, valueSize uint32) uint32 {
			keyPtr := wasmBytePtr(ctx, mod, keyData, keySize)
			valuePtr := wasmBytePtr(ctx, mod, valueData, valueSize)
			return uint32(internal.ProxyAddHeaderMapValue(internal.MapType(mapType), keyPtr, int(keySize), valuePtr, int(valueSize)))
		}).
		Export("proxy_add_header_map_value").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, mapType uint32, keyData uint32, keySize uint32, valueData uint32, valueSize uint32) uint32 {
			keyPtr := wasmBytePtr(ctx, mod, keyData, keySize)
			valuePtr := wasmBytePtr(ctx, mod, valueData, valueSize)
			return uint32(internal.ProxyReplaceHeaderMapValue(internal.MapType(mapType), keyPtr, int(keySize), valuePtr, int(valueSize)))
		}).
		Export("proxy_replace_header_map_value").
		NewFunctionBuilder().
		WithFunc(func(streamType uint32) uint32 {
			return uint32(internal.ProxyContinueStream(internal.StreamType(streamType)))
		}).
		Export("proxy_continue_stream").
		NewFunctionBuilder().
		WithFunc(func(streamType uint32) uint32 {
			return uint32(internal.ProxyCloseStream(internal.StreamType(streamType)))
		}).
		Export("proxy_close_stream").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, mapType uint32, keyData uint32, keySize uint32) uint32 {
			keyPtr := wasmBytePtr(ctx, mod, keyData, keySize)
			return uint32(internal.ProxyRemoveHeaderMapValue(internal.MapType(mapType), keyPtr, int(keySize)))
		}).
		Export("proxy_remove_header_map_value").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, mapType uint32, returnValueData uint32, returnValueSize uint32) uint32 {
			var returnValueHostPtr *byte
			var returnValueSizePtr int
			ret := uint32(internal.ProxyGetHeaderMapPairs(internal.MapType(mapType), &returnValueHostPtr, &returnValueSizePtr))
			copyBytesToWasm(ctx, mod, returnValueHostPtr, returnValueSizePtr, returnValueData, returnValueSize)
			return ret
		}).
		Export("proxy_get_header_map_pairs").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, mapType uint32, mapData uint32, mapSize uint32) uint32 {
			mapPtr := wasmBytePtr(ctx, mod, mapData, mapSize)
			return uint32(internal.ProxySetHeaderMapPairs(internal.MapType(mapType), mapPtr, int(mapSize)))
		}).
		Export("proxy_set_header_map_pairs").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, bufferType uint32, start uint32, maxSize uint32, returnBufferData uint32, returnBufferSize uint32) uint32 {
			var returnBufferDataHostPtr *byte
			var returnBufferSizePtr int
			ret := uint32(internal.ProxyGetBufferBytes(internal.BufferType(bufferType), int(start), int(maxSize), &returnBufferDataHostPtr, &returnBufferSizePtr))
			copyBytesToWasm(ctx, mod, returnBufferDataHostPtr, returnBufferSizePtr, returnBufferData, returnBufferSize)
			return ret
		}).
		Export("proxy_get_buffer_bytes").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, bufferType uint32, start uint32, maxSize uint32, bufferData uint32, bufferSize uint32) uint32 {
			bufferPtr := wasmBytePtr(ctx, mod, bufferData, bufferSize)
			return uint32(internal.ProxySetBufferBytes(internal.BufferType(bufferType), int(start), int(maxSize), bufferPtr, int(bufferSize)))
		}).
		Export("proxy_set_buffer_bytes").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, upstreamData uint32, upstreamSize uint32, headerData uint32, headerSize uint32, bodyData uint32, bodySize uint32, trailersData uint32, trailersSize uint32, timeout uint32, calloutIDPtr uint32) uint32 {
			upstreamPtr := wasmBytePtr(ctx, mod, upstreamData, upstreamSize)
			headerPtr := wasmBytePtr(ctx, mod, headerData, headerSize)
			bodyPtr := wasmBytePtr(ctx, mod, bodyData, bodySize)
			trailersPtr := wasmBytePtr(ctx, mod, trailersData, trailersSize)
			var calloutID uint32
			ret := uint32(internal.ProxyHttpCall(upstreamPtr, int(upstreamSize), headerPtr, int(headerSize), bodyPtr, int(bodySize), trailersPtr, int(trailersSize), timeout, &calloutID))
			handleMemoryStatus(mod.Memory().WriteUint32Le(ctx, calloutIDPtr, calloutID))

			// Finishing proxy_http_call executes a callback, not a plugin lifecycle method, unlike every other host function which would then end up in wasm.
			// We can work around this by registering a callback here to go back to the wasm.
			internal.RegisterHttpCallout(calloutID, func(numHeaders, bodySize, numTrailers int) {
				proxyOnHttpCallResponse := mod.ExportedFunction("proxy_on_http_call_response")
				_, err := proxyOnHttpCallResponse.Call(ctx, uint64(getPluginContextID(ctx)), uint64(calloutID), uint64(numHeaders), uint64(bodySize), uint64(numTrailers))
				if err != nil {
					panic(err)
				}
			})

			return ret
		}).
		Export("proxy_http_call").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, funcNamePtr uint32, funcNameSize uint32, paramPtr uint32, paramSize uint32, returnData uint32, returnSize uint32) uint32 {
			funcName := wasmBytePtr(ctx, mod, funcNamePtr, funcNameSize)
			paramHostPtr := wasmBytePtr(ctx, mod, paramPtr, paramSize)
			var returnDataHostPtr *byte
			var returnDataSizePtr int
			ret := uint32(internal.ProxyCallForeignFunction(funcName, int(funcNameSize), paramHostPtr, int(paramSize), &returnDataHostPtr, &returnDataSizePtr))
			copyBytesToWasm(ctx, mod, returnDataHostPtr, returnDataSizePtr, returnData, returnSize)
			return ret
		}).
		Export("proxy_call_foreign_function").
		NewFunctionBuilder().
		WithFunc(func(period uint32) uint32 {
			return uint32(internal.ProxySetTickPeriodMilliseconds(period))
		}).
		Export("proxy_set_tick_period_milliseconds").
		NewFunctionBuilder().
		WithFunc(func(contextID uint32) uint32 {
			return uint32(internal.ProxySetEffectiveContext(contextID))
		}).
		Export("proxy_set_effective_context").
		NewFunctionBuilder().
		WithFunc(func() uint32 {
			return uint32(internal.ProxyDone())
		}).
		Export("proxy_done").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, metricType uint32, metricNameData uint32, metricNameSize uint32, returnMetricIDPtr uint32) uint32 {
			metricName := wasmBytePtr(ctx, mod, metricNameData, metricNameSize)
			var returnMetricID uint32
			ret := uint32(internal.ProxyDefineMetric(internal.MetricType(metricType), metricName, int(metricNameSize), &returnMetricID))
			handleMemoryStatus(mod.Memory().WriteUint32Le(ctx, returnMetricIDPtr, returnMetricID))
			return ret
		}).
		Export("proxy_define_metric").
		NewFunctionBuilder().
		WithFunc(func(metricID uint32, offset int64) uint32 {
			return uint32(internal.ProxyIncrementMetric(metricID, offset))
		}).
		Export("proxy_increment_metric").
		NewFunctionBuilder().
		WithFunc(func(metricID uint32, value uint64) uint32 {
			return uint32(internal.ProxyRecordMetric(metricID, value))
		}).
		Export("proxy_record_metric").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, metricID uint32, returnMetricValue uint32) uint32 {
			var returnMetricValuePtr uint64
			ret := uint32(internal.ProxyGetMetric(metricID, &returnMetricValuePtr))
			handleMemoryStatus(mod.Memory().WriteUint64Le(ctx, returnMetricValue, returnMetricValuePtr))
			return ret
		}).
		Export("proxy_get_metric").
		Instantiate(ctx, r)
	return err
}

type pluginContextIDKeyType struct{}

var pluginContextIDKey = pluginContextIDKeyType{}

func withPluginContextID(ctx context.Context, id uint32) context.Context {
	return context.WithValue(ctx, pluginContextIDKey, id)
}

func getPluginContextID(ctx context.Context) uint32 {
	id, _ := ctx.Value(pluginContextIDKey).(uint32)
	return id
}
