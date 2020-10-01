## http_auth_random

this example authorize requests depending on the hash values of a response from http://httpbin.org/uuid.

### message on clients
```
$ curl localhost:18000/uuid -v
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 18000 (#0)
> GET /uuid HTTP/1.1
> Host: localhost:18000
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< date: Wed, 25 Mar 2020 09:06:33 GMT
< content-type: application/json
< content-length: 53
< server: envoy
< access-control-allow-origin: *
< access-control-allow-credentials: true
< x-envoy-upstream-service-time: 1056
<
{
  "uuid": "e1020f65-f97a-47cd-9b31-368ba2063b6a"
}



# curl localhost:18000/uuid -v
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 18000 (#0)
> GET /uuid HTTP/1.1
> Host: localhost:18000
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 403 Forbidden
< content-length: 16
< content-type: text/plain
< powered-by: proxy-wasm-go-sdk!!
< date: Wed, 25 Mar 2020 09:07:36 GMT
< server: envoy
<
* Connection #0 to host localhost left intact
access forbidden

```

### message on Envoy

```
wasm log: request header from: :authority: localhost:18000
wasm log: request header from: :path: /uuid
wasm log: request header from: :method: GET
wasm log: request header from: user-agent: curl/7.68.0
wasm log: request header from: accept: */*
wasm log: request header from: x-forwarded-proto: http
wasm log: request header from: x-request-id: fddeac7b-db59-453c-9956-7f1050dbf6d5
wasm log: http call dispatched to httpbin
wasm log: response header from httpbin: :status: 200
wasm log: response header from httpbin: date: Thu, 01 Oct 2020 09:07:32 GMT
wasm log: response header from httpbin: content-type: application/json
wasm log: response header from httpbin: content-length: 53
wasm log: response header from httpbin: connection: keep-alive
wasm log: response header from httpbin: server: gunicorn/19.9.0
wasm log: response header from httpbin: access-control-allow-origin: *
wasm log: response header from httpbin: access-control-allow-credentials: true
wasm log: response header from httpbin: x-envoy-upstream-service-time: 340
wasm log: access granted
wasm log: 2 finished
wasm log: request header from: :authority: localhost:18000
wasm log: request header from: :path: /uuid
wasm log: request header from: :method: GET
wasm log: request header from: user-agent: curl/7.68.0
wasm log: request header from: accept: */*
wasm log: request header from: x-forwarded-proto: http
wasm log: request header from: x-request-id: 02628dd2-b985-4c4f-a2d2-164589c16f53
wasm log: http call dispatched to httpbin
wasm log: response header from httpbin: :status: 200
wasm log: response header from httpbin: date: Thu, 01 Oct 2020 09:07:34 GMT
wasm log: response header from httpbin: content-type: application/json
wasm log: response header from httpbin: content-length: 53
wasm log: response header from httpbin: connection: keep-alive
wasm log: response header from httpbin: server: gunicorn/19.9.0
wasm log: response header from httpbin: access-control-allow-origin: *
wasm log: response header from httpbin: access-control-allow-credentials: true
wasm log: response header from httpbin: x-envoy-upstream-service-time: 350
wasm log: access forbidden
wasm log: 3 finished
```
