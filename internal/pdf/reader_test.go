package pdf

import (
	"os"
	"testing"
)

func TestReader(t *testing.T) {
	path := "../../testdata/test.pdf"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("Skipping test: test.pdf not found")
	}

	r, err := NewReader(path)
	if err != nil {
		t.Fatalf("Failed to open PDF: %v", err)
	}
	defer r.Close()

	if r.Root == nil {
		t.Fatal("Root dictionary is nil")
	}

	t.Logf("Root: %v", r.Root)

	// Check /Type /Catalog
	if typeName, ok := r.Root[Name("Type")].(Name); !ok || typeName != "Catalog" {
		t.Errorf("Expected /Type /Catalog, got %v", r.Root[Name("Type")])
	}

	// Check /Pages
	pagesRef, ok := r.Root[Name("Pages")].(IndirectRef)
	if !ok {
		t.Errorf("Expected /Pages to be IndirectRef, got %T", r.Root[Name("Pages")])
	} else {
		t.Logf("Pages Ref: %v", pagesRef)

		// Try to read Pages object
		pagesObj, err := r.ReadObject(pagesRef.ObjectNumber)
		if err != nil {
			t.Errorf("Failed to read Pages object: %v", err)
		} else {
			t.Logf("Pages Object: %v", pagesObj)
		}
	}

	// Check Page Count
	count, err := r.GetPageCount()
	if err != nil {
		t.Errorf("Failed to get page count: %v", err)
	}
	if count != 75 {
		t.Errorf("Expected 75 pages, got %d", count)
	}

	// Get Page 1
	page1, err := r.GetPage(1)
	if err != nil {
		t.Errorf("Failed to get page 1: %v", err)
	} else {
		if page1[Name("Type")] != Name("Page") {
			t.Errorf("Expected /Type /Page for page 1, got %v", page1[Name("Type")])
		}

		// Extract Content
		content, err := r.ExtractContent(page1)
		if err != nil {
			t.Errorf("Failed to extract content from page 1: %v", err)
		} else {
			t.Logf("Page 1 Content Length: %d", len(content))
			if len(content) == 0 {
				t.Errorf("Page 1 content is empty")
			}

			// Load Fonts
			fm := NewFontManager(r)
			resObj := page1[Name("Resources")]
			if resRef, ok := resObj.(IndirectRef); ok {
				obj, err := r.ReadObject(resRef.ObjectNumber)
				if err != nil {
					t.Errorf("Failed to read resources: %v", err)
				} else if resDict, ok := obj.(Dictionary); ok {
					if err := fm.LoadFonts(resDict); err != nil {
						t.Errorf("Failed to load fonts: %v", err)
					}
				}
			} else if resDict, ok := resObj.(Dictionary); ok {
				if err := fm.LoadFonts(resDict); err != nil {
					t.Errorf("Failed to load fonts: %v", err)
				}
			}

			// Process
			interpreter := NewInterpreter(nil, nil)
			blocks, _, _, err := interpreter.Process(content)
			if err != nil {
				t.Fatalf("Process failed: %v", err)
			} else {
				t.Logf("Extracted %d text blocks", len(blocks))
				if len(blocks) > 0 {
					t.Logf("First block: %+v", blocks[0])
				}
			}
		}
	}

	// Get Page 75
	page75, err := r.GetPage(75)
	if err != nil {
		t.Errorf("Failed to get page 75: %v", err)
	} else {
		if page75[Name("Type")] != Name("Page") {
			t.Errorf("Expected /Type /Page for page 75, got %v", page75[Name("Type")])
		}
	}
}
