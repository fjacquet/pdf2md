package pdf

import (
	"encoding/hex"
	"fmt"
	"math"

	"github.com/fjacquet/pdf2md/internal/types"
)

// handleTextObject handles BT and ET
func (in *Interpreter) handleTextObject(op string) error {
	switch op {
	case "BT": // Begin Text
		in.Tm = IdentityMatrix()
		in.Tlm = IdentityMatrix()
		in.Stack = nil
	case "ET": // End Text
		in.Stack = nil
	}
	return nil
}

// handleTextState handles Tf, Tc, Tw, Tz, TL, Ts
func (in *Interpreter) handleTextState(op string) error {
	switch op {
	case "Tf": // Set Text Font and Size: font size Tf
		if len(in.Stack) < 2 {
			return fmt.Errorf("stack underflow for Tf")
		}
		size := toFloat(in.Stack[len(in.Stack)-1])
		fontName, ok := in.Stack[len(in.Stack)-2].(Name)
		if !ok {
			return fmt.Errorf("invalid operand for Tf: expected Name")
		}
		in.Stack = in.Stack[:len(in.Stack)-2]

		in.State.Tf = fontName
		in.State.Tfs = size

	case "Tc": // Set Character Spacing
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for Tc")
		}
		in.State.Tc = toFloat(in.Stack[len(in.Stack)-1])
		in.Stack = in.Stack[:len(in.Stack)-1]

	case "Tw": // Set Word Spacing
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for Tw")
		}
		in.State.Tw = toFloat(in.Stack[len(in.Stack)-1])
		in.Stack = in.Stack[:len(in.Stack)-1]

	case "Tz": // Set Horizontal Scaling
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for Tz")
		}
		in.State.Th = toFloat(in.Stack[len(in.Stack)-1])
		in.Stack = in.Stack[:len(in.Stack)-1]

	case "TL": // Set text leading
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for TL")
		}
		val := toFloat(in.Stack[len(in.Stack)-1])
		in.Stack = in.Stack[:len(in.Stack)-1]
		in.State.Tl = val

	case "Ts": // Set text rise
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for Ts")
		}
		in.State.Ts = toFloat(in.Stack[len(in.Stack)-1])
		in.Stack = in.Stack[:len(in.Stack)-1]
	}
	return nil
}

// handleTextPosition handles Td, Tm, T*
func (in *Interpreter) handleTextPosition(op string) error {
	switch op {
	case "Td": // Move text position: tx ty Td
		if len(in.Stack) < 2 {
			return fmt.Errorf("stack underflow for Td")
		}
		ty := toFloat(in.Stack[len(in.Stack)-1])
		tx := toFloat(in.Stack[len(in.Stack)-2])
		in.Stack = in.Stack[:len(in.Stack)-2]

		// Update Tlm
		// Tlm = [1 0 0 1 tx ty] * Tlm
		m := Matrix{1, 0, 0, 1, tx, ty}
		in.Tlm = m.Multiply(in.Tlm)
		in.Tm = in.Tlm

	case "Tm": // Set text matrix: a b c d e f Tm
		if len(in.Stack) < 6 {
			return fmt.Errorf("stack underflow for Tm")
		}
		f := toFloat(in.Stack[len(in.Stack)-1])
		e := toFloat(in.Stack[len(in.Stack)-2])
		d := toFloat(in.Stack[len(in.Stack)-3])
		c := toFloat(in.Stack[len(in.Stack)-4])
		b := toFloat(in.Stack[len(in.Stack)-5])
		a := toFloat(in.Stack[len(in.Stack)-6])
		in.Stack = in.Stack[:len(in.Stack)-6]

		in.Tm = Matrix{a, b, c, d, e, f}
		in.Tlm = in.Tm

	case "T*": // Move to start of next line
		// T* is equivalent to: 0 -Tl Td
		tx := 0.0
		ty := -in.State.Tl

		m := Matrix{1, 0, 0, 1, tx, ty}
		in.Tlm = m.Multiply(in.Tlm)
		in.Tm = in.Tlm

	}
	return nil
}

