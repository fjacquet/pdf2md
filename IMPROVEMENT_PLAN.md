# pdf2md Improvement Plan: Better Than Adobe

## Executive Summary

After deep analysis of the codebase, PDF 1.5 specification, and competitive landscape (Marker, MinerU, Docling), I've identified critical gaps that prevent pdf2md from achieving best-in-class accuracy. This plan prioritizes improvements by impact.

---

## Current State Analysis

### Strengths

- Pure Go, single binary, zero dependencies
- Clean architecture: PDF Parser → Extractor → Layout → Markdown
- Good PDF 1.5 support: XRef streams, compressed objects, hybrid files
- Basic text extraction with font handling (Type0, Type1, CID, ToUnicode CMap)

### Critical Gaps (by severity)

| Category           | Current                                 | Best-in-class (Marker/MinerU)           |
| ------------------ | --------------------------------------- | --------------------------------------- |
| **Stream Filters** | FlateDecode only (40%)                  | All filters (LZW, ASCII85, JBIG2, etc.) |
| **Tables**         | 4-space gap heuristic (40-50% accuracy) | ML-based cell detection (95%+)          |
| **Layout**         | Rule-based XY-cut                       | LayoutLMv3 / VLM models                 |
| **Math/Formulas**  | Font-name heuristic (50-60%)            | LaTeX recognition models                |
| **OCR**            | None                                    | Surya/Tesseract integration             |
| **Form XObjects**  | Not implemented                         | Full support                            |

---

## Improvement Roadmap

### Phase 1: Core PDF Parser Fixes (High Impact, Low Effort)l

**Goal**: Support more real-world PDFs without ML dependencies

#### 1.1 Stream Filters (Priority: CRITICAL)

Currently only FlateDecode works. Missing filters cause silent content loss.

**Files to modify:**

- `internal/pdf/filters.go` - Add new filter implementations
- `internal/pdf/content.go:93` - Wire up multi-filter chains

**Filters to implement:**

```
Priority 1 (Common):
- LZWDecode - Used in older PDFs, straightforward LZW implementation
- ASCII85Decode - Text encoding, simple Base85 decoder
- ASCIIHexDecode - Hex encoding, trivial

Priority 2 (Images):
- DCTDecode - Already passes through (JPEG), needs color space handling
- CCITTFaxDecode - Fax compression for scanned docs
- JBIG2Decode - Modern scanned document compression
- JPXDecode - JPEG2000
```

#### 1.2 Name Hex Escapes (Priority: HIGH)

`internal/pdf/tokenizer.go:154` - Missing `#xx` escape handling causes character corruption.

```go
// Current: /Name#20With#20Spaces → fails
// Fixed: /Name#20With#20Spaces → "Name With Spaces"
```

#### 1.3 Form XObjects (Priority: HIGH)

`internal/pdf/ops_xobject.go:80` - Form XObjects contain nested content (templates, headers, footers). Currently ignored.

**Implementation**: Recursively interpret Form XObject content streams with inherited graphics state.

#### 1.4 Inline Images (Priority: MEDIUM)

`BI ... ID ... EI` operators not implemented. Common in single-page documents.

---

### Phase 2: Layout Analysis Improvements (High Impact, Medium Effort)

**Goal**: Fix fragile heuristics that cause misclassification

#### 2.1 Table Detection Overhaul (Priority: CRITICAL)

Current approach (4-space gap detection) has ~40-50% accuracy. Equations get misclassified as tables.

**New approach:**

```
1. Grid-based detection:
   - Build candidate grid from text block positions
   - Validate: consistent column widths, row alignment
   - Score: coverage ratio, regularity score

2. Separator detection:
   - Detect horizontal/vertical lines from graphics path data
   - Use lines to define cell boundaries

3. Content validation:
   - Tables have structured data (numbers, short text)
   - Equations have math symbols, operators
```

**Files:**

- `internal/layout/analyzer.go` - New `classifyAsTable()` with grid detection
- `internal/layout/tables.go` - New file for table structure analysis
- `internal/types/table.go` - Table cell/row/column types

#### 2.2 Equation vs Table Disambiguation (Priority: HIGH)

Current issue: `z^2 = 0. (2)` classified as table due to spacing.

**Solution:**

```go
// Before classifying as table, check:
func isLikelyEquation(blocks []TextBlock) bool {
    mathScore := 0
    // Check for math fonts (cmmi, cmsy, symbol, math)
    // Check for operators (=, +, -, ×, ÷, ∫, ∑)
    // Check for equation numbering pattern: (1), (2.3), [1]
    // Check for inline $ delimiters
    return mathScore > threshold
}
```

#### 2.3 Configurable Thresholds (Priority: MEDIUM)

Replace hardcoded values with configurable options:

```go
type LayoutConfig struct {
    HeaderSizeRatio    float64 // default: 1.2
    ColumnGapThreshold float64 // default: 10.0
    WideGapMultiplier  float64 // default: 3.0
    TableMinColumns    int     // default: 2
    TableMinSpaces     int     // default: 4
}
```

**Files:**

- `internal/layout/config.go` - New configuration struct
- `internal/layout/analyzer.go` - Use config instead of constants

#### 2.4 Graphics-Based Structure Detection (Priority: MEDIUM)

Use vector graphics (lines, rectangles) already captured but unused:

```go
// internal/layout/graphics_analysis.go
func DetectTableFromGraphics(graphics []Graphics, textBlocks []TextBlock) *TableStructure {
    // Find horizontal lines → row separators
    // Find vertical lines → column separators
    // Find rectangles → cell boundaries
    // Map text blocks to cells
}
```

