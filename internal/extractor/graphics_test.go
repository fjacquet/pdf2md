package extractor

import (
	"strings"
	"testing"

	"github.com/fjacquet/pdf2md/internal/types"
)

func TestToSVG(t *testing.T) {
	vg := types.VectorGraphic{
		ID: "test_graphic",
		Operations: []types.PathOperation{
			{Type: types.PathOpMoveTo, Points: []types.Point{{X: 10, Y: 10}}},
			{Type: types.PathOpLineTo, Points: []types.Point{{X: 20, Y: 20}}},
			{Type: types.PathOpClose},
		},
		StrokeColor: "#FF0000",
		FillColor:   "#00FF00",
		LineWidth:   2.0,
		IsStroked:   true,
		IsFilled:    true,
		X:           10,
		Y:           10,
		Width:       10,
		Height:      10,
	}

	svg := ToSVG(vg)

	// Basic checks
	if !strings.HasPrefix(svg, "<svg") {
		t.Error("SVG should start with <svg")
	}
	if !strings.HasSuffix(svg, "</svg>") {
		t.Error("SVG should end with </svg>")
	}

	// Check for path commands
	expectedCommands := []string{"M 10.00 10.00", "L 20.00 20.00", "Z"}
	for _, cmd := range expectedCommands {
		if !strings.Contains(svg, cmd) {
			t.Errorf("SVG missing command: %s", cmd)
		}
	}

	// Check for styles
	if !strings.Contains(svg, "stroke=\"#FF0000\"") {
		t.Error("SVG missing stroke color")
	}
	if !strings.Contains(svg, "fill=\"#00FF00\"") {
		t.Error("SVG missing fill color")
	}
	if !strings.Contains(svg, "stroke-width=\"2.00\"") {
		t.Error("SVG missing stroke width")
	}
}

func TestToSVGDefaults(t *testing.T) {
	vg := types.VectorGraphic{
		Operations: []types.PathOperation{
			{Type: types.PathOpMoveTo, Points: []types.Point{{X: 0, Y: 0}}},
		},
		IsStroked: true,
		IsFilled:  true,
		// No colors specified, should default to black
	}

	svg := ToSVG(vg)

	if !strings.Contains(svg, "stroke=\"black\"") {
		t.Error("SVG missing default stroke color")
	}
	if !strings.Contains(svg, "fill=\"black\"") {
		t.Error("SVG missing default fill color")
	}
}
