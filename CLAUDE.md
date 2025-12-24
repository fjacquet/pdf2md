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

# ONNX Layout Detection (ML-enhanced)
./pdf2md input.pdf output.md          # ONNX enabled by default if model available
./pdf2md --no-onnx input.pdf          # Disable ONNX, use rule-based only
./pdf2md --model-path=/path/model.onnx input.pdf  # Custom model path
./pdf2md --runtime-path=/path/libonnxruntime.so input.pdf  # Custom runtime
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
- **ONNX integration** - `analyzer_onnx.go` fuses ML detections with rule-based analysis
- **TOC detection** - `analyzer_toc.go` identifies and removes table of contents
- **Math/formula handling** - `formulas.go` for LaTeX-style math expressions
- **Link merging** - `links.go` and `analyzer_merge_links.go` for hyperlinks
- **Exclusion zones** - Configurable top/bottom regions to skip (headers/footers)

### 3b. ONNX Layout Detection (`internal/onnx/`)

ML-based layout detection using DocLayout-YOLO model:

- **detector.go** - `Detector` interface and onnxruntime-purego implementation
- **preprocessing.go** - Image preparation: letterbox resize, normalize, BCHW tensor
- **postprocessing.go** - NMS, confidence filtering, box scaling to page coordinates
- **model.go** - Auto-download model from Hugging Face if not present
- **classes.go** - Mapping DocLayout classes to ElementType

### 3c. PDF Rendering (`internal/render/`)

PDF to image rendering for ONNX input:

- **renderer.go** - `PageRenderer` interface and configuration
- **renderer_pdfium.go** - go-pdfium WebAssembly/Wazero implementation
- Cross-platform: WebAssembly runtime, no native dependencies

### 4. Markdown Generation (`internal/markdown/`)

- `builder.go` - Converts `Element` arrays to Markdown text
- `formatter.go` - Text formatting utilities
- Handles headers (H1-H6), paragraphs, lists, code blocks, tables, images

### Shared Types (`internal/types/`)

- `TextBlock` - Text with position (X, Y, Width, Height), FontSize, FontName, LinkURI
- `Image` - Extracted image with binary data and format
- `Link` - Hyperlink with URI and bounding rectangle
- `Graphics` - Vector graphics data (`graphics.go`)
- `Element` - Layout element with Type, ONNX metadata (`element.go`)
- `BoundingBox` - ONNX detection box with confidence and class (`detection.go`)
- `PageDetections` - All ONNX detections for a page (`detection.go`)
- `DocLayoutClass` - Enum for DocLayout-YOLO classes (title, table, figure, etc.)

## Key Design Decisions

1. **Pure Go** - Custom PDF parser avoids CGO complexity and licensing issues
2. **Pluggable extractors** - Interface design allows future upgrades without rewriting analyzers
3. **Coordinate-based analysis** - Uses X/Y positions rather than PDF structure hints
4. **Rule-based classification** - Extensible system for element type detection
5. **ONNX Hybrid** - ML detection fused with rule-based for best accuracy, graceful fallback
6. **WebAssembly PDF rendering** - go-pdfium with Wazero for cross-platform compatibility

## Technical Notes

- **Page numbers are 1-indexed** - PDF convention used in `ExtractTextBlocks(pageIndex)`
- **Y coordinates increase bottom-to-top** - PDF coordinate system
- **Font size is in points** - 1 point = 1/72 inch
- **Column detection threshold** - Default 10.0 pixels (`ColumnGapThreshold`)
- **Header detection ratio** - Default 1.2x body font size (`HeaderSizeRatio`)

### ONNX Technical Notes

- **Model**: DocLayout-YOLO from `wybxc/DocLayout-YOLO-DocStructBench-onnx` (Apache 2.0)
- **Input size**: 1024x1024 pixels with letterbox padding
- **Confidence threshold**: 0.25 for detection, 0.5 for override
- **NMS threshold**: 0.45 IoU for overlapping boxes
- **Runtime**: Requires `libonnxruntime` native library (auto-detected or via `--runtime-path`)
- **Model auto-download**: Stored in `~/.pdf2md/models/` or `PDF2MD_MODEL_PATH` env var
- **Fallback**: Graceful degradation to rule-based if ONNX unavailable

### Hardware Acceleration

ONNX Runtime supports hardware acceleration on various platforms:

**macOS (CoreML)**:
- Automatically attempted on macOS (M1/M2/M3/Intel) when available
- Leverages Apple Neural Engine (ANE) on Apple Silicon for fastest inference
- Falls back to CPU if CoreML unavailable (e.g., Homebrew build)
- Enabled by default via `Config.UseCoreML = true`
- **Note**: Homebrew's onnxruntime is built without CoreML; for CoreML, use official releases from Microsoft

**Installation**:
```bash
# macOS (Homebrew) - CPU only, no CoreML
brew install onnxruntime

# macOS with CoreML - download from Microsoft releases
# https://github.com/microsoft/onnxruntime/releases
# Look for: onnxruntime-osx-arm64-*.tgz (Apple Silicon)
# or: onnxruntime-osx-x64-*.tgz (Intel)

# Linux (Ubuntu/Debian)
apt install libonnxruntime-dev

# Custom path
export ORT_LIB_PATH=/path/to/libonnxruntime.dylib
```

**Environment Variables**:
- `PDF2MD_MODEL_PATH` - Custom path to ONNX model file
- `ORT_LIB_PATH` - Custom path to ONNX Runtime library

**Execution Provider Priority**:
1. CoreML (macOS with Apple Silicon/Intel GPU)
2. CUDA (NVIDIA GPUs, if compiled with CUDA support)
3. CPU (fallback, always available)

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

The analyzer classifies blocks into these element types (defined in `internal/types/element.go`):

- `ElementTypeHeader` - H1-H6 headers (detected by font size ratio or numbering patterns)
- `ElementTypeParagraph` - Regular text paragraphs
- `ElementTypeCodeBlock` - Code blocks (detected by monospace fonts or content patterns)
- `ElementTypeList` - Bulleted or numbered lists
- `ElementTypeTable` - Table structures (enhanced by ONNX detection)
- `ElementTypeAdmonition` - NOTE, WARNING, TIP, etc. callouts
- `ElementTypeImage` - Extracted images
- `ElementTypeFigure` - Figures detected by ONNX
- `ElementTypeCaption` - Figure/table captions detected by ONNX
- `ElementTypeEquation` - Mathematical formulas detected by ONNX
- `ElementTypeFootnote` - Table footnotes detected by ONNX
