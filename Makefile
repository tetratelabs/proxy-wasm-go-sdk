.DEFAULT_GOAL := build.examples

ISTIO_VERSION ?= 1.7.3

.PHONY: build.example build.example.docker build.examples build.examples.docker lint test test.sdk test.e2e

build.example:
	tinygo build -o ./examples/${name}/main.go.wasm -scheduler=none -target=wasi -wasm-abi=generic ./examples/${name}/main.go

build.example.docker:
	docker run -it -w /tmp/proxy-wasm-go -v $(shell pwd):/tmp/proxy-wasm-go tinygo/tinygo-dev:latest \
		tinygo build -o /tmp/proxy-wasm-go/examples/${name}/main.go.wasm -scheduler=none -target=wasi \
		-wasm-abi=generic /tmp/proxy-wasm-go/examples/${name}/main.go

build.examples:
	find ./examples -type f -name "main.go" | xargs -Ip tinygo build -o p.wasm -scheduler=none -target=wasi -wasm-abi=generic p

build.examples.docker:
	docker run -it -w /tmp/proxy-wasm-go -v $(shell pwd):/tmp/proxy-wasm-go tinygo/tinygo:latest /bin/bash -c \
		'find /tmp/proxy-wasm-go/examples/ -type f -name "main.go" | xargs -Ip tinygo build -o p.wasm -scheduler=none -target=wasi -wasm-abi=generic p'

lint:
	golangci-lint run --build-tags proxytest

test:
	go test -tags=proxytest $(shell go list ./... | grep -v e2e | sed 's/github.com\/tetratelabs\/proxy-wasm-go-sdk/./g')

test.e2e:
	docker run -it -w /tmp/proxy-wasm-go -v $(shell pwd):/tmp/proxy-wasm-go getenvoy/proxy-wasm-go-sdk-ci:istio-${ISTIO_VERSION} go test -v ./e2e

test.e2e.single:
	docker run -it -w /tmp/proxy-wasm-go -v $(shell pwd):/tmp/proxy-wasm-go getenvoy/proxy-wasm-go-sdk-ci:istio-${ISTIO_VERSION} go test -v ./e2e -run ${name}

run:
	docker run --entrypoint='/usr/local/bin/envoy' \
		-p 18000:18000 -p 8099:8099 \
		-w /tmp/envoy -v $(shell pwd):/tmp/envoy getenvoy/proxy-wasm-go-sdk-ci:istio-${ISTIO_VERSION} \
		-c /tmp/envoy/examples/${name}/envoy.yaml --concurrency 2 \
		--log-format-prefix-with-location '0' --log-format '%v' # --log-format-prefix-with-location will be removed at 1.17.0 release
