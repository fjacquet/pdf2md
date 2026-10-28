package extractor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writePDF writes objs as objects 1..n with a classic xref table; object 1
// must be the catalog.
func writePDF(t *testing.T, objs ...string) string {
	t.Helper()
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
	path := filepath.Join(t.TempDir(), "test.pdf")
	if err := os.WriteFile(path, []byte(b.String()), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func stream(dict, data string) string {
	return fmt.Sprintf("<< %s /Length %d >>\nstream\n%s\nendstream", dict, len(data), data)
}

// Resources and MediaBox on /Pages only: the page used to fail with
// "no resources dictionary" and report a 0×0 page.
func TestExtractTextBlocks_InheritedResources(t *testing.T) {
	content := "q 10 0 0 10 50 50 cm /Im1 Do Q BT /F1 12 Tf 72 700 Td (inherited) Tj ET"
	e, err := NewCustomExtractor(writePDF(t,
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 /MediaBox [0 0 612 792] /Resources << /Font << /F1 5 0 R >> /XObject << /Im1 6 0 R >> >> >>",
		"<< /Type /Page /Parent 2 0 R /Contents 4 0 R >>",
		stream("", content),
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding /WinAnsiEncoding >>",
		stream("/Type /XObject /Subtype /Image /Width 2 /Height 2 /ColorSpace /DeviceRGB /BitsPerComponent 8", strings.Repeat("\xff\x00\x00", 4)),
	))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = e.Close() }()
	c, err := e.ExtractTextBlocks(1)
	if err != nil {
		t.Fatalf("ExtractTextBlocks(1): %v", err)
	}
	if len(c.Images) != 1 || len(c.TextBlocks) != 1 || c.PageHeight != 792 {
		t.Errorf("got %d images, %d blocks, height %v; want 1, 1, 792", len(c.Images), len(c.TextBlocks), c.PageHeight)
	}
}
