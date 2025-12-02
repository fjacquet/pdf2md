package layout

import (
	"strings"
	"testing"

	"github.com/fjacquet/pdf2md/internal/extractor"
)

func TestAnalyzer_EquationNumber_Repro(t *testing.T) {
	analyzer := NewAnalyzer()

	// Simulate: Equation text ... gap ... (1)
	blocks := []extractor.TextBlock{
		{
			Text:     "z^2 - pz - q = 0.",
			X:        100,
			Y:        500,
			Width:    150,
			Height:   12,
			FontSize: 12,
		},
		{
			Text:     "(",
			X:        400,
			Y:        500,
			Width:    5,
			Height:   12,
			FontSize: 12,
		},
		{
			Text:     "2",
			X:        405,
			Y:        500,
			Width:    10,
			Height:   12,
			FontSize: 12,
		},
		{
			Text:     ")",
			X:        415,
			Y:        500,
			Width:    5,
			Height:   12,
			FontSize: 12,
		},
	}

	content := &extractor.PageContent{
		TextBlocks: blocks,
		PageHeight: 800,
	}

	elements := analyzer.Analyze(content)

	if len(elements) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(elements))
	}

	el := elements[0]
	if el.Type == ElementTypeTable {
		t.Errorf("Expected Paragraph (Equation), got Table. Content: %s", el.Content)
	}

	// Check if it merged with space, not 4 spaces
	// Table formatting adds pipes, so if it's NOT a table, it shouldn't have pipes from MergeTableRows
	if strings.Contains(el.Content, "|") {
		t.Errorf("Content seems to be formatted as table: %s", el.Content)
	}
}
