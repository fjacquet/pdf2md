package extractor

import (
	"github.com/fjacquet/pdf2md/internal/pdf"
)

// CustomExtractor implements PDFExtractor using our custom parser
type CustomExtractor struct {
	reader *pdf.Reader
}

// NewCustomExtractor creates a new custom extractor
func NewCustomExtractor(path string) (*CustomExtractor, error) {
	reader, err := pdf.NewReader(path)
	if err != nil {
		return nil, err
	}
	return &CustomExtractor{reader: reader}, nil
}

// GetPageCount returns the number of pages
func (e *CustomExtractor) GetPageCount() int {
	count, err := e.reader.GetPageCount()
	if err != nil {
		return 0
	}
	return count
}

// ExtractTextBlocks extracts text blocks from a page
func (e *CustomExtractor) ExtractTextBlocks(pageIndex int) ([]TextBlock, error) {
	// Get page dictionary
	pageDict, err := e.reader.GetPage(pageIndex)
	if err != nil {
		return nil, err
	}

	// Load fonts
	fm := pdf.NewFontManager(e.reader)
	if resources, ok := pageDict[pdf.Name("Resources")]; ok {
		// Handle indirect reference for Resources
		if ref, ok := resources.(pdf.IndirectRef); ok {
			obj, err := e.reader.ReadObject(ref.ObjectNumber)
			if err == nil {
				if resDict, ok := obj.(pdf.Dictionary); ok {
					_ = fm.LoadFonts(resDict)
				}
			}
		} else if resDict, ok := resources.(pdf.Dictionary); ok {
			_ = fm.LoadFonts(resDict)
		}
	}

	// Extract content stream
	content, err := e.reader.ExtractContent(pageDict)
	if err != nil {
		return nil, err
	}
	if content == nil {
		return nil, nil // Empty page
	}

	// Interpret content
	interpreter := pdf.NewInterpreter(fm)
	blocks, err := interpreter.Process(content)
	if err != nil {
		return nil, err
	}

	return blocks, nil
}

// SetDebug enables debug logging
func (e *CustomExtractor) SetDebug(debug bool) {
	// TODO: Pass debug flag to reader/interpreter
}

// Close closes the PDF file
func (e *CustomExtractor) Close() error {
	return e.reader.Close()
}
