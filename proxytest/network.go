package proxytest

import (
	"log"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/rawhostcall"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type NetworkFilterHost struct {
	baseHost
	newContext       func(contextID uint32) proxywasm.StreamContext
	streams          map[uint32]*streamState
	currentContextID uint32
}

type streamState struct {
	upstream, downstream []byte
	ctx                  proxywasm.StreamContext
}

func NewNetworkFilterHost(f func(contextID uint32) proxywasm.StreamContext) (*NetworkFilterHost, func()) {
	hostMux.Lock() // acquire the lock of host emulation
	host := &NetworkFilterHost{
		newContext: f,
		streams:    map[uint32]*streamState{},
	}
	rawhostcall.RegisterMockWASMHost(host)
	return host, func() {
		hostMux.Unlock()
	}
}

func (n *NetworkFilterHost) PutUpstreamData(contextID uint32, data []byte) {
	stream, ok := n.streams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	if len(data) > 0 {
		stream.upstream = append(stream.upstream, data...)
	}

	n.currentContextID = contextID
	action := stream.ctx.OnUpstreamData(len(stream.upstream), false)
	switch action {
	case types.ActionPause:
		return
	case types.ActionContinue:
		// TODO: verify the behavior is correct
		stream.upstream = []byte{}
	default:
		log.Fatalf("invalid action type: %d", action)
	}
}

func (n *NetworkFilterHost) PutDownstreamData(contextID uint32, data []byte) {
	stream, ok := n.streams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}
	if len(data) > 0 {
		stream.downstream = append(stream.downstream, data...)
	}

	n.currentContextID = contextID
	action := stream.ctx.OnDownstreamData(len(stream.downstream), false)
	switch action {
	case types.ActionPause:
		return
	case types.ActionContinue:
		// TODO: verify the behavior is correct
		stream.downstream = []byte{}
	default:
		log.Fatalf("invalid action type: %d", action)
	}
}

func (n *NetworkFilterHost) InitConnection() (contextID uint32) {
	contextID = uint32(len(n.streams) + 1)
	ctx := n.newContext(contextID)
	n.streams[contextID] = &streamState{ctx: ctx}

	n.currentContextID = contextID
	ctx.OnNewConnection()
	return
}

func (n *NetworkFilterHost) CloseUpstreamConnection(contextID uint32) {
	n.streams[contextID].ctx.OnUpstreamClose(types.PeerTypeLocal) // peerType will be removed in the next ABI
}

func (n *NetworkFilterHost) CloseDownstreamConnection(contextID uint32) {
	n.streams[contextID].ctx.OnDownstreamClose(types.PeerTypeLocal) // peerType will be removed in the next ABI
}

func (n *NetworkFilterHost) ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) types.Status {
	stream := n.streams[n.currentContextID]
	var buf []byte
	switch bt {
	case types.BufferTypeUpstreamData:
		buf = stream.upstream
	case types.BufferTypeDownstreamData:
		buf = stream.downstream
	default:
		// delegate base host implementation
		return n.getBuffer(bt, start, maxSize, returnBufferData, returnBufferSize)
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
