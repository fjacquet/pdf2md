package render

import (
	"testing"
)

func TestDefaultRenderConfig(t *testing.T) {
	cfg := DefaultRenderConfig()

	if cfg == nil {
		t.Fatal("DefaultRenderConfig() returned nil")
	}

	if cfg.DPI != 150.0 {
		t.Errorf("DPI = %v, want 150.0", cfg.DPI)
	}

	if cfg.MaxDimension != 0 {
		t.Errorf("MaxDimension = %v, want 0", cfg.MaxDimension)
	}

	expectedBG := [3]uint8{255, 255, 255}
	if cfg.BackgroundColor != expectedBG {
		t.Errorf("BackgroundColor = %v, want %v", cfg.BackgroundColor, expectedBG)
	}
}

func TestRenderConfigForONNX(t *testing.T) {
	cfg := RenderConfigForONNX()

	if cfg == nil {
		t.Fatal("RenderConfigForONNX() returned nil")
	}

	if cfg.DPI != 150.0 {
		t.Errorf("DPI = %v, want 150.0", cfg.DPI)
	}

	if cfg.MaxDimension != 2048 {
		t.Errorf("MaxDimension = %v, want 2048", cfg.MaxDimension)
	}

	expectedBG := [3]uint8{255, 255, 255}
	if cfg.BackgroundColor != expectedBG {
		t.Errorf("BackgroundColor = %v, want %v", cfg.BackgroundColor, expectedBG)
	}
}

func TestDefaultRendererOptions(t *testing.T) {
	opts := DefaultRendererOptions()

	if opts == nil {
		t.Fatal("DefaultRendererOptions() returned nil")
	}

	if opts.PoolSize != 1 {
		t.Errorf("PoolSize = %v, want 1", opts.PoolSize)
	}
}

func TestRendererFactoryWithConfig(t *testing.T) {
	cfg := &RenderConfig{
		DPI:             200.0,
		MaxDimension:    1024,
		BackgroundColor: [3]uint8{0, 0, 0},
	}

	factory := RendererFactoryWithConfig(cfg)

	if factory == nil {
		t.Fatal("RendererFactoryWithConfig() returned nil")
	}

	// Factory should be a function type
	// We can't actually call it without a valid PDF, but we verify it's not nil
}

func TestErrorConstants(t *testing.T) {
	// Verify error constants are properly defined
	if ErrRendererNotInitialized == nil {
		t.Error("ErrRendererNotInitialized is nil")
	}

	if ErrInvalidPage == nil {
		t.Error("ErrInvalidPage is nil")
	}

	if ErrRenderFailed == nil {
		t.Error("ErrRenderFailed is nil")
	}

	if ErrPDFNotLoaded == nil {
		t.Error("ErrPDFNotLoaded is nil")
	}

	// Verify error messages
	expectedMsgs := map[error]string{
		ErrRendererNotInitialized: "renderer not initialized",
		ErrInvalidPage:            "invalid page number",
		ErrRenderFailed:           "failed to render page",
		ErrPDFNotLoaded:           "PDF not loaded",
	}

	for err, msg := range expectedMsgs {
		if err.Error() != msg {
			t.Errorf("Error message = %q, want %q", err.Error(), msg)
		}
	}
}
