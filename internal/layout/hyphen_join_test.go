package layout

import "testing"

func TestJoinWrapped(t *testing.T) {
	tests := []struct{ cur, next, want string }{
		{"it is impor-", "tant to", "it is impor\u00ADtant to"}, // marked for ResolveWordBreaks
		{"it is impor\u2011", "tant", "it is impor\u00ADtant"},  // non-breaking hyphen
		{"it is impor\u2010", "tant", "it is impor\u00ADtant"},
		{"it is impor\u00AD", "tant", "it is impor\u00ADtant"},
		{"the input-", "Output pairs", "the input-Output pairs"}, // capital: keep hyphen
		{"version 2-", "3 only", "version 2-3 only"},
		{"ends with space ", "next", "ends with space next"},
		{"plain", "next", "plain next"},
	}
	for _, tt := range tests {
		if got := joinWrapped(tt.cur, tt.next); got != tt.want {
			t.Errorf("joinWrapped(%q, %q) = %q, want %q", tt.cur, tt.next, got, tt.want)
		}
	}
}

func TestMergeParagraphLines_MarksWordBreakHyphen(t *testing.T) {
	got := NewAnalyzer().MergeParagraphLines([]Element{
		{Type: ElementTypeParagraph, Y: 700, FontSize: 10, Content: "it is impor\u2011"},
		{Type: ElementTypeParagraph, Y: 688, FontSize: 10, Content: "tant to note."},
	})
	if len(got) != 1 || got[0].Content != "it is impor\u00ADtant to note." {
		t.Fatalf("got %+v", got)
	}
}

func TestResolveWordBreaks(t *testing.T) {
	got := NewAnalyzer().ResolveWordBreaks([]Element{
		{Content: "It is impor\u00ADtant and well\u00ADstudied."},
		{Content: "Important results."},
	})
	if want := "It is important and well-studied."; got[0].Content != want {
		t.Fatalf("got %q, want %q", got[0].Content, want)
	}
	if got[1].Content != "Important results." {
		t.Fatalf("untouched element changed: %q", got[1].Content)
	}
}
