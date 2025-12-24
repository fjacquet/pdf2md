package render

import (
	"errors"
	"image"
	"os"
	"path/filepath"
	"testing"
)

func getTestPDFPath(t *testing.T) string {
	t.Helper()
	// Find the project root by looking for go.mod
	cwd, err := os.Getwd()
	if err != nil {
		t.Skip("cannot get working directory")
	}

	// Navigate up to find testdata
	for dir := cwd; dir != "/"; dir = filepath.Dir(dir) {
		testPath := filepath.Join(dir, "testdata", "test.pdf")
		if _, err := os.Stat(testPath); err == nil {
			return testPath
		}
	}

	t.Skip("testdata/test.pdf not found")
	return ""
}

func TestNewPdfiumRenderer(t *testing.T) {
	pdfPath := getTestPDFPath(t)

	renderer, err := NewPdfiumRenderer(pdfPath, nil)
	if err != nil {
		t.Fatalf("NewPdfiumRenderer() error = %v", err)
	}
	defer renderer.Close()

	if renderer.PageCount() <= 0 {
		t.Errorf("PageCount() = %d, want > 0", renderer.PageCount())
	}
}

func TestNewPdfiumRenderer_WithConfig(t *testing.T) {
	pdfPath := getTestPDFPath(t)

	config := &RenderConfig{
		DPI:             72.0,
		MaxDimension:    1024,
		BackgroundColor: [3]uint8{255, 255, 255},
	}

	renderer, err := NewPdfiumRenderer(pdfPath, config)
	if err != nil {
		t.Fatalf("NewPdfiumRenderer() error = %v", err)
	}
	defer renderer.Close()

	if renderer.PageCount() <= 0 {
		t.Errorf("PageCount() = %d, want > 0", renderer.PageCount())
	}
}

func TestNewPdfiumRenderer_InvalidPath(t *testing.T) {
	_, err := NewPdfiumRenderer("/nonexistent/path.pdf", nil)
	if err == nil {
		t.Error("NewPdfiumRenderer() expected error for invalid path")
	}
}

func TestPdfiumRenderer_PageSize(t *testing.T) {
	pdfPath := getTestPDFPath(t)

	renderer, err := NewPdfiumRenderer(pdfPath, nil)
	if err != nil {
		t.Fatalf("NewPdfiumRenderer() error = %v", err)
	}
	defer renderer.Close()

	width, height, err := renderer.PageSize(0)
	if err != nil {
		t.Fatalf("PageSize() error = %v", err)
	}

	// Standard PDF pages should have positive dimensions
	if width <= 0 {
		t.Errorf("PageSize() width = %v, want > 0", width)
	}
	if height <= 0 {
		t.Errorf("PageSize() height = %v, want > 0", height)
	}
}

func TestPdfiumRenderer_PageSize_InvalidPage(t *testing.T) {
	pdfPath := getTestPDFPath(t)

	renderer, err := NewPdfiumRenderer(pdfPath, nil)
	if err != nil {
		t.Fatalf("NewPdfiumRenderer() error = %v", err)
	}
	defer renderer.Close()

	// Test negative page index
	_, _, err = renderer.PageSize(-1)
	if err == nil {
		t.Error("PageSize(-1) expected error")
	}
	if !errors.Is(err, ErrInvalidPage) {
		t.Errorf("PageSize(-1) error = %v, want ErrInvalidPage", err)
	}

	// Test page index beyond count
	_, _, err = renderer.PageSize(renderer.PageCount())
	if err == nil {
		t.Error("PageSize(beyond count) expected error")
	}
}

func TestPdfiumRenderer_RenderPage(t *testing.T) {
	pdfPath := getTestPDFPath(t)

	renderer, err := NewPdfiumRenderer(pdfPath, nil)
	if err != nil {
		t.Fatalf("NewPdfiumRenderer() error = %v", err)
	}
	defer renderer.Close()

	// Render first page at default DPI
	img, err := renderer.RenderPage(0, 0) // 0 uses config default
	if err != nil {
		t.Fatalf("RenderPage() error = %v", err)
	}

	if img == nil {
		t.Fatal("RenderPage() returned nil image")
	}

	bounds := img.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Errorf("RenderPage() image size = %dx%d, want > 0", bounds.Dx(), bounds.Dy())
	}
}

