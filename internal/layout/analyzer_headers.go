package layout

import (
	"cmp"
	"math"
	"slices"
	"strings"
)

// RankHeaderLevels sets every header's level from the document's own
// hierarchy: its largest heading style is #, the next ##, and so on, capped at
// ######. A style is the font size rounded to 0.5 pt plus the bold flag; at
// the same size bold ranks first. Within one style, deeper section numbers
// ("1.1" under "1") rank lower; unnumbered headers sit with the shallowest.
func (a *Analyzer) RankHeaderLevels(elements []Element) []Element {
	type style struct {
		size  float64
		bold  bool
		depth int
	}
	// ONNX and table-row headers skip classifyElement, so Bold is only in the markers.
	base := func(e Element) style { return style{math.Round(e.FontSize*2) / 2, e.Bold || fullyBold(e.Content), 0} }
	shallowest := map[style]int{}
	for _, e := range elements {
		if d := numberingDepth(stripBold(e.Content)); e.Type == ElementTypeHeader && d > 0 {
			if m, ok := shallowest[base(e)]; !ok || d < m {
				shallowest[base(e)] = d
			}
		}
	}
	key := func(e Element) style {
		s := base(e)
		if s.depth = numberingDepth(stripBold(e.Content)); s.depth == 0 {
			s.depth = shallowest[base(e)]
		}
		return s
	}
	var styles []style
	for _, e := range elements {
		if e.Type == ElementTypeHeader && !slices.Contains(styles, key(e)) {
			styles = append(styles, key(e))
		}
	}
	slices.SortFunc(styles, func(x, y style) int {
		if x.size != y.size {
			return cmp.Compare(y.size, x.size)
		}
		if x.bold != y.bold {
			if x.bold {
				return -1
			}
			return 1
		}
		return cmp.Compare(x.depth, y.depth)
	})
	for i := range elements {
		if elements[i].Type == ElementTypeHeader {
			elements[i].Level = min(slices.Index(styles, key(elements[i]))+1, 6)
		}
	}
	return elements
}

// numberingDepth is the depth of a leading section number ("2" is 1, "2.1."
// is 2), or 0 when the text does not start with one.
func numberingDepth(text string) int {
	marker, _, _ := strings.Cut(text, " ")
	parts := strings.Split(strings.TrimSuffix(marker, "."), ".")
	for _, p := range parts {
		if p == "" || strings.Trim(p, "0123456789") != "" {
			return 0
		}
	}
	return len(parts)
}

// FilterPageNumbers removes standalone page numbers from the output
func (a *Analyzer) FilterPageNumbers(elements []Element) []Element {
	result := make([]Element, 0, len(elements))
	for _, el := range elements {
		content := strings.TrimSpace(el.Content)
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

// FilterRepeatedPageHeaders removes repeated page headers/footers
// These are text elements that appear multiple times with the exact same content,
// typically at similar Y positions (running headers/footers)
func (a *Analyzer) FilterRepeatedPageHeaders(elements []Element) []Element {
	if len(elements) < 2 {
		return elements
	}

	// Count occurrences of each text content
	contentCounts := make(map[string]int)
	for _, el := range elements {
		content := strings.TrimSpace(el.Content)
		if content != "" && len(content) < 100 { // Only consider shorter text for page headers
			contentCounts[content]++
		}
	}

	// Find texts that appear 3+ times (likely page headers/footers)
	repeatedContent := make(map[string]bool)
	for content, count := range contentCounts {
		if count >= 3 {
			// Additional check: should look like a header (e.g., "CHAPTER X. TITLE" pattern)
			if isLikelyPageHeader(content) {
				repeatedContent[content] = true
			}
		}
	}

	// If no repeated content found, return as-is
	if len(repeatedContent) == 0 {
		return elements
	}

	// Filter out repeated page headers, keeping only the first occurrence
	firstOccurrence := make(map[string]bool)
	result := make([]Element, 0, len(elements))

	for _, el := range elements {
		content := strings.TrimSpace(el.Content)
		if repeatedContent[content] {
			if !firstOccurrence[content] {
				// Keep the first occurrence
				firstOccurrence[content] = true
				result = append(result, el)
			}
			// Skip subsequent occurrences
			continue
		}
		result = append(result, el)
	}

	return result
}

// isLikelyPageHeader checks if text looks like a page header
func isLikelyPageHeader(text string) bool {
	text = strings.TrimSpace(text)

	// Pattern: "CHAPTER X. TITLE" or similar
	if strings.HasPrefix(strings.ToUpper(text), "CHAPTER ") {
		return true
	}

	// Pattern: "SECTION X. TITLE"
	if strings.HasPrefix(strings.ToUpper(text), "SECTION ") {
		return true
	}

	// Pattern: "PART X. TITLE"
	if strings.HasPrefix(strings.ToUpper(text), "PART ") {
		return true
	}

	// All uppercase short text (common for headers)
	if text == strings.ToUpper(text) && len(text) > 5 && len(text) < 60 {
		// Check if it's not a regular sentence (no common lowercase patterns)
		wordCount := len(strings.Fields(text))
		if wordCount >= 2 && wordCount <= 8 {
			return true
		}
	}

	return false
}
