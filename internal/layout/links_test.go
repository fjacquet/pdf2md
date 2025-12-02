package layout

import (
	"testing"

	"github.com/fjacquet/pdf2md/internal/extractor"
	"github.com/fjacquet/pdf2md/internal/types"
)

func TestApplyLinksToBlocks(t *testing.T) {
	blocks := []extractor.TextBlock{
		{Text: "Click here for more info", X: 10, Y: 10, Width: 100, Height: 10},
		{Text: "No link here", X: 10, Y: 30, Width: 100, Height: 10},
	}

	links := []types.Link{
		{
			URI:  "https://example.com",
			Rect: []float64{10, 10, 110, 20}, // Covers the first block roughly
		},
	}

	// Apply links
	newBlocks := ApplyLinksToBlocks(blocks, links)

	if len(newBlocks) != 2 {
		t.Fatalf("Expected 2 blocks, got %d", len(newBlocks))
	}

	// Check first block
	if newBlocks[0].LinkURI != "https://example.com" {
		t.Errorf("Expected LinkURI 'https://example.com', got '%s'", newBlocks[0].LinkURI)
	}
	if newBlocks[0].Text != "Click here for more info" {
		t.Errorf("Expected text 'Click here for more info', got '%s'", newBlocks[0].Text)
	}

	// Check second block
	if newBlocks[1].Text != "No link here" {
		t.Errorf("Expected 'No link here', got '%s'", newBlocks[1].Text)
	}
}
