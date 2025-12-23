package layout

import (
	"strings"
	"testing"

	"github.com/fjacquet/pdf2md/internal/extractor"
)

func TestTableDetector_GroupByY(t *testing.T) {
	td := NewTableDetector()

	blocks := []extractor.TextBlock{
		{Text: "A1", X: 10, Y: 100},
		{Text: "B1", X: 50, Y: 100},
		{Text: "C1", X: 90, Y: 100},
		{Text: "A2", X: 10, Y: 80},
		{Text: "B2", X: 50, Y: 80},
		{Text: "C2", X: 90, Y: 80},
	}

	rows := td.groupByY(blocks)

	if len(rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(rows))
	}

	// Each row should have 3 blocks
	for i, row := range rows {
		if len(row) != 3 {
			t.Errorf("row %d: expected 3 blocks, got %d", i, len(row))
		}
	}
}

func TestTableDetector_FindColumnPositions(t *testing.T) {
	td := NewTableDetector()

	blocks := []extractor.TextBlock{
		{Text: "A1", X: 10, Y: 100},
		{Text: "B1", X: 50, Y: 100},
		{Text: "C1", X: 90, Y: 100},
		{Text: "A2", X: 10, Y: 80},
		{Text: "B2", X: 50, Y: 80},
		{Text: "C2", X: 90, Y: 80},
	}

	rows := td.groupByY(blocks)
	positions := td.findColumnPositionsFromBlocks(blocks, rows)

	if len(positions) != 3 {
		t.Errorf("expected 3 column positions, got %d", len(positions))
	}
}

func TestTableDetector_DetectTables(t *testing.T) {
	td := NewTableDetector()

	// Create a simple 3x3 table
	blocks := []extractor.TextBlock{
		{Text: "Header A", X: 10, Y: 100, Width: 30, Height: 12},
		{Text: "Header B", X: 60, Y: 100, Width: 30, Height: 12},
		{Text: "Header C", X: 110, Y: 100, Width: 30, Height: 12},
		{Text: "Row1 A", X: 10, Y: 85, Width: 30, Height: 12},
		{Text: "Row1 B", X: 60, Y: 85, Width: 30, Height: 12},
		{Text: "Row1 C", X: 110, Y: 85, Width: 30, Height: 12},
		{Text: "Row2 A", X: 10, Y: 70, Width: 30, Height: 12},
		{Text: "Row2 B", X: 60, Y: 70, Width: 30, Height: 12},
		{Text: "Row2 C", X: 110, Y: 70, Width: 30, Height: 12},
	}

	tables, consumed := td.DetectTables(blocks, nil)

	if len(tables) != 1 {
		t.Errorf("expected 1 table, got %d", len(tables))
		return
	}

	table := tables[0]
	if table.NumRows != 3 {
		t.Errorf("expected 3 rows, got %d", table.NumRows)
	}
	if table.NumCols != 3 {
		t.Errorf("expected 3 cols, got %d", table.NumCols)
	}

	// All blocks should be consumed
	if len(consumed) != 9 {
		t.Errorf("expected 9 consumed blocks, got %d", len(consumed))
	}
}

func TestTableStructure_ToMarkdown(t *testing.T) {
	table := &TableStructure{
		Cells: [][]TableCell{
			{{Content: "A"}, {Content: "B"}, {Content: "C"}},
			{{Content: "1"}, {Content: "2"}, {Content: "3"}},
		},
		NumRows:   2,
		NumCols:   3,
		HasHeader: true,
	}

	md := table.ToMarkdown()

	// Should have header, separator, and data row
	lines := strings.Split(md, "\n")
	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines, got %d", len(lines))
	}

	// Check separator line exists
	hasSeparator := false
	for _, line := range lines {
		if strings.Contains(line, "---") {
			hasSeparator = true
			break
		}
	}
	if !hasSeparator {
		t.Error("expected separator line with ---")
	}

	// Check content
	if !strings.Contains(md, "A") || !strings.Contains(md, "B") {
		t.Error("expected header content A and B")
	}
}

func TestTableStructure_ToMarkdown_NoHeader(t *testing.T) {
	table := &TableStructure{
		Cells: [][]TableCell{
			{{Content: "1"}, {Content: "2"}},
			{{Content: "3"}, {Content: "4"}},
		},
		NumRows:   2,
		NumCols:   2,
		HasHeader: false,
	}

	md := table.ToMarkdown()

	// Should still have separator (markdown requirement)
	if !strings.Contains(md, "---") {
		t.Error("expected separator line even without header")
	}
}

func TestIsLikelyEquation(t *testing.T) {
	tests := []struct {
		name   string
		blocks []extractor.TextBlock
		want   bool
	}{
		{
			name: "equation with math symbols",
			blocks: []extractor.TextBlock{
				{Text: "x²", FontName: "CMMI10"},
				{Text: "="},
				{Text: "y", FontName: "CMMI10"},
				{Text: "(1)"},
			},
			want: true,
		},
		{
			name: "simple text table",
			blocks: []extractor.TextBlock{
				{Text: "Name", FontName: "Arial"},
				{Text: "Age", FontName: "Arial"},
				{Text: "John", FontName: "Arial"},
				{Text: "25", FontName: "Arial"},
			},
			want: false,
		},
		{
			name: "equation with integral",
			blocks: []extractor.TextBlock{
				{Text: "∫", FontName: "Symbol"},
				{Text: "f(x)dx", FontName: "CMMI10"},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsLikelyEquation(tt.blocks)
			if got != tt.want {
				t.Errorf("IsLikelyEquation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEquationNumber(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"(1)", true},
		{"(2.3)", true},
		{"(1.2.3)", true},
		{"[1]", true},
		{"[42]", true},
		{"(a)", false},
		{"()", false},
		{"1", false},
		{"(1", false},
		{"1)", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := isEquationNumber(tt.input)
			if got != tt.want {
				t.Errorf("isEquationNumber(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsNumericContent(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"123", true},
		{"45.67", true},
		{"$100", true},
		{"50%", true},
		{"-42", true},
		{"Hello", false},
		{"ABC123", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := isNumericContent(tt.input)
			if got != tt.want {
				t.Errorf("isNumericContent(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
