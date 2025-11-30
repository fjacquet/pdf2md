package types

// Point represents a 2D point
type Point struct {
	X, Y float64
}

// PathOperationType represents the type of path operation
type PathOperationType string

const (
	PathOpMoveTo  PathOperationType = "M"
	PathOpLineTo  PathOperationType = "L"
	PathOpCurveTo PathOperationType = "C"
	PathOpClose   PathOperationType = "Z"
)

// PathOperation represents a single operation in a path
type PathOperation struct {
	Type   PathOperationType
	Points []Point // Control points and end point
}

// VectorGraphic represents a captured vector graphic (path)
type VectorGraphic struct {
	ID          string
	Operations  []PathOperation
	StrokeColor string // Hex or name
	FillColor   string // Hex or name
	LineWidth   float64
	IsStroked   bool
	IsFilled    bool

	// Bounding Box
	X, Y, Width, Height float64
}
