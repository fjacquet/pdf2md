package types //nolint:revive // Package name "types" is a common pattern for shared types

// WordGapEm is the smallest gap, in em, read as a word space: half of a
// typical 0.25 em space glyph.
// shortcut: assumes a 0.25 em space for every font; use the font's own space width once font widths follow the spec.
const WordGapEm = 0.125

// TextBlock represents a block of text with its position and font information
type TextBlock struct {
	Text     string
	X, Y     float64
	Width    float64
	Height   float64
	FontSize float64
	FontName string
	LinkURI  string
	// RenderMode is the Tr text rendering mode (0–7); 3 is invisible text, as in OCR layers.
	RenderMode int
}

// Image represents an extracted image
type Image struct {
	ID     string
	Data   []byte
	Format string // "png", "jpeg", etc.
	X, Y   float64
	Width  float64
	Height float64
}

// Link represents a hyperlink in the document
type Link struct {
	URI  string
	Rect []float64 // [x1, y1, x2, y2]
}
