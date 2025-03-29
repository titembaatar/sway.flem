# Makefile for sway.flem

# Variables
BINARY_NAME=swayflem
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags="-s -w"
BUILD_DIR=build

# Targets
.PHONY: all build clean install uninstall

all: build

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/sway-manager

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@$(GO) clean

install: build
	@echo "Installing..."
	@install -D -m 755 $(BUILD_DIR)/$(BINARY_NAME) $(DESTDIR)/usr/local/bin/$(BINARY_NAME)

uninstall:
	@echo "Uninstalling..."
	@rm -f $(DESTDIR)/usr/local/bin/$(BINARY_NAME)

test:
	@echo "Running tests..."
	@$(GO) test ./...

lint:
	@echo "Running linter..."
	@golangci-lint run
