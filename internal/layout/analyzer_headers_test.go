package layout

import (
	"slices"
	"testing"
)

func hdr(size float64, bold bool) Element {
	return Element{Type: ElementTypeHeader, Content: "H", FontSize: size, Bold: bold, Level: 4}
}

func levels(elements []Element) []int {
	var out []int
	for _, e := range elements {
		out = append(out, e.Level)
	}
	return out
}

func TestRankHeaderLevels(t *testing.T) {
	tests := []struct {
		name string
		in   []Element
		want []int
	}{
		{"title and sections", []Element{hdr(24, false), hdr(16, false), hdr(16, false), hdr(12, true)}, []int{1, 2, 2, 3}},
		{"bold above regular", []Element{hdr(14, false), hdr(14, true)}, []int{2, 1}},
		{"rounds size", []Element{hdr(16.1, false), hdr(15.9, false)}, []int{1, 1}},
		{"caps at six", []Element{hdr(30, false), hdr(28, false), hdr(26, false), hdr(24, false), hdr(22, false), hdr(20, false), hdr(18, false), hdr(16, false)}, []int{1, 2, 3, 4, 5, 6, 6, 6}},
		// Accepted limit: a large running label outranks the real title.
		{"running header takes top", []Element{hdr(30, false), hdr(24, false)}, []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := levels(NewAnalyzer().RankHeaderLevels(tt.in)); !slices.Equal(got, tt.want) {
				t.Errorf("levels = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRankHeaderLevels_NoHeaders(t *testing.T) {
	in := []Element{{Type: ElementTypeParagraph, Content: "p", FontSize: 10}}
	if got := NewAnalyzer().RankHeaderLevels(in); got[0].Level != 0 || got[0].Type != ElementTypeParagraph {
		t.Errorf("got %+v, want the paragraph unchanged", got)
	}
}
