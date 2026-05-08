# 📋 Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [2.3.0](https://github.com/pandeptwidyaop/tempfile/compare/v2.2.3...v2.3.0) (2026-05-08)

### ✨ Features

* add paste view template ([5f38fca](https://github.com/pandeptwidyaop/tempfile/commit/5f38fca1784899dbd50a067592e0ed75368c2788))
* add paste.js for tab switching, copy, line numbers ([3bf0af0](https://github.com/pandeptwidyaop/tempfile/commit/3bf0af0651f6397c9869ee9db89d3e6571d1df9e))
* add PasteHandler with create, view, and raw endpoints ([6dd70c1](https://github.com/pandeptwidyaop/tempfile/commit/6dd70c19a125db9d9d6e003ffd8d8355ad9fd156))
* add PasteResponse model and paste fields to WebPageData ([68562bf](https://github.com/pandeptwidyaop/tempfile/commit/68562bf754a2ba2c2e2cf4392ad346736bd401a8))
* add PasteService for text sharing ([90334ff](https://github.com/pandeptwidyaop/tempfile/commit/90334ff376c6a48b780a2ff31530fb57e49a4f05))
* add RenderPastePage helper to template service ([9214dcb](https://github.com/pandeptwidyaop/tempfile/commit/9214dcb59b480342f7c2035b7de08749029407c2))
* add Share Text tab to upload page ([170dce5](https://github.com/pandeptwidyaop/tempfile/commit/170dce55663dd458b64a4ee21ca9c79679cc93b0))
* add styles for paste tab and paste view ([ca160fc](https://github.com/pandeptwidyaop/tempfile/commit/ca160fcf7de70c4918421668c3f71ad001ecc694))
* cleanup also scans pastes/ subdirectory ([7c3e6d0](https://github.com/pandeptwidyaop/tempfile/commit/7c3e6d0b1bc2ce4dcf3e8a050ceda72dfea42f07))
* redesign paste view with share-link section, raw download, and clipboard fallback ([d52c5fe](https://github.com/pandeptwidyaop/tempfile/commit/d52c5fea3e0f46e565ce37eb8352728e16c6d49f))
* wire paste service, handler, and routes ([79bf1b0](https://github.com/pandeptwidyaop/tempfile/commit/79bf1b0d13dbe5d319f43edff2e94a5c62cc7fbc))

### 🐛 Bug Fixes

* trigger build ([86d03be](https://github.com/pandeptwidyaop/tempfile/commit/86d03beed8a21d10d9e8067be529c7f24abef5b4))

### 📚 Documentation

* document text-sharing endpoints in README ([a85a460](https://github.com/pandeptwidyaop/tempfile/commit/a85a460cbf09da47a2606d0adb4256a7271acd91))

## [2.2.3](https://github.com/pandeptwidyaop/tempfile/compare/v2.2.2...v2.2.3) (2026-03-10)

### 🐛 Bug Fixes

* dynamic client-side max file size validation ([80165ea](https://github.com/pandeptwidyaop/tempfile/commit/80165ea151915f330c051b1ff13b21b1e97ba037))

## [2.2.2](https://github.com/pandeptwidyaop/tempfile/compare/v2.2.1...v2.2.2) (2026-03-10)

### 🐛 Bug Fixes

* use dynamic max file size in upload template ([d8f7e13](https://github.com/pandeptwidyaop/tempfile/commit/d8f7e133d54c88702abf22a16ac6155205333df7))

## [2.2.1](https://github.com/pandeptwidyaop/tempfile/compare/v2.2.0...v2.2.1) (2025-07-25)

### 🐛 Bug Fixes

* parse timestamp from filename ([b0f13e0](https://github.com/pandeptwidyaop/tempfile/commit/b0f13e0cab13817c3f969be24e4ee08bf10c135f))
* update temporary filename ([46c81eb](https://github.com/pandeptwidyaop/tempfile/commit/46c81eb41621c2639712c85b9767d644f68325b8))

## [2.2.0](https://github.com/pandeptwidyaop/tempfile/compare/v2.1.0...v2.2.0) (2025-06-14)

### ✨ Features

* Add build time generation and improve Docker tag strategy in workflows ([8c44baf](https://github.com/pandeptwidyaop/tempfile/commit/8c44baf389078066b931bf79b6c224052bdce9e2))
* Enhance Docker workflows and management for GHCR integration ([e2ae8c2](https://github.com/pandeptwidyaop/tempfile/commit/e2ae8c24b765965d7ce2d7d1883d99a13e3999b4))

## [2.1.0](https://github.com/pandeptwidyaop/tempfile/compare/v2.0.0...v2.1.0) (2025-06-14)

### ✨ Features

* implement hybrid loading for static and template files ([89f2bb7](https://github.com/pandeptwidyaop/tempfile/commit/89f2bb71d51d1e7fdda6b5931ad6f41dcb464bfe))

### 🐛 Bug Fixes

* fix running bin ([9678d32](https://github.com/pandeptwidyaop/tempfile/commit/9678d328063d3da7513debf77d4d3331635a5735))

## [2.0.0](https://github.com/pandeptwidyaop/tempfile/compare/v1.1.0...v2.0.0) (2025-06-14)

### ⚠ BREAKING CHANGES

* Add comprehensive rate limiting system with Redis backend,
IP whitelisting, distributed architecture, and advanced security hardening
* Add comprehensive rate limiting system with Redis backend,
IP whitelisting, distributed architecture, and advanced security hardening
* Add comprehensive rate limiting system with Redis backend,
IP whitelisting, distributed architecture, and advanced monitoring support
* Add comprehensive rate limiting system with Redis backend,
IP whitelisting, and distributed architecture support
* Add comprehensive rate limiting system with Redis backend
* Add comprehensive rate limiting system
* Rate limiter interface extended with CheckLimitsForEndpoint method

### ✨ Features

* implement advanced rate limiting with Redis backend and IP whitelisting ([cfd94ae](https://github.com/pandeptwidyaop/tempfile/commit/cfd94aed62243a21e3624bd19816b86395fa50a9))
* implement Redis-based rate limiting and update deployment configuration ([a9f855e](https://github.com/pandeptwidyaop/tempfile/commit/a9f855ec4c3b42e66ed144249982598073efd6ca))
* rate limiting ([951d561](https://github.com/pandeptwidyaop/tempfile/commit/951d5619595d91531c239638abef6517c445c8ef))

### 🐛 Bug Fixes

* go sec ([5c1cf72](https://github.com/pandeptwidyaop/tempfile/commit/5c1cf72e84832a57cd608917f5cf7836140bb36b))
* resolve all CI/CD issues and complete enterprise rate limiting ([96200a5](https://github.com/pandeptwidyaop/tempfile/commit/96200a5c8e7662d0b453fd3cfa8908fef0bebf96))
* resolve all CI/CD issues and complete enterprise rate limiting system ([4dff7ad](https://github.com/pandeptwidyaop/tempfile/commit/4dff7ada2ba6920da1892320fe59c32be1cf611e))
* resolve all CI/CD issues and complete rate limiting system ([473d49a](https://github.com/pandeptwidyaop/tempfile/commit/473d49a0b34908bc6ee465c1d7146b0d1d8e9a11))
* resolve all CI/CD issues and finalize rate limiting system ([6365bf2](https://github.com/pandeptwidyaop/tempfile/commit/6365bf2cea665a012394489d02a199b66627a9ab))
* resolve all CI/CD issues, security vulnerabilities, and complete enterprise rate limiting ([97a69c0](https://github.com/pandeptwidyaop/tempfile/commit/97a69c0c5dea2f73d98cc6324581d8efb345691c))
* resolve all CI/CD issues, security vulnerabilities, and complete enterprise rate limiting ([310a732](https://github.com/pandeptwidyaop/tempfile/commit/310a7327dc900d67a4401a3705c5735d8ca3eda7))

### 🧹 Chores

* update gitignore for compiled binary ([40d26ed](https://github.com/pandeptwidyaop/tempfile/commit/40d26ed961e7f747b577442ae79f1c1ddf0da0ca))

## [1.1.0](https://github.com/pandeptwidyaop/tempfile/compare/v1.0.0...v1.1.0) (2025-06-14)

### ✨ Features

* add comprehensive setup documentation for semantic versioning and automated release ([89f088c](https://github.com/pandeptwidyaop/tempfile/commit/89f088c3987837bb0c270d878090b6ed232cb840))
* restore test release workflow configuration ([821c75b](https://github.com/pandeptwidyaop/tempfile/commit/821c75bb93d0c675e007e8cd3d3ae12bbcd16e70))

## 1.0.0 (2025-06-14)

### ✨ Features

* add .env support and improved configuration ([473a410](https://github.com/pandeptwidyaop/tempfile/commit/473a410bcdcb1d4b983541608cba38e8b1885b04))
* add semantic versioning and auto-release system ([749bc6e](https://github.com/pandeptwidyaop/tempfile/commit/749bc6ee7b7784096b48b7ca4c0bb971ecdf5b74))
* complete project restructure to modern Go architecture ([07c19ec](https://github.com/pandeptwidyaop/tempfile/commit/07c19eca5de9360ba14653ecad3fd1d58192e9c8))
* enhance GitHub Actions workflows with permission fixes and test release process ([ea764e5](https://github.com/pandeptwidyaop/tempfile/commit/ea764e59064211e95628e15285cffac056f54656))
* enhance theme toggle functionality and improve security headers ([3b78b97](https://github.com/pandeptwidyaop/tempfile/commit/3b78b974913dfb372847f69c41c2926e349cf159))
* implement hybrid loading for static assets and templates with embedded support ([eca05f9](https://github.com/pandeptwidyaop/tempfile/commit/eca05f9827061b5d357af0fc559a0b99b9b0e590))

## [Unreleased]

### Added
- 🏷️ Semantic versioning with automated releases
- 🤖 GitHub Actions workflows for CI/CD
- 📦 Multi-platform binary builds (Linux, macOS, Windows)
- 🐳 Docker image builds in releases
- 📋 Comprehensive release documentation
- ✅ Automated testing and quality checks

### Changed
- 🔄 Updated development workflow to use conventional commits
- 📊 Enhanced build process with semantic versioning
- 🏗️ Improved project structure documentation

---

*This changelog will be automatically updated with each release based on conventional commit messages.*
