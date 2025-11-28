# Recommended Stack for pdf2md

This document outlines the technology stack chosen for `pdf2md` and the rationale behind these choices.

## 🏆 The Current Stack

| Component           | Choice                | Why?                                                                          |
| :------------------ | :-------------------- | :---------------------------------------------------------------------------- |
| **CLI Framework**   | **`kong`**            | Simpler and more type-safe than Cobra for this use case.                      |
| **Logging**         | **`log/slog`**        | Standard library (Go 1.21+), structured, and fast.                            |
| **PDF Engine**      | **Custom Pure Go**    | **Best of both worlds**: Portability of pure Go + Control of a custom engine. |
| **Layout Analysis** | **Custom Heuristics** | Tailored specifically for Markdown structure (headers, lists, code).          |

---

## 1. The Core: PDF Parsing Strategy

We evaluated three options for the critical task of PDF parsing:

### Option A: CGO Wrapper (`go-fitz`)

- **Pros**: High fidelity, mature (MuPDF backend).
- **Cons**: Requires CGO, difficult cross-compilation, external dependency.
- **Verdict**: **Rejected**. We initially used this but moved away to ensure portability and ease of distribution.

### Option B: Existing Pure Go Libraries (`ledongthuc/pdf`, `unidoc`)

- **Pros**: Pure Go.
- **Cons**: `ledongthuc/pdf` is unmaintained and has limited extraction capabilities. `unidoc` has a commercial license (AGPL/Commercial).
- **Verdict**: **Rejected**. Existing open-source pure Go libraries were not robust enough for our layout analysis needs.

### Option C: Custom Pure Go Parser (Selected)

- **Pros**:
  - **Zero Dependencies**: No CGO, no external libs.
  - **Full Control**: We implemented exactly what we need (CMap handling, content stream interpretation).
  - **Optimized**: Tuned specifically for text extraction and layout analysis.
- **Cons**: High initial development effort (already completed).
- **Verdict**: **Selected**. This gives us the portability of Go with the precision required for high-quality Markdown generation.

---

## 2. CLI Framework: `kong`

We chose `alecthomas/kong` over `cobra` because:

- **Type Safety**: Flags are mapped directly to struct fields.
- **Simplicity**: Less boilerplate code.
- **Expressiveness**: struct tags make the CLI definition self-documenting.

```go
type CLI struct {
    Input  string `arg:"" help:"Input PDF file"`
    Output string `arg:"" optional:"" help:"Output Markdown file"`
    Debug  bool   `help:"Enable debug logging"`
}
```

---

## 3. Logging: `log/slog`

We use the standard library's `log/slog` package.

- **Structured**: Good for debugging and machine parsing if needed.
- **Zero Dependency**: No need for `zerolog` or `zap` for a CLI tool.
- **Levels**: Easy control over Debug/Info/Error levels.

---

## 4. Layout Analysis

Instead of relying on the PDF parser to give us "paragraphs" (which they often get wrong), we extract raw **Text Blocks** with precise coordinates (`X`, `Y`, `FontSize`) and perform our own analysis:

1.  **Column Detection**: Histogram analysis of X-coordinates.
2.  **Reading Order**: Sorting based on columns and Y-coordinates.
3.  **Element Classification**:
    - **Headers**: Font size > Body font size.
    - **Lists**: Regex patterns (`^•`, `^\d+\.`).
    - **Code**: Monospace font detection (future) or content heuristics.

This approach gives us the flexibility to handle complex documents like academic papers and technical manuals.
