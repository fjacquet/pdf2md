package layout

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode"
)

// MergeDropCaps merges drop cap letters with their following text
// Drop caps are large initial letters that appear at the start of paragraphs
func (a *Analyzer) MergeDropCaps(elements []Element) []Element {
	if len(elements) <= 1 {
		return elements
	}

	merged := make([]Element, 0, len(elements))

	for i := 0; i < len(elements); i++ {
		current := elements[i]

		// Check if this could be a drop cap:
		// 1. Single uppercase letter (or 1-2 chars)
		// 2. Followed by text that starts with lowercase or continues the word
		content := strings.TrimSpace(current.Content)
		if len(content) <= 2 && i+1 < len(elements) {
			next := elements[i+1]
			nextContent := strings.TrimSpace(next.Content)

			// Check if content is uppercase letter(s)
			isUpperStart := len(content) > 0 && unicode.IsUpper(rune(content[0]))

			// Check if next starts with lowercase (continuation)
			startsLower := len(nextContent) > 0 && unicode.IsLower(rune(nextContent[0]))

			// Check proximity - drop caps are usually close horizontally or vertically aligned
			xGap := next.X - (current.X + current.Width)
			yDiff := math.Abs(current.Y - next.Y)

			// Drop cap detection: single/double uppercase followed by lowercase text
			// with reasonable proximity (within 3x font size horizontally, 2x vertically)
			isDropCap := isUpperStart && startsLower &&
				xGap < current.FontSize*3 && xGap > -current.FontSize &&
				yDiff < current.FontSize*2

			// Also check for larger font size (drop caps are often bigger)
			if current.FontSize > next.FontSize*1.3 {
				isDropCap = isUpperStart && startsLower && xGap < current.FontSize*5
			}

			if isDropCap {
				// Merge drop cap with following text
				current.Content = content + nextContent
				current.Width = next.X + next.Width - current.X
				current.FontSize = next.FontSize // Use the body text font size
				i++                              // Skip next element
			}
		}

		merged = append(merged, current)
	}

	return merged
}

// MergeElements merges consecutive text elements on the same line
func (a *Analyzer) MergeElements(elements []Element) []Element {
	if len(elements) <= 1 {
		return elements
	}

	// First, merge drop caps
	elements = a.MergeDropCaps(elements)

	merged := make([]Element, 0, len(elements))
	current := elements[0]

	for i := 1; i < len(elements); i++ {
		next := elements[i]

		// Check if elements are on the same line
		yDiff := math.Abs(current.Y - next.Y)
		if yDiff < 2.0 {
			gap := next.X - (current.X + current.Width)

			// Use a more lenient threshold for merging
			// Small gaps (< 30% of font size) should merge without separator
			// This helps with kerned text and special styling
			threshold := current.FontSize * 0.3
			if threshold == 0 {
				threshold = 2.0
			}

			wideGapThreshold := 30.0
			if current.FontSize > 0 {
				wideGapThreshold = current.FontSize * 3.0
			}

			switch {
			case gap < threshold:
				current = a.mergeElementsWithLinks(current, next, "")
			case gap > wideGapThreshold:
				// Check if current is a bullet character that should merge with next
				isBullet := strings.TrimSpace(current.Content) == "•" ||
					strings.TrimSpace(current.Content) == "-" ||
					strings.TrimSpace(current.Content) == "▪" ||
					strings.TrimSpace(current.Content) == "◦"

				switch {
				case isBullet:
					// Bullet should always merge with following text
					current = a.mergeElementsWithLinks(current, next, " ")
				case len(current.Content) > 40 || len(next.Content) > 40:
					merged = append(merged, current)
					current = next
				default:
					isEqNum := a.isStartOfEqNum(next, elements[i+1:])

					switch {
					case isEqNum:
						current = a.mergeElementsWithLinks(current, next, " ")
					case strings.Contains(current.Content, "$") || strings.Contains(next.Content, "$"):
						current = a.mergeElementsWithLinks(current, next, "   ")
					default:
						current = a.mergeElementsWithLinks(current, next, "    ")
					}
				}
			default:
				if strings.HasSuffix(current.Content, "-") || strings.HasSuffix(current.Content, "‐") {
					current = a.mergeElementsWithLinks(current, next, "")
				} else {
					current = a.mergeElementsWithLinks(current, next, " ")
				}
			}
		} else {
			merged = append(merged, current)
			current = next
		}
	}

	merged = append(merged, current)
	return merged
}

