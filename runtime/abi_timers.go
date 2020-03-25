package runtime

//go:export proxy_on_tick
func ProxyOnTick(rootContextID uint32) {
	ctx, ok := currentState.rootContexts[rootContextID]
	if !ok {
		panic("invalid root_context_id")
	}
	currentState.setActiveContextID(rootContextID)
	ctx.OnTick()
}
