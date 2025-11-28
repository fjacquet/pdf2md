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
