package layout

import (
	"math"
	"strings"
	"unicode"
)

// CleanMathSymbols merges consecutive duplicate math symbols
func (a *Analyzer) CleanMathSymbols(elements []Element) []Element {
	for i := range elements {
		if elements[i].Type == ElementTypeParagraph || elements[i].Type == ElementTypeHeader {
			text := elements[i].Content

			parts := strings.Split(text, "$")
			if len(parts) < 3 {
				continue
			}

			var newParts []string
			newParts = append(newParts, parts[0])

			for j := 1; j < len(parts); j += 2 {
				mathContent := parts[j]

				if j+2 < len(parts) {
					separator := parts[j+1]
					nextMathContent := parts[j+2]

					isSeparatorClean := strings.TrimSpace(separator) == ""

					if isSeparatorClean && mathContent == nextMathContent {
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
							newParts = append(newParts, parts[k-1])
						}

						j = k - 2
						continue
					}
				}

				newParts = append(newParts, mathContent)
				if j+1 < len(parts) {
					newParts = append(newParts, parts[j+1])
				}
			}

			var sb strings.Builder
			for k, p := range newParts {
				sb.WriteString(p)
				if k < len(newParts)-1 {
					sb.WriteString("$")
				}
			}
			elements[i].Content = sb.String()
		}
	}
	return elements
}

// isStartOfEqNum checks if element forms an equation number like "(1)" or "(2.1)"
func (a *Analyzer) isStartOfEqNum(current Element, remaining []Element) bool {
	text := strings.TrimSpace(current.Content)
	if !strings.HasPrefix(text, "(") {
		return false
	}

	if strings.Contains(text, ")") {
		start := strings.Index(text, "(")
		end := strings.LastIndex(text, ")")
		if start < end {
			inner := text[start+1 : end]
			for _, r := range inner {
				if !unicode.IsDigit(r) && r != '.' && r != ' ' {
					return false
				}
			}
			return true
		}
		return false
	}

	acc := text
	for _, el := range remaining {
		if math.Abs(el.Y-current.Y) > 2.0 {
			break
		}

		content := strings.TrimSpace(el.Content)
		acc += content
		if strings.Contains(content, ")") {
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

		if len(acc) > 20 {
			return false
		}
	}

	return false
}

// CleanFragmentedMath applies CleanMathContent to fix fragmented math expressions
func (a *Analyzer) CleanFragmentedMath(elements []Element) []Element {
	for i := range elements {
		if elements[i].Type == ElementTypeParagraph || elements[i].Type == ElementTypeHeader {
			elements[i].Content = CleanMathContent(elements[i].Content)
		}
	}
	return elements
}
