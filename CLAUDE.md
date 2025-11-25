# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

pdf2md is a "Data Liberator" tool that converts PDF documents to clean, structured Markdown. It targets AI engineers building RAG pipelines and knowledge workers using tools like Obsidian.

**Value Proposition**: Single-binary, zero-dependency Go tool that is faster than Python alternatives and easier to deploy in serverless environments.

## Build and Development Commands

```bash
# Build the binary
go build -o pdf2md ./cmd/pdf2md

# Run tests
go test ./...

# Run the tool
./pdf2md input.pdf output.md

# Add dependencies
go get <package>
```

## Architecture

The codebase follows a **pipeline architecture** with three main stages:

### 1. Extraction Layer (`pkg/extractor/`)
- **Interface-based design** - `PDFExtractor` interface allows swapping implementations
- **Current implementation**: `LedongthucExtractor` uses `ledongthuc/pdf` (pure Go, BSD-3 license)
- **Future options**: Can upgrade to UniPDF for better layout analysis if budget allows
- **Key abstraction**: `TextBlock` represents text with positioning (X, Y, Width, Height, FontSize)

### 2. Layout Analysis (`pkg/layout/`)
- **Column detection** - Analyzes X positions to find column boundaries (critical for multi-column PDFs)
- **Reading order** - Sorts blocks by column first, then Y position (top-to-bottom)
- **Header detection** - Uses font size ratios to identify H1-H6 (configurable threshold)
- **Element merging** - Combines consecutive text blocks on the same line

### 3. Markdown Generation (`pkg/markdown/`)
- Converts structured elements to clean Markdown
- Handles headers, paragraphs, lists, code blocks, tables
- Maintains proper spacing between elements

## Key Design Decisions

1. **Pure Go for MVP** - Using `ledongthuc/pdf` avoids CGO complexity and AGPL licensing issues
2. **Pluggable extractors** - Interface design allows upgrading to UniPDF or other libraries without rewriting analyzers
3. **Coordinate-based analysis** - Uses X/Y positions rather than PDF structure hints (more reliable for varied PDFs)
4. **Progressive enhancement** - Phase 1 focuses on text + columns; tables/images come later

## Technical Constraints

- `ledongthuc/pdf` requires `io.ReaderAt` (not just `io.Reader`)
- The library doesn't provide text height, so we estimate from font size
- Column detection threshold is configurable but defaults to 40 pixels
- Header detection uses 1.3x font size ratio (30% larger than body text)

## Testing Strategy

Place test PDFs in `internal/testdata/` organized by type:
- `simple/` - Single-column, basic text
- `academic/` - Two-column papers (IEEE, ACM format)
- `complex/` - Multi-column, mixed layouts
- `edge_cases/` - Unusual PDFs that break assumptions

## Phase Roadmap

**Phase 1 (Current)**: Text extraction with column detection
**Phase 2**: Lists, better header detection via font histogram
**Phase 3**: Tables and image extraction
**Phase 4**: OCR for scanned PDFs, configurable exclusion zones

## Important Notes

- **os.File implements io.ReaderAt** - Can pass directly to extractors
- **Page numbers are 1-indexed** - PDF convention used throughout
- **Y coordinates increase bottom-to-top** - PDF coordinate system (flip for screen display)
- **Font size is in points** - 1 point = 1/72 inch
