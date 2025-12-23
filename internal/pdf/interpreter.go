package pdf

import (
	"bytes"

	"github.com/fjacquet/pdf2md/internal/types"
)

// Interpreter interprets content streams
type Interpreter struct {
	FontManager *FontManager
	Resources   Dictionary
	Stack       []Object

	// Graphics State Stack
	State      GraphicsState
	StateStack []GraphicsState

	// Text Object State (reset at BT)
	Tm  Matrix
	Tlm Matrix

	TextBlocks []types.TextBlock
	Images     []types.Image
	Graphics   []types.VectorGraphic
}

// NewInterpreter creates a new interpreter
func NewInterpreter(fm *FontManager, resources Dictionary) *Interpreter {
	return &Interpreter{
		FontManager: fm,
		Resources:   resources,
		State: GraphicsState{
			CTM:       IdentityMatrix(),
			Th:        100.0, // Default horizontal scaling is 100%
			LineWidth: 1.0,   // Default line width
		},
		Tm:  IdentityMatrix(),
		Tlm: IdentityMatrix(),
	}
}

// Process interprets the content stream
func (in *Interpreter) Process(content []byte) ([]types.TextBlock, []types.Image, []types.VectorGraphic, error) {
	in.TextBlocks = nil
	in.Images = nil
	in.Graphics = nil
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
			if err := in.executeOperatorWithTokenizer(tok.Value, t); err != nil {
				return nil, nil, nil, err
			}
		} else {
			// Operand, push to stack
			// We need to parse it as object
			t.UnreadToken(tok)
			obj, err := p.ParseObject()
			if err != nil {
				return nil, nil, nil, err
			}
			in.Stack = append(in.Stack, obj)
		}
	}

	return in.TextBlocks, in.Images, in.Graphics, nil
}

func (in *Interpreter) executeOperatorWithTokenizer(op string, tokenizer *Tokenizer) error {
	switch op {
	// Text Object
	case "BT", "ET":
		return in.handleTextObject(op)

	// Text State
	case "Tf", "Tc", "Tw", "Tz", "TL", "Ts":
		return in.handleTextState(op)

	// Text Positioning
	case "Td", "Tm", "T*":
		return in.handleTextPosition(op)

	// Text Showing
	case "Tj", "TJ", "'", "\"":
		return in.handleTextShow(op)

	// Graphics State
	case "q", "Q", "cm", "w":
		return in.handleGraphicsState(op)

	// Path Construction
	case "m", "l", "c", "re", "h":
		return in.handlePathConstruction(op)

	// Path Painting
	case "S", "s", "f", "F", "f*", "B", "B*", "b", "b*", "n":
		return in.handlePathPainting(op)

	// XObject
	case "Do":
		return in.handleXObject(op)

	// Inline Image
	case "BI":
		if tokenizer != nil {
			return in.handleInlineImage(tokenizer)
		}
		// If no tokenizer available, skip
		return nil

	default:
		// Unknown operator, clear stack to be safe
		in.Stack = nil
	}
	return nil
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
