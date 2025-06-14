# 🏷️ Conventional Commit Examples

## 📝 Format Dasar
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## 🎯 Contoh Commit Messages

### ✅ PATCH Version (0.0.x)

#### Bug Fixes
```bash
fix: resolve memory leak in file cleanup service
fix(api): handle large file upload timeout
fix(ui): fix responsive layout on mobile devices
fix(docker): correct healthcheck endpoint path
```

#### Documentation
```bash
docs: add API usage examples to README
docs(contributing): update development setup guide
docs: fix typos in installation instructions
```

#### Code Style & Refactoring
```bash
style: format code according to Go standards
refactor: extract upload logic to separate service
refactor(handlers): simplify error handling logic
```

#### Tests & Chores
```bash
test: add unit tests for file cleanup service
test(api): add integration tests for upload endpoint
chore: update Go dependencies to latest versions
chore(ci): improve GitHub Actions caching
```

### ✅ MINOR Version (0.x.0)

#### New Features
```bash
feat: add drag and drop file upload interface
feat(api): implement batch file upload endpoint
feat(ui): add dark theme toggle
feat: add file preview functionality
feat(config): add Redis backend support
```

#### Feature with Scope
```bash
feat(auth): add optional user authentication
feat(storage): implement S3 storage backend
feat(monitoring): add metrics dashboard
feat(api): add file compression option
```

### ✅ MAJOR Version (x.0.0)

#### Breaking Changes (Method 1: exclamation mark)
```bash
feat!: change upload API response format
fix!: rename configuration environment variables
refactor!: restructure internal package layout
```

#### Breaking Changes (Method 2: footer)
```bash
feat: implement new authentication system

BREAKING CHANGE: All API endpoints now require authentication header
```

```bash
refactor: change configuration file format

BREAKING CHANGE: Configuration now uses YAML instead of JSON format.
Migrate your config files from .json to .yaml
```

## 🔧 Contoh dengan Body dan Footer

### Feature dengan Detail
```bash
feat(upload): add resumable file upload support

Implement resumable uploads using tus protocol to handle
large files and unstable network connections.

- Add tus server implementation
- Update web UI to support resumable uploads
- Add progress tracking and resume capability

Closes #123
```

### Bug Fix dengan Detail
```bash
fix(cleanup): resolve race condition in file deletion

The cleanup service was experiencing race conditions when
multiple goroutines tried to delete the same file simultaneously.

- Add mutex lock for file deletion operations
- Implement proper error handling for concurrent access
- Add unit tests for concurrent cleanup scenarios

Fixes #456
Co-authored-by: John Doe <john@example.com>
```

### Breaking Change dengan Migration Guide
```bash
feat!: implement new configuration system

Replace environment variables with YAML configuration file
for better organization and validation.

BREAKING CHANGE: Environment-based configuration is no longer supported.

Migration guide:
1. Create config.yaml file
2. Move environment variables to YAML format
3. Update deployment scripts

Before:
PORT=3000
MAX_FILE_SIZE=104857600

After:
server:
  port: 3000
  max_file_size: 104857600

Closes #789
```

## 🚨 Contoh yang SALAH

### ❌ Format Tidak Benar
```bash
# Huruf kapital salah
Fix: upload bug
FEAT: new feature
Update readme

# Missing type
added new upload feature
resolve memory issue
update dependencies

# Terlalu umum
fix bug
update
changes
improvements
```

### ❌ Deskripsi Buruk
```bash
# Terlalu singkat
fix: bug
feat: stuff
docs: update

# Terlalu panjang (lebih dari 50 karakter untuk subject)
feat: add a very long description that explains every single detail about the feature implementation
```

## 🎯 Tips untuk Commit Messages yang Baik

### 1. **Gunakan Present Tense**
```bash
✅ fix: resolve upload issue
❌ fix: resolved upload issue
❌ fix: resolves upload issue
```

### 2. **Mulai dengan Huruf Kecil**
```bash
✅ feat: add dark theme
❌ feat: Add dark theme
```

### 3. **Maksimal 50 Karakter untuk Subject**
```bash
✅ feat: add file compression support
❌ feat: add comprehensive file compression support with multiple algorithms
```

### 4. **Gunakan Scope untuk Clarifikasi**
```bash
✅ feat(api): add batch upload endpoint
✅ fix(ui): resolve mobile layout issue
✅ docs(readme): update installation guide
```

### 5. **Jelaskan "What" dan "Why", bukan "How"**
```bash
✅ fix: resolve memory leak in cleanup service
❌ fix: add defer statement and close file handle
```

## 🔄 Workflow dengan Conventional Commits

### 1. **Development**
```bash
git checkout -b feature/file-compression
# ... make changes ...
git add .
git commit -m "feat: add file compression support"
```

### 2. **Bug Fix**
```bash
git checkout -b fix/memory-leak
# ... fix the issue ...
git add .
git commit -m "fix: resolve memory leak in cleanup service"
```

### 3. **Breaking Change**
```bash
git checkout -b breaking/new-api
# ... implement breaking changes ...
git add .
git commit -m "feat!: change API response format

BREAKING CHANGE: All API responses now include metadata object"
```

### 4. **Multiple Commits**
```bash
# Feature development dengan multiple commits
git commit -m "feat(ui): add file upload progress bar"
git commit -m "test: add tests for upload progress tracking"
git commit -m "docs: update API documentation for progress endpoint"
```

## 🛠️ Using Interactive Commit Helper

Gunakan script helper yang sudah disediakan:

```bash
make commit
```

Atau langsung:
```bash
./dev.sh commit
```

Script ini akan:
1. **Check staged changes** - Memastikan ada perubahan yang di-stage
2. **Interactive menu** - Pilih tipe commit dari menu
3. **Input scope** - Masukkan scope opsional  
4. **Input description** - Masukkan deskripsi
5. **Breaking change check** - Tanyakan apakah breaking change
6. **Preview & confirm** - Tampilkan preview dan konfirmasi

## 📊 Version Bumping Examples

### Scenario 1: Patch Release
```bash
# Current version: 1.2.3
git commit -m "fix: resolve upload timeout issue"
# Next version: 1.2.4
```

### Scenario 2: Minor Release  
```bash
# Current version: 1.2.4
git commit -m "feat: add file preview functionality"
# Next version: 1.3.0
```

### Scenario 3: Major Release
```bash
# Current version: 1.3.0
git commit -m "feat!: change API authentication method"
# Next version: 2.0.0
```

### Scenario 4: Multiple Commits
```bash
# Current version: 2.0.0
git commit -m "fix: resolve UI layout issue"        # patch
git commit -m "feat: add dark theme support"        # minor
git commit -m "docs: update README"                 # patch
# Next version: 2.1.0 (highest level wins)
```

## 🎉 Kesimpulan

Conventional commits memberikan:
- ✅ **Consistent** commit history
- ✅ **Automatic** semantic versioning  
- ✅ **Generated** changelogs
- ✅ **Clear** communication dalam tim
- ✅ **Automated** release process

Mulai gunakan format ini untuk semua commit dan nikmati release automation yang sudah disiapkan! 🚀
