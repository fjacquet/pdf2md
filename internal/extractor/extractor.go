package extractor

import "io"

// TextBlock represents a piece of text with positioning information
type TextBlock struct {
	Text   string
	X      float64 // X coordinate (left)
	Y      float64 // Y coordinate (top)
	Width  float64
	Height float64
	FontSize float64
	FontName string
}

// PDFExtractor defines the interface for extracting content from PDFs
// This abstraction allows swapping implementations (ledongthuc/pdf, unidoc, etc.)
type PDFExtractor interface {
	// Open initializes the PDF from a reader
	Open(r io.ReadSeeker) error

	// GetPageCount returns the total number of pages
	GetPageCount() int

	// ExtractTextBlocks extracts text with positioning from a page (1-indexed)
	ExtractTextBlocks(page int) ([]TextBlock, error)

	// SetDebug enables/disables debug logging
	SetDebug(debug bool)

	// Close releases resources
	Close() error
}
