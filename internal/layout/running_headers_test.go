package layout

import (
	"slices"
	"testing"
)

func TestFilterRunningHeaders(t *testing.T) {
	p := func(y float64, s string) Element { return Element{Type: ElementTypeParagraph, Y: y, Content: s} }
	h := func(y float64, s string) Element { return Element{Type: ElementTypeHeader, Y: y, Content: s} }
	page := func(top, body, bottom string) []Element {
		return []Element{p(780, top), p(400, body), p(40, bottom)}
	}
	tests := []struct {
		name  string
		pages [][]Element
		want  []string
	}{
		{
			"header and numbered footer on every page go",
			[][]Element{
				page("3D ERT A PREPRINT", "one", "Page 1 of 4"),
				page("3D ERT A PREPRINT", "two", "Page 2 of 4"),
				page("3D ERT A PREPRINT", "three", "Page 3 of 4"),
				page("3D ERT A PREPRINT", "four", "Page 4 of 4"),
			},
			[]string{"one", "two", "three", "four"},
		},
		{
			"alternating even and odd headers go",
			[][]Element{
				page("Title", "a", "x"), page("Authors", "b", "y"), page("Title", "c", "z"),
				page("Authors", "d", "w"), page("Title", "e", "v"), page("Authors", "f", "u"),
			},
			[]string{"a", "x", "b", "y", "c", "z", "d", "w", "e", "v", "f", "u"},
		},
		{
			"repeated text away from the page edge stays",
			[][]Element{
				{p(780, "Alpha"), p(400, "See text"), p(40, "Omega")},
				{p(780, "Beta"), p(400, "See text"), p(40, "Psi")},
				{p(780, "Gamma"), p(400, "See text"), p(40, "Chi")},
			},
			[]string{"Alpha", "See text", "Omega", "Beta", "See text", "Psi", "Gamma", "See text", "Chi"},
		},
		{
			"a title heading matching the running header stays",
			[][]Element{
				{h(780, "Deep Nets"), p(400, "intro")},
				{p(780, "Deep Nets"), p(400, "b")},
				{p(780, "Deep Nets"), p(400, "c")},
				{p(780, "Deep Nets"), p(400, "d")},
			},
			[]string{"Deep Nets", "intro", "b", "c", "d"},
		},
		{
			"twice is not a running header",
			[][]Element{page("Note", "a", "1"), page("Note", "b", "2")},
			[]string{"Note", "a", "1", "Note", "b", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			for _, e := range NewAnalyzer().FilterRunningHeaders(tt.pages) {
				got = append(got, e.Content)
			}
			if !slices.Equal(got, tt.want) {
				t.Fatalf("got %q\nwant %q", got, tt.want)
			}
		})
	}
}
