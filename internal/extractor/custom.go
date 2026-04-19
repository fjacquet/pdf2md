// Package extractor provides PDF content extraction functionality.
package extractor

import (
	"fmt"

	"github.com/fjacquet/pdf2md/internal/pdf"
	"github.com/fjacquet/pdf2md/internal/types"
)

// PageContent holds all extracted content from a PDF page
type PageContent struct {
	TextBlocks []TextBlock
	Images     []types.Image
	Graphics   []types.VectorGraphic
	Links      []types.Link
	PageWidth  float64
	PageHeight float64
}

// CustomExtractor implements PDFExtractor using our custom parser
type CustomExtractor struct {
	reader *pdf.Reader
	debug  bool
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
func (e *CustomExtractor) GetPageCount() (int, error) {
	return e.reader.GetPageCount()
}

// ExtractTextBlocks extracts text blocks from a page
func (e *CustomExtractor) ExtractTextBlocks(pageIndex int) (*PageContent, error) {
	// Get page dictionary
	pageDict, err := e.reader.GetPage(pageIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to get page %d: %w", pageIndex, err)
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
			} else if e.debug {
				fmt.Printf("Warning: failed to read resources object: %v\n", err)
			}
		} else if resDict, ok := res.(pdf.Dictionary); ok {
			resources = resDict
		}
	}

	// Load fonts
	fm := pdf.NewFontManager(e.reader)
	if resources != nil {
		if err := fm.LoadFonts(resources); err != nil && e.debug {
			fmt.Printf("Warning: failed to load fonts: %v\n", err)
		}
	}

	// Extract content stream
	content, err := e.reader.ExtractContent(pageDict)
	if err != nil {
		return nil, fmt.Errorf("failed to extract content: %w", err)
	}

	// Extract links
	links, err := e.reader.ExtractLinks(pageDict)
	if err != nil && e.debug {
		fmt.Printf("Warning: failed to extract links: %v\n", err)
	}

	width, height := e.readPageDimensions(pageDict)

	if content == nil {
		return &PageContent{
			Links:      links,
			PageWidth:  width,
			PageHeight: height,
		}, nil
	}

	// Process content
	interpreter := pdf.NewInterpreter(fm, resources)
	textBlocks, images, graphics, err := interpreter.Process(content)
	if err != nil {
		return nil, fmt.Errorf("failed to process content: %w", err)
	}

	return &PageContent{
		TextBlocks: textBlocks,
		Images:     images,
		Graphics:   graphics,
		Links:      links,
		PageWidth:  width,
		PageHeight: height,
	}, nil
}

// readPageDimensions resolves the page's MediaBox and returns its width/height.
// Returns (0, 0) if MediaBox is missing or malformed.
func (e *CustomExtractor) readPageDimensions(pageDict pdf.Dictionary) (float64, float64) {
	mediaBoxObj, ok := pageDict[pdf.Name("MediaBox")]
	if !ok {
		return 0, 0
	}

	if ref, ok := mediaBoxObj.(pdf.IndirectRef); ok {
		obj, err := e.reader.ReadObject(ref.ObjectNumber)
		if err == nil {
			mediaBoxObj = obj
		}
	}

	arr, ok := mediaBoxObj.(pdf.Array)
	if !ok || len(arr) != 4 {
		return 0, 0
	}

	x1 := toFloat(arr[0])
	y1 := toFloat(arr[1])
	x2 := toFloat(arr[2])
	y2 := toFloat(arr[3])
	return abs(x2 - x1), abs(y2 - y1)
}

func toFloat(obj pdf.Object) float64 {
	if i, ok := obj.(pdf.Integer); ok {
		return float64(i)
	}
	if r, ok := obj.(pdf.Real); ok {
		return float64(r)
	}
	return 0
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// SetDebug enables debug logging
func (e *CustomExtractor) SetDebug(debug bool) {
	e.debug = debug
}

// Close closes the PDF file
func (e *CustomExtractor) Close() error {
	return e.reader.Close()
}
