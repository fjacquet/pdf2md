package pdf

import (
	"fmt"
	"os"
	"testing"
)

func TestReader_IncrementalUpdate(t *testing.T) {
	// Simulate a PDF with an incremental update.
	// 1. Initial PDF with Object 1 (Root) and Object 2 (Pages).
	// 2. Incremental update adding Object 3 (Page) and updating Root to point to it?
	// Actually, simpler:
	// Initial PDF defines Object 1 (Root).
	// Incremental update defines Object 2 (Something else) but points back to Object 1 as Root.
	// If we only read the last XRef, we might miss Object 1 if it's not redefined.
	// Standard XRef behavior: The last XRef table contains entries for *new* or *updated* objects.
	// It has a /Prev key pointing to the previous XRef offset.
	// To find Object 1, we must follow /Prev.

	header := "%PDF-1.4\n"

	// --- Revision 1 ---
	// Object 1: Root
	obj1 := "1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n"
	// Object 2: Pages
	obj2 := "2 0 obj\n<< /Type /Pages /Count 0 /Kids [] >>\nendobj\n"

	xref1 := "xref\n0 3\n0000000000 65535 f \n"
	offset1 := len(header)
	offset2 := offset1 + len(obj1)
	xrefOffset1 := offset2 + len(obj2)

	entry1 := fmt.Sprintf("%010d 00000 n \n", offset1)
	entry2 := fmt.Sprintf("%010d 00000 n \n", offset2)

	trailer1 := fmt.Sprintf("trailer\n<< /Size 3 /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", xrefOffset1)

	rev1 := header + obj1 + obj2 + xref1 + entry1 + entry2 + trailer1

	// --- Revision 2 (Incremental Update) ---
	// We add Object 3 (Info dictionary, irrelevant but valid object)
	// We do NOT redefine Object 1 or 2.
	// The new XRef table only has entry for Object 3.
	// The Trailer has /Prev pointing to xrefOffset1.

	obj3 := "3 0 obj\n<< /Producer (Test) >>\nendobj\n"
	offset3 := len(rev1)
	xrefOffset2 := offset3 + len(obj3)

	xref2 := "xref\n0 1\n0000000000 65535 f \n3 1\n" // Section for obj 3
	entry3 := fmt.Sprintf("%010d 00000 n \n", offset3)

	// Trailer must contain /Prev
	trailer2 := fmt.Sprintf("trailer\n<< /Size 4 /Root 1 0 R /Prev %d >>\nstartxref\n%d\n%%%%EOF\n", xrefOffset1, xrefOffset2)

	fileContent := []byte(rev1 + obj3 + xref2 + entry3 + trailer2)

	tmpfile, err := os.CreateTemp("", "testpdf_incremental_*.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(fileContent); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Attempt to read
	r, err := NewReader(tmpfile.Name())
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}
	defer r.Close()

	// If we didn't follow /Prev, we won't find Object 1 (Root) in the XRef table,
	// because the last XRef only has Object 3.
	// So NewReader (which calls extractRoot -> ReadObject(1)) should fail or Root will be missing.

	if r.Root == nil {
		t.Fatal("Root dictionary is nil")
	}

	// Verify we can read Object 1
	obj1Read, err := r.ReadObject(1)
	if err != nil {
		t.Fatalf("Failed to read Object 1: %v", err)
	}
	if _, ok := obj1Read.(Dictionary); !ok {
		t.Errorf("Object 1 is not a dictionary")
	}
}
