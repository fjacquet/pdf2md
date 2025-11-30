package layout

import "strings"

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
		"eufm", // Euler Fraktur
		"stix", // STIX Fonts
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
	// TODO: Implement symbol density check if needed
	return false
}
