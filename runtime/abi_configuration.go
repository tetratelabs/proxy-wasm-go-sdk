package runtime

//go:export proxy_on_vm_start
func proxyOnVMStart(rootContextID uint32, vmConfigurationSize int) bool {
	ctx, ok := currentState.rootContexts[rootContextID]
	if !ok {
		panic("invalid context on proxy_on_vm_start")
	}
	currentState.setActiveContextID(rootContextID)
	return ctx.OnVMStart(vmConfigurationSize)
}

//go:export proxy_on_configure
func proxyOnConfigure(rootContextID uint32, pluginConfigurationSize int) bool {
	ctx, ok := currentState.rootContexts[rootContextID]
	if !ok {
		panic("invalid context on proxy_on_configure")
	}
	currentState.setActiveContextID(rootContextID)
	return ctx.OnConfigure(pluginConfigurationSize)
}
