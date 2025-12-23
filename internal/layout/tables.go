package layout

import (
	"math"
	"sort"
	"strings"

	"github.com/fjacquet/pdf2md/internal/extractor"
	"github.com/fjacquet/pdf2md/internal/types"
)

// TableCell represents a single cell in a table
type TableCell struct {
	Content  string
	Row      int
	Col      int
	RowSpan  int
	ColSpan  int
	X, Y     float64
	Width    float64
	Height   float64
	IsHeader bool
}

// TableStructure represents a detected table
type TableStructure struct {
	Cells      [][]TableCell // [row][col]
	NumRows    int
	NumCols    int
	HasHeader  bool
	X, Y       float64
	Width      float64
	Height     float64
	Confidence float64 // 0.0 to 1.0
}

// TableDetector detects tables from text blocks and optional graphics
type TableDetector struct {
	MinColumns       int     // Minimum columns to be considered a table (default: 2)
	MinRows          int     // Minimum rows to be considered a table (default: 2)
	ColumnTolerance  float64 // X-position tolerance for column alignment (default: 5.0)
	RowTolerance     float64 // Y-position tolerance for row alignment (default: 3.0)
	MinConfidence    float64 // Minimum confidence score (default: 0.5)
	GapThreshold     float64 // Minimum gap ratio between columns (default: 2.0)
}

// NewTableDetector creates a new table detector with default settings
func NewTableDetector() *TableDetector {
	return &TableDetector{
		MinColumns:      2,
		MinRows:         2,
		ColumnTolerance: 5.0,
		RowTolerance:    3.0,
		MinConfidence:   0.5,
		GapThreshold:    2.0,
	}
}

// DetectTables finds tables in a set of text blocks
// Returns detected tables and the blocks that were consumed by tables
func (td *TableDetector) DetectTables(blocks []extractor.TextBlock, graphics []types.VectorGraphic) ([]TableStructure, map[int]bool) {
	if len(blocks) < td.MinRows*td.MinColumns {
		return nil, nil
	}

	// Try graphics-based detection first (most accurate)
	tables, consumedBlocks := td.detectFromGraphics(blocks, graphics)
	if len(tables) > 0 {
		return tables, consumedBlocks
	}

	// Fall back to grid-based detection
	return td.detectFromGrid(blocks)
}

// detectFromGraphics uses vector graphics (lines) to detect table structure
func (td *TableDetector) detectFromGraphics(blocks []extractor.TextBlock, graphics []types.VectorGraphic) ([]TableStructure, map[int]bool) {
	// Find horizontal and vertical lines
	hLines, vLines := classifyLines(graphics)

	if len(hLines) < 2 || len(vLines) < 2 {
		return nil, nil
	}

	var tables []TableStructure
	consumedBlocks := make(map[int]bool)

	// Find grid intersections
	gridRegions := findGridRegions(hLines, vLines)

	for _, region := range gridRegions {
		// Find blocks within this region
		var regionBlocks []struct {
			index int
			block extractor.TextBlock
		}

		for i, block := range blocks {
			if consumedBlocks[i] {
				continue
			}
			if isBlockInRegion(block, region) {
				regionBlocks = append(regionBlocks, struct {
					index int
					block extractor.TextBlock
				}{i, block})
			}
		}

		if len(regionBlocks) < td.MinRows*td.MinColumns {
			continue
		}

		// Build table structure from grid
		table := td.buildTableFromGrid(regionBlocks, region)
		if table != nil && table.Confidence >= td.MinConfidence {
			tables = append(tables, *table)
			for _, rb := range regionBlocks {
				consumedBlocks[rb.index] = true
			}
		}
	}

	return tables, consumedBlocks
}

// detectFromGrid uses text block positions to detect table structure
func (td *TableDetector) detectFromGrid(blocks []extractor.TextBlock) ([]TableStructure, map[int]bool) {
	if len(blocks) == 0 {
		return nil, nil
	}

	// Group blocks by Y position (rows)
	rows := td.groupByY(blocks)
	if len(rows) < td.MinRows {
		return nil, nil
	}

	// Find consistent column positions across rows
	columnPositions := td.findColumnPositionsFromBlocks(blocks, rows)
	if len(columnPositions) < td.MinColumns {
		return nil, nil
	}

	// Validate that rows have consistent structure
	table, consumedIndices := td.buildTableFromPositions(blocks, rows, columnPositions)
	if table == nil {
		return nil, nil
	}

	consumedBlocks := make(map[int]bool)
	for _, idx := range consumedIndices {
		consumedBlocks[idx] = true
	}

	return []TableStructure{*table}, consumedBlocks
}

