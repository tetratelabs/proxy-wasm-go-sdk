## http_body

this example demonstrates how to perform operation on a request or response body like append/prepend/replace.

To modify the request:
```
$ curl -XPUT localhost:18000 --data '[initial body]' -H "buffer-operation: prepend"
[this is prepended body][initial body]

$ curl -XPUT localhost:18000 --data '[initial body]' -H "buffer-operation: append"
[initial body][this is appended body]

$ curl -XPUT localhost:18000 --data '[initial body]' -H "buffer-operation: replace"
[this is replaced body]
```

To modify the response:
```
$ curl -XPUT localhost:18000 --data '[initial body]' -H "buffer-operation: prepend" -H "buffer-replace-at: response"
[this is prepended body][initial body]

$ curl -XPUT localhost:18000 --data '[initial body]' -H "buffer-operation: append" -H "buffer-replace-at: response"
[initial body][this is appended body]

$ curl -XPUT localhost:18000 --data '[initial body]' -H "buffer-operation: replace" -H "buffer-replace-at: response"
[this is replaced body]
```
