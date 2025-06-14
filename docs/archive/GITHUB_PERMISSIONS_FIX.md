# 🔧 Mengatasi GitHub Actions Permission Issues

## 🚨 Problem: EGITNOPERMISSION Error

Error yang Anda alami terjadi karena `GITHUB_TOKEN` default tidak memiliki permission untuk push tag dan commit ke repository.

## ✅ Solutions

### Opsi 1: Enable GitHub Actions Permissions (Recommended)

1. **Buka repository settings di GitHub**
2. **Go to Settings → Actions → General**
3. **Scroll ke "Workflow permissions"**
4. **Pilih "Read and write permissions"**
5. **Check "Allow GitHub Actions to create and approve pull requests"**
6. **Save**

### Opsi 2: Create Personal Access Token (Alternative)

Jika opsi 1 tidak berhasil, buat PAT:

1. **Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)**
2. **Generate new token (classic)**
3. **Name**: `TempFile Release Token`
4. **Expiration**: Choose appropriate time
5. **Scopes**:
   - ✅ `repo` (Full control)
   - ✅ `write:packages` (jika pakai Docker registry)
6. **Generate token**
7. **Copy token value**

8. **Add to repository secrets**:
   - Go to repository **Settings → Secrets and variables → Actions**
   - Click **New repository secret**
   - Name: `PAT_TOKEN`
   - Value: paste your token
   - **Add secret**

## 🔄 Updated Workflow

Workflow sudah diupdate untuk mengatasi masalah ini:

### ✅ Yang Sudah Diperbaiki:

1. **Added permissions block**:
   ```yaml
   permissions:
     contents: write
     issues: write
     pull-requests: write
     actions: read
   ```

2. **Fallback token strategy**:
   ```yaml
   token: ${{ secrets.PAT_TOKEN || secrets.GITHUB_TOKEN }}
   ```

3. **Proper Git configuration**:
   ```yaml
   env:
     GIT_AUTHOR_NAME: github-actions[bot]
     GIT_AUTHOR_EMAIL: github-actions[bot]@users.noreply.github.com
   ```

4. **Simplified workflow** - mengurangi complexity yang bisa menyebabkan permission issues

## 🧪 Testing the Fix

### Test 1: Simple Release
```bash
# Create a simple commit to test
git add .
git commit -m "fix: resolve GitHub Actions permissions issue

- Add proper permissions to workflow
- Add PAT token fallback
- Simplify release process"

git push origin main
```

### Test 2: Monitor Workflow
1. Go to **Actions** tab
2. Watch the **🚀 Release & Semantic Versioning** workflow
3. Check for permission errors

### Expected Results:
- ✅ Workflow completes successfully
- ✅ New release is created
- ✅ Artifacts are uploaded
- ✅ Changelog is updated

## 🔍 Debugging Steps

### If Still Getting Permission Error:

1. **Check repository settings**:
   - Settings → Actions → General → Workflow permissions
   - Ensure "Read and write permissions" is selected

2. **Verify branch protection rules**:
   - Settings → Branches
   - Check if `main` branch has restrictions

3. **Check token scopes** (if using PAT):
   - Token must have `repo` scope
   - Token must not be expired

4. **Repository ownership**:
   - Ensure you have admin access to repository
   - Personal repositories vs Organization repositories have different permission models

## 🛠️ Alternative: Manual Release Process

Jika automation masih bermasalah, Anda bisa menggunakan manual process sementara:

```bash
# 1. Determine next version manually
# Check commits since last release
git log --oneline $(git describe --tags --abbrev=0)..HEAD

# 2. Create tag manually
git tag v1.0.0
git push origin v1.0.0

# 3. Build release manually
make build-all

# 4. Create GitHub release manually
# Go to GitHub → Releases → Create new release
# Upload binaries from dist/ folder
```

## 📋 Troubleshooting Checklist

- [ ] Repository permissions set to "Read and write"
- [ ] PAT token created with `repo` scope (if needed)
- [ ] PAT token added as `PAT_TOKEN` secret
- [ ] Branch protection rules allow Actions to push
- [ ] You have admin access to repository
- [ ] Workflow file syntax is correct
- [ ] No typos in secret names

## 🎯 Quick Test Command

Test if permissions are working:

```bash
# Add this simple change and commit
echo "# Test Release" >> TEST_RELEASE.md
git add TEST_RELEASE.md
git commit -m "test: verify release automation permissions"
git push origin main
```

This should trigger a **PATCH** release (0.0.x) if everything is configured correctly.

## 📞 Need Help?

Jika masalah masih berlanjut:

1. **Check GitHub Actions logs** untuk error message spesifik
2. **Share error logs** jika masih ada masalah
3. **Verify** semua steps diatas sudah diikuti

Once permissions are fixed, automatic releases akan berjalan lancar! 🚀
