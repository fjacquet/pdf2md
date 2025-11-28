package extractor

import (
	"github.com/fjacquet/pdf2md/internal/types"
)

// TextBlock is an alias for types.TextBlock to maintain compatibility
type TextBlock = types.TextBlock

// PDFExtractor defines the interface for extracting content from PDFs
// This abstraction allows swapping implementations (ledongthuc/pdf, unidoc, etc.)
type PDFExtractor interface {
	// Open initializes the PDF from a reader
	// Open(r io.ReadSeeker) error

	// GetPageCount returns the total number of pages
	GetPageCount() int

	// ExtractTextBlocks extracts text with positioning from a page (1-indexed)
	ExtractTextBlocks(page int) ([]types.TextBlock, error)

	// SetDebug enables/disables debug logging
	SetDebug(debug bool)

	// Close releases resources
	Close() error
}
