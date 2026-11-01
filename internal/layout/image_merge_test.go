package layout

import "testing"

// Images are never merged: side-by-side subfigures stay separate links and an
// image beside a text line (or just above its baseline) is neither text nor a
// superscript.
func TestMergeElements_ImagesStaySeparate(t *testing.T) {
	img := func(x, y float64, id string) Element {
		return Element{Type: ElementTypeImage, X: x, Y: y, Width: 100, Height: 80, Content: id}
	}
	text := Element{Type: ElementTypeParagraph, X: 310, Y: 300, Width: 50, FontSize: 10, Height: 10, Content: "Figure"}
	tests := map[string][]Element{
		"two subfigures":        {img(100, 300, "a.svg"), img(200, 300, "b.svg")},
		"image then text":       {img(210, 300, "a.svg"), text},
		"drop cap beside image": {{Type: ElementTypeParagraph, X: 90, Y: 300, Width: 8, FontSize: 20, Height: 20, Content: "A"}, img(100, 300, "images/a.svg")},
		"image above text":      {text, img(360, 303, "a.svg")},
	}
	for name, in := range tests {
		t.Run(name, func(t *testing.T) {
			if got := NewAnalyzer().MergeElements(in); len(got) != 2 {
				t.Fatalf("got %d elements, want 2: %+v", len(got), got)
			}
		})
	}
}
