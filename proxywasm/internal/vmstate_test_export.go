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

// +build proxytest

package internal

import "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"

func VMStateReset() {
	// (@mathetake) I assume that the currentState be protected by lock on hostMux
	currentState = &state{
		pluginContexts:    make(map[uint32]*pluginContextState),
		httpContexts:      make(map[uint32]types.HttpContext),
		tcpContexts:       make(map[uint32]types.TcpContext),
		contextIDToRootID: make(map[uint32]uint32),
	}
}

func VMStateGetActiveContextID() uint32 {
	return currentState.activeContextID
}
