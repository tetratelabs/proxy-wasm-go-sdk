
## shared_queue

This example describes how to use a shared queue to communicate between seprate Wasm VMs.

There are two Wasm VMs are configured (See `envoy.yaml` for detail):
1. The one with `vm_id="receiver"` and the binary of `receiver/main.go`.
2. Another one with `vm_id="sender"` and the binary of `sender/main.go`.

`receiver` VM runs as a singleton [Wasm Service](https://www.envoyproxy.io/docs/envoy/latest/configuration/other_features/wasm_service.html) which runs in the main thread, and there are two plugin configurations are given.
One is `http_response_headers` and another is `http_request_headers`. Each of these plugin registers a shared queue whose name equals that configuration respectively.

`sender` VM runs in the http filter chain on worker threads, and enqueue request headers to the shared queue resolved by the `ResolveSharedQueue` with the args of (`vm_id=receiver`,`name=http_request_headers`) and (`vm_id=receiver`,`name=http_response_headers`).

See [this talk](https://www.youtube.com/watch?v=XdWmm_mtVXI&t=1171s) for detail.


```bash
wasm log sender: enqueued data: {"key": ":authority","value": "localhost:18000"}
wasm log receiver: dequeued data: {"key": ":authority","value": "localhost:18000"}
wasm log sender: enqueued data: {"key": ":path","value": "/"}
wasm log receiver: dequeued data: {"key": ":path","value": "/"}
wasm log sender: enqueued data: {"key": ":method","value": "GET"}
wasm log receiver: dequeued data: {"key": ":method","value": "GET"}
wasm log sender: enqueued data: {"key": ":scheme","value": "http"}
wasm log receiver: dequeued data: {"key": ":scheme","value": "http"}
wasm log sender: enqueued data: {"key": "user-agent","value": "curl/7.68.0"}
wasm log receiver: dequeued data: {"key": "user-agent","value": "curl/7.68.0"}
wasm log sender: enqueued data: {"key": "accept","value": "*/*"}
wasm log receiver: dequeued data: {"key": "accept","value": "*/*"}
wasm log sender: enqueued data: {"key": "x-forwarded-proto","value": "http"}
wasm log receiver: dequeued data: {"key": "x-forwarded-proto","value": "http"}
wasm log sender: enqueued data: {"key": "x-request-id","value": "73a13840-6ca2-4f9f-a639-91b014c8d485"}
wasm log receiver: dequeued data: {"key": "x-request-id","value": "73a13840-6ca2-4f9f-a639-91b014c8d485"}
```
