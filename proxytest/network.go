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

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type networkHostEmulator struct {
	streamStates map[uint32]*streamState
}

type streamState struct {
	upstream, downstream []byte
}

func newNetworkHostEmulator() *networkHostEmulator {
	host := &networkHostEmulator{
		streamStates: map[uint32]*streamState{},
	}

	return host
}

// impl rawhostcall.ProxyWASMHost: delegated from hostEmulator
func (n *networkHostEmulator) networkHostEmulatorProxyGetBufferBytes(bt types.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) types.Status {

	active := proxywasm.VMStateGetActiveContextID()
	stream := n.streamStates[active]
	var buf []byte
	switch bt {
	case types.BufferTypeUpstreamData:
		buf = stream.upstream
	case types.BufferTypeDownstreamData:
		buf = stream.downstream
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	if len(buf) == 0 {
		return types.StatusNotFound
	} else if start >= len(buf) {
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

// impl HostEmulator
func (n *networkHostEmulator) CallOnUpstreamData(contextID uint32, data []byte) {
	stream, ok := n.streamStates[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	if len(data) > 0 {
		stream.upstream = append(stream.upstream, data...)
	}

	action := proxywasm.ProxyOnUpstreamData(contextID, len(stream.upstream), false)
	switch action {
	case types.ActionPause:
		return
	case types.ActionContinue:
		stream.upstream = []byte{}
	default:
		log.Fatalf("invalid action type: %d", action)
	}
}

// impl HostEmulator
func (n *networkHostEmulator) CallOnDownstreamData(contextID uint32, data []byte) {
	stream, ok := n.streamStates[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}
	if len(data) > 0 {
		stream.downstream = append(stream.downstream, data...)
	}

	action := proxywasm.ProxyOnDownstreamData(contextID, len(stream.downstream), false)
	switch action {
	case types.ActionPause:
		return
	case types.ActionContinue:
		stream.downstream = []byte{}
	default:
		log.Fatalf("invalid action type: %d", action)
	}
}

// impl HostEmulator
func (n *networkHostEmulator) InitializeConnection() (contextID uint32) {
	contextID = getNextContextID()
	proxywasm.ProxyOnContextCreate(contextID, RootContextID)
	proxywasm.ProxyOnNewConnection(contextID)
	n.streamStates[contextID] = &streamState{}
	return
}

// impl HostEmulator
func (n *networkHostEmulator) CloseUpstreamConnection(contextID uint32) {
	proxywasm.ProxyOnUpstreamConnectionClose(contextID, types.PeerTypeLocal) // peerType will be removed in the next ABI
}

// impl HostEmulator
func (n *networkHostEmulator) CloseDownstreamConnection(contextID uint32) {
	proxywasm.ProxyOnDownstreamConnectionClose(contextID, types.PeerTypeLocal) // peerType will be removed in the next ABI
}

// impl HostEmulator
func (n *networkHostEmulator) CompleteConnection(contextID uint32) {
	// https://github.com/envoyproxy/envoy/blob/867b9e23d2e48350bd1b0d1fbc392a8355f20e35/source/extensions/common/wasm/context.cc#L169-L171
	proxywasm.ProxyOnDone(contextID)
	proxywasm.ProxyOnLog(contextID)
	proxywasm.ProxyOnDelete(contextID)
	delete(n.streamStates, contextID)
}
