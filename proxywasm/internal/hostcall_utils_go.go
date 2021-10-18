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

//go:build proxytest

// Since the difference of the types in SliceHeader.{Len, Cap} between tinygo and go,
// we have to have separated functions for converting bytes
// https://github.com/tinygo-org/tinygo/issues/1284

package internal

import (
	"reflect"
	"unsafe"
)

func RawBytePtrToString(raw *byte, size int) string {
	//nolint
	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(raw)),
		Len:  size,
		Cap:  size,
	}))
}

func RawBytePtrToByteSlice(raw *byte, size int) []byte {
	//nolint
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(raw)),
		Len:  size,
		Cap:  size,
	}))
}
