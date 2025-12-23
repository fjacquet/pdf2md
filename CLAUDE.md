# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

pdf2md is a PDF to Markdown converter written in pure Go. It targets AI engineers building RAG pipelines and knowledge workers using tools like Obsidian.

**Value Proposition**: Single-binary, zero-dependency Go tool with a custom PDF parser—no CGO or external libraries required.

## Build and Development Commands

```bash
# Build the binary
make build
# or: go build -o pdf2md cmd/pdf2md/main.go

# Run all tests
make test
# or: go test ./...

# Run a single test
go test -run TestFunctionName ./internal/layout/

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Format code
make fmt

# Run linter
make lint

# Run benchmark
make benchmark

# Test with a sample PDF
make run-test
```

## CLI Usage

```bash
./pdf2md input.pdf                    # Output to stdout
./pdf2md input.pdf output.md          # Save to file
./pdf2md --debug input.pdf            # Enable debug logging
./pdf2md --exclude-top=50 input.pdf   # Exclude header region
./pdf2md --exclude-bottom=30 input.pdf # Exclude footer region
```

## Architecture

The codebase follows a **pipeline architecture** with four main stages:

### 1. PDF Parser (`internal/pdf/`)

Custom pure Go PDF parser that handles:

- Cross-reference table parsing (`reader.go`)
- Object reading: dictionaries, streams, arrays (`objects.go`, `parser.go`)
- Content stream interpretation (`interpreter.go`, `content.go`)
- Font handling: CMap, Type1, TrueType encodings (`font.go`, `encodings.go`)
- Graphics state management (`graphics_state.go`, `matrix.go`)
- Stream filters: FlateDecode, ASCII85, etc. (`filters.go`)

### 2. Extraction Layer (`internal/extractor/`)

- **Interface**: `PDFExtractor` in `extractor.go` allows swapping implementations
- **Implementation**: `CustomExtractor` in `custom.go` uses the internal PDF parser
- **Output**: `PageContent` containing `TextBlock`, `Image`, and `Graphics` arrays
- Image/graphics saving utilities in `images.go` and `graphics.go`

### 3. Layout Analysis (`internal/layout/`)

- **Column detection** - Analyzes X positions to find column boundaries
- **Reading order** - `sorter.go` sorts blocks by column then Y position (top-to-bottom)
- **Header detection** - Uses font size ratios (default 1.2x body text)
- **Element classification** - Rule-based system in `analyzer.go` for headers, code blocks, lists, tables, admonitions
- **TOC detection** - `analyzer_toc.go` identifies and removes table of contents
- **Math/formula handling** - `formulas.go` for LaTeX-style math expressions
- **Link merging** - `links.go` and `analyzer_merge_links.go` for hyperlinks
- **Exclusion zones** - Configurable top/bottom regions to skip (headers/footers)

### 4. Markdown Generation (`internal/markdown/`)

- `builder.go` - Converts `Element` arrays to Markdown text
- `formatter.go` - Text formatting utilities
- Handles headers (H1-H6), paragraphs, lists, code blocks, tables, images

### Shared Types (`internal/types/`)

- `TextBlock` - Text with position (X, Y, Width, Height), FontSize, FontName, LinkURI
- `Image` - Extracted image with binary data and format
- `Link` - Hyperlink with URI and bounding rectangle
- `Graphics` - Vector graphics data (`graphics.go`)

## Key Design Decisions

1. **Pure Go** - Custom PDF parser avoids CGO complexity and licensing issues
2. **Pluggable extractors** - Interface design allows future upgrades without rewriting analyzers
3. **Coordinate-based analysis** - Uses X/Y positions rather than PDF structure hints
4. **Rule-based classification** - Extensible system for element type detection

## Technical Notes

- **Page numbers are 1-indexed** - PDF convention used in `ExtractTextBlocks(pageIndex)`
- **Y coordinates increase bottom-to-top** - PDF coordinate system
- **Font size is in points** - 1 point = 1/72 inch
- **Column detection threshold** - Default 10.0 pixels (`ColumnGapThreshold`)
- **Header detection ratio** - Default 1.2x body font size (`HeaderSizeRatio`)

## Testing

Test PDFs go in `testdata/`. Tests use the `_test.go` convention:

- `analyzer_test.go` - Layout analysis tests
- `analyzer_fuzz_test.go` - Fuzz testing for robustness
- `builder_test.go` - Markdown generation tests
- `reader_test.go` - PDF parsing tests

## Library Usage

See `examples/basic_usage.go` for programmatic usage. Key steps:

```go
ext, _ := extractor.NewCustomExtractor(pdfPath)
defer ext.Close()

analyzer := layout.NewAnalyzer()
builder := markdown.NewBuilder()

for page := 1; page <= pageCount; page++ {
    content, _ := ext.ExtractTextBlocks(page)  // 1-indexed
    elements := analyzer.Analyze(content)
    elements = analyzer.MergeElements(elements)
    // ...
}

md := builder.Build(allElements)
```

## Element Types

The analyzer classifies blocks into these element types (defined in `internal/layout/analyzer.go`):

- `ElementTypeHeader` - H1-H6 headers (detected by font size ratio or numbering patterns)
- `ElementTypeParagraph` - Regular text paragraphs
- `ElementTypeCodeBlock` - Code blocks (detected by monospace fonts or content patterns)
- `ElementTypeList` - Bulleted or numbered lists
- `ElementTypeTable` - Table structures
- `ElementTypeAdmonition` - NOTE, WARNING, TIP, etc. callouts
- `ElementTypeImage` - Extracted images
