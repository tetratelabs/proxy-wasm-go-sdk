package proxywasm

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

// SetNewRootContextFn is the entrypoint for setting up this entire Wasm VM.
// The given function is responsible for creating types.RootContext and is called when
// the host initializes a plugin for each plugin configuration.
// Please make sure that this entrypoint be called during "main()" function, otherwise
// this VM would fail.
func SetNewRootContextFn(f func(contextID uint32) types.RootContext) {
	internal.SetNewRootContextFn(f)
}
