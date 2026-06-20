# pdf2md - PDF to Markdown Converter

[![CI](https://github.com/fjacquet/pdf2md/actions/workflows/ci.yml/badge.svg)](https://github.com/fjacquet/pdf2md/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/fjacquet/pdf2md?sort=semver)](https://github.com/fjacquet/pdf2md/releases/latest)
[![License](https://img.shields.io/github/license/fjacquet/pdf2md)](https://github.com/fjacquet/pdf2md/blob/HEAD/LICENSE)

A fast, single-binary PDF to Markdown converter written in Go. Designed for AI engineers, researchers, and knowledge workers who need clean, structured Markdown from PDF documents.

## How it Works

- **Single Binary** - Zero dependencies, just download and run
- **Custom PDF Parser** - Pure Go implementation, no CGO or external libraries required
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

### Phase 4: Advanced (Planned)

- [ ] OCR support for scanned PDFs
- [x] Configurable exclusion zones (headers/footers)
- [x] Custom formatting rules

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

## License

MIT License - See LICENSE file for details

## Contributing

Contributions welcome! Please open an issue before submitting major changes.
