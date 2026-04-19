package ocr

import "testing"

// buildProbMap returns a w×h probability map with the given rectangle set to
// `fg` and everything else to zero.
func buildProbMap(w, h int, rect DetectionBox, fg float32) []float32 {
	m := make([]float32, w*h)
	for y := rect.MinY; y < rect.MaxY; y++ {
		for x := rect.MinX; x < rect.MaxX; x++ {
			if x < 0 || y < 0 || x >= w || y >= h {
				continue
			}
			m[y*w+x] = fg
		}
	}
	return m
}

func TestPostprocessDet_SingleRectangle(t *testing.T) {
	w, h := 50, 40
	rect := DetectionBox{MinX: 10, MinY: 5, MaxX: 30, MaxY: 20}
	probs := buildProbMap(w, h, rect, 0.9)

	boxes := PostprocessDet(probs, w, h, 1.0, 1.0, 0.3, 0.5, 3)
	if len(boxes) != 1 {
		t.Fatalf("want 1 box, got %d", len(boxes))
	}
	got := boxes[0]
	if got.MinX != rect.MinX || got.MinY != rect.MinY {
		t.Errorf("min corner = (%d,%d), want (%d,%d)", got.MinX, got.MinY, rect.MinX, rect.MinY)
	}
	if got.MaxX != rect.MaxX || got.MaxY != rect.MaxY {
		t.Errorf("max corner = (%d,%d), want (%d,%d)", got.MaxX, got.MaxY, rect.MaxX, rect.MaxY)
	}
	if got.Score < 0.5 {
		t.Errorf("score = %v, want ≥ 0.5", got.Score)
	}
}

func TestPostprocessDet_MultipleComponents(t *testing.T) {
	w, h := 40, 40
	probs := make([]float32, w*h)
	// Two separated rectangles.
	for _, r := range []DetectionBox{
		{MinX: 2, MinY: 2, MaxX: 10, MaxY: 10},
		{MinX: 20, MinY: 20, MaxX: 30, MaxY: 30},
	} {
		for y := r.MinY; y < r.MaxY; y++ {
			for x := r.MinX; x < r.MaxX; x++ {
				probs[y*w+x] = 0.95
			}
		}
	}

	boxes := PostprocessDet(probs, w, h, 1.0, 1.0, 0.3, 0.5, 3)
	if len(boxes) != 2 {
		t.Fatalf("want 2 boxes, got %d", len(boxes))
	}
}

func TestPostprocessDet_FilterSmallBox(t *testing.T) {
	w, h := 20, 20
	// Tiny 2×2 box — below minSize=3.
	rect := DetectionBox{MinX: 5, MinY: 5, MaxX: 7, MaxY: 7}
	probs := buildProbMap(w, h, rect, 0.9)

	boxes := PostprocessDet(probs, w, h, 1.0, 1.0, 0.3, 0.5, 3)
	if len(boxes) != 0 {
		t.Errorf("want 0 boxes (filtered), got %d", len(boxes))
	}
}

func TestPostprocessDet_FilterLowScore(t *testing.T) {
	w, h := 30, 30
	rect := DetectionBox{MinX: 5, MinY: 5, MaxX: 15, MaxY: 15}
	// Just above dbThreshold but below boxThreshold.
	probs := buildProbMap(w, h, rect, 0.35)

	boxes := PostprocessDet(probs, w, h, 1.0, 1.0, 0.3, 0.5, 3)
	if len(boxes) != 0 {
		t.Errorf("want 0 boxes (filtered), got %d", len(boxes))
	}
}

func TestPostprocessDet_ScaleBack(t *testing.T) {
	w, h := 40, 40
	rect := DetectionBox{MinX: 10, MinY: 10, MaxX: 20, MaxY: 20}
	probs := buildProbMap(w, h, rect, 0.9)

	// scale 2x — original image was half the size.
	boxes := PostprocessDet(probs, w, h, 2.0, 2.0, 0.3, 0.5, 3)
	if len(boxes) != 1 {
		t.Fatalf("want 1 box, got %d", len(boxes))
	}
	got := boxes[0]
	if got.MinX != 5 || got.MaxX != 10 || got.MinY != 5 || got.MaxY != 10 {
		t.Errorf("scaled box = %+v, want (5,5)-(10,10)", got)
	}
}

func TestAlignUp(t *testing.T) {
	cases := []struct {
		in, stride, want int
	}{
		{100, 32, 128},
		{128, 32, 128},
		{1, 32, 32},
		{0, 32, 0},
	}
	for _, c := range cases {
		if got := alignUp(c.in, c.stride); got != c.want {
			t.Errorf("alignUp(%d, %d) = %d, want %d", c.in, c.stride, got, c.want)
		}
	}
}
