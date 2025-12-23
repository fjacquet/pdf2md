package extractor

import (
	"fmt"
	"math"
	"strings"

	"github.com/fjacquet/pdf2md/internal/types"
)

// IsTrivialGraphic checks if a vector graphic is a simple decorative element
// Returns true for simple rectangles or single-path shapes that are likely decorative
func IsTrivialGraphic(vg types.VectorGraphic) bool {
	ops := vg.Operations
	if len(ops) == 0 {
		return true
	}

	// Check if it's a simple rectangle (M, L, L, L, Z or M, L, L, L, L, Z pattern)
	if isSimpleRectangle(ops) {
		return true
	}

	// Check if it's a very simple path with just lines (likely a divider or border)
	if isSimpleLine(ops) {
		return true
	}

	return false
}

// isSimpleRectangle checks if operations form a simple rectangle
func isSimpleRectangle(ops []types.PathOperation) bool {
	// Rectangle pattern: M L L L Z (4 corners) or M L L L L Z (with extra close)
	if len(ops) < 4 || len(ops) > 6 {
		return false
	}

	// First op must be MoveTo
	if ops[0].Type != types.PathOpMoveTo {
		return false
	}

	// Count line operations and close
	lineCount := 0
	hasClose := false
	for i := 1; i < len(ops); i++ {
		switch ops[i].Type {
		case types.PathOpLineTo:
			lineCount++
		case types.PathOpClose:
			hasClose = true
		case types.PathOpCurveTo:
			// Has curves - not a simple rectangle
			return false
		}
	}

	// Simple rectangle: 3-4 lines with a close
	if lineCount >= 3 && lineCount <= 4 && hasClose {
		// Verify it's actually rectangular (all angles are 90 degrees)
		return isRectangularPath(ops)
	}

	return false
}

// isRectangularPath checks if the path forms a rectangle (orthogonal lines)
func isRectangularPath(ops []types.PathOperation) bool {
	// Extract points
	var points []types.Point
	for _, op := range ops {
		if len(op.Points) > 0 {
			// For LineTo and MoveTo, use the point
			// For CurveTo, use the end point (last point)
			points = append(points, op.Points[len(op.Points)-1])
		}
	}

	if len(points) < 4 {
		return false
	}

	// Check if all edges are horizontal or vertical (within tolerance)
	const tolerance = 0.5
	for i := 0; i < len(points)-1; i++ {
		p1, p2 := points[i], points[i+1]
		dx := math.Abs(p2.X - p1.X)
		dy := math.Abs(p2.Y - p1.Y)

		// Edge must be either horizontal or vertical
		isHorizontal := dy < tolerance
		isVertical := dx < tolerance

		if !isHorizontal && !isVertical {
			return false // Diagonal line - not a simple rectangle
		}
	}

	return true
}

// isSimpleLine checks if it's just a simple line or two lines (divider, border)
func isSimpleLine(ops []types.PathOperation) bool {
	if len(ops) < 2 || len(ops) > 3 {
		return false
	}

	// M L pattern (simple line)
	if len(ops) == 2 {
		return ops[0].Type == types.PathOpMoveTo && ops[1].Type == types.PathOpLineTo
	}

	// M L L pattern (two-segment line) or M L Z (line with close)
	if len(ops) == 3 {
		if ops[0].Type != types.PathOpMoveTo {
			return false
		}
		return ops[1].Type == types.PathOpLineTo &&
			(ops[2].Type == types.PathOpLineTo || ops[2].Type == types.PathOpClose)
	}

	return false
}

// ToSVG converts a VectorGraphic to an SVG string
func ToSVG(vg types.VectorGraphic) string {
	var sb strings.Builder

	// Add some padding
	padding := 5.0
	x := vg.X - padding
	y := vg.Y - padding
	width := vg.Width + 2*padding
	height := vg.Height + 2*padding

	// SVG Header
	sb.WriteString(fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"%.2f %.2f %.2f %.2f\" width=\"%.2f\" height=\"%.2f\">\n",
		x, y, width, height, width, height))

	// Path
	sb.WriteString("  <path d=\"")
	for _, op := range vg.Operations {
		switch op.Type {
		case types.PathOpMoveTo:
			sb.WriteString(fmt.Sprintf("M %.2f %.2f ", op.Points[0].X, op.Points[0].Y))
		case types.PathOpLineTo:
			sb.WriteString(fmt.Sprintf("L %.2f %.2f ", op.Points[0].X, op.Points[0].Y))
		case types.PathOpCurveTo:
			sb.WriteString(fmt.Sprintf("C %.2f %.2f %.2f %.2f %.2f %.2f ",
				op.Points[0].X, op.Points[0].Y,
				op.Points[1].X, op.Points[1].Y,
				op.Points[2].X, op.Points[2].Y))
		case types.PathOpClose:
			sb.WriteString("Z ")
		}
	}
	sb.WriteString("\"")

	// Style
	if vg.IsStroked {
		strokeColor := vg.StrokeColor
		if strokeColor == "" {
			strokeColor = "black"
		}
		sb.WriteString(fmt.Sprintf(" stroke=\"%s\" stroke-width=\"%.2f\"", strokeColor, vg.LineWidth))
	} else {
		sb.WriteString(" stroke=\"none\"")
	}

	if vg.IsFilled {
		fillColor := vg.FillColor
		if fillColor == "" {
			fillColor = "black"
		}
		sb.WriteString(fmt.Sprintf(" fill=\"%s\"", fillColor))
	} else {
		sb.WriteString(" fill=\"none\"")
	}

	sb.WriteString(" />\n")
	sb.WriteString("</svg>")

	return sb.String()
}
