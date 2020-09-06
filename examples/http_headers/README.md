## http_headers

this example handles http request/response headers events and log all headers.


```bash
proxy_1  | [2020-03-25 09:09:24.937][16][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1077] wasm log my_plugin my_root_id: request header: :authority: localhost:18000
proxy_1  | [2020-03-25 09:09:24.937][16][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1077] wasm log my_plugin my_root_id: request header: :path: /
proxy_1  | [2020-03-25 09:09:24.937][16][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1077] wasm log my_plugin my_root_id: request header: :method: GET
proxy_1  | [2020-03-25 09:09:24.937][16][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1077] wasm log my_plugin my_root_id: request header: user-agent: curl/7.54.0
proxy_1  | [2020-03-25 09:09:24.937][16][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1077] wasm log my_plugin my_root_id: request header: accept: */*
proxy_1  | [2020-03-25 09:09:24.937][16][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1077] wasm log my_plugin my_root_id: request header: hello: Go
proxy_1  | [2020-03-25 09:09:24.937][16][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1077] wasm log my_plugin my_root_id: request header: x-forwarded-proto: http
proxy_1  | [2020-03-25 09:09:24.937][16][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1077] wasm log my_plugin my_root_id: request header: x-request-id: 6542a4ca-b6b3-4667-b41d-7ef4c3392946
proxy_1  | [2020-03-25 09:09:24.942][16][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1077] wasm log my_plugin my_root_id: response header: :status: 200
proxy_1  | [2020-03-25 09:09:24.942][16][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1077] wasm log my_plugin my_root_id: response header: content-length: 13
proxy_1  | [2020-03-25 09:09:24.942][16][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1077] wasm log my_plugin my_root_id: response header: content-type: text/plain
proxy_1  | [2020-03-25 09:09:24.942][16][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1077] wasm log my_plugin my_root_id: response header: date: Wed, 25 Mar 2020 09:09:24 GMT
proxy_1  | [2020-03-25 09:09:24.942][16][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1077] wasm log my_plugin my_root_id: response header: server: envoy
proxy_1  | [2020-03-25 09:09:24.942][16][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1077] wasm log my_plugin my_root_id: response header: x-envoy-upstream-service-time: 2
```
