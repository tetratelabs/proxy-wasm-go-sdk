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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/log"
)

type logHost struct {
	internal.DefaultProxyWAMSHost
	t           *testing.T
	expMessage  string
	expLogLevel log.Level
	wasCalled   bool
}

func (l *logHost) ProxyLog(logLevel log.Level, messageData *byte, messageSize int) internal.Status {
	l.wasCalled = true
	actual := internal.RawBytePtrToString(messageData, messageSize)
	assert.Equal(l.t, l.expMessage, actual)
	assert.Equal(l.t, l.expLogLevel, logLevel)
	return internal.StatusOK
}

func TestHostCall_ForeignFunction(t *testing.T) {
	defer internal.RegisterMockWasmHost(internal.DefaultProxyWAMSHost{})()

	ret, err := CallForeignFunction("testFunc", []byte(""))
	require.NoError(t, err)
	require.Equal(t, []byte(nil), ret)
}

func TestHostCall_Logging(t *testing.T) {
	t.Run("trace", func(t *testing.T) {
		LogLevel = log.LevelTrace

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "trace",
			expLogLevel:          log.LevelTrace,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogTrace("trace")

		assert.True(t, lh.wasCalled)
	})

	t.Run("trace disabled", func(t *testing.T) {
		LogLevel = log.LevelDebug

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogTrace("trace")

		assert.False(t, lh.wasCalled)
	})

	t.Run("tracef", func(t *testing.T) {
		LogLevel = log.LevelTrace

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "trace: log",
			expLogLevel:          log.LevelTrace,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogTracef("trace: %s", "log")

		assert.True(t, lh.wasCalled)
	})

	t.Run("tracef disabled", func(t *testing.T) {
		LogLevel = log.LevelDebug

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogTracef("trace: %s", "log")

		assert.False(t, lh.wasCalled)
	})

	t.Run("debug", func(t *testing.T) {
		LogLevel = log.LevelDebug

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "abc",
			expLogLevel:          log.LevelDebug,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogDebug("abc")

		assert.True(t, lh.wasCalled)
	})

	t.Run("debug disabled", func(t *testing.T) {
		LogLevel = log.LevelInfo

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogDebug("abc")

		assert.False(t, lh.wasCalled)
	})

	t.Run("debugf", func(t *testing.T) {
		LogLevel = log.LevelDebug

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "debug: log",
			expLogLevel:          log.LevelDebug,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogDebugf("debug: %s", "log")

		assert.True(t, lh.wasCalled)
	})

	t.Run("debugf disabled", func(t *testing.T) {
		LogLevel = log.LevelInfo

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogDebugf("debug: %s", "log")

		assert.False(t, lh.wasCalled)
	})

	t.Run("info", func(t *testing.T) {
		LogLevel = log.LevelInfo

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "info",
			expLogLevel:          log.LevelInfo,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogInfo("info")

		assert.True(t, lh.wasCalled)
	})

	t.Run("info disabled", func(t *testing.T) {
		LogLevel = log.LevelWarn

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogInfo("info")

		assert.False(t, lh.wasCalled)
	})

	t.Run("infof", func(t *testing.T) {
		LogLevel = log.LevelInfo

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "info: log: 10",
			expLogLevel:          log.LevelInfo,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogInfof("info: %s: %d", "log", 10)

		assert.True(t, lh.wasCalled)
	})

	t.Run("infof disabled", func(t *testing.T) {
		LogLevel = log.LevelWarn

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogInfof("info: %s: %d", "log", 10)

		assert.False(t, lh.wasCalled)
	})

	t.Run("warn", func(t *testing.T) {
		LogLevel = log.LevelWarn

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "warn",
			expLogLevel:          log.LevelWarn,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogWarn("warn")

		assert.True(t, lh.wasCalled)
	})

	t.Run("warn disabled", func(t *testing.T) {
		LogLevel = log.LevelError

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogWarn("warn")

		assert.False(t, lh.wasCalled)
	})

	t.Run("warnf", func(t *testing.T) {
		LogLevel = log.LevelWarn

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "warn: log: 10",
			expLogLevel:          log.LevelWarn,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogWarnf("warn: %s: %d", "log", 10)

		assert.True(t, lh.wasCalled)
	})

	t.Run("warnf disabled", func(t *testing.T) {
		LogLevel = log.LevelError

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogWarnf("warn: %s: %d", "log", 10)

		assert.False(t, lh.wasCalled)
	})

	t.Run("error", func(t *testing.T) {
		LogLevel = log.LevelError

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "error",
			expLogLevel:          log.LevelError,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogError("error")

		assert.True(t, lh.wasCalled)
	})

	t.Run("error disabled", func(t *testing.T) {
		LogLevel = log.LevelCritical

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogError("error")

		assert.False(t, lh.wasCalled)
	})

	t.Run("errorf", func(t *testing.T) {
		LogLevel = log.LevelError

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "error: log: 10",
			expLogLevel:          log.LevelError,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogErrorf("error: %s: %d", "log", 10)

		assert.True(t, lh.wasCalled)
	})

	t.Run("errorf disabled", func(t *testing.T) {
		LogLevel = log.LevelCritical

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogErrorf("error: %s: %d", "log", 10)

		assert.False(t, lh.wasCalled)
	})

	t.Run("critical", func(t *testing.T) {
		LogLevel = log.LevelCritical

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "critical error",
			expLogLevel:          log.LevelCritical,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogCritical("critical error")

		assert.True(t, lh.wasCalled)
	})

	t.Run("critical disabled", func(t *testing.T) {
		LogLevel = log.LevelDisabled

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogCritical("critical error")

		assert.False(t, lh.wasCalled)
	})

	t.Run("criticalf", func(t *testing.T) {
		LogLevel = log.LevelCritical

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "critical: log: 10",
			expLogLevel:          log.LevelCritical,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogCriticalf("critical: %s: %d", "log", 10)

		assert.True(t, lh.wasCalled)
	})

	t.Run("criticalf disabled", func(t *testing.T) {
		LogLevel = log.LevelDisabled

		lh := &logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
		}
		release := internal.RegisterMockWasmHost(lh)
		defer release()

		LogCriticalf("critical: %s: %d", "log", 10)

		assert.False(t, lh.wasCalled)
	})
}

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
