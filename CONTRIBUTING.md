# 💡 Contributing Guide

Thank you for considering contributing to TempFiles! This guide will help you get started and ensure smooth collaboration.

## 🚀 Quick Start for Contributors

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/YOUR_USERNAME/tempfile.git`
3. **Create** a feature branch: `git checkout -b feature/amazing-feature`
4. **Make** your changes
5. **Test** your changes: `go test ./...`
6. **Commit** using conventional format: `git commit -m "feat: add amazing feature"`
7. **Push** to your branch: `git push origin feature/amazing-feature`
8. **Create** a Pull Request

## 📝 Conventional Commits Guide

We use [Conventional Commits](https://www.conventionalcommits.org/) for automated semantic versioning. Here's how to format your commit messages:

### Format
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Description | Version Bump |
|------|-------------|--------------|
| `feat` | New feature | MINOR |
| `fix` | Bug fix | PATCH |
| `docs` | Documentation only | PATCH |
| `style` | Code style (formatting, etc) | PATCH |
| `refactor` | Code refactoring | PATCH |
| `perf` | Performance improvement | PATCH |
| `test` | Adding/updating tests | PATCH |
| `build` | Build system changes | PATCH |
| `ci` | CI/CD changes | PATCH |
| `chore` | Maintenance tasks | PATCH |
| `revert` | Revert previous commit | PATCH |

### Breaking Changes
Add `!` after type or `BREAKING CHANGE:` in footer for MAJOR version bump:
```bash
feat!: change upload API format
# or
feat: change upload API

BREAKING CHANGE: Upload endpoint now requires multipart/form-data
```

### Examples

#### ✅ Good Commit Messages
```bash
# Features
feat: add drag and drop file upload
feat(ui): implement dark theme toggle
feat(api): add batch file upload endpoint
feat(ratelimit): implement Redis-based rate limiting
feat(middleware): add IP whitelisting support

# Bug fixes
fix: resolve memory leak in cleanup service
fix(upload): handle large file uploads correctly
fix(ui): fix responsive layout on mobile
fix(ratelimit): handle rate limit edge cases correctly

# Documentation
docs: add API usage examples
docs(readme): update installation instructions
docs(deployment): add rate limiting configuration guide

# Refactoring
refactor: extract upload logic to service layer
refactor(handlers): simplify error handling
refactor(ratelimit): optimize sliding window algorithm

# Tests
test: add unit tests for file cleanup
test(api): add integration tests for upload endpoint
test(ratelimit): add comprehensive rate limiter tests

# Breaking changes
feat!: change API response format

BREAKING CHANGE: All API responses now include metadata object
```

#### ❌ Bad Commit Messages
```bash
# Too vague
fix: bug
update readme
changes

