package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// run is one line of text placed at (x, y) on a US Letter page.
type run struct {
	font string // Helvetica, Helvetica-Bold or Courier
	size float64
	x, y float64
	text string
}

var pdfFonts = []string{"Helvetica", "Helvetica-Bold", "Courier"}

func escapePDF(s string) string {
	// "•" is 0x95 in WinAnsiEncoding, which the fonts below declare.
	return strings.NewReplacer(`\`, `\\`, "(", `\(`, ")", `\)`, "•", `\225`).Replace(s)
}

// writePDF writes a one-page PDF made of the given runs and returns its path.
func writePDF(t *testing.T, runs []run) string {
	t.Helper()

	var content strings.Builder
	for _, r := range runs {
		fontKey := slices.Index(pdfFonts, r.font) + 1
		fmt.Fprintf(&content, "BT /F%d %g Tf %g %g Td (%s) Tj ET\n", fontKey, r.size, r.x, r.y, escapePDF(r.text))
	}
	stream := content.String()

	objs := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R " +
			"/Resources << /Font << /F1 5 0 R /F2 6 0 R /F3 7 0 R >> >> >>",
		fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(stream), stream),
	}
	for _, f := range pdfFonts {
		objs = append(objs, fmt.Sprintf("<< /Type /Font /Subtype /Type1 /BaseFont /%s /Encoding /WinAnsiEncoding >>", f))
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
