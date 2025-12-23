package layout

import (
	"strings"
	"unicode"
)

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
	if strings.HasPrefix(text, "**") && strings.HasSuffix(text, "**") {
		innerText := strings.TrimPrefix(strings.TrimSuffix(text, "**"), "**")
		innerText = strings.TrimSpace(innerText)

		// Skip if it looks like a code line number
		if len(innerText) > 0 && innerText[0] >= '0' && innerText[0] <= '9' {
			// Don't treat as header
		} else if len(innerText) < 100 && len(innerText) > 0 &&
			!strings.HasSuffix(innerText, ".") &&
			!strings.HasSuffix(innerText, ",") &&
			!strings.HasSuffix(innerText, ";") &&
			!strings.Contains(innerText, "\n") &&
			len(innerText) > 2 {
			element.Type = ElementTypeHeader
			element.Level = 4
			element.Content = innerText
			return
		}
	}

	// 3. Check for implicit list items (indentation-based)
	if prev != nil {
		prevText := strings.TrimSpace(prev.Content)
		if strings.HasSuffix(prevText, ":") {
			if element.X > prev.X+5.0 {
				element.Type = ElementTypeList
				if prev.Type == ElementTypeList {
					element.Level = prev.Level + 1
				} else {
					element.Level = 1
				}
				return
			}
		}

		if prev.Type == ElementTypeList {
			if element.X-prev.X < 5.0 && element.X-prev.X > -5.0 {
				element.Type = ElementTypeList
				element.Level = prev.Level
				return
			}
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
	}

	marker := parts[0]
	if len(marker) == 0 {
		return false
	}

	// Check for Roman numeral headers
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

	hasPeriod := strings.HasSuffix(marker, ".")
	hasInternalDot := strings.Contains(strings.TrimSuffix(marker, "."), ".")

	if !hasPeriod && !hasInternalDot {
		return false
	}

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
	cleanText := text
	if strings.HasPrefix(cleanText, "**") {
		cleanText = strings.TrimPrefix(cleanText, "**")
	} else if strings.HasPrefix(cleanText, "*") {
		cleanText = strings.TrimPrefix(cleanText, "*")
	}

	// Check for specific prefixes
	prefixes := []string{
		"func ", "package ", "import ",
		"const ", "type ", "struct ", "interface ",
		"return ", "if ", "else ", "range ",
		"$ ", "# ",
		"kubectl ", "oc ", "docker ", "podman ", "sudo ",
		"apiVersion:", "kind:", "metadata:", "spec:", "status:",
	}

	for _, p := range prefixes {
		if strings.HasPrefix(cleanText, p) {
			if p == "$ " {
				if strings.HasSuffix(strings.TrimSpace(cleanText), "$") {
					continue
				}
			}
			if p == "if " {
				rest := strings.TrimPrefix(cleanText, "if ")
				if !strings.ContainsAny(rest, "({=><!") {
					continue
				}
			}
			return true
		}
	}

	// Check for unique substrings
	containsKeywords := []string{
		"apiVersion:", "kind: Pod", "kind: Service", "kind: Deployment",
		"namespace: ",
	}
	for _, kw := range containsKeywords {
		if strings.Contains(cleanText, kw) {
			return true
		}
	}

	// Check for exact matches of brackets
	if cleanText == "{" || cleanText == "}" || cleanText == "[]" || cleanText == "()" {
		return true
	}

	// Check for JSON/YAML patterns
	if strings.Contains(cleanText, "\":") {
		return true
	}
	if strings.HasPrefix(strings.TrimSpace(cleanText), "},") || strings.HasPrefix(strings.TrimSpace(cleanText), "],") {
		return true
	}
	if strings.HasSuffix(strings.TrimSpace(cleanText), "{") || strings.HasSuffix(strings.TrimSpace(cleanText), "[") {
		if strings.Contains(cleanText, ":") {
			return true
		}
	}

	// Check for SSH keys
	if strings.Contains(cleanText, "ssh-rsa ") {
		return true
	}

	// Check for lines that are just brackets/braces/commas
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

// isLikelyEquationText checks if text content looks like a mathematical equation
func (a *Analyzer) isLikelyEquationText(text string) bool {
	score := 0

	mathSymbols := []string{"=", "+", "−", "×", "÷", "∫", "∑", "∏", "√", "∂", "∇", "±", "≈", "≠", "≤", "≥", "∞", "→", "←", "↔", "⇒", "⇐", "⇔"}
	for _, sym := range mathSymbols {
		count := strings.Count(text, sym)
		if count > 0 {
			score += count
		}
	}

	dollarCount := strings.Count(text, "$")
	if dollarCount >= 2 {
		score += dollarCount
	}

	if hasEquationNumberAtEnd(text) {
		score += 3
	}

	if strings.Contains(text, "^") || strings.Contains(text, "_") {
		if !strings.Contains(text, "http") && !strings.Contains(text, "www") {
			score++
		}
	}

	greekLetters := []string{"α", "β", "γ", "δ", "ε", "ζ", "η", "θ", "ι", "κ", "λ", "μ", "ν", "ξ", "π", "ρ", "σ", "τ", "υ", "φ", "χ", "ψ", "ω", "Γ", "Δ", "Θ", "Λ", "Ξ", "Π", "Σ", "Φ", "Ψ", "Ω"}
	for _, letter := range greekLetters {
		if strings.Contains(text, letter) {
			score++
		}
	}

	if hasMathFractionPattern(text) {
		score++
	}

	tableHeaders := []string{"name", "value", "description", "type", "size", "date", "count", "total", "average", "status"}
	textLower := strings.ToLower(text)
	for _, header := range tableHeaders {
		if strings.Contains(textLower, header) {
			score--
		}
	}

	segments := strings.Split(text, "    ")
	if len(segments) > 3 {
		score -= 2
	}

	threshold := 2
	if a.Config != nil {
		threshold = a.Config.EquationScoreThreshold
	}
	return score >= threshold
}

// hasEquationNumberAtEnd checks if text ends with equation numbering like (1), (2.3), [5]
func hasEquationNumberAtEnd(text string) bool {
	text = strings.TrimSpace(text)
	if len(text) < 3 {
		return false
	}

	lastChar := text[len(text)-1]
	if lastChar != ')' && lastChar != ']' {
		return false
	}

	openChar := '('
	if lastChar == ']' {
		openChar = '['
	}

	searchStart := len(text) - 10
	if searchStart < 0 {
		searchStart = 0
	}

	for i := len(text) - 2; i >= searchStart; i-- {
		if rune(text[i]) == openChar {
			inner := text[i+1 : len(text)-1]
			if isNumericWithDots(inner) {
				return true
			}
			break
		}
	}

	return false
}

// hasMathFractionPattern checks for fraction-like patterns
func hasMathFractionPattern(text string) bool {
	for i := 1; i < len(text)-1; i++ {
		if text[i] == '/' {
			prevEnd := i
			prevStart := i - 1
			for prevStart > 0 && text[prevStart-1] != ' ' && text[prevStart-1] != '$' {
				prevStart--
			}
			prevToken := text[prevStart:prevEnd]

			nextStart := i + 1
			nextEnd := i + 1
			for nextEnd < len(text) && text[nextEnd] != ' ' && text[nextEnd] != '$' {
				nextEnd++
			}
			nextToken := text[nextStart:nextEnd]

			if len(prevToken) <= 3 && len(nextToken) <= 3 && len(prevToken) > 0 && len(nextToken) > 0 {
				return true
			}
		}
	}
	return false
}

// isNumericWithDots checks if string contains only digits and dots
func isNumericWithDots(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r != '.' && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

// isListItem checks if the block starts with a list marker
func (a *Analyzer) isListItem(text string) bool {
	if len(text) == 0 {
		return false
	}

	cleanText := text
	if strings.HasPrefix(cleanText, "**") {
		cleanText = strings.TrimPrefix(cleanText, "**")
	} else if strings.HasPrefix(cleanText, "*") {
		cleanText = strings.TrimPrefix(cleanText, "*")
	}

	if strings.HasPrefix(cleanText, "•") || strings.HasPrefix(cleanText, "- ") || strings.HasPrefix(cleanText, "* ") {
		return true
	}

	if len(cleanText) >= 3 && cleanText[0] >= '0' && cleanText[0] <= '9' && cleanText[1] == '.' && cleanText[2] == ' ' {
		return true
	}

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

	h1Ratio := 2.0
	h2Ratio := 1.75
	h3Ratio := 1.5
	h4Ratio := 1.3
	h5Ratio := 1.15

	if a.Config != nil {
		h1Ratio = a.Config.H1Ratio
		h2Ratio = a.Config.H2Ratio
		h3Ratio = a.Config.H3Ratio
		h4Ratio = a.Config.H4Ratio
		h5Ratio = a.Config.H5Ratio
	}

	if ratio >= h1Ratio {
		return 1
	} else if ratio >= h2Ratio {
		return 2
	} else if ratio >= h3Ratio {
		return 3
	} else if ratio >= h4Ratio {
		return 4
	} else if ratio >= h5Ratio {
		return 5
	}
	return 6
}

// isBold checks if font name indicates bold text
func isBold(fontName string) bool {
	lower := strings.ToLower(fontName)
	return strings.Contains(lower, "bold") || strings.Contains(lower, "black") || strings.Contains(lower, "heavy")
}

// isItalic checks if font name indicates italic text
func isItalic(fontName string) bool {
	lower := strings.ToLower(fontName)
	return strings.Contains(lower, "italic") || strings.Contains(lower, "oblique")
}

// isLikelyAuthorBlock detects author information blocks in academic papers
func (a *Analyzer) isLikelyAuthorBlock(text string) bool {
	textLower := strings.ToLower(text)

	// Check for email patterns
	if strings.Contains(text, "@") && (strings.Contains(textLower, ".edu") ||
		strings.Contains(textLower, ".com") ||
		strings.Contains(textLower, ".org") ||
		strings.Contains(textLower, ".ac.") ||
		strings.Contains(textLower, ".net")) {
		return true
	}

	// Check for common affiliation keywords
	affiliationKeywords := []string{
		"university", "institute", "department", "laboratory", "lab",
		"school of", "college of", "faculty of", "centre", "center",
		"research", "sciences", "engineering", "computer science",
	}
	for _, kw := range affiliationKeywords {
		if strings.Contains(textLower, kw) {
			return true
		}
	}

	// Check for superscript affiliation markers
	affiliationMarkers := []string{"¹", "²", "³", "⁴", "⁵", "†", "‡", "§", "∗", "⋆"}
	for _, marker := range affiliationMarkers {
		if strings.Contains(text, marker) {
			return true
		}
	}

	// Check for author name patterns with titles
	titlePatterns := []string{"dr.", "prof.", "ph.d", "m.d.", "mr.", "ms.", "mrs."}
	for _, title := range titlePatterns {
		if strings.Contains(textLower, title) {
			return true
		}
	}

	// Check for "and" between short segments
	segments := strings.Split(text, "    ")
	if len(segments) >= 2 {
		for _, seg := range segments {
			if strings.Contains(strings.ToLower(seg), " and ") {
				words := strings.Fields(seg)
				capitalizedCount := 0
				for _, word := range words {
					if len(word) > 0 && unicode.IsUpper(rune(word[0])) {
						capitalizedCount++
					}
				}
				if capitalizedCount > len(words)/2 {
					return true
				}
			}
		}
	}

	return false
}

// containsSignificantNumbers checks if text contains numbers that suggest tabular data
func containsSignificantNumbers(text string) bool {
	digitCount := 0
	consecutiveDigits := 0
	maxConsecutive := 0

	for _, r := range text {
		if unicode.IsDigit(r) {
			digitCount++
			consecutiveDigits++
			if consecutiveDigits > maxConsecutive {
				maxConsecutive = consecutiveDigits
			}
		} else {
			consecutiveDigits = 0
		}
	}

	return digitCount >= 2 && maxConsecutive >= 2
}
