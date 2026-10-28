OUT=telia
SRC=./cmd/compiler
GOLINES_VERSION=v0.12.2
GOLANGCI_LINT_VERSION=v1.64.8
# Pin the tool build compiler too; lint must understand the project's Go version.
TOOLS_GO_VERSION=go1.23.12

.PHONY: all
all: build

.PHONY: build
build:
	go build -o $(OUT) $(SRC)

.PHONY: build-dev
build-dev:
	go build -ldflags="-X main.DevMode=1" -o $(OUT) $(SRC)

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt: fmt-code fmt-lines

.PHONY: fmt-code
fmt-code:
	go fmt ./...

.PHONY: fmt-lines
fmt-lines:
	golines --base-formatter=gofmt -w .

.PHONY: fmt-check
fmt-check:
	python3 scripts/check-format.py

.PHONY: install-format-tools
install-format-tools:
	GOTOOLCHAIN=$(TOOLS_GO_VERSION) go install github.com/segmentio/golines@$(GOLINES_VERSION)

.PHONY: install-lint-tools
install-lint-tools:
	GOTOOLCHAIN=$(TOOLS_GO_VERSION) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: lint
lint:
	golangci-lint run --config .golangci.yml
