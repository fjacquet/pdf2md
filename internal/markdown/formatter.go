package markdown

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fjacquet/pdf2md/internal/types"
)

// Formatter handles the formatting of document elements into Markdown
type Formatter struct {
	HeaderTemplate     *template.Template
	ParagraphTemplate  *template.Template
	CodeBlockTemplate  *template.Template
	ListTemplate       *template.Template
	TableTemplate      *template.Template
	AdmonitionTemplate *template.Template
	ImageTemplate      *template.Template
}

// NewFormatter creates a new Formatter with default templates
func NewFormatter() *Formatter {
	f := &Formatter{}
	f.HeaderTemplate = parseTemplate("header", "{{repeat \"#\" .Level}} {{cleanText .Content}}\n\n")
	f.ParagraphTemplate = parseTemplate("paragraph", "{{cleanText .Content}}\n\n")
	f.CodeBlockTemplate = parseTemplate("code_block", "```\n{{.Content}}\n```\n\n")
	// List template handles ordered vs unordered
	f.ListTemplate = parseTemplate("list", "{{indent .Level}}{{listMarker .Content}} {{cleanList .Content}}\n")
	// Table template uses formatTable to ensure proper markdown table syntax
	f.TableTemplate = parseTemplate("table", "{{formatTable .Content}}\n\n")
	// Admonition template handles keyword bolding
	f.AdmonitionTemplate = parseTemplate("admonition", "{{formatAdmonition .Content}}\n\n")
	f.ImageTemplate = parseTemplate("image", "![Image]({{.Content}})\n\n")
	return f
}

func parseTemplate(name, tmpl string) *template.Template {
	t := template.New(name)
	t.Funcs(template.FuncMap{
		"cleanText": func(content string) string {
			// Remove soft hyphens (U+00AD)
			content = strings.ReplaceAll(content, "\u00AD", "")
			// Replace non-breaking hyphen with regular hyphen
			content = strings.ReplaceAll(content, "‑", "-")
			// Clean up corrupted bullet characters that might appear in text
			content = strings.ReplaceAll(content, "�", "")
			return content
		},
		"repeat": func(s string, count int) string {
			result := ""
			for i := 0; i < count; i++ {
				result += s
			}
			return result
		},
		"indent": func(level int) string {
			result := ""
			for i := 1; i < level; i++ {
				result += "  "
			}
			return result
		},
		"listMarker": func(content string) string {
			if isOrdered(content) {
				return "" // Content already has the marker
			}
			return "-"
		},
		"cleanList": func(content string) string {
			// Clean soft hyphens first
			content = strings.ReplaceAll(content, "\u00AD", "") // soft hyphen
			content = strings.ReplaceAll(content, "‑", "-")     // non-breaking hyphen to regular

			if isOrdered(content) {
				return content
			}
			// Strip existing bullets (various unicode bullet characters)
			bulletPrefixes := []string{
				"• ", "* ", "- ", "● ", "○ ", "◦ ", "▪ ", "▫ ",
				"◆ ", "◇ ", "★ ", "☆ ", "→ ", "➤ ", "► ",
				"$$ ", "� ", // corrupted/malformed bullets from PDFs
			}
			for _, prefix := range bulletPrefixes {
				if strings.HasPrefix(content, prefix) {
					return strings.TrimPrefix(content, prefix)
				}
			}
			return content
		},
		"formatAdmonition": func(content string) string {
			lines := strings.Split(content, "\n")
			var sb strings.Builder
			for i, line := range lines {
				if i == 0 {
					keywords := []string{"IMPORTANT", "WARNING", "NOTE", "TIP", "CAUTION"}
					for _, kw := range keywords {
						if strings.HasPrefix(line, kw) {
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
			return sb.String()
		},
		"formatTable": formatTable,
	})
	return template.Must(t.Parse(tmpl))
}

func isOrdered(content string) bool {
	content = strings.TrimSpace(content)
	parts := strings.SplitN(content, " ", 2)
	if len(parts) >= 2 && strings.HasSuffix(parts[0], ".") {
		prefix := strings.TrimSuffix(parts[0], ".")
		for _, c := range prefix {
			if c < '0' || c > '9' {
				return false
			}
		}
		return len(prefix) > 0
	}
	return false
}

// formatTable ensures proper markdown table syntax with separator line
func formatTable(content string) string {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return content
	}

	// Check if table already has a separator line
	hasSeparator := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if isSeparatorLine(trimmed) {
			hasSeparator = true
			break
		}
	}

	if hasSeparator {
		// Already properly formatted
		return normalizeTable(lines)
	}

	// Need to add separator after first row
	if len(lines) < 1 {
		return content
	}

	// Parse first row to determine column count
	firstRow := lines[0]
	numCols := countColumns(firstRow)
	if numCols < 1 {
		return content
	}

	// Build separator line
	separator := buildSeparator(numCols)

	// Reconstruct table with separator
	var sb strings.Builder
	sb.WriteString(normalizeRow(lines[0]))
	sb.WriteString("\n")
	sb.WriteString(separator)

	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) != "" {
			sb.WriteString("\n")
			sb.WriteString(normalizeRow(lines[i]))
		}
	}

	return sb.String()
}

