# Text extraction: rework against ISO 32000

Date: 2026-10-28
Status: draft, awaiting approval (scope widened to reading order)

## Intent

pdf2md has a lot of engineering on top of text extraction (layout rules,
ONNX, OCR, tables), but the base is not reliable: on the READoc benchmark
plain `pdftotext` beats pdf2md on 19 of 20 documents. The user's position:
before inventing anything new, go back to the PDF specification and make
the existing extraction correct. A PDF does not contain glued words; we
lose information while extracting it.

Success criterion (agreed with the user, the step towards v1.0):
**pdf2md's average READoc score ≥ `pdftotext`'s average on the same
scorer.** Today: pdf2md 0.773, `pdftotext` 0.861 (poppler
`pdftotext version 26.10.0`, default mode, no `-layout`).

## Evidence (read-only audit, 2026-10-28)

Measured on the 20 READoc PDFs (READoc commit `ec2999f`). The comparison
harness re-scored the committed pdf2md outputs and matched
`benchmark_baseline.json` within 0.0002 per document.

| Tool | Average | Worst |
|---|---|---|
| pdf2md (`--no-onnx --no-ocr`) | 0.773 | 0.509 |
| `pdftotext` | 0.861 | 0.550 |
| pymupdf4llm | 0.776 | 0.530 |

The scorer (`cmd/benchmark`) is lowercase + whitespace-collapsed
Levenshtein on the raw Markdown, markup included. It rewards correct text
and order, barely rewards structure, and slightly penalises plain-text
tools. `pdftotext` still wins, so the loss is in extraction.

Split by subset, with an order-insensitive check (word-multiset F1
against ground truth, `\w+` tokens, lowercase) next to the Levenshtein
score:

| Subset | Levenshtein pdf2md / pdftotext | Word F1 pdf2md / pdftotext |
|---|---|---|
| GitHub (10) | 0.821 / 0.943 | 0.598 / 0.975 |
| arXiv (10) | 0.724 / 0.778 | 0.894 / 0.917 |

Reading of it:

- **GitHub gap = extraction.** Words are wrong (glued), not misordered.
  In scope.
- **arXiv gap = mostly reading order.** Words are nearly right (0.894 vs
  0.917), order is not. Example: `0709.4466` has word F1 0.932 vs 0.962
  but Levenshtein 0.639 vs 0.864, because the two columns are interleaved.
- **Reachability:** closing the GitHub gap fully gives an overall average
  of about 0.834. The remaining arXiv gap is mostly reading order. The
  user decided (2026-10-28) to widen the scope to reading order so the
  criterion stays 0.861.

Reading-order evidence:

- No READoc PDF is tagged (no `/StructTreeRoot`, no `/MarkInfo`), so the
  logical structure of §14.7–14.8 gives nothing here.
- Content-stream order is a strong baseline. `pdftotext -raw`, which keeps
  the stream order, scores 0.802 on arXiv: better than `pdftotext`'s
  geometric default (0.778) and pdf2md (0.724). On GitHub, `-raw` scores
  0.888 against 0.943 for the default.
- `layout/sorter.go` (recursive XY-cut, 553 lines) finds gutters from
  block widths, and those widths are wrong (deviations 2, 3, 4). A block
  that is too wide crosses the gutter and is classed as "spanning". Part
  of the order problem may therefore disappear with the positioning fix.
  The sorter also carries many fixed constants: 200 pt minimum width,
  16 pt gutter, 10 % density, 40 % column width, flow scores.

Corpus shape:

- **GitHub subset (10 docs):** Chrome/Skia output. CID Type 0C fonts,
  `Identity-H`, embedded subsets, all with `ToUnicode`. Each line is one
  `Td` + one `TJ`. Word spaces are encoded only as `TJ` offsets of about
  `-199`/`-200` (the space width of the font); there is no space glyph.
- **arXiv subset (10 docs):** TeX output. Type 1 fonts (some Type 3),
  almost never a `ToUnicode` CMap, so text depends on `/Encoding`,
  `/Differences` and glyph names. Positioning is `Td`/`Tm`/`T*` + `TJ`.
- No document uses `TD`, `Tz`, `Tc`, `Tw`, `Ts` or `"` in page content
  streams (form XObjects not counted).

## Deviations from the spec found in the current code

Each item cites the code and the documents it is known to affect.

1. **Fixed `TJ` space threshold** — `internal/pdf/ops_text.go:166`. A
   space is inserted only when an offset is `< -200`. Chrome writes
   `-199`/`-200`, so those spaces are lost
   (`Arepositoryfilledwith…` in `github/524987682`; lines using `-211`
   come out correct). The threshold ignores the font's space width and the
   font size. Affects the whole GitHub subset. Also explains the
   "missing words" seen earlier: 123 vs 272 words in `524987682` is
   glued words, not missing text.