// mergeElementsWithLinks merges two elements handling their LinkURIs
func (a *Analyzer) mergeElementsWithLinks(e1, e2 Element, separator string) Element {
	if e1.LinkURI == e2.LinkURI {
		e1.Content = mergeWithStyle(e1.Content, separator+e2.Content)
		e1.Width = e2.X + e2.Width - e1.X
		return e1
	}

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
	e1.LinkURI = ""

	return e1
}

// mergeWithStyle merges two strings, handling markdown style markers
func mergeWithStyle(a, b string) string {
	separator := ""
	cleanB := b
	if strings.HasPrefix(b, " ") {
		separator = " "
		cleanB = b[1:]
	}

	// Check bold
	if strings.HasSuffix(a, "**") && strings.HasPrefix(cleanB, "**") {
		return a[:len(a)-2] + separator + cleanB[2:]
	}
	// Check italic
	if strings.HasSuffix(a, "*") && strings.HasPrefix(cleanB, "*") {
		return a[:len(a)-1] + separator + cleanB[1:]
	}
	// Check math
	if strings.HasSuffix(a, "$") && strings.HasPrefix(cleanB, "$") {
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
				newElem := element
				current = &newElem
			} else {
				current.Content += "\n" + element.Content
				current.Height += element.Height
			}
		} else {
			if current != nil {
				merged = append(merged, *current)
				current = nil
			}
			merged = append(merged, element)
		}
	}

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
				newElem := element
				newElem.Content = strings.ReplaceAll(newElem.Content, "    ", " | ")
				newElem.Content = "| " + newElem.Content + " |"
				current = &newElem
			} else {
				row := strings.ReplaceAll(element.Content, "    ", " | ")
				row = strings.ReplaceAll(row, "   ", " | ")
				current.Content += "\n| " + row + " |"
				current.Height += element.Height
			}
		} else {
			if current != nil {
				current = a.validateTable(current)
				merged = append(merged, *current)
				current = nil
			}
			merged = append(merged, element)
		}
	}

	if current != nil {
		current = a.validateTable(current)
		merged = append(merged, *current)
	}

	return merged
}

// validateTable checks if a table is valid (has multiple columns)
func (a *Analyzer) validateTable(current *Element) *Element {
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
			trimmed = strings.TrimPrefix(trimmed, "| ")
			trimmed = strings.TrimSuffix(trimmed, " |")
			newContent.WriteString(trimmed)
		}
		current.Content = newContent.String()
		current.Type = ElementTypeParagraph
		if a.isNumberedHeader(current.Content) {
			current.Type = ElementTypeHeader
			current.Level = 2
		}
	}

	return current
}

// MergeListContinuations merges paragraph elements that follow list items into those list items
// This handles cases where multi-line list items are split across multiple elements
func (a *Analyzer) MergeListContinuations(elements []Element) []Element {
	if len(elements) <= 1 {
		return elements
	}

	merged := make([]Element, 0, len(elements))
	i := 0

	for i < len(elements) {
		current := elements[i]

		// If this is a list item, check for continuation paragraphs
		if current.Type == ElementTypeList {
			// Collect continuation lines
			j := i + 1
			for j < len(elements) {
				next := elements[j]

				// Stop if we hit another list item, header, code block, table, or image
				if next.Type != ElementTypeParagraph {
					break
				}

				// Stop if this paragraph starts with a list marker (it's a new item)
				if a.isListItem(next.Content) {
					break
				}

				// Check if this looks like a continuation based on content patterns
				if !a.isListContinuation(current.Content, next.Content) {
					break
				}

				// Merge the continuation into the list item
				content := strings.TrimSpace(next.Content)
				if content != "" {
					if !strings.HasSuffix(current.Content, " ") && !strings.HasSuffix(current.Content, "-") {
						current.Content += " "
					}
					current.Content += content
				}

				j++
			}

			merged = append(merged, current)
			i = j
		} else {
			merged = append(merged, current)
			i++
		}
	}

	return merged
}

