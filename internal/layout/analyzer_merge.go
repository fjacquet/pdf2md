package layout

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/fjacquet/pdf2md/internal/types"
)

// softHyphen marks a line-end hyphen that may be a word break, settled by
// ResolveWordBreaks once the whole document is known.
const softHyphen = "\u00AD"

// joinWrapped joins a line to the next wrapped line. A line-end hyphen
// (-, U+2010, U+2011, U+00AD) between lowercase letters becomes a soft hyphen
// ("impor-" + "tant"); other hyphens stay, unspaced ("2-3").
func joinWrapped(cur, next string) string {
	for _, h := range []string{"-", "\u2010", "\u2011", softHyphen} {
		if word, ok := strings.CutSuffix(cur, h); ok {
			last, _ := utf8.DecodeLastRuneInString(word)
			first, _ := utf8.DecodeRuneInString(next)
			if unicode.IsLower(last) && unicode.IsLower(first) {
				return word + softHyphen + next
			}
			return cur + next
		}
	}
	if strings.HasSuffix(cur, " ") {
		return mergeWithStyle(cur, next)
	}
	return mergeWithStyle(cur, " "+next)
}

var brokenWord = regexp.MustCompile(`\pL+(?:\x{00AD}\pL+)+`)

// ResolveWordBreaks settles the soft hyphens joinWrapped left: a word the
// document also spells unbroken loses its hyphen ("important"), any other
// keeps it as a compound ("well-studied").
func (a *Analyzer) ResolveWordBreaks(elements []Element) []Element {
	known := map[string]bool{}
	for _, e := range elements {
		for _, w := range strings.FieldsFunc(e.Content, func(r rune) bool { return !unicode.IsLetter(r) }) {
			known[strings.ToLower(w)] = true
		}
	}
	for i := range elements {
		elements[i].Content = brokenWord.ReplaceAllStringFunc(elements[i].Content, func(w string) string {
			if joined := strings.ReplaceAll(w, softHyphen, ""); known[strings.ToLower(joined)] {
				return joined
			}
			return strings.ReplaceAll(w, softHyphen, "-")
		})
	}
	return elements
}

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
		if len(content) <= 2 && i+1 < len(elements) && current.Type != ElementTypeImage && elements[i+1].Type != ElementTypeImage {
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

		// An image is never text: no line merge, no script.
		if current.Type == ElementTypeImage || next.Type == ElementTypeImage {
			merged = append(merged, current)
			current = next
			continue
		}

		// A superscript or subscript joins its line as $^{…}$ / $_{…}$.
		if mark := scriptMark(current, next); mark != "" {
			current = a.mergeElementsWithLinks(current, asScript(next, mark), "")
			current.Content = joinScripts.Replace(current.Content)
			continue
		}
		if mark := scriptMark(next, current); mark != "" {
			host := next
			current = a.mergeElementsWithLinks(asScript(current, mark), next, "")
			current.Y, current.Height, current.FontSize, current.SpaceWidth, current.Bold = host.Y, host.Height, host.FontSize, host.SpaceWidth, host.Bold
			continue
		}

		// Check if elements are on the same line
		yDiff := math.Abs(current.Y - next.Y)
		if yDiff < 2.0 {
			gap := next.X - (current.X + current.Width)

			sw := current.SpaceWidth
			if sw == 0 {
				// OCR blocks carry no font metrics: assume a 0.25 em space.
				sw = current.FontSize * 0.25
			}
			if sw == 0 {
				sw = 4.0 // no font size either: the former 2 pt threshold
			}

			wideGapThreshold := 30.0
			if current.FontSize > 0 {
				wideGapThreshold = current.FontSize * 3.0
			}

			switch {
			case !types.IsWordGap(gap, sw):
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
				if strings.HasSuffix(current.Content, "-") || strings.HasSuffix(current.Content, "\u2010") {
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

// scriptMark returns "^" when s is a superscript of host, "_" when it is a
// subscript, "" otherwise: a smaller font (at most 0.85x), its baseline
// 0.12-0.45 em off host's, and touching host horizontally (within 1 em).
func scriptMark(host, s Element) string {
	em := host.FontSize
	if em <= 0 || s.FontSize > 0.85*em || scriptText(s.Content) == "" {
		return ""
	}
	dy := s.Y - host.Y
	if math.Abs(dy) < 0.12*em || math.Abs(dy) > 0.45*em {
		return ""
	}
	gap := max(s.X-(host.X+host.Width), host.X-(s.X+s.Width))
	if gap > em || gap < -0.5*em {
		return ""
	}
	if dy > 0 {
		return "^"
	}
	return "_"
}

// joinScripts joins a script drawn in pieces: "^{−}^{1}" -> "^{−1}".
var joinScripts = strings.NewReplacer("}^{", "", "}_{", "")

// asScript rewrites e's text as a LaTeX script: "9" -> "$^{9}$".
func asScript(e Element, mark string) Element {
	e.Content = "$" + mark + "{" + scriptText(e.Content) + "}$"
	return e
}

// scriptText is a script's printable text without math markers; glyphs with
// no text mapping (control characters) leave it empty.
func scriptText(s string) string {
	return strings.TrimSpace(strings.Map(func(r rune) rune {
		if r == '$' || !unicode.IsGraphic(r) {
			return -1
		}
		return r
	}, s))
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
	// A space drawn as its own block trails a: "**a** " + "**b**".
	cleanA := strings.TrimRight(a, " ")
	if cleanA != a {
		separator = " "
	}

	// Check bold
	if strings.HasSuffix(cleanA, "**") && strings.HasPrefix(cleanB, "**") {
		return cleanA[:len(cleanA)-2] + separator + cleanB[2:]
	}
	// Check italic
	if strings.HasSuffix(cleanA, "*") && strings.HasPrefix(cleanB, "*") {
		return cleanA[:len(cleanA)-1] + separator + cleanB[1:]
	}
	// Check math
	if strings.HasSuffix(cleanA, "$") && strings.HasPrefix(cleanB, "$") {
		return cleanA[:len(cleanA)-1] + separator + cleanB[1:]
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
					current.Content = joinWrapped(current.Content, content)
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
	// The gap to the next line is measured from current's last line, while
	// current.Y stays the paragraph's top.
	lastY := func(e *Element) float64 {
		if e.LastY != 0 {
			return e.LastY
		}
		return e.Y
	}

	for i := 1; i < len(elements); i++ {
		next := elements[i]

		if current.Type == ElementTypeParagraph && next.Type == ElementTypeParagraph {
			// Skip duplicate
			if current.Content == next.Content {
				continue
			}

			gap := math.Abs(lastY(current) - next.Y)

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
				current.Content = joinWrapped(current.Content, next.Content)
				current.LastY = lastY(&next)
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
