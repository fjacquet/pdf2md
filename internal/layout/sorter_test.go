package layout

import (
	"testing"

	"github.com/fjacquet/pdf2md/internal/extractor"
)

func TestSortBlocks(t *testing.T) {
	tests := []struct {
		name     string
		blocks   []extractor.TextBlock
		expected []string // Expected text order
	}{
		{
			name: "Single Column",
			blocks: []extractor.TextBlock{
				{Text: "Line 2", X: 10, Y: 80, Width: 100, Height: 10},
				{Text: "Line 1", X: 10, Y: 90, Width: 100, Height: 10},
				{Text: "Line 3", X: 10, Y: 70, Width: 100, Height: 10},
			},
			expected: []string{"Line 1", "Line 2", "Line 3"},
		},
		{
			name: "Two Columns (Simple)",
			blocks: []extractor.TextBlock{
				{Text: "Col1 Line1-", X: 50, Y: 500, Width: 200, Height: 10},
				{Text: "Col2 Line1-", X: 300, Y: 500, Width: 200, Height: 10},
				{Text: "Col1 Line2", X: 50, Y: 480, Width: 200, Height: 10},
				{Text: "Col2 Line2", X: 300, Y: 480, Width: 200, Height: 10},
			},
			// Expect Column 1 then Column 2
			expected: []string{"Col1 Line1-", "Col1 Line2", "Col2 Line1-", "Col2 Line2"},
		},
		{
			name: "Header and Two Columns",
			blocks: []extractor.TextBlock{
				{Text: "Header", X: 50, Y: 600, Width: 450, Height: 10}, // Spans both
				{Text: "Col1 Line1", X: 50, Y: 500, Width: 200, Height: 10},
				{Text: "Col2 Line1", X: 300, Y: 500, Width: 200, Height: 10},
			},
			// Expect Header, then Col1, then Col2
			expected: []string{"Header", "Col1 Line1", "Col2 Line1"},
		},
		{
			name: "Interleaved Input",
			blocks: []extractor.TextBlock{
				{Text: "Col1 Line1-", X: 50, Y: 500, Width: 200, Height: 10},
				{Text: "Col2 Line1-", X: 300, Y: 500, Width: 200, Height: 10},
				{Text: "Col1 Line2", X: 50, Y: 480, Width: 200, Height: 10},
				{Text: "Col2 Line2", X: 300, Y: 480, Width: 200, Height: 10},
			},
			expected: []string{"Col1 Line1-", "Col1 Line2", "Col2 Line1-", "Col2 Line2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sorted := SortBlocks(tt.blocks)
			if len(sorted) != len(tt.expected) {
				t.Fatalf("Expected %d blocks, got %d", len(tt.expected), len(sorted))
			}

			for i, block := range sorted {
				if block.Text != tt.expected[i] {
					t.Errorf("Index %d: expected %s, got %s", i, tt.expected[i], block.Text)
				}
			}
		})
	}
}