func TestPdfiumRenderer_RenderPage_CustomDPI(t *testing.T) {
	pdfPath := getTestPDFPath(t)

	renderer, err := NewPdfiumRenderer(pdfPath, nil)
	if err != nil {
		t.Fatalf("NewPdfiumRenderer() error = %v", err)
	}
	defer renderer.Close()

	// Render at low DPI
	imgLow, err := renderer.RenderPage(0, 72)
	if err != nil {
		t.Fatalf("RenderPage(72 DPI) error = %v", err)
	}

	// Render at high DPI
	imgHigh, err := renderer.RenderPage(0, 150)
	if err != nil {
		t.Fatalf("RenderPage(150 DPI) error = %v", err)
	}

	// Higher DPI should produce larger image
	lowBounds := imgLow.Bounds()
	highBounds := imgHigh.Bounds()

	// 150/72 ~ 2.08x, but allow some tolerance
	if highBounds.Dx() <= lowBounds.Dx() || highBounds.Dy() <= lowBounds.Dy() {
		t.Errorf("Higher DPI should produce larger image: low=%dx%d, high=%dx%d",
			lowBounds.Dx(), lowBounds.Dy(), highBounds.Dx(), highBounds.Dy())
	}
}

func TestPdfiumRenderer_RenderPage_InvalidPage(t *testing.T) {
	pdfPath := getTestPDFPath(t)

	renderer, err := NewPdfiumRenderer(pdfPath, nil)
	if err != nil {
		t.Fatalf("NewPdfiumRenderer() error = %v", err)
	}
	defer renderer.Close()

	// Test negative page index
	_, err = renderer.RenderPage(-1, 150)
	if err == nil {
		t.Error("RenderPage(-1) expected error")
	}
	if !errors.Is(err, ErrInvalidPage) {
		t.Errorf("RenderPage(-1) error = %v, want ErrInvalidPage", err)
	}

	// Test page index beyond count
	_, err = renderer.RenderPage(renderer.PageCount(), 150)
	if err == nil {
		t.Error("RenderPage(beyond count) expected error")
	}
}

func TestPdfiumRenderer_Close(t *testing.T) {
	pdfPath := getTestPDFPath(t)

	renderer, err := NewPdfiumRenderer(pdfPath, nil)
	if err != nil {
		t.Fatalf("NewPdfiumRenderer() error = %v", err)
	}

	// Close should not error
	err = renderer.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// After close, operations should fail gracefully
	_, err = renderer.RenderPage(0, 150)
	if err == nil {
		// May or may not error, but should not panic
		t.Log("RenderPage after Close() did not error (may be acceptable)")
	}
}

func TestDefaultRendererFactory(t *testing.T) {
	pdfPath := getTestPDFPath(t)

	renderer, err := DefaultRendererFactory(pdfPath)
	if err != nil {
		t.Fatalf("DefaultRendererFactory() error = %v", err)
	}
	defer renderer.Close()

	// Verify it's a working renderer
	if renderer.PageCount() <= 0 {
		t.Errorf("PageCount() = %d, want > 0", renderer.PageCount())
	}
}

// Verify PdfiumRenderer implements PageRenderer interface
func TestPdfiumRenderer_ImplementsPageRenderer(t *testing.T) {
	var _ PageRenderer = (*PdfiumRenderer)(nil)
}

// Verify the rendered image is actually an image.Image
func TestPdfiumRenderer_RenderPage_ImageType(t *testing.T) {
	pdfPath := getTestPDFPath(t)

	renderer, err := NewPdfiumRenderer(pdfPath, nil)
	if err != nil {
		t.Fatalf("NewPdfiumRenderer() error = %v", err)
	}
	defer renderer.Close()

	img, err := renderer.RenderPage(0, 150)
	if err != nil {
		t.Fatalf("RenderPage() error = %v", err)
	}

	// Verify it's a valid image
	var _ image.Image = img

	// Check that ColorModel is not nil
	if img.ColorModel() == nil {
		t.Error("image.ColorModel() is nil")
	}

	// Check bounds are valid
	bounds := img.Bounds()
	if bounds.Empty() {
		t.Error("image bounds are empty")
	}
}
