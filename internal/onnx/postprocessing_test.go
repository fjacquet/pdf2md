package onnx

import (
	"math"
	"testing"

	"github.com/fjacquet/pdf2md/internal/types"
)

func TestParseDetections(t *testing.T) {
	raw := [][]float32{
		{10, 20, 100, 200, 0.9, 0},   // High confidence
		{50, 50, 150, 150, 0.3, 1},   // Medium confidence
		{0, 0, 50, 50, 0.1, 2},       // Low confidence
	}

	tests := []struct {
		name          string
		confThreshold float64
		wantCount     int
	}{
		{"high_threshold", 0.8, 1},
		{"medium_threshold", 0.25, 2},
		{"low_threshold", 0.05, 3},
		{"zero_threshold", 0.0, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseDetections(raw, tt.confThreshold)
			if len(result) != tt.wantCount {
				t.Errorf("ParseDetections count = %d, want %d", len(result), tt.wantCount)
			}
		})
	}
}

func TestParseDetections_Values(t *testing.T) {
	raw := [][]float32{
		{10, 20, 100, 200, 0.9, 0},
	}

	result := ParseDetections(raw, 0.0)
	if len(result) != 1 {
		t.Fatalf("expected 1 detection, got %d", len(result))
	}

	box := result[0]
	if box.X1 != 10 || box.Y1 != 20 || box.X2 != 100 || box.Y2 != 200 {
		t.Errorf("coordinates = (%v, %v, %v, %v), want (10, 20, 100, 200)",
			box.X1, box.Y1, box.X2, box.Y2)
	}
	// Use approximate comparison due to float32→float64 conversion
	if math.Abs(box.Confidence-0.9) > 0.001 {
		t.Errorf("confidence = %v, want ~0.9", box.Confidence)
	}
	if box.ClassID != 0 {
		t.Errorf("classID = %v, want 0", box.ClassID)
	}
}

func TestApplyNMS(t *testing.T) {
	tests := []struct {
		name         string
		boxes        []types.BoundingBox
		iouThreshold float64
		wantCount    int
	}{
		{
			name: "no_overlap",
			boxes: []types.BoundingBox{
				{X1: 0, Y1: 0, X2: 50, Y2: 50, Confidence: 0.9, ClassID: 0},
				{X1: 100, Y1: 100, X2: 150, Y2: 150, Confidence: 0.8, ClassID: 0},
			},
			iouThreshold: 0.5,
			wantCount:    2,
		},
		{
			name: "full_overlap_same_class",
			boxes: []types.BoundingBox{
				{X1: 0, Y1: 0, X2: 100, Y2: 100, Confidence: 0.9, ClassID: 0},
				{X1: 0, Y1: 0, X2: 100, Y2: 100, Confidence: 0.8, ClassID: 0},
			},
			iouThreshold: 0.5,
			wantCount:    1,
		},
		{
			name: "full_overlap_different_class",
			boxes: []types.BoundingBox{
				{X1: 0, Y1: 0, X2: 100, Y2: 100, Confidence: 0.9, ClassID: 0},
				{X1: 0, Y1: 0, X2: 100, Y2: 100, Confidence: 0.8, ClassID: 1},
			},
			iouThreshold: 0.5,
			wantCount:    2, // Different classes should not suppress each other
		},
		{
			name: "partial_overlap_below_threshold",
			boxes: []types.BoundingBox{
				{X1: 0, Y1: 0, X2: 100, Y2: 100, Confidence: 0.9, ClassID: 0},
				{X1: 80, Y1: 0, X2: 180, Y2: 100, Confidence: 0.8, ClassID: 0},
				// IoU = 20*100 / (10000 + 10000 - 2000) = 2000/18000 = 0.11
			},
			iouThreshold: 0.5,
			wantCount:    2, // IoU < threshold, both kept
		},
		{
			name:         "empty_input",
			boxes:        []types.BoundingBox{},
			iouThreshold: 0.5,
			wantCount:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyNMS(tt.boxes, tt.iouThreshold)
			if len(result) != tt.wantCount {
				t.Errorf("ApplyNMS count = %d, want %d", len(result), tt.wantCount)
			}
		})
	}
}

func TestApplyNMS_KeepsHighestConfidence(t *testing.T) {
	boxes := []types.BoundingBox{
		{X1: 0, Y1: 0, X2: 100, Y2: 100, Confidence: 0.5, ClassID: 0},
		{X1: 0, Y1: 0, X2: 100, Y2: 100, Confidence: 0.9, ClassID: 0},
		{X1: 0, Y1: 0, X2: 100, Y2: 100, Confidence: 0.7, ClassID: 0},
	}

	result := ApplyNMS(boxes, 0.5)
	if len(result) != 1 {
		t.Fatalf("expected 1 box, got %d", len(result))
	}

	if result[0].Confidence != 0.9 {
		t.Errorf("expected highest confidence (0.9), got %v", result[0].Confidence)
	}
}

