// Copyright 2021 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxywasm

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
)

type metricProxyWasmHost struct {
	internal.DefaultProxyWAMSHost
	idToValue map[uint32]uint64
	idToType  map[uint32]internal.MetricType
	nameToID  map[string]uint32
}

func (m metricProxyWasmHost) ProxyDefineMetric(metricType internal.MetricType,
	metricNameData *byte, metricNameSize int, returnMetricIDPtr *uint32) internal.Status {
	name := internal.RawBytePtrToString(metricNameData, metricNameSize)
	id, ok := m.nameToID[name]
	if !ok {
		id = uint32(len(m.nameToID))
		m.nameToID[name] = id
		m.idToValue[id] = 0
		m.idToType[id] = metricType
	}
	*returnMetricIDPtr = id
	return internal.StatusOK
}

func (m metricProxyWasmHost) ProxyIncrementMetric(metricID uint32, offset int64) internal.Status {
	val, ok := m.idToValue[metricID]
	if !ok {
		return internal.StatusBadArgument
	}

	m.idToValue[metricID] = val + uint64(offset)
	return internal.StatusOK
}

func (m metricProxyWasmHost) ProxyRecordMetric(metricID uint32, value uint64) internal.Status {
	_, ok := m.idToValue[metricID]
	if !ok {
		return internal.StatusBadArgument
	}
	m.idToValue[metricID] = value
	return internal.StatusOK
}

func (m metricProxyWasmHost) ProxyGetMetric(metricID uint32, returnMetricValue *uint64) internal.Status {
	value, ok := m.idToValue[metricID]
	if !ok {
		return internal.StatusBadArgument
	}
	*returnMetricValue = value
	return internal.StatusOK
}

func TestHostCall_Metric(t *testing.T) {
	host := metricProxyWasmHost{
		internal.DefaultProxyWAMSHost{},
		map[uint32]uint64{},
		map[uint32]internal.MetricType{},
		map[string]uint32{},
	}
	release := internal.RegisterMockWasmHost(host)
	defer release()

	t.Run("counter", func(t *testing.T) {
		for _, c := range []struct {
			name   string
			offset uint64
		}{
			{name: "requests", offset: 100},
		} {
			t.Run(c.name, func(t *testing.T) {
				// define metric
				m := DefineCounterMetric(c.name)

				// increment
				m.Increment(c.offset)

				// get
				require.Equal(t, c.offset, m.Value())
			})
		}
	})

	t.Run("gauge", func(t *testing.T) {
		for _, c := range []struct {
			name   string
			offset int64
		}{
			{name: "rate", offset: -50},
		} {
			t.Run(c.name, func(t *testing.T) {
				// define metric
				m := DefineGaugeMetric(c.name)

				// increment
				m.Add(c.offset)

				// get
				require.Equal(t, c.offset, m.Value())
			})
		}
	})

	t.Run("histogram", func(t *testing.T) {
		for _, c := range []struct {
			name  string
			value uint64
		}{
			{name: "request count", value: 10000},
		} {
			t.Run(c.name, func(t *testing.T) {
				// define metric
				m := DefineHistogramMetric(c.name)

				// record
				m.Record(c.value)

				// get
				require.Equal(t, c.value, m.Value())
			})
		}
	})
}
