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
	logs [types.LogLevelMax][]string
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

	// should implement http callout response, vm configuration, plugin configuration
	panic("unimplemented")
}

// TODO: implement http callouts, metrics, times, plugins, queue, shared data
