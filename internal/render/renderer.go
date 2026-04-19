// Package render provides PDF-to-image rendering via go-pdfium.
// Used by ONNX layout detection and OCR fallback to produce page images.
package render

import (
	"errors"
	"image"
)

// Common errors for the render package.
var (
	ErrRendererNotInitialized = errors.New("renderer not initialized")
	ErrInvalidPage            = errors.New("invalid page number")
	ErrRenderFailed           = errors.New("failed to render page")
	ErrPDFNotLoaded           = errors.New("PDF not loaded")
)

// PageRenderer renders PDF pages to images.
// This interface abstracts the PDF rendering implementation,
// allowing for different backends (go-pdfium, poppler, etc.).
type PageRenderer interface {
	// RenderPage renders a specific page to an image.
	// pageIndex is 0-based.
	// dpi controls the resolution (72 = 1:1, 150 = 2x, etc.)
	RenderPage(pageIndex int, dpi float64) (image.Image, error)

	// PageCount returns the total number of pages in the document.
	PageCount() int

	// PageSize returns the size of a page in PDF points (1/72 inch).
	// pageIndex is 0-based.
	PageSize(pageIndex int) (width, height float64, err error)

	// Close releases all resources.
	Close() error
}

// RenderConfig holds configuration for page rendering.
//
//nolint:revive // name is exported API; renaming would churn callers
type RenderConfig struct {
	// DPI for rendering (default: 150 for good ONNX input quality)
	DPI float64

	// MaxDimension limits the maximum width or height of rendered images.
	// 0 means no limit. This helps manage memory for large pages.
	MaxDimension int

	// BackgroundColor for transparent pages (default: white)
	BackgroundColor [3]uint8
}

// DefaultRenderConfig returns the default rendering configuration.
func DefaultRenderConfig() *RenderConfig {
	return &RenderConfig{
		DPI:             150.0,
		MaxDimension:    0,
		BackgroundColor: [3]uint8{255, 255, 255}, // White
	}
}

// RenderConfigForONNX returns configuration optimized for ONNX model input.
// Uses 150 DPI which typically produces images around 1200x1500 pixels
// for standard letter/A4 pages, suitable for the 1024x1024 YOLO input.
//
//nolint:revive // name is exported API; renaming would churn callers
func RenderConfigForONNX() *RenderConfig {
	return &RenderConfig{
		DPI:             150.0,
		MaxDimension:    2048, // Cap to prevent memory issues
		BackgroundColor: [3]uint8{255, 255, 255},
	}
}

// RendererFactory creates PageRenderer instances.
// This is a function type for dependency injection.
type RendererFactory func(pdfPath string) (PageRenderer, error)

// NewRendererOptions holds options for creating a new renderer.
type NewRendererOptions struct {
	// PoolSize is the number of renderer instances to pool.
	// go-pdfium benefits from pooling for concurrent rendering.
	PoolSize int
}

// DefaultRendererOptions returns the default renderer options.
func DefaultRendererOptions() *NewRendererOptions {
	return &NewRendererOptions{
		PoolSize: 1,
	}
}
