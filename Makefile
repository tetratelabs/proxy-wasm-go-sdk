# bingo manages go binaries needed for building the project
include .bingo/Variables.mk

.PHONY: build.example build.example.docker build.examples build.examples.docker lint test test.sdk test.e2e format check

build.example:
	find ./examples -type f -name "main.go" | grep ${name} | xargs -Ip tinygo build -o p.wasm -scheduler=none -target=wasi p

build.examples:
	find ./examples -type f -name "main.go" | xargs -Ip tinygo build -o p.wasm -scheduler=none -target=wasi p

test:
	go test -tags=proxytest $(shell go list ./... | grep -v e2e)

test.e2e:
	go test -v ./e2e

test.e2e.single:
	go test -v ./e2e -run '/${name}'

run:
	envoy -c ./examples/${name}/envoy.yaml --concurrency 2 --log-format '%v'

lint: $(GOLANGCI_LINT)
	@$(GOLANGCI_LINT) run --build-tags proxytest

format: $(GOIMPORTS)
	@find . -type f -name '*.go' | xargs gofmt -s -w
	@for f in `find . -name '*.go'`; do \
	    awk '/^import \($$/,/^\)$$/{if($$0=="")next}{print}' $$f > /tmp/fmt; \
	    mv /tmp/fmt $$f; \
	    $(GOIMPORTS) -w -local github.com/tetratelabs/proxy-wasm-go-sdk $$f; \
	done

check:
	@$(MAKE) format
	@go mod tidy
	@if [ ! -z "`git status -s`" ]; then \
		echo "The following differences will fail CI until committed:"; \
		git diff --exit-code; \
	fi
