package layout

import (
	"math"
	"sort"

	"github.com/fjacquet/pdf2md/internal/extractor"
	"github.com/fjacquet/pdf2md/internal/types"
)

// Line represents a detected line (horizontal or vertical)
type Line struct {
	X1, Y1       float64
	X2, Y2       float64
	Length       float64
	IsHorizontal bool
	IsVertical   bool
}

// GridRegion represents a region bounded by grid lines
type GridRegion struct {
	X, Y            float64
	Width, Height   float64
	HorizontalLines []float64 // Y positions of horizontal lines
	VerticalLines   []float64 // X positions of vertical lines
}

// classifyLines extracts horizontal and vertical lines from vector graphics
func classifyLines(graphics []types.VectorGraphic) (horizontal, vertical []Line) {
	const angleThreshold = 5.0 // degrees
	const minLength = 20.0     // minimum line length to consider

	for _, g := range graphics {
		// Only consider stroked paths (lines, not fills)
		if !g.IsStroked {
			continue
		}

		// Extract lines from path operations
		var currentX, currentY float64
		var startX, startY float64

		for _, op := range g.Operations {
			switch op.Type {
			case types.PathOpMoveTo:
				if len(op.Points) > 0 {
					currentX = op.Points[0].X
					currentY = op.Points[0].Y
					startX = currentX
					startY = currentY
				}

			case types.PathOpLineTo:
				if len(op.Points) > 0 {
					endX := op.Points[0].X
					endY := op.Points[0].Y

					line := classifyLine(currentX, currentY, endX, endY, angleThreshold, minLength)
					if line != nil {
						if line.IsHorizontal {
							horizontal = append(horizontal, *line)
						} else if line.IsVertical {
							vertical = append(vertical, *line)
						}
					}

					currentX = endX
					currentY = endY
				}

			case types.PathOpClose:
				// Close path - draw line back to start
				line := classifyLine(currentX, currentY, startX, startY, angleThreshold, minLength)
				if line != nil {
					if line.IsHorizontal {
						horizontal = append(horizontal, *line)
					} else if line.IsVertical {
						vertical = append(vertical, *line)
					}
				}
				currentX = startX
				currentY = startY
			}
		}
	}

	return horizontal, vertical
}

// classifyLine determines if a line segment is horizontal or vertical
func classifyLine(x1, y1, x2, y2 float64, angleThreshold, minLength float64) *Line {
	dx := x2 - x1
	dy := y2 - y1
	length := math.Sqrt(dx*dx + dy*dy)

	if length < minLength {
		return nil
	}

	// Calculate angle in degrees
	angle := math.Atan2(math.Abs(dy), math.Abs(dx)) * 180 / math.Pi

	line := &Line{
		X1:     x1,
		Y1:     y1,
		X2:     x2,
		Y2:     y2,
		Length: length,
	}

	if angle < angleThreshold {
		line.IsHorizontal = true
	} else if angle > 90-angleThreshold {
		line.IsVertical = true
	}

	return line
}

// findGridRegions finds rectangular regions defined by intersecting lines
func findGridRegions(hLines, vLines []Line) []GridRegion {
	if len(hLines) < 2 || len(vLines) < 2 {
		return nil
	}

	// Cluster lines by position (merge nearby lines)
	hPositions := clusterLinePositions(hLines, true)  // Y positions
	vPositions := clusterLinePositions(vLines, false) // X positions

	if len(hPositions) < 2 || len(vPositions) < 2 {
		return nil
	}

	// Sort positions
	sort.Float64s(hPositions)
	sort.Float64s(vPositions)

	// Find the largest continuous grid
	// For simplicity, we take all positions as one grid region
	region := GridRegion{
		X:               vPositions[0],
		Y:               hPositions[0],
		Width:           vPositions[len(vPositions)-1] - vPositions[0],
		Height:          hPositions[len(hPositions)-1] - hPositions[0],
		HorizontalLines: hPositions,
		VerticalLines:   vPositions,
	}

	return []GridRegion{region}
}

