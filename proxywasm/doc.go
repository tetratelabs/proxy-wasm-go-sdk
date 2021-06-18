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

/*
Proxy-Wasm Go SDK is for developers who wants to write an Proxy-Wasm plugin in Go.
Proxy-Wasm(https://github.com/proxy-wasm) is the project for extending network proxies with WebAssembly modules,
and this SDK can be used for building Proxy-Wasm compatible Wasm binaries easily without knowing its
technical detail and specification.
This SDK leverages TinyGo(https://tinygo.org) and cannnot be used with the official Go compiler toolchain.

Please visit Proxy-Wasm Spec(https://github.com/proxy-wasm), C++ SDK(https://github.com/proxy-wasm/proxy-wasm-cpp-sdk),
Rust SDK(https://github.com/proxy-wasm/proxy-wasm-rust-sdk) for further references.

Overview

There are two main packages in this SDK. The one is this "proxywasm" package and "types" package under /types subdirectory.
proxywasm package depends on types package, and the types package contains the intefaces you are supposed to implement
in order to extend your network proxies.

In types package, there are three types of these intefaces which we call "contexts".
They are called RootContext, TcpContext and HttpContext, and their relationship can be described as the following diagram:

                                              ╱ TcpContext = handling each Tcp stream
                                             ╱
                                            ╱ 1: N
   each plugin configuration = RootContext
                                            ╲ 1: N
                                             ╲
                                              ╲ Http = handling each Http stream

In other words, RootContex is the parent of others, and responsible for creating Tcp and Http contexts
corresponding to each streams if it is configured for running as a Http/Tcp stream plugin.
Given that, RootContext is the primary interface everyone has to implement.

Here we "plugin configuration" means, for example, "http filter configuration" in Envoy proxy's terminology.
That means the same Wasm VM can be run at multiple "http filter configuration" (e.g. multiple http listeners),
and for each configuration, an user implemented instance of RootContext inteface is created in side the Wasm VM and used for
creating corresponding stream contexts.

Please refer to types package's documentation for the detail of interfaces.

Entrypoint

You must call "proxywasm.SetNewRootContextFn" in "main()" function so that the hosts (e.g. Envoy)
can initialize your instances of RootContexts for each plugin configurations. E.g.

 func main() {
 	proxywasm.SetNewRootContextFn(newRootContext)
 }

 func newRootContext(uint32) types.RootContext { return &myRootContext{ ... } }

 type myRootContext struct {
	 ...
 }

 // My implementations of types.RootContext interface..

Test Framework

You can test your Proxy-Wasm program with the "proxytest" build tag without actually running proxies and Wasm VMs.
"proxytest" package is under /proxytest subdirectory. It contains the host emulators which enables you to compile
and run "go test" as usual on your native CPU.

Examples

There are tens of examples here(https://github.com/tetratelabs/proxy-wasm-go-sdk/tree/main/examples),
so please take a look there and get the sense of how to implement your own plugins!
*/
package proxywasm
