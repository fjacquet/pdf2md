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

	// ONNX detection metadata (set when ONNX detection is used)
	ONNXClassID    int     // Class ID from ONNX model (0 if not detected)
	ONNXConfidence float64 // Detection confidence (0.0-1.0)
}

// ElementType defines the type of a document element
type ElementType string

// Element type constants for document elements.
const (
	ElementTypeUnknown    ElementType = ""
	ElementTypeHeader     ElementType = "header"
	ElementTypeParagraph  ElementType = "paragraph"
	ElementTypeCodeBlock  ElementType = "code_block"
	ElementTypeList       ElementType = "list"
	ElementTypeTable      ElementType = "table"
	ElementTypeAdmonition ElementType = "admonition"
	ElementTypeImage      ElementType = "image"
	ElementTypeEquation   ElementType = "equation"
	ElementTypeFigure     ElementType = "figure"
	ElementTypeCaption    ElementType = "caption"
	ElementTypeFootnote   ElementType = "footnote"
)
