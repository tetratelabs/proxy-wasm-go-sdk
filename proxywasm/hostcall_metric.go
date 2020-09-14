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

type (
	MetricCounter   uint32
	MetricGauge     uint32
	MetricHistogram uint32
)

// counter

func DefineCounterMetric(name string) (MetricCounter, error) {
	var id uint32
	ptr := stringBytePtr(name)
	st := rawhostcall.ProxyDefineMetric(types.MetricTypeCounter, ptr, len(name), &id)
	return MetricCounter(id), types.StatusToError(st)
}

func (m MetricCounter) ID() uint32 {
	return uint32(m)
}

func (m MetricCounter) Get() (uint64, error) {
	var val uint64
	st := rawhostcall.ProxyGetMetric(m.ID(), &val)
	return val, types.StatusToError(st)
}

func (m MetricCounter) Increment(offset uint64) error {
	return types.StatusToError(rawhostcall.ProxyIncrementMetric(m.ID(), int64(offset)))
}

// gauge

func DefineGaugeMetric(name string) (MetricGauge, error) {
	var id uint32
	ptr := stringBytePtr(name)
	st := rawhostcall.ProxyDefineMetric(types.MetricTypeGauge, ptr, len(name), &id)
	return MetricGauge(id), types.StatusToError(st)
}

func (m MetricGauge) ID() uint32 {
	return uint32(m)
}

func (m MetricGauge) Get() (int64, error) {
	var val uint64
	st := rawhostcall.ProxyGetMetric(m.ID(), &val)
	return int64(val), types.StatusToError(st)
}

func (m MetricGauge) Add(offset int64) error {
	return types.StatusToError(rawhostcall.ProxyIncrementMetric(m.ID(), offset))
}

// histogram

func DefineHistogramMetric(name string) (MetricHistogram, error) {
	var id uint32
	ptr := stringBytePtr(name)
	st := rawhostcall.ProxyDefineMetric(types.MetricTypeHistogram, ptr, len(name), &id)
	return MetricHistogram(id), types.StatusToError(st)
}

func (m MetricHistogram) ID() uint32 {
	return uint32(m)
}

func (m MetricHistogram) Get() (uint64, error) {
	var val uint64
	st := rawhostcall.ProxyGetMetric(m.ID(), &val)
	return val, types.StatusToError(st)
}

func (m MetricHistogram) Record(value uint64) error {
	return types.StatusToError(rawhostcall.ProxyRecordMetric(m.ID(), value))
}
