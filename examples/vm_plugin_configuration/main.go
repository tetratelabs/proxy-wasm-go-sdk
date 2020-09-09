package main

import (
	"github.com/mathetake/proxy-wasm-go-sdk/proxywasm"
)

func main() {
	proxywasm.SetNewRootContext(func(uint32) proxywasm.RootContext { return &context{} })
}

type context struct{ proxywasm.DefaultContext }

// override
func (ctx *context) OnVMStart(vmConfigurationSize int) bool {
	data, err := proxywasm.HostCallGetVMConfiguration(vmConfigurationSize)
	if err != nil {
		proxywasm.LogCritical("error reading vm configuration", err.Error())
	}

	proxywasm.LogInfo("vm config: \n", string(data))
	return true
}

func (ctx *context) OnConfigure(pluginConfigurationSize int) bool {
	data, err := proxywasm.HostCallGetPluginConfiguration(pluginConfigurationSize)
	if err != nil {
		proxywasm.LogCritical("error reading plugin configuration", err.Error())
	}

	proxywasm.LogInfo("plugin config: \n", string(data))
	return true
}
