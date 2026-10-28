package layout

import (
	"math"

	"github.com/fjacquet/pdf2md/internal/extractor"
)

// OrderBlocks keeps the content-stream order, the order the producer wrote
// the page in, and falls back to XY-cut only when that order is incoherent.
func OrderBlocks(blocks []extractor.TextBlock) []extractor.TextBlock {
	if streamOrderCoherent(blocks) {
		return blocks
	}
	return SortBlocks(blocks)
}

// streamOrderCoherent reports whether at most 10% of line transitions jump
// back up within the same column: a column-ordered stream only moves up when
// it changes column.
func streamOrderCoherent(blocks []extractor.TextBlock) bool {
	type span struct{ y, x0, x1, size float64 }
	var lines []span
	for _, b := range blocks {
		if n := len(lines); n > 0 && math.Abs(b.Y-lines[n-1].y) < 0.5*lines[n-1].size {
			lines[n-1].x0 = math.Min(lines[n-1].x0, b.X)
			lines[n-1].x1 = math.Max(lines[n-1].x1, b.X+b.Width)
			continue
		}
		lines = append(lines, span{b.Y, b.X, b.X + b.Width, b.FontSize})
	}
	jumps := 0
	for i := 1; i < len(lines); i++ {
		prev, next := lines[i-1], lines[i]
		if next.y > prev.y+prev.size && next.x0 < prev.x1 && prev.x0 < next.x1 {
			jumps++
		}
	}
	return jumps*10 <= len(lines)-1
}
