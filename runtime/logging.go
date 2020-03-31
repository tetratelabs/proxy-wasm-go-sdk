// Copyright 2020 Takeshi Yoneda(@mathetake)
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
