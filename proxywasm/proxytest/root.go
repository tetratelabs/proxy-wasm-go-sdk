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
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type (
	rootHostEmulator struct {
		activeCalloutID  uint32
		logs             [internal.LogLevelMax][]string
		tickPeriod       uint32
		foreignFunctions map[string]func([]byte) []byte

		queues        map[uint32][][]byte
		queueNameID   map[string]uint32
		sharedDataKVS map[string]*sharedData

		httpContextIDToCalloutInfos map[uint32][]HttpCalloutAttribute // key: contextID
		httpCalloutIDToContextID    map[uint32]uint32                 // key: calloutID
		httpCalloutResponse         map[uint32]struct {               // key: calloutID
			headers  [][2]string
			trailers [][2]string
			body     []byte
		}

		metricIDToType  map[uint32]internal.MetricType
		metricNameToID  map[string]uint32
		metricIDToValue map[uint32]uint64

		pluginConfiguration, vmConfiguration []byte
	}

	HttpCalloutAttribute struct {
		CalloutID uint32
		Upstream  string
		Headers   [][2]string
		Trailers  [][2]string
		Body      []byte
	}

	sharedData struct {
		data []byte
		cas  uint32
	}
)

func newRootHostEmulator(pluginConfiguration, vmConfiguration []byte) *rootHostEmulator {
	host := &rootHostEmulator{
		foreignFunctions:            map[string]func([]byte) []byte{},
		queues:                      map[uint32][][]byte{},
		queueNameID:                 map[string]uint32{},
		sharedDataKVS:               map[string]*sharedData{},
		metricIDToValue:             map[uint32]uint64{},
		metricIDToType:              map[uint32]internal.MetricType{},
		metricNameToID:              map[string]uint32{},
		httpContextIDToCalloutInfos: map[uint32][]HttpCalloutAttribute{},
		httpCalloutIDToContextID:    map[uint32]uint32{},
		httpCalloutResponse: map[uint32]struct {
			headers  [][2]string
			trailers [][2]string
			body     []byte
		}{},

		pluginConfiguration: pluginConfiguration,
		vmConfiguration:     vmConfiguration,
	}
	return host
}

