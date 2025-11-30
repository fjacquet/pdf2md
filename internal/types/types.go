package types

// TextBlock represents a block of text with its position and font information
type TextBlock struct {
	Text     string
	X, Y     float64
	Width    float64
	Height   float64
	FontSize float64
	FontName string
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
