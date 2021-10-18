
## metrics

this example creates simple request counter

```
wasm log: previous value of proxy_wasm_go.request_counter: 0
wasm log: incremented
wasm log: previous value of proxy_wasm_go.request_counter: 1
wasm log: incremented
wasm log: previous value of proxy_wasm_go.request_counter: 2
wasm log: incremented
wasm log: previous value of proxy_wasm_go.request_counter: 3
wasm log: incremented
wasm log: previous value of proxy_wasm_go.request_counter: 4
wasm log: incremented
wasm log: previous value of proxy_wasm_go.request_counter: 5
wasm log: incremented
```

```
$ curl -s 'localhost:8001/stats/prometheus'| grep proxy
# TYPE proxy_wasm_go_request_counter counter
proxy_wasm_go_request_counter{} 5
```
