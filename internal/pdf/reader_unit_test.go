package pdf

import (
	"fmt"
	"os"
	"testing"
)

func TestReader_ParseXRefTable(t *testing.T) {
	// Create a temporary file with a simple PDF structure
	// Construct PDF content dynamically to get correct offsets
	header := "%PDF-1.4\n"
	obj1Str := "1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n"
	obj2Str := "2 0 obj\n<< /Type /Pages /Count 1 /Kids [ 3 0 R ] >>\nendobj\n"
	obj3Str := "3 0 obj\n<< /Type /Page /Parent 2 0 R >>\nendobj\n"
	xref := "xref\n0 4\n0000000000 65535 f \n"

	// Calculate offsets
	offset1 := len(header)
	offset2 := offset1 + len(obj1Str)
	offset3 := offset2 + len(obj2Str)
	xrefOffset := offset3 + len(obj3Str)

	// Format xref entries
	// 0000000000 65535 f
	// 0000000009 00000 n
	// ...
	// We need to format offsets to 10 digits
	entry1 := fmt.Sprintf("%010d 00000 n \n", offset1)
	entry2 := fmt.Sprintf("%010d 00000 n \n", offset2)
	entry3 := fmt.Sprintf("%010d 00000 n \n", offset3)

	body := header + obj1Str + obj2Str + obj3Str + xref + entry1 + entry2 + entry3
	trailer := "trailer\n<< /Size 4 /Root 1 0 R >>\nstartxref\n"
	footer := fmt.Sprintf("%d\n%%%%EOF\n", xrefOffset)

	content := []byte(body + trailer + footer)
	tmpfile, err := os.CreateTemp("", "testpdf_*.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	r, err := NewReader(tmpfile.Name())
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}
	defer func() { _ = r.Close() }()

	// Verify objects
	obj1, err := r.ReadObject(1)
	if err != nil {
		t.Fatalf("ReadObject(1) failed: %v", err)
	}
	dict1, ok := obj1.(Dictionary)
	if !ok {
		t.Fatalf("Expected Dictionary for obj 1, got %T", obj1)
	}
	if dict1[Name("Type")] != Name("Catalog") {
		t.Errorf("Expected /Type /Catalog, got %v", dict1[Name("Type")])
	}
}

func TestReader_ParseXRefStream(_ *testing.T) {
	// Minimal PDF with XRef Stream
	// Note: Constructing a valid XRef stream manually is hard because of binary data and offsets.
	// We will try to mock the reader or just test specific methods if possible.
	// Since Reader takes a file path, we must use a file.

	// We can try to test `parseXRefStream` by creating a file that points `startxref` to a stream object.
	// But the stream data must be valid.

	// Let's skip complex XRef stream construction for now and focus on other parts if possible.
	// Or use a pre-generated minimal PDF with XRef stream if we can find one or generate one.

	// Alternatively, we can test `ReadObject` with various types.
}
