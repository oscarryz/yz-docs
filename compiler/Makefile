.PHONY: build test test-race clean

# Build the yzc compiler binary
build:
	go build -o bin/yzc ./cmd/yzc

# Run all tests
test:
	go test ./...

# Run all tests with the race detector (use before committing concurrency changes)
test-race:
	go test -race -count=1 ./...

# Remove build artifacts
clean:
	rm -rf bin/
