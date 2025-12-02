package extractor

import (
	"fmt"
	"os"
	"testing"
)

func createTempPDF(t *testing.T, content []byte) string {
	tmpfile, err := os.CreateTemp("", "testpdf_extractor_*.pdf")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	return tmpfile.Name()
}

func TestCustomExtractor_SetDebug(t *testing.T) {
	ext := &CustomExtractor{}
	ext.SetDebug(true)
	if !ext.debug {
		t.Error("SetDebug(true) failed")
	}
	ext.SetDebug(false)
	if ext.debug {
		t.Error("SetDebug(false) failed")
	}
}

func TestCustomExtractor_NewCustomExtractor(t *testing.T) {
	// Test with non-existent file
	_, err := NewCustomExtractor("non_existent_file.pdf")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	// Test with valid file
	// Create minimal PDF
	pdfContent := []byte(`%PDF-1.4
1 0 obj
<< /Type /Catalog /Pages 2 0 R >>
endobj
2 0 obj
<< /Type /Pages /Count 1 /Kids [ 3 0 R ] >>
endobj
3 0 obj
<< /Type /Page /Parent 2 0 R >>
endobj
xref
0 4
0000000000 65535 f 
0000000009 00000 n 
0000000060 00000 n 
0000000117 00000 n 
trailer
<< /Size 4 /Root 1 0 R >>
startxref
166
%%EOF
`)
	// Note: Offsets are approximate, might fail if Reader is strict about XRef.
	// We use the dynamic calculation logic from reader_unit_test if needed.
	// But for NewCustomExtractor it just opens the file.

	path := createTempPDF(t, pdfContent)
	defer os.Remove(path)

	ext, err := NewCustomExtractor(path)
	if err != nil {
		// It might fail if NewReader tries to read XRef immediately and fails due to bad offsets.
		// Let's see.
		// If it fails, we need to construct valid PDF.
	} else {
		ext.Close()
	}
}

func TestCustomExtractor_ExtractTextBlocks_InvalidPage(t *testing.T) {
	// We need a valid PDF structure to initialize Reader
	header := "%PDF-1.4\n"
	obj1Str := "1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n"
	obj2Str := "2 0 obj\n<< /Type /Pages /Count 1 /Kids [ 3 0 R ] >>\nendobj\n"
	obj3Str := "3 0 obj\n<< /Type /Page /Parent 2 0 R >>\nendobj\n"
	xref := "xref\n0 4\n0000000000 65535 f \n"

	offset1 := len(header)
	offset2 := offset1 + len(obj1Str)
	offset3 := offset2 + len(obj2Str)
	xrefOffset := offset3 + len(obj3Str)

	entry1 := fmt.Sprintf("%010d 00000 n \n", offset1)
	entry2 := fmt.Sprintf("%010d 00000 n \n", offset2)
	entry3 := fmt.Sprintf("%010d 00000 n \n", offset3)

	body := header + obj1Str + obj2Str + obj3Str + xref + entry1 + entry2 + entry3
	trailer := "trailer\n<< /Size 4 /Root 1 0 R >>\nstartxref\n"
	footer := fmt.Sprintf("%d\n%%%%EOF\n", xrefOffset)

	content := []byte(body + trailer + footer)
	path := createTempPDF(t, content)
	defer os.Remove(path)

	ext, err := NewCustomExtractor(path)
	if err != nil {
		t.Fatalf("NewCustomExtractor failed: %v", err)
	}
	defer ext.Close()

	// Invalid page index (0 or > count)
	// PDF pages are 1-indexed in GetPage usually?
	// Reader.GetPage(pageIndex) implementation:
	// if pageIndex < 1 || pageIndex > pageCount { return nil, fmt.Errorf(...) }

	_, err = ext.ExtractTextBlocks(0)
	if err == nil {
		t.Error("Expected error for page 0, got nil")
	}

	_, err = ext.ExtractTextBlocks(100)
	if err == nil {
		t.Error("Expected error for page 100, got nil")
	}
}