// handleTextShow handles Tj, TJ, ', "
func (in *Interpreter) handleTextShow(op string) error {
	switch op {
	case "Tj": // Show text: string Tj
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for Tj")
		}
		strObj := in.Stack[len(in.Stack)-1]
		in.Stack = in.Stack[:len(in.Stack)-1]

		if err := in.showText(strObj); err != nil {
			return err
		}

	case "TJ": // Show text with positioning: array TJ
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for TJ")
		}
		arrObj := in.Stack[len(in.Stack)-1]
		in.Stack = in.Stack[:len(in.Stack)-1]

		if arr, ok := arrObj.(Array); ok {
			// Process TJ array more intelligently:
			// - Small adjustments (< 200 units) are kerning within words
			// - Large adjustments (>= 200 units) are word spaces
			// Collect text and only split on large adjustments
			var textBuffer string
			var rawBuffer string
			startX := in.Tm[4]
			startY := in.Tm[5]

			for _, item := range arr {
				if s, ok := item.(StringLiteral); ok {
					rawBuffer += string(s)
					textBuffer += in.decodeText(string(s))
				} else if s, ok := item.(HexString); ok {
					b, _ := hex.DecodeString(string(s))
					rawBuffer += string(b)
					textBuffer += in.decodeText(string(b))
				} else if n, ok := item.(Integer); ok {
					// Adjustment value in thousandths of em
					// Positive = move left (tightening), Negative = move right (spacing)
					// Large negative values (< -200) typically indicate word space
					if n < -200 {
						// This is a word space - insert space in text
						textBuffer += " "
					}
					// Apply position adjustment
					adj := -float64(n) / 1000.0 * in.State.Tfs
					in.Tm[4] += adj * in.Tm[0]
					in.Tm[5] += adj * in.Tm[1]
				} else if f, ok := item.(Real); ok {
					if f < -200 {
						textBuffer += " "
					}
					adj := -float64(f) / 1000.0 * in.State.Tfs
					in.Tm[4] += adj * in.Tm[0]
					in.Tm[5] += adj * in.Tm[1]
				}
			}

			// Create single text block for entire TJ array
			if textBuffer != "" {
				// Restore start position for block creation
				savedTmX, savedTmY := in.Tm[4], in.Tm[5]
				in.Tm[4], in.Tm[5] = startX, startY

				width := 0.0
				if in.FontManager != nil && rawBuffer != "" {
					width = in.FontManager.CalculateWidth(in.State.Tf, rawBuffer) * in.State.Tfs / 1000.0
				}
				in.addTextBlock(textBuffer, width)

				// Restore end position
				in.Tm[4], in.Tm[5] = savedTmX, savedTmY
			}
		}

	case "'": // Move to next line and show text: string '
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for '")
		}
		text := in.Stack[len(in.Stack)-1]
		in.Stack = in.Stack[:len(in.Stack)-1]

		// T*
		tx := 0.0
		ty := -in.State.Tl
		m := Matrix{1, 0, 0, 1, tx, ty}
		in.Tlm = m.Multiply(in.Tlm)
		in.Tm = in.Tlm

		if err := in.showText(text); err != nil {
			return err
		}

	case "\"": // Move to next line and show text: aw ac string "
		if len(in.Stack) < 3 {
			return fmt.Errorf("stack underflow for \"")
		}
		text := in.Stack[len(in.Stack)-1]
		ac := toFloat(in.Stack[len(in.Stack)-2])
		aw := toFloat(in.Stack[len(in.Stack)-3])
		in.Stack = in.Stack[:len(in.Stack)-3]

		in.State.Tw = aw
		in.State.Tc = ac

		// T*
		tx := 0.0
		ty := -in.State.Tl
		m := Matrix{1, 0, 0, 1, tx, ty}
		in.Tlm = m.Multiply(in.Tlm)
		in.Tm = in.Tlm

		if err := in.showText(text); err != nil {
			return err
		}
	}
	return nil
}

func (in *Interpreter) showText(obj Object) error {
	text := ""
	rawText := ""
	if s, ok := obj.(StringLiteral); ok {
		rawText = string(s)
	} else if s, ok := obj.(HexString); ok {
		b, _ := hex.DecodeString(string(s))
		rawText = string(b)
	} else {
		return fmt.Errorf("invalid operand for text show operator: expected StringLiteral or HexString")
	}

	text = in.decodeText(rawText)

	width := 0.0
	if in.FontManager != nil {
		width = in.FontManager.CalculateWidth(in.State.Tf, rawText) * in.State.Tfs / 1000.0
	}

	realWidth := in.addTextBlock(text, width)

	// Update Tm
	// tx = width, ty = 0
	// Tm is [a b c d e f]
	// e += tx * a + ty * c = width * a
	// f += tx * b + ty * d = width * b
	in.Tm[4] += realWidth * in.Tm[0]
	in.Tm[5] += realWidth * in.Tm[1]
	return nil
}

func (in *Interpreter) addTextBlock(text string, width float64) float64 {
	// Calculate position
	// Text space origin (0,0) -> User space
	// First transform by Tm (Text Matrix)
	tx, ty := in.Tm.Transform(0, 0)

	// Then transform by CTM (Current Transformation Matrix)
	x, y := in.State.CTM.Transform(tx, ty)

	estimatedWidth := width
	if estimatedWidth == 0 {
		// Estimate width since we don't have font metrics
		// Average char width is usually ~0.25-0.3 * FontSize for variable width fonts
		estimatedWidth = float64(len(text)) * in.State.Tfs * 0.25
	}

	// Transform width to user space
	// We need to transform the vector (width, 0) by Tm then CTM
	// Vector transform: [w 0 0] * Tm * CTM
	// Or simply: p1 = (0,0), p2 = (width, 0). Transform both. Dist = p2 - p1.

	// Point (width, 0) in text space
	tx2, _ := in.Tm.Transform(estimatedWidth, 0)
	x2, _ := in.State.CTM.Transform(tx2, 0)

	finalWidth := x2 - x
	// finalHeight := y2 - y // Not really height, but change in Y

	// Normalize X and Width if Width is negative (flipped axis)
	if finalWidth < 0 {
		x += finalWidth
		finalWidth = -finalWidth
	}

	// Height is usually font size in user space
	// Transform (0, Tfs)
	tx3, ty3 := in.Tm.Transform(0, in.State.Tfs)
	x3, y3 := in.State.CTM.Transform(tx3, ty3)
	dx, dy := x3-x, y3-y
	finalFontSize := math.Sqrt(dx*dx + dy*dy)

	// Get real font name
	fontName := string(in.State.Tf)
	if in.FontManager != nil {
		if font, ok := in.FontManager.Fonts[in.State.Tf]; ok {
			if font.BaseFont != "" {
				fontName = font.BaseFont
			}
		}
	}

	in.TextBlocks = append(in.TextBlocks, types.TextBlock{
		Text:     text,
		X:        x,
		Y:        y,
		Width:    finalWidth,
		Height:   finalFontSize,
		FontSize: finalFontSize,
		FontName: fontName,
	})

	return estimatedWidth
}

func (in *Interpreter) decodeText(s string) string {
	// Use FontManager to decode string based on current font
	if in.FontManager != nil {
		return in.FontManager.DecodeString(in.State.Tf, s)
	}
	return s
}
