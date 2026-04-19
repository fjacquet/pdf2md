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

# OCR fallback (scanned / image-only pages)
./pdf2md input.pdf                            # OCR auto-used if models present in ~/.pdf2md/models/ocr/
./pdf2md --no-ocr input.pdf                   # Disable OCR fallback entirely
./pdf2md --ocr-det-model-path=/det.onnx input.pdf
./pdf2md --ocr-rec-model-path=/rec.onnx input.pdf
./pdf2md --ocr-dict-path=/dict.txt input.pdf  # Required only if model uses a non-Latin dict
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
- Stream filters (`filters.go`):
  - **Full decoders**: FlateDecode, ASCIIHexDecode, ASCII85Decode, LZWDecode,
    RunLengthDecode
  - **Passthrough** (bytes already a valid image file): DCTDecode (JPEG),
    JPXDecode (JPEG2000, tagged `format="jp2"`)
  - **Explicit `ErrFilterUnsupported`**: JBIG2Decode, CCITTFaxDecode (G4
    decoder port is deferred; see `internal/pdf/filters.go::decodeCCITTFax`)
  - **No-op**: `/Crypt` — decryption happens upstream in the security handler,
    so the stream is plaintext by the time it reaches the filter dispatch
- Document encryption (`crypt.go`):
  - **Supported**: standard security handler V=1/V=2 (RC4 40-/128-bit, R=2/3)
    and V=4 R=4 (RC4 or AES-128, `/CFM` = `V2` or `AESV2`)
  - **Constraint**: only the empty user password is tried — covers the common
    "owner-password-only" PDFs exported by Word. Non-empty user passwords and
    V=5/R=6 (AES-256) return `ErrEncrypted`.
  - Wired automatically when the trailer contains `/Encrypt`; strings and
    stream data are decrypted inside `Reader.ReadObject`. Metadata streams
    are exempt when `/EncryptMetadata false`.

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

### 3d. OCR (`internal/ocr/`)

ONNX-based OCR fallback for scanned / image-only pages. Pipeline mirrors
`internal/onnx/`: pure preprocessing/postprocessing + stateful `PaddleRecognizer`
that owns the ONNX sessions, dependency injection via `Deps` for testability.

- **recognizer.go** - `Recognizer` interface, `Config`, `Deps` (DI)
- **paddle.go** - `PaddleRecognizer` orchestrating det + rec ONNX sessions
- **detection.go** - DB text-detection pre/postprocessing (pure functions,
  connected-components labelling)
- **recognition.go** - CRNN recognition preprocessing: crop + resize to 48×N
  + normalize to [-1, 1]
- **ctc.go** - Greedy CTC decode (pure function)
- **model.go** - Path resolution (no auto-download: users install once,
  see `docs/ocr-setup.md`)
- **dict.go** - Embedded Latin character dictionary (FR/EN/DE + more)

Shared cache primitives live in `internal/modelcache/` (used by both
`internal/onnx/` and `internal/ocr/`).

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

### OCR Technical Notes

- **Models**: PaddleOCR PP-OCRv5 mobile (det + rec), Apache 2.0
- **No auto-download**: user installs once at `~/.pdf2md/models/ocr/{det,rec}.onnx`
  (see `docs/ocr-setup.md`). Stable community ONNX mirrors do not exist.
- **Languages**: Latin scripts (FR/EN/DE + Spanish, Italian, Nordic, etc.)
  via the embedded `latin_dict.txt`
- **Detection**: DB (Differentiable Binarization), short-side resize to 960,
  connected-components postprocessing
- **Recognition**: CRNN + CTC greedy decode, fixed input height 48, max width 320
- **Runs only on image-only pages**: zero cost when the PDF has extractable text
- **Fallback**: if models are missing, the recognizer stays unavailable and
  scanned pages pass through untouched (no error)

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
- `PDF2MD_MODEL_PATH` - Custom path to ONNX layout model file
- `PDF2MD_OCR_DET_PATH` - Custom path to OCR text-detection ONNX model
- `PDF2MD_OCR_REC_PATH` - Custom path to OCR text-recognition ONNX model
- `PDF2MD_OCR_DICT_PATH` - Custom path to OCR character dictionary
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

<!-- code-review-graph MCP tools -->
## MCP Tools: code-review-graph

**IMPORTANT: This project has a knowledge graph. ALWAYS use the
code-review-graph MCP tools BEFORE using Grep/Glob/Read to explore
the codebase.** The graph is faster, cheaper (fewer tokens), and gives
you structural context (callers, dependents, test coverage) that file
scanning cannot.

### When to use graph tools FIRST

- **Exploring code**: `semantic_search_nodes` or `query_graph` instead of Grep
- **Understanding impact**: `get_impact_radius` instead of manually tracing imports
- **Code review**: `detect_changes` + `get_review_context` instead of reading entire files
- **Finding relationships**: `query_graph` with callers_of/callees_of/imports_of/tests_for
- **Architecture questions**: `get_architecture_overview` + `list_communities`

Fall back to Grep/Glob/Read **only** when the graph doesn't cover what you need.

### Key Tools

| Tool | Use when |
|------|----------|
| `detect_changes` | Reviewing code changes — gives risk-scored analysis |
| `get_review_context` | Need source snippets for review — token-efficient |
| `get_impact_radius` | Understanding blast radius of a change |
| `get_affected_flows` | Finding which execution paths are impacted |
| `query_graph` | Tracing callers, callees, imports, tests, dependencies |
| `semantic_search_nodes` | Finding functions/classes by name or keyword |
| `get_architecture_overview` | Understanding high-level codebase structure |
| `refactor_tool` | Planning renames, finding dead code |

### Workflow

1. The graph auto-updates on file changes (via hooks).
2. Use `detect_changes` for code review.
3. Use `get_affected_flows` to understand impact.
4. Use `query_graph` pattern="tests_for" to check coverage.
