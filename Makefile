LLVM_GO_TAG=llvm18
OUT=telia
SRC=./cmd/compiler
CGO_CFLAGS=$(shell llvm-config-18 --cflags)
CGO_LDFLAGS=$(shell llvm-config-18 --ldflags)
CGO_CXXFLAGS=$(shell llvm-config-18 --cxxflags)

.PHONY: all
all: build

.PHONY: build
build:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" CGO_CXXFLAGS="$(CGO_CXXFLAGS)" go build -tags=$(LLVM_GO_TAG) -o $(OUT) $(SRC)

.PHONY: build-dev
build-dev:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" CGO_CXXFLAGS="$(CGO_CXXFLAGS)" go build -tags=$(LLVM_GO_TAG) -ldflags="-X main.DevMode=1" -o $(OUT) $(SRC)

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
