# Headings: right text, right level, nothing else promoted

## Intent

The user: headings "must be correctly present". After the text extraction
rework (v0.4.1) the READoc average is 0.8365 against `pdftotext`'s 0.8609.
arXiv is level (0.7797 vs 0.7785); the gap is GitHub (0.8932 vs 0.9434), and
reading the outputs shows it is Markdown structure, not text. Headings are
the part of that structure the user wants fixed: every heading of the PDF
comes out as a Markdown heading at its level in the document's hierarchy
(document title `#`, sections `##`, subsections `###`), and nothing that is
not a heading is promoted to one.

## Evidence (READoc, v0.4.1, measured 2026-10-28)

Headings of the ground truth matched to the output by normalized text
(fuzzy ratio ≥ 85):

| Ground-truth heading becomes | GitHub | arXiv |
|---|---|---|
| heading, same level | 18 | 1 |
| heading, other level | 107 | 88 |
| plain text or bold | 14 | 25 |
| not found (text absent or much changed) | 29 | 18 |
| output headings matching no ground-truth heading | 12 | 137 |

Causes, read in `internal/layout`:

1. **Levels are absolute.** `calculateHeaderLevel` maps the size ratio to a
   level with fixed cut-offs (2.0 → H1, 1.75 → H2, 1.5 → H3, 1.3 → H4,
   1.15 → H5). Bold short text is always H4 (`analyzer_classify.go`, rule
   2.5). `NormalizeHeaderLevels` (`main.go:300`) only shifts the smallest
   level found to H1. Result: sections at H1 next to the title (108110:
   `## Drivers` → `# Drivers`), or a whole document starting at `####`
   (24053, 670351997, 727813).
2. **Numbered list items become headings.** The "Header (Numbered)" rule runs
   before "List Item" (`analyzer.go:170-185`), and `isNumberedHeader` accepts
   `1.` followed by text: "1. Fork it", "2. Create your feature branch" come
   out as `#` (24053, 670351997, 727813). A second path does the same in
   `analyzer_merge.go:301`.
3. **Bold is taken for a heading, and a heading for bold.** arXiv abstracts
   set in bold come out one `####` per line (0709.4466: 10 lines). The other
   way, 524987682's title ".NET MAUI Learning Resources" is rejected as a
   heading because it starts with "." (rule: first character must be upper
   case or a digit) and stays `**…**`.

## Goals

- Heading levels follow the document's own hierarchy. The rank of a
  heading style within the document gives the level: the largest style is
  `#`, the next `##`, and so on, capped at `######`.
- A numbered line is a heading when it is styled as one (larger than body
  text or bold) or when it stands alone between paragraphs; adjacent
  numbered lines in body style are a list. *(Revised during Task 2:
  1004.3799 sets its section titles "1. Sums of Like Powers" in body style,
  so style alone would lose them; adjacency is what separates a list.)*
- Bold text is a heading only when it is short and stands alone. Consecutive
  bold lines that read as one block (an abstract, an emphasised paragraph)
  stay a paragraph.
- A heading may start with any character that is not a lower-case letter
  (".NET", "(a) …", "§ 3").
- The benchmark reports a heading score so that every PR shows its effect on
  headings, next to the Levenshtein similarity and the word F1.

## Non-goals

- Bullet characters (`-` vs `*`), quote characters (`’` vs `'`), emphasis
  markers. The Markdown is valid as is; changing these would only copy the
  ground truth's style.
- Link markup. The user chose to keep `[text](url)` (option 1).
- Tables, code blocks, admonitions, ONNX-driven classification.
- Tagged PDFs (no READoc document has a structure tree).

## Anti-goals

- Thresholds tuned until particular READoc documents pass. Every rule is
  stated from typography (size, weight, line count, numbering) and carries a
  one-line comment saying why.
- A heading fix that drops text or reorders it.
- A drop of the READoc average or of either subset's word F1 in any PR.

## Success

- Primary: the overall target stays the READoc average ≥ `pdftotext`'s 0.861.
- Headings: the level-aware heading F1 rises on both subsets, and no PR
  lowers it. *(Author's call: no fixed heading-F1 target, since none was
  stated; the number is reported and gated, not aimed at.)*

## Constraints and contract

- Pure Go, no new modules; CLI flags frozen.
- Element types and `Element.Level` stay as they are; the Markdown builder is
  untouched (`#` × level).
- Levels must work page by page as today, then across the document:
  `convertPDFToMarkdown` already collects all elements before
  `NormalizeHeaderLevels`, which is where the document-wide ranking belongs.
- One PR per task, benchmark gate and golden files as in the text extraction
  plan.

## Decisions made by the author (not discussed with the user)

- Heading score: F1 over the multiset of (level, normalized text) pairs of
  ground truth and output, plus the same F1 ignoring the level. Exact match
  on normalized text (lower case, punctuation and link markup removed) keeps
  it simple and deterministic.
- Heading style key: font size rounded to 0.5 pt, plus bold. Same size and
  bold ranks above same size not bold.
- Numbering depth ("2.3.1") never overrides the style rank. *(Revised after
  the v0.5.0 review, with the user's approval: within one style, a deeper
  number ranks one level lower, so "1", "1.1", "1.1.1" set alike stay
  #, ##, ###; unnumbered headings of that style sit with its shallowest
  number.)*
- Isolation for bold headings: the element is one line and the next element
  is not bold text of the same size continuing the block.
- "Styled as a heading" means font size above `body × HeaderSizeRatio`
  (the existing 1.2 config value) or bold.
- "No PR lowers it" is per subset average, with no tolerance: the heading
  score is deterministic.

## Left to the builder

- Where the style key is carried (on `Element`, or recomputed from
  `FontSize` and the bold markers at ranking time).
- How the page-level classification and the document-level ranking split
  the work, as long as the final level comes from the document-wide rank.
- Test fixtures: synthetic layouts through the existing golden-file builder
  (`cmd/pdf2md/testpdf_test.go`).
