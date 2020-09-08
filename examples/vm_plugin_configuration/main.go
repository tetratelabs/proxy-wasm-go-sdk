package main

import (
	"github.com/mathetake/proxy-wasm-go/runtime"
)

func main() {
	runtime.SetNewRootContext(func(uint32) runtime.RootContext { return &context{} })
	runtime.SetNewHttpContext(func(uint32) runtime.HttpContext { return &context{} })
}

type context struct{ runtime.DefaultContext }

// override
func (ctx *context) OnVMStart(vmConfigurationSize int) bool {
	data, err := runtime.HostCallGetVMConfiguration(vmConfigurationSize)
	if err != nil {
		runtime.LogCritical("error reading vm configuration", err.Error())
	}

	runtime.LogInfo("vm config: \n", string(data))
	return true
}

func (ctx *context) OnConfigure(pluginConfigurationSize int) bool {
	data, err := runtime.HostCallGetPluginConfiguration(pluginConfigurationSize)
	if err != nil {
		runtime.LogCritical("error reading plugin configuration", err.Error())
	}

	runtime.LogInfo("plugin config: \n", string(data))
	return true
}
