# bingo manages go binaries needed for building the project
include .bingo/Variables.mk

.PHONY: build.example build.example.docker build.examples build.examples.docker lint test test.sdk test.e2e format check

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

test:
	go test -tags=proxytest $(shell go list ./... | grep -v e2e | sed 's/github.com\/tetratelabs\/proxy-wasm-go-sdk/./g')

test.e2e:
	go test -v ./e2e

test.e2e.single:
	go test -v ./e2e -run ${name}

run:
	envoy -c ./examples/${name}/envoy.yaml --concurrency 2 --log-format '%v'

lint: $(GOLANGCI_LINT)
	@golangci-lint run --build-tags proxytest

format: $(GOIMPORTS)
	@find . -type f -name '*.go' | xargs gofmt -s -w
	@for f in `find . -name '*.go'`; do \
	    awk '/^import \($$/,/^\)$$/{if($$0=="")next}{print}' $$f > /tmp/fmt; \
	    mv /tmp/fmt $$f; \
	    goimports -w -local github.com/tetratelabs/proxy-wasm-go-sdk $$f; \
	done

check:
	@$(MAKE) format
	@go mod tidy
	@if [ ! -z "`git status -s`" ]; then \
		echo "The following differences will fail CI until committed:"; \
		git diff --exit-code; \
	fi
