# Text Extraction Spec Rework Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make pdf2md's text extraction and reading order follow ISO 32000 closely enough that its READoc average reaches at least `pdftotext`'s (0.861).

**Architecture:** `internal/pdf` computes every glyph's position with the §9.4.4 advance formula. It decides word spaces from those positions against the font's real space width and still emits one `types.TextBlock` per text-showing operator. `internal/layout` stops re-deciding spaces with a fixed font-size ratio and reuses the same metric. Reading order starts from content-stream order (the order blocks already arrive in), and XY-cut is used only on pages where that order is geometrically incoherent.

**Tech Stack:** Go (pure, no new modules), poppler `pdftotext` 26.10.0 as the reference, READoc corpus at `./READoc` (commit `ec2999f`).

**Spec:** `docs/superpowers/specs/2026-10-28-text-extraction-spec-rework-design.md`

## Global Constraints

- No new Go modules; the core stays pure Go.
- `TextBlock` coordinates stay PDF user space, points, Y up, `FontSize` after `Tm` × CTM.
- Blocks stay at word or run granularity (one per text-showing operator), never whole lines.
- Every threshold is derived from font metrics or the spec, with a one-line comment naming its source. No value tuned to make particular READoc documents pass.
- CLI flags and env vars are frozen (`cmd/pdf2md/main_test.go:TestCLIFlags`); add none.
- Benchmark mode is `--no-onnx --no-ocr`. Target: average ≥ `pdftotext` average on the same scorer.
- One PR per task. The average must not drop. The PR lists the documents that moved and why. Refresh `benchmark_baseline.json` with `-update-baseline`, merge only on green CI.
- Golden files (`cmd/pdf2md/testdata/golden/`) are updated with `go test ./cmd/pdf2md -run TestGolden -update` only when the new output is more correct. Say so in the commit.
- Lint: `golangci-lint run --fix ./...` first, then `make lint`. gosec suppressions use the inline `//nolint:gosec // Gxxx: reason` style.
- Commits end with the session's `Co-Authored-By` and `Claude-Session` trailers.

## Review Focus

- **Form XObject with its own `/Resources /Font`.** Today the sub-interpreter reuses the page `FontManager` and never loads the form's fonts, so form text decodes with the wrong font or none. Expected: the form's fonts override page fonts of the same name, inside the form only. Test in Task 1.
- **Rotated text** (`Tm` = `0 1 -1 0 x y`, as in arXiv side margins). Expected: text intact, `Width` ≥ 0, no spurious spaces between glyphs. Test in Task 3.
- **`Tf` naming a font absent from resources.** Expected: glyphs still advance, using a 500/1000 em default width and a 250/1000 em space, so words stay separated. Test in Task 2.
- **Negative font size or `Tz`** (mirrored text, allowed by the spec). Expected: positive `Width` and `FontSize`, text intact. Test in Task 1.
- **`TJ` with only numbers, or empty strings.** Expected: no block emitted, but `Tm` still moves. Test in Task 1.

---

### Task 0: Durable comparison against `pdftotext`

**Files:**
- Modify: `cmd/benchmark/main.go`
- Modify: `cmd/benchmark/main_test.go`
- Create: `benchmark_pdftotext.json`
- Modify: `README.md` (Benchmark section), `CLAUDE.md` (commands)

**Interfaces:**
- Produces:
  - flags `-cmd` (default `./pdf2md --no-onnx --no-ocr {in} {out}`) and `-output` (default `READoc/output/pdf2md`);
  - a `Tool string` field (`json:"tool,omitempty"`) on `baseline`;
  - `wordF1(gen, gt string) float64`, printed next to each Levenshtein score;
  - `expandCmd(tmpl, in, out string) []string`.

- [ ] **Step 1: Write failing tests**

