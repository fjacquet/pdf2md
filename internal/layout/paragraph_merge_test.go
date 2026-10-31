package layout

import "testing"

// Three lines of one paragraph, 16.3 pt apart: the gap to each line is
// measured from the line above it, not from the paragraph's first line.
func TestMergeParagraphLines_ThreeLines(t *testing.T) {
	line := func(y float64, s string) Element {
		return Element{Type: ElementTypeParagraph, Y: y, X: 70.9, Height: 10.9, FontSize: 10.9, Content: s}
	}
	got := NewAnalyzer().MergeParagraphLines([]Element{
		line(457.7, "In each of the subfolders"),
		line(441.4, "you can find a list of resources that"),
		line(425.1, "you can use to start learning."),
	})
	if len(got) != 1 {
		t.Fatalf("got %d elements, want 1: %+v", len(got), got)
	}
	if want := "In each of the subfolders you can find a list of resources that you can use to start learning."; got[0].Content != want {
		t.Errorf("content = %q, want %q", got[0].Content, want)
	}
	if got[0].Y != 457.7 {
		t.Errorf("Y = %v, want the first line's 457.7", got[0].Y)
	}
}

// Side-by-side column paragraphs share their top Y: the global pass must
// compare the left paragraph's last line, not its top, with the right one.
func TestMergeParagraphLines_ColumnsStaySeparate(t *testing.T) {
	left := Element{Type: ElementTypeParagraph, X: 72, Y: 700, LastY: 658, FontSize: 12, Content: "Left column line four."}
	right := Element{Type: ElementTypeParagraph, X: 330, Y: 700, FontSize: 12, Content: "Right column line one."}
	if got := NewAnalyzer().MergeParagraphLines([]Element{left, right}); len(got) != 2 {
		t.Fatalf("got %d elements, want 2: %+v", len(got), got)
	}
}
