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
