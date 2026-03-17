.PHONY: build test clean

# Build the yzc compiler binary
build:
	go build -o bin/yzc ./cmd/yzc

# Run all tests
test:
	go test ./...

# Remove build artifacts
clean:
	rm -rf bin/
