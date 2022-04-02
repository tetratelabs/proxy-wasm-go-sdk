## http_body

此示例演示如何对请求正文(requestBody)执行操作，如追加/前置/替换(append/prepend/replace)。
```
$ curl -XPUT localhost:18000 --data '[initial body]' -H "buffer-operation: prepend"
[this is prepended body][initial body]

$ curl -XPUT localhost:18000 --data '[initial body]' -H "buffer-operation: append"
[initial body][this is appended body]

$ curl -XPUT localhost:18000 --data '[initial body]' -H "buffer-operation: replace"
[this is replaced body]
```

envoy 日志

连续执行以下两个命令
```shell
$ curl -XPUT localhost:18000 --data '[initial body]' -H "buffer-operation: prepend"
[this is prepended body][initial body]
$ curl -XPUT localhost:18000 --data '[initial body]' -H "buffer-operation: aaaa"
[this is replaced body]root
```

```shell
[2022-04-02 06:52:48.275][113934][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: <---- NewHttpContext ----> 
[2022-04-02 06:52:48.275][113934][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: <---- OnHttpRequestHeaders ---->
[2022-04-02 06:52:48.275][113934][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: <---- setBodyContest OnHttpRequestBody ---->
[2022-04-02 06:52:48.275][113934][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: <---- setBodyContest OnHttpRequestBody ---->
[2022-04-02 06:52:48.275][113934][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: original request body: [initial body]
[2022-04-02 06:52:48.275][113935][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: <---- NewHttpContext ----> 
[2022-04-02 06:52:48.276][113935][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: <---- echoBodyContest OnHttpRequestBody ---->
[2022-04-02 06:52:48.276][113935][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: <---- echoBodyContest OnHttpRequestBody ---->
[2022-04-02 06:53:06.477][113935][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: <---- NewHttpContext ----> 
[2022-04-02 06:53:06.477][113935][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: <---- OnHttpRequestHeaders ---->
[2022-04-02 06:53:06.477][113935][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: <---- setBodyContest OnHttpRequestBody ---->
[2022-04-02 06:53:06.477][113935][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: <---- setBodyContest OnHttpRequestBody ---->
[2022-04-02 06:53:06.477][113935][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: original request body: [initial body]
[2022-04-02 06:53:06.477][113932][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: <---- NewHttpContext ----> 
[2022-04-02 06:53:06.477][113932][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: <---- echoBodyContest OnHttpRequestBody ---->
[2022-04-02 06:53:06.477][113932][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: <---- echoBodyContest OnHttpRequestBody ---->
```
从上面的结果简单可知，当 执行 ``curl -XPUT localhost:18000 --data '[initial body]' -H "buffer-operation: prepend"``
的时候,也就是有一次http请求的时候，envoy调用wasm的流程大致是如下的：

- 创建一个 httpContext ，也就是一个 http的上下文对象
- 执行 request header 相关的逻辑。这应该是因为我们请求中带有 header相关信息，所以执行了它
- 执行了两次setBody RequestBody方法，第二次才真正执行了增强的逻辑。为什么是两次？暂时未知，官方代码的注释是说 “这里可能会执行两次，直到获取到了请求流”
  原文： 
  > Note that this is potentially called multiple times until we see end_of_stream = true.
- 执行echo setBody RequestBody方法








