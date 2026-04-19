package types

import (
	"math"
	"testing"
)

func TestBoundingBox_Width(t *testing.T) {
	tests := []struct {
		name string
		box  BoundingBox
		want float64
	}{
		{"simple", BoundingBox{X1: 0, Y1: 0, X2: 100, Y2: 50}, 100},
		{"offset", BoundingBox{X1: 10, Y1: 20, X2: 60, Y2: 80}, 50},
		{"zero", BoundingBox{X1: 50, Y1: 50, X2: 50, Y2: 50}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.box.Width(); got != tt.want {
				t.Errorf("Width() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoundingBox_Height(t *testing.T) {
	tests := []struct {
		name string
		box  BoundingBox
		want float64
	}{
		{"simple", BoundingBox{X1: 0, Y1: 0, X2: 100, Y2: 50}, 50},
		{"offset", BoundingBox{X1: 10, Y1: 20, X2: 60, Y2: 80}, 60},
		{"zero", BoundingBox{X1: 50, Y1: 50, X2: 50, Y2: 50}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.box.Height(); got != tt.want {
				t.Errorf("Height() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoundingBox_Area(t *testing.T) {
	tests := []struct {
		name string
		box  BoundingBox
		want float64
	}{
		{"simple", BoundingBox{X1: 0, Y1: 0, X2: 100, Y2: 50}, 5000},
		{"unit", BoundingBox{X1: 0, Y1: 0, X2: 1, Y2: 1}, 1},
		{"zero", BoundingBox{X1: 50, Y1: 50, X2: 50, Y2: 50}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.box.Area(); got != tt.want {
				t.Errorf("Area() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoundingBox_Center(t *testing.T) {
	tests := []struct {
		name  string
		box   BoundingBox
		wantX float64
		wantY float64
	}{
		{"simple", BoundingBox{X1: 0, Y1: 0, X2: 100, Y2: 50}, 50, 25},
		{"offset", BoundingBox{X1: 10, Y1: 20, X2: 30, Y2: 40}, 20, 30},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := tt.box.Center()
			if gotX != tt.wantX || gotY != tt.wantY {
				t.Errorf("Center() = (%v, %v), want (%v, %v)", gotX, gotY, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestBoundingBox_Contains(t *testing.T) {
	box := BoundingBox{X1: 10, Y1: 10, X2: 100, Y2: 100}
	tests := []struct {
		name string
		x, y float64
		want bool
	}{
		{"inside", 50, 50, true},
		{"on_edge", 10, 10, true},
		{"outside_left", 5, 50, false},
		{"outside_right", 105, 50, false},
		{"outside_top", 50, 5, false},
		{"outside_bottom", 50, 105, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := box.Contains(tt.x, tt.y); got != tt.want {
				t.Errorf("Contains(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestBoundingBox_Overlaps(t *testing.T) {
	box := BoundingBox{X1: 10, Y1: 10, X2: 100, Y2: 100}
	tests := []struct {
		name  string
		other BoundingBox
		want  bool
	}{
		{"full_overlap", BoundingBox{X1: 20, Y1: 20, X2: 80, Y2: 80}, true},
		{"partial_overlap", BoundingBox{X1: 50, Y1: 50, X2: 150, Y2: 150}, true},
		{"edge_touch", BoundingBox{X1: 100, Y1: 10, X2: 150, Y2: 100}, false},
		{"no_overlap", BoundingBox{X1: 200, Y1: 200, X2: 300, Y2: 300}, false},
		{"identical", BoundingBox{X1: 10, Y1: 10, X2: 100, Y2: 100}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := box.Overlaps(tt.other); got != tt.want {
				t.Errorf("Overlaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoundingBox_IoU(t *testing.T) {
	tests := []struct {
		name string
		a, b BoundingBox
		want float64
	}{
		{
			name: "identical",
			a:    BoundingBox{X1: 0, Y1: 0, X2: 100, Y2: 100},
			b:    BoundingBox{X1: 0, Y1: 0, X2: 100, Y2: 100},
			want: 1.0,
		},
		{
			name: "no_overlap",
			a:    BoundingBox{X1: 0, Y1: 0, X2: 50, Y2: 50},
			b:    BoundingBox{X1: 100, Y1: 100, X2: 150, Y2: 150},
			want: 0.0,
		},
		{
			name: "half_overlap",
			a:    BoundingBox{X1: 0, Y1: 0, X2: 100, Y2: 100},
			b:    BoundingBox{X1: 50, Y1: 0, X2: 150, Y2: 100},
			// Intersection: 50x100 = 5000, Union: 10000 + 10000 - 5000 = 15000
			// IoU: 5000/15000 = 1/3
			want: 1.0 / 3.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.IoU(tt.b)
			if math.Abs(got-tt.want) > 0.0001 {
				t.Errorf("IoU() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocLayoutClass_String(t *testing.T) {
	tests := []struct {
		class DocLayoutClass
		want  string
	}{
		{ClassTitle, "title"},
		{ClassPlainText, "plain_text"},
		{ClassFigure, "figure"},
		{ClassTable, "table"},
		{ClassIsolateFormula, "isolate_formula"},
		{DocLayoutClass(999), "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.class.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocLayoutClass_ToElementType(t *testing.T) {
	tests := []struct {
		class DocLayoutClass
		want  ElementType
	}{
		{ClassTitle, ElementTypeHeader},
		{ClassPlainText, ElementTypeParagraph},
		{ClassFigure, ElementTypeFigure},
		{ClassTable, ElementTypeTable},
		{ClassAbandon, ElementTypeUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.class.String(), func(t *testing.T) {
			if got := tt.class.ToElementType(); got != tt.want {
				t.Errorf("ToElementType() = %v, want %v", got, tt.want)
			}
		})
	}
}
