package pdf

import (
	"fmt"

	"github.com/fjacquet/pdf2md/internal/types"
)

// handleXObject handles Do
func (in *Interpreter) handleXObject(op string) error {
	if op == "Do" { // Invoke named XObject: name Do
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for Do")
		}
		name, ok := in.Stack[len(in.Stack)-1].(Name)
		if !ok {
			return fmt.Errorf("invalid operand for Do: expected Name")
		}
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

	switch subtype {
	case "Image":
		return in.extractImage(name, stream)
	case "Form":
		return in.processFormXObject(name, stream)
	default:
		return nil
	}
}

// processFormXObject handles Form XObjects which contain nested content streams
func (in *Interpreter) processFormXObject(name Name, stream *Stream) error {
	// Get the Form's Resources dictionary (may inherit from parent)
	var formResources Dictionary
	if res, ok := stream.Dictionary[Name("Resources")].(Dictionary); ok {
		formResources = res
	} else if ref, ok := stream.Dictionary[Name("Resources")].(IndirectRef); ok {
		// Resolve indirect reference
		if in.FontManager != nil && in.FontManager.reader != nil {
			resolved, err := in.FontManager.reader.ReadObject(ref.ObjectNumber)
			if err == nil {
				if dict, ok := resolved.(Dictionary); ok {
					formResources = dict
				}
			}
		}
	}

	// If no Form Resources, use parent Resources
	if formResources == nil {
		formResources = in.Resources
	}

	// Get the Form's Matrix (transformation matrix)
	formMatrix := IdentityMatrix()
	if matrixArr, ok := stream.Dictionary[Name("Matrix")].(Array); ok && len(matrixArr) == 6 {
		formMatrix = Matrix{
			toFloat(matrixArr[0]), toFloat(matrixArr[1]),
			toFloat(matrixArr[2]), toFloat(matrixArr[3]),
			toFloat(matrixArr[4]), toFloat(matrixArr[5]),
		}
	}

	// Decode the Form's content stream
	data, err := DecodeStreamFromDict(stream.Data, stream.Dictionary)
	if err != nil {
		fmt.Printf("Warning: Failed to decode Form XObject %s: %v\n", name, err)
		return nil // Don't fail the page for a single XObject
	}

	// Save current graphics state
	in.StateStack = append(in.StateStack, in.State)

	// Apply the Form's Matrix to CTM
	// CTM = CTM * FormMatrix
	in.State.CTM = in.State.CTM.Multiply(formMatrix)

	// Create a sub-interpreter for the Form content
	// Note: We reuse the same FontManager as it contains all font definitions
	subInterpreter := NewInterpreter(in.FontManager, formResources)

	// Copy current CTM so Form content is positioned correctly
	subInterpreter.State.CTM = in.State.CTM

	// Process the Form content
	textBlocks, images, graphics, err := subInterpreter.Process(data)
	if err != nil {
		fmt.Printf("Warning: Error processing Form XObject %s: %v\n", name, err)
		// Restore state and continue
		if len(in.StateStack) > 0 {
			in.State = in.StateStack[len(in.StateStack)-1]
			in.StateStack = in.StateStack[:len(in.StateStack)-1]
		}
		return nil
	}

	// Merge results from sub-interpreter
	in.TextBlocks = append(in.TextBlocks, textBlocks...)
	in.Images = append(in.Images, images...)
	in.Graphics = append(in.Graphics, graphics...)

	// Restore graphics state
	if len(in.StateStack) > 0 {
		in.State = in.StateStack[len(in.StateStack)-1]
		in.StateStack = in.StateStack[:len(in.StateStack)-1]
	}

	return nil
}

func (in *Interpreter) extractImage(name Name, stream *Stream) error {
	// Extract image data - apply filters
	data, err := DecodeStreamFromDict(stream.Data, stream.Dictionary)
	if err != nil {
		return err
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
