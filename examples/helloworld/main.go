package main

import (
	"strconv"

	"github.com/mathetake/proxy-wasm-go/runtime"
)

func main() {
	runtime.SetNewRootContext(newHelloWorld)
}

type helloWorld struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	runtime.DefaultContext
	contextID uint32
}

func newHelloWorld(contextID uint32) runtime.RootContext {
	return &helloWorld{contextID: contextID}
}

// override
func (ctx *helloWorld) OnVMStart(_ int) bool {
	runtime.LogInfo("proxy_on_vm_start from Go!")
	if err := runtime.HostCallSetTickPeriodMilliSeconds(1000); err != nil {
		runtime.LogCritical("failed to set tick period: ", err.Error())
	}
	return true
}

// override
func (ctx *helloWorld) OnTick() {
	t := runtime.HostCallGetCurrentTime()
	runtime.LogInfo("OnTick on ", strconv.FormatUint(uint64(ctx.contextID), 10),
		", it's ", strconv.FormatInt(t, 10))
}