```go
func TestExpandCmd(t *testing.T) {
	got := expandCmd("pdftotext {in} {out}", "a b.pdf", "o.md")
	want := []string{"pdftotext", "a b.pdf", "o.md"}
	// reflect.DeepEqual(got, want); placeholders are replaced per field, so paths with spaces survive
}
func TestWordF1(t *testing.T) {
	// wordF1("the cat the", "the cat") == 0.8   (multiset ∩ = 2, sizes 3 and 2)
	// wordF1("", "x") == 0 ; wordF1("A-b", "a b") == 1  (tokens are \w+, lowercased)
}
```

- [ ] **Step 2:** Run `go test ./cmd/benchmark -run 'TestExpandCmd|TestWordF1'`. Expected: FAIL, undefined.
- [ ] **Step 3:** Implement both. `expandCmd` uses `strings.Fields(tmpl)` and replaces `{in}`/`{out}` in each field. In `main`, replace the hard-coded `exec.Command("./pdf2md", …)` with `expandCmd(*cmdTmpl, pdfPath, outPath)`. When the command writes to stdout (no `{out}` in the template), write stdout to `outPath`. `writeBaseline` records `Tool`: the first line of `<cmd[0]> -v` for pdftotext, and `git describe --always --dirty` for `./pdf2md`.
- [ ] **Step 4:** Run `go test ./cmd/benchmark`. Expected: PASS.
- [ ] **Step 5:** Generate the reference:
  `go run ./cmd/benchmark -cmd 'pdftotext {in} {out}' -output READoc/output/pdftotext -baseline benchmark_pdftotext.json -update-baseline`
  Expected: `Average Similarity: 0.861x`, and the file records `pdftotext version 26.10.0`.
- [ ] **Step 6:** Run `go run ./cmd/benchmark` (pdf2md, default flags). Expected: average 0.7726, no regressions.
- [ ] **Step 7:** Commit `feat(benchmark): pluggable command, word F1, committed pdftotext reference`.

### Task 1: Glyph-by-glyph text showing (§9.4.4)

**Files:**
- Modify: `internal/pdf/ops_text.go` (rewrite `handleTextShow`, `showText`, `addTextBlock`; add `TD`, `Tr`)
- Modify: `internal/pdf/interpreter.go:87-92` (dispatch `TD`, `Tr`)
- Modify: `internal/pdf/graphics_state.go` (rename fields: `Tmode` holds the render mode set by `Tr`, `Ts` holds the rise; delete the unused `Tr float64`)
- Modify: `internal/pdf/ops_xobject.go:137` (form fonts)
- Modify: `internal/pdf/font.go` (add `Child`)
- Modify: `internal/types/types.go` (`TextBlock.RenderMode int`)
- Test: `internal/pdf/ops_text_test.go`, `internal/pdf/ops_xobject_test.go`

**Interfaces:**
- Produces:
  - `func (in *Interpreter) showString(raw []byte) (text string, width float64)`: walks the codes, appends decoded text, advances `in.Tm` per glyph, and returns the advance sum in user space.
  - `func (fm *FontManager) Child(resources Dictionary) *FontManager`: a copy of `fm.Fonts` overlaid with fonts loaded from `resources`.
  - `TextBlock.RenderMode` (0–7, the `Tr` value).
  - Until Task 2 lands, code splitting and widths keep using `Font.CalculateWidth` semantics through a temporary per-code helper `in.glyphWidth(code int) float64` (thousandths of em). Task 2 replaces it.
- Test helper (in `ops_text_test.go`):

```go
// runText interprets content with font F1: every byte has width 500 and the
// space (32) has width 250, FirstChar 0. It returns the blocks.
func runText(t *testing.T, content string) []types.TextBlock
```

