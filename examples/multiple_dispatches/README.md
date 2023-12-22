## multiple_dispatches

This example dispatches multiple http calls to remote clusters while pausing the original http response processing from the upstream.
Once the plugin recieved all the responses to all dispatched calls, it adds an http header to the original http response, and resumes it 
inside the dispatched callback.

Note: the same logic can be performed for http *requests* as well with the corresponding functions.

### message on clients
```
$ curl --head localhost:18000
HTTP/1.1 200 OK
date: Tue, 16 Aug 2022 03:21:37 GMT
content-type: text/html; charset=utf-8
content-length: 9593
server: envoy
access-control-allow-origin: *
access-control-allow-credentials: true
x-envoy-upstream-service-time: 362
total-dispatched: 10  <---- added inside the dispatched callback.
```

### message on Envoy

```
wasm log: pending dispatched requests: 9
wasm log: pending dispatched requests: 8
wasm log: pending dispatched requests: 7
wasm log: pending dispatched requests: 6
wasm log: pending dispatched requests: 5
wasm log: pending dispatched requests: 4
wasm log: pending dispatched requests: 3
wasm log: pending dispatched requests: 2
wasm log: pending dispatched requests: 1
wasm log: response resumed after processed 10 dispatched request
```
