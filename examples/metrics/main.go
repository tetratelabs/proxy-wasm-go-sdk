package main

import (
	"strconv"

	"github.com/mathetake/proxy-wasm-go/runtime"
	"github.com/mathetake/proxy-wasm-go/runtime/types"
)

func main() {
	runtime.SetNewRootContext(func(uint32) runtime.RootContext { return &metrics{} })
	runtime.SetNewHttpContext(func(uint32) runtime.HttpContext { return &metrics{} })
}

var counter runtime.Metric

const metricsName = "proxy_wasm_go.request_counter"

type metrics struct{ runtime.DefaultContext }

// override
func (ctx *metrics) OnVMStart(_ int) bool {
	ct, err := runtime.HostCallDefineMetric(types.MetricTypeCounter, metricsName)
	if err != nil {
		runtime.LogCritical("error defining metrics: ", err.Error())
	}
	counter = ct
	return true
}

// override
func (ctx *metrics) OnHttpRequestHeaders(_ int, _ bool) types.Action {
	prev, err := counter.GetMetric()
	if err != nil {
		runtime.LogCritical("error retrieving previous metric: ", err.Error())
	}

	runtime.LogInfo("previous value of ", metricsName+": ", strconv.Itoa(int(prev)))

	if err := counter.Increment(1); err != nil {
		runtime.LogCritical("error incrementing metrics", err.Error())
	}
	runtime.LogInfo("incremented")
	return types.ActionContinue
}
