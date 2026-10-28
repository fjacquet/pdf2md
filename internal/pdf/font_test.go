package pdf

import (
	"bytes"
	"log"
	"os"
	"slices"
	"strings"
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

func TestFontManager_Methods(t *testing.T) {
	fm := NewFontManager(nil)
	f := &Font{
		FirstChar: 65,
		LastChar:  65,
		Widths:    []float64{500},
		Encoding:  map[int]string{65: "A"},
	}
	fm.Fonts[Name("F1")] = f

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

func mustParseFont(t *testing.T, dict Dictionary) *Font {
	t.Helper()
	return NewFontManager(nil).parseFont(dict)
}

func TestWidth_MissingWidth(t *testing.T) {
	f := mustParseFont(t, Dictionary{
		Name("Subtype"): Name("TrueType"), Name("BaseFont"): Name("Custom"),
		Name("FirstChar"): Integer(32), Name("LastChar"): Integer(33),
		Name("Widths"):         Array{Integer(250), Integer(300)},
		Name("FontDescriptor"): Dictionary{Name("MissingWidth"): Integer(400)},
	})
	if w := f.Width(Code{65, 1}); w != 400 {
		t.Errorf("Width(65) = %v, want 400 (MissingWidth)", w)
	}
}

func TestWidth_Type3FontMatrix(t *testing.T) {
	f := mustParseFont(t, Dictionary{
		Name("Subtype"): Name("Type3"), Name("FontMatrix"): Array{Real(0.01), Integer(0), Integer(0), Real(0.01), Integer(0), Integer(0)},
		Name("FirstChar"): Integer(65), Name("LastChar"): Integer(65), Name("Widths"): Array{Integer(50)},
		Name("Encoding"): Dictionary{Name("Differences"): Array{Integer(65), Name("A")}},
	})
	if w := f.Width(Code{65, 1}); w != 500 {
		t.Errorf("Width(65) = %v, want 500 (50 × FontMatrix[0] × 1000)", w)
	}
}

func TestWidth_Std14(t *testing.T) {
	for _, base := range []string{"Helvetica", "ABCDEF+Helvetica"} {
		f := mustParseFont(t, Dictionary{Name("Subtype"): Name("Type1"), Name("BaseFont"): Name(base), Name("Encoding"): Name("WinAnsiEncoding")})
		if w := f.Width(Code{'A', 1}); w != 667 {
			t.Errorf("%s: Width('A') = %v, want 667", base, w)
		}
		if w := f.SpaceWidth(); w != 278 {
			t.Errorf("%s: SpaceWidth() = %v, want 278", base, w)
		}
	}
}

func TestSpaceWidth_ZeroWidthSpaceGlyph(t *testing.T) {
	f := &Font{Subtype: "Type0", ToUnicode: map[int]string{3: " "}, CIDWidths: map[int]float64{3: 0}, DefaultWidth: 1000}
	if w := f.SpaceWidth(); w != 250 {
		t.Errorf("SpaceWidth() = %v, want 250 fallback for a zero-width space glyph", w)
	}
}

func TestWidth_UnknownFont(t *testing.T) {
	b := runText(t, "BT /F9 10 Tf (AB) Tj ET")
	near(t, "Width", b[0].Width, 10)
	var unknown *Font
	if w := unknown.SpaceWidth(); w != 250 {
		t.Errorf("nil font SpaceWidth() = %v, want 250", w)
	}
}

func TestCodes_Codespace(t *testing.T) {
	cmap := []byte("2 begincodespacerange\n<00> <80>\n<8140> <FFFF>\nendcodespacerange")
	f := mustParseFont(t, Dictionary{
		Name("Subtype"): Name("Type0"), Name("BaseFont"): Name("Mixed"),
		Name("Encoding"): Stream{Dictionary: Dictionary{}, Data: cmap},
	})
	got := f.Codes([]byte{0x41, 0x81, 0x40, 0x42})
	want := []Code{{0x41, 1}, {0x8140, 2}, {0x42, 1}}
	if !slices.Equal(got, want) {
		t.Errorf("Codes() = %v, want %v", got, want)
	}
}

func TestCodes_IdentityH(t *testing.T) {
	f := mustParseFont(t, Dictionary{Name("Subtype"): Name("Type0"), Name("Encoding"): Name("Identity-H")})
	got := f.Codes([]byte{0x00, 0x41, 0x00, 0x42})
	want := []Code{{0x41, 2}, {0x42, 2}}
	if !slices.Equal(got, want) {
		t.Errorf("Codes() = %v, want %v", got, want)
	}
}

// A font without a space glyph (pdfTeX: FirstChar above 32, MissingWidth 500)
// must not take MissingWidth as its space width.
func TestSpaceWidth_NoSpaceGlyph(t *testing.T) {
	f := mustParseFont(t, Dictionary{
		Name("Subtype"): Name("Type1"), Name("BaseFont"): Name("ABCDEF+CMTI10"), Name("Encoding"): Name("WinAnsiEncoding"),
		Name("FirstChar"): Integer(34), Name("LastChar"): Integer(34), Name("Widths"): Array{Integer(511)},
		Name("FontDescriptor"): Dictionary{Name("MissingWidth"): Integer(500)},
	})
	if w := f.SpaceWidth(); w != 250 {
		t.Errorf("SpaceWidth() = %v, want 250 (no space glyph)", w)
	}
}

// WinAnsi names both 32 and 160 "space": the space width is code 32's, every time.
func TestSpaceWidth_WinAnsiPicksCode32(t *testing.T) {
	widths := make(Array, 256-32)
	for i := range widths {
		widths[i] = Integer(900)
	}
	widths[0] = Integer(0) // pdfTeX CMR: code 32 present with width 0
	f := mustParseFont(t, Dictionary{
		Name("Subtype"): Name("Type1"), Name("BaseFont"): Name("ABCDEF+CMR10"), Name("Encoding"): Name("WinAnsiEncoding"),
		Name("FirstChar"): Integer(32), Name("LastChar"): Integer(255), Name("Widths"): widths,
	})
	for range 20 {
		if w := f.SpaceWidth(); w != 250 {
			t.Fatalf("SpaceWidth() = %v, want 250 (code 32 width 0, not code 160's 900)", w)
		}
	}
}

func cmap(body string) []byte {
	return []byte("/CIDInit /ProcSet findresource begin\nbegincmap\n" + body + "\nendcmap\n")
}

func TestToUnicode_BfrangeIncrement(t *testing.T) {
	m, _ := parseToUnicodeCMap(cmap("1 beginbfrange\n<0001> <0003> <0041>\nendbfrange"))
	for code, want := range map[int]string{1: "A", 2: "B", 3: "C"} {
		if m[code] != want {
			t.Errorf("code %d = %q, want %q", code, m[code], want)
		}
	}
}

// The last UTF-16 unit of a multi-character destination is incremented (§9.10.3).
func TestToUnicode_BfrangeMultiChar(t *testing.T) {
	m, _ := parseToUnicodeCMap(cmap("1 beginbfrange\n<0005> <0006> <00660069>\nendbfrange"))
	if m[5] != "fi" || m[6] != "fj" {
		t.Errorf("codes 5, 6 = %q, %q; want \"fi\", \"fj\"", m[5], m[6])
	}
}

func TestToUnicode_BfrangeArray(t *testing.T) {
	m, _ := parseToUnicodeCMap(cmap("1 beginbfrange\n<0007> <0008> [<0058> <00590059>]\nendbfrange"))
	if m[7] != "X" || m[8] != "YY" {
		t.Errorf("codes 7, 8 = %q, %q; want \"X\", \"YY\"", m[7], m[8])
	}
}

func TestToUnicode_FourByteCode(t *testing.T) {
	m, cs := parseToUnicodeCMap(cmap("1 begincodespacerange\n<00000000> <FFFFFFFF>\nendcodespacerange\n1 beginbfchar\n<00000041> <0042>\nendbfchar"))
	if m[0x41] != "B" {
		t.Errorf("code 0x41 = %q, want \"B\"", m[0x41])
	}
	if len(cs) != 1 || len(cs[0].lo) != 4 {
		t.Errorf("codespace = %v, want one 4-byte range", cs)
	}
}

func TestToUnicode_Codespace(t *testing.T) {
	_, cs := parseToUnicodeCMap(cmap("2 begincodespacerange\n<00> <80>\n<8140> <FFFF>\nendcodespacerange"))
	if len(cs) != 2 || len(cs[0].lo) != 1 || len(cs[1].lo) != 2 {
		t.Errorf("codespace = %v, want a 1-byte and a 2-byte range", cs)
	}
}

// A composite font without an embedded encoding CMap splits codes by its
// ToUnicode codespace.
func TestCodes_ToUnicodeCodespace(t *testing.T) {
	tu := cmap("2 begincodespacerange\n<00> <80>\n<8140> <FFFF>\nendcodespacerange\n1 beginbfchar\n<8140> <3042>\nendbfchar")
	f := mustParseFont(t, Dictionary{
		Name("Subtype"): Name("Type0"), Name("Encoding"): Name("Identity-H"),
		Name("ToUnicode"): Stream{Dictionary: Dictionary{}, Data: tu},
	})
	if got := f.DecodeString("A\x81\x40"); got != "Aあ" {
		t.Errorf("DecodeString = %q, want \"Aあ\"", got)
	}
}

const type1Cleartext = "%!PS-AdobeFont-1.0: Test\n/Encoding 256 array\n0 1 255 {1 index exch /.notdef put} for\ndup 65 /alpha put\nreadonly def\ncurrentfile eexec\n\x8f\x12binary"

func TestType1BuiltinEncoding(t *testing.T) {
	if enc := type1BuiltinEncoding([]byte(type1Cleartext)); enc[65] != "alpha" {
		t.Errorf("enc[65] = %q, want \"alpha\"", enc[65])
	}
	f := mustParseFont(t, Dictionary{
		Name("Subtype"): Name("Type1"), Name("BaseFont"): Name("ABCDEF+CMMI10"),
		Name("FontDescriptor"): Dictionary{Name("Flags"): Integer(4), Name("FontFile"): Stream{Dictionary: Dictionary{}, Data: []byte(type1Cleartext)}},
	})
	if got := f.DecodeString("A"); got != "α" {
		t.Errorf("DecodeString(\"A\") = %q, want \"α\" from the built-in encoding", got)
	}
}

func TestLoadFonts_BadFontKept(t *testing.T) {
	var logged bytes.Buffer
	log.SetOutput(&logged)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	fm := NewFontManager(nil)
	bad := Dictionary{
		Name("Subtype"): Name("Type1"), Name("BaseFont"): Name("Broken"), Name("Encoding"): Name("WinAnsiEncoding"),
		Name("FirstChar"): Integer(65), Name("LastChar"): Integer(65), Name("Widths"): Array{Integer(600)},
		Name("ToUnicode"): Stream{Dictionary: Dictionary{Name("Filter"): Name("FlateDecode")}, Data: []byte("not zlib")},
	}
	if err := fm.LoadFonts(Dictionary{Name("Font"): Dictionary{Name("F1"): bad}}); err != nil {
		t.Fatal(err)
	}
	f := fm.Fonts["F1"]
	if f == nil || f.Width(Code{65, 1}) != 600 || f.Encoding[65] != "A" {
		t.Fatalf("font F1 = %+v, want it kept with widths and encoding", f)
	}
	if !strings.Contains(logged.String(), "ToUnicode") {
		t.Errorf("log = %q, want a ToUnicode warning", logged.String())
	}
}

// Without /Encoding a Type 3 font's codes stay bytes: guessing names from
// /CharProcs in Go map order gave a different text on every run.
func TestType3_NoCharProcsGuess(t *testing.T) {
	f := mustParseFont(t, Dictionary{
		Name("Subtype"): Name("Type3"), Name("FirstChar"): Integer(65), Name("LastChar"): Integer(66),
		Name("CharProcs"): Dictionary{Name("a"): Stream{}, Name("b"): Stream{}},
	})
	if len(f.Encoding) != 0 || f.DecodeString("AB") != "AB" {
		t.Errorf("Encoding %v, DecodeString %q; want no encoding and \"AB\"", f.Encoding, f.DecodeString("AB"))
	}
}
