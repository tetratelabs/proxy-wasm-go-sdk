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
}

func newBaseHost() *baseHost {
	return &baseHost{
		queues:      map[uint32][][]byte{},
		queueNameID: map[string]uint32{},
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

// TODO: implement http callouts, metrics, shared data

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
