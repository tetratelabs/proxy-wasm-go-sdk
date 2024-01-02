## http_body_chunk

This example demonstrates how to perform operations on a request body, chunk by chunk.

Reading the received body chunk by chunk, it looks for the string `pattern` inside the body. If it finds it, a 403 response is returned providing the number of the chunk where the pattern was found. Logs are printed every time a chunk is received providing also the size of the read chunk.

Build and run the example:
```bash
$ make build.example name=http_body_chunk
$ make run name=http_body_chunk
```

Perform a request with a body containing the string `pattern`:
```bash
$ head -c 700000 /dev/urandom | base64 > /tmp/file.txt && echo "pattern" >> /tmp/file.txt && curl 'localhost:18000/anything' -d @/tmp/file.txt
pattern found in chunk: 2
```

Generated logs:
```
wasm log: OnHttpRequestBody called. BodySize: 114532, totalRequestBodyReadSize: 0, endOfStream: false
wasm log: read chunk size: 114532
wasm log: OnHttpRequestBody called. BodySize: 114532, totalRequestBodyReadSize: 114532, endOfStream: false
wasm log: OnHttpRequestBody called. BodySize: 933343, totalRequestBodyReadSize: 114532, endOfStream: true
wasm log: read chunk size: 818811
wasm log: pattern found in chunk: 2
wasm log: local 403 response sent
```
