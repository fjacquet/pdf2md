# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **High Fidelity Extraction**: Implemented a custom PDF parser (`internal/pdf`) for accurate text, font, and image extraction without CGO.
- **Advanced Layout Analysis**: Implemented Recursive XY-Cut algorithm for robust multi-column detection.
- **Image Support**: Extraction of raster images and vector graphics to Markdown.
- **Benchmarking Tool**: Custom Go-based benchmark comparing output against `READoc` dataset.
- **Fuzz Testing**: Added fuzz tests for `Analyzer` stability.
- **Security**: Added `govulncheck` and SBOM generation.

### Changed

- **Architecture**: Refactored into `cmd/`, `internal/`, `pkg/` structure.
- **Text Merging**: Improved paragraph merging with style awareness (bold/italic).
- **Header Detection**: Refined heuristics using font size metadata.
- **Table Extraction**: Improved table detection and TOC cleaning.

### Fixed

- Fixed list item misclassification.
- Fixed "missing bullets" issue by leveraging custom parser's text operators.
- Fixed compilation errors in `internal/pdf` during refactoring.

### Security

- Move Go toolchain to 1.27.2 to fix GO-2026-6603, GO-2026-6605, GO-2026-6607, GO-2026-6608, GO-2026-6610, GO-2026-6611, GO-2026-6613 and GO-2026-6617 (net/http, crypto/tls, mime/multipart).
- Update dependencies: go-pdfium v1.21.1 -> v1.21.2, golang.org/x/sys v0.48.0 -> v0.49.0.
- Pin `golangci-lint` v2.13.2 -> v2.14.0 (v2.13.2 cannot type-check against the Go 1.27.2 standard library).

## [0.1.7] - 2026-10-03

### Changed

- Go toolchain moved to 1.27.1 (`go.mod`, release workflow `goreleaser-cross` image).
- `golangci-lint` pin bumped v2.12.2 -> v2.13.2, `goreleaser` v2.16.0 -> v2.18.0 and
  `govulncheck` v1.1.4 -> v1.8.0 (needed for Go 1.27).

## [0.1.6] - 2026-10-02

### Changed

- Dependency refresh: `github.com/alecthomas/kong` 1.13.0 -> 1.16.1,
  `github.com/klippa-app/go-pdfium` 1.17.2 -> 1.21.1,
  `github.com/yalue/onnxruntime_go` 1.22.0 -> 1.36.0, `golang.org/x/image` 0.41.0 -> 0.46.0.
- Security workflow added (`go-security` via `fjacquet/ci`).

### Fixed

- Rendering now uses go-pdfium `RenderedImage` instead of the deprecated `Image`.

## [0.1.0] - 2023-10-27

### Added

- Initial release of `pdf2md`.
- Basic text extraction using `ledongthuc/pdf`.
- Simple layout analysis.