// groupByY groups blocks that are on the same row (similar Y position)
func (td *TableDetector) groupByY(blocks []extractor.TextBlock) [][]int {
	if len(blocks) == 0 {
		return nil
	}

	// Create index-Y pairs and sort by Y (descending for top-to-bottom)
	type blockY struct {
		index int
		y     float64
	}
	pairs := make([]blockY, len(blocks))
	for i, b := range blocks {
		pairs[i] = blockY{i, b.Y}
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].y > pairs[j].y // Descending (top first)
	})

	var rows [][]int
	var currentRow []int
	var currentY float64 = math.MaxFloat64

	for _, p := range pairs {
		if math.Abs(p.y-currentY) > td.RowTolerance {
			// New row
			if len(currentRow) > 0 {
				rows = append(rows, currentRow)
			}
			currentRow = []int{p.index}
			currentY = p.y
		} else {
			currentRow = append(currentRow, p.index)
		}
	}
	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	return rows
}

// findColumnPositionsFromBlocks finds consistent column X positions
func (td *TableDetector) findColumnPositionsFromBlocks(blocks []extractor.TextBlock, rows [][]int) []float64 {
	// Collect all X positions
	xPositions := make(map[float64]int) // position -> count

	for _, row := range rows {
		for _, idx := range row {
			if idx < len(blocks) {
				// Round to tolerance
				x := math.Round(blocks[idx].X/td.ColumnTolerance) * td.ColumnTolerance
				xPositions[x]++
			}
		}
	}

	// Keep positions that appear in multiple rows
	minOccurrences := len(rows) / 2
	if minOccurrences < 1 {
		minOccurrences = 1
	}

	var positions []float64
	for pos, count := range xPositions {
		if count >= minOccurrences {
			positions = append(positions, pos)
		}
	}

	sort.Float64s(positions)

	// Merge close positions
	if len(positions) < 2 {
		return positions
	}

	merged := []float64{positions[0]}
	for i := 1; i < len(positions); i++ {
		if positions[i]-merged[len(merged)-1] > td.ColumnTolerance*2 {
			merged = append(merged, positions[i])
		}
	}

	return merged
}

// buildTableFromPositions constructs a table from detected row/column structure
func (td *TableDetector) buildTableFromPositions(blocks []extractor.TextBlock, rows [][]int, _ []float64) (*TableStructure, []int) {
	columnPositions := td.findColumnPositionsFromBlocks(blocks, rows)
	if len(columnPositions) < td.MinColumns {
		return nil, nil
	}

	numCols := len(columnPositions)
	numRows := len(rows)

	// Build cell grid
	cells := make([][]TableCell, numRows)
	var consumedIndices []int
	var minX, minY, maxX, maxY float64 = math.MaxFloat64, math.MaxFloat64, 0, 0

	for rowIdx, row := range rows {
		cells[rowIdx] = make([]TableCell, numCols)

		for _, blockIdx := range row {
			if blockIdx >= len(blocks) {
				continue
			}
			block := blocks[blockIdx]

			// Find column for this block
			colIdx := td.findColumnIndex(block.X, columnPositions)
			if colIdx < 0 || colIdx >= numCols {
				continue
			}

			// Merge content if cell already has content
			if cells[rowIdx][colIdx].Content != "" {
				cells[rowIdx][colIdx].Content += " " + strings.TrimSpace(block.Text)
			} else {
				cells[rowIdx][colIdx] = TableCell{
					Content: strings.TrimSpace(block.Text),
					Row:     rowIdx,
					Col:     colIdx,
					RowSpan: 1,
					ColSpan: 1,
					X:       block.X,
					Y:       block.Y,
					Width:   block.Width,
					Height:  block.Height,
				}
			}

			consumedIndices = append(consumedIndices, blockIdx)

			// Update bounds
			if block.X < minX {
				minX = block.X
			}
			if block.Y < minY {
				minY = block.Y
			}
			if block.X+block.Width > maxX {
				maxX = block.X + block.Width
			}
			if block.Y+block.Height > maxY {
				maxY = block.Y + block.Height
			}
		}
	}

	// Calculate confidence based on cell coverage
	totalCells := numRows * numCols
	filledCells := 0
	for _, row := range cells {
		for _, cell := range row {
			if cell.Content != "" {
				filledCells++
			}
		}
	}

	confidence := float64(filledCells) / float64(totalCells)

	// Additional confidence factors
	// - Rows with consistent cell count
	consistentRows := 0
	for _, row := range cells {
		filled := 0
		for _, cell := range row {
			if cell.Content != "" {
				filled++
			}
		}
		if filled >= numCols/2 {
			consistentRows++
		}
	}
	rowConsistency := float64(consistentRows) / float64(numRows)

	finalConfidence := (confidence + rowConsistency) / 2

	if finalConfidence < td.MinConfidence {
		return nil, nil
	}

	// Detect header row (first row often has different formatting)
	hasHeader := td.detectHeaderRow(cells, blocks)

	return &TableStructure{
		Cells:      cells,
		NumRows:    numRows,
		NumCols:    numCols,
		HasHeader:  hasHeader,
		X:          minX,
		Y:          minY,
		Width:      maxX - minX,
		Height:     maxY - minY,
		Confidence: finalConfidence,
	}, consumedIndices
}

