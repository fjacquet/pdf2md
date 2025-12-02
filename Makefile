.PHONY: all build test benchmark clean fmt

# Binary name
BINARY_NAME=pdf2md

# Build the project
build:
	go build -o $(BINARY_NAME) cmd/pdf2md/main.go

# Run tests
test:
	go test ./...

# Run benchmark
benchmark:
	go run cmd/benchmark/main.go

# Format code
fmt:
	go fmt ./...

# Run linters
lint:
	go vet ./...

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f test.md

# Run the tool on the test file
run-test: build
	./$(BINARY_NAME) internal/testdata/test.pdf test.md
