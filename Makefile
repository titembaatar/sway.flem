.PHONY: build clean test test-coverage test-verbose run fmt lint help

BINARY_NAME=sway.flem
BUILD_DIR=bin
COVERAGE_DIR=coverage

help:
	@echo "Make targets:"
	@echo "  build         - Build the binary"
	@echo "  clean         - Remove build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-verbose  - Run tests in verbose mode"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linters"
	@echo "  run           - Run application"

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/sway.flem

clean:
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)
	go clean

test:
	go test -v ./...

test-verbose:
	go test -v -count=1 ./...

test-coverage:
	mkdir -p $(COVERAGE_DIR)
	go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated at $(COVERAGE_DIR)/coverage.html"
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out

fmt:
	go fmt ./...

lint:
	go vet ./...

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)
