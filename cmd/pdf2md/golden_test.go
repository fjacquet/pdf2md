package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite the golden Markdown files")

// body returns n lines of 12pt body text starting at y, spaced 14pt apart.
func body(x, y float64, lines ...string) []run {
	var out []run
	for i, l := range lines {
		out = append(out, run{"Helvetica", 12, x, y - float64(i)*14, l})
	}
	return out
}

// goldenFixtures are synthetic one-page documents, one per layout feature.
// Their Markdown output is versioned in testdata/golden; regenerate it with
// `go test ./cmd/pdf2md -run TestGolden -update` and review the diff.
var goldenFixtures = map[string][]run{
	"headings": append(append(
		[]run{{"Helvetica-Bold", 24, 72, 740, "Annual Report"}},
		body(72, 700, "This report summarizes the results of the year.", "It is split into short sections.")...),
		append([]run{{"Helvetica-Bold", 16, 72, 640, "Results"}},
			body(72, 615, "Revenue grew steadily over the four quarters.")...)...),
	"lists": append(
		[]run{{"Helvetica-Bold", 20, 72, 740, "Shopping"}},
		body(72, 700, "• Apples", "• Bread", "• Cheese", "1. Wash", "2. Slice", "3. Serve")...),
	"code": append(append(
		body(72, 740, "Run the following program:"),
		run{"Courier", 10, 72, 700, "func main() {"},
		run{"Courier", 10, 72, 688, "    fmt.Println(\"hi\")"},
		run{"Courier", 10, 72, 676, "}"}),
		body(72, 640, "It prints a greeting.")...),
	"columns": append(
		body(72, 700, "Left column line one.", "Left column line two.", "Left column line three.", "Left column line four."),
		body(330, 700, "Right column line one.", "Right column line two.", "Right column line three.", "Right column line four.")...),
}

func TestGolden(t *testing.T) {
	for name, runs := range goldenFixtures {
		t.Run(name, func(t *testing.T) {
			got, err := convertPDFToMarkdown(writePDF(t, runs), "", convertOptions{})
			if err != nil {
				t.Fatalf("convertPDFToMarkdown() error = %v", err)
			}

			path := filepath.Join("testdata", "golden", name+".md")
			if *update {
				if err := os.WriteFile(path, []byte(got), 0o600); err != nil {
					t.Fatal(err)
				}
			}

			want, err := os.ReadFile(path) //nolint:gosec // G304: fixed testdata path
			if err != nil {
				t.Fatalf("missing golden file (run with -update): %v", err)
			}
			if strings.ReplaceAll(string(want), "\r\n", "\n") != got {
				t.Errorf("output differs from %s\n got:\n%s\nwant:\n%s", path, got, want)
			}
		})
	}
}
