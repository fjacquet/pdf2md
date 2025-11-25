package layout

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/fjacquet/pdf2md/internal/extractor"
)

// Element represents a logical document element (header, paragraph, etc.)
type Element struct {
	Type     ElementType
	Content  string
	Level    int     // For headers: 1-6, for lists: indent level
	X        float64 // Bounding box
	Y        float64
	Width    float64
	Height   float64
	FontSize float64
}

// ElementType defines the type of a document element
type ElementType string

const (
	ElementTypeHeader    ElementType = "header"
	ElementTypeParagraph ElementType = "paragraph"
	ElementTypeCodeBlock ElementType = "code_block"
	ElementTypeList      ElementType = "list"
	ElementTypeTable     ElementType = "table"
)

// Rule defines a classification rule for the layout analyzer
type Rule struct {
	Name      string
	Condition func(text string, fontSize float64) bool
	Type      ElementType
}

// Analyzer analyzes the layout of text blocks
type Analyzer struct {
	ColumnGapThreshold float64
	HeaderSizeRatio    float64
	Rules              []Rule
}

// NewAnalyzer creates a new Analyzer with default configuration
func NewAnalyzer() *Analyzer {
	a := &Analyzer{
		ColumnGapThreshold: 10.0, // Default gap to consider as new column
		HeaderSizeRatio:    1.2,  // Header is 1.2x larger than body
	}

	// Initialize Rules
	a.Rules = []Rule{
		{
			Name: "Code Block (Keywords)",
			Condition: func(text string, fontSize float64) bool {
				return a.isCodeBlock(text)
			},
			Type: ElementTypeCodeBlock,
		},
		{
			Name: "Table Row (Wide Gaps)",
			Condition: func(text string, fontSize float64) bool {
				// Heuristic: Line contains wide gaps (4 spaces)
				// Require at least 1 wide gap (2 columns)
				return strings.Count(text, "    ") >= 1
			},
			Type: ElementTypeTable,
		},
		{
			Name: "Header (Numbered)",
			Condition: func(text string, fontSize float64) bool {
				return a.isNumberedHeader(text)
			},
			Type: ElementTypeHeader,
		},
		{
			Name: "List Item (Bullet/Number)",
			Condition: func(text string, fontSize float64) bool {
				return a.isListItem(text)
			},
			Type: ElementTypeList,
		},
	}

	return a
}

// Analyze converts raw text blocks into structured elements
func (a *Analyzer) Analyze(blocks []extractor.TextBlock) []Element {
	if len(blocks) == 0 {
		return nil
	}

	// Step 1: Detect columns by analyzing X positions
	columns := a.detectColumns(blocks)

	// Step 2: Reorder blocks by reading order (column-aware)
	orderedBlocks := a.reorderByReadingOrder(blocks, columns)

	// Step 3: Detect body font size (most common)
	bodyFontSize := a.detectBodyFontSize(orderedBlocks)

	// Step 4: Convert blocks to initial elements (all paragraphs initially)
	elements := a.blocksToElements(orderedBlocks)

	// Step 5: Merge consecutive elements on same line
	elements = a.MergeElements(elements)

	// Step 6: Classify elements (Header, CodeBlock, etc.)
	// We do this after merging lines so we have the full context
	for i := range elements {
		a.classifyElement(&elements[i], bodyFontSize)
	}

	// Step 7: Merge consecutive code blocks (vertical)
	elements = a.MergeCodeBlocks(elements)

	// Step 8: Merge consecutive table rows
	elements = a.MergeTableRows(elements)

	return elements
}

