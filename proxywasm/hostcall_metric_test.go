package proxywasm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mathetake/proxy-wasm-go/proxywasm/rawhostcall"
	"github.com/mathetake/proxy-wasm-go/proxywasm/types"
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
	rawhostcall.RegisterMockWASMHost(host)
	defer hostMutex.Unlock()

	for _, c := range []struct {
		name   string
		offset int64
	}{
		{name: "requests", offset: 100},
		{name: "rate", offset: -100},
	} {
		t.Run(c.name, func(t *testing.T) {
			// define metric
			m, err := HostCallDefineMetric(types.MetricTypeCounter, c.name)
			require.NoError(t, err)

			// increment
			require.NoError(t, m.Increment(c.offset))

			// get
			value, err := m.GetMetric()
			require.NoError(t, err)
			assert.Equal(t, uint64(c.offset), value)
		})
	}

}
