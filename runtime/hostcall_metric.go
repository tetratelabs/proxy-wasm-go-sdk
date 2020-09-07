package runtime

import (
	"github.com/mathetake/proxy-wasm-go/runtime/rawhostcall"
	"github.com/mathetake/proxy-wasm-go/runtime/types"
)

type Metric uint32

func HostCallDefineMetric(metricType types.MetricType, name string) (Metric, error) {
	var id uint32
	ptr := unsafeGetStringBytePtr(name)
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
