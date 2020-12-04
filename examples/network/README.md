## network

this example handles tcp connections and output data into the log stream.


```bash
wasm log: new connection!
wasm log: >>>>>> downstream data received >>>>>>
GET /uuid HTTP/1.1
Host: localhost:18000
User-Agent: curl/7.68.0
Accept: */*


wasm log: remote address: 127.0.0.1:8099
wasm log: <<<<<< upstream data received <<<<<<
HTTP/1.1 200 OK
content-length: 13
content-type: text/plain
date: Thu, 01 Oct 2020 09:16:33 GMT
server: envoy

example body

wasm log: downstream connection close!
wasm log: connection complete!
```
