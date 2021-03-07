.DEFAULT_GOAL := build.examples

.PHONY: build.example build.example.docker build.examples build.examples.docker lint test test.sdk test.e2e

build.example:
	tinygo build -o ./examples/${name}/main.go.wasm -scheduler=none -target=wasi ./examples/${name}/main.go

build.example.docker:
	docker run -it -w /tmp/proxy-wasm-go -v $(shell pwd):/tmp/proxy-wasm-go tinygo/tinygo:0.17.0 \
		tinygo build -o /tmp/proxy-wasm-go/examples/${name}/main.go.wasm -scheduler=none -target=wasi \
		/tmp/proxy-wasm-go/examples/${name}/main.go

build.examples:
	find ./examples -type f -name "main.go" | xargs -Ip tinygo build -o p.wasm -scheduler=none -target=wasi p

build.examples.docker:
	docker run -it -w /tmp/proxy-wasm-go -v $(shell pwd):/tmp/proxy-wasm-go tinygo/tinygo:0.17.0 /bin/bash -c \
		'find /tmp/proxy-wasm-go/examples/ -type f -name "main.go" | xargs -Ip tinygo build -o p.wasm -scheduler=none -target=wasi p'

lint:
	golangci-lint run --build-tags proxytest

test:
	go test -tags=proxytest $(shell go list ./... | grep -v e2e | sed 's/github.com\/tetratelabs\/proxy-wasm-go-sdk/./g')

test.e2e:
	go test -v ./e2e

test.e2e.single:
	go test -v ./e2e -run ${name}

run:
	envoy -c ./examples/${name}/envoy.yaml --concurrency 2 --log-format '%v'
