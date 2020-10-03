package proxytest

import (
	"sync"
	"time"

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
	HttpFilterPutResponseHeaders(contextID uint32, headers [][2]string)
	HttpFilterPutRequestTrailers(contextID uint32, headers [][2]string)
	HttpFilterPutResponseTrailers(contextID uint32, headers [][2]string)
	HttpFilterPutRequestBody(contextID uint32, body []byte)
	HttpFilterPutResponseBody(contextID uint32, body []byte)
	HttpFilterCompleteHttpStream(contextID uint32)
	HttpFilterGetCurrentStreamAction(contextID uint32) types.Action
	HttpFilterGetSentLocalResponse(contextID uint32) *LocalHttpResponse
}

const (
	rootContextID uint32 = 1 // TODO: support multiple rootContext
)

var (
	hostMux       = sync.Mutex{}
	nextContextID = rootContextID + 1
)

func NewHostEmulator(pluginConfiguration,
	vmConfiguration []byte,
	newRootContext func(uint32) proxywasm.RootContext,
	newStreamContext func(uint32) proxywasm.StreamContext,
	newHttpContext func(uint32) proxywasm.HttpContext,
) HostEmulator {
	root := newRootHostEmulator(pluginConfiguration, vmConfiguration)
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
	proxywasm.SetNewRootContext(newRootContext)
	proxywasm.SetNewStreamContext(newStreamContext)
	proxywasm.SetNewHttpContext(newHttpContext)

	// create root context: TODO: support multiple root contexts
	proxywasm.ProxyOnContextCreate(rootContextID, 0)

	return emulator
}

func getNextContextID() (ret uint32) {
	ret = nextContextID
	nextContextID++
	return
}

type hostEmulator struct {
	*rootHostEmulator
	*networkHostEmulator
	*httpHostEmulator

	effectiveContextID uint32
}

// impl host HostEmulator
func (*hostEmulator) Done() {
	hostMux.Unlock()
	proxywasm.VMStateReset()
}

// impl host rawhostcall.ProxyWASMHost
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

// impl host rawhostcall.ProxyWASMHost
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

// impl host rawhostcall.ProxyWASMHost
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

// impl host rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxyGetCurrentTimeNanoseconds(returnTime *int64) types.Status {
	*returnTime = time.Now().UnixNano()
	return types.StatusOK
}

// impl host rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxySetEffectiveContext(contextID uint32) types.Status {
	h.effectiveContextID = contextID
	return types.StatusOK
}

// impl host rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxySetProperty(*byte, int, *byte, int) types.Status {
	panic("unimplemented")
}

// impl host rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxyGetProperty(*byte, int, **byte, *int) types.Status {
	panic("unimplemented")
}

// impl host rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxyResolveSharedQueue(vmIDData *byte, vmIDSize int, nameData *byte, nameSize int, returnID *uint32) types.Status {
	panic("unimplemented")
}

// impl host rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxyCloseStream(streamType types.StreamType) types.Status {
	panic("unimplemented")
}

// impl host rawhostcall.ProxyWASMHost
func (h *hostEmulator) ProxyDone() types.Status {
	panic("unimplemented")
}
