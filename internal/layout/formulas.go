package layout

import (
	"strings"
	"unicode"
)

// isMathFont checks if the font name corresponds to a known mathematical font
func isMathFont(fontName string) bool {
	lower := strings.ToLower(fontName)

	// Common TeX/Math fonts
	mathFonts := []string{
		"cmmi", // Computer Modern Math Italic
		"cmsy", // Computer Modern Math Symbols
		"cmex", // Computer Modern Math Extension
		"msbm", // AMS Blackboard Bold
		"msam", // AMS Symbols A
		"symbol",
		"math",
		"eufm",       // Euler Fraktur
		"stix",       // STIX Fonts
		"cambria",    // Cambria Math
		"latinmodern", // Latin Modern
		"lmmath",     // Latin Modern Math
	}

	for _, mf := range mathFonts {
		if strings.Contains(lower, mf) {
			return true
		}
	}

	return false
}

// isMathSymbol checks if the text contains primarily math symbols
// This is a heuristic for when font information is missing or generic
func isMathSymbol(text string) bool {
	if len(text) == 0 {
		return false
	}

	mathSymbols := "∫∑∏√∂∇±≈≠≤≥∞→←↔⇒⇐⇔∈∉⊂⊃∪∩∀∃∅αβγδεζηθικλμνξπρστυφχψωΓΔΘΛΞΠΣΦΨΩ×÷°′″∝∴∵⊥∠∡∢"
	operatorSymbols := "=+-×÷<>≤≥≠≈"

	symbolCount := 0
	operatorCount := 0
	totalChars := 0

	for _, r := range text {
		if unicode.IsSpace(r) {
			continue
		}
		totalChars++

		if strings.ContainsRune(mathSymbols, r) {
			symbolCount++
		}
		if strings.ContainsRune(operatorSymbols, r) {
			operatorCount++
		}
	}

	if totalChars == 0 {
		return false
	}

	// If more than 30% of characters are math symbols, it's likely math
	symbolRatio := float64(symbolCount) / float64(totalChars)
	if symbolRatio > 0.3 {
		return true
	}

	// If it's a single character that's a common math variable (single letter with operators nearby)
	if totalChars <= 3 && operatorCount > 0 {
		return true
	}

	return false
}

// CleanMathContent attempts to clean up fragmented math expressions
// e.g., "$x$$y$" -> "$xy$" and "$a$ = $b$" -> "$a = b$"
func CleanMathContent(text string) string {
	// First pass: merge adjacent $...$ blocks that were incorrectly separated
	// Pattern: $a$$b$ -> $ab$
	for strings.Contains(text, "$$") {
		text = strings.ReplaceAll(text, "$$", "$")
	}

	// Second pass: merge math blocks separated by just whitespace or operators
	// Pattern: $a$ $b$ -> $a b$ if they look like they should be together
	// This is tricky - we need to identify equation patterns

	// Split by $ to find math segments
	parts := strings.Split(text, "$")
	if len(parts) < 3 {
		return text
	}

	var result strings.Builder
	inMath := false
	pendingMath := ""

	for i, part := range parts {
		if i == 0 {
			// First part is always text (before first $)
			result.WriteString(part)
			inMath = true
			continue
		}

		if inMath {
			// This is math content
			if pendingMath == "" {
				pendingMath = part
			} else {
				// Check if previous separator was just whitespace or operator
				// In that case, merge
				pendingMath += part
			}
			inMath = false
		} else {
			// This is text between math blocks
			trimmed := strings.TrimSpace(part)

			// Check if this looks like it should be inside math
			// (operators, short connectors, or empty)
			if shouldMergeWithMath(trimmed) && len(parts) > i+1 {
				// Merge with next math block
				pendingMath += " " + part + " "
				inMath = true
			} else {
				// Output pending math and this text
				if pendingMath != "" {
					result.WriteString("$")
					result.WriteString(strings.TrimSpace(pendingMath))
					result.WriteString("$")
					pendingMath = ""
				}
				result.WriteString(part)
				inMath = true
			}
		}
	}

	// Flush pending math
	if pendingMath != "" {
		result.WriteString("$")
		result.WriteString(strings.TrimSpace(pendingMath))
		result.WriteString("$")
	}

	return result.String()
}

// shouldMergeWithMath checks if a text segment between math blocks
// should be merged into the math expression
func shouldMergeWithMath(text string) bool {
	text = strings.TrimSpace(text)

	// Empty or whitespace only -> merge
	if text == "" {
		return true
	}

	// Single operator -> merge
	operators := []string{"=", "+", "-", "×", "÷", "<", ">", "≤", "≥", "≠", "≈", ",", ";", "(", ")", "[", "]", "{", "}"}
	for _, op := range operators {
		if text == op {
			return true
		}
	}

	// Short text that looks like equation continuation
	// e.g., "and", "or", "where", "for"
	mathWords := []string{"and", "or", "where", "for", "if", "then", "with"}
	lower := strings.ToLower(text)
	for _, word := range mathWords {
		if lower == word {
			return true
		}
	}

	return false
}
