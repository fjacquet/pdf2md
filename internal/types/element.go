// Package types defines shared data structures used across pdf2md packages.
package types

// Element represents a logical document element (header, paragraph, etc.)
type Element struct {
	Type     ElementType
	Content  string
	Level    int     // For headers: 1-6, for lists: indent level
	X        float64 // Bounding box
	Y        float64
	Width    float64
	Height   float64
	FontSize float64
	LinkURI  string
}

// ElementType defines the type of a document element
type ElementType string

// Element type constants for document elements.
const (
	ElementTypeHeader     ElementType = "header"
	ElementTypeParagraph  ElementType = "paragraph"
	ElementTypeCodeBlock  ElementType = "code_block"
	ElementTypeList       ElementType = "list"
	ElementTypeTable      ElementType = "table"
	ElementTypeAdmonition ElementType = "admonition"
	ElementTypeImage      ElementType = "image"
)
