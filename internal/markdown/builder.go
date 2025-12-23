// Package markdown provides Markdown generation from document elements.
package markdown

import (
	"strings"

	"github.com/fjacquet/pdf2md/internal/types"
)

// Builder converts structured elements to Markdown
type Builder struct {
	// Configuration
	EnableSmartFormatting bool
	Formatter             *Formatter
}

// NewBuilder creates a new Markdown builder
func NewBuilder() *Builder {
	return &Builder{
		EnableSmartFormatting: true,
		Formatter:             NewFormatter(),
	}
}

// Build converts elements to Markdown string
func (b *Builder) Build(elements []types.Element) string {
	var sb strings.Builder

	for i, element := range elements {
		formatted, err := b.Formatter.FormatElement(element)
		if err != nil {
			// Fallback or log error? For now, just write content
			sb.WriteString(element.Content + "\n\n")
		} else {
			sb.WriteString(formatted)
		}

		// Add spacing between elements (handled by templates mostly, but we might want to ensure consistency)
		// The templates currently add \n\n or \n.
		// If we rely on templates, we don't need extra spacing logic here,
		// UNLESS the templates don't include trailing newlines.
		// My templates DO include trailing newlines (\n\n or \n).
		// So I should remove the manual spacing logic.

		// However, the original logic had conditional spacing.
		// "if element.Type == layout.ElementTypeHeader || nextElement.Type == layout.ElementTypeHeader { sb.WriteString("\n\n") }"
		// My templates have \n\n for headers.
		// Paragraphs have \n\n.
		// Lists have \n.

		// Let's rely on templates for now.
		// If we need dynamic spacing based on *next* element, templates can't do that easily.
		// But standard markdown usually is fine with \n\n everywhere except lists.
		// List items in my template have \n.
		// If the next element is NOT a list item, we might need an extra \n.

		if i < len(elements)-1 {
			nextElement := elements[i+1]
			// If current is list and next is NOT list, add extra newline to break list
			if element.Type == types.ElementTypeList && nextElement.Type != types.ElementTypeList {
				sb.WriteString("\n")
			}
		}
	}

	return sb.String()
}

// Escape escapes special Markdown characters
func (b *Builder) Escape(text string) string {
	// Characters that need escaping in Markdown
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"`", "\\`",
		"*", "\\*",
		"_", "\\_",
		"{", "\\{",
		"}", "\\}",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}