// detectColumns identifies column boundaries by analyzing X positions
func (a *Analyzer) detectColumns(blocks []extractor.TextBlock) []float64 {
	if len(blocks) == 0 {
		return []float64{0}
	}

	// Collect all X positions
	xPositions := make([]float64, 0, len(blocks))
	for _, block := range blocks {
		xPositions = append(xPositions, block.X)
	}

	// Sort X positions
	sort.Float64s(xPositions)

	// Find gaps larger than threshold
	columns := []float64{0} // Start of first column
	for i := 1; i < len(xPositions); i++ {
		gap := xPositions[i] - xPositions[i-1]
		if gap > a.ColumnGapThreshold {
			// Found a column boundary
			columns = append(columns, xPositions[i])
		}
	}

	return columns
}

// reorderByReadingOrder sorts blocks by reading order (top-to-bottom, left-to-right, column-aware)
func (a *Analyzer) reorderByReadingOrder(blocks []extractor.TextBlock, columns []float64) []extractor.TextBlock {
	// Assign each block to a column
	type blockWithColumn struct {
		block  extractor.TextBlock
		column int
	}

	blocksWithCols := make([]blockWithColumn, 0, len(blocks))
	for _, block := range blocks {
		col := a.getColumnIndex(block.X, columns)
		blocksWithCols = append(blocksWithCols, blockWithColumn{block, col})
	}

	// Sort by: column (left to right), then Y (top to bottom), then X (left to right)
	sort.SliceStable(blocksWithCols, func(i, j int) bool {
		if blocksWithCols[i].column != blocksWithCols[j].column {
			return blocksWithCols[i].column < blocksWithCols[j].column
		}

		// Check if on same line (within threshold)
		yDiff := math.Abs(blocksWithCols[i].block.Y - blocksWithCols[j].block.Y)
		if yDiff < 2.0 {
			// If X is significantly different, sort by X
			if math.Abs(blocksWithCols[i].block.X-blocksWithCols[j].block.X) > 1.0 {
				return blocksWithCols[i].block.X < blocksWithCols[j].block.X
			}
			// If X is identical (or very close), preserve original order (return false)
			// Since we use SliceStable, returning false for equal items preserves order.
			return false
		}

		return blocksWithCols[i].block.Y > blocksWithCols[j].block.Y
	})

	// Extract sorted blocks
	result := make([]extractor.TextBlock, len(blocksWithCols))
	for i, bwc := range blocksWithCols {
		result[i] = bwc.block
	}

	return result
}

// getColumnIndex determines which column a given X position belongs to
func (a *Analyzer) getColumnIndex(x float64, columns []float64) int {
	for i := len(columns) - 1; i >= 0; i-- {
		if x >= columns[i] {
			return i
		}
	}
	return 0
}

// detectBodyFontSize finds the most common font size (assumed to be body text)
func (a *Analyzer) detectBodyFontSize(blocks []extractor.TextBlock) float64 {
	if len(blocks) == 0 {
		return 12.0 // default
	}

	// Count font size occurrences
	fontSizes := make(map[float64]int)
	for _, block := range blocks {
		if block.FontSize > 0 {
			fontSizes[block.FontSize]++
		}
	}

	// Find most common
	var maxCount int
	var bodySize float64 = 12.0
	for size, count := range fontSizes {
		if count > maxCount {
			maxCount = count
			bodySize = size
		}
	}

	return bodySize
}

// blocksToElements converts text blocks to initial elements
func (a *Analyzer) blocksToElements(blocks []extractor.TextBlock) []Element {
	elements := make([]Element, 0, len(blocks))

	for _, block := range blocks {
		element := Element{
			Content:  block.Text,
			X:        block.X,
			Y:        block.Y,
			Width:    block.Width,
			Height:   block.Height,
			FontSize: block.FontSize,
			Type:     ElementTypeParagraph, // Default
		}
		elements = append(elements, element)
	}

	return elements
}

// classifyElement determines the type of an element based on its content and properties
func (a *Analyzer) classifyElement(element *Element, bodyFontSize float64) {
	text := strings.TrimSpace(element.Content)
	if text == "" {
		return
	}

	// 1. Check Rules
	for _, rule := range a.Rules {
		if rule.Condition(text, element.FontSize) {
			element.Type = rule.Type
			return
		}
	}

	// 2. Fallback: Font Size for Headers (if available)
	if element.FontSize > bodyFontSize*a.HeaderSizeRatio {
		element.Type = ElementTypeHeader
		return
	}

	// Default is Paragraph
	element.Type = ElementTypeParagraph
}

