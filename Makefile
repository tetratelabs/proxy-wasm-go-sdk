goimports := golang.org/x/tools/cmd/goimports@v0.1.10
golangci_lint := github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2

.PHONY: build.example
build.example:
	@find ./examples -type f -name "main.go" | grep ${name}\
	| xargs -I {} bash -c 'dirname {}' \
	| xargs -I {} bash -c 'cd {} && tinygo build -o main.wasm -scheduler=none -target=wasi ./main.go'

.PHONY: build.examples
build.examples:
	@find ./examples -mindepth 1 -type f -name "main.go" \
	| xargs -I {} bash -c 'dirname {}' \
	| xargs -I {} bash -c 'cd {} && tinygo build -o main.wasm -scheduler=none -target=wasi ./main.go'

.PHONY: test
test:
	@# First we test the main module because the iteration through the modules
	@# in the lines above is very inconvenient when the folder is ".".
	@go test -tags=proxytest $(shell go list ./... | grep -v e2e)

	@# Now we go through the examples.
	@find . -name "go.mod" \
	| grep -v "\.\/go\.mod" \
	| xargs -I {} bash -c 'dirname {}' \
	| xargs -I {} bash -c 'cd {}; go test -tags=proxytest ./...'

.PHONY: test.e2e
test.e2e:
	@go test -v ./e2e -count=1

.PHONY: test.e2e.single
test.e2e.single:
	@go test -v ./e2e -run '/${name}' -count=1

.PHONY: run
run:
	@envoy -c ./examples/${name}/envoy.yaml --concurrency 2 --log-format '%v'

.PHONY: lint
lint:
	@go run $(golangci_lint) run --build-tags proxytest

.PHONY: format
format:
	@find . -type f -name '*.go' | xargs gofmt -s -w
	@for f in `find . -name '*.go'`; do \
	    awk '/^import \($$/,/^\)$$/{if($$0=="")next}{print}' $$f > /tmp/fmt; \
	    mv /tmp/fmt $$f; \
	done
	@go run $(goimports) -w -local github.com/tetratelabs/proxy-wasm-go-sdk `find . -name '*.go'`

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
		docker build -t $$ref . -f examples/wasm-image.Dockerfile --build-arg WASM_BINARY_PATH=$$(dirname $$f)/main.wasm; \
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
		buildah bud -f examples/wasm-image.Dockerfile --build-arg WASM_BINARY_PATH=$$(dirname $$f)/main.wasm -t $$ref .; \
		buildah push $$ref; \
	done

.PHONY: tidy
tidy: ## Runs go mod tidy on every module
	@find . -name "go.mod" \
	| grep go.mod \
	| xargs -I {} bash -c 'dirname {}' \
	| xargs -I {} bash -c 'echo "=> {}"; cd {}; go mod tidy -v; '
