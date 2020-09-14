package proxywasm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/rawhostcall"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type metricProxyWASMHost struct {
	rawhostcall.DefaultProxyWAMSHost
	idToValue map[uint32]uint64
	idToType  map[uint32]types.MetricType
	nameToID  map[string]uint32
}

func (m metricProxyWASMHost) ProxyDefineMetric(metricType types.MetricType,
	metricNameData *byte, metricNameSize int, returnMetricIDPtr *uint32) types.Status {
	name := rawBytePtrToString(metricNameData, metricNameSize)
	id, ok := m.nameToID[name]
	if !ok {
		id = uint32(len(m.nameToID))
		m.nameToID[name] = id
		m.idToValue[id] = 0
		m.idToType[id] = metricType
	}
	*returnMetricIDPtr = id
	return types.StatusOK
}

func (m metricProxyWASMHost) ProxyIncrementMetric(metricID uint32, offset int64) types.Status {
	val, ok := m.idToValue[metricID]
	if !ok {
		return types.StatusBadArgument
	}

	m.idToValue[metricID] = val + uint64(offset)
	return types.StatusOK
}

func (m metricProxyWASMHost) ProxyRecordMetric(metricID uint32, value uint64) types.Status {
	_, ok := m.idToValue[metricID]
	if !ok {
		return types.StatusBadArgument
	}
	m.idToValue[metricID] = value
	return types.StatusOK
}

func (m metricProxyWASMHost) ProxyGetMetric(metricID uint32, returnMetricValue *uint64) types.Status {
	value, ok := m.idToValue[metricID]
	if !ok {
		return types.StatusBadArgument
	}
	*returnMetricValue = value
	return types.StatusOK
}

func TestHostCall_Metric(t *testing.T) {
	host := metricProxyWASMHost{
		rawhostcall.DefaultProxyWAMSHost{},
		map[uint32]uint64{},
		map[uint32]types.MetricType{},
		map[string]uint32{},
	}
	hostMutex.Lock()
	defer hostMutex.Unlock()
	rawhostcall.RegisterMockWASMHost(host)

	t.Run("counter", func(t *testing.T) {
		for _, c := range []struct {
			name   string
			offset uint64
		}{
			{name: "requests", offset: 100},
		} {
			t.Run(c.name, func(t *testing.T) {
				// define metric
				m, err := DefineCounterMetric(c.name)
				require.NoError(t, err)

				// increment
				require.NoError(t, m.Increment(c.offset))

				// get
				value, err := m.Get()
				require.NoError(t, err)
				assert.Equal(t, c.offset, value)
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
				m, err := DefineGaugeMetric(c.name)
				require.NoError(t, err)

				// increment
				require.NoError(t, m.Add(c.offset))

				// get
				value, err := m.Get()
				require.NoError(t, err)
				assert.Equal(t, c.offset, value)
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
				m, err := DefineHistogramMetric(c.name)
				require.NoError(t, err)

				// record
				require.NoError(t, m.Record(c.value))

				// get
				value, err := m.Get()
				require.NoError(t, err)
				assert.Equal(t, c.value, value)
			})
		}
	})
}
