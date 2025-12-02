package markdown

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fjacquet/pdf2md/internal/layout"
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
	f.TableTemplate = parseTemplate("table", "{{.Content}}\n\n")
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

// FormatElement formats a single element using the appropriate template
func (f *Formatter) FormatElement(element layout.Element) (string, error) {
	var tmpl *template.Template
	switch element.Type {
	case layout.ElementTypeHeader:
		tmpl = f.HeaderTemplate
	case layout.ElementTypeParagraph:
		tmpl = f.ParagraphTemplate
	case layout.ElementTypeCodeBlock:
		tmpl = f.CodeBlockTemplate
	case layout.ElementTypeList:
		tmpl = f.ListTemplate
	case layout.ElementTypeTable:
		tmpl = f.TableTemplate
	case layout.ElementTypeAdmonition:
		tmpl = f.AdmonitionTemplate
	case layout.ElementTypeImage:
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
