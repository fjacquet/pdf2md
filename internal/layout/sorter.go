package layout

import (
	"math"
	"sort"

	"github.com/fjacquet/pdf2md/internal/extractor"
)

// SortBlocks sorts text blocks using the Recursive XY-Cut algorithm
// This handles multi-column layouts by recursively splitting the page into
// horizontal rows and vertical columns.
func SortBlocks(blocks []extractor.TextBlock) []extractor.TextBlock {
	if len(blocks) <= 1 {
		return blocks
	}

	// Initial bounds
	var minX, minY, maxX, maxY float64
	minX = blocks[0].X
	maxX = blocks[0].X + blocks[0].Width
	minY = blocks[0].Y
	maxY = blocks[0].Y + blocks[0].Height

	for _, b := range blocks {
		if b.X < minX {
			minX = b.X
		}
		if b.X+b.Width > maxX {
			maxX = b.X + b.Width
		}
		if b.Y < minY {
			minY = b.Y
		}
		if b.Y+b.Height > maxY {
			maxY = b.Y + b.Height
		}
	}

	// 1. Recursive XY-Cut
	// We start by trying to detect columns. If that fails (e.g. headers),
	// it will fallback to row splitting.
	return detectAndSortColumns(blocks)
}

func splitByRows(blocks []extractor.TextBlock, gapThreshold float64) [][]extractor.TextBlock {
	if len(blocks) == 0 {
		return nil
	}

	// Sort by Top Y Descending
	sorted := make([]extractor.TextBlock, len(blocks))
	copy(sorted, blocks)
	sort.Slice(sorted, func(i, j int) bool {
		return (sorted[i].Y + sorted[i].Height) > (sorted[j].Y + sorted[j].Height)
	})

	var rows [][]extractor.TextBlock
	var currentRow []extractor.TextBlock

	currentRow = append(currentRow, sorted[0])
	currentBottom := sorted[0].Y

	for i := 1; i < len(sorted); i++ {
		b := sorted[i]
		top := b.Y + b.Height

		// Check for gap
		// If the top of the current block is significantly below the bottom of the previous cluster
		if top < currentBottom-gapThreshold {
			// Gap found!
			rows = append(rows, currentRow)
			currentRow = []extractor.TextBlock{b}
			currentBottom = b.Y
		} else {
			// Overlap or close enough
			currentRow = append(currentRow, b)
			if b.Y < currentBottom {
				currentBottom = b.Y
			}
		}
	}
	rows = append(rows, currentRow)

	return rows
}

