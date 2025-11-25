package extractor

import (
	"fmt"
	"io"

	"github.com/ledongthuc/pdf"
)

// LedongthucExtractor implements PDFExtractor using ledongthuc/pdf library
type LedongthucExtractor struct {
	reader *pdf.Reader
	debug  bool
}

// SetDebug enables/disables debug logging
func (e *LedongthucExtractor) SetDebug(debug bool) {
	e.debug = debug
}

// NewLedongthucExtractor creates a new extractor using ledongthuc/pdf
func NewLedongthucExtractor() *LedongthucExtractor {
	return &LedongthucExtractor{}
}

// Open initializes the PDF reader
func (e *LedongthucExtractor) Open(r io.ReadSeeker) error {
	// Get file size
	size, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// pdf.NewReader requires io.ReaderAt, check if we have it
	readerAt, ok := r.(io.ReaderAt)
	if !ok {
		return ErrNotReaderAt
	}

	reader, err := pdf.NewReader(readerAt, size)
	if err != nil {
		return err
	}

	e.reader = reader
	return nil
}

// GetPageCount returns the number of pages in the PDF
func (e *LedongthucExtractor) GetPageCount() int {
	if e.reader == nil {
		return 0
	}
	return e.reader.NumPage()
}

// ExtractTextBlocks extracts text content with positioning information
func (e *LedongthucExtractor) ExtractTextBlocks(page int) ([]TextBlock, error) {
	if e.reader == nil {
		return nil, ErrNotInitialized
	}

	if page < 1 || page > e.reader.NumPage() {
		return nil, ErrInvalidPage
	}

	p := e.reader.Page(page)
	if p.V.IsNull() {
		return nil, ErrInvalidPage
	}

	// Extract text with positioning using GetTextByRow
	rows, err := p.GetTextByRow()
	if err != nil {
		return nil, err
	}

	if e.debug {
		// fmt.Printf("Page %d: Found %d rows\n", page, len(rows))
	}

	var blocks []TextBlock
	count := 0
	for _, row := range rows {
		for _, word := range row.Content {
			if count < 20 {
				fmt.Printf("DEBUG: Extractor Word: '%s' X=%.2f Y=%.2f W=%.2f\n", word.S, word.X, word.Y, word.W)
			}
			blocks = append(blocks, TextBlock{
				Text:     word.S,
				X:        word.X,
				Y:        word.Y,
				Width:    word.W,
				Height:   0, // ledongthuc/pdf doesn't provide height
				FontSize: word.FontSize,
				FontName: word.Font,
			})
			count++
		}
	}

	return blocks, nil
}

// Close releases resources
func (e *LedongthucExtractor) Close() error {
	e.reader = nil
	return nil
}
