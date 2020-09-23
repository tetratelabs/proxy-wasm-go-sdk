## helloworld

this example handles tcp connections and output data into the log stream.


```bash
[2020-09-23 15:19:57.054][1008407][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: new connection!
[2020-09-23 15:19:57.055][1008407][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: downstream data received: GET /uuid HTTP/1.1
Host: localhost:18000
User-Agent: curl/7.68.0
Accept: */*


[2020-09-23 15:19:57.055][1008407][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: remote address: 127.0.0.1:8099
[2020-09-23 15:19:57.055][1008407][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: upstream data received: HTTP/1.1 200 OK
content-length: 13
content-type: text/plain
date: Wed, 23 Sep 2020 06:19:56 GMT
server: envoy

example body

[2020-09-23 15:19:57.055][1008407][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: remote address: 127.0.0.1:8099
[2020-09-23 15:19:57.055][1008407][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: downstream connection close!
[2020-09-23 15:19:57.056][1008407][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:997] wasm log: connection complete!
```
