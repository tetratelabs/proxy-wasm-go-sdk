package main

import (
	"strconv"

	"github.com/mathetake/proxy-wasm-go/runtime"
	"github.com/mathetake/proxy-wasm-go/runtime/types"
)

func main() {
	runtime.SetNewRootContext(func(uint32) runtime.RootContext { return &queue{} })
	runtime.SetNewHttpContext(func(uint32) runtime.HttpContext { return &queue{} })
}

type queue struct{ runtime.DefaultContext }

const (
	queueName        = "proxy_wasm_go.queue"
	tickMilliseconds = 1000
)

var queueID uint32

// override
func (ctx *queue) OnVMStart(_ int) bool {
	qID, err := runtime.HostCallRegisterSharedQueue(queueName)
	if err != nil {
		panic(err.Error())
	}
	queueID = qID
	runtime.LogInfo("queue registered, name: ", queueName, ", id: ", strconv.Itoa(int(qID)))

	if err := runtime.HostCallSetTickPeriodMilliSeconds(tickMilliseconds); err != nil {
		runtime.LogCritical("failed to set tick period: ", err.Error())
	}
	runtime.LogInfo("set tick period milliseconds: ", strconv.Itoa(tickMilliseconds))
	return true
}

// override
func (ctx *queue) OnQueueReady(queueID uint32) {
	runtime.LogInfo("queue ready: ", strconv.Itoa(int(queueID)))
}

// override
func (ctx *queue) OnHttpRequestHeaders(int, bool) types.Action {
	for _, msg := range []string{"hello", "world", "hello", "proxy-wasm"} {
		if err := runtime.HostCallEnqueueSharedQueue(queueID, []byte(msg)); err != nil {
			runtime.LogCritical("error queueing: ", err.Error())
		}
	}
	return types.ActionContinue
}

// override
func (ctx *queue) OnTick() {
	data, err := runtime.HostCallDequeueSharedQueue(queueID)
	switch err {
	case types.ErrorStatusEmpty:
		return
	case nil:
		runtime.LogInfo("dequed data: ", string(data))
	default:
		runtime.LogCritical("error retrieving data from queue ", strconv.Itoa(int(queueID)), ", ", err.Error())
	}
}
