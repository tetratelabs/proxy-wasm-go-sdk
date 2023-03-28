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
	"os"
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

	compiled, err := r.CompileModule(ctx, wasm)
	if err != nil {
		return nil, err
	}

	wazeroconfig := wazero.NewModuleConfig().WithStdout(os.Stderr).WithStderr(os.Stderr)
	mod, err := r.InstantiateModule(ctx, compiled, wazeroconfig)
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

func wasmBytePtr(mod api.Module, off uint32, size uint32) *byte {
	if size == 0 {
		return nil

	}
	buf, ok := mod.Memory().Read(off, size)
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
	buf, ok := mod.Memory().Read(uint32(res[0]), uint32(size))
	handleMemoryStatus(ok)

	copy(buf, hostSlice)
	ok = mod.Memory().WriteUint32Le(wasmPtrPtr, uint32(res[0]))
	handleMemoryStatus(ok)

	ok = mod.Memory().WriteUint32Le(wasmSizePtr, uint32(size))
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
		// proxy_log logs a message at the given log_level.
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_log
		NewFunctionBuilder().
		WithParameterNames("log_level", "message_data", "message_size").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, logLevel, messageData, messageSize uint32) uint32 {
			messageDataPtr := wasmBytePtr(mod, messageData, messageSize)
			return uint32(internal.ProxyLog(internal.LogLevel(logLevel), messageDataPtr, int(messageSize)))
		}).
		Export("proxy_log").
		// proxy_set_property sets a property value.
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_set_property
		NewFunctionBuilder().
		WithParameterNames("property_path_data", "property_path_size", "property_value_data",
			"property_value_size").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, pathData, pathSize, valueData, valueSize uint32) uint32 {
			pathDataPtr := wasmBytePtr(mod, pathData, pathSize)
			valueDataPtr := wasmBytePtr(mod, valueData, valueSize)
			return uint32(internal.ProxySetProperty(pathDataPtr, int(pathSize), valueDataPtr, int(valueSize)))
		}).
		Export("proxy_set_property").
		// proxy_get_property gets a property value.
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_get_property
		NewFunctionBuilder().
		WithParameterNames("property_path_data", "property_path_size", "return_property_value_data",
			"return_property_value_size").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, pathData, pathSize, returnValueData,
			returnValueSize uint32) uint32 {
			pathDataPtr := wasmBytePtr(mod, pathData, pathSize)
			var returnValueHostPtr *byte
			var returnValueSizePtr int
			ret := uint32(internal.ProxyGetProperty(pathDataPtr, int(pathSize), &returnValueHostPtr, &returnValueSizePtr))
			copyBytesToWasm(ctx, mod, returnValueHostPtr, returnValueSizePtr, returnValueData, returnValueSize)
			return ret
		}).
		Export("proxy_get_property").
		// proxy_send_local_response sends an HTTP response without forwarding request to the upstream.
		//
		// Note: proxy-wasm spec calls this proxy_send_http_response. See
		// https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_send_http_response
		NewFunctionBuilder().
		WithParameterNames("response_code", "response_code_details_data", "response_code_details_size",
			"response_body_data", "response_body_size", "additional_headers_map_data", "additional_headers_size",
			"grpc_status").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, statusCode, statusCodeDetailData, statusCodeDetailsSize,
			bodyData, bodySize, headersData, headersSize, grpcStatus uint32) uint32 {
			statusCodeDetailDataPtr := wasmBytePtr(mod, statusCodeDetailData, statusCodeDetailsSize)
			bodyDataPtr := wasmBytePtr(mod, bodyData, bodySize)
			headersDataPtr := wasmBytePtr(mod, headersData, headersSize)
			return uint32(internal.ProxySendLocalResponse(statusCode, statusCodeDetailDataPtr,
				int(statusCodeDetailsSize), bodyDataPtr, int(bodySize), headersDataPtr, int(headersSize), int32(grpcStatus)))
		}).
		Export("proxy_send_local_response").
		// proxy_get_shared_data gets shared data identified by a key. The compare-and-switch value is returned and can
		// be used when updating the value with proxy_set_shared_data.
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_get_shared_data
		NewFunctionBuilder().
		WithParameterNames("key_data", "key_size", "return_value_data", "return_value_size", "return_cas").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, keyData, keySize, returnValueData, returnValueSize,
			returnCas uint32) uint32 {
			keyDataPtr := wasmBytePtr(mod, keyData, keySize)
			var returnValueHostPtr *byte
			var returnValueSizePtr int
			var returnCasPtr uint32
			ret := uint32(internal.ProxyGetSharedData(keyDataPtr, int(keySize), &returnValueHostPtr,
				&returnValueSizePtr, &returnCasPtr))
			copyBytesToWasm(ctx, mod, returnValueHostPtr, returnValueSizePtr, returnValueData, returnValueSize)
			handleMemoryStatus(mod.Memory().WriteUint32Le(returnCas, returnCasPtr))
			return ret
		}).
		Export("proxy_get_shared_data").
		// proxy_set_shared_data sets the value of shared data using its key. If compare-and-switch value is set, it
		// must match the current value.
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_set_shared_data
		NewFunctionBuilder().
		WithParameterNames("key_data", "key_size", "value_data", "value_size", "cas").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, keyData, keySize, valueData, valueSize, cas uint32) uint32 {
			keyDataPtr := wasmBytePtr(mod, keyData, keySize)
			valueDataPtr := wasmBytePtr(mod, valueData, valueSize)
			return uint32(internal.ProxySetSharedData(keyDataPtr, int(keySize), valueDataPtr, int(valueSize), cas))
		}).
		Export("proxy_set_shared_data").
		// proxy_register_shared_queue registers a shared queue using a given name. It can be referred to in
		// proxy_enqueue_shared_queue and proxy_dequeue_shared_queue using the returned ID.
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_register_shared_queue
		NewFunctionBuilder().
		WithParameterNames("queue_name_data", "queue_name_size", "return_queue_id").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, nameData, nameSize, returnID uint32) uint32 {
			namePtr := wasmBytePtr(mod, nameData, nameSize)
			var returnIDPtr uint32
			ret := uint32(internal.ProxyRegisterSharedQueue(namePtr, int(nameSize), &returnIDPtr))
			handleMemoryStatus(mod.Memory().WriteUint32Le(returnID, returnIDPtr))
			return ret
		}).
		Export("proxy_register_shared_queue").
		// proxy_resolve_shared_queue resolves existing shared queue using a given name. It can be referred to in
		// proxy_enqueue_shared_queue and proxy_dequeue_shared_queue using the returned ID.
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_resolve_shared_queue
		//
		// Note: The "vm_id_data" and "vm_id_size" parameters are not documented in proxy-wasm spec.
		NewFunctionBuilder().
		WithParameterNames("vm_id_data", "vm_id_size", "queue_name_data", "queue_name_size", "return_queue_id").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, vmIDData, vmIDSize, nameData, nameSize, returnID uint32) uint32 {
			vmID := wasmBytePtr(mod, vmIDData, vmIDSize)
			namePtr := wasmBytePtr(mod, nameData, nameSize)
			var returnIDPtr uint32
			ret := uint32(internal.ProxyResolveSharedQueue(vmID, int(vmIDSize), namePtr, int(nameSize), &returnIDPtr))
			handleMemoryStatus(mod.Memory().WriteUint32Le(returnID, returnIDPtr))
			return ret
		}).
		Export("proxy_resolve_shared_queue").
		// proxy_dequeue_shared_queue gets data from the end of the queue.
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_dequeue_shared_queue
		NewFunctionBuilder().
		WithParameterNames("queue_id", "payload_data", "payload_size").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, queueID, returnValueData, returnValueSize uint32) uint32 {
			var returnValueHostPtr *byte
			var returnValueSizePtr int
			ret := uint32(internal.ProxyDequeueSharedQueue(queueID, &returnValueHostPtr, &returnValueSizePtr))
			copyBytesToWasm(ctx, mod, returnValueHostPtr, returnValueSizePtr, returnValueData, returnValueSize)
			return ret
		}).
		Export("proxy_dequeue_shared_queue").
		// proxy_enqueue_shared_queue adds data to the front of the queue.
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_enqueue_shared_queue
		NewFunctionBuilder().
		WithParameterNames("queue_id", "payload_data", "payload_size").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, queueID, valueData, valueSize uint32) uint32 {
			valuePtr := wasmBytePtr(mod, valueData, valueSize)
			return uint32(internal.ProxyEnqueueSharedQueue(queueID, valuePtr, int(valueSize)))
		}).
		Export("proxy_enqueue_shared_queue").
		// proxy_get_header_map_value gets the content of key from a given map.
		//
		// Note: proxy-wasm-spec calls this proxy_get_map_value. See
		// https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_get_map_value
		NewFunctionBuilder().
		WithParameterNames("map_type", "key_data", "key_size", "return_value_data", "return_value_size").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, mapType, keyData, keySize, returnValueData,
			returnValueSize uint32) uint32 {
			keyPtr := wasmBytePtr(mod, keyData, keySize)
			var retValDataHostPtr *byte
			var retValSizePtr int
			ret := uint32(internal.ProxyGetHeaderMapValue(internal.MapType(mapType), keyPtr, int(keySize), &retValDataHostPtr, &retValSizePtr))
			copyBytesToWasm(ctx, mod, retValDataHostPtr, retValSizePtr, returnValueData, returnValueSize)
			return ret
		}).
		Export("proxy_get_header_map_value").
		// proxy_add_header_map_value adds a value to the key of a given map.
		//
		// Note: proxy-wasm-spec calls this proxy_add_map_value. See
		// https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_add_map_value
		NewFunctionBuilder().
		WithParameterNames("map_type", "key_data", "key_size", "value_data", "value_size").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, mapType, keyData, keySize, valueData, valueSize uint32) uint32 {
			keyPtr := wasmBytePtr(mod, keyData, keySize)
			valuePtr := wasmBytePtr(mod, valueData, valueSize)
			return uint32(internal.ProxyAddHeaderMapValue(internal.MapType(mapType), keyPtr, int(keySize), valuePtr, int(valueSize)))
		}).
		Export("proxy_add_header_map_value").
		// proxy_replace_header_map_value replaces any value of the key in a given map.
		//
		// Note: proxy-wasm-spec calls this proxy_set_map_value. See
		// https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_set_map_value
		NewFunctionBuilder().
		WithParameterNames("map_type", "key_data", "key_size", "value_data", "value_size").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, mapType, keyData, keySize, valueData, valueSize uint32) uint32 {
			keyPtr := wasmBytePtr(mod, keyData, keySize)
			valuePtr := wasmBytePtr(mod, valueData, valueSize)
			return uint32(internal.ProxyReplaceHeaderMapValue(internal.MapType(mapType), keyPtr, int(keySize), valuePtr, int(valueSize)))
		}).
		Export("proxy_replace_header_map_value").
		// proxy_continue_stream resume processing of paused stream.
		//
		// Note: This is similar to proxy_resume_downstream, proxy_resume_upstream, proxy_resume_http_request and
		// proxy_resume_http_response in proxy-wasm spec. See
		// https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_continue_stream
		NewFunctionBuilder().
		WithParameterNames("stream_type").
		WithResultNames("call_result").
		WithFunc(func(streamType uint32) uint32 {
			return uint32(internal.ProxyContinueStream(internal.StreamType(streamType)))
		}).
		Export("proxy_continue_stream").
		// proxy_close_stream closes a stream.
		//
		// Note: This is undocumented in proxy-wasm spec.
		NewFunctionBuilder().
		WithParameterNames("stream_type").
		WithResultNames("call_result").
		WithFunc(func(streamType uint32) uint32 {
			return uint32(internal.ProxyCloseStream(internal.StreamType(streamType)))
		}).
		Export("proxy_close_stream").
		// proxy_remove_header_map_value removes all values of the key in a given map.
		//
		// Note: proxy-wasm-spec calls this proxy_remove_map_value. See
		// https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_remove_map_value
		NewFunctionBuilder().
		WithParameterNames("map_type", "key_data", "key_size").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, mapType, keyData, keySize uint32) uint32 {
			keyPtr := wasmBytePtr(mod, keyData, keySize)
			return uint32(internal.ProxyRemoveHeaderMapValue(internal.MapType(mapType), keyPtr, int(keySize)))
		}).
		Export("proxy_remove_header_map_value").
		// proxy_get_header_map_pairs gets all key-value pairs from a given map.
		//
		// Note: proxy-wasm-spec calls this proxy_get_map. See
		// https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_get_map
		NewFunctionBuilder().
		WithParameterNames("map_type", "return_map_data", "return_map_size").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, mapType, returnValueData, returnValueSize uint32) uint32 {
			var returnValueHostPtr *byte
			var returnValueSizePtr int
			ret := uint32(internal.ProxyGetHeaderMapPairs(internal.MapType(mapType), &returnValueHostPtr, &returnValueSizePtr))
			copyBytesToWasm(ctx, mod, returnValueHostPtr, returnValueSizePtr, returnValueData, returnValueSize)
			return ret
		}).
		Export("proxy_get_header_map_pairs").
		// proxy_set_header_map_pairs gets all key-value pairs from a given map.
		//
		// Note: proxy-wasm-spec calls this proxy_set_map. See
		// https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_set_map
		NewFunctionBuilder().
		WithParameterNames("map_type", "map_data", "map_size").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, mapType, mapData, mapSize uint32) uint32 {
			mapPtr := wasmBytePtr(mod, mapData, mapSize)
			return uint32(internal.ProxySetHeaderMapPairs(internal.MapType(mapType), mapPtr, int(mapSize)))
		}).
		Export("proxy_set_header_map_pairs").
		// proxy_get_buffer_bytes gets up to max_size bytes from the buffer, starting from offset.
		//
		// Note: proxy-wasm-spec calls this proxy_get_buffer, but the signature is incompatible. See
		// https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_get_buffer
		NewFunctionBuilder().
		WithParameterNames("buffer_type", "offset", "max_size", "return_buffer_data", "return_buffer_size").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, bufferType, start, maxSize, returnBufferData,
			returnBufferSize uint32) uint32 {
			var returnBufferDataHostPtr *byte
			var returnBufferSizePtr int
			ret := uint32(internal.ProxyGetBufferBytes(internal.BufferType(bufferType), int(start), int(maxSize), &returnBufferDataHostPtr, &returnBufferSizePtr))
			copyBytesToWasm(ctx, mod, returnBufferDataHostPtr, returnBufferSizePtr, returnBufferData, returnBufferSize)
			return ret
		}).
		Export("proxy_get_buffer_bytes").
		// proxy_set_buffer_bytes replaces a byte range of the given buffer type.
		//
		// Note: proxy-wasm-spec calls this proxy_set_buffer, but the signature is incompatible. See
		// https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_set_buffer
		NewFunctionBuilder().
		WithParameterNames("buffer_type", "offset", "size", "buffer_data", "buffer_size").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, bufferType, start, maxSize, bufferData,
			bufferSize uint32) uint32 {
			bufferPtr := wasmBytePtr(mod, bufferData, bufferSize)
			return uint32(internal.ProxySetBufferBytes(internal.BufferType(bufferType), int(start), int(maxSize), bufferPtr, int(bufferSize)))
		}).
		Export("proxy_set_buffer_bytes").
		// proxy_http_call dispatches an HTTP call to upstream. Once the response is returned to the host,
		// proxy_on_http_call_response will be called with a unique call identifier (return_callout_id).
		//
		// Note: proxy-wasm-spec calls this proxy_dispatch_http_call. See
		// https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_dispatch_http_call
		NewFunctionBuilder().
		WithParameterNames("upstream_name_data", "upstream_name_size", "headers_map_data", "headers_map_size",
			"body_data", "body_size", "trailers_map_data", "trailers_map_size", "timeout_milliseconds",
			"return_callout_id").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, upstreamData, upstreamSize, headerData, headerSize, bodyData,
			bodySize, trailersData, trailersSize, timeout, calloutIDPtr uint32) uint32 {
			upstreamPtr := wasmBytePtr(mod, upstreamData, upstreamSize)
			headerPtr := wasmBytePtr(mod, headerData, headerSize)
			bodyPtr := wasmBytePtr(mod, bodyData, bodySize)
			trailersPtr := wasmBytePtr(mod, trailersData, trailersSize)
			var calloutID uint32
			ret := uint32(internal.ProxyHttpCall(upstreamPtr, int(upstreamSize), headerPtr, int(headerSize), bodyPtr, int(bodySize), trailersPtr, int(trailersSize), timeout, &calloutID))
			handleMemoryStatus(mod.Memory().WriteUint32Le(calloutIDPtr, calloutID))

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
		// proxy_call_foreign_function calls a registered foreign function.
		//
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_call_foreign_function
		NewFunctionBuilder().
		WithParameterNames("function_name_data", "function_name_size", "parameters_data", "parameters_size",
			"return_results_data", "return_results_size").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, funcNamePtr, funcNameSize, paramPtr, paramSize, returnData,
			returnSize uint32) uint32 {
			funcName := wasmBytePtr(mod, funcNamePtr, funcNameSize)
			paramHostPtr := wasmBytePtr(mod, paramPtr, paramSize)
			var returnDataHostPtr *byte
			var returnDataSizePtr int
			ret := uint32(internal.ProxyCallForeignFunction(funcName, int(funcNameSize), paramHostPtr, int(paramSize), &returnDataHostPtr, &returnDataSizePtr))
			copyBytesToWasm(ctx, mod, returnDataHostPtr, returnDataSizePtr, returnData, returnSize)
			return ret
		}).
		Export("proxy_call_foreign_function").
		// proxy_set_tick_period_milliseconds sets the timer period. Once set, the host environment will call
		// proxy_on_tick every tick_period milliseconds.
		//
		// Note: proxy-wasm spec calls this proxy_set_tick_period. See
		// https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_set_tick_period
		NewFunctionBuilder().
		WithParameterNames("tick_period").
		WithResultNames("call_result").
		WithFunc(func(period uint32) uint32 {
			return uint32(internal.ProxySetTickPeriodMilliseconds(period))
		}).
		Export("proxy_set_tick_period_milliseconds").
		// proxy_set_effective_context changes the effective context. This function is usually used to change the
		// context after receiving proxy_on_http_call_response, proxy_on_grpc_call_response or proxy_on_queue_ready.
		//
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_set_effective_context
		NewFunctionBuilder().
		WithParameterNames("context_id").
		WithResultNames("call_result").
		WithFunc(func(contextID uint32) uint32 {
			return uint32(internal.ProxySetEffectiveContext(contextID))
		}).
		Export("proxy_set_effective_context").
		// proxy_done indicates to the host environment that Wasm VM side is done processing current context. This can
		// be used after returning false in proxy_on_done.
		//
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_done
		NewFunctionBuilder().
		WithResultNames("call_result").
		WithFunc(func() uint32 {
			return uint32(internal.ProxyDone())
		}).
		Export("proxy_done").
		// proxy_define_metric defines a metric using a given name. It can be referred to in proxy_get_metric,
		// proxy_increment_metric and proxy_record_metric using returned ID.
		//
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_define_metric
		NewFunctionBuilder().
		WithParameterNames("metric_type", "metric_name_data", "metric_name_size", "return_metric_id").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, metricType, metricNameData, metricNameSize,
			returnMetricIDPtr uint32) uint32 {
			metricName := wasmBytePtr(mod, metricNameData, metricNameSize)
			var returnMetricID uint32
			ret := uint32(internal.ProxyDefineMetric(internal.MetricType(metricType), metricName, int(metricNameSize), &returnMetricID))
			handleMemoryStatus(mod.Memory().WriteUint32Le(returnMetricIDPtr, returnMetricID))
			return ret
		}).
		Export("proxy_define_metric").
		// proxy_increment_metric increments or decrements a metric value using an offset.
		//
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_increment_metric
		NewFunctionBuilder().
		WithParameterNames("metric_id", "offset").
		WithResultNames("call_result").
		WithFunc(func(metricID uint32, offset int64) uint32 {
			return uint32(internal.ProxyIncrementMetric(metricID, offset))
		}).
		Export("proxy_increment_metric").
		// proxy_record_metric sets the value of a metric.
		//
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_record_metric
		NewFunctionBuilder().
		WithParameterNames("metric_id", "value").
		WithResultNames("call_result").
		WithFunc(func(metricID uint32, value uint64) uint32 {
			return uint32(internal.ProxyRecordMetric(metricID, value))
		}).
		Export("proxy_record_metric").
		// proxy_get_metric gets the value of a metric.
		//
		// See https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT#proxy_get_metric
		NewFunctionBuilder().
		WithParameterNames("metric_id", "return_value").
		WithResultNames("call_result").
		WithFunc(func(ctx context.Context, mod api.Module, metricID, returnMetricValue uint32) uint32 {
			var returnMetricValuePtr uint64
			ret := uint32(internal.ProxyGetMetric(metricID, &returnMetricValuePtr))
			handleMemoryStatus(mod.Memory().WriteUint64Le(returnMetricValue, returnMetricValuePtr))
			return ret
		}).
		Export("proxy_get_metric").
		Instantiate(ctx)
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
