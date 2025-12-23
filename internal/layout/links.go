package layout

import (
	"fmt"
	"math"

	"github.com/fjacquet/pdf2md/internal/extractor"
	"github.com/fjacquet/pdf2md/internal/types"
)

// ProcessLinks maps links to text blocks and formats them as Markdown links
func ProcessLinks(elements []Element, links []types.Link) []Element {
	if len(links) == 0 {
		return elements
	}

	// For each element, check if it intersects with any link
	for i := range elements {
		el := &elements[i]

		// Only process paragraphs and headers for now
		if el.Type != ElementTypeParagraph && el.Type != ElementTypeHeader && el.Type != ElementTypeList {
			continue
		}

		// Check all links
		for _, link := range links {
			// Element bounding box
			ex1, ey1 := el.X, el.Y
			ex2, ey2 := el.X+el.Width, el.Y+el.Height

			// Link rect (x1, y1, x2, y2)
			lx1, ly1, lx2, ly2 := link.Rect[0], link.Rect[1], link.Rect[2], link.Rect[3]

			// Normalize Link Rect (ensure x1<x2, y1<y2)
			if lx1 > lx2 {
				lx1, lx2 = lx2, lx1
			}
			if ly1 > ly2 {
				ly1, ly2 = ly2, ly1
			}

			// Calculate intersection rectangle
			ix1 := math.Max(ex1, lx1)
			iy1 := math.Max(ey1, ly1)
			ix2 := math.Min(ex2, lx2)
			iy2 := math.Min(ey2, ly2)

			// Check if there's a valid intersection
			if ix1 >= ix2 || iy1 >= iy2 {
				continue // No intersection
			}

			// Calculate areas
			intersectionArea := (ix2 - ix1) * (iy2 - iy1)
			elementArea := el.Width * el.Height
			linkArea := (lx2 - lx1) * (ly2 - ly1)

			// Match if intersection covers significant portion of either rect
			// - Link is mostly inside element (link is anchor text)
			// - Element is mostly inside link (element is clickable)
			linkCoverage := 0.0
			elementCoverage := 0.0
			if linkArea > 0 {
				linkCoverage = intersectionArea / linkArea
			}
			if elementArea > 0 {
				elementCoverage = intersectionArea / elementArea
			}

			// Match if >30% of either rectangle is covered
			// (lower threshold than 50% to catch partial links)
			if linkCoverage > 0.3 || elementCoverage > 0.3 {
				// Match! Format as markdown link
				el.Content = fmt.Sprintf("[%s](%s)", el.Content, link.URI)
				break // Only one link per element for now
			}
		}
	}

	return elements
}

// ApplyLinksToBlocks applies links to raw text blocks before they are merged into elements.
// This allows for more granular linking (e.g. linking just "here" in "Click here").
func ApplyLinksToBlocks(blocks []extractor.TextBlock, links []types.Link) []extractor.TextBlock {
	if len(links) == 0 {
		return blocks
	}

	newBlocks := make([]extractor.TextBlock, len(blocks))
	copy(newBlocks, blocks)

	for i := range newBlocks {
		block := &newBlocks[i]

		// Block bounding box
		bx1, by1 := block.X, block.Y
		bx2, by2 := block.X+block.Width, block.Y+block.Height

		for _, link := range links {
			lx1, ly1, lx2, ly2 := link.Rect[0], link.Rect[1], link.Rect[2], link.Rect[3]

			// Normalize link rect
			if lx1 > lx2 {
				lx1, lx2 = lx2, lx1
			}
			if ly1 > ly2 {
				ly1, ly2 = ly2, ly1
			}

			// Calculate intersection rectangle
			ix1 := math.Max(bx1, lx1)
			iy1 := math.Max(by1, ly1)
			ix2 := math.Min(bx2, lx2)
			iy2 := math.Min(by2, ly2)

			// Check if there's a valid intersection
			if ix1 >= ix2 || iy1 >= iy2 {
				continue // No intersection
			}

			// Calculate areas
			intersectionArea := (ix2 - ix1) * (iy2 - iy1)
			blockArea := block.Width * block.Height
			linkArea := (lx2 - lx1) * (ly2 - ly1)

			// Calculate coverage ratios
			linkCoverage := 0.0
			blockCoverage := 0.0
			if linkArea > 0 {
				linkCoverage = intersectionArea / linkArea
			}
			if blockArea > 0 {
				blockCoverage = intersectionArea / blockArea
			}

			// Match if >30% of either rectangle is covered
			if linkCoverage > 0.3 || blockCoverage > 0.3 {
				block.LinkURI = link.URI
				break
			}
		}
	}

	return newBlocks
}