func detectAndSortColumns(blocks []extractor.TextBlock) []extractor.TextBlock {
	if len(blocks) <= 1 {
		return blocks
	}

	// 1. Calculate Bounds
	var minX, maxX float64
	minX = blocks[0].X
	maxX = blocks[0].X + blocks[0].Width
	for _, b := range blocks {
		if b.X < minX {
			minX = b.X
		}
		if b.X+b.Width > maxX {
			maxX = b.X + b.Width
		}
	}
	width := maxX - minX
	// If width is too small, don't try to split (avoid splitting columns into words)
	if width < 200.0 {
		return detectAndSortRows(blocks)
	}

	// 2. Find the best vertical gutter (X-Cut)
	// We look for a vertical strip with minimal text density.
	// We use a histogram approach.
	minGutterWidthBins := 8 // Reduced to 8 (approx 16 units) to detect narrower gaps
	if int(maxX-minX)/2 < minGutterWidthBins {
		return detectAndSortRows(blocks)
	}

	hist := make([]int, int((maxX-minX)/2.0)+1)
	for _, b := range blocks {
		startBin := int((b.X - minX) / 2.0)
		endBin := int((b.X + b.Width - minX) / 2.0)
		for i := startBin; i <= endBin; i++ {
			if i >= 0 && i < len(hist) {
				hist[i]++
			}
		}
	}

	// Find the window of size minGutterWidthBins with minimum sum
	minWindowDensity := 999999
	bestWindowStart := -1

	for i := 0; i <= len(hist)-minGutterWidthBins; i++ {
		sum := 0
		for j := 0; j < minGutterWidthBins; j++ {
			sum += hist[i+j]
		}
		if sum < minWindowDensity {
			minWindowDensity = sum
			bestWindowStart = i
		}
	}

	// If the best window has high density, abort X-Cut and try Y-Cut.
	// High density means many blocks cross this "gutter".
	// If density > len(blocks) * 0.1, it's probably not a gutter.
	// If the best window has high density, check if it's just Header/Footer.
	// If the blocks in the gutter leave a large vertical gap (e.g. > 50% of height), it's a valid column split.
	if minWindowDensity > int(float64(len(blocks))*0.1) {
		// Check for vertical gap in gutter
		gutterBlocks := make([]extractor.TextBlock, 0)
		gutterStart := minX + float64(bestWindowStart)*2.0
		gutterEnd := gutterStart + float64(minGutterWidthBins)*2.0

		for _, b := range blocks {
			// Check if block overlaps with gutter
			if b.X+b.Width > gutterStart && b.X < gutterEnd {
				gutterBlocks = append(gutterBlocks, b)
			}
		}

		if len(gutterBlocks) > 0 {
			// Sort by Y
			sort.Slice(gutterBlocks, func(i, j int) bool {
				return gutterBlocks[i].Y < gutterBlocks[j].Y
			})

			// Find max gap
			maxGap := 0.0
			// Check gap from top of page to first block
			minY := blocks[0].Y
			maxY := blocks[0].Y + blocks[0].Height
			for _, b := range blocks {
				if b.Y < minY {
					minY = b.Y
				}
				if b.Y+b.Height > maxY {
					maxY = b.Y + b.Height
				}
			}

			if len(gutterBlocks) > 0 {
				maxGap = math.Max(maxGap, gutterBlocks[0].Y-minY)
				for i := 0; i < len(gutterBlocks)-1; i++ {
					gap := gutterBlocks[i+1].Y - (gutterBlocks[i].Y + gutterBlocks[i].Height)
					if gap > maxGap {
						maxGap = gap
					}
				}
				maxGap = math.Max(maxGap, maxY-(gutterBlocks[len(gutterBlocks)-1].Y+gutterBlocks[len(gutterBlocks)-1].Height))
			} else {
				maxGap = maxY - minY
			}

			pageHeight := maxY - minY
			// If max gap is less than 50% of page height, then the gutter is blocked by content throughout.
			if maxGap < pageHeight*0.5 {
				return detectAndSortRows(blocks)
			}
		}
	}

	// Center of the best window
	bestGutterBin := bestWindowStart + minGutterWidthBins/2

	gutterX := minX + float64(bestGutterBin)*2.0

	// 3. Identify Spanning Blocks vs Left/Right
	var left, right, spanning []extractor.TextBlock

	for _, b := range blocks {
		if b.X < gutterX && b.X+b.Width > gutterX {
			spanning = append(spanning, b)
		} else if b.X+b.Width <= gutterX {
			left = append(left, b)
		} else {
			right = append(right, b)
		}
	}

	// If one side is empty, it's single column (or gutter is at edge)
	if len(left) == 0 || len(right) == 0 {
		return detectAndSortRows(blocks)
	}

	// 4. Balanced Column Check (Stability Fix)
	// If one column is very narrow compared to the total width, it's likely a list or indentation, not a column.
	// Real columns usually take up at least 30% of the width each (for 2 columns).
	// We'll be conservative and say if a column is less than 20% of the total width, it's not a column.

	getGroupWidth := func(blks []extractor.TextBlock) float64 {
		if len(blks) == 0 {
			return 0
		}
		minX := blks[0].X
		maxX := blks[0].X + blks[0].Width
		for _, b := range blks {
			if b.X < minX {
				minX = b.X
			}
			if b.X+b.Width > maxX {
				maxX = b.X + b.Width
			}
		}
		return maxX - minX
	}

	leftWidth := getGroupWidth(left)
	rightWidth := getGroupWidth(right)

	// If either side is too narrow (e.g. < 40% of total width), treat as single column
	// Real columns (2 cols) are usually ~50%. Rivers are usually off-center.
	if leftWidth < width*0.40 || rightWidth < width*0.40 {
		return detectAndSortRows(blocks)
	}

	// 6. Text Flow Check (Semantic Safeguard)
	// Check if text appears to flow vertically (Columnar) or Horizontally/Independently (Table/Code).
	// We check if lines in Left column continue to the next line in Left column (hyphens, lowercase).
	verticalFlowScore := checkVerticalFlow(left)
	horizontalFlowScore := checkHorizontalFlow(left, right)

	// Check for row alignment
	alignmentScore := checkRowAlignment(left, right)

	// If Horizontal Flow is significant (>= Vertical * 2), it's likely a single column.
	if horizontalFlowScore > verticalFlowScore*2 && verticalFlowScore < 5 {
		return detectAndSortRows(blocks)
	}

	// Also keep the alignment check for Tables/Code that might not have flow (e.g. numbers)
	if verticalFlowScore == 0 && alignmentScore > 0.8 {
		return detectAndSortRows(blocks)
	}

	// 7. Process Segments
	// Sort each group
	sortedLeft := detectAndSortColumns(left)      // Recurse
	sortedRight := detectAndSortColumns(right)    // Recurse
	sortedSpanning := detectAndSortRows(spanning) // Spanning blocks might need Y-sorting/splitting

	return append(append(sortedLeft, sortedRight...), sortedSpanning...)
}

