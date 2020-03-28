package runtime

//go:export proxy_on_queue_ready
func ProxyOnQueueReady(contextID, queueID uint32) {
	ctx, ok := currentState.rootContexts[contextID]
	if !ok {
		panic("invalid context")
	}

	currentState.setActiveContextID(contextID)
	ctx.OnQueueReady(queueID)
}
