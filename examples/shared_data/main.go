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
	"strconv"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetNewRootContext(func(uint32) proxywasm.RootContext { return data{} })
	proxywasm.SetNewHttpContext(func(uint32) proxywasm.HttpContext { return data{} })
}

type data struct{ proxywasm.DefaultContext }

const sharedDataKey = "shared_data_key"

// override
func (ctx data) OnVMStart(vid int) bool {
	_, cas, err := proxywasm.HostCallGetSharedData(sharedDataKey)
	if err != nil {
		proxywasm.LogWarn("error getting shared data on OnVMStart: ", err.Error())
	}

	if err = proxywasm.HostCallSetSharedData(sharedDataKey, []byte{0}, cas); err != nil {
		proxywasm.LogWarn("error setting shared data on OnVMStart: ", err.Error())
	}
	return true
}

// override
func (ctx data) OnHttpRequestHeaders(int, bool) types.Action {
	value, cas, err := proxywasm.HostCallGetSharedData(sharedDataKey)
	if err != nil {
		proxywasm.LogWarn("error getting shared data on OnHttpRequestHeaders: ", err.Error())
		return types.ActionContinue
	}

	value[0]++
	if err := proxywasm.HostCallSetSharedData(sharedDataKey, value, cas); err != nil {
		proxywasm.LogWarn("error setting shared data on OnHttpRequestHeaders: ", err.Error())
		return types.ActionContinue
	}

	proxywasm.LogInfo("shared value: ", strconv.Itoa(int(value[0])))
	return types.ActionContinue
}