// isNumberedHeader checks if the text looks like a numbered header
func (a *Analyzer) isNumberedHeader(text string) bool {
	parts := strings.SplitN(text, " ", 2)
	if len(parts) < 2 {
		return false
	} // Must have text after number

	// Check if first part is a number like "1." or "1.2" or "1.2.3"
	marker := parts[0]
	if len(marker) == 0 {
		return false
	}

	// Must start with digit
	if marker[0] < '0' || marker[0] > '9' {
		return false
	}

	validChars := "0123456789."
	for _, c := range marker {
		if !strings.ContainsRune(validChars, c) {
			return false
		}
	}

	return true
}

// isCodeBlock checks if the text looks like code
func (a *Analyzer) isCodeBlock(text string) bool {
	// Stricter checks to avoid false positives in English text

	// 1. Check for specific prefixes (start of line)
	prefixes := []string{
		"func ", "package ", "import ",
		"const ", "type ", "struct ", "interface ",
		"return ", "if ", "else ", "for ", "range ",
		"$ ", "# ", // Shell prompts
		"kubectl ", "oc ", "docker ", "podman ", "sudo ",
		"apiVersion:", "kind:", "metadata:", "spec:", "status:",
	}

	for _, p := range prefixes {
		if strings.HasPrefix(text, p) {
			// fmt.Printf("DEBUG: Found code block (prefix): '%s' (prefix: '%s')\n", text, p)
			return true
		}
	}

	// 2. Check for unique substrings anywhere (very specific)
	containsKeywords := []string{
		"apiVersion:", "kind: Pod", "kind: Service", "kind: Deployment",
		"namespace: ",
	}

	for _, kw := range containsKeywords {
		if strings.Contains(text, kw) {
			// fmt.Printf("DEBUG: Found code block (contains): '%s' (keyword: '%s')\n", text, kw)
			return true
		}
	}

	// 3. Check for exact matches of brackets (unlikely in normal text unless math)
	if text == "{" || text == "}" || text == "[]" || text == "()" {
		return true
	}

	return false
}

// isListItem checks if the block starts with a list marker
func (a *Analyzer) isListItem(text string) bool {
	if len(text) == 0 {
		return false
	}

	// Check for bullet points
	if strings.HasPrefix(text, "•") || strings.HasPrefix(text, "- ") || strings.HasPrefix(text, "* ") {
		return true
	}

	// Check for numbered lists (1., 2., etc.)
	// Simple check for digit + dot + space
	if len(text) >= 3 && text[0] >= '0' && text[0] <= '9' && text[1] == '.' && text[2] == ' ' {
		return true
	}

	return false
}

// calculateHeaderLevel determines header level (1-6) based on font size
func (a *Analyzer) calculateHeaderLevel(fontSize, bodyFontSize float64) int {
	ratio := fontSize / bodyFontSize

	// Map font size ratios to header levels
	// H1: 2x or more
	// H2: 1.75x
	// H3: 1.5x
	// H4: 1.3x
	// H5: 1.15x
	// H6: 1.0x (same as body, but might be bold)

	if ratio >= 2.0 {
		return 1
	} else if ratio >= 1.75 {
		return 2
	} else if ratio >= 1.5 {
		return 3
	} else if ratio >= 1.3 {
		return 4
	} else if ratio >= 1.15 {
		return 5
	}
	return 6
}

