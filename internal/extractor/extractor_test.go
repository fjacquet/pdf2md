package extractor

import (
	"os"
	"testing"
)

func TestCustomExtractor(_ *testing.T) {
	// Create a simple test to ensure the extractor compiles
	// We need a dummy file or just check the type
	var _ PDFExtractor = (*CustomExtractor)(nil)
}

// TestWithRealPDF is a manual test for use with actual PDF files
// Run with: go test -v -run TestWithRealPDF
// Place a test PDF at internal/testdata/test.pdf
func TestWithRealPDF(t *testing.T) {
	pdfPath := "../../testdata/test.pdf"

	// Skip if test file doesn't exist
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		t.Skip("Skipping test: test.pdf not found")
	}

	ext, err := NewCustomExtractor(pdfPath)
	if err != nil {
		t.Fatalf("Failed to open PDF: %v", err)
	}
	defer func() { _ = ext.Close() }()

	pageCount, err := ext.GetPageCount()
	if err != nil {
		t.Fatalf("Failed to get page count: %v", err)
	}
	t.Logf("PDF has %d pages", pageCount)

	if pageCount == 0 {
		t.Fatal("Expected at least 1 page")
	}

	// Test first page	// Extract
	content, err := ext.ExtractTextBlocks(1)
	if err != nil {
		t.Fatalf("ExtractTextBlocks failed: %v", err)
	}

	t.Logf("Extracted %d text blocks and %d images from page 1", len(content.TextBlocks), len(content.Images))

	// Print first few blocks for inspection
	for i, block := range content.TextBlocks {
		if i >= 5 {
			break
		}
		t.Logf("Block %d: '%s' at (%.2f, %.2f) size=%.2f",
			i, block.Text, block.X, block.Y, block.FontSize)
	}
}
