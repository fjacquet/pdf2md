package markdown

import (
	"strings"
	"testing"

	"github.com/fjacquet/pdf2md/internal/layout"
)

func TestBuilder(t *testing.T) {
	builder := NewBuilder()

	elements := []layout.Element{
		{Type: layout.ElementTypeHeader, Level: 1, Content: "Title"},
		{Type: layout.ElementTypeParagraph, Content: "This is a paragraph."},
		{Type: layout.ElementTypeHeader, Level: 2, Content: "Section"},
		{Type: layout.ElementTypeList, Content: "Item 1"},
		{Type: layout.ElementTypeList, Content: "Item 2"},
		{Type: layout.ElementTypeCodeBlock, Content: "func main() {}"},
	}

	md := builder.Build(elements)

	expectedParts := []string{
		"# Title",
		"This is a paragraph.",
		"## Section",
		"- Item 1",
		"- Item 2",
		"```",
		"func main() {}",
		"```",
	}

	for _, part := range expectedParts {
		if !strings.Contains(md, part) {
			t.Errorf("Expected output to contain %q, but it didn't.\nOutput:\n%s", part, md)
		}
	}
}