# Not following format
Fix: upload bug (wrong capitalization)
added new feature (missing type)
FIX: Bug in upload (all caps)
```

## 🧪 Testing Requirements

Before submitting a PR, ensure:

1. **All tests pass**: `go test ./...`
2. **Code is formatted**: `go fmt ./...`
3. **No linting errors**: `golangci-lint run` (if installed)
4. **Build succeeds**: `go build ./cmd/server`

### Writing Tests

- Add tests for new features
- Maintain or improve test coverage
- Use table-driven tests for multiple test cases
- Mock external dependencies

```go
func TestUploadHandler(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected int
    }{
        {"valid file", "test.txt", 200},
        {"empty file", "", 400},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## 🏗️ Development Setup

### Local Development

```bash
# Clone repository
git clone https://github.com/pandeptwidyaop/tempfile.git
cd tempfile

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Optional: Start Redis for rate limiting testing
docker run -d --name redis-dev -p 6379:6379 redis:7-alpine

# Run in development mode
go run cmd/server/main.go

# Or with live reload (install air first)
air
```

### Testing Rate Limiting Features

```bash
# Test with memory store (default)
ENABLE_RATE_LIMIT=true go run cmd/server/main.go

# Test with Redis store
ENABLE_RATE_LIMIT=true RATE_LIMIT_STORE=redis go run cmd/server/main.go

# Run rate limiting tests
go test ./internal/ratelimit/... -v

# Run integration tests
go test ./internal/middleware/... -v
```

### Development Tools

**Recommended tools:**
- [Air](https://github.com/cosmtrek/air) - Live reload for Go
- [golangci-lint](https://golangci-lint.run/) - Go linter
- [gosec](https://github.com/securecodewarrior/gosec) - Security scanner

```bash
# Install development tools
go install github.com/cosmtrek/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
```

## 📋 Code Style Guidelines

### Go Code Standards

1. **Follow Go conventions**:
   - Use `gofmt` for formatting
   - Follow naming conventions (`camelCase` for private, `PascalCase` for public)
   - Write clear, self-documenting code

2. **Error handling**:
   ```go
   // Good
   if err != nil {
       return fmt.Errorf("failed to upload file: %w", err)
   }
   
   // Avoid
   if err != nil {
       panic(err)
   }
   ```

3. **Comments**:
   ```go
   // UploadHandler handles file uploads and returns download URL
   func UploadHandler(c *fiber.Ctx) error {
       // Implementation
   }
   ```

### Project Structure

```
internal/
├── config/          # Configuration management
├── handlers/        # HTTP route handlers
├── middleware/      # HTTP middleware (rate limiting, etc.)
├── models/          # Data structures
├── ratelimit/       # Rate limiting implementation
├── services/        # Business logic
└── utils/           # Utility functions
```

- Keep business logic in `services/`
- HTTP-specific code in `handlers/`
- Middleware components in `middleware/`
- Rate limiting logic in `ratelimit/`
- Shared utilities in `utils/`
- Configuration in `config/`

### Rate Limiting Architecture

The rate limiting system follows a modular design:

```
ratelimit/
├── interface.go     # Core interfaces and types
├── limiter.go       # Main rate limiter implementation
├── memory.go        # In-memory storage backend
├── redis.go         # Redis storage backend
├── ipdetector.go    # IP detection and whitelisting
└── errors.go        # Rate limiting specific errors
```

**Key principles:**
- **Interface-based design** - Easy to swap storage backends
- **Thread-safe operations** - Safe for concurrent use
- **Configurable limits** - Per-endpoint and global limits
- **IP detection** - Handles reverse proxy scenarios
- **Graceful degradation** - Falls back when storage unavailable

## 🔍 Pull Request Guidelines

### Before Submitting

- [ ] Tests pass locally
- [ ] Code follows style guidelines  
- [ ] Commit messages follow conventional format
- [ ] Documentation updated (if needed)
- [ ] No sensitive information in code

### PR Description Template

```markdown
## 📝 Description
Brief description of changes

## 🎯 Type of Change
- [ ] Bug fix (non-breaking change)
- [ ] New feature (non-breaking change)  
- [ ] Breaking change (fix/feature causing existing functionality to change)
- [ ] Documentation update

## 🧪 Testing
- [ ] Tests pass
- [ ] New tests added (for features)
- [ ] Manual testing completed

## 📋 Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
```

### Review Process

1. **Automated checks** run first (tests, linting, security)
2. **Code review** by maintainers
3. **Feedback** addressed by contributor  
4. **Approval** and merge to main
5. **Automatic release** (if merged to main)

## 🐛 Reporting Issues

### Bug Reports

Use the bug report template:

```markdown
**Describe the bug**
Clear description of the issue

**To Reproduce**
Steps to reproduce:
1. Go to '...'
2. Click on '....'
3. See error

**Expected behavior**
What should happen

**Environment:**
- OS: [e.g. Ubuntu 20.04]
- Go version: [e.g. 1.24]
- Browser: [e.g. Chrome 120]
```

### Feature Requests

```markdown
**Feature Description**
Clear description of the feature

**Use Case**
Why is this feature needed?

**Proposed Solution**
How should this work?

**Additional Context**
Any other relevant information
```

## 🎯 Areas for Contribution

### High Priority
- 🐛 **Bug fixes** - Always welcome
- 📚 **Documentation** - Improve guides and examples
- 🧪 **Tests** - Increase test coverage
- 🔒 **Security** - Security improvements
- ⚡ **Rate limiting** - Performance optimizations and edge cases

### Medium Priority
- ✨ **Features** - New functionality (discuss first)
- 🎨 **UI/UX** - Web interface improvements
- ⚡ **Performance** - Optimization improvements
- 🔧 **Monitoring** - Add metrics and observability

### Rate Limiting Specific Contributions
- 🚀 **Storage backends** - Add support for other databases (MongoDB, PostgreSQL)
- 📊 **Metrics** - Add Prometheus metrics for rate limiting
- 🔧 **Configuration** - Dynamic rate limit configuration
- 🌐 **Distributed** - Improve distributed rate limiting algorithms
- 🧪 **Benchmarks** - Performance benchmarking and optimization

### Ideas for New Contributors
- 📖 Improve documentation and examples
- 🧪 Add more comprehensive tests for rate limiting
- 🐛 Fix "good first issue" labeled bugs
- 🌐 Add internationalization support
- 📱 Mobile UI improvements
- 📊 Add rate limiting dashboard/monitoring

## 🤝 Getting Help

- 💬 **Discussions**: [GitHub Discussions](https://github.com/pandeptwidyaop/tempfile/discussions)
- 📧 **Email**: pandeptwidyaop@gmail.com
- 🐛 **Issues**: [GitHub Issues](https://github.com/pandeptwidyaop/tempfile/issues)

## 📜 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to TempFiles!** 🚀
