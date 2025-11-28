package layout

import (
	"math"
	"strings"
	"unicode"

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
	ElementTypeHeader     ElementType = "header"
	ElementTypeParagraph  ElementType = "paragraph"
	ElementTypeCodeBlock  ElementType = "code_block"
	ElementTypeList       ElementType = "list"
	ElementTypeTable      ElementType = "table"
	ElementTypeAdmonition ElementType = "admonition"
)

// Rule defines a classification rule for the layout analyzer
// Rule defines a classification rule for the layout analyzer
type Rule struct {
	Name      string
	Condition func(text string, fontSize, bodyFontSize float64) bool
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
			Name: "Admonition",
			Condition: func(text string, fontSize, bodyFontSize float64) bool {
				cleanText := text
				if strings.HasPrefix(cleanText, "**") {
					cleanText = strings.TrimPrefix(cleanText, "**")
				} else if strings.HasPrefix(cleanText, "*") {
					cleanText = strings.TrimPrefix(cleanText, "*")
				}
				keywords := []string{"IMPORTANT", "WARNING", "NOTE", "TIP", "CAUTION"}
				for _, kw := range keywords {
					if strings.HasPrefix(cleanText, kw) {
						return true
					}
				}
				return false
			},
			Type: ElementTypeAdmonition,
		},
		{
			Name: "Code Block (Keywords)",
			Condition: func(text string, fontSize, bodyFontSize float64) bool {
				return a.isCodeBlock(text)
			},
			Type: ElementTypeCodeBlock,
		},
		{
			Name: "Table Row (Wide Gaps)",
			Condition: func(text string, fontSize, bodyFontSize float64) bool {
				// Heuristic: Line contains wide gaps (4 spaces)
				// Require at least 1 wide gap (2 columns)
				return strings.Count(text, "    ") >= 1
			},
			Type: ElementTypeTable,
		},
		{
			Name: "Header (Numbered)",
			Condition: func(text string, fontSize, bodyFontSize float64) bool {
				return a.isNumberedHeader(text)
			},
			Type: ElementTypeHeader,
		},
		{
			Name: "List Item (Bullet/Number)",
			Condition: func(text string, fontSize, bodyFontSize float64) bool {
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

	// Step 1: Use blocks directly (assuming extractor provides reasonable order)
	orderedBlocks := blocks

	// Step 3: Detect body font size (most common)
	bodyFontSize := a.detectBodyFontSize(orderedBlocks)

	// Step 4: Convert blocks to initial elements (all paragraphs initially)
	elements := a.blocksToElements(orderedBlocks)

	// Step 5: Merge consecutive elements on same line
	elements = a.MergeElements(elements)

	// Step 6: Classify elements
	var prev *Element
	for i := range elements {
		a.classifyElement(&elements[i], prev, bodyFontSize)
		prev = &elements[i]
	}

	// Step 7: Merge consecutive code blocks (vertical)
	elements = a.MergeCodeBlocks(elements)

	// Step 8: Merge consecutive table rows
	elements = a.MergeTableRows(elements)

	return elements
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
		content := block.Text

		// Apply formatting based on font name
		if isBold(block.FontName) {
			content = "**" + content + "**"
		} else if isItalic(block.FontName) {
			content = "*" + content + "*"
		}

		element := Element{
			Content:  content,
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

func isBold(fontName string) bool {
	lower := strings.ToLower(fontName)
	return strings.Contains(lower, "bold") || strings.Contains(lower, "black") || strings.Contains(lower, "heavy")
}

func isItalic(fontName string) bool {
	lower := strings.ToLower(fontName)
	return strings.Contains(lower, "italic") || strings.Contains(lower, "oblique")
}

// classifyElement determines the type of an element based on its content and properties
func (a *Analyzer) classifyElement(element *Element, prev *Element, bodyFontSize float64) {
	text := strings.TrimSpace(element.Content)
	if text == "" {
		return
	}

	// 1. Check Rules
	for _, rule := range a.Rules {
		if rule.Condition(text, element.FontSize, bodyFontSize) {
			element.Type = rule.Type
			return
		}
	}

	// 2. Fallback: Font Size for Headers (if available)
	if element.FontSize > bodyFontSize*a.HeaderSizeRatio {
		element.Type = ElementTypeHeader
		return
	}

	// 3. Check for implicit list items (indentation-based)
	// This handles lists where bullets are missing or it's just indented text
	if prev != nil {
		// Case A: Previous line ends with ":" and current is indented
		prevText := strings.TrimSpace(prev.Content)
		if strings.HasSuffix(prevText, ":") {
			// Check indentation.
			// If element.X is significantly larger than prev.X.
			// Using 5.0 as a threshold (approx 1 char width or half-indent)
			if element.X > prev.X+5.0 {
				element.Type = ElementTypeList
				// If prev was also a list item, increment level
				if prev.Type == ElementTypeList {
					element.Level = prev.Level + 1
				} else {
					element.Level = 1
				}
				return
			}
		}

		// Case B: Previous line was a list item, and current is indented similarly
		if prev.Type == ElementTypeList {
			// If X is similar (aligned), treat as list item
			if math.Abs(element.X-prev.X) < 5.0 {
				element.Type = ElementTypeList
				element.Level = prev.Level
				return
			}
			// If X is larger (nested), treat as nested list item
			if element.X > prev.X+5.0 {
				element.Type = ElementTypeList
				element.Level = prev.Level + 1
				return
			}
		}
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

// isCodeBlock checks if the text looks like a code block
func (a *Analyzer) isCodeBlock(text string) bool {
	// Strip markdown for checks
	cleanText := text
	if strings.HasPrefix(cleanText, "**") {
		cleanText = strings.TrimPrefix(cleanText, "**")
	} else if strings.HasPrefix(cleanText, "*") {
		cleanText = strings.TrimPrefix(cleanText, "*")
	}

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
		if strings.HasPrefix(cleanText, p) {
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
		if strings.Contains(cleanText, kw) {
			// fmt.Printf("DEBUG: Found code block (contains): '%s' (keyword: '%s')\n", text, kw)
			return true
		}
	}

	// 3. Check for exact matches of brackets (unlikely in normal text unless math)
	if cleanText == "{" || cleanText == "}" || cleanText == "[]" || cleanText == "()" {
		return true
	}

	// 4. Check for JSON/YAML patterns
	// - Key-value pairs: "key": value
	// - Array/Object closers: }, ],
	if strings.Contains(cleanText, "\":") { // JSON key
		return true
	}
	if strings.HasPrefix(strings.TrimSpace(cleanText), "},") || strings.HasPrefix(strings.TrimSpace(cleanText), "],") {
		return true
	}
	if strings.HasSuffix(strings.TrimSpace(cleanText), "{") || strings.HasSuffix(strings.TrimSpace(cleanText), "[") {
		// Only if it also looks like code context (e.g. ends with open bracket)
		// But be careful with normal text ending in [
		// Check if it has a colon before it? e.g. "users": [
		if strings.Contains(cleanText, ":") {
			return true
		}
	}

	// 5. Check for SSH keys
	if strings.Contains(cleanText, "ssh-rsa ") {
		return true
	}

	// 6. Check for lines that are just brackets/braces/commas (JSON structure)
	trimmed := strings.TrimSpace(cleanText)
	isStructure := true
	if len(trimmed) > 0 {
		for _, r := range trimmed {
			if !strings.ContainsRune("{}[](), ", r) {
				isStructure = false
				break
			}
		}
		if isStructure {
			return true
		}
	}

	return false
}

// isListItem checks if the block starts with a list marker
func (a *Analyzer) isListItem(text string) bool {
	if len(text) == 0 {
		return false
	}

	// Strip markdown for checks
	cleanText := text
	if strings.HasPrefix(cleanText, "**") {
		cleanText = strings.TrimPrefix(cleanText, "**")
	} else if strings.HasPrefix(cleanText, "*") {
		cleanText = strings.TrimPrefix(cleanText, "*")
	}

	// Check for bullet points
	if strings.HasPrefix(cleanText, "•") || strings.HasPrefix(cleanText, "- ") || strings.HasPrefix(cleanText, "* ") {
		return true
	}

	// Check for numbered lists (1., 2., etc.)
	// Simple check for digit + dot + space
	if len(cleanText) >= 3 && cleanText[0] >= '0' && cleanText[0] <= '9' && cleanText[1] == '.' && cleanText[2] == ' ' {
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

			if gap < threshold {
				current.Content += next.Content
				// Update width
				current.Width = next.X + next.Width - current.X
			} else if gap > wideGapThreshold {
				// Wide gap, likely a table column or visual separation
				// Do NOT merge. Treat as separate elements.
				merged = append(merged, current)
				current = next
			} else {
				current.Content += " " + next.Content
				// Update width
				current.Width = next.X + next.Width - current.X
			}
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

// MergeParagraphLines merges consecutive paragraph lines
func (a *Analyzer) MergeParagraphLines(elements []Element) []Element {
	if len(elements) == 0 {
		return nil
	}

	var merged []Element
	current := &elements[0]

	for i := 1; i < len(elements); i++ {
		next := elements[i]

		if current.Type == ElementTypeParagraph && next.Type == ElementTypeParagraph {
			// Check if duplicate (exact match)
			if current.Content == next.Content {
				// Skip duplicate
				continue
			}

			// Check vertical gap (baseline to baseline)
			gap := math.Abs(current.Y - next.Y)

			// Check for sentence break (terminator + uppercase start)
			isSentenceBreak := false
			trimmedCurrent := strings.TrimSpace(current.Content)
			trimmedNext := strings.TrimSpace(next.Content)
			if (strings.HasSuffix(trimmedCurrent, ".") || strings.HasSuffix(trimmedCurrent, "?") || strings.HasSuffix(trimmedCurrent, "!")) &&
				len(trimmedNext) > 0 && unicode.IsUpper(rune(trimmedNext[0])) {
				isSentenceBreak = true
			}

			// Determine threshold
			// If sentence break, use tighter threshold (1.2x) to avoid merging distinct paragraphs
			// If not sentence break, use looser threshold (1.6x) to handle split lines
			threshold := current.FontSize * 1.6
			if isSentenceBreak {
				threshold = current.FontSize * 1.2
			}

			// Also check if next looks like a list item (don't merge if so)
			isList := a.isListItem(next.Content)

			// Debug
			// fmt.Printf("DEBUG: MergeParagraphLines: '%s' vs '%s'. Gap=%.2f. Threshold=%.2f. isList=%v, isSentenceBreak=%v\n", current.Content, next.Content, gap, threshold, isList, isSentenceBreak)

			// If gap is small (less than threshold), merge
			// Standard line height is ~1.2. Paragraph spacing is usually > 1.4.
			// Also check if next looks like a list item (don't merge if so)

			if gap < threshold && !isList {
				// Add space if needed
				if !strings.HasSuffix(current.Content, " ") && !strings.HasSuffix(current.Content, "-") {
					current.Content += " "
				}
				current.Content += next.Content
				// Update height? For paragraphs, height is usually just the block height.
				// If we merge lines, the height should cover both.
				// But Y is baseline.
				// Let's keep the logic simple for now.
			} else {
				merged = append(merged, *current)
				current = &elements[i]
			}
		} else {
			merged = append(merged, *current)
			current = &elements[i]
		}
	}
	merged = append(merged, *current)

	return merged
}
