# Go SDK for WebAssembly-based Envoy extensions
[![Build](https://github.com/tetratelabs/proxy-wasm-go-sdk/workflows/build-test/badge.svg)](https://github.com/tetratelabs/proxy-wasm-go-sdk/actions)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

The Go sdk for
 [proxy-wasm](https://github.com/proxy-wasm/spec), enabling developers to write Envoy extensions in Go.

proxy-wasm-go-sdk is powered by [TinyGo](https://tinygo.org/) and does not support the official Go compiler.


```golang

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

var counter proxywasm.MetricCounter

const metricName = "proxy_wasm_go.request_counter"

type context struct{ proxywasm.DefaultContext }

func (ctx *context) OnVMStart(int) bool {
	// initialize the new metric	
	counter, _ = proxywasm.DefineCounterMetric(metricName)
	return true
}

func (ctx *context) OnHttpRequestHeaders(int, bool) types.Action {
	// increment the request counter when we receive request headers
	counter.Increment(1)  
	return types.ActionContinue
}

```




## requirements

- TinyGo(0.14.0+): https://tinygo.org/
- GetEnvoy: https://www.getenvoy.io/install/ (for running examples)

To download compatible envoyproxy, run
```bash
getenvoy fetch wasm:1.15
```

The target Envoy version is `release/v1.15`
 branch on [envoyproxy/envoy-wasm](https://github.com/envoyproxy/envoy-wasm/tree/release/v1.15).

## run examples

build:

```bash
make build.examples # build all examples

make build.example name=helloworld # build a specific example
```

run:

```bash
make run name=helloworld
``` 

## sdk development

```bash
make test # run local tests without running envoy processes

make test.e2e # run e2e tests
```

## language and compiler limitations/considerations

- You can only use really limited set of existing libraries. 
    - There are two reasons for this:
        1. TinyGo [uses](https://github.com/tinygo-org/tinygo/blob/release/loader/loader.go#L79-L83) the official parser, lexer, etc., 
    which forces your program to implicitly import `syscall/js` package if you import packages using system calls. 
    The package expects the host environment to have the `syscall/js` specific ABI as in 
    [wasm_exec.js](https://github.com/tinygo-org/tinygo/blob/154d4a781f6121bd6f584cca4a88909e0b091f63/targets/wasm_exec.js) 
    which is not available outside of that javascript.
        2. TinyGo does not implement all of reflect package([examples](https://github.com/tinygo-org/tinygo/blob/v0.14.1/src/reflect/value.go#L299-L305)).
    - The syscall problem can be solved by one of the followings:
        - emulate `syscall/js` function through WASI interface (which is implemented by V8 engine running on Envoy)
        - support WASI target in TinyGo: https://github.com/tinygo-org/tinygo/pull/1373
- There's performance overhead in using Go/TinyGo due to GC
    - runtime.GC() is called whenever heap allocation happens (see [1](https://tinygo.org/lang-support/#garbage-collection), 
    [2](https://github.com/tinygo-org/tinygo/blob/v0.14.1/src/runtime/gc_conservative.go#L218-L239)).
    - TinyGo allows us to disable GC, but we cannot do that since we need to use maps (implicitly causes allocation)
     for saving the plugin's [state](https://github.com/tetratelabs/proxy-wasm-go-sdk/blob/master/proxywasm/vmstate.go#L17-L22).
    - Theoretically, we can implement our own GC algorithms tailored for proxy-wasm through `alloc(uintptr)` [interface](https://github.com/tinygo-org/tinygo/blob/v0.14.1/src/runtime/gc_none.go#L13) 
    with `-gc=none` option. This is the future TODO.
- `recover` is [not implemented](https://github.com/tinygo-org/tinygo/issues/891) in TinyGo, and there's no way to prevent the WASM virtual machine from aborting.
- Be careful about using Goroutine
    - In Tinygo, Goroutine is implmeneted through LLVM's coroutine (see [this blog post](https://aykevl.nl/2019/02/tinygo-goroutines)).
    - Make every goroutine exit as soon as possible, otherwise it will block the proxy's worker thread. That is too bad for processing realtime network requests.
    - We strongly recommend that you implement the `OnTick` function for any asynchronous task instead of using Goroutine so we do not block requests.


## references

- https://github.com/proxy-wasm/spec
- https://github.com/proxy-wasm/proxy-wasm-cpp-sdk
- https://github.com/proxy-wasm/proxy-wasm-rust-sdk
- https://github.com/tetratelabs/envoy-wasm-rust-sdk
- https://tinygo.org/


Special thanks to TinyGo folks:)