// findColumnIndex finds which column a block belongs to based on X position
func (td *TableDetector) findColumnIndex(x float64, columnPositions []float64) int {
	for i, pos := range columnPositions {
		if math.Abs(x-pos) <= td.ColumnTolerance*2 {
			return i
		}
	}

	// Find nearest column
	minDist := math.MaxFloat64
	nearestCol := -1
	for i, pos := range columnPositions {
		dist := math.Abs(x - pos)
		if dist < minDist {
			minDist = dist
			nearestCol = i
		}
	}

	// Only return if reasonably close
	if minDist <= td.ColumnTolerance*5 {
		return nearestCol
	}

	return -1
}

// detectHeaderRow checks if the first row looks like a header
func (td *TableDetector) detectHeaderRow(cells [][]TableCell, blocks []extractor.TextBlock) bool {
	if len(cells) < 2 {
		return false
	}

	firstRow := cells[0]
	otherRows := cells[1:]

	// Check if first row has different characteristics:
	// 1. Bold font (would be in FontName)
	// 2. All filled cells
	// 3. Text content (not numbers)

	firstRowFilled := 0
	firstRowAllText := true
	for _, cell := range firstRow {
		if cell.Content != "" {
			firstRowFilled++
			// Check if content is numeric
			if isNumericContent(cell.Content) {
				firstRowAllText = false
			}
		}
	}

	// First row should be mostly filled
	if firstRowFilled < len(firstRow)/2 {
		return false
	}

	// Check if other rows have more numeric content
	otherRowsNumeric := 0
	for _, row := range otherRows {
		for _, cell := range row {
			if isNumericContent(cell.Content) {
				otherRowsNumeric++
			}
		}
	}

	// If first row is all text and other rows have numbers, likely a header
	totalOtherCells := len(otherRows) * len(firstRow)
	if totalOtherCells > 0 && firstRowAllText && float64(otherRowsNumeric)/float64(totalOtherCells) > 0.3 {
		return true
	}

	return false
}

// isNumericContent checks if content is primarily numeric
func isNumericContent(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}

	numDigits := 0
	numOther := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			numDigits++
		} else if r != '.' && r != ',' && r != '-' && r != '+' && r != '%' && r != '$' && r != ' ' {
			numOther++
		}
	}

	if numDigits == 0 {
		return false
	}

	return float64(numDigits)/float64(numDigits+numOther) > 0.5
}

// buildTableFromGrid builds a table structure from a grid region defined by graphics
func (td *TableDetector) buildTableFromGrid(regionBlocks []struct {
	index int
	block extractor.TextBlock
}, region gridRegion) *TableStructure {
	// Extract just the blocks
	var blocks []extractor.TextBlock
	for _, rb := range regionBlocks {
		blocks = append(blocks, rb.block)
	}

	rows := td.groupByY(blocks)
	if len(rows) < td.MinRows {
		return nil
	}

	// Use grid lines to define columns
	columnPositions := region.verticalLines
	if len(columnPositions) < td.MinColumns+1 {
		return nil
	}

	numCols := len(columnPositions) - 1 // N lines = N-1 columns
	numRows := len(rows)

	cells := make([][]TableCell, numRows)
	for rowIdx, row := range rows {
		cells[rowIdx] = make([]TableCell, numCols)
		for _, blockIdx := range row {
			if blockIdx >= len(blocks) {
				continue
			}
			block := blocks[blockIdx]

			// Find column from grid lines
			colIdx := -1
			for i := 0; i < len(columnPositions)-1; i++ {
				left := columnPositions[i]
				right := columnPositions[i+1]
				if block.X >= left-td.ColumnTolerance && block.X < right+td.ColumnTolerance {
					colIdx = i
					break
				}
			}

			if colIdx < 0 || colIdx >= numCols {
				continue
			}

			if cells[rowIdx][colIdx].Content != "" {
				cells[rowIdx][colIdx].Content += " " + strings.TrimSpace(block.Text)
			} else {
				cells[rowIdx][colIdx] = TableCell{
					Content: strings.TrimSpace(block.Text),
					Row:     rowIdx,
					Col:     colIdx,
					RowSpan: 1,
					ColSpan: 1,
					X:       block.X,
					Y:       block.Y,
					Width:   block.Width,
					Height:  block.Height,
				}
			}
		}
	}

	return &TableStructure{
		Cells:      cells,
		NumRows:    numRows,
		NumCols:    numCols,
		HasHeader:  td.detectHeaderRow(cells, blocks),
		X:          region.x,
		Y:          region.y,
		Width:      region.width,
		Height:     region.height,
		Confidence: 0.9, // Graphics-based detection is high confidence
	}
}

