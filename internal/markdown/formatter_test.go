package markdown

import (
	"strings"
	"testing"
	"text/template"

	"github.com/fjacquet/pdf2md/internal/layout"
)

func TestCustomFormatter(t *testing.T) {
	builder := NewBuilder()

	// Customize Header Template to use Setext style for H1
	// e.g. "Title\n====="
	// Note: We need a helper to repeat "=" length of content
	tmpl := template.New("header")
	tmpl.Funcs(template.FuncMap{
		"repeat": func(s string, count int) string {
			result := ""
			for i := 0; i < count; i++ {
				result += s
			}
			return result
		},
		"len": func(s string) int {
			return len(s)
		},
	})
	// If Level is 1, use Setext. Else use Atx (#)
	templateString := `{{if eq .Level 1}}{{.Content}}
{{repeat "=" (len .Content)}}
{{else}}{{repeat "#" .Level}} {{.Content}}
{{end}}
`
	var err error
	builder.Formatter.HeaderTemplate, err = tmpl.Parse(templateString)
	if err != nil {
		t.Fatalf("Failed to parse custom template: %v", err)
	}

	elements := []layout.Element{
		{Type: layout.ElementTypeHeader, Content: "My Title", Level: 1},
		{Type: layout.ElementTypeHeader, Content: "Subtitle", Level: 2},
	}

	output := builder.Build(elements)

	expectedH1 := "My Title\n========"
	expectedH2 := "## Subtitle"

	if !strings.Contains(output, expectedH1) {
		t.Errorf("Expected output to contain:\n%s\nGot:\n%s", expectedH1, output)
	}
	if !strings.Contains(output, expectedH2) {
		t.Errorf("Expected output to contain:\n%s\nGot:\n%s", expectedH2, output)
	}
}
