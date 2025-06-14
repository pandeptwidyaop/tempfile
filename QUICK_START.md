# 🚀 Quick Start Guide

Get TempFiles up and running in minutes!

## 🎯 For Users

### Option 1: Download Pre-built Binary (Recommended)

```bash
# Linux (AMD64)
wget https://github.com/pandeptwidyaop/tempfile/releases/latest/download/tempfile-linux-amd64.tar.gz
tar -xzf tempfile-linux-amd64.tar.gz
./tempfile-linux-amd64

# macOS (Apple Silicon)
wget https://github.com/pandeptwidyaop/tempfile/releases/latest/download/tempfile-darwin-arm64.tar.gz
tar -xzf tempfile-darwin-arm64.tar.gz
./tempfile-darwin-arm64

# Windows (Download from releases page)
# https://github.com/pandeptwidyaop/tempfile/releases/latest
```

### Option 2: Docker

```bash
# Quick run
docker run -p 3000:3000 ghcr.io/pandeptwidyaop/tempfile:latest

# With persistent storage
docker run -p 3000:3000 -v $(pwd)/uploads:/app/uploads ghcr.io/pandeptwidyaop/tempfile:latest
```

### Option 3: Build from Source

```bash
git clone https://github.com/pandeptwidyaop/tempfile.git
cd tempfile
make setup
make dev
```

## 🛠️ For Developers

### First Time Setup

```bash
# Clone repository
git clone https://github.com/pandeptwidyaop/tempfile.git
cd tempfile

# Setup development environment (installs tools, creates .env, etc.)
make setup

# Start development server with live reload
make dev
```

### Daily Development Workflow

```bash
# 1. Create feature branch
git checkout -b feature/awesome-feature

# 2. Make changes and test
make test

# 3. Run quality checks
make check

# 4. Commit with conventional format
make commit
# This opens interactive commit helper:
# - Select type: feat, fix, docs, etc.
# - Enter scope (optional): api, ui, etc.
# - Enter description
# - Mark breaking changes if any

# 5. Push and create PR
git push origin feature/awesome-feature
```

### Available Commands

| Command | Description |
|---------|-------------|
| `make setup` | First-time development setup |
| `make dev` | Start development server |
| `make test` | Run tests |
| `make check` | Run all quality checks |
| `make commit` | Interactive commit helper |
| `make build` | Build binary |
| `make help` | Show all commands |

## 🏷️ Commit Message Examples

### ✅ Good Examples
```bash
# New feature
feat: add drag and drop upload
feat(api): add batch upload endpoint

# Bug fix
fix: resolve memory leak in cleanup
fix(ui): fix mobile responsive layout

# Documentation
docs: add API usage examples
docs(readme): update installation guide

# Breaking change
feat!: change upload API response format

BREAKING CHANGE: All responses now include metadata
```

### ❌ Avoid These
```bash
# Too vague
update stuff
fix bug
changes

# Wrong format
Fix: upload issue    # Wrong case
added feature        # Missing type
FIX: Bug            # All caps
```

## 🔄 Release Process

### For Maintainers

```bash
# 1. Ensure you're on main branch
git checkout main
git pull origin main

# 2. Run release preparation
make release-prep
# This runs all quality checks and ensures everything is ready

# 3. Push to main (triggers automatic release)
git push origin main
```

### Automatic Process
1. **Push to main** triggers GitHub Actions
2. **Analyzes commits** using conventional commit format
3. **Determines version** bump (MAJOR.MINOR.PATCH)
4. **Builds** multi-platform binaries
5. **Creates** GitHub release with changelog
6. **Uploads** artifacts (Linux, macOS, Windows, Docker)

## 🧪 Testing

```bash
# Run all tests
make test

# Test with coverage
make test-coverage

# Security scan
make security

# All quality checks at once
make check
```

## 🐳 Docker Development

```bash
# Build Docker image
make docker-build

# Run with Docker Compose
make docker-up

# Stop services
make docker-down

# Test health endpoint
make docker-health
```

## 📱 Usage Examples

### Web Interface
1. Open http://localhost:3000
2. Drag & drop file or click to select
3. File automatically expires in 1 hour
4. Share the download link

### API Usage
```bash
# Upload file
curl -X POST -F "file=@document.pdf" http://localhost:3000/

# Download file
curl -O http://localhost:3000/1718270400.pdf

# Check health
curl http://localhost:3000/health
```

## 🔧 Configuration

Create `.env` file (copied from `.env.example` during setup):

```bash
# Basic configuration
PORT=3000
PUBLIC_URL=http://localhost:3000
MAX_FILE_SIZE=104857600  # 100MB
FILE_EXPIRY_HOURS=1

# Advanced options
ENABLE_CORS=true
ENABLE_WEB_UI=true
DEBUG=false
```

## ❗ Troubleshooting

### Common Issues

**Port already in use:**
```bash
# Change port in .env
PORT=8080
```

**Permission denied:**
```bash
# Make sure script is executable
chmod +x dev.sh
```

**Go not found:**
```bash
# Install Go 1.24+ from https://golang.org/
```

**Tests failing:**
```bash
# Clean and reinstall dependencies
make clean
make deps
make test
```

### Getting Help

- 📖 **Full Documentation**: [README.md](README.MD)
- 🏷️ **Versioning Guide**: [SEMANTIC_VERSIONING.md](SEMANTIC_VERSIONING.md)
- 🤝 **Contributing**: [CONTRIBUTING.md](CONTRIBUTING.md)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/pandeptwidyaop/tempfile/discussions)
- 🐛 **Issues**: [GitHub Issues](https://github.com/pandeptwidyaop/tempfile/issues)
- 📧 **Email**: pandeptwidyaop@gmail.com

---

**🎉 That's it! You're ready to use TempFiles!**

For more detailed information, check out the [full README](README.MD) or [contributing guide](CONTRIBUTING.md).