2. **Strings inside `TJ` do not advance the text matrix** —
   `ops_text.go:155-185`. Only numeric offsets move `Tm`. The end position
   of a `TJ` is wrong, so anything shown after it on the same line
   (another `TJ`/`Tj` without a new `Td`) is placed too far left.
3. **`showText` advance is double-scaled** — `ops_text.go:271`. `Tm` is
   advanced by a width that `addTextBlock` already transformed through
   `Tm` and the CTM, and the result is multiplied by `Tm[0]` again.
4. **`Tc`, `Tw`, `Tz` never enter the glyph advance.** They are parsed
   (`handleTextState`) and ignored. §9.4.4:
   `tx = ((w0 − Tj/1000) × Tfs + Tc + Tw) × Th`, with `Tw` applied only to
   single-byte code 32.
5. **`TD` is not implemented** — falls through to the `default` branch of
   `executeOperatorWithTokenizer` (`interpreter.go:121`) and is silently
   dropped. `docs/pdf_parser_spec.md` claims it is supported. Not in this
   corpus, common elsewhere.
6. **`Tr` (render mode) is not handled**, so invisible text (mode 3,
   e.g. OCR layers) cannot be told apart. `GraphicsState` also mixes up
   `Ts`/`Tr` field names (`graphics_state.go`).
7. **No metrics for the standard 14 fonts.** `Font.GetWidth` returns 0
   when `/Widths` is absent, and `addTextBlock` falls back to
   `len(text) × Tfs × 0.25` (`ops_text.go:291`).
8. **Width lookups ignore `MissingWidth`** (`FontDescriptor`), and the
   CID default width is hard-coded to 1000 when `/DW` is absent (correct
   per spec, but `/W` parsing silently skips malformed entries).
9. **Composite fonts assume 2-byte codes.** `DecodeString` and
   `CalculateWidth` always read 2 bytes for `Type0`, ignoring the
   encoding CMap's `codespacerange`. Correct for `Identity-H`, wrong for
   other CMaps. Without `ToUnicode`, unknown codes fall back to "ASCII if
   printable" (`font.go:488`).
10. **`ToUnicode` parsing is partial** (`font.go:557`):
    - codes are read as 1 or 2 bytes only;
    - `bfrange` with a multi-character destination is skipped
      (`len(uniBase) == 1`);
    - `begincodespacerange` is ignored.
11. **Simple fonts default to `StandardEncoding`** when `/Encoding` is
    absent (`font.go:157`). The spec says the font program's built-in
    encoding applies (symbolic TeX fonts such as CMSY, CMMI). The embedded
    font program is never read.
12. **Type 3 `FontMatrix` is ignored** (`font.go:275`). Type 3 widths are
    in glyph space, not 1/1000 em, so widths and positions are wrong.
    Present in `1201.2692` and `2206.11906`. The `CharProcs` fallback maps
    names to codes in Go map order (random).
13. **Font parse errors drop the font** and print to stdout
    (`font.go:79`). The text then decodes as raw bytes.
14. **Layout re-glues spaced text** — `internal/layout/analyzer_merge.go:90`.
    Blocks on the same line merge with no separator when
    `gap < 0.3 × FontSize`. A real space is about 0.2–0.28 em, so even
    correctly separated words merge without a space.
    `Config.WordGapMultiplier` (0.2) exists but is never used.

## Goals

- Text positioning follows ISO 32000-2 §9.2–9.4: text object, text state
  (`Tc`, `Tw`, `Tz`, `TL`, `Tf`, `Ts`, `Tr`), positioning (`Td`, `TD`,
  `Tm`, `T*`) and showing (`Tj`, `TJ`, `'`, `"`), with the exact glyph
  advance of §9.4.4 combined with the CTM.
- Each shown glyph gets a correct position and advance from §9.6–9.7
  metrics: `/Widths` with `FirstChar`, `/W` with `/DW`, `MissingWidth`,
  the standard 14 AFM widths, and the Type 3 `FontMatrix`.
- Character codes are split per the font: 1 byte for simple fonts, the
  encoding CMap's `codespacerange` for composite fonts.
- Code → Unicode follows §9.10.2 order: `ToUnicode` (full CMap syntax for
  `bfchar`/`bfrange`, multi-byte codes, multi-char destinations), then
  `/Encoding` + `/Differences` + glyph names, then the built-in encoding
  of the embedded font program where it is needed.
- Word boundaries are decided once, from positions: a space is inserted
  when the gap between consecutive glyphs exceeds a fraction of the
  font's actual space width (or of the font size when the font has no
  space glyph). Explicit space glyphs are kept.
- Reading order: blocks come out in reading order on single- and
  multi-column pages. Content-stream order is the starting point. Geometry
  is used to detect and repair the cases where the stream order is not the
  reading order, not to rebuild the order from scratch.
- Reach the success criterion above.

## Non-goals

- Tables, headings, lists, Markdown structure, ONNX and OCR, except where
  they must adapt to the new output.
- Tagged-PDF structure trees (§14.8) as an order source: no document in
  the corpus has one. Worth adding later; not needed for the criterion.
