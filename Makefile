.PHONY: all
all: generate test

.PHONY: generate
generate:
	@echo "Generating code from OpenAPI spec..."
	@go run cmd/generate/main.go

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v -race -cover ./...

.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	@go test -v -race -tags=integration ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

.PHONY: vet
vet:
	@echo "Running go vet..."
	@go vet ./...

.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

.PHONY: clean
clean:
	@echo "Cleaning generated files..."
	@rm -rf pkg/client/generated/
	@rm -f pkg/models/generated*.go
	@rm -f coverage.txt coverage.html
	@echo "Clean complete"

.PHONY: build-examples
build-examples:
	@echo "Building examples..."
	@go build -o bin/basic examples/basic/main.go
	@go build -o bin/deployment examples/deployment/main.go
	@go build -o bin/monitoring examples/monitoring/main.go
	@go build -o bin/pagination examples/pagination/main.go
	@echo "Examples built in bin/"

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all              - Generate code and run tests"
	@echo "  generate         - Generate code from OpenAPI spec"
	@echo "  test             - Run unit tests"
	@echo "  test-integration - Run integration tests (requires running Prefect server)"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  lint             - Run golangci-lint"
	@echo "  fmt              - Format code"
	@echo "  vet              - Run go vet"
	@echo "  tidy             - Tidy dependencies"
	@echo "  clean            - Remove generated files"
	@echo "  build-examples   - Build example programs"
	@echo "  help             - Show this help message"
