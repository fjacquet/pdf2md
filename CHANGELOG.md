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

## [0.1.0] - 2023-10-27

### Added

- Initial release of `pdf2md`.
- Basic text extraction using `ledongthuc/pdf`.
- Simple layout analysis.
