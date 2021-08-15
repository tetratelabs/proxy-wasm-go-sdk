# bingo manages go binaries needed for building the project
include .bingo/Variables.mk

.PHONY: build.example
build.example:
	@find ./examples -type f -name "main.go" | grep ${name} | xargs -Ip tinygo build -o p.wasm -scheduler=none -target=wasi p

.PHONY: build.examples
build.examples:
	@find ./examples -type f -name "main.go" | xargs -Ip tinygo build -o p.wasm -scheduler=none -target=wasi p

.PHONY: test
test:
	go test -tags=proxytest $(shell go list ./... | grep -v e2e)

.PHONY: test.e2e
test.e2e:
	@cd ./e2e && go test -v . -count=1

.PHONY: test.e2e.single
test.e2e.single:
	@cd ./e2e && go test -v . -run '${name}' -count=1

.PHONY: test.e2e.loadtest
test.e2e.loadtest:
	@cd ./e2e && go test -v ./loadtest -count=1

.PHONY: run
run:
	envoy -c ./examples/${name}/envoy.yaml --concurrency 2 --log-format '%v'

.PHONY: lint
lint: $(GOLANGCI_LINT)
	@$(GOLANGCI_LINT) run --build-tags proxytest

.PHONY: format
format: $(GOIMPORTS)
	@find . -type f -name '*.go' | xargs gofmt -s -w
	@for f in `find . -name '*.go'`; do \
	    awk '/^import \($$/,/^\)$$/{if($$0=="")next}{print}' $$f > /tmp/fmt; \
	    mv /tmp/fmt $$f; \
	    $(GOIMPORTS) -w -local github.com/tetratelabs/proxy-wasm-go-sdk $$f; \
	done

.PHONY: check
check:
	@$(MAKE) format
	@go mod tidy
	@if [ ! -z "`git status -s`" ]; then \
		echo "The following differences will fail CI until committed:"; \
		git diff --exit-code; \
	fi

# Build docker images of *compat* variant of Wasm Image Specification with built example binaries,
# and push to ghcr.io/tetratelabs/proxy-wasm-go-sdk-examples.
# See https://github.com/solo-io/wasm/blob/master/spec/spec-compat.md for details.
# Only-used in github workflow on the main branch, and not for developers.
.PHONY: wasm_image.build_push
wasm_image.build_push:
	@for f in `find ./examples -type f -name "main.go"`; do \
		name=`echo $$f | sed -e 's/\\//-/g' | sed -e 's/\.-examples-//g' -e 's/\-main\.go//g'` ; \
		ref=ghcr.io/tetratelabs/proxy-wasm-go-sdk-examples:$$name; \
		docker build -t $$ref . -f examples/wasm-image.Dockerfile --build-arg WASM_BINARY_PATH=$$f.wasm; \
		docker push $$ref; \
	done

# Build OCI images of *compat* variant of Wasm Image Specification with built example binaries,
# and push to ghcr.io/tetratelabs/proxy-wasm-go-sdk-examples.
# See https://github.com/solo-io/wasm/blob/master/spec/spec-compat.md for details.
# Only-used in github workflow on the main branch, and not for developers.
# Requires "buildah" CLI.
.PHONY: wasm_image.build_push_oci
wasm_image.build_push_oci:
	@for f in `find ./examples -type f -name "main.go"`; do \
		name=`echo $$f | sed -e 's/\\//-/g' | sed -e 's/\.-examples-//g' -e 's/\-main\.go//g'` ; \
		ref=ghcr.io/tetratelabs/proxy-wasm-go-sdk-examples:$$name-oci; \
		buildah bud -f examples/wasm-image.Dockerfile --build-arg WASM_BINARY_PATH=$$f.wasm -t $$ref .; \
		buildah push $$ref; \
	done
