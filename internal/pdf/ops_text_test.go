package pdf

import (
	"math"
	"testing"

	"github.com/fjacquet/pdf2md/internal/types"
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

	// TJ - now merges text with spacing adjustments into single block
	in.Stack = []Object{Array{StringLiteral("A"), Integer(-1000), StringLiteral("B")}}
	if err := in.handleTextShow("TJ"); err != nil {
		t.Errorf("TJ failed: %v", err)
	}
	// Should produce 1 block "A B" (large negative adjustment inserts space)
	if len(in.TextBlocks) != 2 { // Hello, "A B"
		t.Errorf("Expected 2 text blocks, got %d", len(in.TextBlocks))
	}
	if in.TextBlocks[1].Text != "A B" {
		t.Errorf("Expected text 'A B', got '%s'", in.TextBlocks[1].Text)
	}

	// ' (Tick)
	in.Stack = []Object{StringLiteral("Next Line")}
	if err := in.handleTextShow("'"); err != nil {
		t.Errorf("' failed: %v", err)
	}
	if len(in.TextBlocks) != 3 {
		t.Errorf("Expected 3 text blocks, got %d", len(in.TextBlocks))
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

// TestInterpreter_TJWordGap pins the TJ offset read as a word space:
// Chrome writes spaces as -166..-265, kerning stays within a few dozen units.
func TestInterpreter_TJWordGap(t *testing.T) {
	tests := []struct {
		name string
		arr  Array
		want string
	}{
		{"chrome space", Array{StringLiteral("A"), Integer(-199), StringLiteral("B")}, "A B"},
		{"justified space", Array{StringLiteral("A"), Real(-166.4), StringLiteral("B")}, "A B"},
		{"kerning", Array{StringLiteral("A"), Integer(-48), StringLiteral("B")}, "AB"},
		{"tightening", Array{StringLiteral("A"), Integer(50), StringLiteral("B")}, "AB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := NewInterpreter(nil, nil)
			_ = in.handleTextObject("BT")
			in.State.Tfs = 12
			in.Stack = []Object{tt.arr}
			if err := in.handleTextShow("TJ"); err != nil {
				t.Fatalf("TJ failed: %v", err)
			}
			if got := in.TextBlocks[0].Text; got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// fixedWidthFont returns a simple font where every code has width w and the
// space (32) has width 250, FirstChar 0.
func fixedWidthFont(w float64) *Font {
	widths := make([]float64, 256)
	for i := range widths {
		widths[i] = w
	}
	widths[32] = 250
	return &Font{Subtype: "Type1", ToUnicode: map[int]string{}, Encoding: map[int]string{32: "space"}, Widths: widths, FirstChar: 0, LastChar: 255}
}

// runText interprets content with font F1: every byte has width 500 and the
// space (32) has width 250, FirstChar 0. It returns the blocks.
func runText(t *testing.T, content string) []types.TextBlock {
	t.Helper()
	fm := NewFontManager(nil)
	fm.Fonts["F1"] = fixedWidthFont(500)
	blocks, _, _, err := NewInterpreter(fm, nil).Process([]byte(content))
	if err != nil {
		t.Fatalf("Process(%q): %v", content, err)
	}
	return blocks
}

func near(t *testing.T, what string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", what, got, want)
	}
}

func TestAdvance_TcTz(t *testing.T) {
	b := runText(t, "BT /F1 10 Tf 2 Tc 50 Tz 100 200 Td (AB) Tj (C) Tj ET")
	near(t, "blocks[0].Width", b[0].Width, 7) // per glyph (0.5*10 + 2) * 0.5 = 3.5
	near(t, "blocks[1].X", b[1].X, 107)
}

func TestAdvance_TwOnlyOnSpace(t *testing.T) {
	b := runText(t, "BT /F1 10 Tf 3 Tw 0 0 Td (A A) Tj (B) Tj ET")
	near(t, "blocks[1].X", b[1].X, 15.5) // 5 + (2.5 + 3) + 5
}

func TestAdvance_TJStringsAdvance(t *testing.T) {
	b := runText(t, "BT /F1 10 Tf 0 0 Td [(A) -1000 (B)] TJ (C) Tj ET")
	near(t, "blocks[1].X", b[1].X, 20)
}

func TestAdvance_NoDoubleScaling(t *testing.T) {
	b := runText(t, "BT /F1 10 Tf 2 0 0 2 0 0 Tm (A) Tj (B) Tj ET")
	near(t, "blocks[1].X", b[1].X, 10)
	near(t, "blocks[0].FontSize", b[0].FontSize, 20)
}

func TestAdvance_CTM(t *testing.T) {
	b := runText(t, "2 0 0 2 0 0 cm BT /F1 10 Tf (A) Tj (B) Tj ET")
	near(t, "blocks[1].X", b[1].X, 10)
}

func TestTD_SetsLeading(t *testing.T) {
	b := runText(t, "BT /F1 10 Tf 0 100 Td 0 -14 TD (A) Tj T* (B) Tj ET")
	near(t, "blocks[0].Y", b[0].Y, 86)
	near(t, "blocks[1].Y", b[1].Y, 72)
}

func TestRise(t *testing.T) {
	b := runText(t, "BT /F1 10 Tf 5 Ts (A) Tj ET")
	near(t, "blocks[0].Y", b[0].Y, 5)
}

func TestRenderMode(t *testing.T) {
	b := runText(t, "BT /F1 10 Tf 3 Tr (A) Tj ET")
	if b[0].RenderMode != 3 || b[0].Text != "A" {
		t.Errorf("got RenderMode %d text %q, want 3 and \"A\"", b[0].RenderMode, b[0].Text)
	}
}

func TestTJOnlyNumbers(t *testing.T) {
	b := runText(t, "BT /F1 10 Tf [-1000] TJ (A) Tj ET")
	if len(b) != 1 {
		t.Fatalf("got %d blocks, want 1", len(b))
	}
	near(t, "blocks[0].X", b[0].X, 10)
}

func TestNegativeFontSize(t *testing.T) {
	b := runText(t, "BT /F1 -10 Tf (AB) Tj ET")
	near(t, "Width", b[0].Width, 10)
	near(t, "FontSize", b[0].FontSize, 10)
	if b[0].Text != "AB" {
		t.Errorf("Text = %q, want \"AB\"", b[0].Text)
	}
}

func TestSpace_FontSpaceWidth(t *testing.T) {
	fm := NewFontManager(nil)
	wide := fixedWidthFont(500)
	wide.Widths[32] = 600
	fm.Fonts["F1"] = wide
	blocks, _, _, err := NewInterpreter(fm, nil).Process([]byte("BT /F1 10 Tf [(A) -199 (B)] TJ ET"))
	if err != nil {
		t.Fatal(err)
	}
	if blocks[0].Text != "AB" {
		t.Errorf("Text = %q, want \"AB\": -199 is under half of a 600 space", blocks[0].Text)
	}
	near(t, "SpaceWidth", blocks[0].SpaceWidth, 6)
}

func TestSpace_NoDoubleSpace(t *testing.T) {
	b := runText(t, "BT /F1 10 Tf [(A ) -1000 (B)] TJ ET")
	if b[0].Text != "A B" {
		t.Errorf("Text = %q, want \"A B\"", b[0].Text)
	}
}

func TestSpace_Rotated(t *testing.T) {
	b := runText(t, "BT /F1 10 Tf 0 1 -1 0 100 100 Tm [(A) -199 (B)] TJ ET")
	if b[0].Text != "A B" || b[0].Width < 0 {
		t.Errorf("got %q width %v, want \"A B\" and width >= 0", b[0].Text, b[0].Width)
	}
}
