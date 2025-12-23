package layout

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	// Verify default values
	if cfg.ColumnGapThreshold != 10.0 {
		t.Errorf("Expected ColumnGapThreshold=10.0, got %f", cfg.ColumnGapThreshold)
	}
	if cfg.HeaderSizeRatio != 1.2 {
		t.Errorf("Expected HeaderSizeRatio=1.2, got %f", cfg.HeaderSizeRatio)
	}
	if cfg.H1Ratio != 2.0 {
		t.Errorf("Expected H1Ratio=2.0, got %f", cfg.H1Ratio)
	}
	if cfg.H2Ratio != 1.75 {
		t.Errorf("Expected H2Ratio=1.75, got %f", cfg.H2Ratio)
	}
	if cfg.TableMinColumns != 2 {
		t.Errorf("Expected TableMinColumns=2, got %d", cfg.TableMinColumns)
	}
	if cfg.MinImageSize != 50.0 {
		t.Errorf("Expected MinImageSize=50.0, got %f", cfg.MinImageSize)
	}
	if cfg.EquationScoreThreshold != 2 {
		t.Errorf("Expected EquationScoreThreshold=2, got %d", cfg.EquationScoreThreshold)
	}
}

func TestConfigForAcademicPapers(t *testing.T) {
	cfg := ConfigForAcademicPapers()
	if cfg == nil {
		t.Fatal("ConfigForAcademicPapers() returned nil")
	}

	// Academic config should have lower equation threshold
	if cfg.EquationScoreThreshold != 1 {
		t.Errorf("Expected EquationScoreThreshold=1 for academic, got %d", cfg.EquationScoreThreshold)
	}
	// And wider column gap
	if cfg.ColumnGapThreshold != 15.0 {
		t.Errorf("Expected ColumnGapThreshold=15.0 for academic, got %f", cfg.ColumnGapThreshold)
	}
}

func TestConfigForTechnicalDocs(t *testing.T) {
	cfg := ConfigForTechnicalDocs()
	if cfg == nil {
		t.Fatal("ConfigForTechnicalDocs() returned nil")
	}

	// Technical config should have lower header ratio
	if cfg.HeaderSizeRatio != 1.15 {
		t.Errorf("Expected HeaderSizeRatio=1.15 for technical, got %f", cfg.HeaderSizeRatio)
	}
}

func TestConfigForScannedDocuments(t *testing.T) {
	cfg := ConfigForScannedDocuments()
	if cfg == nil {
		t.Fatal("ConfigForScannedDocuments() returned nil")
	}

	// Scanned config should have larger tolerances
	if cfg.ColumnGapThreshold != 20.0 {
		t.Errorf("Expected ColumnGapThreshold=20.0 for scanned, got %f", cfg.ColumnGapThreshold)
	}
	if cfg.IndentThreshold != 10.0 {
		t.Errorf("Expected IndentThreshold=10.0 for scanned, got %f", cfg.IndentThreshold)
	}
	if cfg.MinImageSize != 100.0 {
		t.Errorf("Expected MinImageSize=100.0 for scanned, got %f", cfg.MinImageSize)
	}
}

func TestNewAnalyzerWithConfig(t *testing.T) {
	// Test with nil config (should use defaults)
	a := NewAnalyzerWithConfig(nil)
	if a == nil {
		t.Fatal("NewAnalyzerWithConfig(nil) returned nil")
	}
	if a.Config == nil {
		t.Fatal("Analyzer.Config should not be nil")
	}
	if a.ColumnGapThreshold != 10.0 {
		t.Errorf("Expected default ColumnGapThreshold=10.0, got %f", a.ColumnGapThreshold)
	}

	// Test with custom config
	customCfg := &LayoutConfig{
		ColumnGapThreshold:     25.0,
		HeaderSizeRatio:        1.5,
		H1Ratio:                2.5,
		H2Ratio:                2.0,
		H3Ratio:                1.75,
		H4Ratio:                1.5,
		H5Ratio:                1.25,
		EquationScoreThreshold: 3,
		MinImageSize:           75.0,
	}
	a2 := NewAnalyzerWithConfig(customCfg)
	if a2.ColumnGapThreshold != 25.0 {
		t.Errorf("Expected custom ColumnGapThreshold=25.0, got %f", a2.ColumnGapThreshold)
	}
	if a2.HeaderSizeRatio != 1.5 {
		t.Errorf("Expected custom HeaderSizeRatio=1.5, got %f", a2.HeaderSizeRatio)
	}
}

func TestCalculateHeaderLevelWithConfig(t *testing.T) {
	// Test with custom header ratios
	customCfg := &LayoutConfig{
		H1Ratio: 2.5,
		H2Ratio: 2.0,
		H3Ratio: 1.75,
		H4Ratio: 1.5,
		H5Ratio: 1.25,
	}
	a := NewAnalyzerWithConfig(customCfg)

	// With body size 10, and custom ratios:
	// H1 requires 25+ (2.5x)
	// H2 requires 20+ (2.0x)
	// H3 requires 17.5+ (1.75x)
	bodySize := 10.0

	tests := []struct {
		fontSize float64
		expected int
	}{
		{25.0, 1}, // 2.5x = H1
		{20.0, 2}, // 2.0x = H2
		{17.5, 3}, // 1.75x = H3
		{15.0, 4}, // 1.5x = H4
		{12.5, 5}, // 1.25x = H5
		{10.0, 6}, // 1.0x = H6
	}

	for _, tt := range tests {
		level := a.calculateHeaderLevel(tt.fontSize, bodySize)
		if level != tt.expected {
			t.Errorf("Font %.1f / body %.1f with custom config: expected H%d, got H%d",
				tt.fontSize, bodySize, tt.expected, level)
		}
	}
}
