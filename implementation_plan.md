# Kong CLI Implementation Plan

**Goal**: Replace manual argument parsing with `alecthomas/kong` for a robust CLI.

## 1. Dependencies

- Add `github.com/alecthomas/kong`.

## 2. CLI Structure

Define a struct to map command line arguments:

```go
var CLI struct {
    Input  string `arg:"" help:"Input PDF file path." type:"path"`
    Output string `arg:"" optional:"" help:"Output Markdown file path." type:"path"`
    Debug  bool   `help:"Enable debug logging."`
}
```

## 3. Changes

- **Modify `cmd/pdf2md/main.go`**:
  - Remove `os.Args` logic.
  - Initialize `kong`.
  - Use `CLI.Input` and `CLI.Output`.
  - Configure `slog` based on `CLI.Debug`.

## 4. Verification

- Run `go run cmd/pdf2md/main.go --help` to verify help message.
- Run `go run cmd/pdf2md/main.go input.pdf` to verify functionality.
