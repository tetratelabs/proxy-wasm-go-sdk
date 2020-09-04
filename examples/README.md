
Theses are the proxy-wasm-go reimplementation of examples in https://github.com/proxy-wasm/proxy-wasm-rust-sdk/tree/master/examples.

## requirements

- TinyGo(0.14.0+): https://tinygo.org/
- GetEnvoy: https://www.getenvoy.io/install/

To download compatible envoyproxy, run
```bash
getenvoy fetch wasm:1.15
```

## build

```bash
tinygo build -o ./${example}/wasm.wasm -wasm-abi=generic -target wasm ./${example}/main.go
```

## run

```bash
getenvoy run wasm:1.15 -- -c ./${example}/envoy.yaml
``` 
