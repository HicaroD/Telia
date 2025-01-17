LLVM_GO_TAG=llvm18
OUT=telia

# Default target
.PHONY: all
all: build

.PHONY: build
build:
	go build -tags=$(LLVM_GO_TAG) -o $(OUT)

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	go fmt ./...
