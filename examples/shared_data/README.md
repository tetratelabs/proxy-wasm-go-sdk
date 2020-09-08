
## shared_queue

this example uses the shared key value store (across VMs) 
and increments the value in response to http requests atomically.

```
[2020-09-08 13:19:18.197][97057][warning][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1000] wasm log my_root_id: error setting shared data on OnVMStart: cas mismatch
[2020-09-08 13:19:18.198][97059][warning][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1000] wasm log my_root_id: error setting shared data on OnVMStart: cas mismatch
[2020-09-08 13:19:18.199][97067][warning][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1000] wasm log my_root_id: error setting shared data on OnVMStart: cas mismatch
[2020-09-08 13:19:29.913][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 1
[2020-09-08 13:19:29.922][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 2
[2020-09-08 13:19:29.930][97067][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 3
[2020-09-08 13:19:29.935][97067][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 4
[2020-09-08 13:19:29.940][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 5
[2020-09-08 13:19:29.948][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 6
[2020-09-08 13:19:29.956][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 7
[2020-09-08 13:19:29.964][97067][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 8
[2020-09-08 13:19:29.971][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 9
[2020-09-08 13:19:29.976][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 10
[2020-09-08 13:19:29.982][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 11
[2020-09-08 13:19:29.990][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 12
[2020-09-08 13:19:29.997][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 13
[2020-09-08 13:19:30.005][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 14
[2020-09-08 13:19:30.013][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 15
[2020-09-08 13:19:30.021][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 16
[2020-09-08 13:19:30.029][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 17
[2020-09-08 13:19:30.036][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 18
[2020-09-08 13:19:30.044][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 19
[2020-09-08 13:19:30.051][97067][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 20
[2020-09-08 13:19:40.191][97067][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 21
[2020-09-08 13:19:40.199][97105][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 22
[2020-09-08 13:19:40.208][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 23
[2020-09-08 13:19:40.215][97107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: shared value: 24
```