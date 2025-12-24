package render

import (
	"fmt"
	"image"
	"log"
	"sync"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
)

// pdfiumInstance is the singleton go-pdfium instance.
// go-pdfium uses a pool internally, so we maintain a single global instance.
var (
	pdfiumPool     pdfium.Pool
	pdfiumInitOnce sync.Once
	pdfiumInitErr  error
)

// initPdfium initializes the go-pdfium WebAssembly runtime.
// This is called once on first use.
func initPdfium() error {
	pdfiumInitOnce.Do(func() {
		var err error
		pdfiumPool, err = webassembly.Init(webassembly.Config{
			MinIdle:  1,
			MaxIdle:  1,
			MaxTotal: 4,
		})
		if err != nil {
			pdfiumInitErr = fmt.Errorf("failed to init pdfium: %w", err)
		}
	})
	return pdfiumInitErr
}

// PdfiumRenderer implements PageRenderer using go-pdfium WebAssembly backend.
type PdfiumRenderer struct {
	instance  pdfium.Pdfium
	document  references.FPDF_DOCUMENT
	pageCount int
	pdfPath   string
	config    *RenderConfig
}

// NewPdfiumRenderer creates a new renderer for the given PDF file.
// Uses go-pdfium with WebAssembly/Wazero backend for cross-platform compatibility.
func NewPdfiumRenderer(pdfPath string, config *RenderConfig) (*PdfiumRenderer, error) {
	if config == nil {
		config = DefaultRenderConfig()
	}

	// Initialize pdfium pool if not already done
	if err := initPdfium(); err != nil {
		return nil, err
	}

	// Get an instance from the pool
	instance, err := pdfiumPool.GetInstance(30 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to get pdfium instance: %w", err)
	}

	// Open the PDF document
	doc, err := instance.OpenDocument(&requests.OpenDocument{
		FilePath: &pdfPath,
	})
	if err != nil {
		instance.Close()
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}

	// Get page count
	pageCountResp, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: doc.Document,
	})
	if err != nil {
		instance.Close()
		return nil, fmt.Errorf("failed to get page count: %w", err)
	}

	return &PdfiumRenderer{
		instance:  instance,
		document:  doc.Document,
		pageCount: pageCountResp.PageCount,
		pdfPath:   pdfPath,
		config:    config,
	}, nil
}

// RenderPage renders a specific page to an image.
// pageIndex is 0-based.
func (r *PdfiumRenderer) RenderPage(pageIndex int, dpi float64) (image.Image, error) {
	if r.instance == nil {
		return nil, ErrRendererNotInitialized
	}

	if pageIndex < 0 || pageIndex >= r.pageCount {
		return nil, fmt.Errorf("%w: %d (document has %d pages)", ErrInvalidPage, pageIndex, r.pageCount)
	}

	if dpi <= 0 {
		dpi = r.config.DPI
	}

	// Convert DPI to int for the API
	dpiInt := int(dpi)
	if dpiInt < 1 {
		dpiInt = 150
	}

	// Render the page
	renderResp, err := r.instance.RenderPageInDPI(&requests.RenderPageInDPI{
		DPI: dpiInt,
		Page: requests.Page{
			ByIndex: &requests.PageByIndex{
				Document: r.document,
				Index:    pageIndex,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRenderFailed, err)
	}

	if renderResp.Result.Image == nil {
		return nil, fmt.Errorf("%w: no image returned", ErrRenderFailed)
	}

	return renderResp.Result.Image, nil
}

// PageCount returns the total number of pages in the document.
func (r *PdfiumRenderer) PageCount() int {
	return r.pageCount
}

// PageSize returns the size of a page in PDF points (1/72 inch).
// pageIndex is 0-based.
func (r *PdfiumRenderer) PageSize(pageIndex int) (width, height float64, err error) {
	if r.instance == nil {
		return 0, 0, ErrRendererNotInitialized
	}

	if pageIndex < 0 || pageIndex >= r.pageCount {
		return 0, 0, fmt.Errorf("%w: %d", ErrInvalidPage, pageIndex)
	}

	sizeResp, err := r.instance.GetPageSize(&requests.GetPageSize{
		Page: requests.Page{
			ByIndex: &requests.PageByIndex{
				Document: r.document,
				Index:    pageIndex,
			},
		},
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get page size: %w", err)
	}

	return sizeResp.Width, sizeResp.Height, nil
}

// Close releases all resources.
func (r *PdfiumRenderer) Close() error {
	if r.instance != nil && r.document != "" {
		_, err := r.instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
			Document: r.document,
		})
		if err != nil {
			log.Printf("warning: failed to close PDF document: %v", err)
		}
		r.instance.Close()
	}

	return nil
}

// DefaultRendererFactory creates a PdfiumRenderer.
// This is the default factory for PageRenderer creation.
func DefaultRendererFactory(pdfPath string) (PageRenderer, error) {
	return NewPdfiumRenderer(pdfPath, DefaultRenderConfig())
}

// RendererFactoryWithConfig creates a factory with specific config.
func RendererFactoryWithConfig(config *RenderConfig) RendererFactory {
	return func(pdfPath string) (PageRenderer, error) {
		return NewPdfiumRenderer(pdfPath, config)
	}
}

// ClosePdfiumPool closes the global pdfium pool.
// Call this at program exit to cleanly release resources.
func ClosePdfiumPool() {
	if pdfiumPool != nil {
		pdfiumPool.Close()
	}
}
