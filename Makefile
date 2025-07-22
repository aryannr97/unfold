.PHONY: lint test build clean help

# Default target
help:
	@echo "Available targets:"
	@echo "  lint      - Run golangci-lint"
	@echo "  test      - Run tests"
	@echo "  build     - Build the application"
	@echo "  clean     - Clean build artifacts"

# Run linter
lint:
	@echo "+ $@"
	@$(golangci-lint run)

# Run tests
test:
	@echo "+ $@"
	@$(go test -v ./...)

# Run tests with coverage
test-coverage:
	@echo "+ $@"
	@$(go test -v -coverprofile=coverage.out ./...)

# Build the application
build:
	@echo "+ $@"
	@$(go build -v ./...)

# Clean build artifacts
clean:
	@echo "+ $@"
	@$(go clean)
	@$(rm -f coverage.out)