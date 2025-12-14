# Go files to format
__GOFMT_FILES != find . -name "*.go"
GOFMT_FILES ?= ${__GOFMT_FILES}

all: fmt

.PHONY: fmt
fmt:
	gofmt -w ${GOFMT_FILES}

.PHONY: clean
clean:
	go clean -testcache

.PHONY: test
test: clean
	go test -v ./...
