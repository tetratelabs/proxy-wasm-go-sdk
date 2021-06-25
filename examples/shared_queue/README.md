
## shared_queue

This example describes how to use a shared queue to communicate between seprate Wasm VMs.

There are two Wasm VMs are configured (See `envoy.yaml` for detail):
1. The one with `vm_id="receiver"` and the binary of `receiver/main.go`.
2. Another one with `vm_id="sender"` and the binary of `sender/main.go`.

`receiver` VM runs as a singleton [Wasm Service](https://www.envoyproxy.io/docs/envoy/latest/configuration/other_features/wasm_service.html) which runs in the main thread, and there are **three** plugin configurations are given. These configuration values are `http_response_headers`, `http_request_headers` and `tcp_data_hashes`.
Each of these corresponding PluginContext registers a shared queue whose name equals that configuration respectively.

`sender` VM runs in a http filter and a network filter chain on worker threads, and 
- enqueue request headers to the shared queue resolved by the `ResolveSharedQueue` with the args of (`vm_id=receiver`,`name=http_request_headers`) and (`vm_id=receiver`,`name=http_response_headers`).
- enqueue hash values of tcp data frames to the shared queue resolved by the `ResolveSharedQueue` with the args of (`vm_id=receiver`,`name=tcp_data_hashes`).

See [this talk](https://www.youtube.com/watch?v=XdWmm_mtVXI&t=1171s) for detail.


```bash
wasm log receiver: queue "http_request_headers" registered as queueID=1 by contextID=1
wasm log receiver: queue "http_request_headers" registered as queueID=1 by contextID=1
wasm log receiver: queue "http_response_headers" registered as queueID=2 by contextID=2
wasm log receiver: queue "tcp_data_hashes" registered as queueID=3 by contextID=3
all clusters initialized. initializing init manager
all dependencies initialized. starting workers
wasm log sender: contextID=1 is configured for http
wasm log sender: contextID=2 is configured for tcp
wasm log sender: contextID=1 is configured for http
wasm log sender: contextID=2 is configured for tcp

....

# curl localhost:18000

wasm log sender: enqueued data: {"key": ":authority","value": "localhost:18000"}
wasm log sender: enqueued data: {"key": ":path","value": "/"}
wasm log sender: enqueued data: {"key": ":method","value": "GET"}
wasm log sender: enqueued data: {"key": ":scheme","value": "http"}
wasm log sender: enqueued data: {"key": "user-agent","value": "curl/7.68.0"}
wasm log sender: enqueued data: {"key": "accept","value": "*/*"}
wasm log sender: enqueued data: {"key": "x-forwarded-proto","value": "http"}
wasm log sender: enqueued data: {"key": "x-request-id","value": "57d77551-02e8-455c-bf86-45a0d7308a0e"}
wasm log receiver: (contextID=1) dequeued data from http_request_headers(queueID=1): {"key": ":authority","value": "localhost:18000"}
wasm log receiver: (contextID=1) dequeued data from http_request_headers(queueID=1): {"key": ":path","value": "/"}
wasm log receiver: (contextID=1) dequeued data from http_request_headers(queueID=1): {"key": ":method","value": "GET"}
wasm log receiver: (contextID=1) dequeued data from http_request_headers(queueID=1): {"key": ":scheme","value": "http"}
wasm log receiver: (contextID=1) dequeued data from http_request_headers(queueID=1): {"key": "user-agent","value": "curl/7.68.0"}
wasm log receiver: (contextID=1) dequeued data from http_request_headers(queueID=1): {"key": "accept","value": "*/*"}
wasm log receiver: (contextID=1) dequeued data from http_request_headers(queueID=1): {"key": "x-forwarded-proto","value": "http"}
wasm log receiver: (contextID=1) dequeued data from http_request_headers(queueID=1): {"key": "x-request-id","value": "57d77551-02e8-455c-bf86-45a0d7308a0e"}
wasm log sender: (contextID=3) enqueued data: {"key": ":status","value": "200"}
wasm log receiver: (contextID=2) dequeued data from http_response_headers(queueID=2): {"key": ":status","value": "200"}
wasm log sender: (contextID=3) enqueued data: {"key": "content-length","value": "13"}
wasm log receiver: (contextID=2) dequeued data from http_response_headers(queueID=2): {"key": "content-length","value": "13"}
wasm log sender: (contextID=3) enqueued data: {"key": "content-type","value": "text/plain"}
wasm log receiver: (contextID=2) dequeued data from http_response_headers(queueID=2): {"key": "content-type","value": "text/plain"}

# curl localhost:18001

wasm log sender: (contextID=4) enqueued data: 7d1a184bc958cdb9f1fee6591a3f2ae2
wasm log receiver: (contextID=3) dequeued data from tcp_data_hashes(queueID=3): 7d1a184bc958cdb9f1fee6591a3f2ae2
```
