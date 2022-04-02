## Requirements

- [Go](https://go.dev/dl/) 1.17 or higher.
- [TinyGo](https://tinygo.org/) - This SDK depends on TinyGo and leverages its [WASI](https://github.com/WebAssembly/WASI) (WebAssembly System Interface) target. Please follow the official instruction [here](https://tinygo.org/getting-started/) for installing TinyGo.
- [Envoy](https://www.envoyproxy.io) - To run compiled examples, you need to have Envoy binary. We recommend using [func-e](https://func-e.io) as the easiest way to get started with Envoy. Alternatively, you can follow [the official instruction](https://www.envoyproxy.io/docs/envoy/latest/start/install).

## Installation

Install `go1.17` or higher
```yaml
tar -C /usr/local/ -xzf go1.17.1.linux-amd64.tar.gz
 
vim /etc/profile.d/go.sh
### 加入以下内容
 
export GOROOT=/usr/local/go
export GOPATH=/data/go
export PATH=$PATH:$GOROOT/bin:$GOPATH
export GO111MODULE="on" # 开启 Go moudles 特性
export GOPROXY=https://goproxy.cn,direct # 安装 Go 模块时，国内代理服务器设置
 
# 让配置生效
source /etc/profile

```
```shell
$ go version
go version go1.17.1 linux/amd64
```

Ubuntu install `tinygo`
```shell
wget https://github.com/tinygo-org/tinygo/releases/download/v0.21.0/tinygo_0.21.0_amd64.deb
sudo dpkg -i tinygo_0.21.0_amd64.deb

export PATH=$PATH:/usr/local/bin
``` 
```shell
$ tinygo version
tinygo version 0.21.0 linux/amd64 (using go version go1.17.1 and LLVM version 11.0.0)
```

Install ``envoy``
```yaml
sudo apt update
sudo apt install apt-transport-https gnupg2 curl lsb-release
curl -sL 'https://deb.dl.getenvoy.io/public/gpg.8115BA8E629CC074.key' | sudo gpg --dearmor -o /usr/share/keyrings/getenvoy-keyring.gpg
Verify the keyring - this should yield "OK"
echo a077cb587a1b622e03aa4bf2f3689de14658a9497a9af2c427bba5f4cc3c4723 /usr/share/keyrings/getenvoy-keyring.gpg | sha256sum --check
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/getenvoy-keyring.gpg] https://deb.dl.getenvoy.io/public/deb/ubuntu $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/getenvoy.list
sudo apt update
sudo apt install -y getenvoy-envoy
```

## Run Examples
```shell
### 编译得到wasm文件
tinygo build -o main.wasm -scheduler=none -target=wasi main.go
### 运行envoy
envoy -c envoy.yaml
### 访问envoy expose listener
curl -v localhost:port ## example/helloworld curl -v localhost:18000
```


