package pdf

import (
	"fmt"

	"github.com/fjacquet/pdf2md/internal/types"
)

// handleXObject handles Do
func (in *Interpreter) handleXObject(op string) error {
	switch op {
	case "Do": // Invoke named XObject: name Do
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for Do")
		}
		name, ok := in.Stack[len(in.Stack)-1].(Name)
		if !ok {
			return fmt.Errorf("invalid operand for Do: expected Name")
		}
		fmt.Printf("DEBUG: Do operator called for %s\n", name)
		in.Stack = in.Stack[:len(in.Stack)-1]
		if err := in.processXObject(name); err != nil {
			return err
		}
	}
	return nil
}

func (in *Interpreter) processXObject(name Name) error {
	// Look up XObject in Resources
	if in.Resources == nil {
		// Some PDFs might call Do without resources if it's a standard name? Unlikely.
		return fmt.Errorf("no resources dictionary")
	}

	key := Name("XObject")
	xObjectDict, ok := in.Resources[key].(Dictionary)
	if !ok {
		// If Do is called but no XObject dict, it's an error in the PDF or our parsing.
		// But we shouldn't crash or fail the page if possible.
		// Just log warning and return.
		fmt.Printf("Warning: Do operator called for %s but no XObject dictionary found\n", name)
		return nil
	}

	// Resolve indirect reference if needed
	obj := xObjectDict[name]
	if obj == nil {
		return fmt.Errorf("XObject %s not found in resources", name)
	}

	if ref, ok := obj.(IndirectRef); ok {
		// Need access to reader to resolve reference.
		if in.FontManager != nil && in.FontManager.reader != nil {
			resolved, err := in.FontManager.reader.ReadObject(ref.ObjectNumber)
			if err != nil {
				return err
			}
			obj = resolved
		} else {
			return fmt.Errorf("cannot resolve XObject reference: no reader available")
		}
	}

	var stream *Stream
	if s, ok := obj.(Stream); ok {
		stream = &s
	} else if s, ok := obj.(*Stream); ok {
		stream = s
	} else {
		return fmt.Errorf("XObject %s is not a stream (type %T)", name, obj)
	}

	subtype, _ := stream.Dictionary[Name("Subtype")].(Name)
	fmt.Printf("DEBUG: XObject %s Subtype=%s\n", name, subtype)

	if subtype == "Image" {
		return in.extractImage(name, stream)
	} else if subtype == "Form" {
		// TODO: Handle Form XObjects (nested content)
	}

	return nil
}

func (in *Interpreter) extractImage(name Name, stream *Stream) error {
	// Extract image data
	// Apply filters
	data := stream.Data

	// Check filter
	if filter, ok := stream.Dictionary[Name("Filter")].(Name); ok {
		decoded, err := DecodeStream(data, filter)
		if err != nil {
			return err
		}
		data = decoded
	} else if filters, ok := stream.Dictionary[Name("Filter")].(Array); ok {
		// Multiple filters
		for _, f := range filters {
			if filterName, ok := f.(Name); ok {
				decoded, err := DecodeStream(data, filterName)
				if err != nil {
					return err
				}
				data = decoded
			}
		}
	}

	// Get dimensions from CTM
	// In PDF, images are drawn in a 1x1 unit square at (0,0) and transformed by CTM.
	// So the position is CTM.Transform(0,0) and size is determined by CTM scaling factors.
	// Width = CTM[0], Height = CTM[3] (assuming no rotation/skew for simplicity)

	x, y := in.State.CTM.Transform(0, 0)
	width := in.State.CTM[0]
	height := in.State.CTM[3]

	// Determine format
	format := "png" // Default
	// If Filter was DCTDecode, it's JPEG
	if filter, ok := stream.Dictionary[Name("Filter")].(Name); ok && filter == "DCTDecode" {
		format = "jpeg"
	}

	in.Images = append(in.Images, types.Image{
		ID:     string(name),
		Data:   data,
		Format: format,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	})

	return nil
}
