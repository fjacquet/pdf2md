package pdf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"testing"
)

func TestReader_XRefStm(t *testing.T) {
	// Create a synthetic PDF with XRefStm
	// Object 1: Root (Catalog) - defined in XRefStm
	// Object 2: Pages - defined in standard XRef
	// Object 3: Page - defined in standard XRef
	// Object 4: XRef Stream - defined in standard XRef

	// We need to calculate offsets carefully.
	// Let's build it dynamically.

	var buf bytes.Buffer
	buf.WriteString("%PDF-1.7\n")

	// Object 1: Catalog
	obj1Offset := buf.Len()
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")

	// Object 2: Pages
	obj2Offset := buf.Len()
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")

	// Object 3: Page
	obj3Offset := buf.Len()
	buf.WriteString("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] >>\nendobj\n")

	// Object 4: XRef Stream
	// It should contain entry for Object 1.
	// W [1 2 1]
	// Index [1 1] (Start 1, Count 1)
	// Entry for Obj 1: Type 1, Offset obj1Offset, Gen 0
	obj4Offset := buf.Len()

	// Prepare stream data
	var streamData bytes.Buffer
	// Obj 1: Type 1, Offset, Gen 0
	streamData.WriteByte(1)
	streamData.WriteByte(byte(obj1Offset >> 8))
	streamData.WriteByte(byte(obj1Offset))
	streamData.WriteByte(0)

	compressedData := compress(streamData.Bytes())

	buf.WriteString(fmt.Sprintf("4 0 obj\n<< /Type /XRef /Size 5 /W [1 2 1] /Index [1 1] /Filter /FlateDecode /Length %d >>\nstream\n", len(compressedData)))
	buf.Write(compressedData)
	buf.WriteString("\nendstream\nendobj\n")

	// Standard XRef Table
	xrefOffset := buf.Len()
	buf.WriteString("xref\n")
	buf.WriteString("0 1\n0000000000 65535 f \n")
	buf.WriteString("2 3\n") // Objects 2, 3, 4
	buf.WriteString(fmt.Sprintf("%010d 00000 n \n", obj2Offset))
	buf.WriteString(fmt.Sprintf("%010d 00000 n \n", obj3Offset))
	buf.WriteString(fmt.Sprintf("%010d 00000 n \n", obj4Offset))

	// Trailer
	buf.WriteString("trailer\n")
	buf.WriteString(fmt.Sprintf("<< /Size 5 /Root 1 0 R /XRefStm %d >>\n", obj4Offset))
	buf.WriteString("startxref\n")
	buf.WriteString(fmt.Sprintf("%d\n", xrefOffset))
	buf.WriteString("%%EOF\n")

	// Write to temp file
	tmpFile, err := os.CreateTemp("", "xrefstm_test_*.pdf")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	if _, err := tmpFile.Write(buf.Bytes()); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Open with Reader
	r, err := NewReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to open reader: %v", err)
	}
	defer func() { _ = r.Close() }()

	// Verify Root
	if r.Root == nil {
		t.Fatal("Root is nil")
	}

	// Check if Object 1 (Root) was found (it was in XRefStm)
	if _, ok := r.XrefTable[1]; !ok {
		t.Errorf("Object 1 not found in XRef table")
	}

	// Check if Object 2 (Pages) was found (it was in standard XRef)
	if _, ok := r.XrefTable[2]; !ok {
		t.Errorf("Object 2 not found in XRef table")
	}
}

func compress(data []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, _ = w.Write(data)
	_ = w.Close()
	return b.Bytes()
}
