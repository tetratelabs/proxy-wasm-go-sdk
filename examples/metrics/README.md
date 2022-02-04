
## metrics

this example creates simple request counter with prometheus tags.

```
$ curl localhost:18000 -v -H "my-custom-header: foo" 

$ curl -s 'localhost:8001/stats/prometheus'| grep proxy
# TYPE custom_header_value_counts counter
custom_header_value_counts{value="foo",reporter="wasmgosdk"} 1
```
