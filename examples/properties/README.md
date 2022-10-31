## properties

this example prevalidates the authentication header via the usage of properties fetched from the proxy (i.e. Envoy metadata here).

### message on clients
```
curl localhost:18000/one -v
*   Trying 127.0.0.1:18000...
* Connected to localhost (127.0.0.1) port 18000 (#0)
> GET /one HTTP/1.1
> Host: localhost:18000
> User-Agent: curl/7.82.0
> Accept: */*
> 
< HTTP/1.1 401 Unauthorized
< date: Mon, 31 Oct 2022 00:53:01 GMT
< server: envoy
< content-length: 0
< 

curl localhost:18000/one -v -H 'cookie: value'
*   Trying 127.0.0.1:18000...
* Connected to localhost (127.0.0.1) port 18000 (#0)
> GET /one HTTP/1.1
> Host: localhost:18000
> User-Agent: curl/7.82.0
> Accept: */*
> cookie: value
> 
< HTTP/1.1 200 OK
< content-length: 13
< content-type: text/plain
< date: Mon, 31 Oct 2022 00:54:59 GMT
< server: envoy
< x-envoy-upstream-service-time: 0
< 
example body

curl localhost:18000/two -v
*   Trying 127.0.0.1:18000...
* Connected to localhost (127.0.0.1) port 18000 (#0)
> GET /two HTTP/1.1
> Host: localhost:18000
> User-Agent: curl/7.82.0
> Accept: */*
> 
< HTTP/1.1 401 Unauthorized
< date: Mon, 31 Oct 2022 00:53:30 GMT
< server: envoy
< content-length: 0
< 

curl localhost:18000/two -v -H 'authorization: token'
*   Trying 127.0.0.1:18000...
* Connected to localhost (127.0.0.1) port 18000 (#0)
> GET /two HTTP/1.1
> Host: localhost:18000
> User-Agent: curl/7.82.0
> Accept: */*
> authorization: token
> 
< HTTP/1.1 200 OK
< content-length: 13
< content-type: text/plain
< date: Mon, 31 Oct 2022 00:53:52 GMT
< server: envoy
< x-envoy-upstream-service-time: 0
< 
example body

curl localhost:18000/three -v
*   Trying 127.0.0.1:18000...
* Connected to localhost (127.0.0.1) port 18000 (#0)
> GET /three HTTP/1.1
> Host: localhost:18000
> User-Agent: curl/7.82.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< content-length: 13
< content-type: text/plain
< date: Mon, 31 Oct 2022 00:55:27 GMT
< server: envoy
< x-envoy-upstream-service-time: 0
< 
example body
```

### message on Envoy
```
wasm log: auth header is "cookie"
wasm log: 2 finished
wasm log: auth header is "cookie"
wasm log: 3 finished
wasm log: auth header is "authorization"
wasm log: 2 finished
wasm log: auth header is "authorization"
wasm log: 3 finished
wasm log: no auth header for route
wasm log: 4 finished
```
