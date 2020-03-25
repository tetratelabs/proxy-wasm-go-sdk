# proxy-wasm-go

proxy-wasm-go is an experimental runtime sdk for
 [proxy-wasm](https://github.com/proxy-wasm/spec) for Gophers which implements
 the low-level Application Binary Interface(ABI) called __Proxy-Wasm ABI__.
proxy-wasm-go is powered by [TinyGo](https://tinygo.org/), a go compiler for small places.
 
Note that Proxy-wasm ABI itself is in a very early stage 
so does proxy-wasm-go.

Any comments, questions, issues reports and PRs are welcome.
Please open issues/PRs or reach me [@mathetake](https://twitter.com/mathetake).

## TODOs
- docs
- support get/set shared queue
- support get/set shared data
- support get/set property
- support gRPC
- support enhance error handling
- support stream context

## references

- https://github.com/proxy-wasm/spec
- https://github.com/proxy-wasm/proxy-wasm-cpp-sdk
- https://github.com/proxy-wasm/proxy-wasm-rust-sdk
- https://tinygo.org/


Special thanks to TinyGo folks:)
