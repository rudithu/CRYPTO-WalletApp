# Project metadata
APP_NAME := CRYPTO-WalletApp2
PKG := ./test/...
CMD := ./main.go

# Go environment
GO := go

# Go tools
LINT := golangci-lint run

.PHONY: all build run test lint clean

# Default target
all: build

# Build the binary
build:
	$(GO) build -o $(APP_NAME) $(CMD)

# Run the app
run:
	$(GO) run $(CMD)

# Run unit tests
test:
	$(GO) test -v $(PKG)

# Lint the code (requires golangci-lint installed)
lint:
	$(LINT)

# Clean up binaries
clean:
	@rm -f $(APP_NAME)