// clusterLinePositions groups lines by their position and returns unique positions
func clusterLinePositions(lines []Line, useY bool) []float64 {
	const tolerance = 3.0 // Cluster lines within 3 points

	var positions []float64
	for _, line := range lines {
		var pos float64
		if useY {
			// For horizontal lines, use Y position (average of Y1 and Y2)
			pos = (line.Y1 + line.Y2) / 2
		} else {
			// For vertical lines, use X position (average of X1 and X2)
			pos = (line.X1 + line.X2) / 2
		}
		positions = append(positions, pos)
	}

	if len(positions) == 0 {
		return nil
	}

	sort.Float64s(positions)

	// Cluster nearby positions
	var clustered []float64
	clustered = append(clustered, positions[0])

	for i := 1; i < len(positions); i++ {
		if positions[i]-clustered[len(clustered)-1] > tolerance {
			clustered = append(clustered, positions[i])
		}
	}

	return clustered
}

// isBlockInRegion checks if a text block is within a grid region
func isBlockInRegion(block extractor.TextBlock, region GridRegion) bool {
	// Block center point
	cx := block.X + block.Width/2
	cy := block.Y + block.Height/2

	// Add tolerance
	const tolerance = 5.0

	return cx >= region.X-tolerance &&
		cx <= region.X+region.Width+tolerance &&
		cy >= region.Y-tolerance &&
		cy <= region.Y+region.Height+tolerance
}

// DetectTableFromGraphics attempts to detect a table structure from graphics and text blocks
func DetectTableFromGraphics(graphics []types.VectorGraphic, textBlocks []extractor.TextBlock) *TableStructure {
	detector := NewTableDetector()
	tables, _ := detector.detectFromGraphics(textBlocks, graphics)
	if len(tables) > 0 {
		return &tables[0]
	}
	return nil
}

// FindTableCellsFromLines uses horizontal and vertical lines to define cell boundaries
// and maps text blocks to those cells
func FindTableCellsFromLines(hLines, vLines []Line, blocks []extractor.TextBlock) [][]TableCell {
	if len(hLines) < 2 || len(vLines) < 2 {
		return nil
	}

	// Get sorted unique Y positions (row boundaries)
	hPositions := clusterLinePositions(hLines, true)
	sort.Float64s(hPositions)

	// Get sorted unique X positions (column boundaries)
	vPositions := clusterLinePositions(vLines, false)
	sort.Float64s(vPositions)

	numRows := len(hPositions) - 1
	numCols := len(vPositions) - 1

	if numRows < 1 || numCols < 1 {
		return nil
	}

	// Create cell grid
	cells := make([][]TableCell, numRows)
	for i := range cells {
		cells[i] = make([]TableCell, numCols)
		for j := range cells[i] {
			cells[i][j] = TableCell{
				Row:     i,
				Col:     j,
				RowSpan: 1,
				ColSpan: 1,
				X:       vPositions[j],
				Y:       hPositions[i],
				Width:   vPositions[j+1] - vPositions[j],
				Height:  hPositions[i+1] - hPositions[i],
			}
		}
	}

	// Map blocks to cells
	for _, block := range blocks {
		// Find which cell this block belongs to
		row := findRowForBlock(block, hPositions)
		col := findColForBlock(block, vPositions)

		if row >= 0 && row < numRows && col >= 0 && col < numCols {
			if cells[row][col].Content != "" {
				cells[row][col].Content += " " + block.Text
			} else {
				cells[row][col].Content = block.Text
			}
		}
	}

	return cells
}

