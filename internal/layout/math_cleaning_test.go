package layout

import (
	"testing"
)

func TestCleanMathSymbols(t *testing.T) {
	analyzer := NewAnalyzer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple duplicate",
			input:    "$n$ $n$",
			expected: "$n$",
		},
		{
			name:     "No duplicate",
			input:    "$n$ $m$",
			expected: "$n$ $m$",
		},
		{
			name:     "Duplicate with text",
			input:    "Let $x$ $x$ be a variable",
			expected: "Let $x$ be a variable",
		},
		{
			name:     "Triple duplicate",
			input:    "$x$ $x$ $x$",
			expected: "$x$",
		},
		// New cases we want to handle
		{
			name:     "Duplicate without space",
			input:    "$n$$n$",
			expected: "$n$",
		},
		{
			name:     "Duplicate mixed",
			input:    "$n$ $n$$n$",
			expected: "$n$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements := []Element{
				{Type: ElementTypeParagraph, Content: tt.input},
			}
			cleaned := analyzer.CleanMathSymbols(elements)
			if cleaned[0].Content != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, cleaned[0].Content)
			}
		})
	}
}
