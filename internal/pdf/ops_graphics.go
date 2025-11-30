package pdf

import (
	"fmt"

	"github.com/fjacquet/pdf2md/internal/types"
)

// handleGraphicsState handles q, Q, cm, w
func (in *Interpreter) handleGraphicsState(op string) error {
	switch op {
	case "cm": // Concatenate matrix to CTM: a b c d e f cm
		if len(in.Stack) < 6 {
			return fmt.Errorf("stack underflow for cm")
		}
		f := toFloat(in.Stack[len(in.Stack)-1])
		e := toFloat(in.Stack[len(in.Stack)-2])
		d := toFloat(in.Stack[len(in.Stack)-3])
		c := toFloat(in.Stack[len(in.Stack)-4])
		b := toFloat(in.Stack[len(in.Stack)-5])
		a := toFloat(in.Stack[len(in.Stack)-6])
		in.Stack = in.Stack[:len(in.Stack)-6]

		m := Matrix{a, b, c, d, e, f}
		in.State.CTM = m.Multiply(in.State.CTM)

	case "q": // Save graphics state
		in.StateStack = append(in.StateStack, in.State)

	case "Q": // Restore graphics state
		if len(in.StateStack) > 0 {
			in.State = in.StateStack[len(in.StateStack)-1]
			in.StateStack = in.StateStack[:len(in.StateStack)-1]
		}

	case "w": // Set line width
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for w")
		}
		in.State.LineWidth = toFloat(in.Stack[len(in.Stack)-1])
		in.Stack = in.Stack[:len(in.Stack)-1]
	}
	return nil
}

// handlePathConstruction handles m, l, c, re, h
func (in *Interpreter) handlePathConstruction(op string) error {
	switch op {
	case "m": // Move to: x y m
		if len(in.Stack) < 2 {
			return fmt.Errorf("stack underflow for m")
		}
		y := toFloat(in.Stack[len(in.Stack)-1])
		x := toFloat(in.Stack[len(in.Stack)-2])
		in.Stack = in.Stack[:len(in.Stack)-2]

		// Transform point by CTM
		tx, ty := in.State.CTM.Transform(x, y)

		in.State.CurrentPath = append(in.State.CurrentPath, types.PathOperation{
			Type:   types.PathOpMoveTo,
			Points: []types.Point{{X: tx, Y: ty}},
		})

	case "l": // Line to: x y l
		if len(in.Stack) < 2 {
			return fmt.Errorf("stack underflow for l")
		}
		y := toFloat(in.Stack[len(in.Stack)-1])
		x := toFloat(in.Stack[len(in.Stack)-2])
		in.Stack = in.Stack[:len(in.Stack)-2]

		tx, ty := in.State.CTM.Transform(x, y)

		in.State.CurrentPath = append(in.State.CurrentPath, types.PathOperation{
			Type:   types.PathOpLineTo,
			Points: []types.Point{{X: tx, Y: ty}},
		})

	case "c": // Curve to: x1 y1 x2 y2 x3 y3 c
		if len(in.Stack) < 6 {
			return fmt.Errorf("stack underflow for c")
		}
		y3 := toFloat(in.Stack[len(in.Stack)-1])
		x3 := toFloat(in.Stack[len(in.Stack)-2])
		y2 := toFloat(in.Stack[len(in.Stack)-3])
		x2 := toFloat(in.Stack[len(in.Stack)-4])
		y1 := toFloat(in.Stack[len(in.Stack)-5])
		x1 := toFloat(in.Stack[len(in.Stack)-6])
		in.Stack = in.Stack[:len(in.Stack)-6]

		tx1, ty1 := in.State.CTM.Transform(x1, y1)
		tx2, ty2 := in.State.CTM.Transform(x2, y2)
		tx3, ty3 := in.State.CTM.Transform(x3, y3)

		in.State.CurrentPath = append(in.State.CurrentPath, types.PathOperation{
			Type: types.PathOpCurveTo,
			Points: []types.Point{
				{X: tx1, Y: ty1},
				{X: tx2, Y: ty2},
				{X: tx3, Y: ty3},
			},
		})

	case "re": // Rectangle: x y w h re
		if len(in.Stack) < 4 {
			return fmt.Errorf("stack underflow for re")
		}
		h := toFloat(in.Stack[len(in.Stack)-1])
		w := toFloat(in.Stack[len(in.Stack)-2])
		y := toFloat(in.Stack[len(in.Stack)-3])
		x := toFloat(in.Stack[len(in.Stack)-4])
		in.Stack = in.Stack[:len(in.Stack)-4]

		// Convert rect to move/line operations
		// Move to (x, y)
		tx, ty := in.State.CTM.Transform(x, y)
		in.State.CurrentPath = append(in.State.CurrentPath, types.PathOperation{
			Type:   types.PathOpMoveTo,
			Points: []types.Point{{X: tx, Y: ty}},
		})

		// Line to (x+w, y)
		tx, ty = in.State.CTM.Transform(x+w, y)
		in.State.CurrentPath = append(in.State.CurrentPath, types.PathOperation{
			Type:   types.PathOpLineTo,
			Points: []types.Point{{X: tx, Y: ty}},
		})

		// Line to (x+w, y+h)
		tx, ty = in.State.CTM.Transform(x+w, y+h)
		in.State.CurrentPath = append(in.State.CurrentPath, types.PathOperation{
			Type:   types.PathOpLineTo,
			Points: []types.Point{{X: tx, Y: ty}},
		})

		// Line to (x, y+h)
		tx, ty = in.State.CTM.Transform(x, y+h)
		in.State.CurrentPath = append(in.State.CurrentPath, types.PathOperation{
			Type:   types.PathOpLineTo,
			Points: []types.Point{{X: tx, Y: ty}},
		})

		// Close
		in.State.CurrentPath = append(in.State.CurrentPath, types.PathOperation{
			Type: types.PathOpClose,
		})

	case "h": // Close path
		in.State.CurrentPath = append(in.State.CurrentPath, types.PathOperation{
			Type: types.PathOpClose,
		})
	}
	return nil
}

