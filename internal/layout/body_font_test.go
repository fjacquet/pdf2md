package layout

import (
	"strings"
	"testing"

	"github.com/fjacquet/pdf2md/internal/extractor"
)

func TestBodyFontSize(t *testing.T) {
	tb := func(fs float64, s string) extractor.TextBlock { return extractor.TextBlock{FontSize: fs, Text: s} }
	line := strings.Repeat("x", 80)
	tests := []struct {
		name   string
		blocks []extractor.TextBlock
		want   float64
	}{
		{"empty", nil, 0},
		{
			"long body lines outweigh many short fragments",
			[]extractor.TextBlock{tb(10, line), tb(10, line), tb(8, "a"), tb(8, "b"), tb(8, "c"), tb(8, "d"), tb(8, "e")},
			10,
		},
		{"tie picks the smaller size", []extractor.TextBlock{tb(12, "abcd"), tb(10, "wxyz")}, 10},
		{"unknown sizes are ignored", []extractor.TextBlock{tb(0, line), tb(9, "body")}, 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BodyFontSize(tt.blocks); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

// A page full of small reference text keeps the document's body size, so its
// ordinary lines are not headings.
func TestAnalyze_UsesDocumentBodyFontSize(t *testing.T) {
	refs := strings.Repeat("Smith J. Some cited paper title, 2010. ", 3)
	content := &extractor.PageContent{TextBlocks: []extractor.TextBlock{
		{X: 72, Y: 700, Width: 450, FontSize: 10, Text: "We use the MICCAI 2018 data set and train for 100 epochs"},
		{X: 72, Y: 600, Width: 450, FontSize: 8, Text: refs},
		{X: 72, Y: 590, Width: 450, FontSize: 8, Text: refs},
		{X: 72, Y: 580, Width: 450, FontSize: 8, Text: refs},
	}}
	a := NewAnalyzer()
	a.BodyFontSize = 10
	for _, e := range a.Analyze(content) {
		if e.Type == ElementTypeHeader {
			t.Fatalf("body line became a heading: %q", e.Content)
		}
	}
}
