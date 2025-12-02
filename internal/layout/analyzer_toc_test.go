package layout

import (
	"testing"
)

func TestAnalyzer_RemoveTableOfContentsRange(t *testing.T) {
	analyzer := &Analyzer{}

	tests := []struct {
		name     string
		elements []Element
		expected []string // Content of remaining elements
	}{
		{
			name: "Standard TOC",
			elements: []Element{
				{Content: "Title"},
				{Content: "Table of Contents"},
				{Content: "1. Intro"},
				{Content: "2. Body"},
				{Content: "CHAPTER 1"},
				{Content: "Introduction text"},
			},
			expected: []string{"Title", "CHAPTER 1", "Introduction text"},
		},
		{
			name: "Abstract Fallback",
			elements: []Element{
				{Content: "Title"},
				{Content: "Abstract"},
				{Content: "This is the abstract text."},
				{Content: "1. Intro"},
				{Content: "CHAPTER 1"},
				{Content: "Introduction text"},
			},
			expected: []string{"Title", "Abstract", "This is the abstract text.", "CHAPTER 1", "Introduction text"},
		},
		{
			name: "No TOC",
			elements: []Element{
				{Content: "Title"},
				{Content: "Introduction"},
			},
			expected: []string{"Title", "Introduction"},
		},
		{
			name: "Start Found End Not Found",
			elements: []Element{
				{Content: "Title"},
				{Content: "Table of Contents"},
				{Content: "1. Intro"},
			},
			expected: []string{"Title", "Table of Contents", "1. Intro"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.RemoveTableOfContentsRange(tt.elements)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d elements, got %d", len(tt.expected), len(result))
			}
			for i, el := range result {
				if i < len(tt.expected) && el.Content != tt.expected[i] {
					t.Errorf("Index %d: expected %q, got %q", i, tt.expected[i], el.Content)
				}
			}
		})
	}
}
