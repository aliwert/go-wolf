# Contributing to go-wolf

We welcome contributions to go-wolf! Here's how you can help:

## Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/go-wolf.git`
3. Install dependencies: `go mod tidy`
4. Create a feature branch: `git checkout -b feature/your-feature`

## Guidelines

### Code Style

- Follow Go conventions and use `gofmt`
- Write clear, self-documenting code
- Add comments for exported functions and types
- Use descriptive variable names

### Testing

- Write tests for new features
- Ensure all tests pass: `go test ./...`
- Aim for good test coverage
- Include both unit and integration tests

### Documentation

- Update README.md for new features
- Add examples for new functionality
- Document breaking changes

### Performance

- Profile code for performance-critical paths
- Minimize memory allocations
- Use object pooling where appropriate
- Benchmark new features

## Submitting Changes

1. Write clear commit messages
2. Include tests for new features
3. Update documentation
4. Submit a pull request with description of changes
5. Respond to code review feedback

## Development Commands

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...

# Build examples
go build ./examples/...

# Check formatting
gofmt -l .

# Run linter
golangci-lint run
```

## Areas for Contribution

- [ ] WebSocket support
- [ ] GraphQL integration
- [ ] Database adapters
- [ ] CLI tool for scaffolding
- [ ] OpenAPI/Swagger generation
- [ ] More middleware (JWT, OAuth, etc.)
- [ ] Performance optimizations
- [ ] Documentation improvements
- [ ] Examples and tutorials

## Questions?

- Open an issue for questions
- Check existing issues before creating new ones
- Be respectful and constructive

Thank you for contributing to go-wolf 🐺
