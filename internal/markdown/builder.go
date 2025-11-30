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
		case layout.ElementTypeAdmonition:
			b.writeAdmonition(&sb, element)
		case layout.ElementTypeImage:
			b.writeImage(&sb, element)
		}

		// Add spacing between elements
		if i < len(elements)-1 {
			nextElement := elements[i+1]
			if element.Type == layout.ElementTypeHeader || nextElement.Type == layout.ElementTypeHeader {
				sb.WriteString("\n\n")
			} else if element.Type == layout.ElementTypeParagraph && nextElement.Type == layout.ElementTypeParagraph {
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

	content := strings.TrimSpace(element.Content)

	// Check if content already starts with a list marker
	isOrdered := false
	if len(content) > 0 {
		// Check for numbered list pattern (digits + dot + space)
		parts := strings.SplitN(content, " ", 2)
		if len(parts) >= 2 && strings.HasSuffix(parts[0], ".") {
			// Check if prefix is a number
			prefix := strings.TrimSuffix(parts[0], ".")
			isAllDigits := true
			for _, c := range prefix {
				if c < '0' || c > '9' {
					isAllDigits = false
					break
				}
			}
			if isAllDigits && len(prefix) > 0 {
				isOrdered = true
			}
		}
	}

	if !isOrdered {
		// If it's not an ordered list, add a bullet point
		// Also strip existing bullets if any (to normalize)
		if strings.HasPrefix(content, "• ") || strings.HasPrefix(content, "* ") || strings.HasPrefix(content, "- ") {
			content = content[2:]
		}
		sb.WriteString("- ")
	}

	sb.WriteString(content)
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

// writeAdmonition writes a Markdown admonition (blockquote)
func (b *Builder) writeAdmonition(sb *strings.Builder, element layout.Element) {
	// Format: > **KEYWORD**
	//         > Content...

	content := strings.TrimSpace(element.Content)
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if i == 0 {
			// Check if the first line starts with the keyword
			keywords := []string{"IMPORTANT", "WARNING", "NOTE", "TIP", "CAUTION"}
			for _, kw := range keywords {
				if strings.HasPrefix(line, kw) {
					// Bold the keyword
					line = "**" + kw + "**" + strings.TrimPrefix(line, kw)
					break
				}
			}
		}
		sb.WriteString("> " + line)
		if i < len(lines)-1 {
			sb.WriteString("\n")
		}
	}
}

// writeImage writes a Markdown image
func (b *Builder) writeImage(sb *strings.Builder, element layout.Element) {
	// Format: ![Image](path/to/image.png)
	// Content holds the image path
	sb.WriteString("![Image](" + element.Content + ")")
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
