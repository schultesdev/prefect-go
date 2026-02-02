# Contributing to prefect-go

Thank you for your interest in contributing to prefect-go!

## Development Setup

### Prerequisites

- Go 1.21 or later
- A running Prefect server (for integration tests)

### Clone the Repository

```bash
git clone https://git.schultes.dev/schultesdev/prefect-go.git
cd prefect-go
```

### Install Dependencies

```bash
go mod download
```

## Development Workflow

### Running Tests

```bash
# Run unit tests
make test

# Run integration tests (requires running Prefect server)
make test-integration

# Generate coverage report
make test-coverage
```

### Code Formatting

```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet
```

### Building Examples

```bash
make build-examples
```

## Project Structure

```
prefect-go/
├── api/                    # OpenAPI specifications
├── cmd/                    # Command-line tools
│   └── generate/          # Code generator
├── examples/              # Example programs
│   ├── basic/            # Basic usage
│   ├── deployment/       # Deployment examples
│   ├── monitoring/       # Monitoring examples
│   └── pagination/       # Pagination examples
├── pkg/                   # Library code
│   ├── client/           # Main client and services
│   ├── errors/           # Error types
│   ├── models/           # Data models
│   ├── pagination/       # Pagination support
│   └── retry/            # Retry logic
└── tools/                 # Build tools
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Write tests for new functionality
- Add GoDoc comments for exported types and functions
- Keep functions focused and small
- Use meaningful variable names

## Commit Messages

- Use clear, descriptive commit messages
- Start with a verb in present tense (e.g., "Add", "Fix", "Update")
- Reference issue numbers when applicable

Example:
```
Add support for work pool queues

Implements work pool queue operations including create, get, and list.
Closes #123
```

## Pull Request Process

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Make** your changes
4. **Add** tests for new functionality
5. **Ensure** all tests pass (`make test`)
6. **Run** code formatting (`make fmt`)
7. **Commit** your changes with a clear message
8. **Push** to your fork (`git push origin feature/amazing-feature`)
9. **Open** a Pull Request

### Pull Request Guidelines

- Provide a clear description of the changes
- Reference any related issues
- Ensure CI checks pass
- Request review from maintainers
- Be responsive to feedback

## Adding New Features

### Adding a New Service

1. Add the service struct in `pkg/client/services.go` or a new file
2. Add the service to the `Client` struct in `pkg/client/client.go`
3. Initialize the service in `NewClient()`
4. Implement service methods following existing patterns
5. Add tests for the new service
6. Update documentation

Example:
```go
// NewService represents operations for new resources
type NewService struct {
    client *Client
}

func (s *NewService) Get(ctx context.Context, id uuid.UUID) (*models.NewResource, error) {
    var resource models.NewResource
    path := joinPath("/new_resources", id.String())
    if err := s.client.get(ctx, path, &resource); err != nil {
        return nil, fmt.Errorf("failed to get resource: %w", err)
    }
    return &resource, nil
}
```

### Adding New Models

1. Add model structs in `pkg/models/models.go`
2. Follow existing naming conventions (`Resource`, `ResourceCreate`, `ResourceUpdate`)
3. Use proper JSON tags
4. Add documentation comments

## Testing

### Unit Tests

- Write unit tests for all new functions
- Use table-driven tests when appropriate
- Mock external dependencies
- Aim for >80% code coverage

### Integration Tests

- Add integration tests for API operations
- Use the `integration` build tag
- Clean up resources after tests
- Handle test failures gracefully

## Documentation

### GoDoc

- Add package-level documentation
- Document all exported types and functions
- Include examples in documentation when helpful
- Use proper formatting (code blocks, links, etc.)

### Examples

- Add examples for new features in `examples/`
- Make examples runnable and self-contained
- Include error handling
- Add comments explaining key concepts

## Questions?

Feel free to open an issue if you have questions or need clarification.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
