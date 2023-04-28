# WebAssembly for Proxies (Go SDK) [![Build](https://github.com/tetratelabs/proxy-wasm-go-sdk/workflows/Test/badge.svg)](https://github.com/tetratelabs/proxy-wasm-go-sdk/actions) [![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

The Go SDK for
 [Proxy-Wasm](https://github.com/proxy-wasm/spec), enabling developers to write Proxy-Wasm plugins in Go. 
This SDK is powered by [TinyGo](https://tinygo.org/) and does not support the official Go compiler.

## Getting Started

- [examples](examples) directory contains the example codes on top of this SDK.
- [OVERVIEW.md](doc/OVERVIEW.md) the overview of Proxy-Wasm, the API of this SDK, and the things you should know when writing plugins.

## Requirements

- [TinyGo](https://tinygo.org/) - This SDK depends on TinyGo and leverages its [WASI](https://github.com/WebAssembly/WASI) (WebAssembly System Interface) target. Please follow the official instruction [here](https://tinygo.org/getting-started/) for installing TinyGo.
- [Envoy](https://www.envoyproxy.io) - To run compiled examples, you need to have Envoy binary. We recommend using [func-e](https://func-e.io) as the easiest way to get started with Envoy. Alternatively, you can follow [the official instruction](https://www.envoyproxy.io/docs/envoy/latest/start/install).


## Dealing with memory issues

TinyGo's default memory allocator (Garbage Collector) is known to have some issues when it's used in the high workload environment (e.g. [1](https://github.com/tetratelabs/proxy-wasm-go-sdk/issues/349),[2](https://github.com/tetratelabs/proxy-wasm-go-sdk/issues/375)).
There's an alternative GC called [nottinygc](https://github.com/wasilibs/nottinygc) which not only resolves the memory related issues, but
also improves the performance on production usage.

The following images are an end user's observation on the perf of their Go SDK-compiled plugin on a high-workload environment.
This clearly indicates that nottinygc performs pretty well compared to the default setting of TinyGo.

![img](https://user-images.githubusercontent.com/13513977/235026482-ff8dcc3b-a7dc-444d-a1af-8137c64e1d53.png)
![img](https://user-images.githubusercontent.com/13513977/235026493-97122fe3-9de0-4417-93a0-dd3a32bebce7.png)

It can be enabled by adding a single line in your source code. Please refer to https://github.com/wasilibs/nottinygc for detail.

## Installation

```
go get github.com/tetratelabs/proxy-wasm-go-sdk
```

## Build and run Examples

```bash
# Build all examples.
make build.examples

# Build a specific example.
make build.example name=helloworld

# Run a specific example.
make run name=helloworld
```

## Compatible Envoy builds

Envoy is the first host side implementation of Proxy-Wasm ABI, 
and we run end-to-end tests with multiple versions of Envoy and Envoy-based [istio/proxy](https://github.com/istio/proxy) in order to verify Proxy-Wasm Go SDK works as expected.

Please refer to [workflow.yaml](.github/workflows/workflow.yaml) for which version is used for End-to-End tests.

## Build tags

The following build tags can be used to customize the behavior of the built plugin:

- `proxywasm_timing`: Enables logging of time spent in invocation of the plugin's exported functions. This can be useful for debugging performance issues.

## Contributing

We welcome contributions from the community! See [CONTRIBUTING.md](doc/CONTRIBUTING.md) for how to contribute to this repository.

## External links

- [WebAssembly for Proxies (ABI specification)](https://github.com/proxy-wasm/spec)
- [WebAssembly for Proxies (AssemblyScript SDK)](https://github.com/solo-io/proxy-runtime)
- [WebAssembly for Proxies (C++ SDK)](https://github.com/proxy-wasm/proxy-wasm-cpp-sdk)
- [WebAssembly for Proxies (Rust SDK)](https://github.com/proxy-wasm/proxy-wasm-rust-sdk)
- [WebAssembly for Proxies (Zig SDK)](https://github.com/mathetake/proxy-wasm-zig-sdk)
- [TinyGo - Go compiler for small places](https://tinygo.org/)