func detectAndSortRows(blocks []extractor.TextBlock) []extractor.TextBlock {
	if len(blocks) <= 1 {
		return blocks
	}

	// 1. Calculate Bounds
	minY, maxY := 99999.0, -99999.0
	for _, b := range blocks {
		if b.Y < minY {
			minY = b.Y
		}
		if b.Y+b.Height > maxY {
			maxY = b.Y + b.Height
		}
	}

	// 2. Build Histogram of "crossing" blocks (Y-axis)
	// Bin size 2.0
	numBins := int((maxY - minY) / 2.0)
	if numBins < 1 {
		return sortY(blocks)
	}
	histogram := make([]int, numBins)

	for _, b := range blocks {
		startBin := int((b.Y - minY) / 2.0)
		endBin := int((b.Y + b.Height - minY) / 2.0)

		if startBin < 0 {
			startBin = 0
		}
		if endBin >= numBins {
			endBin = numBins - 1
		}

		for i := startBin; i <= endBin; i++ {
			histogram[i]++
		}
	}

	// Find the BEST horizontal gutter (lowest density)
	// We want a gap where NO blocks cross (density 0) if possible.
	// Or minimal density.
	minGutterHeightBins := 4 // 8 units height
	minWindowDensity := 999999
	bestWindowStart := -1

	for i := 0; i <= numBins-minGutterHeightBins; i++ {
		currentDensity := 0
		for j := 0; j < minGutterHeightBins; j++ {
			currentDensity += histogram[i+j]
		}

		if currentDensity < minWindowDensity {
			minWindowDensity = currentDensity
			bestWindowStart = i
		}
	}

	// If density is high, we can't split horizontally.
	// Threshold: 0? We really want clear gaps for rows.
	// But maybe some noise?
	// Let's use strict 0 for now, or very low.
	if minWindowDensity > 0 {
		return sortY(blocks)
	}

	// Split Upper / Lower (Y-axis)
	// PDF Y increases upwards.
	// Lower Y = Bottom of page.
	// Higher Y = Top of page.
	bestGutterBin := bestWindowStart + minGutterHeightBins/2
	gutterY := minY + float64(bestGutterBin)*2.0

	var lower, upper, spanning []extractor.TextBlock
	for _, b := range blocks {
		if b.Y < gutterY && b.Y+b.Height > gutterY {
			spanning = append(spanning, b)
		} else if b.Y+b.Height <= gutterY {
			lower = append(lower, b)
		} else {
			upper = append(upper, b)
		}
	}

	// Recurse
	// Upper and Lower might be columns, so call detectAndSortColumns
	sortedUpper := detectAndSortColumns(upper)
	sortedLower := detectAndSortColumns(lower)
	sortedSpanning := sortY(spanning)

	// Order: Upper (Top), Spanning, Lower (Bottom)
	return append(append(sortedUpper, sortedSpanning...), sortedLower...)
}