// ToMarkdown converts a TableStructure to markdown format
func (ts *TableStructure) ToMarkdown() string {
	if ts.NumRows == 0 || ts.NumCols == 0 {
		return ""
	}

	var sb strings.Builder

	// Calculate column widths
	colWidths := make([]int, ts.NumCols)
	for _, row := range ts.Cells {
		for colIdx, cell := range row {
			contentLen := len(cell.Content)
			if contentLen > colWidths[colIdx] {
				colWidths[colIdx] = contentLen
			}
		}
	}

	// Minimum width of 3 for separator
	for i := range colWidths {
		if colWidths[i] < 3 {
			colWidths[i] = 3
		}
	}

	// Output rows
	for rowIdx, row := range ts.Cells {
		sb.WriteString("|")
		for colIdx, cell := range row {
			content := cell.Content
			// Pad content
			padded := content + strings.Repeat(" ", colWidths[colIdx]-len(content))
			sb.WriteString(" " + padded + " |")
		}
		sb.WriteString("\n")

		// Add separator after header row
		if rowIdx == 0 && ts.HasHeader {
			sb.WriteString("|")
			for colIdx := range row {
				sb.WriteString(" " + strings.Repeat("-", colWidths[colIdx]) + " |")
			}
			sb.WriteString("\n")
		}
	}

	// If no header was detected, add separator after first row anyway (markdown requirement)
	if !ts.HasHeader && ts.NumRows > 0 {
		// Insert separator after first row
		lines := strings.Split(sb.String(), "\n")
		if len(lines) >= 1 {
			var result strings.Builder
			result.WriteString(lines[0] + "\n")
			result.WriteString("|")
			for colIdx := 0; colIdx < ts.NumCols; colIdx++ {
				result.WriteString(" " + strings.Repeat("-", colWidths[colIdx]) + " |")
			}
			result.WriteString("\n")
			for i := 1; i < len(lines); i++ {
				if lines[i] != "" {
					result.WriteString(lines[i] + "\n")
				}
			}
			return strings.TrimSuffix(result.String(), "\n")
		}
	}

	return strings.TrimSuffix(sb.String(), "\n")
}

// IsLikelyEquation checks if blocks look more like an equation than a table
// This helps disambiguate tables from mathematical expressions
func IsLikelyEquation(blocks []extractor.TextBlock) bool {
	if len(blocks) == 0 {
		return false
	}

	mathScore := 0
	totalBlocks := len(blocks)

	mathSymbols := []string{"=", "+", "-", "×", "÷", "∫", "∑", "∏", "√", "∂", "∇", "±", "≈", "≠", "≤", "≥", "∞"}
	mathFonts := []string{"cmmi", "cmsy", "cmex", "symbol", "math", "cambria math", "stix"}

	for _, block := range blocks {
		text := block.Text
		fontLower := strings.ToLower(block.FontName)

		// Check for math fonts
		for _, mf := range mathFonts {
			if strings.Contains(fontLower, mf) {
				mathScore += 2
				break
			}
		}

		// Check for math symbols
		for _, sym := range mathSymbols {
			if strings.Contains(text, sym) {
				mathScore++
			}
		}

		// Check for equation numbering pattern: (1), (2.3), [1]
		if isEquationNumber(text) {
			mathScore += 3
		}

		// Check for inline math delimiters
		if strings.Count(text, "$") >= 2 {
			mathScore += 2
		}

		// Check for superscripts/subscripts patterns like x^2, a_n
		if strings.Contains(text, "^") || strings.Contains(text, "_") {
			mathScore++
		}
	}

	// Normalize score
	normalizedScore := float64(mathScore) / float64(totalBlocks)

	return normalizedScore > 1.5 // Threshold for equation detection
}

// isEquationNumber checks if text matches equation numbering patterns
func isEquationNumber(text string) bool {
	text = strings.TrimSpace(text)

	// Pattern: (n) or (n.m) or (n.m.k)
	if strings.HasPrefix(text, "(") && strings.HasSuffix(text, ")") {
		inner := text[1 : len(text)-1]
		return isNumericWithDots(inner)
	}

	// Pattern: [n]
	if strings.HasPrefix(text, "[") && strings.HasSuffix(text, "]") {
		inner := text[1 : len(text)-1]
		return isNumericWithDots(inner)
	}

	return false
}

// Note: isNumericWithDots is defined in analyzer.go and shared between files