- [ ] **Step 1: Write failing tests**, one per rule. Each is a `runText` call plus assertions on `X`, `Y`, `Width`, `Text`, with tolerance 1e-9:
  - `TestAdvance_TcTz`: `BT /F1 10 Tf 2 Tc 50 Tz 100 200 Td (AB) Tj (C) Tj ET` gives `blocks[0].Width == 7` and `blocks[1].X == 107`. Per glyph: (0.5×10 + 2) × 0.5 = 3.5.
  - `TestAdvance_TwOnlyOnSpace`: `BT /F1 10 Tf 3 Tw 0 0 Td (A A) Tj (B) Tj ET` gives `blocks[1].X == 15.5` (5 + 2.5 + 3 + 5).
  - `TestAdvance_TJStringsAdvance`: `BT /F1 10 Tf 0 0 Td [(A) -1000 (B)] TJ (C) Tj ET` gives `blocks[1].X == 20`.
  - `TestAdvance_NoDoubleScaling`: `BT /F1 10 Tf 2 0 0 2 0 0 Tm (A) Tj (B) Tj ET` gives `blocks[1].X == 10` and `blocks[0].FontSize == 20`.
  - `TestAdvance_CTM`: `2 0 0 2 0 0 cm BT /F1 10 Tf (A) Tj (B) Tj ET` gives `blocks[1].X == 10`.
  - `TestTD_SetsLeading`: `BT /F1 10 Tf 0 100 Td 0 -14 TD (A) Tj T* (B) Tj ET` gives `Y` = 86 and 72.
  - `TestRise`: `BT /F1 10 Tf 5 Ts (A) Tj ET` gives `blocks[0].Y == 5`.
  - `TestRenderMode`: `BT /F1 10 Tf 3 Tr (A) Tj ET` gives `RenderMode == 3` and text `"A"`, still extracted.
  - `TestTJOnlyNumbers`: `BT /F1 10 Tf [-1000] TJ (A) Tj ET` gives exactly one block, at `X == 10`.
  - `TestNegativeFontSize`: `BT /F1 -10 Tf (AB) Tj ET` gives `Width == 10`, `FontSize == 10` and text `"AB"`.
  - `TestFormFonts` (in `ops_xobject_test.go`): a page font `F1` (widths 500) and a form whose resources define `F1` with widths 1000. `(AB) Tj` inside the form gives `Width == 20`, and after `Do` the page's `(AB) Tj` gives `Width == 10`.
- [ ] **Step 2:** Run `go test ./internal/pdf -run 'TestAdvance|TestTD|TestRise|TestRenderMode|TestTJOnly|TestNegative|TestFormFonts'`. Expected: FAIL.
- [ ] **Step 3: Implement.** Per glyph, as in §9.4.4:
  - `tx = ((w0/1000)·Tfs + Tc + Tw·[single-byte code 32]) · Th/100`, then `Tm = [1 0 0 1 tx 0] × Tm`.
  - A `TJ` number `n` gives `Tm = [1 0 0 1 (−n/1000)·Tfs·Th/100 0] × Tm`.
  - Block origin: `(0, Ts)` through `Tm` then CTM, taken before the first glyph.
  - `Width`: the distance from the start to the end origin in user space, absolute value.
  - `FontSize`: `|(0,Tfs)|` through `Tm` and CTM, as today.
  - Keep the current "offset < −200 inserts a space" rule unchanged in this task. Task 3 replaces it.
  - `TD tx ty` = `-ty TL` then `tx ty Td`.
  - Delete the `0.25 × Tfs` width estimate in `addTextBlock`.
