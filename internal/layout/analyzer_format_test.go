package layout

import (
	"testing"

	"github.com/fjacquet/pdf2md/internal/extractor"
)

func TestBlocksToElements_Formatting(t *testing.T) {
	analyzer := NewAnalyzer()

	blocks := []extractor.TextBlock{
		{Text: "Regular text", FontName: "/F1", FontSize: 12},
		{Text: "Bold text", FontName: "Arial-Bold", FontSize: 12},
		{Text: "Italic text", FontName: "Times-Italic", FontSize: 12},
		{Text: "BoldItalic text", FontName: "Helvetica-BoldOblique", FontSize: 12},
	}

	elements := analyzer.blocksToElements(blocks)

	if len(elements) != 4 {
		t.Fatalf("Expected 4 elements, got %d", len(elements))
	}

	expected := []string{
		"Regular text",
		"**Bold text**",
		"*Italic text*",
		"**BoldItalic text**", // Bold takes precedence in my logic
	}

	for i, elem := range elements {
		if elem.Content != expected[i] {
			t.Errorf("Element %d: expected '%s', got '%s'", i, expected[i], elem.Content)
		}
	}
}

// A PDF that draws each word space as its own block must still yield one bold
// run, not "**Services** **réservés**" (Reka booking confirmation).
func TestBlocksToElements_SpaceBlockKeepsBoldRun(t *testing.T) {
	analyzer := NewAnalyzer()
	blocks := []extractor.TextBlock{
		{Text: "Services", X: 60, Width: 48.4, Y: 590.8, FontName: "ArialBold", FontSize: 12, SpaceWidth: 2.71},
		{Text: " ", X: 108.4, Width: 2.7, Y: 590.8, FontName: "ArialBold", FontSize: 12, SpaceWidth: 2.71},
		{Text: "réservés", X: 111.1, Width: 50, Y: 590.8, FontName: "ArialBold", FontSize: 12, SpaceWidth: 2.71},
	}

	elements := analyzer.MergeElements(analyzer.blocksToElements(blocks))

	if len(elements) != 1 || elements[0].Content != "**Services réservés**" {
		t.Fatalf("got %+v, want one element \"**Services réservés**\"", elements)
	}

	// Glyphs overlap the space's advance: the gap alone is no word gap, so
	// the space block itself must carry the word break.
	blocks = []extractor.TextBlock{
		{Text: "Room", X: 51.38, Width: 24.91, Y: 252.9, FontName: "ArialBold", FontSize: 10, SpaceWidth: 2.56},
		{Text: " ", X: 75.01, Width: 2.56, Y: 252.9, FontName: "ArialBold", FontSize: 10, SpaceWidth: 2.56},
		{Text: "1:", X: 77.32, Width: 7.67, Y: 252.9, FontName: "ArialBold", FontSize: 10, SpaceWidth: 2.56},
	}

	elements = analyzer.MergeElements(analyzer.blocksToElements(blocks))

	if len(elements) != 1 || elements[0].Content != "**Room 1:**" {
		t.Fatalf("got %+v, want one element \"**Room 1:**\"", elements)
	}
}

func TestIsListItem_Markdown(t *testing.T) {
	analyzer := NewAnalyzer()

	cases := []struct {
		text     string
		expected bool
	}{
		{"- Item", true},
		{"**- Item**", true},
		{"*• Item*", true},
		{"1. Item", true},
		{"**1. Item**", true},
		{"Regular text", false},
		{"**Regular**", false},
	}

	for _, c := range cases {
		if got := analyzer.isListItem(c.text); got != c.expected {
			t.Errorf("isListItem('%s'): expected %v, got %v", c.text, c.expected, got)
		}
	}
}

func TestAdmonition_Markdown(t *testing.T) {
	analyzer := NewAnalyzer()

	// Find Admonition rule
	var rule Rule
	for _, r := range analyzer.Rules {
		if r.Name == "Admonition" {
			rule = r
			break
		}
	}

	if rule.Name == "" {
		t.Fatal("Admonition rule not found")
	}

	cases := []struct {
		text     string
		expected bool
	}{
		{"IMPORTANT: check this", true},
		{"**IMPORTANT**: check this", true},
		{"*NOTE*: note this", true},
		{"Regular text", false},
	}

	for _, c := range cases {
		if got := rule.Condition(c.text, 12, 12); got != c.expected {
			t.Errorf("Admonition('%s'): expected %v, got %v", c.text, c.expected, got)
		}
	}
}
