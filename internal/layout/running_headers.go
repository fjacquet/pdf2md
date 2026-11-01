package layout

import (
	"strings"
	"unicode"
)

// FilterRunningHeaders drops running headers and footers and flattens the
// pages. A paragraph that is a page's top or bottom element goes when the same
// text, ignoring digits, spaces and case ("Page 3 of 9"), sits at a page edge
// on at least 3 pages and a third of them (so even/odd headers count too).
// Headings stay: a title on page 1 may repeat as the running header.
func (a *Analyzer) FilterRunningHeaders(pages [][]Element) []Element {
	key := func(s string) string {
		return strings.Map(func(r rune) rune {
			if unicode.IsDigit(r) || unicode.IsSpace(r) {
				return -1
			}
			return unicode.ToLower(r)
		}, s)
	}
	// edges[i] holds the indexes of page i's top and bottom text elements.
	edges := make([][]int, len(pages))
	counts := map[string]int{}
	for i, page := range pages {
		top, bottom := -1, -1
		for j, e := range page {
			if e.Type == ElementTypeImage || strings.TrimSpace(e.Content) == "" {
				continue
			}
			if top < 0 || e.Y > page[top].Y {
				top = j
			}
			if bottom < 0 || e.Y < page[bottom].Y {
				bottom = j
			}
		}
		if top < 0 {
			continue
		}
		edges[i] = []int{top}
		if bottom != top {
			edges[i] = append(edges[i], bottom)
		}
		seen := map[string]bool{}
		for _, j := range edges[i] {
			if k := key(page[j].Content); k != "" && !seen[k] {
				seen[k] = true
				counts[k]++
			}
		}
	}
	var out []Element
	for i, page := range pages {
		for j, e := range page {
			drop := false
			if e.Type == ElementTypeParagraph {
				for _, edge := range edges[i] {
					if c := counts[key(e.Content)]; edge == j && c >= 3 && 3*c >= len(pages) {
						drop = true
					}
				}
			}
			if !drop {
				out = append(out, e)
			}
		}
	}
	return out
}