// isSeparatorLine checks if a line is a markdown table separator (|---|---|)
func isSeparatorLine(line string) bool {
	if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
		return false
	}
	// Remove pipes and check if only dashes, colons, and spaces remain
	inner := strings.Trim(line, "|")
	for _, r := range inner {
		if r != '-' && r != ':' && r != ' ' && r != '|' {
			return false
		}
	}
	// Must have at least some dashes
	return strings.Contains(inner, "-")
}

// countColumns counts the number of columns in a table row
func countColumns(row string) int {
	row = strings.TrimSpace(row)
	if !strings.HasPrefix(row, "|") {
		// Not a pipe-delimited row, count by spaces or assume single column
		return 1
	}

	// Count pipes (columns = pipes - 1, but we have outer pipes)
	// | A | B | C | has 4 pipes, 3 columns
	pipes := strings.Count(row, "|")
	if pipes < 2 {
		return 1
	}
	return pipes - 1
}

// buildSeparator creates a markdown table separator line
func buildSeparator(numCols int) string {
	var sb strings.Builder
	sb.WriteString("|")
	for i := 0; i < numCols; i++ {
		sb.WriteString(" --- |")
	}
	return sb.String()
}

// normalizeRow ensures a row has proper pipe formatting
func normalizeRow(row string) string {
	row = strings.TrimSpace(row)
	if row == "" {
		return row
	}

	// Ensure leading pipe
	if !strings.HasPrefix(row, "|") {
		row = "| " + row
	}

	// Ensure trailing pipe
	if !strings.HasSuffix(row, "|") {
		row = row + " |"
	}

	return row
}

// normalizeTable normalizes all rows in a table
func normalizeTable(lines []string) string {
	var sb strings.Builder
	first := true
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !first {
			sb.WriteString("\n")
		}
		sb.WriteString(normalizeRow(trimmed))
		first = false
	}
	return sb.String()
}

// FormatElement formats a single element using the appropriate template
func (f *Formatter) FormatElement(element types.Element) (string, error) {
	var tmpl *template.Template
	switch element.Type {
	case types.ElementTypeHeader:
		tmpl = f.HeaderTemplate
	case types.ElementTypeParagraph:
		tmpl = f.ParagraphTemplate
	case types.ElementTypeCodeBlock:
		tmpl = f.CodeBlockTemplate
	case types.ElementTypeList:
		tmpl = f.ListTemplate
	case types.ElementTypeTable:
		tmpl = f.TableTemplate
	case types.ElementTypeAdmonition:
		tmpl = f.AdmonitionTemplate
	case types.ElementTypeImage:
		tmpl = f.ImageTemplate
	default:
		return element.Content + "\n\n", nil
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, element); err != nil {
		return "", err
	}
	return buf.String(), nil
}
