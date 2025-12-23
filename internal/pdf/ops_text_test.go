package pdf

import (
	"testing"
)

func TestInterpreter_HandleTextObject(t *testing.T) {
	in := NewInterpreter(nil, nil)

	// BT
	if err := in.handleTextObject("BT"); err != nil {
		t.Errorf("BT failed: %v", err)
	}
	if in.Tm[0] != 1 || in.Tlm[0] != 1 {
		t.Error("BT did not reset matrices")
	}

	// ET
	in.Stack = []Object{Integer(1)}
	if err := in.handleTextObject("ET"); err != nil {
		t.Errorf("ET failed: %v", err)
	}
	if len(in.Stack) != 0 {
		t.Error("ET did not clear stack")
	}
}

func TestInterpreter_HandleTextState(t *testing.T) {
	in := NewInterpreter(nil, nil)

	// Tf
	in.Stack = []Object{Name("F1"), Integer(12)}
	if err := in.handleTextState("Tf"); err != nil {
		t.Errorf("Tf failed: %v", err)
	}
	if in.State.Tf != "F1" || in.State.Tfs != 12 {
		t.Errorf("Tf did not set state: %v", in.State)
	}

	// Tc
	in.Stack = []Object{Real(1.5)}
	if err := in.handleTextState("Tc"); err != nil {
		t.Errorf("Tc failed: %v", err)
	}
	if in.State.Tc != 1.5 {
		t.Errorf("Tc did not set state: %v", in.State)
	}

	// Tw
	in.Stack = []Object{Real(2.5)}
	if err := in.handleTextState("Tw"); err != nil {
		t.Errorf("Tw failed: %v", err)
	}
	if in.State.Tw != 2.5 {
		t.Errorf("Tw did not set state: %v", in.State)
	}

	// Tz
	in.Stack = []Object{Real(100)}
	if err := in.handleTextState("Tz"); err != nil {
		t.Errorf("Tz failed: %v", err)
	}
	if in.State.Th != 100 {
		t.Errorf("Tz did not set state: %v", in.State)
	}

	// TL
	in.Stack = []Object{Real(14)}
	if err := in.handleTextState("TL"); err != nil {
		t.Errorf("TL failed: %v", err)
	}
	if in.State.Tl != 14 {
		t.Errorf("TL did not set state: %v", in.State)
	}

	// Ts
	in.Stack = []Object{Real(5)}
	if err := in.handleTextState("Ts"); err != nil {
		t.Errorf("Ts failed: %v", err)
	}
	if in.State.Ts != 5 {
		t.Errorf("Ts did not set state: %v", in.State)
	}
}

func TestInterpreter_HandleTextPosition(t *testing.T) {
	in := NewInterpreter(nil, nil)
	_ = in.handleTextObject("BT") // Reset matrices

	// Td
	in.Stack = []Object{Integer(10), Integer(20)}
	if err := in.handleTextPosition("Td"); err != nil {
		t.Errorf("Td failed: %v", err)
	}
	// Tm should be updated. Tlm starts as Identity.
	// Tlm = [1 0 0 1 10 20] * Identity = [1 0 0 1 10 20]
	if in.Tm[4] != 10 || in.Tm[5] != 20 {
		t.Errorf("Td did not update Tm correctly: %v", in.Tm)
	}

	// Tm
	in.Stack = []Object{Integer(1), Integer(0), Integer(0), Integer(1), Integer(50), Integer(60)}
	if err := in.handleTextPosition("Tm"); err != nil {
		t.Errorf("Tm failed: %v", err)
	}
	if in.Tm[4] != 50 || in.Tm[5] != 60 {
		t.Errorf("Tm did not set matrix correctly: %v", in.Tm)
	}

	// T*
	in.State.Tl = 14
	if err := in.handleTextPosition("T*"); err != nil {
		t.Errorf("T* failed: %v", err)
	}
	// T* moves down by Tl. Previous Tm was (50, 60).
	// T* -> 0 -14 Td relative to Tlm lines start.
	// Tlm was set to (50, 60) by Tm operator.
	// New Tlm = [1 0 0 1 0 -14] * [1 0 0 1 50 60]
	// = [1 0 0 1 50 46]
	if in.Tm[4] != 50 || in.Tm[5] != 46 {
		t.Errorf("T* did not update Tm correctly: %v", in.Tm)
	}
}

func TestInterpreter_HandleTextShow(t *testing.T) {
	in := NewInterpreter(nil, nil)
	_ = in.handleTextObject("BT")
	in.State.Tf = "F1"
	in.State.Tfs = 12

	// Tj
	in.Stack = []Object{StringLiteral("Hello")}
	if err := in.handleTextShow("Tj"); err != nil {
		t.Errorf("Tj failed: %v", err)
	}
	if len(in.TextBlocks) != 1 {
		t.Errorf("Expected 1 text block, got %d", len(in.TextBlocks))
	}
	if in.TextBlocks[0].Text != "Hello" {
		t.Errorf("Expected text 'Hello', got '%s'", in.TextBlocks[0].Text)
	}

	// TJ
	in.Stack = []Object{Array{StringLiteral("A"), Integer(-1000), StringLiteral("B")}}
	if err := in.handleTextShow("TJ"); err != nil {
		t.Errorf("TJ failed: %v", err)
	}
	// Should produce 2 blocks (A and B)
	if len(in.TextBlocks) != 3 { // Hello, A, B
		t.Errorf("Expected 3 text blocks, got %d", len(in.TextBlocks))
	}

	// ' (Tick)
	in.Stack = []Object{StringLiteral("Next Line")}
	if err := in.handleTextShow("'"); err != nil {
		t.Errorf("' failed: %v", err)
	}
	if len(in.TextBlocks) != 4 {
		t.Errorf("Expected 4 text blocks, got %d", len(in.TextBlocks))
	}

	// " (Quote)
	in.Stack = []Object{Real(2), Real(3), StringLiteral("Quote Line")}
	if err := in.handleTextShow("\""); err != nil {
		t.Errorf("\" failed: %v", err)
	}
	if in.State.Tw != 2 || in.State.Tc != 3 {
		t.Errorf("Quote did not set spacing: %v", in.State)
	}
}
