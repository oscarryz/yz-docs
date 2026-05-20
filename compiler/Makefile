.PHONY: build test test-race clean

# Build the yzc compiler binary
VERSION := $(shell git tag --sort=-version:refname 2>/dev/null | head -1 || echo "v0.0.0")-$(shell git rev-parse --short HEAD)

build:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/yzc ./cmd/yzc

# Run all tests
test:
	go test ./...

# Run all tests with the race detector (use before committing concurrency changes)
test-race:
	go test -race -count=1 ./...

# Remove build artifacts
clean:
	rm -rf bin/
