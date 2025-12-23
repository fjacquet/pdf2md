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
