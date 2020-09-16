package proxytest

import (
	"log"
	"reflect"
	"sync"
	"unsafe"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/rawhostcall"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

var hostMux = sync.Mutex{}

type baseHost struct {
	rawhostcall.DefaultProxyWAMSHost
	logs       [types.LogLevelMax][]string
	tickPeriod uint32

	queues      map[uint32][][]byte
	queueNameID map[string]uint32

	sharedDataKVS map[string]*sharedData
}

type sharedData struct {
	data []byte
	cas  uint32
}

func newBaseHost() *baseHost {
	return &baseHost{
		queues:        map[uint32][][]byte{},
		queueNameID:   map[string]uint32{},
		sharedDataKVS: map[string]*sharedData{},
	}
}

func (b *baseHost) ProxyLog(logLevel types.LogLevel, messageData *byte, messageSize int) types.Status {
	str := *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(messageData)),
		Len:  messageSize,
		Cap:  messageSize,
	}))

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

func (b *baseHost) getBuffer(bt types.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) types.Status {

	// should implement http callout response
	panic("unimplemented")
}

func (b *baseHost) ProxySetTickPeriodMilliseconds(period uint32) types.Status {
	b.tickPeriod = period
	return types.StatusOK
}

func (b *baseHost) GetTickPeriod() uint32 {
	return b.tickPeriod
}

// TODO: implement http callouts, metrics

func (b *baseHost) ProxyRegisterSharedQueue(nameData *byte, nameSize int, returnID *uint32) types.Status {
	name := *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(nameData)),
		Len:  nameSize,
		Cap:  nameSize,
	}))

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

	b.queues[queueID] = append(queue, *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(valueData)),
		Len:  valueSize,
		Cap:  valueSize,
	})))

	// TODO: should call OnQueueReady?

	return types.StatusOK
}

func (b *baseHost) GetQueueSize(queueID uint32) int {
	return len(b.queues[queueID])
}

func (b *baseHost) ProxyGetSharedData(keyData *byte, keySize int,
	returnValueData **byte, returnValueSize *int, returnCas *uint32) types.Status {
	key := *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(keyData)),
		Len:  keySize,
		Cap:  keySize,
	}))

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
	key := *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(keyData)),
		Len:  keySize,
		Cap:  keySize,
	}))
	value := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(valueData)),
		Len:  valueSize,
		Cap:  valueSize,
	}))

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
