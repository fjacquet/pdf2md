package layout

import (
	"math"
	"strings"
)

// MergeSameLinkElements merges consecutive elements that share the same LinkURI
func (a *Analyzer) MergeSameLinkElements(elements []Element) []Element {
	merged := make([]Element, 0, len(elements))
	current := elements[0]

	for i := 1; i < len(elements); i++ {
		next := elements[i]

		// Check if elements are on the same line (similar Y coordinate)
		yDiff := math.Abs(current.Y - next.Y)

		// Check if they have the same LinkURI (and it's not empty)
		sameLink := current.LinkURI != "" && current.LinkURI == next.LinkURI

		if yDiff < 2.0 && sameLink {
			// Merge them!
			// Calculate gap
			gap := next.X - (current.X + current.Width)

			// If gap is small, merge without space or with space depending on context
			// Since they are part of the same link, they are likely part of the same text flow.
			// Use standard merging logic for content.

			// Check for Hyphenation (unlikely in links but possible)
			if strings.HasSuffix(current.Content, "-") || strings.HasSuffix(current.Content, "‐") {
				current.Content = mergeWithStyle(current.Content, next.Content)
			} else {
				// If gap is very small, merge directly. Otherwise space.
				// Links usually don't have wide gaps inside them unless it's multiple words.
				// Let's assume space if gap > 0.5 char width
				threshold := current.FontSize * 0.2
				if gap > threshold {
					current.Content = mergeWithStyle(current.Content, " "+next.Content)
				} else {
					current.Content = mergeWithStyle(current.Content, next.Content)
				}
			}

			// Update width
			current.Width = next.X + next.Width - current.X
		} else {
			// Save current and move to next
			merged = append(merged, current)
			current = next
		}
	}

	// Add the last element
	merged = append(merged, current)

	return merged
}
