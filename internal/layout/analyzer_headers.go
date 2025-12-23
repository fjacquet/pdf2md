package layout

import "strings"

// NormalizeHeaderLevels adjusts header levels so the smallest level becomes H1
func (a *Analyzer) NormalizeHeaderLevels(elements []Element) []Element {
	// Find the minimum header level
	minLevel := 7
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
			elements[i].Level -= shift
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
