package layout

import (
	"fmt"

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
			// Check intersection
			// Link rect is [x1, y1, x2, y2] (bottom-left, top-right usually)
			// Element has X, Y, Width, Height (Y is usually bottom or top depending on coord system)
			// PDF coords: Y=0 at bottom.
			// Our extractor seems to normalize?
			// Let's assume standard PDF coords for Link Rect.
			// Our TextBlock Y is usually bottom-left of the text?
			// We need to be careful about coordinate systems.
			// The extractor (pdf/interpreter.go) converts coords.

			// Let's assume simple bounding box intersection for now.
			// Element: X, Y, Width, Height.
			// Link: Rect[0], Rect[1], Rect[2], Rect[3] (x1, y1, x2, y2)

			// Check if Element center is inside Link Rect
			cx := el.X + el.Width/2
			cy := el.Y + el.Height/2 // This depends on Y origin

			// If Y is top-down (like HTML/Canvas), then Y increases downwards.
			// If Y is bottom-up (PDF), then Y increases upwards.
			// Our extractor seems to keep PDF coords (Y=0 at bottom) or flips them?
			// Let's check interpreter.go or just try.
			// Usually PDF extractors normalize to top-left origin.
			// If so, Y increases downwards.

			// Let's assume the Link Rect matches the Element coordinate system
			// because both come from the same PDF reader/extractor.

			lx1, ly1, lx2, ly2 := link.Rect[0], link.Rect[1], link.Rect[2], link.Rect[3]

			// Normalize Link Rect (ensure x1<x2, y1<y2)
			if lx1 > lx2 {
				lx1, lx2 = lx2, lx1
			}
			if ly1 > ly2 {
				ly1, ly2 = ly2, ly1
			}

			// Check if element is roughly within link rect
			// Using center point check
			// But wait, a link might cover only PART of a text block if the block is a full line.
			// But we are processing Elements which are merged blocks.
			// Ideally we should do this BEFORE merging elements.
			// But `Analyze` calls `blocksToElements` then `MergeElements`.
			// If we do it after merging, we lose granularity.
			// A link might be on "click here" but the element is "Please click here for more".
			// If we just wrap the whole element, it becomes "[Please click here for more](url)".
			// This is acceptable for a first version, but not ideal.

			// Better approach:
			// If we want precise links, we need to apply this at `blocksToElements` stage
			// where blocks are smaller (often individual words or small phrases).
			// But `ProcessLinks` takes `[]Element`.
			// Maybe we should move this logic to `blocksToElements`?

			// For now, let's implement checking if the element is *mostly* covered by the link.
			// Or if the link is *inside* the element?

			// Intersection logic:
			// If intersection area / link area > 0.5 -> The link targets this element.
			// If intersection area / element area > 0.5 -> The element is the link.

			// Let's try a simple containment check of the center.
			if cx >= lx1 && cx <= lx2 && cy >= ly1 && cy <= ly2 {
				// Match!
				// Format as markdown link
				// Avoid double linking if already linked?
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

		cx := block.X + block.Width/2
		cy := block.Y + block.Height/2

		for _, link := range links {
			lx1, ly1, lx2, ly2 := link.Rect[0], link.Rect[1], link.Rect[2], link.Rect[3]

			// Normalize
			if lx1 > lx2 {
				lx1, lx2 = lx2, lx1
			}
			if ly1 > ly2 {
				ly1, ly2 = ly2, ly1
			}

			if cx >= lx1 && cx <= lx2 && cy >= ly1 && cy <= ly2 {
				// Set LinkURI instead of modifying Text
				block.LinkURI = link.URI
				break
			}
		}
	}

	return newBlocks
}
