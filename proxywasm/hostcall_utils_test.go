package proxywasm

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func Test_stringBytePtr(t *testing.T) {
	exp := "abcd"
	ptr := stringBytePtr(exp)

	actual := *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(ptr)),
		Len:  len(exp),
		Cap:  len(exp),
	}))
	assert.Equal(t, exp, actual)
}
