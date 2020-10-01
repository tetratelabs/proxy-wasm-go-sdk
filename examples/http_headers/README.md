## http_headers

this example handles http request/response headers events and log all headers.

```
wasm log: request header --> :authority: localhost:18000
wasm log: request header --> :path: /uuid
wasm log: request header --> :method: GET
wasm log: request header --> user-agent: curl/7.68.0
wasm log: request header --> accept: */*
wasm log: request header --> x-forwarded-proto: http
wasm log: request header --> x-request-id: 5692b633-fd9c-4700-b4dd-7a58e2853eb4
wasm log: response header <-- :status: 200
wasm log: response header <-- content-length: 13
wasm log: response header <-- content-type: text/plain
wasm log: response header <-- date: Thu, 01 Oct 2020 09:10:09 GMT
wasm log: response header <-- server: envoy
wasm log: response header <-- x-envoy-upstream-service-time: 0
wasm log: 2 finished
```
