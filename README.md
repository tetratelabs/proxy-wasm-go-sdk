> [!WARNING] 
> We are no longer recommending this SDK or Wasm in general for anyone due to the fundamental memory issue of TinyGo (See [the detailed explanation](https://github.com/tetratelabs/proxy-wasm-go-sdk/issues/450#issuecomment-2253729297) by a long-time community member)
> as well as [the project state of Proxy-Wasm in general](https://github.com/envoyproxy/envoy/issues/35420).
> If you are not in a position where you have to run untrusted binaries (like for example, you run Envoy proxies while your client gives you the binaries to run), we recommend using other extension mechanism 
> such as Lua or External Processing which should be comparable or better or worse depending on the use case. 
> 
> If you are already using this SDK, but still want to continue using Wasm for some reason instead of Lua or External Processing, 
> we strongly recommend migrating to the Rust or C++ SDK due to the memory issue of TinyGo described in the link above.
> 
> We keep this repository open and not archived for the existing users, but we cannot provide any support or guarantee for the future development of this SDK.

# WebAssembly for Proxies (Go SDK) [![Build](https://github.com/tetratelabs/proxy-wasm-go-sdk/workflows/Test/badge.svg)](https://github.com/tetratelabs/proxy-wasm-go-sdk/actions) [![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

The Go SDK for
 [Proxy-Wasm](https://github.com/proxy-wasm/spec), enabling developers to write Proxy-Wasm plugins in Go. 
This SDK is powered by [TinyGo](https://tinygo.org/) and does not support the official Go compiler.
