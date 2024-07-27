goimports := golang.org/x/tools/cmd/goimports@v0.21.0
golangci_lint := github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.0


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
	@go test $(shell go list ./... | grep -v e2e)
	@go test -tags "proxywasm_timing" ./proxywasm/proxytest

.PHONY: test.examples
test.examples:
	@find ./examples -mindepth 1 -type f -name "main.go" \
	| xargs -I {} bash -c 'dirname {}' \
	| xargs -I {} bash -c 'cd {} && go test ./...'

.PHONY: run
run:
	@envoy -c ./examples/${name}/envoy.yaml --concurrency 2 --log-format '%v'

.PHONY: lint
lint:
	@find . -name "go.mod" \
	| grep go.mod \
	| xargs -I {} bash -c 'dirname {}' \
	| xargs -I {} bash -c 'echo "=> {}"; cd {}; go run $(golangci_lint) run; '

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

.PHONY: tidy
tidy: ## Runs go mod tidy on every module
	@find . -name "go.mod" \
	| grep go.mod \
	| xargs -I {} bash -c 'dirname {}' \
	| xargs -I {} bash -c 'echo "=> {}"; cd {}; go mod tidy -v; '
