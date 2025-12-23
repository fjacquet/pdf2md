package pdf

import (
	"fmt"
	"os"
	"testing"
)

func TestReader_Resolve(t *testing.T) {
	// Create a simple PDF
	header := "%PDF-1.4\n"
	// Object 1: Root Catalog (Dictionary)
	obj1Str := "1 0 obj\n<< /Type /Catalog >>\nendobj\n"
	// Object 2: Test String
	obj2Str := "2 0 obj\n(Object 2)\nendobj\n"

	xref := "xref\n0 3\n0000000000 65535 f \n"

	offset1 := len(header)
	offset2 := offset1 + len(obj1Str)
	xrefOffset := offset2 + len(obj2Str)

	entry1 := fmt.Sprintf("%010d 00000 n \n", offset1)
	entry2 := fmt.Sprintf("%010d 00000 n \n", offset2)

	body := header + obj1Str + obj2Str + xref + entry1 + entry2
	trailer := "trailer\n<< /Size 3 /Root 1 0 R >>\nstartxref\n"
	footer := fmt.Sprintf("%d\n%%%%EOF\n", xrefOffset)

	content := []byte(body + trailer + footer)
	tmpfile, err := os.CreateTemp("", "testpdf_resolve_*.pdf")
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

	// Test Resolve with direct object
	directObj := Integer(123)
	res, err := r.Resolve(directObj)
	if err != nil {
		t.Errorf("Resolve(direct) failed: %v", err)
	}
	if res != directObj {
		t.Errorf("Resolve(direct) returned different object")
	}

	// Test Resolve with indirect reference to Object 2
	ref := IndirectRef{ObjectNumber: 2, GenerationNumber: 0}
	res, err = r.Resolve(ref)
	if err != nil {
		t.Errorf("Resolve(ref) failed: %v", err)
	}
	if s, ok := res.(StringLiteral); !ok || string(s) != "Object 2" {
		t.Errorf("Resolve(ref) returned wrong object: %v", res)
	}
}
