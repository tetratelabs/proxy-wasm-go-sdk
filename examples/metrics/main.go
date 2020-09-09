package main

import (
	"strconv"

	"github.com/mathetake/proxy-wasm-go-sdk/proxywasm"
	"github.com/mathetake/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetNewRootContext(func(uint32) proxywasm.RootContext { return &metrics{} })
	proxywasm.SetNewHttpContext(func(uint32) proxywasm.HttpContext { return &metrics{} })
}

var counter proxywasm.Metric

const metricsName = "proxy_wasm_go.request_counter"

type metrics struct{ proxywasm.DefaultContext }

// override
func (ctx *metrics) OnVMStart(int) bool {
	ct, err := proxywasm.HostCallDefineMetric(types.MetricTypeCounter, metricsName)
	if err != nil {
		proxywasm.LogCritical("error defining metrics: ", err.Error())
	}
	counter = ct
	return true
}

// override
func (ctx *metrics) OnHttpRequestHeaders(int, bool) types.Action {
	prev, err := counter.GetMetric()
	if err != nil {
		proxywasm.LogCritical("error retrieving previous metric: ", err.Error())
	}

	proxywasm.LogInfo("previous value of ", metricsName, ": ", strconv.Itoa(int(prev)))

	if err := counter.Increment(1); err != nil {
		proxywasm.LogCritical("error incrementing metrics", err.Error())
	}
	proxywasm.LogInfo("incremented")
	return types.ActionContinue
}
