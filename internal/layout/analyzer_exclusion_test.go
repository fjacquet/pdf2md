package layout

import (
	"testing"

	"github.com/fjacquet/pdf2md/internal/extractor"
)

func TestExclusionZones(t *testing.T) {
	analyzer := NewAnalyzer()
	// Set exclusion zones: Top 100, Bottom 100
	analyzer.Exclusion = ExclusionZone{
		Top:    100,
		Bottom: 100,
	}

	// Page height 1000
	// Top margin: Y > 1000 - 100 = 900 excluded
	// Bottom margin: Y < 100 excluded
	content := &extractor.PageContent{
		PageHeight: 1000,
		TextBlocks: []extractor.TextBlock{
			{Text: "Header", Y: 950}, // Should be excluded
			{Text: "Body 1", Y: 800}, // Keep
			{Text: "Body 2", Y: 500}, // Keep
			{Text: "Footer", Y: 50},  // Should be excluded
		},
	}

	elements := analyzer.Analyze(content)

	if len(elements) != 2 {
		t.Fatalf("Expected 2 elements, got %d", len(elements))
	}

	if elements[0].Content != "Body 1" {
		t.Errorf("Expected 'Body 1', got '%s'", elements[0].Content)
	}
	if elements[1].Content != "Body 2" {
		t.Errorf("Expected 'Body 2', got '%s'", elements[1].Content)
	}
}
