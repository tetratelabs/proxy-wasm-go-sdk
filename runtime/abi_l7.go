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

//export proxy_on_request_headers
func proxyOnRequestHeaders(contextID uint32, numHeaders int, endOfStream bool) Action {
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_headers")
	}

	currentState.setActiveContextID(contextID)
	return ctx.OnHttpRequestHeaders(numHeaders, endOfStream)
}

//export proxy_on_request_body
func proxyOnRequestBody(contextID uint32, bodySize int, endOfStream bool) Action {
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_body")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnHttpRequestBody(bodySize, endOfStream)
}

//export proxy_on_request_trailers
func proxyOnRequestTrailers(contextID uint32, numTrailers int) Action {
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_trailers")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnHttpRequestTrailers(numTrailers)
}

//export proxy_on_response_headers
func proxyOnHttpResponseHeaders(contextID uint32, numHeaders int, endOfStream bool) Action {
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnHttpResponseHeaders(numHeaders, endOfStream)
}

//export proxy_on_response_body
func proxyOnResponseBody(contextID uint32, bodySize int, endOfStream bool) Action {
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnHttpResponseBody(bodySize, endOfStream)
}

//export proxy_on_response_trailers
func proxyOnResponseTrailers(contextID uint32, numTrailers int) Action {
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnHttpResponseTrailers(numTrailers)
}

//export proxy_on_http_call_response
func proxyOnHttpCallResponse(_, calloutID uint32, numHeaders, bodySize, numTrailers int) {
	ctxID, ok := currentState.callOuts[calloutID]
	if !ok {
		panic("invalid callout id")
	}

	delete(currentState.callOuts, calloutID)

	if ctx, ok := currentState.streamContexts[ctxID]; ok {
		currentState.setActiveContextID(ctxID)
		setEffectiveContext(ctxID)
		ctx.OnHttpCallResponse(calloutID, numHeaders, bodySize, numTrailers)
	} else if ctx, ok := currentState.httpContexts[ctxID]; ok {
		currentState.setActiveContextID(ctxID)
		setEffectiveContext(ctxID)
		ctx.OnHttpCallResponse(calloutID, numHeaders, bodySize, numTrailers)
	} else if ctx, ok := currentState.rootContexts[ctxID]; ok {
		currentState.activeContextID = ctxID
		setEffectiveContext(ctxID)
		ctx.OnHttpCallResponse(calloutID, numHeaders, bodySize, numTrailers)
	} else {
		panic("invalid context on proxy_on_http_call_response")
	}
}
