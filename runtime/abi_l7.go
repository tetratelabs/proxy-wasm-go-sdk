package runtime

//go:export proxy_on_request_headers
func proxyOnRequestHeaders(contextID uint32, numHeaders int) Action {
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_headers")
	}

	currentState.setActiveContextID(contextID)
	return ctx.OnHttpRequestHeaders(numHeaders)
}

//go:export proxy_on_request_body
func proxyOnRequestBody(contextID uint32, bodySize int, endOfStream bool) Action {
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_body")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnHttpRequestBody(bodySize, endOfStream)
}

//go:export proxy_on_request_trailers
func proxyOnRequestTrailers(contextID uint32, numTrailers int) Action {
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_trailers")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnHttpRequestTrailers(numTrailers)
}

//go:export proxy_on_response_headers
func proxyOnHttpResponseHeaders(contextID uint32, numHeaders int) Action {
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnHttpResponseHeaders(numHeaders)
}

//go:export proxy_on_response_body
func proxyOnResponseBody(contextID uint32, bodySize int, endOfStream bool) Action {
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnHttpResponseBody(bodySize, endOfStream)
}

//go:export proxy_on_response_trailers
func proxyOnResponseTrailers(contextID uint32, numTrailers int) Action {
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnHttpResponseTrailers(numTrailers)
}

//go:export proxy_on_http_call_response
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

//go:export proxy_on_queue_ready
func ProxyOnQueueReady(contextID uint32, queueID uint32) {
	return
}
