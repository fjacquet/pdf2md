package markdown

import (
	"strings"

	"github.com/fjacquet/pdf2md/internal/layout"
)

// Builder converts structured elements to Markdown
type Builder struct {
	// Configuration
	EnableSmartFormatting bool
}

// NewBuilder creates a new Markdown builder
func NewBuilder() *Builder {
	return &Builder{
		EnableSmartFormatting: true,
	}
}

// Build converts elements to Markdown string
func (b *Builder) Build(elements []layout.Element) string {
	var sb strings.Builder

	for i, element := range elements {
		switch element.Type {
		case layout.ElementTypeHeader:
			b.writeHeader(&sb, element)
		case layout.ElementTypeParagraph:
			b.writeParagraph(&sb, element)
		case layout.ElementTypeList:
			b.writeList(&sb, element)
		case layout.ElementTypeCodeBlock:
			b.writeCodeBlock(&sb, element)
		case layout.ElementTypeTable:
			b.writeTable(&sb, element)
		}

		// Add spacing between elements
		if i < len(elements)-1 {
			nextElement := elements[i+1]
			if element.Type == layout.ElementTypeHeader || nextElement.Type == layout.ElementTypeHeader {
				sb.WriteString("\n\n")
			} else if element.Type != nextElement.Type {
				sb.WriteString("\n\n")
			} else {
				sb.WriteString("\n")
			}
		}
	}

	return sb.String()
}

// writeHeader writes a Markdown header
func (b *Builder) writeHeader(sb *strings.Builder, element layout.Element) {
	level := element.Level
	if level < 1 {
		level = 1
	}
	if level > 6 {
		level = 6
	}

	sb.WriteString(strings.Repeat("#", level))
	sb.WriteString(" ")
	sb.WriteString(strings.TrimSpace(element.Content))
}

// writeParagraph writes a Markdown paragraph
func (b *Builder) writeParagraph(sb *strings.Builder, element layout.Element) {
	content := strings.TrimSpace(element.Content)
	if content != "" {
		sb.WriteString(content)
	}
}

// writeList writes a Markdown list item
func (b *Builder) writeList(sb *strings.Builder, element layout.Element) {
	indent := strings.Repeat("  ", element.Level)
	sb.WriteString(indent)
	sb.WriteString("- ")
	sb.WriteString(strings.TrimSpace(element.Content))
}

// writeCodeBlock writes a Markdown code block
func (b *Builder) writeCodeBlock(sb *strings.Builder, element layout.Element) {
	sb.WriteString("```\n")
	sb.WriteString(element.Content)
	sb.WriteString("\n```")
}

// writeTable writes a Markdown table
func (b *Builder) writeTable(sb *strings.Builder, element layout.Element) {
	// Table formatting will be implemented in Phase 3
	sb.WriteString(element.Content)
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
