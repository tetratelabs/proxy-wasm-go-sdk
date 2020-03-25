package runtime

//go:export proxy_on_new_connection
func proxyOnNewConnection(contextID uint32) Action {
	// TODO:
	return ActionContinue
}

//go:export proxy_on_downstream_data
func proxyOnDownstreamData(contextID uint32, dataSize uint, endOfStream bool) Action {
	// TODO:
	return ActionContinue
}

//go:export on_downstream_connection_close
func proxyOnDownstreamConnectionClose(contextID uint32, pType PeerType) {
	// TODO:
}

//go:export proxy_on_upstream_data
func proxyOnUpstreamData(contextID uint, dataSize uint, endOfStream bool) Action {
	// TODO:
	return ActionContinue
}

//go:export proxy_on_upstream_connection_close
func proxyOnUpstreamConnectionClose(contextID uint32, pType PeerType) {}