// handlePathPainting handles S, s, f, F, B, b, n
func (in *Interpreter) handlePathPainting(op string) error {
	switch op {
	case "S", "s": // Stroke
		if op == "s" {
			// Close path first
			in.State.CurrentPath = append(in.State.CurrentPath, types.PathOperation{
				Type: types.PathOpClose,
			})
		}
		in.capturePath(true, false)

	case "f", "F", "f*": // Fill
		in.capturePath(false, true)

	case "B", "B*", "b", "b*": // Fill and Stroke
		if op == "b" || op == "b*" {
			// Close path first
			in.State.CurrentPath = append(in.State.CurrentPath, types.PathOperation{
				Type: types.PathOpClose,
			})
		}
		in.capturePath(true, true)

	case "n": // End path (no op)
		in.State.CurrentPath = nil
	}
	return nil
}

func (in *Interpreter) capturePath(stroke, fill bool) {
	if len(in.State.CurrentPath) == 0 {
		return
	}

	// Calculate bounding box
	minX, minY := 1e9, 1e9
	maxX, maxY := -1e9, -1e9

	for _, op := range in.State.CurrentPath {
		for _, p := range op.Points {
			if p.X < minX {
				minX = p.X
			}
			if p.Y < minY {
				minY = p.Y
			}
			if p.X > maxX {
				maxX = p.X
			}
			if p.Y > maxY {
				maxY = p.Y
			}
		}
	}

	// Filter out small paths and lines (noise reduction)
	// If width or height is very small, it's likely a line or dot
	if (maxX-minX) < 5.0 || (maxY-minY) < 5.0 {
		in.State.CurrentPath = nil
		return
	}

	// Create VectorGraphic
	vg := types.VectorGraphic{
		ID:         fmt.Sprintf("Path%d", len(in.Graphics)+1),
		Operations: make([]types.PathOperation, len(in.State.CurrentPath)),
		LineWidth:  in.State.LineWidth,
		IsStroked:  stroke,
		IsFilled:   fill,
		X:          minX,
		Y:          minY,
		Width:      maxX - minX,
		Height:     maxY - minY,
	}
	copy(vg.Operations, in.State.CurrentPath)

	in.Graphics = append(in.Graphics, vg)

	// Clear path
	in.State.CurrentPath = nil
}
