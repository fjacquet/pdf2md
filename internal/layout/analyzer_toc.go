package layout

import "strings"

// RemoveTableOfContentsRange removes the entire TOC section by finding the start and end markers.
func (a *Analyzer) RemoveTableOfContentsRange(elements []Element) []Element {
	startIndex := -1
	endIndex := -1

	// Find start: "Table of Contents" or "Contents"
	for i, el := range elements {
		content := strings.TrimSpace(el.Content)
		if strings.EqualFold(content, "Table of Contents") || strings.EqualFold(content, "Contents") {
			startIndex = i
			break
		}
	}

	// Fallback: Look for "Abstract"
	if startIndex == -1 {
		for i, el := range elements {
			content := strings.TrimSpace(el.Content)
			if strings.EqualFold(content, "Abstract") {
				// Found Abstract. The TOC usually starts after the abstract text.
				// We'll assume the abstract is 1 or 2 paragraphs.
				// We'll start deleting from the first element that looks like a TOC entry (starts with number)
				// or just skip a fixed number of elements.
				// Let's look ahead for "CHAPTER 1" and "1.1"
				startIndex = i + 1 // Start looking after Abstract header
				break
			}
		}
	}

	if startIndex == -1 {
		return elements
	}

	// Find end: "CHAPTER 1" or similar
	// We search AFTER the start index
	for i := startIndex; i < len(elements); i++ {
		el := elements[i]
		content := strings.TrimSpace(el.Content)
		// Check for Chapter 1
		if strings.HasPrefix(strings.ToUpper(content), "CHAPTER 1") {
			endIndex = i
			break
		}
	}

	if endIndex == -1 {
		return elements
	}

	// Refine Start Index if we used Abstract fallback
	// We want to keep the Abstract text.
	// We iterate from startIndex (Abstract + 1) to endIndex.
	// If we find a paragraph that doesn't look like a TOC entry, we keep it.
	// A TOC entry usually starts with a number or "Chapter".
	// Or we can just say: keep the first element after Abstract if it doesn't start with a number.
	if startIndex > 0 && strings.EqualFold(strings.TrimSpace(elements[startIndex-1].Content), "Abstract") {
		// Check if the current startIndex element is the abstract text
		// If it doesn't start with a number, assume it's text and keep it.
		// TOC entries in this file start with "1.1", "1.2", etc.
		// Abstract text starts with "This document..."
		content := strings.TrimSpace(elements[startIndex].Content)
		firstChar := ""
		if len(content) > 0 {
			firstChar = string(content[0])
		}

		// If it's not a digit, it's likely the abstract text.
		if !strings.ContainsAny(firstChar, "0123456789") {
			startIndex++ // Skip this paragraph
		}
	}

	// Remove the range [startIndex, endIndex)
	cleaned := make([]Element, 0, len(elements)-(endIndex-startIndex))
	cleaned = append(cleaned, elements[:startIndex]...)
	cleaned = append(cleaned, elements[endIndex:]...)

	return cleaned
}
