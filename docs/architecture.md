# pdf2md Architecture

## Overview

`pdf2md` is designed as a modular pipeline that transforms a PDF document into a structured Markdown file. The architecture emphasizes separation of concerns, testability, and extensibility.

## Pipeline

The conversion process follows a linear pipeline:

```mermaid
graph LR
    PDF[PDF File] --> Reader[PDF Reader]
    Reader --> Interpreter[Interpreter]
    Interpreter --> Extractor[Extractor]
    Extractor --> Analyzer[Layout Analyzer]
    Analyzer --> Builder[Markdown Builder]
    Builder --> MD[Markdown Output]
```

### 1. PDF Reader & Interpreter (`internal/pdf`)

- **Responsibility**: Low-level parsing of the PDF file format.
- **Components**:
  - `Reader`: Parses file structure (Trailer, Xref), reads objects.
  - `Tokenizer`: Lexical analysis of PDF streams.
  - `Parser`: Syntactic analysis of PDF objects.
  - `Interpreter`: Executes content stream operators to generate text blocks.
  - `FontManager`: Handles font loading and text decoding (CMap).

### 2. Extractor (`internal/extractor`)

- **Responsibility**: Abstraction layer that provides a uniform interface for text extraction.
- **Interface**: `PDFExtractor`
- **Implementation**: `CustomExtractor` (uses `internal/pdf`).
- **Output**: A list of `types.TextBlock` for each page.

### 3. Layout Analyzer (`internal/layout`)

- **Responsibility**: Reconstructs the logical structure of the document from raw text blocks.
- **Process**:
  1.  **Column Detection**: Identifies multi-column layouts based on X-coordinates.
  2.  **Reading Order**: Sorts blocks top-to-bottom, left-to-right (respecting columns).
  3.  **Merging**: Merges fragmented text blocks into coherent lines and paragraphs.
  4.  **Classification**: Classifies blocks into elements (Header, Paragraph, CodeBlock, List, Table) based on heuristics (font size, content patterns).

### 4. Markdown Builder (`internal/markdown`)

- **Responsibility**: Converts structured elements into Markdown syntax.
- **Features**:
  - Generates headers (`#`, `##`).
  - Formats lists (`-`, `1.`).
  - Formats code blocks (```).
  - Handles tables.

## Key Data Structures

### `types.TextBlock` (`internal/types`)

Represents a raw piece of text extracted from the PDF.

```go
type TextBlock struct {
    Text     string
    X, Y     float64
    Width    float64
    Height   float64
    FontSize float64
    FontName string
}
```

### `layout.Element` (`internal/layout`)

Represents a logical document element after analysis.

```go
type Element struct {
    Type     ElementType // Header, Paragraph, etc.
    Content  string
    Level    int         // For headers
    // ... positioning info
}
```

## Design Decisions

### Pure Go Implementation

We replaced the initial `go-fitz` (CGO) dependency with a custom pure Go parser to ensure:

- **Portability**: Easy to cross-compile.
- **Simplicity**: No external C library dependencies.
- **Control**: Ability to fix bugs and add features (like custom CMap handling) directly in the parser.

### `internal/types` Package

Introduced to resolve import cycles between `internal/pdf` and `internal/extractor`. It holds shared types like `TextBlock` that are used across the pipeline.
