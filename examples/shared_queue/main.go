package main

import (
	"strconv"

	"github.com/mathetake/proxy-wasm-go-sdk/proxywasm"
	"github.com/mathetake/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetNewRootContext(func(uint32) proxywasm.RootContext { return &queue{} })
	proxywasm.SetNewHttpContext(func(uint32) proxywasm.HttpContext { return &queue{} })
}

type queue struct{ proxywasm.DefaultContext }

const (
	queueName        = "proxy_wasm_go.queue"
	tickMilliseconds = 1000
)

var queueID uint32

// override
func (ctx *queue) OnVMStart(int) bool {
	qID, err := proxywasm.HostCallRegisterSharedQueue(queueName)
	if err != nil {
		panic(err.Error())
	}
	queueID = qID
	proxywasm.LogInfo("queue registered, name: ", queueName, ", id: ", strconv.Itoa(int(qID)))

	if err := proxywasm.HostCallSetTickPeriodMilliSeconds(tickMilliseconds); err != nil {
		proxywasm.LogCritical("failed to set tick period: ", err.Error())
	}
	proxywasm.LogInfo("set tick period milliseconds: ", strconv.Itoa(tickMilliseconds))
	return true
}

// override
func (ctx *queue) OnQueueReady(queueID uint32) {
	proxywasm.LogInfo("queue ready: ", strconv.Itoa(int(queueID)))
}

// override
func (ctx *queue) OnHttpRequestHeaders(int, bool) types.Action {
	for _, msg := range []string{"hello", "world", "hello", "proxy-wasm"} {
		if err := proxywasm.HostCallEnqueueSharedQueue(queueID, []byte(msg)); err != nil {
			proxywasm.LogCritical("error queueing: ", err.Error())
		}
	}
	return types.ActionContinue
}

// override
func (ctx *queue) OnTick() {
	data, err := proxywasm.HostCallDequeueSharedQueue(queueID)
	switch err {
	case types.ErrorStatusEmpty:
		return
	case nil:
		proxywasm.LogInfo("dequed data: ", string(data))
	default:
		proxywasm.LogCritical("error retrieving data from queue ", strconv.Itoa(int(queueID)), ", ", err.Error())
	}
}
