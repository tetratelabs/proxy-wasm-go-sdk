## http_routing

this example proxies http requests and randomly route them to primary/canary clusters by manipulating :authorty header.

```
$ curl localhost:18000
hello from primary!

$ curl localhost:18000
hello from canary!
```