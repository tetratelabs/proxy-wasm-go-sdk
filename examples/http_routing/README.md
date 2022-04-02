## http_routing

this example proxies http requests and randomly route them to primary/canary clusters by manipulating :authorty header.

此示例代理 http 请求并通过操作 :authorty 标头将它们随机路由到主/金丝雀集群

### message on clients
```shell
$ curl localhost:18000
hello from primary!
$ curl localhost:18000
hello from primary!
$ curl localhost:18000
hello from canary!
$ curl localhost:18000
hello from canary!
$ curl localhost:18000
hello from primary!
$ curl localhost:18000
hello from canary!
$ curl localhost:18000
hello from canary!
$ curl localhost:18000
hello from primary!
$ curl localhost:18000
hello from canary!
$ curl localhost:18000
hello from primary!
$ curl localhost:18000
hello from primary!
$ curl localhost:18000
hello from primary!
$ curl localhost:18000
hello from canary!

```

### message on Envoy

```shell
wasm log:  <---- New HttpContext ----> 
wasm log:  <---- OnHttpRequestHeaders ----> 
wasm log: value: 219793847

wasm log:  <---- New HttpContext ----> 
wasm log:  <---- OnHttpRequestHeaders ----> 
wasm log: value: 920072585

wasm log:  <---- New HttpContext ----> 
wasm log:  <---- OnHttpRequestHeaders ----> 
wasm log: value: 3241992484

wasm log:  <---- New HttpContext ----> 
wasm log:  <---- OnHttpRequestHeaders ----> 
wasm log: value: 2532829830
...... 
```








