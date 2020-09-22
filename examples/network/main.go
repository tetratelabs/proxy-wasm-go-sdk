// Copyright 2020 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

var (
	connectionCounterName = "proxy_wasm_go.connection_counter"
	counter               proxywasm.MetricCounter
)

func main() {
	proxywasm.SetNewStreamContext(func(contextID uint32) proxywasm.StreamContext { return context{} })
	proxywasm.SetNewRootContext(func(contextID uint32) proxywasm.RootContext { return context{} })
}

type context struct{ proxywasm.DefaultContext }

func (ctx context) OnVMStart(int) bool {
	var err error
	counter, err = proxywasm.DefineCounterMetric(connectionCounterName)
	if err != nil {
		proxywasm.LogCritical("failed to initialize connection counter: ", err.Error())
	}
	return true
}

func (ctx context) OnNewConnection() types.Action {
	proxywasm.LogInfo("new connection!")
	return types.ActionContinue
}

func (ctx context) OnDownstreamData(dataSize int, _ bool) types.Action {
	// TODO: dispatch http call

	if dataSize == 0 {
		return types.ActionContinue
	}

	data, err := proxywasm.HostCallGetDownStreamData(0, dataSize)
	if err != nil && err != types.ErrorStatusNotFound {
		proxywasm.LogCritical(err.Error())
	}

	proxywasm.LogInfo("downstream data received: ", string(data))
	return types.ActionContinue
}

func (ctx context) OnDownstreamClose(types.PeerType) {
	proxywasm.LogInfo("downstream connection close!")
	return
}

func (ctx context) OnUpstreamData(dataSize int, _ bool) types.Action {

	if dataSize == 0 {
		return types.ActionContinue
	}

	data, err := proxywasm.HostCallGetUpstreamData(0, dataSize)
	if err != nil && err != types.ErrorStatusNotFound {
		proxywasm.LogCritical(err.Error())
	}

	proxywasm.LogInfo("upstream data received: ", string(data))
	return types.ActionContinue
}

func (ctx context) OnDone() bool {
	err := counter.Increment(1)
	if err != nil {
		proxywasm.LogCritical("failed to increment connection counter: ", err.Error())
	}
	proxywasm.LogInfo("connection complete!")
	return true
}
