package layout

import (
	"testing"

	"github.com/fjacquet/pdf2md/internal/types"
)

func TestClassifyLine(t *testing.T) {
	tests := []struct {
		name        string
		x1, y1      float64
		x2, y2      float64
		wantHoriz   bool
		wantVert    bool
		angleThresh float64
		minLength   float64
	}{
		{
			name: "horizontal line",
			x1:   0, y1: 100,
			x2: 100, y2: 100,
			wantHoriz:   true,
			wantVert:    false,
			angleThresh: 5.0,
			minLength:   20.0,
		},
		{
			name: "vertical line",
			x1:   50, y1: 0,
			x2: 50, y2: 100,
			wantHoriz:   false,
			wantVert:    true,
			angleThresh: 5.0,
			minLength:   20.0,
		},
		{
			name: "diagonal line",
			x1:   0, y1: 0,
			x2: 100, y2: 100,
			wantHoriz:   false,
			wantVert:    false,
			angleThresh: 5.0,
			minLength:   20.0,
		},
		{
			name: "too short",
			x1:   0, y1: 0,
			x2: 5, y2: 0,
			wantHoriz:   false, // nil returned
			wantVert:    false,
			angleThresh: 5.0,
			minLength:   20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := classifyLine(tt.x1, tt.y1, tt.x2, tt.y2, tt.angleThresh, tt.minLength)

			if line == nil {
				if tt.wantHoriz || tt.wantVert {
					t.Error("expected line, got nil")
				}
				return
			}

			if line.IsHorizontal != tt.wantHoriz {
				t.Errorf("IsHorizontal = %v, want %v", line.IsHorizontal, tt.wantHoriz)
			}
			if line.IsVertical != tt.wantVert {
				t.Errorf("IsVertical = %v, want %v", line.IsVertical, tt.wantVert)
			}
		})
	}
}

func TestClassifyLines(t *testing.T) {
	// Create a simple rectangle graphic
	graphics := []types.VectorGraphic{
		{
			IsStroked: true,
			Operations: []types.PathOperation{
				{Type: types.PathOpMoveTo, Points: []types.Point{{X: 0, Y: 0}}},
				{Type: types.PathOpLineTo, Points: []types.Point{{X: 100, Y: 0}}},
				{Type: types.PathOpLineTo, Points: []types.Point{{X: 100, Y: 50}}},
				{Type: types.PathOpLineTo, Points: []types.Point{{X: 0, Y: 50}}},
				{Type: types.PathOpClose},
			},
		},
	}

	hLines, vLines := classifyLines(graphics)

	// Should have 2 horizontal lines (top and bottom)
	if len(hLines) != 2 {
		t.Errorf("expected 2 horizontal lines, got %d", len(hLines))
	}

	// Should have 2 vertical lines (left and right)
	if len(vLines) != 2 {
		t.Errorf("expected 2 vertical lines, got %d", len(vLines))
	}
}

func TestClusterLinePositions(t *testing.T) {
	lines := []Line{
		{Y1: 100, Y2: 100, IsHorizontal: true},
		{Y1: 102, Y2: 102, IsHorizontal: true}, // Close to first, should be clustered
		{Y1: 200, Y2: 200, IsHorizontal: true},
	}

	positions := clusterLinePositions(lines, true)

	// Should have 2 unique positions (100~102 clustered, and 200)
	if len(positions) != 2 {
		t.Errorf("expected 2 clustered positions, got %d: %v", len(positions), positions)
	}
}

func TestFindGridRegions(t *testing.T) {
	hLines := []Line{
		{Y1: 0, Y2: 0, IsHorizontal: true},
		{Y1: 50, Y2: 50, IsHorizontal: true},
		{Y1: 100, Y2: 100, IsHorizontal: true},
	}
	vLines := []Line{
		{X1: 0, X2: 0, IsVertical: true},
		{X1: 50, X2: 50, IsVertical: true},
		{X1: 100, X2: 100, IsVertical: true},
	}

	regions := findGridRegions(hLines, vLines)

	if len(regions) != 1 {
		t.Errorf("expected 1 grid region, got %d", len(regions))
		return
	}

	region := regions[0]
	if len(region.HorizontalLines) < 2 {
		t.Errorf("expected at least 2 horizontal lines in region")
	}
	if len(region.VerticalLines) < 2 {
		t.Errorf("expected at least 2 vertical lines in region")
	}
}

func TestDetectHorizontalSeparators(t *testing.T) {
	graphics := []types.VectorGraphic{
		{
			IsStroked: true,
			Operations: []types.PathOperation{
				{Type: types.PathOpMoveTo, Points: []types.Point{{X: 0, Y: 100}}},
				{Type: types.PathOpLineTo, Points: []types.Point{{X: 500, Y: 100}}},
			},
		},
		{
			IsStroked: true,
			Operations: []types.PathOperation{
				{Type: types.PathOpMoveTo, Points: []types.Point{{X: 0, Y: 50}}},
				{Type: types.PathOpLineTo, Points: []types.Point{{X: 100, Y: 50}}}, // Too short
			},
		},
	}

	pageWidth := 600.0
	separators := DetectHorizontalSeparators(graphics, pageWidth)

	// Only the first line (500px) should qualify (>=30% of 600)
	if len(separators) != 1 {
		t.Errorf("expected 1 separator, got %d", len(separators))
	}
}

func TestExtractRectangle(t *testing.T) {
	ops := []types.PathOperation{
		{Type: types.PathOpMoveTo, Points: []types.Point{{X: 10, Y: 10}}},
		{Type: types.PathOpLineTo, Points: []types.Point{{X: 100, Y: 10}}},
		{Type: types.PathOpLineTo, Points: []types.Point{{X: 100, Y: 50}}},
		{Type: types.PathOpLineTo, Points: []types.Point{{X: 10, Y: 50}}},
		{Type: types.PathOpClose},
	}

	rect := extractRectangle(ops)

	if rect == nil {
		t.Fatal("expected rectangle, got nil")
	}

	if rect.X != 10 {
		t.Errorf("expected x=10, got %f", rect.X)
	}
	if rect.Y != 10 {
		t.Errorf("expected y=10, got %f", rect.Y)
	}
	if rect.Width != 90 {
		t.Errorf("expected width=90, got %f", rect.Width)
	}
	if rect.Height != 40 {
		t.Errorf("expected height=40, got %f", rect.Height)
	}
}
