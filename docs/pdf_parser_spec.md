# Custom PDF Parser Technical Specification

## 1. Overview

This document outlines the technical specification for the custom PDF parser implemented in `pdf2md`. The goal was to replace the `go-fitz` dependency (which required CGO) with a pure Go implementation that is robust enough to handle common PDF structures, including cross-reference tables, content streams, and font encodings.

## 2. Architecture

The parser is located in `internal/pdf` and consists of the following components:

- **Reader (`reader.go`)**: The entry point. Handles file opening, reading the trailer, parsing the Cross-Reference Table (xref), and providing access to objects.
- **Tokenizer (`tokenizer.go`)**: Scans the PDF byte stream and produces tokens (keywords, literals, numbers, strings).
- **Parser (`parser.go`)**: Consumes tokens to parse PDF objects (Dictionaries, Arrays, Streams, Indirect Objects).
- **Interpreter (`interpreter.go`)**: Processes content streams (page content) by executing operators (e.g., `BT`, `Tj`, `Td`).
- **Font Manager (`font_manager.go`)**: Handles font loading, CMap parsing, and text decoding (mapping raw bytes to Unicode).
- **CMap Parser (`cmap_parser.go`)**: Parses ToUnicode CMaps to support custom character mappings.

## 3. Key Features Implemented

### 3.1. File Structure Parsing

- **Trailer Parsing**: Locates the `trailer` dictionary and the `startxref` offset.
- **Xref Table**: Parses standard xref tables to map Object IDs to file offsets.
- **Indirect Objects**: Reads objects like `1 0 obj ... endobj`.

### 3.2. Object Types

- **Booleans**: `true`, `false`
- **Integers & Reals**: `123`, `12.34`
- **Strings**: Literal `(Hello)` and Hex `<48656C6C6F>`
- **Names**: `/Name`
- **Arrays**: `[1 2 3]`
- **Dictionaries**: `<< /Key /Value >>`
- **Streams**: `<< ... >> stream ... endstream`
- **Null**: `null`
- **Indirect References**: `1 0 R`

### 3.3. Content Stream Interpretation

- **Text State**: Tracks Text Matrix (`Tm`), Text Line Matrix (`Tlm`), Font (`Tf`), Font Size (`Tfs`), and Leading (`Tl`).
- **Operators**:
  - `BT`, `ET`: Begin/End Text object.
  - `Tf`: Set font and size.
  - `Td`, `TD`, `T*`: Move text position.
  - `Tm`: Set text matrix.
  - `Tj`, `TJ`, `'`, `"`: Show text (with and without positioning/spacing).
  - `q`, `Q`: Save/Restore graphics state (partial support).
  - `cm`: Concatenate matrix (partial support).

### 3.4. Font Handling

- **Simple Fonts**: Supports Type1 and TrueType fonts.
- **Composite Fonts**: Basic support for Type0 fonts.
- **Encodings**:
  - `WinAnsiEncoding`
  - `MacRomanEncoding`
  - `StandardEncoding`
- **ToUnicode CMaps**: Parses embedded CMap streams to map character codes to Unicode strings. This is crucial for correctly extracting text from PDFs with custom encodings.

## 4. Usage

The parser is integrated via the `CustomExtractor` in `internal/extractor/custom.go`.

```go
reader, err := pdf.NewReader("file.pdf")
if err != nil {
    return err
}
defer reader.Close()

// Get total pages
count, err := reader.GetPageCount()

// Extract text from a page
content, err := reader.GetPageContent(pageIndex)
interpreter := pdf.NewInterpreter(reader.FontManager)
blocks, err := interpreter.Process(content)
```

## 5. Future Improvements

- **Graphics State**: More complete support for graphics state operators (color, line width, etc.) to aid in layout analysis (e.g., detecting table borders).
- **Images**: Extracting inline images (`BI`...`EI`) and XObject images.
- **Filters**: Support for more stream filters (currently supports FlateDecode; need LZWDecode, ASCII85Decode, etc.).
- **Object Streams**: Support for compressed object streams (PDF 1.5+).
