# Contributing to DeepSeek Go Client

We love your input! We want to make contributing to the DeepSeek Go client as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes
5. Make sure your code lints
6. Issue that pull request!

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/deepseek-go.git
   cd deepseek-go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment:
   ```bash
   cp .env.example .env
   # Edit .env with your DeepSeek API key
   ```

4. Install development tools:
   ```bash
   # Install golangci-lint
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

## Testing

### Running Tests

1. Run all tests:
   ```bash
   go test -v ./...
   ```

2. Run tests in short mode (skips API calls):
   ```bash
   go test -v -short ./...
   ```

3. Run with coverage:
   ```bash
   go test -v -coverprofile=coverage.txt -covermode=atomic ./...
   go tool cover -html=coverage.txt  # View in browser
   ```

### Writing Tests

- Add tests for any new code you write
- Tests should be in `*_test.go` files next to the code they test
- Use table-driven tests when possible
- Test both success and error cases
- Use descriptive test names
- Keep tests focused and independent

## Code Style

We use `golangci-lint` to enforce code style. Run it locally:

```bash
golangci-lint run
```

Key style points:

- Follow standard Go conventions
- Use meaningful variable names
- Add comments for exported functions and types
- Keep functions focused and concise
- Use proper error handling
- Add context to errors when wrapping them

## Pull Request Process

1. Update the README.md with details of changes to the interface
2. Update the CHANGELOG.md with a note describing your changes
3. The PR will be merged once you have the sign-off of at least one maintainer

## Any contributions you make will be under the MIT Software License

In short, when you submit code changes, your submissions are understood to be under the same [MIT License](LICENSE) that covers the project. Feel free to contact the maintainers if that's a concern.

## Report bugs using GitHub's [issue tracker]

We use GitHub issues to track public bugs. Report a bug by [opening a new issue]().

## Write bug reports with detail, background, and sample code

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

## License

By contributing, you agree that your contributions will be licensed under its MIT License.

## References

This document was adapted from the open-source contribution guidelines for [Facebook's Draft](https://github.com/facebook/draft-js/blob/a9316a723f9e918afde44dea68b5f9f39b7d9b00/CONTRIBUTING.md).
