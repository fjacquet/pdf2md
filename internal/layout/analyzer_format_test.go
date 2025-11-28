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
