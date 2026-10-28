package layout

import (
	"slices"
	"testing"

	"github.com/fjacquet/pdf2md/internal/extractor"
)

func texts(blocks []extractor.TextBlock) []string {
	out := make([]string, len(blocks))
	for i, b := range blocks {
		out[i] = b.Text
	}
	return out
}

func line(text string, x, y float64) extractor.TextBlock {
	return extractor.TextBlock{Text: text, X: x, Y: y, Width: 200, Height: 10, FontSize: 10}
}

func TestStreamOrder_TwoColumnsInStreamOrder(t *testing.T) {
	blocks := []extractor.TextBlock{
		line("c1 l1", 50, 500), line("c1 l2", 50, 480),
		line("c2 l1", 300, 500), line("c2 l2", 300, 480),
	}
	if got := texts(OrderBlocks(blocks)); !slices.Equal(got, texts(blocks)) {
		t.Errorf("OrderBlocks() = %v, want stream order %v", got, texts(blocks))
	}
}

func TestStreamOrder_IncoherentFallsBack(t *testing.T) {
	blocks := []extractor.TextBlock{line("y100", 50, 100), line("y200", 50, 200), line("y300", 50, 300)}
	want := []string{"y300", "y200", "y100"}
	if got := texts(OrderBlocks(blocks)); !slices.Equal(got, want) {
		t.Errorf("OrderBlocks() = %v, want %v", got, want)
	}
}

// A header written after the body is one backward jump in 20 transitions:
// the stream order is kept and the header stays last. This is the accepted limit.
func TestStreamOrder_HeaderWrittenLast(t *testing.T) {
	var blocks []extractor.TextBlock
	for i := range 20 {
		blocks = append(blocks, line("body", 50, 600-float64(i)*14))
	}
	blocks = append(blocks, line("header", 50, 700))
	if !streamOrderCoherent(blocks) {
		t.Fatal("streamOrderCoherent() = false, want true (1 backward jump in 20 transitions)")
	}
	if got := OrderBlocks(blocks); got[len(got)-1].Text != "header" {
		t.Errorf("header moved to %v, want stream order kept", texts(got))
	}
}
