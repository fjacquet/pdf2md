package layout

import "github.com/fjacquet/pdf2md/internal/types"

// dropFigureText removes plot labels (axis ticks, legends): text under 0.7x
// the body size inside a figure-sized graphic. A graphic that also holds text
// of 0.9x body or more is a container (framed box, page background), not a
// figure, and keeps everything. A block is tested at its horizontal centre on
// the baseline, so rotated labels (width 0) are caught too.
func dropFigureText(blocks []types.TextBlock, graphics []types.VectorGraphic, body, minSize float64) []types.TextBlock {
	if body <= 0 {
		return blocks
	}
	inside := func(g types.VectorGraphic, b types.TextBlock) bool {
		x := b.X + b.Width/2
		return x >= g.X && x <= g.X+g.Width && b.Y >= g.Y && b.Y <= g.Y+g.Height
	}
	var figures []types.VectorGraphic
	for _, g := range graphics {
		if g.Width < minSize || g.Height < minSize {
			continue
		}
		container := false
		for _, b := range blocks {
			if b.FontSize >= 0.9*body && inside(g, b) {
				container = true
				break
			}
		}
		if !container {
			figures = append(figures, g)
		}
	}
	if len(figures) == 0 {
		return blocks
	}
	kept := make([]types.TextBlock, 0, len(blocks))
	for _, b := range blocks {
		drop := false
		if b.FontSize < 0.7*body {
			for _, g := range figures {
				if inside(g, b) {
					drop = true
					break
				}
			}
		}
		if !drop {
			kept = append(kept, b)
		}
	}
	return kept
}
