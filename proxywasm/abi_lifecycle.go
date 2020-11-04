// Copyright 2020 Tetrate
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

package proxywasm

//export proxy_on_context_create
func proxyOnContextCreate(contextID uint32, rootContextID uint32) {
	if rootContextID == 0 {
		currentState.createRootContext(contextID)
	} else if currentState.newHttpContext != nil {
		currentState.createHttpContext(contextID, rootContextID)
	} else if currentState.newStreamContext != nil {
		currentState.createStreamContext(contextID, rootContextID)
	} else {
		panic("invalid context id on proxy_on_context_create")
	}
}

//export proxy_on_done
func proxyOnDone(contextID uint32) bool {
	defer func() {
		delete(currentState.contextIDToRootID, contextID)
	}()
	if ctx, ok := currentState.streams[contextID]; ok {
		currentState.setActiveContextID(contextID)
		delete(currentState.streams, contextID)
		ctx.OnStreamDone()
		return true
	} else if ctx, ok := currentState.httpStreams[contextID]; ok {
		currentState.setActiveContextID(contextID)
		ctx.OnHttpStreamDone()
		delete(currentState.httpStreams, contextID)
		return true
	} else if ctx, ok := currentState.rootContexts[contextID]; ok {
		currentState.setActiveContextID(contextID)
		response := ctx.context.OnVMDone()
		delete(currentState.rootContexts, contextID)
		return response
	} else {
		panic("invalid context on proxy_on_done")
	}
}
