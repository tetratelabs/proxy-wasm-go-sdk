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

type metricRootContext struct { proxywasm.DefaultRootContext }

func (ctx *metricRootContext) OnVMStart(int) bool {
	// initialize the metric
	counter = proxywasm.DefineCounterMetric("proxy_wasm_go.request_counter")
	return true
}

type metricHttpContext struct { proxywasm.DefaultHttpContext }

func (ctx *metricHttpContext) OnHttpRequestHeaders(int, bool) types.Action {
	// increment the request counter when we receive request headers
	counter.Increment(1)
	return types.ActionContinue
}
```

### requirements

proxy-wasm-go-sdk depends on TinyGo's [WASI](https://github.com/WebAssembly/WASI) (WebAssembly System Interface) target
which is introduced in [v0.16.0](https://github.com/tinygo-org/tinygo/releases/tag/v0.16.0).

Please follow the official instruction [here](https://tinygo.org/getting-started/).

### compatible ABI / Envoy builds (verified on CI)

| proxy-wasm-go-sdk| proxy-wasm ABI version |istio/proxyv2| Envoy upstream|
|:-------------:|:-------------:|:-------------:|:-------------:|
| main |  0.2.0|   v1.8.1 | v1.7.x |
| v0.0.13 |  0.2.0|   v1.8.1 | [5932dfd1f6b80eda9f56](https://github.com/envoyproxy/envoy/tree/5932dfd1f6b80eda9f56415ceff056a15a699c5e) |


## run examples

build:

```bash
make build.examples        # build all examples
make build.examples.docker # in docker

make build.example name=helloworld        # build a specific example
make build.example.docker name=helloworld # in docker
```

run:

```bash
make run name=helloworld
``` 

## sdk development

```bash
make test # run local tests without running envoy processes

## requires you to have Envoy binary locally
make test.e2e # run e2e tests

## requires you to have Envoy binary locally
make test.e2e.single name=helloworld # run e2e tests
```

## limitations and considerations

- Some of existing libraries are not available (importable but runtime panic / non-importable)
    - There are several reasons for this:
        1. TinyGo's WASI target does not support some of syscall: For example, we cannot import `crypto/rand` package.
        2. TinyGo does not implement all of reflect package([examples](https://github.com/tinygo-org/tinygo/blob/v0.14.1/src/reflect/value.go#L299-L305)).
        3. [proxy-wasm-cpp-host](https://github.com/proxy-wasm/proxy-wasm-cpp-host) has not supported some of WASI APIs yet 
        (see the [supported functions](https://github.com/proxy-wasm/proxy-wasm-cpp-host/blob/master/include/proxy-wasm/exports.h#L134-L147)).
    - These issues will be mitigated as TinyGo and proxy-wasm-cpp-host evolve.
- There's performance overhead of using Go/TinyGo due to GC
    - `runtime.GC` is called whenever the heap runs out (see [1](https://tinygo.org/lang-support/#garbage-collection),
    [2](https://github.com/tinygo-org/tinygo/blob/v0.14.1/src/runtime/gc_conservative.go#L218-L239)).
    - TinyGo allows us to disable GC, but we cannot do that since we need to use maps (implicitly causes allocation)
     for saving the plugin's [state](https://github.com/tetratelabs/proxy-wasm-go-sdk/blob/cf6ad74ed58b284d3d8ceeb8c5dba2280d5b1007/proxywasm/vmstate.go#L41-L46).
    - Theoretically, we can implement our own GC algorithms tailored for proxy-wasm through `alloc(uintptr)` [interface](https://github.com/tinygo-org/tinygo/blob/v0.14.1/src/runtime/gc_none.go#L13) 
    with `-gc=none` option. This is the future TODO.
- `recover` is [not implemented](https://github.com/tinygo-org/tinygo/issues/891) in TinyGo, and there's no way to prevent the WASM virtual machine from aborting.
- Goroutine support
    - In Tinygo, Goroutine is implmeneted through LLVM's coroutine (see [this blog post](https://aykevl.nl/2019/02/tinygo-goroutines)).
    - In Envoy, WASM modules are run in the event driven manner, and therefore the "scheduler" is not executed once the main function exits. 
        That means you cannot have the expected behavior of Goroutine as in ordinary host environments.
        - The question "How to deal with Goroutine in a thread local WASM VM executed in the event drive manner" has yet to be answered.
    - We strongly recommend that you implement the `OnTick` function for any asynchronous task instead of using Goroutine.
    - The scheduler can be disabled with `-scheduler=none` option of TinyGo.

## references

- https://github.com/proxy-wasm/spec
- https://github.com/proxy-wasm/proxy-wasm-cpp-sdk
- https://github.com/proxy-wasm/proxy-wasm-rust-sdk
- https://github.com/tetratelabs/envoy-wasm-rust-sdk
- https://tinygo.org/


Special thanks to TinyGo folks:)
