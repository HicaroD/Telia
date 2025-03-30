LLVM_GO_TAG=llvm18
OUT=telia
SRC=./cmd/compiler

.PHONY: all
all: build

.PHONY: build
build:
	go build -tags=$(LLVM_GO_TAG) -o $(OUT) $(SRC)

.PHONY: build-dev
build-dev:
	go build -tags=$(LLVM_GO_TAG) -ldflags="-X main.DevMode=1" -o $(OUT) $(SRC)

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
	golines -w .

.PHONY: lint
lint:
	golangci-lint run
