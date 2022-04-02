## helloworld

首先，编译此project，生成`wasm`文件
```shell
tinygo build -o main.wasm -scheduler=none -target=wasi main.go
```
然后，运行`envoy`
```shell
envoy -c envoy.yaml
```

envoy的日志输出大致如下

```
wasm log: OnPluginStart from Go!
[2022-04-02 02:25:37.016][33107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: OnPluginStart from Go!
[2022-04-02 02:25:37.016][33109][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: OnPluginStart from Go!
[2022-04-02 02:25:37.016][33111][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: OnPluginStart from Go!
[2022-04-02 02:25:38.010][33096][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: It's 1648866338010816000: random value: 17149606892496637177
[2022-04-02 02:25:38.010][33096][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: OnTick called
[2022-04-02 02:25:38.018][33107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: It's 1648866338018833000: random value: 14386222021184593187
[2022-04-02 02:25:38.018][33106][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: It's 1648866338018846000: random value: 13568281049617375916
[2022-04-02 02:25:38.018][33106][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: OnTick called
[2022-04-02 02:25:38.018][33107][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: OnTick called
[2022-04-02 02:25:38.018][33109][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: It's 1648866338018853000: random value: 7393372814440656961
[2022-04-02 02:25:38.019][33109][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: OnTick called
[2022-04-02 02:25:38.018][33111][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: It's 1648866338018854000: random value: 468198035833561196
[2022-04-02 02:25:38.019][33111][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: OnTick called

```