// findRowForBlock finds which row a block belongs to based on Y position
func findRowForBlock(block extractor.TextBlock, hPositions []float64) int {
	cy := block.Y + block.Height/2

	for i := 0; i < len(hPositions)-1; i++ {
		if cy >= hPositions[i] && cy < hPositions[i+1] {
			return i
		}
	}

	// Check if above first row
	if cy < hPositions[0] {
		return 0
	}
	// Check if below last row
	if cy >= hPositions[len(hPositions)-1] {
		return len(hPositions) - 2
	}

	return -1
}

// findColForBlock finds which column a block belongs to based on X position
func findColForBlock(block extractor.TextBlock, vPositions []float64) int {
	cx := block.X + block.Width/2

	for i := 0; i < len(vPositions)-1; i++ {
		if cx >= vPositions[i] && cx < vPositions[i+1] {
			return i
		}
	}

	// Check if left of first column
	if cx < vPositions[0] {
		return 0
	}
	// Check if right of last column
	if cx >= vPositions[len(vPositions)-1] {
		return len(vPositions) - 2
	}

	return -1
}

// DetectHorizontalSeparators finds horizontal lines that could be table row separators
func DetectHorizontalSeparators(graphics []types.VectorGraphic, pageWidth float64) []float64 {
	hLines, _ := classifyLines(graphics)

	// Filter to lines that span a significant portion of the page
	minWidth := pageWidth * 0.3

	var separators []float64
	for _, line := range hLines {
		if line.Length >= minWidth {
			y := (line.Y1 + line.Y2) / 2
			separators = append(separators, y)
		}
	}

	sort.Float64s(separators)
	return separators
}

// DetectVerticalSeparators finds vertical lines that could be table column separators
func DetectVerticalSeparators(graphics []types.VectorGraphic, pageHeight float64) []float64 {
	_, vLines := classifyLines(graphics)

	// Filter to lines that span a significant portion vertically
	minHeight := pageHeight * 0.1

	var separators []float64
	for _, line := range vLines {
		if line.Length >= minHeight {
			x := (line.X1 + line.X2) / 2
			separators = append(separators, x)
		}
	}

	sort.Float64s(separators)
	return separators
}

// MergeRectanglesIntoCells detects rectangles (from graphics fill operations)
// and uses them as cell boundaries
func MergeRectanglesIntoCells(graphics []types.VectorGraphic) []GridRegion {
	var regions []GridRegion

	for _, g := range graphics {
		// Look for filled rectangles
		if !g.IsFilled {
			continue
		}

		// A rectangle typically has 4 or 5 operations: M, L, L, L, (Z)
		if len(g.Operations) < 4 {
			continue
		}

		// Check if it forms a rectangle
		rect := extractRectangle(g.Operations)
		if rect != nil {
			regions = append(regions, *rect)
		}
	}

	return regions
}

// extractRectangle extracts rectangle bounds from path operations
func extractRectangle(ops []types.PathOperation) *GridRegion {
	if len(ops) < 4 {
		return nil
	}

	var points []types.Point

	for _, op := range ops {
		switch op.Type {
		case types.PathOpMoveTo, types.PathOpLineTo:
			if len(op.Points) > 0 {
				points = append(points, op.Points[0])
			}
		}
	}

	if len(points) < 4 {
		return nil
	}

	// Check if points form a rectangle (axis-aligned)
	// A rectangle has 4 corners with 2 unique X and 2 unique Y values
	xVals := make(map[float64]bool)
	yVals := make(map[float64]bool)

	for _, p := range points[:4] {
		xVals[math.Round(p.X)] = true
		yVals[math.Round(p.Y)] = true
	}

	if len(xVals) != 2 || len(yVals) != 2 {
		return nil // Not a rectangle
	}

	// Find bounds
	minX, maxX := math.MaxFloat64, -math.MaxFloat64
	minY, maxY := math.MaxFloat64, -math.MaxFloat64

	for _, p := range points[:4] {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	return &GridRegion{
		X:      minX,
		Y:      minY,
		Width:  maxX - minX,
		Height: maxY - minY,
	}
}