- [ ] **Step 4:** Run `go test ./internal/pdf/...`. Expected: PASS. The old `TestInterpreter_HandleTextShow` must still pass: `"A B"` comes from the −1000 rule.
- [ ] **Step 5:** Run `go test ./...`, `make lint`, `make build`, then `go run ./cmd/benchmark`. Record per-document deltas for the PR. Expected: average ≥ 0.7726. If a document drops more than 0.02, find the compensating heuristic and name it in the PR, as the design allows.
- [ ] **Step 5b: Position check against poppler.** Run `pdftotext -bbox-layout -f 1 -l 1 READoc/data/github/pdf/524987682.pdf -` and compare the first 20 `<word xMin=…>` values with the X of the matching pdf2md blocks. Read the blocks with a throwaway `go test` that prints `ExtractTextBlocks(1)`, and do not commit it. Poppler's Y is top-down, so `y_pdf = pageHeight − yMax`. Expected: block starts within 0.5 pt of poppler's word starts.
- [ ] **Step 6:** Fix `docs/pdf_parser_spec.md` §3.3: list `TD` and `Tr`, and drop "partial" for `q`/`Q`/`cm`. Commit `fix(pdf): glyph advance per ISO 32000 §9.4.4, TD, Tr, form fonts`.

### Task 2: Glyph widths and code splitting (§9.6–9.7)

**Files:**
- Create: `internal/pdf/std14.go`, generated by `internal/pdf/gen_std14/main.go` (`//go:build ignore`). It reads the Adobe Core14 AFM files (`C <code> ; WX <w> ; N <name>`) and emits `var std14Widths = map[string]map[string]int{ "Helvetica": {"space": 278, …}, … }` plus the standard aliases `Arial`→`Helvetica`, `TimesNewRoman`→`Times-Roman`, `CourierNew`→`Courier`. Source directory: `$(go env GOMODCACHE)/github.com/pdfcpu/pdfcpu@v0.16.1/internal/corefont/Core14_AFMs/`. Put it in the `go:generate` comment. The AFM files themselves are not committed.
- Modify: `internal/pdf/font.go` (width lookup, `MissingWidth`, Type 3 `FontMatrix`, codespace)
- Test: `internal/pdf/font_test.go`

**Interfaces:**
- Consumes: `in.glyphWidth` call sites from Task 1.
- Produces:
  - `type codespace struct{ lo, hi []byte }` and `func parseCodespace(cmap []byte) []codespace`, which reads `begincodespacerange … endcodespacerange`. Task 4 reuses both.
  - `func (f *Font) Codes(raw []byte) []Code` with `type Code struct{ Value, Len int }`. Simple fonts use 1 byte. Composite fonts use the codespace of the embedded `/Encoding` CMap stream, else 2 bytes for `Identity-H`/`Identity-V`. Other predefined CMaps also get 2 bytes, with a `shortcut:` comment.
  - `func (f *Font) Width(c Code) float64`, in thousandths of em of text space. Order:
    1. `/W` (composite) or `/Widths` at `Value−FirstChar` (simple);
    2. `FontDescriptor /MissingWidth`;
    3. std14 by `BaseFont` (subset prefix `ABCDEF+` stripped) via the glyph name from `Encoding`;
    4. `/DW` (composite, default 1000);
    5. 500.
    For Type 3, `/Widths` values are multiplied by `FontMatrix[0]×1000`.
  - `func (f *Font) SpaceWidth() float64`: width of the code whose Unicode is `" "` (ToUnicode or encoding), else std14 `space`, else 250.
  - A nil font (unknown `Tf` name) behaves as width 500 and space 250.
  - `in.glyphWidth` is deleted. `Font.CalculateWidth` and `FontManager.CalculateWidth` are deleted if they are no longer called.
- [ ] **Step 1: Write failing tests**:
  - `TestWidth_MissingWidth`: Widths 32..33 = [250, 300], `MissingWidth` 400, so `Width({65,1}) == 400`.
  - `TestWidth_Type3FontMatrix`: `FontMatrix [0.01 0 0 0.01 0 0]`, Widths A = 50, so `Width == 500`.
  - `TestWidth_Std14`: `BaseFont /Helvetica`, no Widths, WinAnsi, so `Width({'A',1}) == 667` and `SpaceWidth() == 278`.
  - `TestWidth_SubsetPrefix`: `ABCDEF+Helvetica` gives the same widths.
  - `TestWidth_UnknownFont`: `runText` with `/F9 10 Tf (AB) Tj` (F9 absent) gives `Width == 10`.
  - `TestCodes_Codespace`: codespace `<00> <80>` and `<8140> <FFFF>`, bytes `41 81 40 42` give codes `0x41/1, 0x8140/2, 0x42/1`.
  - `TestCodes_IdentityH`: `00 41 00 42` gives `0x41/2, 0x42/2`.