// MergeElements merges consecutive text elements on the same line
func (a *Analyzer) MergeElements(elements []Element) []Element {
	if len(elements) <= 1 {
		return elements
	}

	merged := make([]Element, 0, len(elements))
	current := elements[0]

	for i := 1; i < len(elements); i++ {
		next := elements[i]

		// Check if elements are on the same line (similar Y coordinate)
		yDiff := math.Abs(current.Y - next.Y)
		// We only merge if they are on the same line.
		// Since we classify AFTER merging, we don't need to check Type here (they are all Paragraphs or default)
		if yDiff < 2.0 {
			// Calculate gap between elements
			// Note: Width might be 0 if not provided by extractor, so be careful
			gap := next.X - (current.X + current.Width)

			// If gap is small (relative to font size or absolute), don't add space
			// Using a heuristic here: if gap is less than 20% of font size, assume it's part of the same word
			// Or if gap is negative (overlap)
			threshold := current.FontSize * 0.2
			if threshold == 0 {
				threshold = 2.0 // Fallback
			}

			// Determine threshold for "wide gap"
			// If FontSize is available, use relative. Otherwise use absolute.
			wideGapThreshold := 30.0 // Increased from 20.0
			if current.FontSize > 0 {
				wideGapThreshold = current.FontSize * 3.0
			}

			// Debug specific text to understand gap
			if strings.Contains(current.Content, "Component") {
				fmt.Printf("DEBUG: Merging '%s' (X=%.2f, W=%.2f) with '%s' (X=%.2f). Gap=%.2f. Threshold=%.2f, Wide=%.2f\n",
					current.Content, current.X, current.Width, next.Content, next.X, gap, threshold, wideGapThreshold)
			}

			if gap < threshold {
				current.Content += next.Content
			} else if gap > wideGapThreshold {
				// Wide gap, likely a table column or visual separation
				current.Content += "    " + next.Content
			} else {
				current.Content += " " + next.Content
			}

			// Update width
			current.Width = next.X + next.Width - current.X
		} else {
			// Save current and move to next
			merged = append(merged, current)
			current = next
		}
	}

	// Add the last element
	merged = append(merged, current)

	return merged
}

// MergeCodeBlocks merges consecutive code block elements into a single block
func (a *Analyzer) MergeCodeBlocks(elements []Element) []Element {
	if len(elements) <= 1 {
		return elements
	}

	merged := make([]Element, 0, len(elements))
	var current *Element

	for i := 0; i < len(elements); i++ {
		element := elements[i]

		if element.Type == ElementTypeCodeBlock {
			if current == nil {
				// Start new code block
				// Make a copy to avoid modifying original if needed
				newElem := element
				current = &newElem
			} else {
				// Append to current code block
				current.Content += "\n" + element.Content
				current.Height += element.Height // Approx
			}
		} else {
			// Not a code block
			if current != nil {
				// Flush current code block
				merged = append(merged, *current)
				current = nil
			}
			merged = append(merged, element)
		}
	}

	// Flush last code block if exists
	if current != nil {
		merged = append(merged, *current)
	}

	return merged
}

// MergeTableRows merges consecutive table row elements into a single table block
func (a *Analyzer) MergeTableRows(elements []Element) []Element {
	if len(elements) <= 1 {
		return elements
	}

	merged := make([]Element, 0, len(elements))
	var current *Element

	for i := 0; i < len(elements); i++ {
		element := elements[i]

		if element.Type == ElementTypeTable {
			if current == nil {
				// Start new table
				newElem := element
				// Convert 4 spaces to Markdown table separator
				newElem.Content = strings.ReplaceAll(newElem.Content, "    ", " | ")
				newElem.Content = "| " + newElem.Content + " |"
				current = &newElem
			} else {
				// Append to current table
				row := strings.ReplaceAll(element.Content, "    ", " | ")
				current.Content += "\n| " + row + " |"
				current.Height += element.Height
			}
		} else {
			// Not a table row
			if current != nil {
				// Flush current table
				merged = append(merged, *current)
				current = nil
			}
			merged = append(merged, element)
		}
	}

	// Flush last table if exists
	if current != nil {
		merged = append(merged, *current)
	}

	return merged
}
