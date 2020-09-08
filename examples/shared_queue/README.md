
## shared_queue

this example queues data to the shared queue on every request,
 and periodically processes the queued data on OnTick function.


```bash
[2020-09-08 10:55:48.699][644874][info][main] [external/envoy/source/server/server.cc:652] starting main dispatch loop
[2020-09-08 10:55:48.706][644897][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: queue registered, name: proxy_wasm_go.queue, id: 1
[2020-09-08 10:55:48.706][644897][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: set tick period milliseconds: 1000
[2020-09-08 10:55:48.711][644900][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: queue registered, name: proxy_wasm_go.queue, id: 1
[2020-09-08 10:55:48.712][644900][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: set tick period milliseconds: 1000
[2020-09-08 10:55:48.713][644898][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: queue registered, name: proxy_wasm_go.queue, id: 1
[2020-09-08 10:55:48.713][644898][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: set tick period milliseconds: 1000
[2020-09-08 10:55:48.716][644899][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: queue registered, name: proxy_wasm_go.queue, id: 1
[2020-09-08 10:55:48.716][644899][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: set tick period milliseconds: 1000
[2020-09-08 10:55:54.707][644874][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:hello
[2020-09-08 10:55:54.725][644897][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:world
[2020-09-08 10:55:54.726][644900][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:hello
[2020-09-08 10:55:54.732][644898][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:proxy-wasm
[2020-09-08 10:55:59.717][644874][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:hello
[2020-09-08 10:55:59.744][644900][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:world
[2020-09-08 10:55:59.744][644897][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:hello
[2020-09-08 10:55:59.744][644898][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:proxy-wasm
[2020-09-08 10:56:01.719][644874][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:hello
[2020-09-08 10:56:01.749][644900][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:world
[2020-09-08 10:56:01.749][644897][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:hello
[2020-09-08 10:56:01.749][644898][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:proxy-wasm
[2020-09-08 10:56:09.727][644874][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:hello
[2020-09-08 10:56:09.768][644898][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:world
[2020-09-08 10:56:09.768][644900][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:hello
[2020-09-08 10:56:09.770][644899][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:proxy-wasm
[2020-09-08 10:56:10.732][644874][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:hello
[2020-09-08 10:56:10.773][644900][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:hello
[2020-09-08 10:56:10.773][644899][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:proxy-wasm
[2020-09-08 10:56:10.773][644898][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:world
[2020-09-08 11:00:12.269][644874][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:hello
[2020-09-08 11:00:12.276][644900][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:world
[2020-09-08 11:00:12.280][644898][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:hello
[2020-09-08 11:00:12.293][644897][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log my_root_id: dequed data:proxy-wasm
```