- [ ] **Step 2:** Run `go test ./internal/pdf -run 'TestWidth|TestCodes'`. Expected: FAIL.
- [ ] **Step 3:** Implement. Run `go generate ./internal/pdf`, which produces `std14.go`. `parseFont` reads `FontDescriptor` (direct or indirect) for `MissingWidth`, and Type 3 `FontMatrix`. `DecodeString` iterates `Codes()` instead of its own 2-byte guess.
- [ ] **Step 4:** Run `go test ./internal/pdf/...`. Expected: PASS.
- [ ] **Step 5:** Run `go test ./...`, `make lint`, `make build`, then `go run ./cmd/benchmark`. Expected: average ≥ Task 1's.
- [ ] **Step 6:** Commit `feat(pdf): spec-ordered glyph widths, std14 metrics, codespace-based code splitting`.

### Task 3: Word boundaries from positions

**Files:**
- Modify: `internal/pdf/ops_text.go` (space decision inside a block)
- Modify: `internal/types/types.go` (`TextBlock.SpaceWidth float64`), `internal/types/element.go` (`Element.SpaceWidth float64`)
- Modify: `internal/layout/analyzer.go` (`blocksToElements` copies `SpaceWidth`)
- Modify: `internal/layout/analyzer_merge.go:68-140` (`MergeElements` gap rule)
- Test: `internal/pdf/ops_text_test.go`, `internal/layout/analyzer_merge_test.go`

**Interfaces:**
- Consumes: `Font.SpaceWidth()`, `Font.Width()` from Task 2.
- Produces:
  - `TextBlock.SpaceWidth`: space width in user space for the block's font and size.
  - `func isWordGap(gap, spaceWidth float64) bool`, true when `gap > 0.5·spaceWidth`. It lives in `internal/types` so both packages share it. Comment: "half the font's own space width: kerning stays below it, an inter-word space (even a justified, shrunk one) stays above".
- [ ] **Step 1: Write failing tests**:
  - `TestSpace_TJOffsetChrome`: `[(A) -199 (B)] TJ` gives `"A B"`. The threshold is 125 thousandths, since the space width is 250.
  - `TestSpace_TJKerning`: `[(A) -100 (B)] TJ` gives `"AB"`, and `[(A) 50 (B)] TJ` also gives `"AB"`.
  - `TestSpace_NoDoubleSpace`: `[(A ) -1000 (B)] TJ` gives `"A B"`.
  - `TestSpace_Rotated`: `BT /F1 10 Tf 0 1 -1 0 100 100 Tm [(A) -199 (B)] TJ ET` gives `"A B"` and `Width ≥ 0`. The gap is measured along the text direction in text space, not in X.
  - `TestMergeElements_SpaceWidth` (layout): elements `A` at X 0, Width 5 and `B` at X 8, `SpaceWidth` 2.5, same Y, merge to `"A B"`. With `B` at X 6 they merge to `"AB"`. With `SpaceWidth == 0` (OCR blocks), the fallback is `0.25·FontSize`, with a comment.
- [ ] **Step 2:** Run `go test ./internal/pdf ./internal/layout -run 'TestSpace|TestMergeElements_SpaceWidth'`. Expected: FAIL.
- [ ] **Step 3:** Implement.
  - In `showString`, measure the gap in text space before each glyph: the expected pen position minus the actual one, after `TJ` offsets.
  - Insert `" "` when `isWordGap` is true and neither side is already a space.
  - Delete the `n < -200` rule.
  - In `MergeElements`, replace `threshold := current.FontSize * 0.3` with `isWordGap(gap, sw)`, where `sw` is `current.SpaceWidth` or the fallback. A word gap merges with `" "`; anything else merges with `""`.
  - Keep the hyphen and wide-gap branches unchanged.
