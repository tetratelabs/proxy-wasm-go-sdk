.DEFAULT_GOAL := build.examples

.PHONY: help build.examples lint test test.sdk test.e2e
help:
	grep -E '^[a-z0-9A-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build.example:
	tinygo build -o ./examples/${name}/main.go.wasm -target=wasm -wasm-abi=generic ./examples/${name}/main.go

build.examples:
	find ./examples -type f -name "main.go" | xargs -Ip tinygo build -o p.wasm -target=wasm -wasm-abi=generic p

lint:
	golangci-lint run --build-tags proxytest

test:
	go test -tags=proxytest -race -v ./...

test.sdk:
	go test -tags=proxytest -race -v ./proxywasm/...

test.e2e:
	go test -tags=proxytest -race -v ./e2e

run:
	getenvoy run wasm:1.15 -- -c ./examples/${name}/envoy.yaml
