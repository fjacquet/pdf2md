package types //nolint:revive // Package name "types" is a common pattern for shared types

// IsWordGap reports whether a gap between two glyphs is a word space: more
// than half the font's own space width. Kerning stays below it, an
// inter-word space (even a justified, shrunk one) stays above.
func IsWordGap(gap, spaceWidth float64) bool {
	return gap > 0.5*spaceWidth
}

// TextBlock represents a block of text with its position and font information
type TextBlock struct {
	Text     string
	X, Y     float64
	Width    float64
	Height   float64
	FontSize float64
	FontName string
	LinkURI  string
	// SpaceWidth is the font's space width in user space at this size, 0 when unknown (OCR).
	SpaceWidth float64
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
