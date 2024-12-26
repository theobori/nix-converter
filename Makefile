# Go files to format
GOFMT_FILES ?= $(shell find . -name "*.go")

all: fmt

.PHONY: fmt
fmt:
	gofmt -w $(GOFMT_FILES)

.PHONY: clean
clean:
	go clean -testcache

.PHONY: test
test: clean
	go test -v ./...
