package proxywasm

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func SetNewRootContextFn(f func(contextID uint32) types.RootContext) {
	internal.SetNewRootContextFn(f)
}
