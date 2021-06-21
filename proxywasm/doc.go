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

Contexts and Wasm VM

There are two main packages in this SDK. The one is this "proxywasm" package and "types" package under /types subdirectory.
proxywasm package depends on types package, and the types package contains the intefaces you are supposed to implement
in order to extend your network proxies.

In types package, there are four types of these intefaces which we call "contexts".
They are VMContext, PluginContext, TcpContext and HttpContext, and their relationship can be described as the following diagram:

                         Wasm Virtual Machine(VM)
                    (corresponds to VM configuration)
 ┌────────────────────────────────────────────────────────────────────────────┐
 │                                                      TcpContext            │
 │                                                  ╱ (Each Tcp stream)       │
 │                                                 ╱                          │
 │                      1: N                      ╱ 1: N                      │
 │       VMContext  ──────────  PluginContext                                 │
 │  (VM configuration)     (Plugin configuration) ╲ 1: N                      │
 │                                                 ╲                          │
 │                                                  ╲   HttpContext           │
 │                                                   (Each Http stream)       │
 └────────────────────────────────────────────────────────────────────────────┘

To summarize,

1) VMContext corresponds to each Wasm Virtual Machine, and only one VMContext exists in each VM.
Note that in Envoy, Wasm VMs are created per "vm_config" field in envoy.yaml. For example having different "vm_config.configuration" fields
results in multiple VMs being created and each of them corresponds to each "vm_config.configuration".

2) VMContext is parent of PluginContext, and is responsible for creating arbitrary number of PluginContexts.

3) PluginContext corresponds to each plugin configurations in the host. In Envoy, each plugin configuration is given at HttpFilter or NetworkFilter
on listeners. That said, a plugin context corresponds to a Http or Network filter on a litener and is in charge of creating "filter instances" for
each Http or Tcp streams. And these "filter instances" are HttpContexts or TcpContexts.

4) PluginContext is parent of TcpContext and HttpContexts, and is responsible for creating arbitrary number of these contexts.

5) TcpContext is responsible for handling each Tcp stream events.

6) HttpContext is responsible for handling each Http stream events.

Please refer to types package's documentation for the detail of interfaces.

Entrypoint

You must call "proxywasm.SetVMContext" in "main()" function so that hosts can call OnVMStart and
the program becomes ready for creating PluginContexts. E.g.

 func main() {
 	proxywasm.SetVMContext(&myVMContext{})
 }

 type myVMContext struct {
	 ...
 }

 // My implementations of types.VMContext interface..

Test Framework

You can test your Proxy-Wasm program with the "proxytest" build tag without actually running proxies and Wasm VMs.
"proxytest" package is under /proxytest subdirectory. It contains the host emulators which enables you to compile
and run "go test" as usual on your native CPU.

Examples

There are tens of examples here(https://github.com/tetratelabs/proxy-wasm-go-sdk/tree/main/examples),
so please take a look there and get the sense of how to implement your own plugins!
*/
package proxywasm
