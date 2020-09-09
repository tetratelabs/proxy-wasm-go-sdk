# Go SDK for WebAssembly-based Envoy extensions

The Go sdk for
 [proxy-wasm](https://github.com/proxy-wasm/spec), enabling developers to write Envoy extensions in Go.

proxy-wasm-go-sdk is powered by [TinyGo](https://tinygo.org/) and does not support the official Go compiler.


## requirements

- TinyGo(0.14.0+): https://tinygo.org/
- GetEnvoy: https://www.getenvoy.io/install/ (for running examples)

To download compatible envoyproxy, run
```bash
getenvoy fetch wasm:1.15
```

The targe Envoy version is `release/v1.15`
 branch on [envoyproxy/envoy-wasm](https://github.com/envoyproxy/envoy-wasm/tree/release/v1.15).

## run examples

build:

```bash
find ./examples -type f -name "main.go" | xargs -Ip tinygo build -o p.wasm -target=wasm -wasm-abi=generic p
```

run:

```bash
getenvoy run wasm:1.15 -- -c ./examples/${name}/envoy.yaml
``` 

## sdk development

To run tests:

```bash
go test -tags=proxytest -v -race ./...
```

## limitations / considerations

TODO

## references

- https://github.com/proxy-wasm/spec
- https://github.com/proxy-wasm/proxy-wasm-cpp-sdk
- https://github.com/proxy-wasm/proxy-wasm-rust-sdk
- https://github.com/tetratelabs/envoy-wasm-rust-sdk
- https://tinygo.org/


Special thanks to TinyGo folks:)