func TestScaleBoxToPage(t *testing.T) {
	box := types.BoundingBox{
		X1: 100, Y1: 100, X2: 200, Y2: 200,
		Confidence: 0.9, ClassID: 0, ClassName: "test",
	}

	// Assume scale=2.0, padding=(50, 50), page=(612, 792)
	result := ScaleBoxToPage(box, 2.0, 50, 50, 612, 792)

	// Expected: x1 = (100-50)/2 = 25, y1 = (100-50)/2 = 25
	//           x2 = (200-50)/2 = 75, y2 = (200-50)/2 = 75
	if result.X1 != 25 || result.Y1 != 25 || result.X2 != 75 || result.Y2 != 75 {
		t.Errorf("scaled box = (%v, %v, %v, %v), want (25, 25, 75, 75)",
			result.X1, result.Y1, result.X2, result.Y2)
	}

	// Metadata should be preserved
	if result.Confidence != 0.9 || result.ClassID != 0 || result.ClassName != "test" {
		t.Error("metadata not preserved")
	}
}

func TestScaleBoxToPage_Clamping(t *testing.T) {
	box := types.BoundingBox{
		X1: -100, Y1: -100, X2: 2000, Y2: 2000,
		Confidence: 0.9, ClassID: 0,
	}

	result := ScaleBoxToPage(box, 1.0, 0, 0, 612, 792)

	// Values should be clamped to page bounds
	if result.X1 != 0 || result.Y1 != 0 {
		t.Errorf("min coordinates = (%v, %v), want (0, 0)", result.X1, result.Y1)
	}
	if result.X2 != 612 || result.Y2 != 792 {
		t.Errorf("max coordinates = (%v, %v), want (612, 792)", result.X2, result.Y2)
	}
}

func TestFilterByConfidence(t *testing.T) {
	boxes := []types.BoundingBox{
		{Confidence: 0.9},
		{Confidence: 0.5},
		{Confidence: 0.3},
		{Confidence: 0.1},
	}

	result := FilterByConfidence(boxes, 0.4)
	if len(result) != 2 {
		t.Errorf("filtered count = %d, want 2", len(result))
	}
}

func TestFilterByClass(t *testing.T) {
	boxes := []types.BoundingBox{
		{ClassID: 0},
		{ClassID: 1},
		{ClassID: 0},
		{ClassID: 2},
	}

	result := FilterByClass(boxes, 0)
	if len(result) != 2 {
		t.Errorf("filtered count = %d, want 2", len(result))
	}
}

func TestSortByConfidence(t *testing.T) {
	boxes := []types.BoundingBox{
		{Confidence: 0.3},
		{Confidence: 0.9},
		{Confidence: 0.5},
	}

	result := SortByConfidence(boxes)

	// Check descending order
	if result[0].Confidence != 0.9 || result[1].Confidence != 0.5 || result[2].Confidence != 0.3 {
		t.Error("boxes not sorted by confidence (descending)")
	}

	// Check original is unchanged
	if boxes[0].Confidence != 0.3 {
		t.Error("original slice was modified")
	}
}

func TestSortByPosition(t *testing.T) {
	boxes := []types.BoundingBox{
		{X1: 100, Y1: 100},
		{X1: 50, Y1: 50},
		{X1: 50, Y1: 100},
	}

	result := SortByPosition(boxes)

	// Expected order: (50,50), (50,100), (100,100)
	if result[0].Y1 != 50 || result[1].Y1 != 100 || result[2].Y1 != 100 {
		t.Error("boxes not sorted by position (top to bottom)")
	}
}

func TestPostprocess(t *testing.T) {
	raw := [][]float32{
		{100, 100, 200, 200, 0.9, 0},
		{110, 110, 210, 210, 0.8, 0}, // Overlaps with first
		{500, 500, 600, 600, 0.7, 1}, // Different location and class
	}

	// scale=1.0, no padding, page 1000x1000
	result := Postprocess(raw, 0.5, 0.5, 1.0, 0, 0, 1000, 1000)

	// Should have 2 boxes after NMS (first two overlap, third is separate)
	if len(result) != 2 {
		t.Errorf("postprocess count = %d, want 2", len(result))
	}

	// First should be the highest confidence one
	found09 := false
	found07 := false
	for _, box := range result {
		if math.Abs(box.Confidence-0.9) < 0.01 {
			found09 = true
		}
		if math.Abs(box.Confidence-0.7) < 0.01 {
			found07 = true
		}
	}
	if !found09 || !found07 {
		t.Error("expected boxes with confidence 0.9 and 0.7")
	}
}
