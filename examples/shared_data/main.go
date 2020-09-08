package main

import (
	"strconv"

	"github.com/mathetake/proxy-wasm-go/runtime"
	"github.com/mathetake/proxy-wasm-go/runtime/types"
)

func main() {
	runtime.SetNewRootContext(func(uint32) runtime.RootContext { return &data{} })
	runtime.SetNewHttpContext(func(uint32) runtime.HttpContext { return &data{} })
}

type data struct{ runtime.DefaultContext }

const sharedDataKey = "shared_data_key"

// override
func (ctx *data) OnVMStart(vid int) bool {
	_, cas, err := runtime.HostCallGetSharedData(sharedDataKey)
	if err != nil {
		runtime.LogWarn("error getting shared data on OnVMStart: ", err.Error())
	}

	if err = runtime.HostCallSetSharedData(sharedDataKey, []byte{0}, cas); err != nil {
		runtime.LogWarn("error setting shared data on OnVMStart: ", err.Error())
	}
	return true
}

// override
func (ctx *data) OnHttpRequestHeaders(int, bool) types.Action {
	value, cas, err := runtime.HostCallGetSharedData(sharedDataKey)
	if err != nil {
		runtime.LogWarn("error getting shared data on OnHttpRequestHeaders: ", err.Error())
		return types.ActionContinue
	}

	value[0]++
	if err := runtime.HostCallSetSharedData(sharedDataKey, value, cas); err != nil {
		runtime.LogWarn("error setting shared data on OnHttpRequestHeaders: ", err.Error())
		return types.ActionContinue
	}

	runtime.LogInfo("shared value: ", strconv.Itoa(int(value[0])))
	return types.ActionContinue
}
