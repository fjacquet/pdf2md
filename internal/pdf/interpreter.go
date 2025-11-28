package pdf

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/fjacquet/pdf2md/internal/types"
)

// Matrix represents a 3x3 affine transformation matrix [a b c d e f]
// a b 0
// c d 0
// e f 1
type Matrix [6]float64

// IdentityMatrix returns [1 0 0 1 0 0]
func IdentityMatrix() Matrix {
	return Matrix{1, 0, 0, 1, 0, 0}
}

// Multiply multiplies two matrices: m1 x m2
func (m1 Matrix) Multiply(m2 Matrix) Matrix {
	// [a1 b1 0]   [a2 b2 0]
	// [c1 d1 0] x [c2 d2 0]
	// [e1 f1 1]   [e2 f2 1]

	return Matrix{
		m1[0]*m2[0] + m1[1]*m2[2],         // a
		m1[0]*m2[1] + m1[1]*m2[3],         // b
		m1[2]*m2[0] + m1[3]*m2[2],         // c
		m1[2]*m2[1] + m1[3]*m2[3],         // d
		m1[4]*m2[0] + m1[5]*m2[2] + m2[4], // e (tx)
		m1[4]*m2[1] + m1[5]*m2[3] + m2[5], // f (ty)
	}
}

// TransformPoint transforms a point (x, y)
func (m Matrix) Transform(x, y float64) (float64, float64) {
	// [x y 1] * m
	tx := x*m[0] + y*m[2] + m[4]
	ty := x*m[1] + y*m[3] + m[5]
	return tx, ty
}

// Interpreter interprets content streams
type Interpreter struct {
	CTM Matrix // Current Transformation Matrix

	// Text State
	Tm  Matrix  // Text Matrix
	Tlm Matrix  // Text Line Matrix
	Tf  Name    // Text Font Name
	Tfs float64 // Text Font Size
	Tl  float64 // Text Leading

	// Resources
	FontManager *FontManager

	// Stacks
	Stack []Object

	// Output
	TextBlocks []types.TextBlock
}

// NewInterpreter creates a new interpreter
func NewInterpreter(fm *FontManager) *Interpreter {
	return &Interpreter{
		CTM:         IdentityMatrix(),
		Tm:          IdentityMatrix(),
		Tlm:         IdentityMatrix(),
		FontManager: fm,
		Tl:          0, // Default leading is 0
	}
}

// Process interprets the content stream
func (in *Interpreter) Process(content []byte) ([]types.TextBlock, error) {
	in.TextBlocks = nil
	in.Stack = nil

	// Create tokenizer for content
	t := NewTokenizer(bytes.NewReader(content))
	p := NewParser(t)

	for {
		tok, err := t.NextToken()
		if err != nil {
			break // EOF
		}

		if tok.Type == TokenKeyword {
			// Operator
			if err := in.executeOperator(tok.Value); err != nil {
				return nil, err
			}
		} else {
			// Operand, push to stack
			// We need to parse it as object
			t.UnreadToken(tok)
			obj, err := p.ParseObject()
			if err != nil {
				return nil, err
			}
			in.Stack = append(in.Stack, obj)
		}
	}

	return in.TextBlocks, nil
}

