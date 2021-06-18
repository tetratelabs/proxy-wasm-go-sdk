// Copyright 2020-2021 Tetrate
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
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
)

type (
	MetricCounter   uint32
	MetricGauge     uint32
	MetricHistogram uint32
)

func DefineCounterMetric(name string) MetricCounter {
	var id uint32
	ptr := internal.StringBytePtr(name)
	st := internal.ProxyDefineMetric(internal.MetricTypeCounter, ptr, len(name), &id)
	if err := internal.StatusToError(st); err != nil {
		LogCriticalf("define metric of name %s: %v", name, internal.StatusToError(st))
		panic("") // abort
	}
	return MetricCounter(id)
}

func (m MetricCounter) ID() uint32 {
	return uint32(m)
}

func (m MetricCounter) Get() uint64 {
	var val uint64
	st := internal.ProxyGetMetric(m.ID(), &val)
	if err := internal.StatusToError(st); err != nil {
		LogCriticalf("get metric of  %d: %v", m.ID(), internal.StatusToError(st))
		panic("") // abort
	}
	return val
}

func (m MetricCounter) Increment(offset uint64) {
	if err := internal.StatusToError(internal.ProxyIncrementMetric(m.ID(), int64(offset))); err != nil {
		LogCriticalf("increment %d by %d: %v", m.ID(), offset, err)
		panic("") // abort
	}
}

func DefineGaugeMetric(name string) MetricGauge {
	var id uint32
	ptr := internal.StringBytePtr(name)
	st := internal.ProxyDefineMetric(internal.MetricTypeGauge, ptr, len(name), &id)
	if err := internal.StatusToError(st); err != nil {
		LogCriticalf("error define metric of name %s: %v", name, internal.StatusToError(st))
		panic("") // abort
	}
	return MetricGauge(id)
}

func (m MetricGauge) ID() uint32 {
	return uint32(m)
}

func (m MetricGauge) Get() int64 {
	var val uint64
	if err := internal.StatusToError(internal.ProxyGetMetric(m.ID(), &val)); err != nil {
		LogCriticalf("get metric of  %d: %v", m.ID(), err)
		panic("") // abort
	}
	return int64(val)
}

func (m MetricGauge) Add(offset int64) {
	if err := internal.StatusToError(internal.ProxyIncrementMetric(m.ID(), offset)); err != nil {
		LogCriticalf("error adding %d by %d: %v", m.ID(), offset, err)
		panic("") // abort
	}
}

func DefineHistogramMetric(name string) MetricHistogram {
	var id uint32
	ptr := internal.StringBytePtr(name)
	st := internal.ProxyDefineMetric(internal.MetricTypeHistogram, ptr, len(name), &id)
	if err := internal.StatusToError(st); err != nil {
		LogCriticalf("error define metric of name %s: %v", name, internal.StatusToError(st))
		panic("") // abort
	}
	return MetricHistogram(id)
}

func (m MetricHistogram) ID() uint32 {
	return uint32(m)
}

func (m MetricHistogram) Get() uint64 {
	var val uint64
	st := internal.ProxyGetMetric(m.ID(), &val)
	if err := internal.StatusToError(st); err != nil {
		LogCriticalf("get metric of  %d: %v", m.ID(), internal.StatusToError(st))
		panic("") // abort
	}
	return val
}

func (m MetricHistogram) Record(value uint64) {
	if err := internal.StatusToError(internal.ProxyRecordMetric(m.ID(), value)); err != nil {
		LogCriticalf("error adding %d: %v", m.ID(), err)
		panic("") // abort
	}
}
