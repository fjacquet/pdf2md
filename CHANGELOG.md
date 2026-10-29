# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- Bold headings promoted by ONNX layout detection or from a numbered table row rank with the other bold headings of the same size instead of one level deeper.
- Numbered headings set in one style keep their hierarchy: "1 Introduction", "1.1 Scope", "1.1.1 Terms" in the same font come out `#`, `##`, `###` instead of all `#`. Unnumbered headings of that style sit with its shallowest number.
- `examples/basic_usage.go` and the library usage notes call `RankHeaderLevels`, so library output gets the same heading levels as the CLI.

## [0.5.0] - 2026-10-28

### Fixed

- Heading levels follow the document's own hierarchy: heading styles (size rounded to 0.5 pt, bold) are ranked across the whole document, the largest is `#`, the next `##`, up to `######`. They used to come from fixed size ratios, so sections sat at `#` next to the title or whole documents started at `####`. READoc heading F1 (level and text): GitHub 0.18 -> 0.89, arXiv 0.00 -> 0.29; average 0.8367 -> 0.8383.
- Numbered lines in a row ("1. Fork it", "2. Create your feature branch") are list items, not `#` headings; a numbered line standing alone, or larger or bold, is still a heading. Consecutive bold lines (an abstract) stay a paragraph instead of one heading per line. A heading may start with any character but a lower-case letter, so ".NET MAUI Learning Resources" is no longer left in bold. Ordered list items no longer get a leading space.

### Added

- `cmd/benchmark` reports a heading F1 per document (level and text, then text only) and per-subset averages of heading and word F1.

## [0.4.1] - 2026-10-28

### Fixed

- ToUnicode CMaps are read completely (ISO 32000 §9.10.3): multi-character `bfrange` destinations increment their last UTF-16 unit, array destinations, codes of up to 4 bytes, and the codespace, which now splits composite codes when the font embeds no encoding CMap. An unreadable ToUnicode stream no longer drops the whole font.
- Simple fonts without `/Encoding` use the built-in encoding of their embedded Type 1 program, and StandardEncoding only when the font is not symbolic. Type 3 fonts without `/Encoding` keep their bytes instead of a glyph-name guess made in random map order. READoc 0.8357 -> 0.8365, arXiv 0.7797 (pdftotext 0.7785).

## [0.4.0] - 2026-10-28

### Fixed

- Form XObjects: nesting is bounded (32 levels), so a form that draws itself no longer crashes with a stack overflow; the form `/Matrix` is applied before the CTM (§8.10.1), so form text lands where viewers draw it; an indirect `/XObject` dictionary is resolved, so its images and forms are no longer dropped; an unknown or broken XObject is skipped with a warning instead of losing the whole page.
- **Security:** a FlateDecode or LZWDecode stream that decodes to more than 256 MiB is refused (`ErrStreamTooLarge`). A 1 MB PDF holding 1 GiB of zeros used to take 2.16 GB of memory; it now stops at the cap and skips that stream.
- Page attributes inherited from the page tree (`/Resources`, `/MediaBox`, `/CropBox`, `/Rotate`, ISO 32000 §7.7.3.4) are honoured. A page whose `/Resources` sat on its `/Pages` parent was lost entirely ("no resources dictionary"), and an inherited `/MediaBox` gave a 0×0 page, which silently disabled `--exclude-top` and fed ONNX/OCR a zero page size.
- A `/Kids` cycle in the page tree, or an object stream claiming to sit inside an object stream, now returns an error instead of crashing with a stack overflow.
- Library warnings (unparsable font, missing XObject, image encode failure, …) go to stderr. Without an output file the Markdown is printed on stdout, and these warnings used to land inside it.
- A page whose only content is an unsaved image or a too-small graphic no longer panics.
- **Security:** saved images and SVG graphics are named `<prefix>_p<page>_<n>.<ext>`. The PDF's XObject name no longer reaches the file system: a name such as `/..#2F..#2Fevil` could write outside the `images/` directory. Pages no longer overwrite each other's `Im1` / `Path1` files, and an unsaved or trivial graphic no longer turns into a broken Markdown image link.
- Word spaces are decided against each font's own space width (half of it), for TJ offsets and same-line gaps alike, instead of a fixed 0.125 em. Fonts without a space glyph (pdfTeX, FirstChar above 32) fall back to 250/1000 em instead of taking `/MissingWidth`, and WinAnsi's second `space` (code 160) no longer wins at random. `TextBlock.SpaceWidth` carries the width in user space.
- Glyph widths follow the spec order: `/W` and `/DW` for composite fonts; `/Widths`, `/MissingWidth`, then the standard 14 metrics (generated from Adobe's Core14 AFM files) for simple fonts; Type 3 widths scaled by `FontMatrix`. Composite fonts split codes by their encoding CMap's codespace. Standard fonts without `/Widths` (e.g. Times-Roman in TeX output) no longer overlap the next word.
- Text positions follow ISO 32000 §9.4.4: every glyph advances by its width plus `Tc`, `Tw` (on single-byte code 32) and `Tz`, strings inside `TJ` advance the pen, and `TD`, `Ts` (rise) and `Tr` (render mode, kept on `TextBlock.RenderMode`) are honoured. Line extents now match `pdftotext -bbox-layout` within 0.01 pt. READoc average 0.834 -> 0.836.
- Form XObjects that declare their own `/Resources /Font` now decode their text with those fonts instead of the page's.

## [0.3.0] - 2026-10-28

### Changed

- Reading order follows the content stream, the order the producer wrote the page in. XY-cut is used only on pages where the stream jumps back up within a column on more than 10% of line transitions. READoc average 0.807 -> 0.834 (arXiv 0.724 -> 0.773); two-column text no longer comes out as a table.

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
