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

### requirements

proxy-wasm-go-sdk depends on the latest TinyGo which supports WASI target([tinygo-org/tinygo#1373](https://github.com/tinygo-org/tinygo/pull/1373)).
In order to install that, simply run (Ubuntu/Debian):

```shell
# this corresponds to https://github.com/tinygo-org/tinygo/commit/f50ad3585d084b17f7754f4b3cb0d42661fee036
wget https://19227-136505169-gh.circle-artifacts.com/0/tmp/tinygo_amd64.deb
dpkg -i tinygo_amd64.deb
```

Alternatively, you can use the pre-built docker container `getenvoy/extention-tingyo-builder:wasi-dev` for any platform.

TinyGo's official release of WASI target will come soon, and after that you could
 just follow https://tinygo.org/getting-started/ to install the requirement on any platform. Stay tuned!


### compatible Envoy builds

| proxy-wasm-go-sdk| proxy-wasm ABI version | envoyproxy/envoy-wasm| istio/proxyv2|
|:-------------:|:-------------:|:-------------:|:-------------:|
| main |  0.2.0|  N/A  |   v1.17.x |
| v0.0.4 |  0.2.0|  N/A  |   v1.17.x |
| v0.0.3 |  0.2.0|  N/A  |   v1.17.x |
| v0.0.2 | 0.1.0|release/v1.15 | N/A |


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

## compiler limitations and considerations

- Some of existing libraries are not available (importable but runtime panic / non-importable)
    - There are two reasons for this:
        1. TinyGo's WASI target does not support some of syscall: For example, we cannot import `crypto/rand` package.
        2. TinyGo does not implement all of reflect package([examples](https://github.com/tinygo-org/tinygo/blob/v0.14.1/src/reflect/value.go#L299-L305)).
    - These issues will be mitigated as the TinyGo improves.
- There's performance overhead in using Go/TinyGo due to GC
    - runtime.GC() is called whenever the heap runs out (see [1](https://tinygo.org/lang-support/#garbage-collection),
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
