package runtime

import (
	"reflect"
	"unsafe"
)

// any better way to string formatting without system calls....?

func LogTrace(msg string) {
	proxyLog(LogLevelTrace, unsafeGetStringBytePtr(msg), len(msg))
}

func LogDebug(msg string) {
	proxyLog(LogLevelDebug, unsafeGetStringBytePtr(msg), len(msg))
}

func LogInfo(msg string) {
	proxyLog(LogLevelInfo, unsafeGetStringBytePtr(msg), len(msg))
}

func LogWarn(msg string) {
	proxyLog(LogLevelWarn, unsafeGetStringBytePtr(msg), len(msg))
}

func LogError(msg string) {
	proxyLog(LogLevelWarn, unsafeGetStringBytePtr(msg), len(msg))
}

func LogCritical(msg string) {
	proxyLog(LogLevelWarn, unsafeGetStringBytePtr(msg), len(msg))
}

func unsafeGetStringBytePtr(msg string) *byte {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&msg))
	bt := *(*[]byte)(unsafe.Pointer(&reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}))
	return &bt[0]
}