---

### Phase 3: Font & Encoding Completeness (Medium Impact, Medium Effort)

#### 3.1 Type3 Fonts (Priority: MEDIUM)

Custom user-defined fonts using drawing operators. Currently ignored.

**Files:** `internal/pdf/font.go`

#### 3.2 Extended Glyph Mapping (Priority: MEDIUM)

Current `GlyphToUnicode` has 262 entries. Many documents use extended glyphs.

**Solution:** Embed Adobe Glyph List (4,000+ mappings) from AGL specification.

#### 3.3 Ligatures (Priority: LOW)

Handle fi, fl, ff, ffi, ffl ligatures that appear as single glyphs.

---

### Phase 4: Future Considerations (Out of Scope for Now)

These are documented for future reference but **not part of current plan** (Pure Go constraint):

- OCR integration (would require external subprocess)
- ML-assisted layout (requires ONNX or external models)
- Reading order models (neural network based)

**Current focus**: Maximize heuristic accuracy to reach 85% target with zero dependencies.

---

### Phase 5: Output Quality (Medium Impact, Low Effort)

#### 5.1 Markdown Table Generation (Priority: HIGH)

Current: Tables output as text with spaces
Target: Proper markdown tables with alignment

```markdown
| Header 1 | Header 2 | Header 3 |
| -------- | -------- | -------- |
| Cell 1   | Cell 2   | Cell 3   |
```

**Files:** `internal/markdown/builder.go` - New `formatTable()` function

#### 5.2 Math Output Modes (Priority: MEDIUM)

```go
type MathOutputMode string
const (
    MathPassthrough MathOutputMode = "passthrough" // Keep as-is
    MathLatex       MathOutputMode = "latex"       // Wrap in $...$
    MathUnicode     MathOutputMode = "unicode"     // Convert to Unicode symbols
)
```

#### 5.3 Link Handling Improvements (Priority: LOW)

Current center-point intersection misses partial overlaps. Use proper rectangle intersection.

---

## Implementation Priority Matrix (Pure Go, 85% Target)

### Sprint 1: Foundation (Parallel Tracks)

**Track A - Filters:**
| Item | File | Effort |
|------|------|--------|
| LZWDecode filter | `internal/pdf/filters.go` | 2-3 hours |
| ASCII85Decode filter | `internal/pdf/filters.go` | 1 hour |
| ASCIIHexDecode filter | `internal/pdf/filters.go` | 30 min |
| Name hex escapes (#xx) | `internal/pdf/tokenizer.go` | 1 hour |

**Track B - Tables:**
| Item | File | Effort |
|------|------|--------|
| Grid-based table detection | `internal/layout/tables.go` (new) | 4-5 hours |
| Graphics line detection | `internal/layout/graphics_analysis.go` (new) | 3 hours |
| Markdown table output | `internal/markdown/builder.go` | 2 hours |

### Sprint 2: Disambiguation & Structure

| Item                         | File                              | Effort  |
| ---------------------------- | --------------------------------- | ------- |
| Equation vs table scoring    | `internal/layout/analyzer.go`     | 3 hours |
| Form XObjects support        | `internal/pdf/ops_xobject.go`     | 4 hours |
| Extended glyph mapping (AGL) | `internal/pdf/encodings.go`       | 2 hours |
| Configurable thresholds      | `internal/layout/config.go` (new) | 2 hours |

### Sprint 3: Polish & Edge Cases

| Item                       | File                          | Effort  |
| -------------------------- | ----------------------------- | ------- |
| Inline images (BI/EI)      | `internal/pdf/interpreter.go` | 3 hours |
| Multi-filter chains        | `internal/pdf/content.go`     | 2 hours |
| Link rect intersection fix | `internal/layout/links.go`    | 1 hour  |
| Type3 font basic support   | `internal/pdf/font.go`        | 4 hours |

---

## Success Metrics

### Benchmark Targets (READoc dataset)

| Metric                 | Current (estimated) | Target                   |
| ---------------------- | ------------------- | ------------------------ |
| Levenshtein similarity | ~60-70%             | >85%                     |
| Table accuracy         | ~40-50%             | >80%                     |
| Header detection       | ~70-80%             | >90%                     |
| Math preservation      | ~50-60%             | >75%                     |
| Processing speed       | Fast                | Maintain (no regression) |

### Testing Strategy

1. **Unit tests** for each new filter and feature
2. **Integration tests** with diverse PDF corpus:
   - Academic papers (arXiv subset)
   - Technical documentation (GitHub subset)
   - Scanned documents (OCR validation)
   - Complex tables (financial reports)
3. **Benchmark regression** on every PR

---

## User Requirements (Confirmed)

- **Approach**: Balanced - work on filters and tables in parallel
- **Dependencies**: Pure Go only - no ML models, improve heuristics
- **Target**: 85% Levenshtein similarity on READoc benchmark

---

## References

- [PDF 1.5 Object Streams](https://qpdf.readthedocs.io/en/stable/object-streams.html)
- [How Marker Works](https://kevinhu.io/notes/how-marker-works/)
- [MinerU GitHub](https://github.com/opendatalab/MinerU)
- [PDF Converter Benchmarks](https://ai.gopubby.com/benchmarking-pdf-to-markdown-document-converters-fc65a2c73bf2)
- [Marker Architecture](https://journal.hexmos.com/marker-pdf-document-ai/)
