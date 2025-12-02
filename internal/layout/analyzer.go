package layout

import (
	"fmt"
	"math"
	"sort"
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
	LinkURI  string
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
	ElementTypeImage      ElementType = "image"
)

// Rule defines a classification rule for the layout analyzer
// Rule defines a classification rule for the layout analyzer
type Rule struct {
	Name      string
	Condition func(text string, fontSize, bodyFontSize float64) bool
	Type      ElementType
}

// ExclusionZone defines areas to ignore (e.g., headers, footers)
type ExclusionZone struct {
	Top    float64 // Height from top to ignore
	Bottom float64 // Height from bottom to ignore
}

// Analyzer analyzes the layout of text blocks
type Analyzer struct {
	ColumnGapThreshold float64
	HeaderSizeRatio    float64
	Rules              []Rule
	Exclusion          ExclusionZone
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
				if strings.Count(text, "    ") >= 1 {
					return true
				}
				// Also allow 3 spaces if there are multiple gaps (likely a table with math)
				// But exclude single gap (likely equation + number)
				// We check for "   " (3 spaces) which is used for math lines
				if strings.Count(text, "   ") >= 2 {
					return true
				}
				return false
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
func (a *Analyzer) Analyze(content *extractor.PageContent) []Element {
	if content == nil {
		return nil
	}
	blocks := content.TextBlocks
	images := content.Images
	graphics := content.Graphics
	// links := content.Links // Unused for now in layout analysis?

	if len(blocks) == 0 && len(images) == 0 && len(graphics) == 0 {
		return nil
	}

	// Step 0: Filter blocks based on exclusion zones
	if a.Exclusion.Top > 0 || a.Exclusion.Bottom > 0 {
		filteredBlocks := make([]extractor.TextBlock, 0, len(blocks))
		pageHeight := content.PageHeight
		if pageHeight == 0 {
			// Fallback: try to estimate from blocks? Or just ignore bottom exclusion?
			// If PageHeight is 0, we can't do bottom exclusion properly unless we assume standard A4/Letter?
			// Let's assume 0 means "unknown" and skip bottom exclusion if so, or warn.
			// But Top exclusion works regardless of page height (Y starts at bottom? No, PDF Y starts at bottom usually).
			// Wait, if PDF Y starts at bottom (0), then Top is Y > PageHeight - TopMargin.
			// And Bottom is Y < BottomMargin.
			// So we DO need PageHeight for Top exclusion if Y is bottom-up.
			// Let's check coordinate system.
			// In `custom.go`, we extract using `pdf.Reader`.
			// `interpreter.go` usually normalizes?
			// Let's assume standard PDF coordinates: (0,0) is bottom-left.
			// So Top margin means Y > (PageHeight - TopMargin).
			// Bottom margin means Y < BottomMargin.
		}

		for _, block := range blocks {
			// Check Bottom margin
			if a.Exclusion.Bottom > 0 && block.Y < a.Exclusion.Bottom {
				continue
			}
			// Check Top margin
			if a.Exclusion.Top > 0 && pageHeight > 0 && block.Y > (pageHeight-a.Exclusion.Top) {
				continue
			}
			filteredBlocks = append(filteredBlocks, block)
		}
		blocks = filteredBlocks
	}

	// Step 0.5: Apply links to blocks (before merging)
	blocks = ApplyLinksToBlocks(blocks, content.Links)

	// Step 1: Sort blocks using Recursive XY-Cut (Column Detection)
	orderedBlocks := SortBlocks(blocks)

	// Step 3: Detect body font size (most common)
	bodyFontSize := a.detectBodyFontSize(orderedBlocks)

	// Step 4: Convert blocks to initial elements (all paragraphs initially)
	elements := a.blocksToElements(orderedBlocks)

	// Step 4.5: Add images to elements
	for _, img := range images {
		elements = append(elements, Element{
			Type:    ElementTypeImage,
			Content: img.ID, // Use ID as content for now, or path later
			X:       img.X,
			Y:       img.Y,
			Width:   img.Width,
			Height:  img.Height,
		})
	}

	// Step 4.6: Add vector graphics to elements (filter small decorative ones)
	for _, vg := range graphics {
		// Filter out small decorative graphics (bullets, icons, etc.)
		// Minimum size threshold: 50x50 points
		if vg.Width < 50 || vg.Height < 50 {
			continue
		}
		elements = append(elements, Element{
			Type:    ElementTypeImage, // Treat as image for now
			Content: vg.ID,
			X:       vg.X,
			Y:       vg.Y,
			Width:   vg.Width,
			Height:  vg.Height,
		})
	}

	// Step 4.7: Merge consecutive elements with same LinkURI (before general merge)
	// This prevents fragmentation of links
	elements = a.MergeSameLinkElements(elements)

	// Step 5: Merge consecutive elements on same line
	elements = a.MergeElements(elements)

	// Step 6: Classify elements
	var prev *Element
	for i := range elements {
		// Skip classification for images
		if elements[i].Type == ElementTypeImage {
			continue
		}
		a.classifyElement(&elements[i], prev, bodyFontSize)
		prev = &elements[i]
	}

	// Step 5.5: Merge consecutive paragraph lines (vertical merge)
	// Moved after classification to prevent merging Headers/CodeBlocks with Paragraphs
	elements = a.MergeParagraphLines(elements)

	// Step 7: Merge consecutive code blocks (vertical)
	elements = a.MergeCodeBlocks(elements)

	// Step 8: Merge consecutive table rows
	elements = a.MergeTableRows(elements)

	// Step 9: Clean math symbols
	elements = a.CleanMathSymbols(elements)

	// Step 10: Filter page numbers (standalone numbers at page boundaries)
	elements = a.FilterPageNumbers(elements)

	// Step 11: Normalize header levels (make the highest level become H1)
	elements = a.NormalizeHeaderLevels(elements)

	return elements
}

// NormalizeHeaderLevels adjusts header levels so the smallest level (largest header) becomes H1
func (a *Analyzer) NormalizeHeaderLevels(elements []Element) []Element {
	// Find the minimum header level (i.e., the "biggest" header)
	minLevel := 7 // Start above max (6)
	for _, el := range elements {
		if el.Type == ElementTypeHeader && el.Level > 0 && el.Level < minLevel {
			minLevel = el.Level
		}
	}

	// If no headers found or already H1, return as-is
	if minLevel >= 7 || minLevel == 1 {
		return elements
	}

	// Shift all header levels down
	shift := minLevel - 1
	for i := range elements {
		if elements[i].Type == ElementTypeHeader {
			elements[i].Level = elements[i].Level - shift
			if elements[i].Level < 1 {
				elements[i].Level = 1
			}
			if elements[i].Level > 6 {
				elements[i].Level = 6
			}
		}
	}

	return elements
}

// FilterPageNumbers removes standalone page numbers from the output
func (a *Analyzer) FilterPageNumbers(elements []Element) []Element {
	result := make([]Element, 0, len(elements))
	for _, el := range elements {
		content := strings.TrimSpace(el.Content)
		// Skip elements that look like page numbers:
		// - Standalone numbers 1-999
		// - Numbers with common page number patterns like "- 5 -" or "Page 5"
		if a.isPageNumber(content) {
			continue
		}
		result = append(result, el)
	}
	return result
}

// isPageNumber checks if text looks like a page number
func (a *Analyzer) isPageNumber(text string) bool {
	text = strings.TrimSpace(text)

	// Empty text
	if text == "" {
		return false
	}

	// Pure number (1-999)
	if len(text) <= 3 {
		allDigits := true
		for _, c := range text {
			if c < '0' || c > '9' {
				allDigits = false
				break
			}
		}
		if allDigits && text != "" {
			return true
		}
	}

	// Patterns like "- 5 -", "— 5 —"
	text = strings.Trim(text, "-–—")
	text = strings.TrimSpace(text)
	if len(text) <= 3 {
		allDigits := true
		for _, c := range text {
			if c < '0' || c > '9' {
				allDigits = false
				break
			}
		}
		if allDigits && text != "" {
			return true
		}
	}

	return false
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
		// Only apply if content is not just whitespace
		if strings.TrimSpace(content) != "" {
			if isMathFont(block.FontName) {
				content = "$" + content + "$"
			} else if isBold(block.FontName) {
				content = "**" + content + "**"
			} else if isItalic(block.FontName) {
				content = "*" + content + "*"
			}
		}

		element := Element{
			Content:  content,
			X:        block.X,
			Y:        block.Y,
			Width:    block.Width,
			Height:   block.Height,
			FontSize: block.FontSize,
			Type:     ElementTypeParagraph, // Default
			LinkURI:  block.LinkURI,
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
		element.Level = a.calculateHeaderLevel(element.FontSize, bodyFontSize)
		// Strip bold markers from headers (formatting is already implied by header level)
		content := element.Content
		if strings.HasPrefix(content, "**") && strings.HasSuffix(content, "**") {
			element.Content = strings.TrimPrefix(strings.TrimSuffix(content, "**"), "**")
		}
		return
	}

	// 2.5. Bold short text as headers
	// Bold text that is short (< 100 chars) and doesn't end with punctuation
	// is likely a header/section title
	if strings.HasPrefix(text, "**") && strings.HasSuffix(text, "**") {
		innerText := strings.TrimPrefix(strings.TrimSuffix(text, "**"), "**")
		innerText = strings.TrimSpace(innerText)

		// Skip if it looks like a code line number (starts with a digit)
		if len(innerText) > 0 && innerText[0] >= '0' && innerText[0] <= '9' {
			// Don't treat as header
		} else if len(innerText) < 100 && len(innerText) > 0 &&
			!strings.HasSuffix(innerText, ".") &&
			!strings.HasSuffix(innerText, ",") &&
			!strings.HasSuffix(innerText, ";") &&
			!strings.Contains(innerText, "\n") &&
			len(innerText) > 2 { // Must be more than 2 chars
			element.Type = ElementTypeHeader
			element.Level = 4 // Set to H4, will be normalized later
			// Strip the bold markers since header formatting will be added
			element.Content = innerText
			return
		}
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

	// Check if first part is a number like "1." or "1.2" or "1.2.3" or Roman numerals
	marker := parts[0]
	if len(marker) == 0 {
		return false
	}

	// Check for Roman numeral headers (I. II. III. IV. V. VI. VII. VIII. IX. X.)
	romanNumerals := []string{"I.", "II.", "III.", "IV.", "V.", "VI.", "VII.", "VIII.", "IX.", "X.",
		"XI.", "XII.", "XIII.", "XIV.", "XV."}
	for _, rn := range romanNumerals {
		if marker == rn {
			return true
		}
	}

	// Must start with digit
	if marker[0] < '0' || marker[0] > '9' {
		return false
	}

	// Check for section number patterns:
	// - Ends with period: "1." "1.2." "1.2.3."
	// - Or contains internal dots (section numbers): "8.3" "1.2.3"
	// This prevents matching single-digit code line numbers like "1 #" but allows "8.3 Accessing"
	hasPeriod := strings.HasSuffix(marker, ".")
	hasInternalDot := strings.Contains(strings.TrimSuffix(marker, "."), ".")

	// If no period at end and no internal dots, it's likely a code line number (reject)
	if !hasPeriod && !hasInternalDot {
		return false
	}

	// Check that all characters (except trailing dot) are digits or dots
	markerToCheck := strings.TrimSuffix(marker, ".")
	if len(markerToCheck) == 0 {
		return false
	}

	validChars := "0123456789."
	for _, c := range markerToCheck {
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
		"return ", "if ", "else ", "range ",
		"$ ", "# ", // Shell prompts
		"kubectl ", "oc ", "docker ", "podman ", "sudo ",
		"apiVersion:", "kind:", "metadata:", "spec:", "status:",
	}

	for _, p := range prefixes {
		if strings.HasPrefix(cleanText, p) {
			// Special check for shell prompt "$ "
			if p == "$ " {
				// If it ends with "$", it's likely math: "$ ... $"
				if strings.HasSuffix(strings.TrimSpace(cleanText), "$") {
					continue
				}
			}
			// Special check for "if " to avoid false positives in text
			if p == "if " {
				// "if" is very common. Only treat as code if it looks like code syntax.
				// e.g. "if (", "if {", "if x >", "if err !="
				// Simple heuristic: check for typical code chars
				rest := strings.TrimPrefix(cleanText, "if ")
				if !strings.ContainsAny(rest, "({=><!") {
					continue
				}
			}
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

	// Check for reference style lists ([1], [12], etc.)
	if strings.HasPrefix(cleanText, "[") {
		end := strings.Index(cleanText, "]")
		if end > 1 {
			inner := cleanText[1:end]
			isNum := true
			for _, r := range inner {
				if !unicode.IsDigit(r) {
					isNum = false
					break
				}
			}
			if isNum {
				return true
			}
		}
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
			// Using a heuristic here: if gap is less than 10% of font size, assume it's part of the same word
			// Or if gap is negative (overlap)
			threshold := current.FontSize * 0.1
			if threshold == 0 {
				threshold = 1.0 // Fallback
			}

			// Determine threshold for "wide gap"
			// If FontSize is available, use relative. Otherwise use absolute.
			// Increased to 3.0x to avoid false positives for tables (was 2.0x)
			wideGapThreshold := 30.0
			if current.FontSize > 0 {
				wideGapThreshold = current.FontSize * 3.0
			}

			if gap < threshold {
				// Merge logic with LinkURI handling
				current = a.mergeElementsWithLinks(current, next, "")
			} else if gap > wideGapThreshold {
				// Wide gap, likely a table column.
				// Check if content is long (likely text columns, not table)
				if len(current.Content) > 40 || len(next.Content) > 40 {
					// Don't merge, treat as separate blocks
					merged = append(merged, current)
					current = next
				} else {
					// Check for Equation Numbering (e.g. "(1)", "(2.1)")
					// If the right side is just a number in parens, treat as Equation, not Table.
					// We check if 'next' is the start of an equation number, potentially split across multiple blocks.
					isEqNum := a.isStartOfEqNum(next, elements[i+1:])

					if isEqNum {
						// Merge with space (or special separator?)
						// Just space to keep it as text/paragraph
						current = a.mergeElementsWithLinks(current, next, " ")
					} else if strings.Contains(current.Content, "$") || strings.Contains(next.Content, "$") {
						// If content looks like math, use 3 spaces for wide gaps instead of 4 spaces
						// This prevents equations from being detected as tables (which requires 4 spaces)
						current = a.mergeElementsWithLinks(current, next, "   ")
					} else {
						// Merge with 4 spaces to allow table detection.
						current = a.mergeElementsWithLinks(current, next, "    ")
					}
				}
			} else {
				// Regular gap (space)

				// Check for Hyphenation
				// If current ends with "-" and gap is not huge, merge without space.
				// This fixes "sur- geons" -> "sur-geons"
				if strings.HasSuffix(current.Content, "-") || strings.HasSuffix(current.Content, "‐") {
					// Check if next starts with lowercase (optional, but safer)
					// Actually, just merging is usually safe for hyphens in this context.
					// But let's ensure gap isn't massive (it's already < wideGapThreshold)
					current = a.mergeElementsWithLinks(current, next, "")
				} else {
					current = a.mergeElementsWithLinks(current, next, " ")
				}
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

// mergeElementsWithLinks merges two elements handling their LinkURIs
func (a *Analyzer) mergeElementsWithLinks(e1, e2 Element, separator string) Element {
	// If both have same LinkURI, keep it
	if e1.LinkURI == e2.LinkURI {
		e1.Content = mergeWithStyle(e1.Content, separator+e2.Content)
		e1.Width = e2.X + e2.Width - e1.X
		return e1
	}

	// If different LinkURIs (or one is empty), we must finalize the link for the one that has it
	// because the merged element can only have one LinkURI (or none if mixed).
	// Actually, if we merge them into one string, we can't keep LinkURI property for the whole string
	// if it only applies to part of it.
	// So we should format the link as Markdown in the Content and clear LinkURI.

	content1 := e1.Content
	if e1.LinkURI != "" {
		content1 = fmt.Sprintf("[%s](%s)", e1.Content, e1.LinkURI)
	}

	content2 := e2.Content
	if e2.LinkURI != "" {
		content2 = fmt.Sprintf("[%s](%s)", e2.Content, e2.LinkURI)
	}

	e1.Content = mergeWithStyle(content1, separator+content2)
	e1.Width = e2.X + e2.Width - e1.X
	e1.LinkURI = "" // Mixed or handled, so clear it

	return e1
}

// mergeWithStyle merges two strings, handling markdown style markers (bold/italic)
// to avoid artifacts like "**A****B**" -> "**AB**"
func mergeWithStyle(a, b string) string {
	separator := ""
	cleanB := b
	if strings.HasPrefix(b, " ") {
		separator = " "
		cleanB = b[1:]
	}

	// Check bold
	if strings.HasSuffix(a, "**") && strings.HasPrefix(cleanB, "**") {
		// a = "...**", cleanB = "**..."
		// Result: "... " + "..." (merged)
		return a[:len(a)-2] + separator + cleanB[2:]
	}
	// Check italic
	if strings.HasSuffix(a, "*") && strings.HasPrefix(cleanB, "*") {
		return a[:len(a)-1] + separator + cleanB[1:]
	}

	// Check math ($...$)
	if strings.HasSuffix(a, "$") && strings.HasPrefix(cleanB, "$") {
		// a = "...$", cleanB = "$..."
		// Result: "... " + "..." (merged, keeping outer $)
		return a[:len(a)-1] + separator + cleanB[1:]
	}

	return a + b
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
				// Also handle 3 spaces (for math tables)
				row = strings.ReplaceAll(row, "   ", " | ")

				current.Content += "\n| " + row + " |"
				current.Height += element.Height
			}
		} else {
			// Not a table row
			if current != nil {
				// Flush current table
				// Check if it's a single-column table (false positive)
				// A single column table will have lines like "| Content |"
				// So each line has exactly 2 pipes.
				lines := strings.Split(current.Content, "\n")
				isSingleCol := true
				for _, line := range lines {
					if strings.Count(line, "|") > 2 {
						isSingleCol = false
						break
					}
				}

				if isSingleCol {
					// Revert to paragraph(s)
					// We need to strip the pipes
					var newContent strings.Builder
					for i, line := range lines {
						if i > 0 {
							newContent.WriteString("\n")
						}
						// Strip leading "| " and trailing " |"
						trimmed := strings.TrimSpace(line)
						if strings.HasPrefix(trimmed, "| ") {
							trimmed = trimmed[2:]
						}
						if strings.HasSuffix(trimmed, " |") {
							trimmed = trimmed[:len(trimmed)-2]
						}
						newContent.WriteString(trimmed)
					}
					current.Content = newContent.String()
					current.Type = ElementTypeParagraph // Revert type

					// Re-classify? It might be a header.
					// But we lost the original font size info for individual lines if we merged them.
					// However, current.FontSize should be the font size of the first element.
					// Let's leave it as Paragraph for now, or maybe Header if it looks like one.
					// Actually, if it was "I. INTRODUCTION", it might be a header.
					if a.isNumberedHeader(current.Content) {
						current.Type = ElementTypeHeader
						current.Level = 2 // Default to H2 for now
					}
				}

				merged = append(merged, *current)
				current = nil
			}
			merged = append(merged, element)
		}
	}

	// Flush last table if exists
	if current != nil {
		// Check single col for last one too
		lines := strings.Split(current.Content, "\n")
		isSingleCol := true
		for _, line := range lines {
			if strings.Count(line, "|") > 2 {
				isSingleCol = false
				break
			}
		}

		if isSingleCol {
			var newContent strings.Builder
			for i, line := range lines {
				if i > 0 {
					newContent.WriteString("\n")
				}
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "| ") {
					trimmed = trimmed[2:]
				}
				if strings.HasSuffix(trimmed, " |") {
					trimmed = trimmed[:len(trimmed)-2]
				}
				newContent.WriteString(trimmed)
			}
			current.Content = newContent.String()
			current.Type = ElementTypeParagraph
			if a.isNumberedHeader(current.Content) {
				current.Type = ElementTypeHeader
				current.Level = 2
			}
		}
		merged = append(merged, *current)
	}

	return merged
}

// MergeParagraphLines merges consecutive paragraph lines
func (a *Analyzer) MergeParagraphLines(elements []Element) []Element {
	if len(elements) == 0 {
		return nil
	}

	// Calculate median gap to detect line spacing
	var gaps []float64
	for i := 0; i < len(elements)-1; i++ {
		curr := elements[i]
		next := elements[i+1]
		if curr.Type == ElementTypeParagraph && next.Type == ElementTypeParagraph {
			gap := math.Abs(curr.Y - next.Y)
			// Filter out huge gaps (likely page breaks or distinct sections)
			if gap < curr.FontSize*5 {
				gaps = append(gaps, gap)
			}
		}
	}

	medianGap := 0.0
	if len(gaps) > 0 {
		sort.Float64s(gaps)
		medianGap = gaps[len(gaps)/2]
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
			// Use median gap as baseline if available and reasonable
			baseThreshold := current.FontSize * 1.6
			if medianGap > current.FontSize*1.5 {
				// Document has wide spacing, use medianGap * 1.2 as base
				baseThreshold = medianGap * 1.2
			}

			threshold := baseThreshold
			if isSentenceBreak {
				// Stricter threshold for sentence breaks - don't merge separate sentences
				// Use a fixed ratio based on font size, ignoring median gap
				// Two sentences on separate lines should generally stay separate
				threshold = current.FontSize * 1.5
			}

			// Also check if next looks like a list item (don't merge if so)
			isList := a.isListItem(next.Content)

			// If gap is small (less than threshold), merge
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

// CleanMathSymbols merges consecutive math symbols to improve readability
// e.g. "$n$" "$n$" -> "$n$"
// CleanMathSymbols merges consecutive math symbols to improve readability
// e.g. "$n$" "$n$" -> "$n$"
func (a *Analyzer) CleanMathSymbols(elements []Element) []Element {
	for i := range elements {
		if elements[i].Type == ElementTypeParagraph || elements[i].Type == ElementTypeHeader {
			text := elements[i].Content

			// Split by "$" to identify math blocks
			// Even indices are text, Odd indices are math content (assuming balanced $)
			parts := strings.Split(text, "$")
			if len(parts) < 3 {
				continue
			}

			var newParts []string
			newParts = append(newParts, parts[0]) // First text part

			// Iterate over math blocks (indices 1, 3, 5...)
			for j := 1; j < len(parts); j += 2 {
				mathContent := parts[j]

				// Check if we have a next math block
				if j+2 < len(parts) {
					separator := parts[j+1]
					nextMathContent := parts[j+2]

					// Check if separator is empty or whitespace
					isSeparatorClean := strings.TrimSpace(separator) == ""

					// Check if content is identical
					if isSeparatorClean && mathContent == nextMathContent {
						// Duplicate found! Skip this block and the separator
						// We effectively merge by ignoring the current one and keeping the next one
						// (or vice versa).
						// Actually, let's keep the current one and skip the next one?
						// No, if we have A A A, we want A.
						// If we keep current, we add A. Then we look at next A.
						// We need to look ahead.

						// If match, we skip adding the separator and the next math content?
						// No, that's complex.

						// Let's use a "skip next" flag or just modify the loop.
						// If duplicate, we merge them.
						// "A" + sep + "A" -> "A"
						// So we add "A", and we want the next iteration to skip the first "A" of the pair?
						// No.

						// Let's restart logic.
						// We have [text0, math1, text2, math3, text4 ...]
						// We want to merge math1 and math3 if text2 is empty/space and math1==math3.

						// We added parts[0].
						// Loop j=1. math1.
						// Check j+2 (math3).
						// If mergeable:
						//   We want to keep math1.
						//   We want to discard text2 and math3.
						//   But wait, if we have A A A.
						//   j=1 (A). Next is A. Merge.
						//   We keep A. We skip text2 and math3.
						//   Next loop should start at j+2? No, j+2 is the one we skipped.
						//   We need to check if the *merged* A is followed by another A.
						//   So we should NOT advance j yet, or we should check recursively.

						// Easier: Loop and only append if NOT duplicate of *previous*?
						// But we need to check separator.

						// Let's use a while loop style.
						// currentMath = parts[j]
						// k = j + 2
						// while k < len(parts) && parts[k] == currentMath && strings.TrimSpace(parts[k-1]) == "" {
						//    k += 2
						// }
						// Now k points to the first non-duplicate math block (or end).
						// We append currentMath.
						// We append parts[k-1] (the text before the next math block).
						// Set j = k - 2 (so loop increment j+=2 makes j=k).

						k := j + 2
						for k < len(parts) {
							sep := parts[k-1]
							next := parts[k]
							if next == mathContent && strings.TrimSpace(sep) == "" {
								k += 2
							} else {
								break
							}
						}

						newParts = append(newParts, mathContent)
						if k-1 < len(parts) {
							newParts = append(newParts, parts[k-1]) // Separator/Text after the last merged block
						}

						// Advance j
						j = k - 2
						continue
					}
				}

				// No duplicate, just add
				newParts = append(newParts, mathContent)
				if j+1 < len(parts) {
					newParts = append(newParts, parts[j+1])
				}
			}

			// Reconstruct string
			// We split by $, so we need to join by $?
			// Wait, newParts contains [text, math, text, math, text].
			// We need to put $ back around math parts.
			// Actually, strings.Split removes $.
			// So we need to add $ back.
			// newParts[0] is text.
			// newParts[1] is math. -> needs $ around it.
			// newParts[2] is text.

			var sb strings.Builder
			for k, p := range newParts {
				sb.WriteString(p)
				// Add $ if this was a math block (odd index in original, but here?)
				// In newParts:
				// 0: text
				// 1: math
				// 2: text
				// 3: math
				// So odd indices are math.
				// Wait, if we merged, we kept the structure [text, math, text, math].
				// So yes, odd indices are math.
				// Except the last one?
				// len(parts) is odd (text, math, text ... math, text).
				// newParts should also be odd length.

				if k < len(newParts)-1 {
					sb.WriteString("$")
				}
			}
			elements[i].Content = sb.String()
		}
	}
	return elements
}

// isStartOfEqNum checks if the current element and subsequent elements form an equation number like "(1)" or "(2.1)"
func (a *Analyzer) isStartOfEqNum(current Element, remaining []Element) bool {
	text := strings.TrimSpace(current.Content)
	if !strings.HasPrefix(text, "(") {
		return false
	}

	// If it already contains ")", check if it's a valid EqNum
	if strings.Contains(text, ")") {
		// Extract content inside parens
		start := strings.Index(text, "(")
		end := strings.LastIndex(text, ")")
		if start < end {
			inner := text[start+1 : end]
			// Allow digits and dots
			for _, r := range inner {
				if !unicode.IsDigit(r) && r != '.' && r != ' ' { // Allow spaces inside parens? e.g. ( 1 )
					return false
				}
			}
			return true
		}
		return false
	}

	// Look ahead
	acc := text
	for _, el := range remaining {
		// Check if element is on same line
		if math.Abs(el.Y-current.Y) > 2.0 {
			break
		}

		content := strings.TrimSpace(el.Content)
		acc += content
		if strings.Contains(content, ")") {
			// Found closing paren
			// Check validity
			start := strings.Index(acc, "(")
			end := strings.LastIndex(acc, ")")
			if start < end {
				inner := acc[start+1 : end]
				for _, r := range inner {
					if !unicode.IsDigit(r) && r != '.' && r != ' ' {
						return false
					}
				}
				return true
			}
			return false
		}

		if len(acc) > 20 { // Sanity limit
			return false
		}
	}

	return false
}