func checkVerticalFlow(blocks []extractor.TextBlock) int {
	// Sort by Y
	sorted := make([]extractor.TextBlock, len(blocks))
	copy(sorted, blocks)
	sort.Slice(sorted, func(i, j int) bool {
		return (sorted[i].Y + sorted[i].Height) > (sorted[j].Y + sorted[j].Height)
	})

	score := 0
	for i := 0; i < len(sorted)-1; i++ {
		curr := sorted[i]
		next := sorted[i+1]

		// Check if next line is immediately below
		// Gap should be small (e.g. < 2 * Height)
		if (curr.Y - (next.Y + next.Height)) > curr.Height*2 {
			continue
		}

		text := curr.Text
		nextText := next.Text

		if len(text) > 0 && len(nextText) > 0 {
			lastChar := text[len(text)-1]
			firstChar := nextText[0]

			// Hyphenation
			if lastChar == '-' {
				score += 2
			}
			// Sentence continuation (lowercase start)
			if firstChar >= 'a' && firstChar <= 'z' {
				score += 1
			}
		}
	}
	return score
}

func checkHorizontalFlow(left, right []extractor.TextBlock) int {
	// Check flow from Left -> Right at similar Y
	score := 0
	for _, l := range left {
		lMid := l.Y + l.Height/2
		for _, r := range right {
			rMid := r.Y + r.Height/2
			if abs(lMid-rMid) < 5.0 { // Tolerance 5pt
				// Found aligned block
				text := l.Text
				nextText := r.Text

				if len(text) > 0 && len(nextText) > 0 {
					lastChar := text[len(text)-1]
					firstChar := nextText[0]

					// Hyphenation: Only count if next line starts with lowercase
					if lastChar == '-' {
						if firstChar >= 'a' && firstChar <= 'z' {
							score += 2
						}
					} else if firstChar >= 'a' && firstChar <= 'z' {
						// Sentence continuation (lowercase start) without hyphen
						score += 1
					}
					// Code syntax flow?
					// e.g. "mkdir" -> "-p"
					// "mongo" -> "javascript"
					// Hard to detect generically without language model.
					// But if it looks like a sentence, we count it.
				}
				break
			}
		}
	}
	return score
}

func checkRowAlignment(left, right []extractor.TextBlock) float64 {
	// Check how many blocks in Left have a corresponding block in Right at similar Y
	matches := 0
	for _, l := range left {
		lMid := l.Y + l.Height/2
		for _, r := range right {
			rMid := r.Y + r.Height/2
			if abs(lMid-rMid) < 5.0 { // Tolerance 5pt
				matches++
				break
			}
		}
	}
	if len(left) == 0 {
		return 0
	}
	return float64(matches) / float64(len(left))
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func sortY(blocks []extractor.TextBlock) []extractor.TextBlock {
	if len(blocks) == 0 {
		return nil
	}

	// 1. Sort by Y descending strictly (Top to Bottom in PDF)
	// We use a temporary slice to avoid modifying the input in place if it matters,
	// but here we return a new slice anyway.
	sorted := make([]extractor.TextBlock, len(blocks))
	copy(sorted, blocks)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Y > sorted[j].Y
	})

	// 2. Group into rows based on Y proximity
	var rows [][]extractor.TextBlock
	currentRow := []extractor.TextBlock{sorted[0]}

	for i := 1; i < len(sorted); i++ {
		// If Y diff is small, add to current row
		if math.Abs(sorted[i].Y-currentRow[0].Y) <= 4.0 { // Increased tolerance to 4.0
			currentRow = append(currentRow, sorted[i])
		} else {
			// New row
			rows = append(rows, currentRow)
			currentRow = []extractor.TextBlock{sorted[i]}
		}
	}
	rows = append(rows, currentRow)

	// 3. Sort each row by X and return combined
	var result []extractor.TextBlock
	for _, row := range rows {
		// Sort row by X (Left to Right)
		sort.Slice(row, func(i, j int) bool {
			return row[i].X < row[j].X
		})
		result = append(result, row...)
	}
	return result
}
