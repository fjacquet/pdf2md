package layout

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/fjacquet/pdf2md/internal/extractor"
	"github.com/fjacquet/pdf2md/internal/types"
)

func TestNewAnalyzer(t *testing.T) {
	analyzer := NewAnalyzer()
	if analyzer == nil {
		t.Fatal("NewAnalyzer() returned nil")
	}
}

func TestAnalyzeEmptyBlocks(t *testing.T) {
	analyzer := NewAnalyzer()
	doc := analyzer.Analyze(&extractor.PageContent{TextBlocks: []extractor.TextBlock{}})
	if len(doc) != 0 {
		t.Errorf("Expected 0 elements, got %d", len(doc))
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

// TestMergeElements_WordGap pins the same-line gap read as a word space:
// a 0.2 em gap (a real space is about 0.25 em) must not glue words.
func TestMergeElements_WordGap(t *testing.T) {
	tests := []struct {
		name string
		gap  float64
		want string
	}{
		{"space-sized gap", 0.2 * 10, "or ENV"},
		{"kerning-sized gap", 0.05 * 10, "orENV"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements := []Element{
				{Type: ElementTypeParagraph, Content: "or", Y: 100, X: 10, Width: 10, FontSize: 10},
				{Type: ElementTypeParagraph, Content: "ENV", Y: 100, X: 20 + tt.gap, Width: 15, FontSize: 10},
			}
			merged := NewAnalyzer().MergeElements(elements)
			if len(merged) != 1 || merged[0].Content != tt.want {
				t.Errorf("got %+v, want one element %q", merged, tt.want)
			}
		})
	}
}

// TestMergeElements_SpaceWidth: the font's own space width decides, and
// blocks without one (OCR) fall back to a 0.25 em space.
func TestMergeElements_SpaceWidth(t *testing.T) {
	tests := []struct {
		name       string
		bX, sw, fs float64
		want       string
	}{
		{"gap over half the space", 8, 2.5, 40, "A B"},
		{"gap under half the space", 6, 2.5, 40, "AB"},
		{"no space width, 0.25 em fallback", 7, 0, 10, "A B"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements := []Element{
				{Type: ElementTypeParagraph, Content: "A", Y: 100, X: 0, Width: 5, FontSize: tt.fs, SpaceWidth: tt.sw},
				{Type: ElementTypeParagraph, Content: "B", Y: 100, X: tt.bX, Width: 5, FontSize: tt.fs, SpaceWidth: tt.sw},
			}
			merged := NewAnalyzer().MergeElements(elements)
			if len(merged) != 1 || merged[0].Content != tt.want {
				t.Errorf("got %+v, want one element %q", merged, tt.want)
			}
		})
	}
}

// An image or graphic with no ID has no file to link (not saved, or trivial):
// it must not become a Markdown image.
func TestAnalyze_UnsavedGraphicsNotLinked(t *testing.T) {
	content := &extractor.PageContent{
		Images:   []types.Image{{ID: "", X: 0, Y: 100, Width: 100, Height: 100}},
		Graphics: []types.VectorGraphic{{ID: "", X: 0, Y: 300, Width: 100, Height: 100}, {ID: "images/doc_p1_2.svg", X: 0, Y: 500, Width: 100, Height: 100}},
	}
	var links []string
	for _, e := range NewAnalyzer().Analyze(content) {
		if e.Type == ElementTypeImage {
			links = append(links, e.Content)
		}
	}
	if len(links) != 1 || links[0] != "images/doc_p1_2.svg" {
		t.Errorf("image elements = %q, want only the saved graphic", links)
	}
}

// A page whose only content yields no element (an unsaved image, a graphic
// under MinImageSize) must not panic.
func TestAnalyze_NoElementsLeft(t *testing.T) {
	content := &extractor.PageContent{
		Images:   []types.Image{{ID: ""}},
		Graphics: []types.VectorGraphic{{ID: "images/x.svg", Width: 1, Height: 1}},
	}
	if got := NewAnalyzer().Analyze(content); len(got) != 0 {
		t.Errorf("Analyze() = %v, want no elements", got)
	}
}

