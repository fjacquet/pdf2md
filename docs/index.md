# pdf2md

A fast, single-binary PDF to Markdown converter written in Go. Designed for AI
engineers, researchers, and knowledge workers who need clean, structured Markdown
from PDF documents.

## Highlights

- **Single binary** — zero runtime dependencies for the core converter.
- **Column detection** — handles multi-column layouts (academic papers, newspapers).
- **Smart layout analysis** — preserves reading order and document structure.
- **RAG-ready output** — clean Markdown suited to retrieval pipelines and Obsidian.

## Documentation

- [Architecture](architecture.md) — internal design and data structures.
- [PDF Parser Spec](pdf_parser_spec.md) — the custom pure-Go PDF parser.
- [OCR Setup](ocr-setup.md) — optional ONNX Runtime / OCR configuration.
- [Stack Recommendation](stack_recommendation.md) — recommended libraries and CLI design.

## Quick start

```bash
pdf2md input.pdf output.md
```

See the [project README](https://github.com/fjacquet/pdf2md) for installation
instructions and the full feature list.
