package runtime

//go:export proxy_on_new_connection
func proxyOnNewConnection(contextID uint32) Action {
	ctx, ok := currentState.streamContexts[contextID]
	if !ok {
		panic("invalid context")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnNewConnection()
}

//go:export proxy_on_downstream_data
func proxyOnDownstreamData(contextID uint32, dataSize int, endOfStream bool) Action {
	ctx, ok := currentState.streamContexts[contextID]
	if !ok {
		panic("invalid context")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnDownstreamData(dataSize, endOfStream)
}

//go:export on_downstream_connection_close
func proxyOnDownstreamConnectionClose(contextID uint32, pType PeerType) {
	ctx, ok := currentState.streamContexts[contextID]
	if !ok {
		panic("invalid context")
	}
	currentState.setActiveContextID(contextID)
	ctx.OnDownStreamClose(pType)
}

//go:export proxy_on_upstream_data
func proxyOnUpstreamData(contextID uint32, dataSize int, endOfStream bool) Action {
	ctx, ok := currentState.streamContexts[contextID]
	if !ok {
		panic("invalid context")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnUpstreamData(dataSize, endOfStream)
}

//go:export proxy_on_upstream_connection_close
func proxyOnUpstreamConnectionClose(contextID uint32, pType PeerType) {
	ctx, ok := currentState.streamContexts[contextID]
	if !ok {
		panic("invalid context")
	}
	currentState.setActiveContextID(contextID)
	ctx.OnUpstreamStreamClose(pType)
}
