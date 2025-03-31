.PHONY: build clean test run fmt lint help

BINARY_NAME=sway.flem
BUILD_DIR=bin

help:
	@echo "Make targets:"
	@echo "  build    - Build the binary"
	@echo "  clean    - Remove build artifacts"
	@echo "  test     - Run tests"
	@echo "  fmt      - Format code"
	@echo "  lint     - Run linters"
	@echo "  run      - Run application"

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/sway.flem

clean:
	rm -rf $(BUILD_DIR)
	go clean

test:
	go test -v ./...

fmt:
	go fmt ./...

lint:
	go vet ./...

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)
