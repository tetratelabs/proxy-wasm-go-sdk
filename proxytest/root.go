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

package proxytest

import (
	"fmt"
	"log"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type (
	rootHostEmulator struct {
		logs       [types.LogLevelMax][]string
		tickPeriod uint32

		queues      map[uint32][][]byte
		queueNameID map[string]uint32

		sharedDataKVS map[string]*sharedData

		metricIDToValue map[uint32]uint64
		metricIDToType  map[uint32]types.MetricType
		metricNameToID  map[string]uint32

		httpContextIDToCalloutInfos map[uint32][]HttpCalloutAttribute // key: contextID
		httpCalloutIDToContextID    map[uint32]uint32                 // key: calloutID
		httpCalloutResponse         map[uint32]struct {               // key: calloutID
			headers  types.Headers
			trailers types.Trailers
			body     []byte
		}

		pluginConfiguration, vmConfiguration []byte

		activeCalloutID uint32
	}

	HttpCalloutAttribute struct {
		CalloutID uint32
		Upstream  string
		Headers   types.Headers
		Trailers  types.Trailers
		Body      []byte
	}
)

type sharedData struct {
	data []byte
	cas  uint32
}

func newRootHostEmulator(pluginConfiguration, vmConfiguration []byte) *rootHostEmulator {
	host := &rootHostEmulator{
		queues:                      map[uint32][][]byte{},
		queueNameID:                 map[string]uint32{},
		sharedDataKVS:               map[string]*sharedData{},
		metricIDToValue:             map[uint32]uint64{},
		metricIDToType:              map[uint32]types.MetricType{},
		metricNameToID:              map[string]uint32{},
		httpContextIDToCalloutInfos: map[uint32][]HttpCalloutAttribute{},
		httpCalloutIDToContextID:    map[uint32]uint32{},
		httpCalloutResponse: map[uint32]struct {
			headers  types.Headers
			trailers types.Trailers
			body     []byte
		}{},

		pluginConfiguration: pluginConfiguration,
		vmConfiguration:     vmConfiguration,
	}
	return host
}

// impl rawhostcall.ProxyWASMHost
func (r *rootHostEmulator) ProxyLog(logLevel types.LogLevel, messageData *byte, messageSize int) types.Status {
	str := proxywasm.RawBytePtrToString(messageData, messageSize)

	log.Printf("proxy_%s_log: %s", logLevel, str)
	r.logs[logLevel] = append(r.logs[logLevel], str)
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func (r *rootHostEmulator) ProxySetTickPeriodMilliseconds(period uint32) types.Status {
	r.tickPeriod = period
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func (r *rootHostEmulator) ProxyRegisterSharedQueue(nameData *byte, nameSize int, returnID *uint32) types.Status {
	name := proxywasm.RawBytePtrToString(nameData, nameSize)
	if id, ok := r.queueNameID[name]; ok {
		*returnID = id
		return types.StatusOK
	}

	id := uint32(len(r.queues))
	r.queues[id] = [][]byte{}
	r.queueNameID[name] = id
	*returnID = id
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func (r *rootHostEmulator) ProxyDequeueSharedQueue(queueID uint32, returnValueData **byte, returnValueSize *int) types.Status {
	queue, ok := r.queues[queueID]
	if !ok {
		log.Printf("queue %d is not found", queueID)
		return types.StatusNotFound
	} else if len(queue) == 0 {
		log.Printf("queue %d is empty", queueID)
		return types.StatusEmpty
	}

	data := queue[0]
	*returnValueData = &data[0]
	*returnValueSize = len(data)
	r.queues[queueID] = queue[1:]
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func (r *rootHostEmulator) ProxyEnqueueSharedQueue(queueID uint32, valueData *byte, valueSize int) types.Status {
	queue, ok := r.queues[queueID]
	if !ok {
		log.Printf("queue %d is not found", queueID)
		return types.StatusNotFound
	}

	r.queues[queueID] = append(queue, proxywasm.RawBytePtrToByteSlice(valueData, valueSize))

	// note that this behavior is not accurate for some old host implementations:
	//	see: https://github.com/proxy-wasm/proxy-wasm-cpp-host/pull/36
	proxywasm.ProxyOnQueueReady(RootContextID, queueID) // Note that this behavior is not accurate on Istio before 1.8.x
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func (r *rootHostEmulator) ProxyGetSharedData(keyData *byte, keySize int,
	returnValueData **byte, returnValueSize *int, returnCas *uint32) types.Status {
	key := proxywasm.RawBytePtrToString(keyData, keySize)

	value, ok := r.sharedDataKVS[key]
	if !ok {
		return types.StatusNotFound
	}

	*returnValueSize = len(value.data)
	*returnValueData = &value.data[0]
	*returnCas = value.cas
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func (r *rootHostEmulator) ProxySetSharedData(keyData *byte, keySize int,
	valueData *byte, valueSize int, cas uint32) types.Status {
	key := proxywasm.RawBytePtrToString(keyData, keySize)
	value := proxywasm.RawBytePtrToByteSlice(valueData, valueSize)

	prev, ok := r.sharedDataKVS[key]
	if !ok {
		r.sharedDataKVS[key] = &sharedData{
			data: value,
			cas:  cas + 1,
		}
		return types.StatusOK
	}

	if prev.cas != cas {
		return types.StatusCasMismatch
	}

	r.sharedDataKVS[key].cas = cas + 1
	r.sharedDataKVS[key].data = value
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func (r *rootHostEmulator) ProxyDefineMetric(metricType types.MetricType,
	metricNameData *byte, metricNameSize int, returnMetricIDPtr *uint32) types.Status {
	name := proxywasm.RawBytePtrToString(metricNameData, metricNameSize)
	id, ok := r.metricNameToID[name]
	if !ok {
		id = uint32(len(r.metricNameToID))
		r.metricNameToID[name] = id
		r.metricIDToValue[id] = 0
		r.metricIDToType[id] = metricType
	}
	*returnMetricIDPtr = id
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func (r *rootHostEmulator) ProxyIncrementMetric(metricID uint32, offset int64) types.Status {
	val, ok := r.metricIDToValue[metricID]
	if !ok {
		return types.StatusBadArgument
	}

	r.metricIDToValue[metricID] = val + uint64(offset)
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func (r *rootHostEmulator) ProxyRecordMetric(metricID uint32, value uint64) types.Status {
	_, ok := r.metricIDToValue[metricID]
	if !ok {
		return types.StatusBadArgument
	}
	r.metricIDToValue[metricID] = value
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func (r *rootHostEmulator) ProxyGetMetric(metricID uint32, returnMetricValue *uint64) types.Status {
	value, ok := r.metricIDToValue[metricID]
	if !ok {
		return types.StatusBadArgument
	}
	*returnMetricValue = value
	return types.StatusOK
}

// impl rawhostcall.ProxyWASMHost
func (r *rootHostEmulator) ProxyHttpCall(upstreamData *byte, upstreamSize int, headerData *byte, headerSize int, bodyData *byte,
	bodySize int, trailersData *byte, trailersSize int, timeout uint32, calloutIDPtr *uint32) types.Status {
	upstream := proxywasm.RawBytePtrToString(upstreamData, upstreamSize)
	body := proxywasm.RawBytePtrToString(bodyData, bodySize)
	headers := proxywasm.DeserializeMap(proxywasm.RawBytePtrToByteSlice(headerData, headerSize))
	trailers := proxywasm.DeserializeMap(proxywasm.RawBytePtrToByteSlice(trailersData, trailersSize))

	log.Printf("[http callout to %s] timeout: %d", upstream, timeout)
	log.Printf("[http callout to %s] headers: %v", upstream, headers)
	log.Printf("[http callout to %s] body: %s", upstream, body)
	log.Printf("[http callout to %s] trailers: %v", upstream, trailers)

	calloutID := uint32(len(r.httpCalloutIDToContextID))
	contextID := proxywasm.VMStateGetActiveContextID()
	r.httpCalloutIDToContextID[calloutID] = contextID
	r.httpContextIDToCalloutInfos[contextID] = append(r.httpContextIDToCalloutInfos[contextID], HttpCalloutAttribute{
		CalloutID: calloutID,
		Upstream:  upstream,
		Headers:   headers,
		Trailers:  trailers,
		Body:      []byte(body),
	})

	*calloutIDPtr = calloutID
	return types.StatusOK
}

// // impl rawhostcall.ProxyWASMHost: delegated from hostEmulator
func (r *rootHostEmulator) rootHostEmulatorProxyGetHeaderMapPairs(mapType types.MapType, returnValueData **byte, returnValueSize *int) types.Status {
	res, ok := r.httpCalloutResponse[r.activeCalloutID]
	if !ok {
		log.Fatalf("callout response unregistered for %d", r.activeCalloutID)
	}

	var raw []byte
	switch mapType {
	case types.MapTypeHttpCallResponseHeaders:
		raw = proxywasm.SerializeMap(res.headers)
	case types.MapTypeHttpCallResponseTrailers:
		raw = proxywasm.SerializeMap(res.trailers)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	*returnValueData = &raw[0]
	*returnValueSize = len(raw)
	return types.StatusOK
}

// // impl rawhostcall.ProxyWASMHost: delegated from hostEmulator
func (r *rootHostEmulator) rootHostEmulatorProxyGetMapValue(mapType types.MapType, keyData *byte,
	keySize int, returnValueData **byte, returnValueSize *int) types.Status {
	res, ok := r.httpCalloutResponse[r.activeCalloutID]
	if !ok {
		log.Fatalf("callout response unregistered for %d", r.activeCalloutID)
	}

	key := proxywasm.RawBytePtrToString(keyData, keySize)

	var hs [][2]string
	switch mapType {
	case types.MapTypeHttpCallResponseHeaders:
		hs = res.headers
	case types.MapTypeHttpCallResponseTrailers:
		hs = res.trailers
	default:
		panic("unimplemented")
	}

	for _, h := range hs {
		if h[0] == key {
			v := []byte(h[1])
			*returnValueData = &v[0]
			*returnValueSize = len(v)
			return types.StatusOK
		}
	}

	return types.StatusNotFound
}

// // impl rawhostcall.ProxyWASMHost: delegated from hostEmulator
func (r *rootHostEmulator) rootHostEmulatorProxyGetBufferBytes(bt types.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) types.Status {
	var buf []byte
	switch bt {
	case types.BufferTypePluginConfiguration:
		buf = r.pluginConfiguration
	case types.BufferTypeVMConfiguration:
		buf = r.vmConfiguration
	case types.BufferTypeHttpCallResponseBody:
		activeID := proxywasm.VMStateGetActiveContextID()
		res, ok := r.httpCalloutResponse[r.activeCalloutID]
		if !ok {
			log.Fatalf("callout response unregistered for %d", activeID)
		}
		buf = res.body
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	if start >= len(buf) {
		log.Printf("start index out of range: %d (start) >= %d ", start, len(buf))
		return types.StatusBadArgument
	}

	*returnBufferData = &buf[start]
	if maxSize > len(buf)-start {
		*returnBufferSize = len(buf) - start
	} else {
		*returnBufferSize = maxSize
	}
	return types.StatusOK
}

// impl HostEmulator
func (r *rootHostEmulator) GetLogs(level types.LogLevel) []string {
	if level >= types.LogLevelMax {
		log.Fatalf("invalid log level: %d", level)
	}
	return r.logs[level]
}

// impl HostEmulator
func (r *rootHostEmulator) GetTickPeriod() uint32 {
	return r.tickPeriod
}

// impl HostEmulator
func (r *rootHostEmulator) Tick() {
	proxywasm.ProxyOnTick(RootContextID)
}

// impl HostEmulator
func (r *rootHostEmulator) GetQueueSize(queueID uint32) int {
	return len(r.queues[queueID])
}

// impl HostEmulator
func (r *rootHostEmulator) GetCalloutAttributesFromContext(contextID uint32) []HttpCalloutAttribute {
	infos := r.httpContextIDToCalloutInfos[contextID]
	return infos
}

// impl HostEmulator
func (r *rootHostEmulator) StartVM() {
	proxywasm.ProxyOnVMStart(RootContextID, len(r.vmConfiguration))
}

// impl HostEmulator
func (r *rootHostEmulator) StartPlugin() {
	proxywasm.ProxyOnConfigure(RootContextID, len(r.pluginConfiguration))
}

// impl HostEmulator
func (r *rootHostEmulator) PutCalloutResponse(calloutID uint32, headers, trailers [][2]string, body []byte) {
	r.httpCalloutResponse[calloutID] = struct {
		headers, trailers [][2]string
		body              []byte
	}{headers: headers, trailers: trailers, body: body}

	// RootContextID, calloutID uint32, numHeaders, bodySize, numTrailers in
	r.activeCalloutID = calloutID
	defer func() {
		r.activeCalloutID = 0
		delete(r.httpCalloutResponse, calloutID)
		delete(r.httpCalloutIDToContextID, calloutID)
	}()
	proxywasm.ProxyOnHttpCallResponse(RootContextID, calloutID, len(headers), len(body), len(trailers))
}

// impl HostEmulator
func (r *rootHostEmulator) FinishVM() {
	proxywasm.ProxyOnDone(RootContextID)
}

func (r *rootHostEmulator) GetCounterMetric(name string) (uint64, error) {
	id, ok := r.metricNameToID[name]
	if !ok {
		return 0, fmt.Errorf("%s not found", name)
	}

	t, ok := r.metricIDToType[id]
	if !ok {
		return 0, fmt.Errorf("%s not found", name)
	}

	if t != types.MetricTypeCounter {
		return 0, fmt.Errorf(
			"%s is not %v metric type but %v", name, types.MetricTypeCounter, t)
	}

	v, ok := r.metricIDToValue[id]
	if !ok {
		return 0, fmt.Errorf("%s not found", name)
	}
	return v, nil
}

func (r *rootHostEmulator) GetGaugeMetric(name string) (uint64, error) {
	id, ok := r.metricNameToID[name]
	if !ok {
		return 0, fmt.Errorf("%s not found", name)
	}

	t, ok := r.metricIDToType[id]
	if !ok {
		return 0, fmt.Errorf("%s not found", name)
	}

	if t != types.MetricTypeGauge {
		return 0, fmt.Errorf(
			"%s is not %v metric type but %v", name, types.MetricTypeGauge, t)
	}

	v, ok := r.metricIDToValue[id]
	if !ok {
		return 0, fmt.Errorf("%s not found", name)
	}
	return v, nil
}

func (r *rootHostEmulator) GetHistogramMetric(name string) (uint64, error) {
	id, ok := r.metricNameToID[name]
	if !ok {
		return 0, fmt.Errorf("%s not found", name)
	}

	t, ok := r.metricIDToType[id]
	if !ok {
		return 0, fmt.Errorf("%s not found", name)
	}

	if t != types.MetricTypeHistogram {
		return 0, fmt.Errorf(
			"%s is not %v metric type but %v", name, types.MetricTypeHistogram, t)
	}

	v, ok := r.metricIDToValue[id]
	if !ok {
		return 0, fmt.Errorf("%s not found", name)
	}
	return v, nil
}
