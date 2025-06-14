# 🔧 Docker Tag Fix - GitHub Actions

## Issue yang Diperbaiki

❌ **Error sebelumnya:**
```
ERROR: invalid tag "ghcr.io/pandeptwidyaop/tempfile:-8c87a02": invalid reference format
```

## ✅ **Perbaikan yang Diterapkan**

### 1. **Fixed Docker Metadata Configuration**

#### Before (❌ Bermasalah):
```yaml
tags: |
  type=sha,prefix={{branch}}-  # Menghasilkan tag seperti "-8c87a02"
```

#### After (✅ Diperbaiki):
```yaml
tags: |
  type=ref,event=branch         # Branch name tags
  type=ref,event=pr            # PR tags
  type=sha,prefix=sha-,format=short  # SHA tags: "sha-8c87a02"
  type=raw,value=dev,enable=${{ github.ref == 'refs/heads/develop' }}
  type=raw,value=latest,enable=${{ github.ref == 'refs/heads/main' }}
```

### 2. **Improved Tag Strategy**

| Event Type | Old Tag | New Tag | Valid? |
|------------|---------|---------|--------|
| PR | `pr-3`, `-8c87a02` | `pr-3`, `sha-8c87a02` | ✅ |
| Branch Push | `develop`, `-8c87a02` | `develop`, `sha-8c87a02` | ✅ |
| Main Branch | `main`, `-8c87a02`, `latest` | `main`, `sha-8c87a02`, `latest` | ✅ |

### 3. **Enhanced Metadata Labels**

```yaml
labels: |
  org.opencontainers.image.title=Tempfile
  org.opencontainers.image.description=A secure file sharing service
  org.opencontainers.image.vendor=${{ github.repository_owner }}
  org.opencontainers.image.source=https://github.com/${{ github.repository }}
  org.opencontainers.image.revision=${{ github.sha }}
  org.opencontainers.image.created=${{ steps.build-time.outputs.time }}
  org.opencontainers.image.url=https://github.com/${{ github.repository }}
  org.opencontainers.image.documentation=https://github.com/${{ github.repository }}/blob/main/README.md
```

### 4. **Fixed Build Args**

#### Before (❌):
```yaml
build-args: |
  BUILD_TIME=${{ github.event.head_commit.timestamp }}  # Bisa kosong
```

#### After (✅):
```yaml
- name: 🕑 Generate Build Time
  id: build-time
  run: echo "time=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" >> $GITHUB_OUTPUT

build-args: |
  BUILD_VERSION=${{ github.sha }}
  BUILD_TIME=${{ steps.build-time.outputs.time }}  # Selalu valid
```

## 🎯 **Expected Results**

### Build Workflow Tags:
- **PR**: `ghcr.io/pandeptwidyaop/tempfile:pr-3`, `ghcr.io/pandeptwidyaop/tempfile:sha-8c87a02`
- **Develop Branch**: `ghcr.io/pandeptwidyaop/tempfile:develop`, `ghcr.io/pandeptwidyaop/tempfile:dev`, `ghcr.io/pandeptwidyaop/tempfile:sha-8c87a02`
- **Feature Branch**: `ghcr.io/pandeptwidyaop/tempfile:feature-xyz`, `ghcr.io/pandeptwidyaop/tempfile:sha-8c87a02`

### Release Workflow Tags:
- **Main Branch**: `ghcr.io/pandeptwidyaop/tempfile:main`, `ghcr.io/pandeptwidyaop/tempfile:latest`, `ghcr.io/pandeptwidyaop/tempfile:sha-8c87a02`
- **After Semantic Release**: `ghcr.io/pandeptwidyaop/tempfile:v1.2.0`, `ghcr.io/pandeptwidyaop/tempfile:1.2.0`

## 🔍 **Validation**

Untuk memvalidasi fix ini:

1. **Push ke branch develop**:
   ```bash
   git push origin develop
   ```
   Expected tags: `develop`, `dev`, `sha-xxxxxxx`

2. **Create Pull Request**:
   Expected tags: `pr-N`, `sha-xxxxxxx`

3. **Push ke main branch**:
   Expected tags: `main`, `latest`, `sha-xxxxxxx`

4. **Check GitHub Packages**:
   - Go to repository → Packages
   - Verify all tags are valid and present

## 🚀 **Test Commands**

```bash
# Test pulling different tags
docker pull ghcr.io/pandeptwidyaop/tempfile:develop
docker pull ghcr.io/pandeptwidyaop/tempfile:sha-8c87a02
docker pull ghcr.io/pandeptwidyaop/tempfile:latest

# Test running
docker run --rm -p 3000:3000 ghcr.io/pandeptwidyaop/tempfile:develop

# Test health check
curl http://localhost:3000/health
```

## 📋 **Checklist**

- [x] Fixed invalid tag format (`-` prefix issue)
- [x] Added consistent build time generation
- [x] Enhanced metadata labels
- [x] Improved tag strategy for different events
- [x] Added proper conditional tags for branches
- [x] Fixed build args references
- [x] Added comprehensive OCI labels

**Status: ✅ READY FOR TESTING**

Push changes dan workflow akan generate tags yang valid!
