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
	"strings"

	"github.com/mathetake/proxy-wasm-go/runtime/rawhostcall"
	"github.com/mathetake/proxy-wasm-go/runtime/types"
)

func LogTrace(msgs ...string) {
	msg := strings.Join(msgs, "")
	rawhostcall.ProxyLog(types.LogLevelTrace, unsafeGetStringBytePtr(msg), len(msg))
}

func LogDebug(msgs ...string) {
	msg := strings.Join(msgs, "")
	rawhostcall.ProxyLog(types.LogLevelDebug, unsafeGetStringBytePtr(msg), len(msg))
}

func LogInfo(msgs ...string) {
	msg := strings.Join(msgs, "")
	rawhostcall.ProxyLog(types.LogLevelInfo, unsafeGetStringBytePtr(msg), len(msg))
}

func LogWarn(msgs ...string) {
	msg := strings.Join(msgs, "")
	rawhostcall.ProxyLog(types.LogLevelWarn, unsafeGetStringBytePtr(msg), len(msg))
}

func LogError(msgs ...string) {
	msg := strings.Join(msgs, "")
	rawhostcall.ProxyLog(types.LogLevelError, unsafeGetStringBytePtr(msg), len(msg))
}

func LogCritical(msgs ...string) {
	msg := strings.Join(msgs, "")
	rawhostcall.ProxyLog(types.LogLevelCritical, unsafeGetStringBytePtr(msg), len(msg))
}
