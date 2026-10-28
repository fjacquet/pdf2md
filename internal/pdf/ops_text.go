package pdf

import (
	"fmt"
	"math"
	"strings"

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

// handleTextState handles Tf, Tc, Tw, Tz, TL, Ts, Tr
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

	case "Tr": // Set text rendering mode
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for Tr")
		}
		in.State.Tmode = int(toFloat(in.Stack[len(in.Stack)-1]))
		in.Stack = in.Stack[:len(in.Stack)-1]
	}
	return nil
}

// handleTextPosition handles Td, TD, Tm, T*
func (in *Interpreter) handleTextPosition(op string) error {
	switch op {
	case "Td", "TD": // Move text position: tx ty Td; TD also sets the leading to -ty
		if len(in.Stack) < 2 {
			return fmt.Errorf("stack underflow for %s", op)
		}
		ty := toFloat(in.Stack[len(in.Stack)-1])
		tx := toFloat(in.Stack[len(in.Stack)-2])
		in.Stack = in.Stack[:len(in.Stack)-2]
		if op == "TD" {
			in.State.Tl = -ty
		}

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
			// One block per TJ. Strings advance glyph by glyph, numbers move
			// the pen by -n/1000 em: a word space when types.IsWordGap holds
			// against the font's space width, kerning otherwise.
			x0, y0 := in.textOrigin()
			var text strings.Builder
			for _, item := range arr {
				switch v := item.(type) {
				case StringLiteral, HexString:
					raw, err := stringBytes(v)
					if err != nil {
						continue
					}
					t, _ := in.showString(raw)
					text.WriteString(t)
				case Integer, Real:
					n := toFloat(v)
					if text.Len() > 0 && !strings.HasSuffix(text.String(), " ") && types.IsWordGap(-n, in.font().SpaceWidth()) {
						text.WriteString(" ")
					}
					in.Tm = Matrix{1, 0, 0, 1, -n / 1000 * in.State.Tfs * in.State.Th / 100, 0}.Multiply(in.Tm)
				}
			}
			if text.Len() > 0 {
				x1, _ := in.textOrigin()
				in.addTextBlock(text.String(), x0, y0, x1-x0)
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
	raw, err := stringBytes(obj)
	if err != nil {
		return fmt.Errorf("invalid operand for text show operator: %w", err)
	}
	x0, y0 := in.textOrigin()
	text, width := in.showString(raw)
	in.addTextBlock(text, x0, y0, width)
	return nil
}

// showString decodes raw and advances Tm glyph by glyph (ISO 32000 §9.4.4):
// tx = (w0/1000·Tfs + Tc + Tw·[single-byte code 32]) · Th/100. It returns the
// text and the signed horizontal advance of the text origin in user space.
func (in *Interpreter) showString(raw []byte) (string, float64) {
	x0, _ := in.textOrigin()
	font := in.font()
	for _, c := range font.Codes(raw) {
		tx := font.Width(c)/1000*in.State.Tfs + in.State.Tc
		if c.Len == 1 && c.Value == 32 {
			tx += in.State.Tw
		}
		in.Tm = Matrix{1, 0, 0, 1, tx * in.State.Th / 100, 0}.Multiply(in.Tm)
	}
	x1, _ := in.textOrigin()
	return in.decodeText(string(raw)), x1 - x0
}

// font is the current font, nil when Tf names a font absent from resources.
func (in *Interpreter) font() *Font {
	if in.FontManager == nil {
		return nil
	}
	return in.FontManager.Fonts[in.State.Tf]
}

// textOrigin is the text origin, raised by Ts, in user space.
func (in *Interpreter) textOrigin() (float64, float64) {
	tx, ty := in.Tm.Transform(0, in.State.Ts)
	return in.State.CTM.Transform(tx, ty)
}

// addTextBlock appends a block starting at (x, y) in user space with signed
// horizontal advance dx (negative for mirrored text).
func (in *Interpreter) addTextBlock(text string, x, y, dx float64) {
	if dx < 0 {
		x += dx
		dx = -dx
	}

	// Font size is the length of the (0, Tfs) text-space vector in user space,
	// the space width that of the (space advance, 0) vector.
	m := in.Tm.Multiply(in.State.CTM)
	fontSize := math.Hypot(in.State.Tfs*m[2], in.State.Tfs*m[3])
	sw := in.font().SpaceWidth() / 1000 * in.State.Tfs * in.State.Th / 100
	spaceWidth := math.Hypot(sw*m[0], sw*m[1])

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
		Text:       text,
		X:          x,
		Y:          y,
		Width:      dx,
		Height:     fontSize,
		FontSize:   fontSize,
		FontName:   fontName,
		SpaceWidth: spaceWidth,
		RenderMode: in.State.Tmode,
	})
}

func (in *Interpreter) decodeText(s string) string {
	// Use FontManager to decode string based on current font
	if in.FontManager != nil {
		return in.FontManager.DecodeString(in.State.Tf, s)
	}
	return s
}
