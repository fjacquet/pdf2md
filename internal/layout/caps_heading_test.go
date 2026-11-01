package layout

import "testing"

// IEEE sets unnumbered section titles in small caps at body size or smaller:
// a short line of capitals alone is a heading, other capital text is not.
func TestClassify_CapitalsHeading(t *testing.T) {
	tests := map[string]bool{
		"REFERENCES":                true,
		"ACKNOWLEDGMENT":            true,
		"RELATED WORK":              true,
		"NASA":                      true,
		"IV":                        false, // too short: a numeral or an acronym
		"REFERENCES [1] Smith":      false,
		"THE QUICK BROWN FOX JUMPS": false, // over three words: a shouted sentence
		"References":                false,
		"U.S.A.":                    false,
		"TABLE II":                  false, // caption label: Roman numeral
		"A PREPRINT":                false, // one-letter word: a banner
		"MIX DESIGN":                true,  // letters of a numeral, not one
	}
	for text, want := range tests {
		e := Element{Type: ElementTypeParagraph, Content: text, FontSize: 8}
		NewAnalyzer().classifyElement(&e, nil, 10)
		if got := e.Type == ElementTypeHeader; got != want {
			t.Errorf("%q: heading = %v, want %v", text, got, want)
		}
	}
}

func TestClassify_CapitalsHeadingDropsBold(t *testing.T) {
	e := Element{Type: ElementTypeParagraph, Content: "**LICENSE**", FontSize: 10}
	NewAnalyzer().classifyElement(&e, nil, 10)
	if e.Type != ElementTypeHeader || e.Content != "LICENSE" || !e.Bold {
		t.Fatalf("got %+v", e)
	}
}
