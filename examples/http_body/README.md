## http_body

this example demonstrates how to perform operation on a request body like append/prepend/replace.


```
$ curl -XPUT localhost:18000 --data '[initial body]' -H "buffer-operation: prepend"
[this is prepended body][initial body]

$ curl -XPUT localhost:18000 --data '[initial body]' -H "buffer-operation: append"
[initial body][this is appended body]

$ curl -XPUT localhost:18000 --data '[initial body]' -H "buffer-operation: replace"
[this is replaced body]
```
