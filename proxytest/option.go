package proxytest

import "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"

type EmulatorOption struct {
	pluginConfiguration, vmConfiguration []byte
	newRootContext                       func(uint32) proxywasm.RootContext
	newStreamContext                     func(uint32) proxywasm.StreamContext
	newHttpContext                       func(uint32) proxywasm.HttpContext
}

func NewEmulatorOption() *EmulatorOption {
	return &EmulatorOption{}
}

func (o *EmulatorOption) WithNewRootContext(f func(uint32) proxywasm.RootContext) *EmulatorOption {
	o.newRootContext = f
	return o
}

func (o *EmulatorOption) WithNewHttpContext(f func(uint32) proxywasm.HttpContext) *EmulatorOption {
	o.newHttpContext = f
	return o
}

func (o *EmulatorOption) WithNewStreamContext(f func(uint32) proxywasm.StreamContext) *EmulatorOption {
	o.newStreamContext = f
	return o
}

func (o *EmulatorOption) WithPluginConfiguration(data []byte) *EmulatorOption {
	o.pluginConfiguration = data
	return o
}

func (o *EmulatorOption) WithVMConfiguration(data []byte) *EmulatorOption {
	o.vmConfiguration = data
	return o
}
