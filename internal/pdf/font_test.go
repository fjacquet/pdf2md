package pdf

import (
	"testing"
)

func TestFont_GetWidth(t *testing.T) {
	// Simple font with Widths array
	f1 := &Font{
		FirstChar: 32,
		LastChar:  34,
		Widths:    []float64{200, 300, 400}, // space, !, "
	}

	if w := f1.GetWidth(32); w != 200 {
		t.Errorf("Expected width 200 for char 32, got %f", w)
	}
	if w := f1.GetWidth(33); w != 300 {
		t.Errorf("Expected width 300 for char 33, got %f", w)
	}
	if w := f1.GetWidth(35); w != 0 {
		t.Errorf("Expected width 0 for char 35 (out of range), got %f", w)
	}

	// CID font with CIDWidths map
	f2 := &Font{
		CIDWidths:    map[int]float64{10: 500, 20: 600},
		DefaultWidth: 1000,
	}

	if w := f2.GetWidth(10); w != 500 {
		t.Errorf("Expected width 500 for CID 10, got %f", w)
	}
	if w := f2.GetWidth(99); w != 1000 {
		t.Errorf("Expected default width 1000 for CID 99, got %f", w)
	}
}

func TestFont_DecodeString(t *testing.T) {
	// Simple font with Encoding
	f1 := &Font{
		Encoding: map[int]string{
			65: "A",
			66: "B",
		},
	}
	// "AB"
	s1 := string([]byte{65, 66})
	if res := f1.DecodeString(s1); res != "AB" {
		t.Errorf("Expected 'AB', got '%s'", res)
	}

	// Font with ToUnicode CMap (Simple)
	f2 := &Font{
		ToUnicode: map[int]string{
			65: "X",
			66: "Y",
		},
	}
	if res := f2.DecodeString(s1); res != "XY" {
		t.Errorf("Expected 'XY', got '%s'", res)
	}

	// Font with ToUnicode CMap (Composite/Type0)
	f3 := &Font{
		Subtype: "Type0",
		ToUnicode: map[int]string{
			0x0041: "Z", // 00 41 -> Z
		},
	}
	// 00 41
	s3 := string([]byte{0, 65})
	if res := f3.DecodeString(s3); res != "Z" {
		t.Errorf("Expected 'Z', got '%s'", res)
	}
}

func TestFont_CalculateWidth(t *testing.T) {
	// Simple font
	f1 := &Font{
		FirstChar: 65, // A
		LastChar:  66, // B
		Widths:    []float64{500, 600},
	}
	// "AB"
	s1 := string([]byte{65, 66})
	if w := f1.CalculateWidth(s1); w != 1100 {
		t.Errorf("Expected width 1100, got %f", w)
	}

	// Composite font
	f2 := &Font{
		Subtype:      "Type0",
		DefaultWidth: 1000,
		CIDWidths: map[int]float64{
			0x0041: 500,
		},
	}
	// 00 41 00 42 (A, B)
	// A has width 500, B has default width 1000
	s2 := string([]byte{0, 65, 0, 66})
	if w := f2.CalculateWidth(s2); w != 1500 {
		t.Errorf("Expected width 1500, got %f", w)
	}
}

func TestFontManager_Methods(t *testing.T) {
	fm := NewFontManager(nil)
	f := &Font{
		FirstChar: 65,
		LastChar:  65,
		Widths:    []float64{500},
		Encoding:  map[int]string{65: "A"},
	}
	fm.Fonts[Name("F1")] = f

	if w := fm.CalculateWidth(Name("F1"), "A"); w != 500 {
		t.Errorf("Expected width 500, got %f", w)
	}
	if w := fm.CalculateWidth(Name("F2"), "A"); w != 0 {
		t.Errorf("Expected width 0 for unknown font, got %f", w)
	}

	if s := fm.DecodeString(Name("F1"), "A"); s != "A" {
		t.Errorf("Expected 'A', got '%s'", s)
	}
	if s := fm.DecodeString(Name("F2"), "A"); s != "A" {
		t.Errorf("Expected 'A' (raw) for unknown font, got '%s'", s)
	}
}

func TestType3Font_DecodeString(t *testing.T) {
	// Type3 font with Encoding mapping character codes to glyph names
	f := &Font{
		Subtype: "Type3",
		Encoding: map[int]string{
			65: "A",       // Standard glyph name
			66: "B",       // Standard glyph name
			67: "percent", // Glyph name that maps to %
		},
		FirstChar: 65,
		LastChar:  67,
		Widths:    []float64{500, 500, 500},
	}

	// Test standard character decoding
	s1 := string([]byte{65, 66}) // AB
	if res := f.DecodeString(s1); res != "AB" {
		t.Errorf("Type3: Expected 'AB', got '%s'", res)
	}

	// Test glyph name that requires GlyphToUnicode lookup
	s2 := string([]byte{67}) // percent -> %
	if res := f.DecodeString(s2); res != "%" {
		t.Errorf("Type3: Expected '%%', got '%s'", res)
	}

	// Type3 with ToUnicode (less common but possible)
	f2 := &Font{
		Subtype: "Type3",
		ToUnicode: map[int]string{
			65: "X",
			66: "Y",
		},
	}
	s3 := string([]byte{65, 66})
	if res := f2.DecodeString(s3); res != "XY" {
		t.Errorf("Type3 with ToUnicode: Expected 'XY', got '%s'", res)
	}
}

func TestType3Font_GetWidth(t *testing.T) {
	// Type3 font uses simple Widths array like Type1
	f := &Font{
		Subtype:   "Type3",
		FirstChar: 65,
		LastChar:  67,
		Widths:    []float64{500, 600, 700},
	}

	if w := f.GetWidth(65); w != 500 {
		t.Errorf("Type3: Expected width 500 for char 65, got %f", w)
	}
	if w := f.GetWidth(66); w != 600 {
		t.Errorf("Type3: Expected width 600 for char 66, got %f", w)
	}
	if w := f.GetWidth(68); w != 0 {
		t.Errorf("Type3: Expected width 0 for out-of-range char 68, got %f", w)
	}
}
