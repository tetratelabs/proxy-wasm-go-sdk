package runtime

import (
	"reflect"
	"unsafe"
)

func unsafeGetStringBytePtr(msg string) *byte {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&msg))

	// TODO: seems redundant and we should use sliceHeader.Data directly (probably possible)
	bt := *(*[]byte)(unsafe.Pointer(&reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}))
	return &bt[0]
}
