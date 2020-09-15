## helloworld

this example handles tcp connections and output data into the log stream.


```bash
# curl localhost:18000

[2020-09-14 14:04:40.804][390417][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: new connection!
[2020-09-14 14:04:40.804][390417][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: downstream data received: GET /uuid HTTP/1.1
Host: localhost:18000
User-Agent: curl/7.68.0
Accept: */*


[2020-09-14 14:04:40.805][390417][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: downstream connection close!
```
