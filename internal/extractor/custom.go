package extractor

import (
	"github.com/fjacquet/pdf2md/internal/pdf"
	"github.com/fjacquet/pdf2md/internal/types"
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
func (e *CustomExtractor) ExtractTextBlocks(pageIndex int) ([]TextBlock, []types.Image, []types.VectorGraphic, error) {
	// Get page dictionary
	pageDict, err := e.reader.GetPage(pageIndex)
	if err != nil {
		return nil, nil, nil, err
	}

	// Get resources dictionary
	var resources pdf.Dictionary
	if res, ok := pageDict[pdf.Name("Resources")]; ok {
		if ref, ok := res.(pdf.IndirectRef); ok {
			obj, err := e.reader.ReadObject(ref.ObjectNumber)
			if err == nil {
				if resDict, ok := obj.(pdf.Dictionary); ok {
					resources = resDict
				}
			}
		} else if resDict, ok := res.(pdf.Dictionary); ok {
			resources = resDict
		}
	}

	// Load fonts
	fm := pdf.NewFontManager(e.reader)
	if resources != nil {
		_ = fm.LoadFonts(resources)
	}

	// Extract content stream
	content, err := e.reader.ExtractContent(pageDict)
	if err != nil {
		return nil, nil, nil, err
	}
	if content == nil {
		return nil, nil, nil, nil // Empty page
	}

	// Process content
	interpreter := pdf.NewInterpreter(fm, resources)
	textBlocks, images, graphics, err := interpreter.Process(content)
	if err != nil {
		return nil, nil, nil, err
	}

	return textBlocks, images, graphics, nil
}

// SetDebug enables debug logging
func (e *CustomExtractor) SetDebug(debug bool) {
	// TODO: Pass debug flag to reader/interpreter
}

// Close closes the PDF file
func (e *CustomExtractor) Close() error {
	return e.reader.Close()
}
