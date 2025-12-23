package layout

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode"
)

// MergeElements merges consecutive text elements on the same line
func (a *Analyzer) MergeElements(elements []Element) []Element {
	if len(elements) <= 1 {
		return elements
	}

	merged := make([]Element, 0, len(elements))
	current := elements[0]

	for i := 1; i < len(elements); i++ {
		next := elements[i]

		// Check if elements are on the same line
		yDiff := math.Abs(current.Y - next.Y)
		if yDiff < 2.0 {
			gap := next.X - (current.X + current.Width)

			threshold := current.FontSize * 0.1
			if threshold == 0 {
				threshold = 1.0
			}

			wideGapThreshold := 30.0
			if current.FontSize > 0 {
				wideGapThreshold = current.FontSize * 3.0
			}

			switch {
			case gap < threshold:
				current = a.mergeElementsWithLinks(current, next, "")
			case gap > wideGapThreshold:
				if len(current.Content) > 40 || len(next.Content) > 40 {
					merged = append(merged, current)
					current = next
				} else {
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
