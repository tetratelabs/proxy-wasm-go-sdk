# proxy-wasm-go

proxy-wasm-go is an experimental runtime sdk for
 [proxy-wasm](https://github.com/proxy-wasm/spec) for Gophers which implements
 the low-level Application Binary Interface(ABI) called __Proxy-Wasm ABI__.
proxy-wasm-go is powered by [TinyGo](https://tinygo.org/), a Go compiler for small places.


## requirements

- TinyGo(0.14.0+): https://tinygo.org/
- GetEnvoy: https://www.getenvoy.io/install/

To download compatible envoyproxy, run
```bash
getenvoy fetch wasm:1.15
```

The Envoy's version is the `release/v1.15`
 branch on [envoyproxy/envoy-wasm](https://github.com/envoyproxy/envoy-wasm/tree/release/v1.15).

## examples
Theses are the proxy-wasm-go reimplementation of examples in https://github.com/proxy-wasm/proxy-wasm-rust-sdk/tree/master/examples.

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

## references

- https://github.com/proxy-wasm/spec
- https://github.com/proxy-wasm/proxy-wasm-cpp-sdk
- https://github.com/proxy-wasm/proxy-wasm-rust-sdk
- https://tinygo.org/


Special thanks to TinyGo folks:)