- [ ] **Step 4:** Run `go test ./internal/pdf/... ./internal/layout/...`. Expected: PASS. Update `TestInterpreter_HandleTextShow` only if it relied on the old −200 rule. With −1000 it still yields `"A B"`.
- [ ] **Step 5:** Run `go test ./...`, `make lint`, `make build`, then `go run ./cmd/benchmark`. Expected: the GitHub subset's word F1 goes from 0.598 to ≥ 0.95, and `github/524987682` comes out spaced. Refresh the golden files if they got more correct.
- [ ] **Step 6: Checkpoint.** Post the per-subset Levenshtein and word-F1 table in the PR, compared with `benchmark_pdftotext.json`.
- [ ] **Step 7:** Commit `fix: decide word spaces from glyph positions and the font's space width`.

### Task 4: Code → Unicode (§9.10)

**Files:**
- Modify: `internal/pdf/font.go` (`parseToUnicodeCMap`, encoding fallback, Type 3 `CharProcs` heuristic, error path in `LoadFonts`)
- Create: `internal/pdf/type1_encoding.go`
- Test: `internal/pdf/font_test.go`

**Interfaces:**
- Produces:
  - `func parseToUnicodeCMap(data []byte) (map[int]string, []codespace)`, a pure function that calls `parseCodespace` from Task 2. When the font has no encoding-CMap codespace, the ToUnicode one is used by `Font.Codes`.
  - `func type1BuiltinEncoding(fontFile []byte) map[int]string`, which reads `dup <code> /<name> put` from the cleartext part of a `FontFile` stream.
- [ ] **Step 1: Write failing tests**:
  - `TestToUnicode_BfrangeIncrement`: `<0001> <0003> <0041>` maps codes 1, 2, 3 to `A`, `B`, `C`.
  - `TestToUnicode_BfrangeMultiChar`: `<0005> <0006> <00660069>` maps 5 to `"fi"` and 6 to `"fj"`. The last UTF-16 unit is incremented.
  - `TestToUnicode_BfrangeArray`: `<0007> <0008> [<0058> <00590059>]` maps 7 to `"X"` and 8 to `"YY"`.
  - `TestToUnicode_FourByteCode`: `<00000041> <0042>` with a 4-byte codespace maps 0x41 to `"B"`.
  - `TestToUnicode_Codespace`: parsing returns the ranges from `begincodespacerange`.
  - `TestType1BuiltinEncoding`: the cleartext `/Encoding 256 array\n0 1 255 {1 index exch /.notdef put} for\ndup 65 /alpha put\nreadonly def` maps 65 to `"alpha"`. A `parseFont` with that `FontFile` and no `/Encoding` decodes byte 65 to `"α"`.
  - `TestLoadFonts_BadFontKept`: a font dict whose `ToUnicode` stream fails to decode is still registered, with widths and encoding. A warning goes to `log.Printf`, not stdout.
  - `TestType3_NoCharProcsGuess`: a Type 3 font without `/Encoding` gets no code→name map from Go map order. Its codes fall back to the bytes as today.
- [ ] **Step 2:** Run `go test ./internal/pdf -run 'TestToUnicode|TestType1|TestLoadFonts_Bad|TestType3_No'`. Expected: FAIL.
- [ ] **Step 3:** Implement.
  - Lookup order per §9.10.2: `ToUnicode`, then the encoding's glyph name through `GlyphToUnicode`, then the built-in Type 1 encoding when there is no `/Encoding`, then the byte.
  - `StandardEncoding` stays the default only for non-symbolic fonts: `FontDescriptor /Flags` bit 3 (symbolic) is clear.
  - For CFF (`FontFile3`) built-in encodings, add `shortcut: CFF built-in encoding not read, upgrade when a corpus document needs it`.
