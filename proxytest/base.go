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

package proxytest

import (
	"log"
	"sync"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/rawhostcall"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

var hostMux = sync.Mutex{}

type baseHost struct {
	rawhostcall.DefaultProxyWAMSHost
	currentContextID uint32

	logs       [types.LogLevelMax][]string
	tickPeriod uint32

	queues      map[uint32][][]byte
	queueNameID map[string]uint32

	sharedDataKVS map[string]*sharedData

	metricIDToValue map[uint32]uint64
	metricIDToType  map[uint32]types.MetricType
	metricNameToID  map[string]uint32

	calloutCallbackCaller func(contextID uint32, numHeaders, bodySize, numTrailers int)
	calloutResponse       map[uint32]struct {
		headers, trailers [][2]string
		body              []byte
	}
	callouts map[uint32]struct{}
}

type sharedData struct {
	data []byte
	cas  uint32
}

func newBaseHost(f func(contextID uint32, numHeaders, bodySize, numTrailers int)) *baseHost {
	return &baseHost{
		queues:                map[uint32][][]byte{},
		queueNameID:           map[string]uint32{},
		sharedDataKVS:         map[string]*sharedData{},
		metricIDToValue:       map[uint32]uint64{},
		metricIDToType:        map[uint32]types.MetricType{},
		metricNameToID:        map[string]uint32{},
		calloutCallbackCaller: f,
		calloutResponse: map[uint32]struct {
			headers, trailers [][2]string
			body              []byte
		}{},
		callouts: map[uint32]struct{}{},
	}
}

func (b *baseHost) ProxyLog(logLevel types.LogLevel, messageData *byte, messageSize int) types.Status {
	str := proxywasm.RawBytePtrToString(messageData, messageSize)

	log.Printf("proxy_log: %s", str)
	// TODO: exit if loglevel == fatal?

	b.logs[logLevel] = append(b.logs[logLevel], str)
	return types.StatusOK
}

func (b *baseHost) GetLogs(level types.LogLevel) []string {
	if level >= types.LogLevelMax {
		log.Fatalf("invalid log level: %d", level)
	}
	return b.logs[level]
}

func (b *baseHost) ProxySetTickPeriodMilliseconds(period uint32) types.Status {
	b.tickPeriod = period
	return types.StatusOK
}

func (b *baseHost) GetTickPeriod() uint32 {
	return b.tickPeriod
}

func (b *baseHost) ProxyRegisterSharedQueue(nameData *byte, nameSize int, returnID *uint32) types.Status {
	name := proxywasm.RawBytePtrToString(nameData, nameSize)
	if id, ok := b.queueNameID[name]; ok {
		*returnID = id
		return types.StatusOK
	}

	id := uint32(len(b.queues))
	b.queues[id] = [][]byte{}
	b.queueNameID[name] = id
	*returnID = id
	return types.StatusOK
}

func (b *baseHost) ProxyDequeueSharedQueue(queueID uint32, returnValueData **byte, returnValueSize *int) types.Status {
	queue, ok := b.queues[queueID]
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
	b.queues[queueID] = queue[1:]
	return types.StatusOK
}

func (b *baseHost) ProxyEnqueueSharedQueue(queueID uint32, valueData *byte, valueSize int) types.Status {
	queue, ok := b.queues[queueID]
	if !ok {
		log.Printf("queue %d is not found", queueID)
		return types.StatusNotFound
	}

	b.queues[queueID] = append(queue, proxywasm.RawBytePtrToByteSlice(valueData, valueSize))

	// TODO: should call OnQueueReady?

	return types.StatusOK
}

func (b *baseHost) GetQueueSize(queueID uint32) int {
	return len(b.queues[queueID])
}

func (b *baseHost) ProxyGetSharedData(keyData *byte, keySize int,
	returnValueData **byte, returnValueSize *int, returnCas *uint32) types.Status {
	key := proxywasm.RawBytePtrToString(keyData, keySize)

	value, ok := b.sharedDataKVS[key]
	if !ok {
		return types.StatusNotFound
	}

	*returnValueSize = len(value.data)
	*returnValueData = &value.data[0]
	*returnCas = value.cas
	return types.StatusOK
}

func (b *baseHost) ProxySetSharedData(keyData *byte, keySize int,
	valueData *byte, valueSize int, cas uint32) types.Status {
	key := proxywasm.RawBytePtrToString(keyData, keySize)
	value := proxywasm.RawBytePtrToByteSlice(valueData, valueSize)

	prev, ok := b.sharedDataKVS[key]
	if !ok {
		b.sharedDataKVS[key] = &sharedData{
			data: value,
			cas:  cas + 1,
		}
		return types.StatusOK
	}

	if prev.cas != cas {
		return types.StatusCasMismatch
	}

	b.sharedDataKVS[key].cas = cas + 1
	b.sharedDataKVS[key].data = value
	return types.StatusOK
}

func (b *baseHost) ProxyDefineMetric(metricType types.MetricType,
	metricNameData *byte, metricNameSize int, returnMetricIDPtr *uint32) types.Status {
	name := proxywasm.RawBytePtrToString(metricNameData, metricNameSize)
	id, ok := b.metricNameToID[name]
	if !ok {
		id = uint32(len(b.metricNameToID))
		b.metricNameToID[name] = id
		b.metricIDToValue[id] = 0
		b.metricIDToType[id] = metricType
	}
	*returnMetricIDPtr = id
	return types.StatusOK
}

func (b *baseHost) ProxyIncrementMetric(metricID uint32, offset int64) types.Status {
	// TODO: check metric type

	val, ok := b.metricIDToValue[metricID]
	if !ok {
		return types.StatusBadArgument
	}

	b.metricIDToValue[metricID] = val + uint64(offset)
	return types.StatusOK
}

func (b *baseHost) ProxyRecordMetric(metricID uint32, value uint64) types.Status {
	// TODO: check metric type

	_, ok := b.metricIDToValue[metricID]
	if !ok {
		return types.StatusBadArgument
	}
	b.metricIDToValue[metricID] = value
	return types.StatusOK
}

func (b *baseHost) ProxyGetMetric(metricID uint32, returnMetricValue *uint64) types.Status {
	value, ok := b.metricIDToValue[metricID]
	if !ok {
		return types.StatusBadArgument
	}
	*returnMetricValue = value
	return types.StatusOK
}

func (b *baseHost) getBuffer(bt types.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) types.Status {
	if bt != types.BufferTypeHttpCallResponseBody {
		panic("unimplemented")
	}

	res, ok := b.calloutResponse[b.currentContextID]
	if !ok {
		log.Fatalf("callout response unregistered for %d", b.currentContextID)
	}

	*returnBufferData = &res.body[0]
	*returnBufferSize = len(res.body)
	return types.StatusOK
}

func (b *baseHost) getMapValue(mapType types.MapType, keyData *byte,
	keySize int, returnValueData **byte, returnValueSize *int) types.Status {
	res, ok := b.calloutResponse[b.currentContextID]
	if !ok {
		log.Fatalf("callout response unregistered for %d", b.currentContextID)
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

func (b *baseHost) ProxyHttpCall(upstreamData *byte, upstreamSize int, headerData *byte, headerSize int, bodyData *byte,
	bodySize int, trailersData *byte, trailersSize int, timeout uint32, _ *uint32) types.Status {
	upstream := proxywasm.RawBytePtrToString(upstreamData, upstreamSize)
	body := proxywasm.RawBytePtrToString(bodyData, bodySize)
	headers := proxywasm.DeserializeMap(proxywasm.RawBytePtrToByteSlice(headerData, headerSize))
	trailers := proxywasm.DeserializeMap(proxywasm.RawBytePtrToByteSlice(trailersData, trailersSize))

	log.Printf("[http callout to %s] timeout: %d", upstream, timeout)
	log.Printf("[http callout to %s] headers: %v", upstream, headers)
	log.Printf("[http callout to %s] body: %s", upstream, body)
	log.Printf("[http callout to %s] trailers: %v", upstream, trailers)

	b.callouts[b.currentContextID] = struct{}{}
	return types.StatusOK
}

func (b *baseHost) PutCalloutResponse(contextID uint32, headers, trailers [][2]string, body []byte) {
	b.calloutResponse[contextID] = struct {
		headers, trailers [][2]string
		body              []byte
	}{headers: headers, trailers: trailers, body: body}

	b.currentContextID = contextID
	b.calloutCallbackCaller(contextID, len(headers), len(body), len(trailers))
	delete(b.calloutResponse, contextID)
}

func (b *baseHost) IsDispatchCalled(contextID uint32) bool {
	_, ok := b.callouts[contextID]
	return ok
}

func (b *baseHost) getMapPairs(mapType types.MapType, returnValueData **byte, returnValueSize *int) types.Status {
	res, ok := b.calloutResponse[b.currentContextID]
	if !ok {
		log.Fatalf("callout response unregistered for %d", b.currentContextID)
	}

	var raw []byte
	switch mapType {
	case types.MapTypeHttpCallResponseHeaders:
		raw = proxywasm.SerializeMap(res.headers)
	case types.MapTypeHttpCallResponseTrailers:
		raw = proxywasm.SerializeMap(res.trailers)
	default:
		panic("unimplemented")
	}

	*returnValueData = &raw[0]
	*returnValueSize = len(raw)
	return types.StatusOK
}
