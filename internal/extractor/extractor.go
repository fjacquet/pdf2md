package extractor

import (
	"github.com/fjacquet/pdf2md/internal/types"
)

// TextBlock is an alias for types.TextBlock to maintain compatibility
type TextBlock = types.TextBlock

// Image is an alias for types.Image
type Image = types.Image

// PDFExtractor defines the interface for extracting content from PDFs
// This abstraction allows swapping implementations (ledongthuc/pdf, unidoc, etc.)
type PDFExtractor interface {
	GetPageCount() int
	ExtractTextBlocks(pageIndex int) ([]TextBlock, []types.Image, []types.VectorGraphic, error)
	Close() error
}