- Rendering or rasterising fonts.
- New dependencies. The core stays pure Go (embedding the standard 14 AFM
  widths as Go data is fine; they are public Adobe metrics).

## Anti-goals

- Thresholds tuned to make these 20 documents pass. Every threshold must
  be derived from font metrics or the spec, and stated as such in a
  comment. The benchmark is 20 documents; overfitting it is the failure
  mode.
- Breaking `internal/layout` consumers silently.

## Constraints and contract with the rest of the code

- `internal/pdf` keeps returning `[]types.TextBlock`. Coordinates stay as
  today: PDF user space, points, Y up, `FontSize` after `Tm` × CTM.
  Consumers that rely on this: `layout/sorter.go`, `layout/tables.go`
  (`groupByY`), `layout/graphics_analysis.go`, `layout/links.go`
  (`ApplyLinksToBlocks`), the exclusion zones in `layout/analyzer.go`,
  `internal/ocr` (producer of blocks), `cmd/pdf2md`.
- Blocks stay at word or run granularity, not whole lines, so links
  (`LinkURI`) and font changes (bold via `FontName`) keep their precision.
- **Where the space decision lives:** in `internal/pdf`, from glyph
  positions. A block's text carries its own spaces. `MergeElements`
  (`analyzer_merge.go:68`) must stop re-deciding spacing with
  `0.3 × FontSize`. It uses the same metric-based rule, or blocks carry
  enough information (a flag or the space width) for it to know a space
  belongs between them. This change is in scope.
- Existing non-regression checks stay green or are updated on purpose,
  with the reason in the commit:
  - golden files in `cmd/pdf2md/testdata/golden/`;
  - `tests/TestQuality_EndToEnd`;
  - fuzz tests;
  - `benchmark_baseline.json`, refreshed with `-update-baseline` only
    when scores go up.

## Verification

- **Unit tests derived from the spec:** small hand-built PDFs (the
  `writePDF` helper pattern from `cmd/pdf2md/testpdf_test.go`), one per
  operator and per width source. They assert glyph positions and spacing
  numerically, e.g. `TJ` with `-199` at a space width of 200 yields a
  space; `Tc`/`Tw`/`Tz` move the next glyph by the §9.4.4 amount; `TD`
  sets `TL`.
- **Position check against poppler:** `pdftotext -bbox-layout` gives word
  boxes. Poppler's origin is top-left with Y down, so convert to PDF user
  space using the page's CropBox before comparing. Use it on the READoc
  PDFs to check each positioning phase on its own before the benchmark
  moves.
- **Benchmark:** `cmd/benchmark` plus the scratch comparison against
  `pdftotext` after each phase. Record per-document scores.

## Order of work (proposed)

0. Make the comparison durable: `cmd/benchmark` can score an arbitrary
   output directory, and the per-document `pdftotext` scores are
   committed with the poppler version. Without per-document targets,
   "≥ pdftotext" cannot be checked one PR at a time.
1. Text-state and glyph-advance core (deviations 2, 3, 4, 5, 6), with the
   per-glyph advance as the single source of truth for positions.
2. Width sources (7, 8, 12) and code splitting (9).
3. Word boundary from positions (1) and the layout merge rule (14).
4. Unicode mapping (10, 11, 13).

5. Reading order. Each `TextBlock` carries its index in content-stream
   order. Measure the stream order alone first, then the existing XY-cut on
   the corrected positions. Then build the stream-order-first approach
   described in the goals, and replace `SortBlocks` with it if it scores
   better.

**Checkpoint after step 3:** measure the word-F1 and Levenshtein split
again. That shows how much order work is left before step 5 starts.

## Decisions made by the author (not discussed with the user)

- The benchmark runs `--no-onnx --no-ocr`, as `cmd/benchmark` already does.
  The criterion applies to that mode.
- One PR per step of the order above. A step may lower individual
  documents when the cause is a compensating heuristic that the next step
  removes. The average must not drop, the PR says which documents moved
  and why, and the baseline is refreshed. Merge only on green CI.
- Form XObjects (`ops_xobject.go`) and any other nested content go
  through the same text path, so they get the same fixes for free.
- Invisible text (`Tr 3`) is still extracted, as `pdftotext` does (it is
  usually an OCR layer). The render mode is recorded so that layout can
  ignore it later if needed.
- Layout heuristics that compensated for wrong positions
  (`MergeDropCaps`, `MergeParagraphLines` thresholds) are touched only
  when a step makes them wrong. Each such change is listed in that PR.

## Left to the builder

- Internal data structures (per-glyph slice vs. building runs directly),
  as long as the `TextBlock` contract above holds.
- Exact fraction of the space width used for word breaks, derived from
  metrics and documented.
- Whether to parse the embedded Type 1 / CFF / TrueType program for the
  built-in encoding in this project or defer it, based on how many
  benchmark documents need it after step 4.
- Updating `docs/pdf_parser_spec.md`, which currently claims `TD` support.
