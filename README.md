# pdf2md - PDF to Markdown Converter

A fast, single-binary PDF to Markdown converter written in Go. Designed for AI engineers, researchers, and knowledge workers who need clean, structured Markdown from PDF documents.

## How it Works

- **Single Binary** - Zero dependencies, just download and run
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
pdf2md input.pdf
```

Save to a file:

```bash
pdf2md input.pdf output.md
```

## Architecture

pdf2md uses a modular pipeline architecture:

1. **Extractor** (`pkg/extractor`) - Extracts raw text blocks with positioning data
2. **Analyzer** (`pkg/layout`) - Detects columns, headers, and reading order
3. **Builder** (`pkg/markdown`) - Generates clean Markdown output

The extractor is abstraction-based, making it easy to swap PDF parsing libraries.

### Current Implementation

- **Phase 1 (MVP)**: Text extraction with column detection
- Uses `ledongthuc/pdf` for pure Go compatibility
- Focuses on digital PDFs (not scanned documents)

## Roadmap

### Phase 1: Clean Text (Current)
- [x] Text extraction
- [x] Column detection
- [x] Header identification (H1-H6)
- [x] Reading order preservation

### Phase 2: Structure (Current)
- [x] List detection (bullets and numbered)
- [x] Code block detection (keyword-based)
- [ ] Font histogram analysis for better header detection
- [ ] Footnote handling

### Phase 3: Visual Elements
- [ ] Table extraction and formatting
- [ ] Image extraction to assets folder

### Phase 4: Advanced
- [ ] OCR support for scanned PDFs
- [ ] Configurable exclusion zones (headers/footers)
- [ ] Custom formatting rules

## Development

### Project Structure

```
pdf2md/
├── cmd/pdf2md/          # Main CLI application
├── pkg/
│   ├── extractor/       # PDF text extraction (pluggable)
│   ├── layout/          # Layout analysis and structure detection
│   └── markdown/        # Markdown generation
├── internal/
│   ├── config/          # Configuration
│   └── testdata/        # Test PDFs
└── go.mod
```

### Running Tests

```bash
go test ./...
```

### Adding a New Extractor

Implement the `PDFExtractor` interface in `pkg/extractor/extractor.go`:

```go
type PDFExtractor interface {
    Open(r io.ReadSeeker) error
    GetPageCount() int
    ExtractTextBlocks(page int) ([]TextBlock, error)
    Close() error
}
```

## License

MIT License - See LICENSE file for details

## Contributing

Contributions welcome! Please open an issue before submitting major changes.

## Use Cases

- **RAG Pipelines** - Convert technical documentation for vector databases
- **Research** - Extract text from academic papers preserving structure
- **Note-Taking** - Import PDFs into Obsidian, Notion, or other Markdown systems
- **Documentation** - Convert PDF manuals to searchable Markdown
