# pdf2md - PDF to Markdown Converter

[![CI](https://github.com/fjacquet/pdf2md/actions/workflows/ci.yml/badge.svg)](https://github.com/fjacquet/pdf2md/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/fjacquet/pdf2md?sort=semver)](https://github.com/fjacquet/pdf2md/releases/latest)
[![License](https://img.shields.io/github/license/fjacquet/pdf2md)](https://github.com/fjacquet/pdf2md/blob/HEAD/LICENSE)

A fast, single-binary PDF to Markdown converter written in Go. Designed for AI engineers, researchers, and knowledge workers who need clean, structured Markdown from PDF documents.

## How it Works

- **Single Binary** - Download and run; the core PDF to Markdown path has no native dependencies
- **Custom PDF Parser** - Pure Go implementation, no CGO or external libraries required for the core path
- **Optional ML** - ONNX layout detection and OCR for scanned pages are opt-in and need `libonnxruntime` on the host (see below)
- **Column Detection** - Handles multi-column layouts (academic papers, newspapers)
- **Smart Layout Analysis** - Preserves reading order and document structure
- **Header Detection** - Automatically identifies H1-H6 based on font size or numbering
- **Code Block Detection** - Identifies code blocks based on content patterns
- **List Detection** - Handles bulleted and numbered lists
- **Privacy-First** - Runs completely locally, no cloud API calls
- **Fast** - Built in Go for high performance

## Installation

### From Source

```bash
go install github.com/fjacquet/pdf2md/cmd/pdf2md@latest
```

### Build Locally

```bash
git clone https://github.com/fjacquet/pdf2md
cd pdf2md
go build -o pdf2md ./cmd/pdf2md
```

## Usage

Convert a PDF to Markdown (output to stdout):

```bash
./pdf2md input.pdf
```

Save to a file:

```bash
./pdf2md input.pdf output.md
```

### Options

| Flag | Purpose |
|------|---------|
| `-V`, `--version` | Print version and exit |
| `--debug` | Enable debug logging |
| `--exclude-top=N` / `--exclude-bottom=N` | Skip a header / footer region of N points |
| `--password=PW` | User password for an encrypted PDF (or set `PDF2MD_PASSWORD`, which keeps it out of the shell history); the blank password is always tried first |
| `--no-onnx` | Disable ONNX layout detection (rule-based only) |
| `--model-path=FILE` | ONNX layout model (auto-downloaded if omitted) |
| `--runtime-path=FILE` | ONNX Runtime library (auto-detected if omitted) |
| `--no-ocr` | Disable the OCR fallback on image-only pages |
| `--ocr-det-model-path=FILE` / `--ocr-rec-model-path=FILE` | OCR detection / recognition models |
| `--ocr-dict-path=FILE` | OCR character dictionary (only for non-Latin models) |

Environment variables: `PDF2MD_MODEL_PATH`, `PDF2MD_OCR_DET_PATH`, `PDF2MD_OCR_REC_PATH`,
`PDF2MD_OCR_DICT_PATH`, `PDF2MD_PASSWORD`, `ORT_LIB_PATH`. OCR models are not downloaded automatically; see
[docs/ocr-setup.md](docs/ocr-setup.md). Without `libonnxruntime` or the models, pdf2md falls back
to rule-based analysis and skips OCR without failing.

### Stability

From 1.0, the command-line flags and environment variables above, and the Markdown written to
stdout or the output file, follow semantic versioning: a rename or removal needs a new major
version. The Go packages live under `internal/` and carry no compatibility promise.

## Architecture

pdf2md uses a modular pipeline architecture:

1. **PDF Parser** (`internal/pdf`) - Custom pure Go PDF parser that handles:
   - Cross-reference table parsing
   - Object reading (dictionaries, streams, arrays)
   - Content stream interpretation
   - Font handling (CMap, Encodings)
2. **Extractor** (`internal/extractor`) - Abstraction layer for text extraction
3. **Analyzer** (`internal/layout`) - Detects columns, headers, and reading order
4. **Builder** (`internal/markdown`) - Generates clean Markdown output

### Current Implementation

- **Pure Go**: Replaced `go-fitz` with a custom PDF parser implementation.
- **Robust**: Handles various PDF versions and structures.
- **Extensible**: Designed to support more PDF features in the future.

## Roadmap

### Phase 1: Core Parsing (Completed)

- [x] Custom PDF Reader & Parser
- [x] Content Stream Interpreter
- [x] Font Handling (Type1, TrueType, CMap)
- [x] Text Extraction

### Phase 2: Layout Analysis (Completed)

- [x] Column detection
- [x] Header identification (H1-H6)
- [x] Reading order preservation
- [x] List detection
- [x] Code block detection

### Phase 3: Visual Elements (Completed)

- [x] Table extraction and formatting (Basic implementation)
- [x] Image extraction to assets folder
- [x] Link extraction
- [x] Vector graphics extraction

### Phase 4: Advanced

- [x] OCR support for scanned PDFs (ONNX, opt-in)
- [x] ONNX layout detection (DocLayout-YOLO, opt-in)
- [x] Configurable exclusion zones (headers/footers)
- [x] Custom formatting rules

### Planned for 1.0

- [x] AES-256 (V=5/R=6) encryption and non-empty user passwords (`--password`)
- [ ] Golden-file regression corpus

### After 1.0

- [ ] CCITT G4 and JBIG2 image decoders (currently `ErrFilterUnsupported`)
- [ ] Automatic download of OCR models

## Development

### Project Structure

```
pdf2md/
├── cmd/pdf2md/          # Main CLI application
├── internal/
│   ├── pdf/             # Custom PDF parser implementation
│   ├── extractor/       # Extraction interface and implementations
│   ├── layout/          # Layout analysis and structure detection
│   ├── markdown/        # Markdown generation
│   ├── types/           # Shared types (TextBlock, etc.)
│   └── testdata/        # Test PDFs
└── go.mod
```

### Running Tests

```bash
go test ./...
```

Golden Markdown files for synthetic layouts live in `cmd/pdf2md/testdata/golden`; after an intended
output change, refresh them with `go test ./cmd/pdf2md -run TestGolden -update` and review the diff.

### Benchmark

With a [READoc](https://github.com/icip-cas/READoc) clone in `./READoc` and the binary built
(`go build -o pdf2md ./cmd/pdf2md`), `make benchmark` scores the output against the ground truth
and fails if any document drops more than 0.02 below `benchmark_baseline.json`. After an intended
improvement, run `go run cmd/benchmark/main.go -update-baseline`.

Each document gets a Levenshtein similarity and a word F1 (order-free, so it separates text
extraction from reading order). `-cmd` scores another converter, with `{in}` and `{out}`
placeholders (stdout is used without `{out}`), and `-output` picks its output directory.
`benchmark_pdftotext.json` holds poppler's per-document scores as the reference to beat:

```bash
go run ./cmd/benchmark -cmd 'pdftotext {in} {out}' -output READoc/output/pdftotext \
  -baseline benchmark_pdftotext.json
```

## License

MIT License - See LICENSE file for details

## Contributing

Contributions welcome! Please open an issue before submitting major changes.
