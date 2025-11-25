# Contributing to pdf2md

Thank you for considering contributing to pdf2md! This document provides guidelines for contributing.

## Development Setup

1. **Prerequisites**
   - Go 1.21 or later
   - Git

2. **Clone and Build**
   ```bash
   git clone https://github.com/fjacquet/pdf2md
   cd pdf2md
   go build -o pdf2md ./cmd/pdf2md
   ```

3. **Run Tests**
   ```bash
   go test ./...
   ```

## Project Structure

- `cmd/pdf2md/` - CLI application entry point
- `pkg/extractor/` - PDF text extraction (pluggable interface)
- `pkg/layout/` - Layout analysis and structure detection
- `pkg/markdown/` - Markdown generation
- `internal/testdata/` - Test PDF files

## Making Changes

1. **Create a branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Follow Go conventions (run `go fmt`)
   - Add tests for new functionality
   - Update documentation if needed

3. **Test your changes**
   ```bash
   go test ./...
   go build -o pdf2md ./cmd/pdf2md
   ```

4. **Commit and push**
   ```bash
   git add .
   git commit -m "feat: describe your changes"
   git push origin feature/your-feature-name
   ```

## Adding a New PDF Extractor

To add support for a new PDF library:

1. Create a new file in `pkg/extractor/` (e.g., `unidoc.go`)
2. Implement the `PDFExtractor` interface:
   ```go
   type PDFExtractor interface {
       Open(r io.ReadSeeker) error
       GetPageCount() int
       ExtractTextBlocks(page int) ([]TextBlock, error)
       Close() error
   }
   ```
3. Add a constructor function (e.g., `NewUnidocExtractor()`)
4. Add tests in `extractor_test.go`
5. Update README with new extractor option

## Testing with Real PDFs

Place test PDFs in `internal/testdata/` organized by category:
- `simple/` - Single-column documents
- `academic/` - Research papers with multi-column layouts
- `complex/` - Mixed layouts, tables, images
- `edge_cases/` - PDFs that expose bugs

Run tests with:
```bash
go test -v ./pkg/extractor -run TestWithRealPDF
```

## Code Style

- Follow standard Go conventions
- Use `go fmt` before committing
- Write descriptive variable names
- Add comments for exported functions and types
- Keep functions focused and small

## Commit Messages

Use conventional commits format:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Adding or updating tests
- `refactor:` - Code refactoring

## Pull Request Process

1. Open an issue first to discuss significant changes
2. Update README.md if you're adding features
3. Add tests for new functionality
4. Ensure all tests pass
5. Update CLAUDE.md if architecture changes

## Questions?

Open an issue for questions or suggestions!
