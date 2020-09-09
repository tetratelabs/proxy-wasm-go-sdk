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

package proxywasm

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/rawhostcall"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type Metric uint32

func HostCallDefineMetric(metricType types.MetricType, name string) (Metric, error) {
	var id uint32
	ptr := stringBytePtr(name)
	st := rawhostcall.ProxyDefineMetric(metricType, ptr, len(name), &id)
	return Metric(id), types.StatusToError(st)
}

func (m Metric) ID() uint32 {
	return uint32(m)
}

func (m Metric) Increment(offset int64) error {
	return types.StatusToError(rawhostcall.ProxyIncrementMetric(m.ID(), offset))
}

func (m Metric) RecordMetric(value uint64) error {
	return types.StatusToError(rawhostcall.ProxyRecordMetric(m.ID(), value))
}

func (m Metric) GetMetric() (uint64, error) {
	var val uint64
	st := rawhostcall.ProxyGetMetric(m.ID(), &val)
	return val, types.StatusToError(st)
}
