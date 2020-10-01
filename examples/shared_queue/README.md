
## shared_queue

this example queues data to the shared queue on every request,
 and periodically processes the queued data on OnTick function.


```bash
wasm log my_root_id: queue registered, name: proxy_wasm_go.queue, id: 1
wasm log my_root_id: set tick period milliseconds: 100
wasm log my_root_id: queue registered, name: proxy_wasm_go.queue, id: 1
wasm log my_root_id: set tick period milliseconds: 100
wasm log my_root_id: dequeued data: hello
wasm log my_root_id: dequeued data: world
wasm log my_root_id: dequeued data: hello
wasm log my_root_id: dequeued data: proxy-wasm
wasm log my_root_id: dequeued data: hello
wasm log my_root_id: dequeued data: world
wasm log my_root_id: dequeued data: hello
wasm log my_root_id: dequeued data: proxy-wasm
wasm log my_root_id: dequeued data: hello
wasm log my_root_id: dequeued data: world
wasm log my_root_id: dequeued data: hello
wasm log my_root_id: dequeued data: proxy-wasm
```
