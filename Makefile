# Minimal Makefile for archive-parser

BINARY_NAME=archive-parser
MAIN_PATH=./cmd

.PHONY: all build test run clean

# Default: build and test
all: build test

# Build the binary
build:
	@echo "Building..."
	go build -o ${BINARY_NAME} ${MAIN_PATH}

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run the program (after building)
run: build
	./${BINARY_NAME}

# Clean artifacts
clean:
	@echo "Cleaning..."
	rm -f ${BINARY_NAME}

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