// impl internal.ProxyWasmHost
func (r *rootHostEmulator) ProxyLog(logLevel internal.LogLevel, messageData *byte, messageSize int) internal.Status {
	str := internal.RawBytePtrToString(messageData, messageSize)

	log.Printf("proxy_%s_log: %s", logLevel, str)
	r.logs[logLevel] = append(r.logs[logLevel], str)
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (r *rootHostEmulator) ProxySetTickPeriodMilliseconds(period uint32) internal.Status {
	r.tickPeriod = period
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (r *rootHostEmulator) ProxyRegisterSharedQueue(nameData *byte, nameSize int, returnID *uint32) internal.Status {
	name := internal.RawBytePtrToString(nameData, nameSize)
	if id, ok := r.queueNameID[name]; ok {
		*returnID = id
		return internal.StatusOK
	}

	id := uint32(len(r.queues))
	r.queues[id] = [][]byte{}
	r.queueNameID[name] = id
	*returnID = id
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (r *rootHostEmulator) ProxyDequeueSharedQueue(queueID uint32, returnValueData **byte, returnValueSize *int) internal.Status {
	queue, ok := r.queues[queueID]
	if !ok {
		log.Printf("queue %d is not found", queueID)
		return internal.StatusNotFound
	} else if len(queue) == 0 {
		log.Printf("queue %d is empty", queueID)
		return internal.StatusEmpty
	}

	data := queue[0]
	*returnValueData = &data[0]
	*returnValueSize = len(data)
	r.queues[queueID] = queue[1:]
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (r *rootHostEmulator) ProxyEnqueueSharedQueue(queueID uint32, valueData *byte, valueSize int) internal.Status {
	queue, ok := r.queues[queueID]
	if !ok {
		log.Printf("queue %d is not found", queueID)
		return internal.StatusNotFound
	}

	r.queues[queueID] = append(queue, internal.RawBytePtrToByteSlice(valueData, valueSize))
	internal.ProxyOnQueueReady(PluginContextID, queueID)
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (r *rootHostEmulator) ProxyGetSharedData(keyData *byte, keySize int,
	returnValueData **byte, returnValueSize *int, returnCas *uint32) internal.Status {
	key := internal.RawBytePtrToString(keyData, keySize)

	value, ok := r.sharedDataKVS[key]
	if !ok {
		return internal.StatusNotFound
	}

	*returnValueSize = len(value.data)
	*returnValueData = &value.data[0]
	*returnCas = value.cas
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (r *rootHostEmulator) ProxySetSharedData(keyData *byte, keySize int,
	valueData *byte, valueSize int, cas uint32) internal.Status {
	// Copy data provided by plugin to keep ownership within host. Otherwise, when
	// plugin deallocates the memory could be modified.
	key := strings.Clone(internal.RawBytePtrToString(keyData, keySize))
	v := internal.RawBytePtrToByteSlice(valueData, valueSize)
	value := make([]byte, len(v))
	copy(value, v)

	prev, ok := r.sharedDataKVS[key]
	if !ok {
		r.sharedDataKVS[key] = &sharedData{
			data: value,
			cas:  cas + 1,
		}
		return internal.StatusOK
	}

	if prev.cas != cas {
		return internal.StatusCasMismatch
	}

	r.sharedDataKVS[key].cas = cas + 1
	r.sharedDataKVS[key].data = value
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (r *rootHostEmulator) ProxyDefineMetric(metricType internal.MetricType,
	metricNameData *byte, metricNameSize int, returnMetricIDPtr *uint32) internal.Status {
	name := internal.RawBytePtrToString(metricNameData, metricNameSize)
	id, ok := r.metricNameToID[name]
	if !ok {
		id = uint32(len(r.metricNameToID))
		r.metricNameToID[name] = id
		r.metricIDToValue[id] = 0
		r.metricIDToType[id] = metricType
	}
	*returnMetricIDPtr = id
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (r *rootHostEmulator) ProxyIncrementMetric(metricID uint32, offset int64) internal.Status {
	val, ok := r.metricIDToValue[metricID]
	if !ok {
		return internal.StatusBadArgument
	}

	r.metricIDToValue[metricID] = val + uint64(offset)
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (r *rootHostEmulator) ProxyRecordMetric(metricID uint32, value uint64) internal.Status {
	_, ok := r.metricIDToValue[metricID]
	if !ok {
		return internal.StatusBadArgument
	}
	r.metricIDToValue[metricID] = value
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (r *rootHostEmulator) ProxyGetMetric(metricID uint32, returnMetricValue *uint64) internal.Status {
	value, ok := r.metricIDToValue[metricID]
	if !ok {
		return internal.StatusBadArgument
	}
	*returnMetricValue = value
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (r *rootHostEmulator) ProxyHttpCall(upstreamData *byte, upstreamSize int, headerData *byte, headerSize int, bodyData *byte,
	bodySize int, trailersData *byte, trailersSize int, timeout uint32, calloutIDPtr *uint32) internal.Status {
	upstream := internal.RawBytePtrToString(upstreamData, upstreamSize)
	body := internal.RawBytePtrToString(bodyData, bodySize)
	headers := deserializeRawBytePtrToMap(headerData, headerSize)
	trailers := deserializeRawBytePtrToMap(trailersData, trailersSize)

	log.Printf("[http callout to %s] timeout: %d", upstream, timeout)
	log.Printf("[http callout to %s] headers: %v", upstream, headers)
	log.Printf("[http callout to %s] body: %s", upstream, body)
	log.Printf("[http callout to %s] trailers: %v", upstream, trailers)

	calloutID := uint32(len(r.httpCalloutIDToContextID))
	contextID := internal.VMStateGetActiveContextID()
	r.httpCalloutIDToContextID[calloutID] = contextID
	r.httpContextIDToCalloutInfos[contextID] = append(r.httpContextIDToCalloutInfos[contextID], HttpCalloutAttribute{
		CalloutID: calloutID,
		Upstream:  upstream,
		Headers:   headers,
		Trailers:  trailers,
		Body:      []byte(body),
	})

	*calloutIDPtr = calloutID
	return internal.StatusOK
}

// impl internal.ProxyWasmHost
func (r *rootHostEmulator) RegisterForeignFunction(name string, f func([]byte) []byte) {
	r.foreignFunctions[name] = f
}

// impl internal.ProxyWasmHost
func (r *rootHostEmulator) ProxyCallForeignFunction(funcNamePtr *byte, funcNameSize int, paramPtr *byte, paramSize int, returnData **byte, returnSize *int) internal.Status {
	funcName := internal.RawBytePtrToString(funcNamePtr, funcNameSize)
	param := internal.RawBytePtrToByteSlice(paramPtr, paramSize)

	log.Printf("[foreign call] funcname: %s", funcName)
	log.Printf("[foreign call] param: %s", param)

	f, ok := r.foreignFunctions[funcName]
	if !ok {
		log.Fatalf("%s not registered as a foreign function", funcName)
	}
	ret := f(param)
	*returnData = &ret[0]
	*returnSize = len(ret)

	return internal.StatusOK
}

// // impl internal.ProxyWasmHost: delegated from hostEmulator
func (r *rootHostEmulator) rootHostEmulatorProxyGetHeaderMapPairs(mapType internal.MapType, returnValueData **byte, returnValueSize *int) internal.Status {
	res, ok := r.httpCalloutResponse[r.activeCalloutID]
	if !ok {
		log.Fatalf("callout response unregistered for %d", r.activeCalloutID)
	}

	var raw []byte
	switch mapType {
	case internal.MapTypeHttpCallResponseHeaders:
		raw = internal.SerializeMap(res.headers)
	case internal.MapTypeHttpCallResponseTrailers:
		raw = internal.SerializeMap(res.trailers)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	*returnValueData = &raw[0]
	*returnValueSize = len(raw)
	return internal.StatusOK
}

// // impl internal.ProxyWasmHost: delegated from hostEmulator
func (r *rootHostEmulator) rootHostEmulatorProxyGetMapValue(mapType internal.MapType, keyData *byte,
	keySize int, returnValueData **byte, returnValueSize *int) internal.Status {
	res, ok := r.httpCalloutResponse[r.activeCalloutID]
	if !ok {
		log.Fatalf("callout response unregistered for %d", r.activeCalloutID)
	}

	var hs [][2]string
	switch mapType {
	case internal.MapTypeHttpCallResponseHeaders:
		hs = res.headers
	case internal.MapTypeHttpCallResponseTrailers:
		hs = res.trailers
	default:
		panic("unimplemented")
	}

	key := strings.ToLower(internal.RawBytePtrToString(keyData, keySize))

	for _, h := range hs {
		if h[0] == key {
			v := []byte(h[1])
			*returnValueData = &v[0]
			*returnValueSize = len(v)
			return internal.StatusOK
		}
	}

	return internal.StatusNotFound
}

// // impl internal.ProxyWasmHost: delegated from hostEmulator
func (r *rootHostEmulator) rootHostEmulatorProxyGetBufferBytes(bt internal.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) internal.Status {
	var buf []byte
	switch bt {
	case internal.BufferTypePluginConfiguration:
		buf = r.pluginConfiguration
	case internal.BufferTypeVMConfiguration:
		buf = r.vmConfiguration
	case internal.BufferTypeHttpCallResponseBody:
		activeID := internal.VMStateGetActiveContextID()
		res, ok := r.httpCalloutResponse[r.activeCalloutID]
		if !ok {
			log.Fatalf("callout response unregistered for %d", activeID)
		}
		buf = res.body
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	if len(buf) == 0 {
		return internal.StatusNotFound
	} else if start >= len(buf) {
		log.Printf("start index out of range: %d (start) >= %d ", start, len(buf))
		return internal.StatusBadArgument
	}

	*returnBufferData = &buf[start]
	if maxSize > len(buf)-start {
		*returnBufferSize = len(buf) - start
	} else {
		*returnBufferSize = maxSize
	}
	return internal.StatusOK
}

// impl HostEmulator
func (r *rootHostEmulator) GetTraceLogs() []string {
	return r.getLogs(internal.LogLevelTrace)
}

// impl HostEmulator
func (r *rootHostEmulator) GetDebugLogs() []string {
	return r.getLogs(internal.LogLevelDebug)
}

// impl HostEmulator
func (r *rootHostEmulator) GetInfoLogs() []string {
	return r.getLogs(internal.LogLevelInfo)
}

// impl HostEmulator
func (r *rootHostEmulator) GetWarnLogs() []string {
	return r.getLogs(internal.LogLevelWarn)
}

// impl HostEmulator
func (r *rootHostEmulator) GetErrorLogs() []string {
	return r.getLogs(internal.LogLevelError)
}

// impl HostEmulator
func (r *rootHostEmulator) GetCriticalLogs() []string {
	return r.getLogs(internal.LogLevelCritical)
}

func (r *rootHostEmulator) getLogs(level internal.LogLevel) []string {
	return r.logs[level]
}

// impl HostEmulator
func (r *rootHostEmulator) GetTickPeriod() uint32 {
	return r.tickPeriod
}

// impl HostEmulator
func (r *rootHostEmulator) Tick() {
	internal.ProxyOnTick(PluginContextID)
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
func (r *rootHostEmulator) StartVM() types.OnVMStartStatus {
	return internal.ProxyOnVMStart(PluginContextID, len(r.vmConfiguration))
}

// impl HostEmulator
func (r *rootHostEmulator) StartPlugin() types.OnPluginStartStatus {
	return internal.ProxyOnConfigure(PluginContextID, len(r.pluginConfiguration))
}

// impl HostEmulator
func (r *rootHostEmulator) CallOnHttpCallResponse(calloutID uint32, headers, trailers [][2]string, body []byte) {
	r.httpCalloutResponse[calloutID] = struct {
		headers, trailers [][2]string
		body              []byte
	}{headers: cloneWithLowerCaseMapKeys(headers), trailers: cloneWithLowerCaseMapKeys(trailers), body: body}

	// PluginContextID, calloutID uint32, numHeaders, bodySize, numTrailers in
	r.activeCalloutID = calloutID
	defer func() {
		r.activeCalloutID = 0
		delete(r.httpCalloutResponse, calloutID)
		delete(r.httpCalloutIDToContextID, calloutID)
	}()
	internal.ProxyOnHttpCallResponse(PluginContextID, calloutID, len(headers), len(body), len(trailers))
}

// impl HostEmulator
func (r *rootHostEmulator) FinishVM() bool {
	return internal.ProxyOnDone(PluginContextID)
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

	if t != internal.MetricTypeCounter {
		return 0, fmt.Errorf(
			"%s is not %v metric type but %v", name, internal.MetricTypeCounter, t)
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

	if t != internal.MetricTypeGauge {
		return 0, fmt.Errorf(
			"%s is not %v metric type but %v", name, internal.MetricTypeGauge, t)
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

	if t != internal.MetricTypeHistogram {
		return 0, fmt.Errorf(
			"%s is not %v metric type but %v", name, internal.MetricTypeHistogram, t)
	}

	v, ok := r.metricIDToValue[id]
	if !ok {
		return 0, fmt.Errorf("%s not found", name)
	}
	return v, nil
}
