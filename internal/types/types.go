package types //nolint:revive // Package name "types" is a common pattern for shared types

// TextBlock represents a block of text with its position and font information
type TextBlock struct {
	Text     string
	X, Y     float64
	Width    float64
	Height   float64
	FontSize float64
	FontName string
	LinkURI  string
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
