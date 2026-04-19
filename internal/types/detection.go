// Package types defines shared data structures used across pdf2md packages.
package types

// BoundingBox represents a detected region in page coordinates.
type BoundingBox struct {
	X1, Y1     float64 // Top-left corner
	X2, Y2     float64 // Bottom-right corner
	Confidence float64 // Detection confidence [0,1]
	ClassID    int     // DocLayout class ID
	ClassName  string  // Human-readable class name
}

// Width returns the width of the bounding box.
func (b BoundingBox) Width() float64 {
	return b.X2 - b.X1
}

// Height returns the height of the bounding box.
func (b BoundingBox) Height() float64 {
	return b.Y2 - b.Y1
}

// Area returns the area of the bounding box.
func (b BoundingBox) Area() float64 {
	return b.Width() * b.Height()
}

// Center returns the center point of the bounding box.
func (b BoundingBox) Center() (x, y float64) {
	return (b.X1 + b.X2) / 2, (b.Y1 + b.Y2) / 2
}

// Contains returns true if point (x, y) is inside the bounding box.
func (b BoundingBox) Contains(x, y float64) bool {
	return x >= b.X1 && x <= b.X2 && y >= b.Y1 && y <= b.Y2
}

// Overlaps returns true if this box overlaps with another.
func (b BoundingBox) Overlaps(other BoundingBox) bool {
	return b.X1 < other.X2 && b.X2 > other.X1 && b.Y1 < other.Y2 && b.Y2 > other.Y1
}

// IoU calculates Intersection over Union with another bounding box.
// Returns a value between 0 (no overlap) and 1 (identical boxes).
func (b BoundingBox) IoU(other BoundingBox) float64 {
	// Calculate intersection
	x1 := max(b.X1, other.X1)
	y1 := max(b.Y1, other.Y1)
	x2 := min(b.X2, other.X2)
	y2 := min(b.Y2, other.Y2)

	if x2 <= x1 || y2 <= y1 {
		return 0 // No intersection
	}

	intersectionArea := (x2 - x1) * (y2 - y1)
	unionArea := b.Area() + other.Area() - intersectionArea

	if unionArea <= 0 {
		return 0
	}

	return intersectionArea / unionArea
}

// PageDetections holds all detections for a single page.
type PageDetections struct {
	PageIndex  int           // 0-indexed page number
	PageWidth  float64       // Page width in points
	PageHeight float64       // Page height in points
	Detections []BoundingBox // Detected regions
}

// DocLayoutClass represents document element classes from DocLayout-YOLO model.
type DocLayoutClass int

// DocLayout-YOLO class constants.
// These match the class indices in the wybxc/DocLayout-YOLO-DocStructBench-onnx model.
const (
	ClassTitle          DocLayoutClass = 0
	ClassPlainText      DocLayoutClass = 1
	ClassAbandon        DocLayoutClass = 2
	ClassFigure         DocLayoutClass = 3
	ClassFigureCaption  DocLayoutClass = 4
	ClassTable          DocLayoutClass = 5
	ClassTableCaption   DocLayoutClass = 6
	ClassTableFootnote  DocLayoutClass = 7
	ClassIsolateFormula DocLayoutClass = 8
	ClassFormulaCaption DocLayoutClass = 9
)

// String returns the human-readable name for a DocLayoutClass.
func (c DocLayoutClass) String() string {
	switch c {
	case ClassTitle:
		return "title"
	case ClassPlainText:
		return "plain_text"
	case ClassAbandon:
		return "abandon"
	case ClassFigure:
		return "figure"
	case ClassFigureCaption:
		return "figure_caption"
	case ClassTable:
		return "table"
	case ClassTableCaption:
		return "table_caption"
	case ClassTableFootnote:
		return "table_footnote"
	case ClassIsolateFormula:
		return "isolate_formula"
	case ClassFormulaCaption:
		return "formula_caption"
	default:
		return "unknown"
	}
}

// ToElementType converts a DocLayoutClass to the corresponding ElementType.
func (c DocLayoutClass) ToElementType() ElementType {
	switch c {
	case ClassTitle:
		return ElementTypeHeader
	case ClassPlainText:
		return ElementTypeParagraph
	case ClassFigure:
		return ElementTypeFigure
	case ClassFigureCaption:
		return ElementTypeCaption
	case ClassTable:
		return ElementTypeTable
	case ClassTableCaption:
		return ElementTypeCaption
	case ClassTableFootnote:
		return ElementTypeFootnote
	case ClassIsolateFormula, ClassFormulaCaption:
		return ElementTypeEquation
	case ClassAbandon:
		return ElementTypeUnknown // Ignored content
	default:
		return ElementTypeParagraph
	}
}
