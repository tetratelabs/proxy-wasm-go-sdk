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

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
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

// impl internal.ProxyWasmHost: delegated from hostEmulator
func (n *networkHostEmulator) networkHostEmulatorProxyGetBufferBytes(bt internal.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) internal.Status {

	active := internal.VMStateGetActiveContextID()
	stream := n.streamStates[active]
	var buf []byte
	switch bt {
	case internal.BufferTypeUpstreamData:
		buf = stream.upstream
	case internal.BufferTypeDownstreamData:
		buf = stream.downstream
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	if len(buf) == 0 {
		return internal.StatusNotFound
	} else if start >= len(buf) {
		log.Printf("start index out of range: %d (start) >= %d ", start, len(buf))
		return internal.StatusBadArgument
	}

	*returnBufferData = &buf[start]
	if maxSize > len(buf)-start {
		*returnBufferSize = len(buf) - start
	} else {
		*returnBufferSize = maxSize
	}
	return internal.StatusOK
}

// impl HostEmulator
func (n *networkHostEmulator) CallOnUpstreamData(contextID uint32, data []byte) types.Action {
	stream, ok := n.streamStates[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}

	if len(data) > 0 {
		stream.upstream = append(stream.upstream, data...)
	}

	action := internal.ProxyOnUpstreamData(contextID, len(stream.upstream), false)
	switch action {
	case types.ActionPause:
	case types.ActionContinue:
		stream.upstream = []byte{}
	default:
		log.Fatalf("invalid action type: %d", action)
	}
	return action
}

// impl HostEmulator
func (n *networkHostEmulator) CallOnDownstreamData(contextID uint32, data []byte) types.Action {
	stream, ok := n.streamStates[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}
	if len(data) > 0 {
		stream.downstream = append(stream.downstream, data...)
	}

	action := internal.ProxyOnDownstreamData(contextID, len(stream.downstream), false)
	switch action {
	case types.ActionPause:
	case types.ActionContinue:
		stream.downstream = []byte{}
	default:
		log.Fatalf("invalid action type: %d", action)
	}
	return action
}

// impl HostEmulator
func (n *networkHostEmulator) InitializeConnection() (contextID uint32, action types.Action) {
	contextID = getNextContextID()
	internal.ProxyOnContextCreate(contextID, PluginContextID)
	action = internal.ProxyOnNewConnection(contextID)
	n.streamStates[contextID] = &streamState{}
	return
}

// impl HostEmulator
func (n *networkHostEmulator) CloseUpstreamConnection(contextID uint32) {
	internal.ProxyOnUpstreamConnectionClose(contextID, types.PeerTypeLocal) // peerType will be removed in the next ABI
}

// impl HostEmulator
func (n *networkHostEmulator) CloseDownstreamConnection(contextID uint32) {
	internal.ProxyOnDownstreamConnectionClose(contextID, types.PeerTypeLocal) // peerType will be removed in the next ABI
}

// impl HostEmulator
func (n *networkHostEmulator) CompleteConnection(contextID uint32) {
	internal.ProxyOnLog(contextID)
	internal.ProxyOnDelete(contextID)
	delete(n.streamStates, contextID)
}
