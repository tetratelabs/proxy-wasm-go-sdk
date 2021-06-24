__This project is in its early stage, and the API is likely to change and not stable.__

# WebAssembly for Proxies (Go SDK) [![Build](https://github.com/tetratelabs/proxy-wasm-go-sdk/workflows/test/badge.svg)](https://github.com/tetratelabs/proxy-wasm-go-sdk/actions) [![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

The Go SDK for
 [Proxy-Wasm](https://github.com/proxy-wasm/spec), enabling developers to write Proxy-Wasm plugins in Go. 
This SDK is powered by [TinyGo](https://tinygo.org/) and does not support the official Go compiler.

## Getting Started

- [examples](examples) directory contains the example codes on top of this SDK.

## Requirements

- [TinyGo](https://tinygo.org/) - This SDK depends on TinyGo and leverages its [WASI](https://github.com/WebAssembly/WASI) (WebAssembly System Interface) target. Please follow the official instruction [here](https://tinygo.org/getting-started/) for installing TinyGo.
- [Envoy](https://www.envoyproxy.io) - To run compiled examples, you need to have Envoy binary. Please follow [the official instruction](https://www.envoyproxy.io/docs/envoy/latest/start/install).

## Build and run Examples

```bash
# Build all examples.
make build.examples

# Build a specific example.
make build.example name=helloworld

# Run a specific example.
make run name=helloworld
```

## Compatible Envoy builds (tested on CI)

Envoy is the first host side implementation of Proxy-Wasm ABI, 
and we run end-to-end tests with multiple Envoy versions in order to verify Proxy-Wasm Go SDK works as expected.

| proxy-wasm-go-sdk| istio/proxyv2| Envoy upstream|
|:-------------:|:-------------:|:-------------:|
| main | 1.9, 1.10 | 1.18 |
| v0.3.0 | 1.9, 1.10 | 1.18 |

## Contributing

We welcome contributions from the community! See [CONTRIBUTING.md](doc/CONTRIBUTING.md) for how to contribute to this repository.

## External links

- [WebAssembly for Proxies (ABI specification)](https://github.com/proxy-wasm/spec)
- [WebAssembly for Proxies (AssemblyScript SDK)](https://github.com/solo-io/proxy-runtime)
- [WebAssembly for Proxies (C++ SDK)](https://github.com/proxy-wasm/proxy-wasm-cpp-sdk)
- [WebAssembly for Proxies (Rust SDK)](https://github.com/proxy-wasm/proxy-wasm-rust-sdk)
- [WebAssembly for Proxies (Zig SDK)](https://github.com/mathetake/proxy-wasm-zig-sdk)
- [TinyGo - Go compiler for small places](https://tinygo.org/)
