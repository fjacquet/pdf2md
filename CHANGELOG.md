# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `cmd/benchmark`: `-cmd` and `-output` flags score any converter, each document also gets an order-free word F1, and `benchmark_pdftotext.json` commits poppler `pdftotext` 26.10.0's per-document scores (average 0.8609) as the reference to beat.

## [0.2.0] - 2026-10-28

### Added

- `--password` flag and `PDF2MD_PASSWORD` for encrypted PDFs, and AES-256 (V=5/R=6) decryption. Owner passwords and R=5 are not supported.
- Golden-file tests for synthetic layouts (`cmd/pdf2md/testdata/golden`) and a READoc score baseline (`benchmark_baseline.json`) checked by `cmd/benchmark`.
- `TestCLIFlags` freezes the public CLI flag set; tests for `cmd/pdf2md` and `internal/modelcache`.
- **High Fidelity Extraction**: Implemented a custom PDF parser (`internal/pdf`) for accurate text, font, and image extraction without CGO.
- **Advanced Layout Analysis**: Implemented Recursive XY-Cut algorithm for robust multi-column detection.
- **Image Support**: Extraction of raster images and vector graphics to Markdown.
- **Benchmarking Tool**: Custom Go-based benchmark comparing output against `READoc` dataset.
- **Fuzz Testing**: Added fuzz tests for `Analyzer` stability.
- **Security**: Added `govulncheck` and SBOM generation.

### Changed

- README: ML features documented as opt-in (they need `libonnxruntime`), flag and environment variable tables, stability policy, current roadmap.
- **Architecture**: Refactored into `cmd/`, `internal/`, `pkg/` structure.
- **Text Merging**: Improved paragraph merging with style awareness (bold/italic).
- **Header Detection**: Refined heuristics using font size metadata.
- **Table Extraction**: Improved table detection and TOC cleaning.

### Fixed

- Words were glued on Chrome-generated PDFs: TJ offsets and same-line gaps wider than 0.125 em (half a typical space) now become a word space; before, only TJ offsets below -200 did, and same-line blocks closer than 0.3 em were glued. READoc average 0.773 -> 0.807.
- Encrypted PDFs whose `/O`, `/U` or `/ID` were written as hex strings (e.g. by qpdf) could never be opened: hex strings are now decoded before use, and hex strings inside encrypted objects are decrypted from their decoded bytes.
- Fixed list item misclassification.
- Fixed "missing bullets" issue by leveraging custom parser's text operators.
- Fixed compilation errors in `internal/pdf` during refactoring.

### Security

- Move Go toolchain to 1.27.2 to fix GO-2026-6603, GO-2026-6605, GO-2026-6607, GO-2026-6608, GO-2026-6610, GO-2026-6611, GO-2026-6613 and GO-2026-6617 (net/http, crypto/tls, mime/multipart).
- Update dependencies: go-pdfium v1.21.1 -> v1.21.2, golang.org/x/sys v0.48.0 -> v0.49.0.
- Pin `golangci-lint` v2.13.2 -> v2.14.0 (v2.13.2 cannot type-check against the Go 1.27.2 standard library).

## [0.1.8] - 2026-10-27

### Fixed

- ONNX layout detection and OCR were silently disabled on Windows: the PDF renderer now opens the
  file on the host and hands go-pdfium a reader, because Wazero rejects Windows volume paths such
  as `D:\` (#9).

### Changed

- CI runs the `internal/render` tests on `windows-latest`; the render tests now generate their own
  PDF instead of skipping when the gitignored `testdata/test.pdf` is absent.

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
