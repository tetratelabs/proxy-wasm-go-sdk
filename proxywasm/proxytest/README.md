##  Test framework for proxy-wasm-go-sdk

Using proxytest, you can test your extension with the official command:

```
go test ./...
```

This framework emulates the expected behavior of Envoyproxy, and you can test your extensions without running Envoy.
For detail, see `examples/*/main_test.go`.


Note that we have not covered all the functionality, and the API is very likely to change in the future.