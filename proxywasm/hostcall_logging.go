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

package proxywasm

import (
	"fmt"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
)

func LogTrace(msg string) {
	internal.ProxyLog(internal.LogLevelTrace, internal.StringBytePtr(msg), len(msg))
}

func LogTracef(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	internal.ProxyLog(internal.LogLevelTrace, internal.StringBytePtr(msg), len(msg))
}

func LogDebug(msg string) {
	internal.ProxyLog(internal.LogLevelDebug, internal.StringBytePtr(msg), len(msg))
}

func LogDebugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	internal.ProxyLog(internal.LogLevelDebug, internal.StringBytePtr(msg), len(msg))
}

func LogInfo(msg string) {
	internal.ProxyLog(internal.LogLevelInfo, internal.StringBytePtr(msg), len(msg))
}

func LogInfof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	internal.ProxyLog(internal.LogLevelInfo, internal.StringBytePtr(msg), len(msg))
}

func LogWarn(msg string) {
	internal.ProxyLog(internal.LogLevelWarn, internal.StringBytePtr(msg), len(msg))
}

func LogWarnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	internal.ProxyLog(internal.LogLevelWarn, internal.StringBytePtr(msg), len(msg))
}

func LogError(msg string) {
	internal.ProxyLog(internal.LogLevelError, internal.StringBytePtr(msg), len(msg))
}

func LogErrorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	internal.ProxyLog(internal.LogLevelError, internal.StringBytePtr(msg), len(msg))
}

func LogCritical(msg string) {
	internal.ProxyLog(internal.LogLevelCritical, internal.StringBytePtr(msg), len(msg))
}

func LogCriticalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	internal.ProxyLog(internal.LogLevelCritical, internal.StringBytePtr(msg), len(msg))
}
