package layout

import (
	"fmt"
	"testing"

	"github.com/fjacquet/pdf2md/internal/extractor"
)

func TestNewAnalyzer(t *testing.T) {
	analyzer := NewAnalyzer()
	if analyzer == nil {
		t.Fatal("NewAnalyzer() returned nil")
	}

	if analyzer.ColumnGapThreshold <= 0 {
		t.Error("ColumnGapThreshold should be positive")
	}

	if analyzer.HeaderSizeRatio <= 1.0 {
		t.Error("HeaderSizeRatio should be > 1.0")
	}
}

func TestAnalyzeEmptyBlocks(t *testing.T) {
	analyzer := NewAnalyzer()
	elements := analyzer.Analyze(nil)
	if elements != nil {
		t.Error("Expected nil for empty blocks")
	}

	elements = analyzer.Analyze([]extractor.TextBlock{})
	if elements != nil {
		t.Error("Expected nil for empty blocks")
	}
}

func TestDetectColumns(t *testing.T) {
	analyzer := NewAnalyzer()
	analyzer.ColumnGapThreshold = 50.0

	// Single column
	blocks := []extractor.TextBlock{
		{X: 10, Y: 100, Text: "Left"},
		{X: 12, Y: 200, Text: "Left"},
		{X: 15, Y: 300, Text: "Left"},
	}

	columns := analyzer.detectColumns(blocks)
	if len(columns) != 1 {
		t.Errorf("Expected 1 column, got %d", len(columns))
	}

	// Two columns
	blocks = []extractor.TextBlock{
		{X: 10, Y: 100, Text: "Left"},
		{X: 12, Y: 200, Text: "Left"},
		{X: 300, Y: 100, Text: "Right"},
		{X: 305, Y: 200, Text: "Right"},
	}

	columns = analyzer.detectColumns(blocks)
	if len(columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(columns))
	}
}

func TestDetectBodyFontSize(t *testing.T) {
	analyzer := NewAnalyzer()

	blocks := []extractor.TextBlock{
		{FontSize: 12.0},
		{FontSize: 12.0},
		{FontSize: 12.0},
		{FontSize: 18.0}, // Header
		{FontSize: 12.0},
	}

	bodySize := analyzer.detectBodyFontSize(blocks)
	if bodySize != 12.0 {
		t.Errorf("Expected body size 12.0, got %.2f", bodySize)
	}
}

func TestCalculateHeaderLevel(t *testing.T) {
	analyzer := NewAnalyzer()

	tests := []struct {
		fontSize float64
		bodySize float64
		expected int
	}{
		{24.0, 12.0, 1}, // 2x = H1
		{21.0, 12.0, 2}, // 1.75x = H2
		{18.0, 12.0, 3}, // 1.5x = H3
		{15.6, 12.0, 4}, // 1.3x = H4
		{13.8, 12.0, 5}, // 1.15x = H5
		{12.0, 12.0, 6}, // 1.0x = H6
	}

	for _, tt := range tests {
		level := analyzer.calculateHeaderLevel(tt.fontSize, tt.bodySize)
		if level != tt.expected {
			t.Errorf("Font %.1f / body %.1f: expected H%d, got H%d",
				tt.fontSize, tt.bodySize, tt.expected, level)
		}
	}
}

func TestMergeElements(t *testing.T) {
	analyzer := NewAnalyzer()

	elements := []Element{
		{Type: ElementTypeParagraph, Content: "Hello", Y: 100, X: 10, Width: 40, FontSize: 12},
		{Type: ElementTypeParagraph, Content: "World", Y: 100.5, X: 60, Width: 40, FontSize: 12}, // Gap = 60 - (10+40) = 10. Threshold ~2.4. Gap > Threshold -> Space
		{Type: ElementTypeParagraph, Content: "Next", Y: 150, X: 10, Width: 40, FontSize: 12},    // Different line
	}

	merged := analyzer.MergeElements(elements)

	if len(merged) != 2 {
		t.Errorf("Expected 2 merged elements, got %d", len(merged))
	}

	if merged[0].Content != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", merged[0].Content)
	}
}

func TestClassifyElement(t *testing.T) {
	analyzer := NewAnalyzer()
	bodySize := 12.0

	tests := []struct {
		content  string
		fontSize float64
		expected ElementType
	}{
		{"Normal paragraph", 12.0, ElementTypeParagraph},
		{"1. Introduction", 12.0, ElementTypeHeader},
		{"8.3 Accessing support", 12.0, ElementTypeHeader},
		{"apiVersion: v1", 12.0, ElementTypeCodeBlock},
		{"kind: Pod", 12.0, ElementTypeCodeBlock},
		{"- List item", 12.0, ElementTypeList},
		{"• Bullet point", 12.0, ElementTypeList},
		{"Big Header", 24.0, ElementTypeHeader},
	}

	for _, tt := range tests {
		element := Element{Content: tt.content, FontSize: tt.fontSize}
		analyzer.classifyElement(&element, bodySize)
		if element.Type != tt.expected {
			t.Errorf("classifyElement(%q) type = %v, want %v", tt.content, element.Type, tt.expected)
		}
	}
}

func ExampleAnalyzer_Analyze() {
	// 1. Create an Analyzer
	analyzer := NewAnalyzer()

	// 2. Define raw text blocks (simulating extraction)
	blocks := []extractor.TextBlock{
		{Text: "1. Introduction", X: 10, Y: 800, FontSize: 12},
		{Text: "This is a paragraph.", X: 10, Y: 780, FontSize: 12},
		{Text: "func main() {", X: 10, Y: 760, FontSize: 10},
	}

	// 3. Analyze layout
	elements := analyzer.Analyze(blocks)

	// 4. Print results
	for _, el := range elements {
		fmt.Printf("Type: %s, Content: %s\n", el.Type, el.Content)
	}

	// Output:
	// Type: header, Content: 1. Introduction
	// Type: paragraph, Content: This is a paragraph.
	// Type: code_block, Content: func main() {
}

func TestMergeCodeBlocks(t *testing.T) {
	analyzer := NewAnalyzer()

	elements := []Element{
		{Type: ElementTypeCodeBlock, Content: "line 1"},
		{Type: ElementTypeCodeBlock, Content: "line 2"},
		{Type: ElementTypeParagraph, Content: "text"},
		{Type: ElementTypeCodeBlock, Content: "line 3"},
	}

	merged := analyzer.MergeCodeBlocks(elements)

	if len(merged) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(merged))
	}

	if merged[0].Type != ElementTypeCodeBlock || merged[0].Content != "line 1\nline 2" {
		t.Errorf("First block mismatch: %v", merged[0])
	}

	if merged[1].Type != ElementTypeParagraph {
		t.Errorf("Second block mismatch: %v", merged[1])
	}

	if merged[2].Type != ElementTypeCodeBlock || merged[2].Content != "line 3" {
		t.Errorf("Third block mismatch: %v", merged[2])
	}
}
