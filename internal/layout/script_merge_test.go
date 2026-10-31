package layout

import (
	"strings"
	"testing"
)

func TestMergeElements_Scripts(t *testing.T) {
	el := func(x, w, y, fs float64, s string) Element {
		return Element{Type: ElementTypeParagraph, X: x, Width: w, Y: y, FontSize: fs, Height: fs, Content: s}
	}
	tests := []struct {
		name   string
		in     []Element
		want   string
		wantY  float64
		wantFS float64
	}{
		{
			"superscript",
			[]Element{el(460.9, 22.9, 204.2, 10.9, "= 10"), el(483.8, 4, 208.1, 8, "9"), el(492.5, 8, 204.2, 10.9, "K.")},
			"= 10$^{9}$ K.", 204.2, 10.9,
		},
		{
			"subscript",
			[]Element{el(431.4, 6.4, 204.2, 10.9, "$T$"), el(437.8, 20, 202.5, 8, "AGN"), el(460.9, 22.9, 204.2, 10.9, "= 10")},
			"$T_{AGN}$ = 10", 204.2, 10.9,
		},
		{
			"leading footnote mark",
			[]Element{el(72, 4, 208, 7, "1"), el(76.5, 100, 204, 10.9, "Footnote text")},
			"$^{1}$Footnote text", 204, 10.9,
		},
		{
			"multi-piece exponent",
			[]Element{el(100, 10, 204.2, 10.9, "s"), el(110, 4, 208.1, 8, "−"), el(114, 4, 208.1, 8, "1")},
			"s$^{−1}$", 204.2, 10.9,
		},
		{
			"small text on its own baseline stays text",
			[]Element{el(72, 30, 204.2, 10.9, "see"), el(105, 20, 204.2, 8, "note")},
			"see note", 204.2, 10.9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAnalyzer().MergeElements(tt.in)
			if len(got) != 1 || got[0].Content != tt.want || got[0].Y != tt.wantY || got[0].FontSize != tt.wantFS {
				t.Fatalf("got %+v, want one element %q at y=%v fs=%v", got, tt.want, tt.wantY, tt.wantFS)
			}
		})
	}
}

// A glyph with no text mapping (a control character) is not a script: no
// empty "$^{}$" is produced.
func TestMergeElements_UnmappedGlyphIsNotAScript(t *testing.T) {
	in := []Element{
		{Type: ElementTypeParagraph, X: 72, Width: 30, Y: 204.2, FontSize: 10.9, Content: "Sch¨on"},
		{Type: ElementTypeParagraph, X: 102, Width: 4, Y: 208.1, FontSize: 8, Content: "\x03"},
	}
	for _, e := range NewAnalyzer().MergeElements(in) {
		if strings.Contains(e.Content, "$^{") {
			t.Fatalf("unmapped glyph became a script: %q", e.Content)
		}
	}
}
