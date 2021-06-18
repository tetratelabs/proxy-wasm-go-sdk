// Copyright 2021 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxywasm

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
)

type logHost struct {
	internal.DefaultProxyWAMSHost
	t           *testing.T
	expMessage  string
	expLogLevel internal.LogLevel
}

func (l logHost) ProxyLog(logLevel internal.LogLevel, messageData *byte, messageSize int) internal.Status {
	actual := internal.RawBytePtrToString(messageData, messageSize)
	require.Equal(l.t, l.expMessage, actual)
	require.Equal(l.t, l.expLogLevel, logLevel)
	return internal.StatusOK
}

func TestHostCall_Logging(t *testing.T) {
	t.Run("trace", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "trace",
			expLogLevel:          internal.LogLevelTrace,
		})
		defer release()
		LogTrace("trace")
	})

	t.Run("tracef", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "trace: log",
			expLogLevel:          internal.LogLevelTrace,
		})
		defer release()
		LogTracef("trace: %s", "log")
	})

	t.Run("debug", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "abc",
			expLogLevel:          internal.LogLevelDebug,
		})
		defer release()
		LogDebug("abc")
	})

	t.Run("debugf", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "debug: log",
			expLogLevel:          internal.LogLevelDebug,
		})
		defer release()
		LogDebugf("debug: %s", "log")
	})

	t.Run("info", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "info",
			expLogLevel:          internal.LogLevelInfo,
		})
		defer release()
		LogInfo("info")
	})

	t.Run("infof", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "info: log: 10",
			expLogLevel:          internal.LogLevelInfo,
		})
		defer release()
		LogInfof("info: %s: %d", "log", 10)
	})

	t.Run("warn", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "warn",
			expLogLevel:          internal.LogLevelWarn,
		})
		defer release()
		LogWarn("warn")
	})

	t.Run("warnf", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "warn: log: 10",
			expLogLevel:          internal.LogLevelWarn,
		})
		defer release()
		LogWarnf("warn: %s: %d", "log", 10)
	})

	t.Run("error", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "error",
			expLogLevel:          internal.LogLevelError,
		})
		defer release()
		LogError("error")
	})

	t.Run("warnf", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "warn: log: 10",
			expLogLevel:          internal.LogLevelWarn,
		})
		defer release()
		LogWarnf("warn: %s: %d", "log", 10)
	})

	t.Run("critical", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "critical error",
			expLogLevel:          internal.LogLevelCritical,
		})
		defer release()
		LogCritical("critical error")
	})

	t.Run("criticalf", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "critical: log: 10",
			expLogLevel:          internal.LogLevelCritical,
		})
		defer release()
		LogCriticalf("critical: %s: %d", "log", 10)
	})
}
