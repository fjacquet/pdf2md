package markdown

import (
	"strings"
	"testing"

	"github.com/fjacquet/pdf2md/internal/layout"
)

func TestBuilder_Build(t *testing.T) {
	builder := NewBuilder()

	tests := []struct {
		name     string
		elements []layout.Element
		expected []string
	}{
		{
			name: "Basic Elements",
			elements: []layout.Element{
				{Type: layout.ElementTypeHeader, Level: 1, Content: "Title"},
				{Type: layout.ElementTypeParagraph, Content: "Paragraph 1"},
				{Type: layout.ElementTypeParagraph, Content: "Paragraph 2"},
			},
			expected: []string{
				"# Title",
				"Paragraph 1",
				"Paragraph 2",
			},
		},
		{
			name: "Lists",
			elements: []layout.Element{
				{Type: layout.ElementTypeList, Content: "Item 1"},
				{Type: layout.ElementTypeList, Content: "Item 2"},
				{Type: layout.ElementTypeList, Content: "1. Ordered Item"},
				{Type: layout.ElementTypeList, Content: "2. Another Ordered"},
			},
			expected: []string{
				"- Item 1",
				"- Item 2",
				"1. Ordered Item",
				"2. Another Ordered",
			},
		},
		{
			name: "Code Block",
			elements: []layout.Element{
				{Type: layout.ElementTypeCodeBlock, Content: "func main() {}"},
			},
			expected: []string{
				"```",
				"func main() {}",
				"```",
			},
		},
		{
			name: "Admonition",
			elements: []layout.Element{
				{Type: layout.ElementTypeAdmonition, Content: "IMPORTANT: Read this."},
				{Type: layout.ElementTypeAdmonition, Content: "Just a note."},
			},
			expected: []string{
				"> **IMPORTANT**: Read this.",
				"> Just a note.",
			},
		},
		{
			name: "Image",
			elements: []layout.Element{
				{Type: layout.ElementTypeImage, Content: "path/to/img.png"},
			},
			expected: []string{
				"![Image](path/to/img.png)",
			},
		},
		{
			name: "Table",
			elements: []layout.Element{
				{Type: layout.ElementTypeTable, Content: "| A | B |\n|---|---|\n| 1 | 2 |"},
			},
			expected: []string{
				"| A | B |",
				"|---|---|",
				"| 1 | 2 |",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := builder.Build(tt.elements)
			for _, part := range tt.expected {
				if !strings.Contains(md, part) {
					t.Errorf("Expected output to contain %q, but it didn't.\nOutput:\n%s", part, md)
				}
			}
		})
	}
}

func TestBuilder_Escape(t *testing.T) {
	builder := NewBuilder()
	input := "Hello [World] *Bold* `Code`"
	expected := "Hello \\[World\\] \\*Bold\\* \\`Code\\`"
	got := builder.Escape(input)
	if got != expected {
		t.Errorf("Escape(%q) = %q, want %q", input, got, expected)
	}
}
