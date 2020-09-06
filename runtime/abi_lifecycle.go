// Copyright 2020 Takeshi Yoneda(@mathetake)
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

package runtime

//export proxy_on_context_create
func proxyOnContextCreate(contextID uint32, rootContextID uint32) {
	if rootContextID == 0 {
		currentState.createRootContext(contextID)
	} else if currentState.newHttpContext != nil {
		currentState.createHttpContext(contextID, rootContextID)
	} else if currentState.newStreamContext != nil {
		currentState.createStreamContext(contextID, rootContextID)
	} else {
		panic("proxy_on_context_create failed")
	}
}

//export proxy_on_done
func proxyOnDone(contextID uint32) bool {
	if ctx, ok := currentState.streamContexts[contextID]; ok {
		currentState.setActiveContextID(contextID)
		return ctx.OnDone()
	} else if ctx, ok := currentState.httpContexts[contextID]; ok {
		currentState.setActiveContextID(contextID)
		return ctx.OnDone()
	} else if ctx, ok := currentState.rootContexts[contextID]; ok {
		currentState.setActiveContextID(contextID)
		return ctx.OnDone()
	} else {
		panic("invalid context on proxy_on_done")
	}
}

//export proxy_on_log
func proxyOnLog(contextID uint32) {
	if ctx, ok := currentState.streamContexts[contextID]; ok {
		currentState.setActiveContextID(contextID)
		ctx.OnLog()
	} else if ctx, ok := currentState.httpContexts[contextID]; ok {
		currentState.setActiveContextID(contextID)
		ctx.OnLog()
	} else if ctx, ok := currentState.rootContexts[contextID]; ok {
		currentState.setActiveContextID(contextID)
		ctx.OnLog()
	} else {
		panic("invalid context on proxy_on_log")
	}
}

//export proxy_on_delete
func proxyOnDelete(contextID uint32) {
	if _, ok := currentState.streamContexts[contextID]; ok {
		delete(currentState.streamContexts, contextID)
	} else if _, ok := currentState.httpContexts[contextID]; ok {
		delete(currentState.httpContexts, contextID)
	} else if _, ok := currentState.rootContexts[contextID]; ok {
		delete(currentState.rootContexts, contextID)
	} else {
		panic("invalid context on proxy_on_delete")
	}
}
