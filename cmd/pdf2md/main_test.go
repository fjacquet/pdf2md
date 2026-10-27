package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
)

// TestCLIFlags freezes the public flag set: renaming or removing a flag is a
// breaking change and must be a deliberate edit of this list.
func TestCLIFlags(t *testing.T) {
	want := []string{
		"debug", "exclude-bottom", "exclude-top", "help", "model-path",
		"no-ocr", "no-onnx", "ocr-det-model-path", "ocr-dict-path",
		"ocr-rec-model-path", "runtime-path", "version",
	}

	var got []string
	for _, f := range kong.Must(&CLI).Model.Flags {
		got = append(got, f.Name)
	}
	slices.Sort(got)

	if !slices.Equal(got, want) {
		t.Errorf("CLI flags changed\n got: %v\nwant: %v", got, want)
	}
}

// writeTextPDF writes a one-page PDF whose only content is the given text.
func writeTextPDF(t *testing.T, text string) string {
	t.Helper()
	stream := fmt.Sprintf("BT /F1 12 Tf 72 700 Td (%s) Tj ET", text)
	objs := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R " +
			"/Resources << /Font << /F1 5 0 R >> >> >>",
		fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(stream), stream),
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>",
	}

	var b strings.Builder
	b.WriteString("%PDF-1.4\n")
	offsets := make([]int, len(objs))
	for i, o := range objs {
		offsets[i] = b.Len()
		fmt.Fprintf(&b, "%d 0 obj\n%s\nendobj\n", i+1, o)
	}
	xref := b.Len()
	fmt.Fprintf(&b, "xref\n0 %d\n0000000000 65535 f \n", len(objs)+1)
	for _, off := range offsets {
		fmt.Fprintf(&b, "%010d 00000 n \n", off)
	}
	fmt.Fprintf(&b, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objs)+1, xref)

	path := filepath.Join(t.TempDir(), "doc.pdf")
	if err := os.WriteFile(path, []byte(b.String()), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestConvertPDFToMarkdown(t *testing.T) {
	pdf := writeTextPDF(t, "Hello pdf2md")

	md, err := convertPDFToMarkdown(pdf, "", convertOptions{})
	if err != nil {
		t.Fatalf("convertPDFToMarkdown() error = %v", err)
	}
	if !strings.Contains(md, "Hello pdf2md") {
		t.Errorf("markdown = %q, want it to contain %q", md, "Hello pdf2md")
	}
}

func TestConvertPDFToMarkdown_MissingFile(t *testing.T) {
	_, err := convertPDFToMarkdown(filepath.Join(t.TempDir(), "nope.pdf"), "", convertOptions{})
	if err == nil {
		t.Fatal("convertPDFToMarkdown() error = nil, want error for missing file")
	}
}
