package layout

import (
	"slices"
	"testing"

	"github.com/fjacquet/pdf2md/internal/types"
)

func TestDropFigureText(t *testing.T) {
	plot := types.VectorGraphic{X: 72, Y: 619, Width: 451, Height: 150}
	page := types.VectorGraphic{X: 0, Y: 0, Width: 595, Height: 842}
	tb := func(x, y, w, fs float64, s string) types.TextBlock {
		return types.TextBlock{X: x, Y: y, Width: w, FontSize: fs, Text: s}
	}
	tests := []struct {
		name     string
		blocks   []types.TextBlock
		graphics []types.VectorGraphic
		want     []string
	}{
		{
			"plot labels inside a figure go, caption below stays",
			[]types.TextBlock{
				tb(97, 634, 97, 5.5, "-0.1 0.1 0.2"),
				tb(378, 683, 0, 5.5, "Sample Quantiles"), // rotated: width 0
				tb(72, 598, 319, 10, "Figure 1: Samples"),
			},
			[]types.VectorGraphic{plot},
			[]string{"Figure 1: Samples"},
		},
		{
			"small script inside a box holding body text stays",
			[]types.TextBlock{tb(100, 650, 200, 10, "Algorithm 1 step"), tb(150, 640, 5, 5, "i")},
			[]types.VectorGraphic{plot},
			[]string{"Algorithm 1 step", "i"},
		},
		{
			"small text on a full-page background stays",
			[]types.TextBlock{tb(72, 700, 400, 10, "Body text"), tb(72, 60, 100, 5.5, "footnote")},
			[]types.VectorGraphic{page},
			[]string{"Body text", "footnote"},
		},
		{
			"table text at 0.8x body stays",
			[]types.TextBlock{tb(100, 650, 40, 8, "0.93"), tb(72, 500, 400, 10, "Body text")},
			[]types.VectorGraphic{plot},
			[]string{"0.93", "Body text"},
		},
		{
			"graphic smaller than a figure is ignored",
			[]types.TextBlock{tb(100, 650, 10, 5.5, "x"), tb(72, 500, 400, 10, "Body text")},
			[]types.VectorGraphic{{X: 95, Y: 645, Width: 30, Height: 30}},
			[]string{"x", "Body text"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			for _, b := range dropFigureText(tt.blocks, tt.graphics, 10, 50) {
				got = append(got, b.Text)
			}
			if !slices.Equal(got, tt.want) {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