func TestClassifyElement(t *testing.T) {
	analyzer := NewAnalyzer()
	// analyzer.ColumnGapThreshold = 10.0 // Removed
	// analyzer.HeaderSizeRatio = 1.2    // Removed
	bodySize := 12.0

	tests := []struct {
		content  string
		fontSize float64
		expected ElementType
	}{
		{"Normal paragraph", 12.0, ElementTypeParagraph},
		{"1. Introduction", 16.0, ElementTypeHeader},           // numbered + larger than body × HeaderSizeRatio
		{"**8.3 Accessing support**", 12.0, ElementTypeHeader}, // numbered + bold
		{"8.3 Accessing support", 12.0, ElementTypeHeader},     // numbered, alone: a heading
		{"apiVersion: v1", 12.0, ElementTypeCodeBlock},
		{"kind: Pod", 12.0, ElementTypeCodeBlock},
		{"- List item", 12.0, ElementTypeList},
		{"• Bullet point", 12.0, ElementTypeList},
		{"Big Header", 24.0, ElementTypeHeader},
	}

	for _, tt := range tests {
		element := Element{Content: tt.content, FontSize: tt.fontSize}
		analyzer.classifyElement(&element, nil, bodySize)
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
		{Text: "1. Introduction", X: 10, Y: 800, FontSize: 16},      // numbered and larger than body: a heading
		{Text: "This is a paragraph.", X: 10, Y: 770, FontSize: 12}, // Gap 30 > Threshold
		{Text: "Another paragraph.", X: 10, Y: 740, FontSize: 12},   // Gap 30 > Threshold
		{Text: "More body text.", X: 10, Y: 710, FontSize: 12},      // Gap 30 > Threshold
		{Text: "func main() {", X: 10, Y: 680, FontSize: 10},
	}

	// 3. Analyze layout
	content := &extractor.PageContent{
		TextBlocks: blocks,
	}
	elements := analyzer.Analyze(content)

	// 4. Print results
	for _, el := range elements {
		fmt.Printf("Type: %s, Content: %s\n", el.Type, el.Content)
	}

	// Output:
	// Type: header, Content: 1. Introduction
	// Type: paragraph, Content: This is a paragraph.
	// Type: paragraph, Content: Another paragraph.
	// Type: paragraph, Content: More body text.
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

func TestIsLikelyEquationText(t *testing.T) {
	analyzer := NewAnalyzer()

	tests := []struct {
		name string
		text string
		want bool
	}{
		{
			name: "equation with equals and number",
			text: "$z^2 = 0$    (2)",
			want: true,
		},
		{
			name: "equation with math symbols",
			text: "∫ f(x)dx = F(x) + C",
			want: true,
		},
		{
			name: "equation with Greek letters",
			text: "α + β = γ",
			want: true,
		},
		{
			name: "simple table row",
			text: "Name    Value    Description",
			want: false,
		},
		{
			name: "table with numbers",
			text: "Total    100    200    300",
			want: false,
		},
		{
			name: "equation with fraction",
			text: "a/b + c/d = e",
			want: true,
		},
		{
			name: "inline math",
			text: "$x^2$ equals $y^2$",
			want: true,
		},
		{
			name: "just text",
			text: "This is a paragraph of text",
			want: false,
		},
		{
			name: "equation number only",
			text: "x = y    (1.2)",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := analyzer.isLikelyEquationText(tt.text)
			if got != tt.want {
				t.Errorf("isLikelyEquationText(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestHasEquationNumberAtEnd(t *testing.T) {
	tests := []struct {
		text string
		want bool
	}{
		{"x = y    (1)", true},
		{"equation    (2.3)", true},
		{"something    [5]", true},
		{"no number at end", false},
		{"text (abc)", false},
		{"", false},
		{"(1)", true},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			got := hasEquationNumberAtEnd(tt.text)
			if got != tt.want {
				t.Errorf("hasEquationNumberAtEnd(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestHasMathFractionPattern(t *testing.T) {
	tests := []struct {
		text string
		want bool
	}{
		{"a/b", true},
		{"1/2", true},
		{"x/y + z", true},
		{"this/is/a/path", false}, // too long tokens
		{"no fraction here", false},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			got := hasMathFractionPattern(tt.text)
			if got != tt.want {
				t.Errorf("hasMathFractionPattern(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

// tb is a text block at y in font (regular "Helvetica" or "Helvetica-Bold").
func tb(text, font string, size, y float64) extractor.TextBlock {
	return extractor.TextBlock{Text: text, X: 72, Y: y, Width: float64(len(text)) * size * 0.5, Height: size, FontSize: size, FontName: font}
}

// typeOf analyses a page whose body is 10 pt regular text and returns the
// element types in order, with their content.
func typeOf(t *testing.T, blocks ...extractor.TextBlock) []string {
	t.Helper()
	body := []extractor.TextBlock{
		tb("Body text that sets the body font size of the page.", "Helvetica", 10, 100),
		tb("More body text that sets the body font size.", "Helvetica", 10, 88),
		tb("Even more body text for the body font size.", "Helvetica", 10, 76),
	}
	var out []string
	for _, e := range NewAnalyzer().Analyze(&extractor.PageContent{TextBlocks: append(blocks, body...)}) {
		out = append(out, fmt.Sprintf("%s:%s", e.Type, e.Content))
	}
	return out
}

func isHeader(got []string, text string) bool {
	return slices.Contains(got, fmt.Sprintf("%s:%s", ElementTypeHeader, text))
}

func TestHeading_NonLetterStart(t *testing.T) {
	if got := typeOf(t, tb(".NET MAUI Learning Resources", "Helvetica-Bold", 14, 700)); !isHeader(got, ".NET MAUI Learning Resources") {
		t.Errorf("got %q, want a header", got)
	}
}

func TestHeading_LowercaseStartRejected(t *testing.T) {
	if got := typeOf(t, tb("and then some", "Helvetica", 14, 700)); isHeader(got, "and then some") {
		t.Errorf("got %q, want no header", got)
	}
}

func TestHeading_NumberedBodyLineIsList(t *testing.T) {
	got := typeOf(t, tb("1. Wash", "Helvetica", 10, 700), tb("2. Slice", "Helvetica", 10, 688))
	for _, e := range got {
		if strings.HasPrefix(e, string(ElementTypeHeader)+":") {
			t.Errorf("got %q, want the numbered body lines as list items", got)
		}
	}
}

func TestHeading_NumberedStyledLineIsHeading(t *testing.T) {
	got := typeOf(t, tb("2. Results", "Helvetica", 14, 700), tb("2.1 Setup", "Helvetica-Bold", 10, 650))
	if !isHeader(got, "2. Results") || !isHeader(got, "2.1 Setup") {
		t.Errorf("got %q, want both numbered lines as headers", got)
	}
}

func TestHeading_NumberedMidSentence(t *testing.T) {
	got := typeOf(t, tb("as shown in 2. The result follows.", "Helvetica", 10, 700))
	if isHeader(got, "as shown in 2. The result follows.") {
		t.Errorf("got %q, want no header", got)
	}
}

func TestHeading_BoldRunIsParagraph(t *testing.T) {
	got := typeOf(t,
		tb("Abstract — We address the error floor", "Helvetica-Bold", 10, 700),
		tb("problem of low-density parity check codes", "Helvetica-Bold", 10, 688),
		tb("on the binary input Gaussian channel", "Helvetica-Bold", 10, 676))
	for _, e := range got {
		if strings.HasPrefix(e, string(ElementTypeHeader)+":") {
			t.Errorf("got %q, want the bold block as a paragraph", got)
		}
	}
}

func TestHeading_SingleBoldLineIsHeading(t *testing.T) {
	if got := typeOf(t, tb("Installation", "Helvetica-Bold", 10, 700)); !isHeader(got, "Installation") {
		t.Errorf("got %q, want a header", got)
	}
}

func TestHeading_BoldSentenceIsParagraph(t *testing.T) {
	if got := typeOf(t, tb("This is important.", "Helvetica-Bold", 10, 700)); isHeader(got, "This is important.") {
		t.Errorf("got %q, want no header", got)
	}
}

// A numbered section title in body style (1004.3799) stands alone: a heading.
func TestHeading_NumberedBodyLineAloneIsHeading(t *testing.T) {
	if got := typeOf(t, tb("1. Sums of Like Powers", "Helvetica", 10, 700)); !isHeader(got, "1. Sums of Like Powers") {
		t.Errorf("got %q, want a header", got)
	}
}

// A section opening straight into its subsection is not a list: list items
// share one numbering depth.
func TestHeading_SectionThenSubsection(t *testing.T) {
	got := typeOf(t,
		tb("2. Methods", "Helvetica", 10, 700),
		tb("2.1. Data Collection", "Helvetica", 10, 688),
		tb("We collected data from the survey.", "Helvetica", 10, 670))
	if !isHeader(got, "2. Methods") || !isHeader(got, "2.1. Data Collection") {
		t.Errorf("got %q, want both numbered lines as headers", got)
	}
}

func TestHeading_BoldSectionThenSubsection(t *testing.T) {
	got := typeOf(t,
		tb("1. Introduction", "Helvetica-Bold", 10, 700),
		tb("1.1. Background", "Helvetica-Bold", 10, 684),
		tb("We start from earlier work.", "Helvetica", 10, 666))
	if !isHeader(got, "1. Introduction") || !isHeader(got, "1.1. Background") {
		t.Errorf("got %q, want both bold numbered lines as headers", got)
	}
}

// A paragraph line that happens to start with "2." continues the sentence.
func TestHeading_NumberedLineWrapsInParagraph(t *testing.T) {
	got := typeOf(t,
		tb("The outcome is as shown in", "Helvetica", 10, 700),
		tb("2. The result was positive and", "Helvetica", 10, 688),
		tb("the experiment ended after that step.", "Helvetica", 10, 676))
	if isHeader(got, "2. The result was positive and") {
		t.Errorf("got %q, want no header", got)
	}
}

// An ornament has no word to title anything, and ranking would give it #.
func TestHeading_LetterlessLargeTextIsNotHeading(t *testing.T) {
	if got := typeOf(t, tb("· · ·", "Helvetica", 24, 700)); isHeader(got, "· · ·") {
		t.Errorf("got %q, want no header", got)
	}
}
