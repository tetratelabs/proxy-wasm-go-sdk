package proxywasm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mathetake/proxy-wasm-go/proxywasm/rawhostcall"
	"github.com/mathetake/proxy-wasm-go/proxywasm/types"
)

type logHost struct {
	rawhostcall.DefaultProxyWAMSHost
	t           *testing.T
	expMessage  string
	expLogLevel types.LogLevel
}

func (l logHost) ProxyLog(logLevel types.LogLevel, messageData *byte, messageSize int) types.Status {
	actual := rawBytePtrToString(messageData, messageSize)
	assert.Equal(l.t, l.expMessage, actual)
	assert.Equal(l.t, l.expLogLevel, logLevel)
	return types.StatusOK
}

func TestHostCall_Logging(t *testing.T) {
	hostMutex.Lock()
	defer hostMutex.Unlock()

	t.Run("trace", func(t *testing.T) {
		rawhostcall.RegisterMockWASMHost(logHost{
			DefaultProxyWAMSHost: rawhostcall.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "trace",
			expLogLevel:          types.LogLevelTrace,
		})
		LogTrace("trace")
	})

	t.Run("debug", func(t *testing.T) {
		rawhostcall.RegisterMockWASMHost(logHost{
			DefaultProxyWAMSHost: rawhostcall.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "abc",
			expLogLevel:          types.LogLevelDebug,
		})
		LogDebug("a", "b", "c")
	})

	t.Run("info", func(t *testing.T) {
		rawhostcall.RegisterMockWASMHost(logHost{
			DefaultProxyWAMSHost: rawhostcall.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "info",
			expLogLevel:          types.LogLevelInfo,
		})
		LogInfo("in", "f", "o")
	})

	t.Run("warn", func(t *testing.T) {
		rawhostcall.RegisterMockWASMHost(logHost{
			DefaultProxyWAMSHost: rawhostcall.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "warn",
			expLogLevel:          types.LogLevelWarn,
		})
		LogWarn("warn")
	})

	t.Run("error", func(t *testing.T) {
		rawhostcall.RegisterMockWASMHost(logHost{
			DefaultProxyWAMSHost: rawhostcall.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "error",
			expLogLevel:          types.LogLevelError,
		})
		LogError("error")
	})

	t.Run("critical", func(t *testing.T) {
		rawhostcall.RegisterMockWASMHost(logHost{
			DefaultProxyWAMSHost: rawhostcall.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "critical error",
			expLogLevel:          types.LogLevelCritical,
		})
		LogCritical("critical", " error")
	})
}