- [ ] **Step 4:** Run `go test ./internal/pdf/...`. Expected: PASS.
- [ ] **Step 5:** Run `go test ./...`, `make lint`, `make build`, then `go run ./cmd/benchmark`. Expected: the arXiv word F1 does not drop and moves towards `pdftotext`'s 0.917.
- [ ] **Step 6:** Commit `fix(pdf): complete ToUnicode CMap parsing and built-in Type 1 encodings`.

### Task 5: Reading order, stream first

**Files:**
- Create: `internal/layout/stream_order.go`
- Modify: `internal/layout/analyzer.go:223` (`SortBlocks` call becomes `OrderBlocks`)
- Test: `internal/layout/stream_order_test.go`

**Interfaces:**
- Consumes: the corrected `X`, `Y`, `Width`, `FontSize` from Tasks 1–3. Blocks reach `Analyze` in content-stream order, because the interpreter appends them as it runs and forms are inlined at `Do`. No sequence field is needed; this departs from the design's "index" wording because the slice order already carries it.
- Produces:
  - `func OrderBlocks(blocks []extractor.TextBlock) []extractor.TextBlock`: returns `blocks` unchanged when `streamOrderCoherent(blocks)`, else `SortBlocks(blocks)`.
  - `func streamOrderCoherent(blocks []extractor.TextBlock) bool`:
    - Group consecutive blocks into lines (|ΔY| < 0.5·FontSize).
    - Count backward jumps: the next line starts more than one `FontSize` above the previous line, and their X ranges overlap (same column).
    - Coherent when backward jumps ≤ 10 % of line transitions.
    - Comment: "a column-ordered stream only moves up when it changes column".
- [ ] **Step 1: Write failing tests**:
  - `TestStreamOrder_TwoColumnsInStreamOrder`: col 1 lines (X 50, Y 500, then 480), then col 2 lines (X 300, Y 500, then 480). The result equals the input.
  - `TestStreamOrder_IncoherentFallsBack`: lines written bottom-to-top in a single column (Y 100, 200, 300, X 50) give `SortBlocks` order, Y 300, 200, 100.
  - `TestStreamOrder_HeaderWrittenLast`: a one-line header at the top, written after 20 body lines, gives 1 backward jump in 20 transitions. That is ≤ 10 %, so the stream order is kept. This documents the accepted limit.
- [ ] **Step 2:** Run `go test ./internal/layout -run TestStreamOrder`. Expected: FAIL.
- [ ] **Step 3:** Implement and switch `Analyze` to `OrderBlocks`.
- [ ] **Step 4:** Run `go test ./internal/layout/...`. Expected: PASS. Existing `TestSortBlocks` is unchanged, because `SortBlocks` stays.
- [ ] **Step 5:** Measure three variants on the same commit by temporarily editing the call in `analyzer.go`, without committing the variants: (a) `SortBlocks`, (b) stream order only, (c) `OrderBlocks`. Keep the variant with the best average. On a tie within 0.005, keep the simpler one, in the order b, then c, then a. Put the three averages in the PR.
- [ ] **Step 6:** Run `go test ./...`, `make lint`, `make build`, then `go run ./cmd/benchmark` and compare with `benchmark_pdftotext.json`. Expected: average ≥ 0.861. If not, stop and report the per-document gap to the user. Do not add thresholds.
- [ ] **Step 7:** Refresh the golden files if they are more correct (`columns.md` should stop being a table). Update `CHANGELOG.md` (Unreleased) for Tasks 1–5, the README roadmap and `CLAUDE.md` (Layout: reading order, Technical Notes: word gap rule). Commit `feat(layout): stream-order-first reading order with XY-cut fallback`.