// isListContinuation checks if nextContent looks like a continuation of currentContent
func (a *Analyzer) isListContinuation(currentContent, nextContent string) bool {
	current := strings.TrimSpace(currentContent)
	next := strings.TrimSpace(nextContent)

	if len(next) == 0 {
		return false
	}

	// Strip markdown formatting from the end of current for checking
	cleanCurrent := current
	if strings.HasSuffix(cleanCurrent, "**") {
		cleanCurrent = strings.TrimSuffix(cleanCurrent, "**")
	} else if strings.HasSuffix(cleanCurrent, "*") {
		cleanCurrent = strings.TrimSuffix(cleanCurrent, "*")
	}
	cleanCurrent = strings.TrimSpace(cleanCurrent)

	// Check if current ends without sentence-ending punctuation
	endsWithPunctuation := strings.HasSuffix(cleanCurrent, ".") ||
		strings.HasSuffix(cleanCurrent, "!") ||
		strings.HasSuffix(cleanCurrent, "?") ||
		strings.HasSuffix(cleanCurrent, ":")

	// Strip markdown formatting from start of next for checking
	cleanNext := next
	if strings.HasPrefix(cleanNext, "**") {
		cleanNext = strings.TrimPrefix(cleanNext, "**")
	} else if strings.HasPrefix(cleanNext, "*") {
		cleanNext = strings.TrimPrefix(cleanNext, "*")
	}
	cleanNext = strings.TrimSpace(cleanNext)

	// Get the first character of next content
	firstChar := rune(0)
	for _, r := range cleanNext {
		firstChar = r
		break
	}

	// Strong continuation indicator: current ends without punctuation and next starts lowercase
	if !endsWithPunctuation && unicode.IsLower(firstChar) {
		return true
	}

	// Also consider continuation if next is short and doesn't look like a new sentence
	// (e.g., "applications." could be the end of a multi-line list item)
	if len(cleanNext) < 100 && !endsWithPunctuation {
		// If current doesn't end with punctuation, continue until we find an ending
		return true
	}

	return false
}

// MergeParagraphLines merges consecutive paragraph lines
func (a *Analyzer) MergeParagraphLines(elements []Element) []Element {
	if len(elements) == 0 {
		return nil
	}

	// Calculate median gap
	var gaps []float64
	for i := 0; i < len(elements)-1; i++ {
		curr := elements[i]
		next := elements[i+1]
		if curr.Type == ElementTypeParagraph && next.Type == ElementTypeParagraph {
			gap := math.Abs(curr.Y - next.Y)
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
			// Skip duplicate
			if current.Content == next.Content {
				continue
			}

			gap := math.Abs(current.Y - next.Y)

			// Check for sentence break
			isSentenceBreak := false
			trimmedCurrent := strings.TrimSpace(current.Content)
			trimmedNext := strings.TrimSpace(next.Content)
			if (strings.HasSuffix(trimmedCurrent, ".") || strings.HasSuffix(trimmedCurrent, "?") || strings.HasSuffix(trimmedCurrent, "!")) &&
				len(trimmedNext) > 0 && unicode.IsUpper(rune(trimmedNext[0])) {
				isSentenceBreak = true
			}

			baseThreshold := current.FontSize * 1.6
			if medianGap > current.FontSize*1.5 {
				baseThreshold = medianGap * 1.2
			}

			threshold := baseThreshold
			if isSentenceBreak {
				threshold = current.FontSize * 1.5
			}

			isList := a.isListItem(next.Content)

			if gap < threshold && !isList {
				if !strings.HasSuffix(current.Content, " ") && !strings.HasSuffix(current.Content, "-") {
					current.Content += " "
				}
				current.Content += next.Content
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
