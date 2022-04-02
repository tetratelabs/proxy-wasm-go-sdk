## http_headers

this example handles http request/response headers events and log all headers.

call envoy
```
curl -v localhost:18000
```
调用后的日志
```shell

root@13:~/proxy-wasm-go-sdk/examples/http_headers# curl -v localhost:18000
*   Trying 127.0.0.1:18000...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 18000 (#0)
> GET / HTTP/1.1
> Host: localhost:18000
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< content-length: 13
< content-type: text/plain
< date: Sat, 02 Apr 2022 03:12:27 GMT
< server: envoy
< x-envoy-upstream-service-time: 0
< 
example body
* Connection #0 to host localhost left intact
```

envoy的日志
```shell
[2022-04-02 03:11:16.480][34362][info][main] [external/envoy/source/server/server.cc:745] all clusters initialized. initializing init manager
[2022-04-02 03:11:16.480][34362][info][config] [external/envoy/source/server/listener_manager_impl.cc:888] all dependencies initialized. starting workers
[2022-04-02 03:11:16.481][34362][info][main] [external/envoy/source/server/server.cc:764] starting main dispatch loop
[2022-04-02 03:12:27.921][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> :authority: localhost:18000
[2022-04-02 03:12:27.921][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> :path: /
[2022-04-02 03:12:27.921][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> :method: GET
[2022-04-02 03:12:27.921][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> :scheme: http
[2022-04-02 03:12:27.921][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> user-agent: curl/7.68.0
[2022-04-02 03:12:27.921][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> accept: */*
[2022-04-02 03:12:27.921][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> x-forwarded-proto: http
[2022-04-02 03:12:27.921][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> x-request-id: de03fd26-3e04-47fa-9180-dd2aa0fc9cbd
[2022-04-02 03:12:27.921][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> test: best
[2022-04-02 03:12:27.922][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: response header <-- :status: 200
[2022-04-02 03:12:27.922][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: response header <-- content-length: 13
[2022-04-02 03:12:27.922][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: response header <-- content-type: text/plain
[2022-04-02 03:12:27.922][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: response header <-- date: Sat, 02 Apr 2022 03:12:27 GMT
[2022-04-02 03:12:27.922][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: response header <-- server: envoy
[2022-04-02 03:12:27.922][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: response header <-- x-envoy-upstream-service-time: 0
[2022-04-02 03:12:27.922][34372][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: 2 finished

```






