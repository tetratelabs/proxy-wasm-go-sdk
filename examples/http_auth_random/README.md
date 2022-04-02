## http_auth_random
this example authorize requests depending on the hash values of a response from http://httpbin.org/uuid.

这里是在envoy中调用一个远程服务（`httpbin uuid接口`）的例子，主要是使用 wasm 的`proxywasm.DispatchHttpCall`方法

### message on clients
```shell
$ curl -v localhost:18000/uuid
*   Trying 127.0.0.1:18000...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 18000 (#0)
> GET /uuid HTTP/1.1
> Host: localhost:18000
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< date: Sat, 02 Apr 2022 08:22:50 GMT
< content-type: application/json
< content-length: 53
< server: envoy
< access-control-allow-origin: *
< access-control-allow-credentials: true
< x-envoy-upstream-service-time: 2396
< 
{
  "uuid": "cc33b3eb-290b-4131-8a8e-a2117f383f0e"
}
* Connection #0 to host localhost left intact
```
```shell
$ curl -v localhost:18000/uuid
*   Trying 127.0.0.1:18000...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 18000 (#0)
> GET /uuid HTTP/1.1
> Host: localhost:18000
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 403 Forbidden
< powered-by: proxy-wasm-go-sdk!!
< content-length: 16
< content-type: text/plain
< date: Sat, 02 Apr 2022 08:34:13 GMT
< server: envoy
< 
* Connection #0 to host localhost left intact
access forbidden

```

### message on Envoy

```shell
wasm log: <---- pluginCx NewHttpContext ---->
wasm log: <---- OnHttpRequestHeaders ---->
wasm log: request header: :authority: localhost:18000
wasm log: request header: :path: /uuid
wasm log: request header: :method: GET
wasm log: request header: :scheme: http
wasm log: request header: user-agent: curl/7.68.0
wasm log: request header: accept: */*
wasm log: request header: x-forwarded-proto: http
wasm log: request header: x-request-id: 2d364422-5c29-4707-9b86-5864a372e501
wasm log: http call dispatched to httpbin
wasm log: <---- httpCallResponseCallBack ---->
wasm log: response header from httpbin: :status: 200
wasm log: response header from httpbin: date: Sat, 02 Apr 2022 08:22:47 GMT
wasm log: response header from httpbin: content-type: application/json
wasm log: response header from httpbin: content-length: 53
wasm log: response header from httpbin: connection: keep-alive
wasm log: response header from httpbin: server: gunicorn/19.9.0
wasm log: response header from httpbin: access-control-allow-origin: *
wasm log: response header from httpbin: access-control-allow-credentials: true
wasm log: response header from httpbin: x-envoy-upstream-service-time: 1403
wasm log: access granted
```
