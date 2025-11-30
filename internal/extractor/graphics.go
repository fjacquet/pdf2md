package extractor

import (
	"fmt"
	"strings"

	"github.com/fjacquet/pdf2md/internal/types"
)

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