func (in *Interpreter) executeOperator(op string) error {
	switch op {
	case "BT": // Begin Text
		in.Tm = IdentityMatrix()
		in.Tlm = IdentityMatrix()
		in.Stack = nil

	case "ET": // End Text
		in.Stack = nil

	case "Tf": // Set Text Font and Size: font size Tf
		if len(in.Stack) < 2 {
			return fmt.Errorf("stack underflow for Tf")
		}
		sizeObj := in.Stack[len(in.Stack)-1]
		fontObj := in.Stack[len(in.Stack)-2]
		in.Stack = in.Stack[:len(in.Stack)-2]

		in.Tf = fontObj.(Name)
		if s, ok := sizeObj.(Integer); ok {
			in.Tfs = float64(s)
		} else if s, ok := sizeObj.(Real); ok {
			in.Tfs = float64(s)
		}

	case "Td": // Move text position: tx ty Td
		if len(in.Stack) < 2 {
			return fmt.Errorf("stack underflow for Td")
		}
		ty := toFloat(in.Stack[len(in.Stack)-1])
		tx := toFloat(in.Stack[len(in.Stack)-2])
		in.Stack = in.Stack[:len(in.Stack)-2]

		// Tlm = [1 0 0 1 tx ty] x Tlm
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

	case "Tj": // Show text: string Tj
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for Tj")
		}
		strObj := in.Stack[len(in.Stack)-1]
		in.Stack = in.Stack[:len(in.Stack)-1]

		text := ""
		rawText := ""
		if s, ok := strObj.(StringLiteral); ok {
			rawText = string(s)
		} else if s, ok := strObj.(HexString); ok {
			// Hex string is already decoded to bytes by Tokenizer?
			// No, Tokenizer returns HexString value (e.g. "001A").
			// We need to decode hex to bytes first.
			b, _ := hex.DecodeString(string(s))
			rawText = string(b)
		}

		// Decode using current font
		if in.FontManager != nil {
			if font, ok := in.FontManager.Fonts[in.Tf]; ok {
				text = font.DecodeString(rawText)
			} else {
				text = rawText // Fallback
			}
		} else {
			text = rawText
		}

		in.addTextBlock(text)

	case "TJ": // Show text with positioning: array TJ
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for TJ")
		}
		arrObj := in.Stack[len(in.Stack)-1]
		in.Stack = in.Stack[:len(in.Stack)-1]

		if arr, ok := arrObj.(Array); ok {
			for _, item := range arr {
				if s, ok := item.(StringLiteral); ok {
					in.addTextBlock(string(s))
					// TODO: Advance width
				} else if s, ok := item.(HexString); ok {
					in.addTextBlock(string(s))
				} else if n, ok := item.(Integer); ok {
					// Adjustment
					_ = n
					// TODO: Apply adjustment to Tm
				} else if n, ok := item.(Real); ok {
					_ = n
				}
			}
		}

	case "T*": // Move to start of next line
		// Td(0, -Tleading)
		// Tlm = [1 0 0 1 0 -Tl] x Tlm
		m := Matrix{1, 0, 0, 1, 0, -in.Tl}
		in.Tlm = m.Multiply(in.Tlm)
		in.Tm = in.Tlm

	case "'": // Move to next line and show text: string '
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for '")
		}
		strObj := in.Stack[len(in.Stack)-1]
		in.Stack = in.Stack[:len(in.Stack)-1]

		// T*
		m := Matrix{1, 0, 0, 1, 0, -in.Tl}
		in.Tlm = m.Multiply(in.Tlm)
		in.Tm = in.Tlm

		// Tj
		text := ""
		if s, ok := strObj.(StringLiteral); ok {
			text = string(s)
		} else if s, ok := strObj.(HexString); ok {
			b, _ := hex.DecodeString(string(s))
			text = string(b)
		}

		// Decode
		if in.FontManager != nil {
			if font, ok := in.FontManager.Fonts[in.Tf]; ok {
				text = font.DecodeString(text)
			}
		}
		in.addTextBlock(text)

	case "\"": // Set word spacing, set character spacing, move to next line and show text: aw cw string "
		if len(in.Stack) < 3 {
			return fmt.Errorf("stack underflow for \"")
		}
		strObj := in.Stack[len(in.Stack)-1]
		// cw := in.Stack[len(in.Stack)-2]
		// aw := in.Stack[len(in.Stack)-3]
		in.Stack = in.Stack[:len(in.Stack)-3]

		// Tw (aw) - TODO: Implement Tw state
		// Tc (cw) - TODO: Implement Tc state

		// T*
		m := Matrix{1, 0, 0, 1, 0, -in.Tl}
		in.Tlm = m.Multiply(in.Tlm)
		in.Tm = in.Tlm

		// Tj
		text := ""
		if s, ok := strObj.(StringLiteral); ok {
			text = string(s)
		} else if s, ok := strObj.(HexString); ok {
			b, _ := hex.DecodeString(string(s))
			text = string(b)
		}

		// Decode
		if in.FontManager != nil {
			if font, ok := in.FontManager.Fonts[in.Tf]; ok {
				text = font.DecodeString(text)
			}
		}
		in.addTextBlock(text)

	case "Tc": // Set character spacing
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for Tc")
		}
		in.Stack = in.Stack[:len(in.Stack)-1] // Ignore for now

	case "Tw": // Set word spacing
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for Tw")
		}
		in.Stack = in.Stack[:len(in.Stack)-1] // Ignore for now

	case "Tz": // Set horizontal scaling
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for Tz")
		}
		in.Stack = in.Stack[:len(in.Stack)-1] // Ignore for now

	case "TL": // Set text leading
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for TL")
		}
		val := toFloat(in.Stack[len(in.Stack)-1])
		in.Stack = in.Stack[:len(in.Stack)-1]
		in.Tl = val

	case "Ts": // Set text rise
		if len(in.Stack) < 1 {
			return fmt.Errorf("stack underflow for Ts")
		}
		in.Stack = in.Stack[:len(in.Stack)-1] // Ignore for now

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
		in.CTM = m.Multiply(in.CTM)

	case "q": // Save graphics state
		// TODO: Push CTM to stack

	case "Q": // Restore graphics state
		// TODO: Pop CTM from stack

	default:
		// Unknown operator, clear stack to be safe
		in.Stack = nil
	}
	return nil
}

func (in *Interpreter) addTextBlock(text string) {
	// Calculate position
	// Text space origin (0,0) -> User space
	x, y := in.Tm.Transform(0, 0)

	// Apply CTM (if we tracked it properly, for now assume Identity or simple scaling)
	// x, y = in.CTM.Transform(x, y)

	// Estimate width since we don't have font metrics
	// Average char width is usually ~0.5-0.6 * FontSize
	estimatedWidth := float64(len(text)) * in.Tfs * 0.6

	// Get real font name
	fontName := string(in.Tf)
	if in.FontManager != nil {
		if font, ok := in.FontManager.Fonts[in.Tf]; ok {
			if font.BaseFont != "" {
				fontName = font.BaseFont
			}
		}
	}

	in.TextBlocks = append(in.TextBlocks, types.TextBlock{
		Text:     text,
		X:        x,
		Y:        y,
		Width:    estimatedWidth,
		FontSize: in.Tfs,
		FontName: fontName,
	})

	// Update Tm?
	// Strictly speaking, we should update Tm by the width of the text.
	// But we don't have font metrics yet.
	// So subsequent Tj in same BT block might be wrong if we don't update.
	// But usually Tj is followed by positioning operators or is a single string.
}

func toFloat(obj Object) float64 {
	if i, ok := obj.(Integer); ok {
		return float64(i)
	}
	if r, ok := obj.(Real); ok {
		return float64(r)
	}
	return 0
}
