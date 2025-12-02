package pdf

import "github.com/fjacquet/pdf2md/internal/types"

// GraphicsState holds the current graphics state parameters
type GraphicsState struct {
	CTM Matrix

	// Text State
	Tc    float64 // Character spacing
	Tw    float64 // Word spacing
	Th    float64 // Horizontal scaling
	Tl    float64 // Leading
	Tf    Name    // Font name
	Tfs   float64 // Font size
	Tmode int     // Text rendering mode
	Tr    float64 // Text rise
	Ts    float64 // Text knockout (unused mostly)

	// Path Construction
	CurrentPath []types.PathOperation
	LineWidth   float64
}
