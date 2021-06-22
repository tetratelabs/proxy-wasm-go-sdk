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
	"fmt"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
)

type (
	// MetricCounter represents a counter metric.
	// Use DefineCounterMetric for initialization.
	MetricCounter uint32
	// MetricGauge represents a gauge metric.
	// Use DefineGaugeMetric for initialization.
	MetricGauge uint32
	// MetricHistogram represents a histogram metric.
	// Use DefineHistogramMetric for initialization.
	MetricHistogram uint32
)

// DefineCounterMetric returnes MetricCounter for a name.
func DefineCounterMetric(name string) MetricCounter {
	var id uint32
	ptr := internal.StringBytePtr(name)
	st := internal.ProxyDefineMetric(internal.MetricTypeCounter, ptr, len(name), &id)
	if err := internal.StatusToError(st); err != nil {
		panic(fmt.Sprintf("define metric of name %s: %v", name, internal.StatusToError(st)))
	}
	return MetricCounter(id)
}

// Value returnes the current value for this counter.
func (m MetricCounter) Value() uint64 {
	var val uint64
	st := internal.ProxyGetMetric(uint32(m), &val)
	if err := internal.StatusToError(st); err != nil {
		panic(fmt.Sprintf("get metric of  %d: %v", uint32(m), internal.StatusToError(st)))
	}
	return val
}

// Increment increments the current value by a offset for this counter.
func (m MetricCounter) Increment(offset uint64) {
	if err := internal.StatusToError(internal.ProxyIncrementMetric(uint32(m), int64(offset))); err != nil {
		panic(fmt.Sprintf("increment %d by %d: %v", uint32(m), offset, err))
	}
}

// DefineCounterMetric returnes MetricGauge for a name.
func DefineGaugeMetric(name string) MetricGauge {
	var id uint32
	ptr := internal.StringBytePtr(name)
	st := internal.ProxyDefineMetric(internal.MetricTypeGauge, ptr, len(name), &id)
	if err := internal.StatusToError(st); err != nil {
		panic(fmt.Sprintf("error define metric of name %s: %v", name, internal.StatusToError(st)))
	}
	return MetricGauge(id)
}

// Value returnes the current value for this gauge.
func (m MetricGauge) Value() int64 {
	var val uint64
	if err := internal.StatusToError(internal.ProxyGetMetric(uint32(m), &val)); err != nil {
		panic(fmt.Sprintf("get metric of  %d: %v", uint32(m), err))
	}
	return int64(val)
}

// Add adds a offset to the current value for this counter.
func (m MetricGauge) Add(offset int64) {
	if err := internal.StatusToError(internal.ProxyIncrementMetric(uint32(m), offset)); err != nil {
		panic(fmt.Sprintf("error adding %d by %d: %v", uint32(m), offset, err))
	}
}

// DefineHistogramMetric returnes MetricHistogram for a name.
func DefineHistogramMetric(name string) MetricHistogram {
	var id uint32
	ptr := internal.StringBytePtr(name)
	st := internal.ProxyDefineMetric(internal.MetricTypeHistogram, ptr, len(name), &id)
	if err := internal.StatusToError(st); err != nil {
		panic(fmt.Sprintf("error define metric of name %s: %v", name, internal.StatusToError(st)))
	}
	return MetricHistogram(id)
}

// Value returnes the current value for this histogram.
func (m MetricHistogram) Value() uint64 {
	var val uint64
	st := internal.ProxyGetMetric(uint32(m), &val)
	if err := internal.StatusToError(st); err != nil {
		panic(fmt.Sprintf("get metric of  %d: %v", uint32(m), internal.StatusToError(st)))
	}
	return val
}

// Record records a value for this histogram.
func (m MetricHistogram) Record(value uint64) {
	if err := internal.StatusToError(internal.ProxyRecordMetric(uint32(m), value)); err != nil {
		panic(fmt.Sprintf("error adding %d: %v", uint32(m), err))
	}
}
