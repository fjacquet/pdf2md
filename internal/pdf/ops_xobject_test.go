package pdf

import (
	"testing"

	"github.com/fjacquet/pdf2md/internal/types"
)

func TestInterpreter_HandleXObject(t *testing.T) {
	// Setup resources - need valid RGB data for 2x2 image (12 bytes for RGB)
	imageData := make([]byte, 12) // 2x2 pixels * 3 components
	for i := range imageData {
		imageData[i] = byte(i * 20) // Some variation in color
	}

	imageStream := Stream{
		Dictionary: Dictionary{
			Name("Type"):       Name("XObject"),
			Name("Subtype"):    Name("Image"),
			Name("Width"):      Integer(2),
			Name("Height"):     Integer(2),
			Name("ColorSpace"): Name("DeviceRGB"),
		},
		Data: imageData,
	}

	resources := Dictionary{
		Name("XObject"): Dictionary{
			Name("Im1"): imageStream,
		},
	}

	in := NewInterpreter(nil, resources)

	// Do Im1
	in.Stack = []Object{Name("Im1")}
	if err := in.handleXObject("Do"); err != nil {
		t.Errorf("Do failed: %v", err)
	}

	if len(in.Images) != 1 {
		t.Errorf("Expected 1 image, got %d", len(in.Images))
	}
	if len(in.Images) > 0 && in.Images[0].ID != "Im1" {
		t.Errorf("Expected image ID 'Im1', got '%s'", in.Images[0].ID)
	}

	// Do Unknown: warns and keeps the page going
	in.Stack = []Object{Name("Unknown")}
	if err := in.handleXObject("Do"); err != nil {
		t.Errorf("unknown XObject: err = %v, want nil (warning only)", err)
	}
	// Images count should remain 1
	if len(in.Images) != 1 {
		t.Errorf("Expected 1 image, got %d", len(in.Images))
	}

	// Missing Resources
	inNoRes := NewInterpreter(nil, nil)
	inNoRes.Stack = []Object{Name("Im1")}
	if err := inNoRes.handleXObject("Do"); err != nil || len(inNoRes.Images) != 0 {
		t.Errorf("missing resources: err = %v, %d images; want nil, 0", err, len(inNoRes.Images))
	}
}

func TestInterpreter_ExtractImage_Filters(t *testing.T) {
	// Test with filters (mock DecodeStream if possible, or use Identity/Flate if supported)
	// Since DecodeStream is in another file and hard to mock without dependency injection,
	// we will test the logic flow.

	// Filter: DCTDecode (JPEG) - JPEG data is passed through directly
	jpegStream := Stream{
		Dictionary: Dictionary{
			Name("Subtype"): Name("Image"),
			Name("Filter"):  Name("DCTDecode"),
			Name("Width"):   Integer(1),
			Name("Height"):  Integer(1),
		},
		Data: []byte{0xFF, 0xD8, 0xFF, 0xE0}, // JPEG magic bytes
	}

	resources := Dictionary{
		Name("XObject"): Dictionary{
			Name("ImgJPEG"): jpegStream,
		},
	}

	in := NewInterpreter(nil, resources)
	in.Stack = []Object{Name("ImgJPEG")}
	if err := in.handleXObject("Do"); err != nil {
		t.Errorf("Do JPEG failed: %v", err)
	}

	if len(in.Images) != 1 {
		t.Fatal("Expected 1 image")
	}
	if in.Images[0].Format != "jpeg" {
		t.Errorf("Expected jpeg format, got %s", in.Images[0].Format)
	}
}

// A form's own /Resources /Font overrides the page font of the same name,
// inside the form only.
func TestFormFonts(t *testing.T) {
	fm := NewFontManager(nil)
	fm.Fonts["F1"] = fixedWidthFont()
	wide := Dictionary{
		Name("Type"): Name("Font"), Name("Subtype"): Name("Type1"), Name("BaseFont"): Name("Wide"),
		Name("FirstChar"): Integer(65), Name("LastChar"): Integer(66),
		Name("Widths"): Array{Integer(1000), Integer(1000)},
	}
	form := Stream{
		Dictionary: Dictionary{
			Name("Type"): Name("XObject"), Name("Subtype"): Name("Form"),
			Name("Resources"): Dictionary{Name("Font"): Dictionary{Name("F1"): wide}},
		},
		Data: []byte("BT /F1 10 Tf (AB) Tj ET"),
	}
	page := Dictionary{Name("XObject"): Dictionary{Name("Fm1"): form}}
	blocks, _, _, err := NewInterpreter(fm, page).Process([]byte("/Fm1 Do BT /F1 10 Tf (AB) Tj ET"))
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(blocks) != 2 {
		t.Fatalf("got %d blocks, want 2", len(blocks))
	}
	if blocks[0].Width != 20 || blocks[1].Width != 10 {
		t.Errorf("widths = %v, %v; want 20 (form font) and 10 (page font)", blocks[0].Width, blocks[1].Width)
	}
}

// form returns a Form XObject stream running content, with optional resources.
func form(content string, resources Dictionary, matrix Array) Stream {
	d := Dictionary{Name("Type"): Name("XObject"), Name("Subtype"): Name("Form")}
	if resources != nil {
		d[Name("Resources")] = resources
	}
	if matrix != nil {
		d[Name("Matrix")] = matrix
	}
	return Stream{Dictionary: d, Data: []byte(content)}
}

func runPage(t *testing.T, content string, resources Dictionary) []types.TextBlock {
	t.Helper()
	fm := NewFontManager(nil)
	fm.Fonts["F1"] = fixedWidthFont()
	blocks, _, _, err := NewInterpreter(fm, resources).Process([]byte(content))
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	return blocks
}

// Run alone: the failure mode is a fatal stack overflow. A form without
// /Resources inherits the page's, which hold the form itself.
func TestForm_SelfReference(t *testing.T) {
	page := Dictionary{Name("XObject"): Dictionary{Name("Fm1"): form("/Fm1 Do", nil, nil)}}
	runPage(t, "/Fm1 Do BT /F1 10 Tf (after) Tj ET", page)
}

// Bounded by depth, not by a seen-set: a form drawn twice renders twice.
func TestForm_DrawnTwice(t *testing.T) {
	page := Dictionary{Name("XObject"): Dictionary{Name("Fm1"): form("BT /F1 10 Tf (logo) Tj ET", nil, nil)}}
	if b := runPage(t, "/Fm1 Do /Fm1 Do", page); len(b) != 2 {
		t.Errorf("got %d blocks, want 2", len(b))
	}
}

// §8.10.1: the form's /Matrix applies before the CTM. poppler puts this text at x = 2·(100+10) = 220.
func TestForm_MatrixOrder(t *testing.T) {
	m := Array{Integer(1), Integer(0), Integer(0), Integer(1), Integer(100), Integer(0)}
	page := Dictionary{Name("XObject"): Dictionary{Name("Fm1"): form("BT /F1 12 Tf 10 300 Td (A) Tj ET", nil, m)}}
	b := runPage(t, "2 0 0 2 0 0 cm /Fm1 Do", page)
	if len(b) != 1 || b[0].X != 220 {
		t.Errorf("blocks = %+v, want one block at X 220", b)
	}
}

// One unknown name used to fail the whole page.
func TestXObject_UnknownNameKeepsPage(t *testing.T) {
	page := Dictionary{Name("XObject"): Dictionary{}}
	if b := runPage(t, "/Nope Do BT /F1 10 Tf (kept) Tj ET", page); len(b) != 1 || b[0].Text != "kept" {
		t.Errorf("blocks = %+v, want the kept text", b)
	}
}

// /XObject may be an indirect reference (ISO 32000 §7.8.3).
func TestXObject_IndirectDictionary(t *testing.T) {
	r, err := NewReader(writeTestPDF(t,
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [] /Count 0 >>",
		"<< /Im1 4 0 R >>",
		"<< /Type /XObject /Subtype /Image /Width 1 /Height 1 /ColorSpace /DeviceRGB /BitsPerComponent 8 /Length 3 >>\nstream\n\xff\x00\x00\nendstream",
	))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()
	page := Dictionary{Name("XObject"): IndirectRef{ObjectNumber: 3}}
	_, images, _, err := NewInterpreter(NewFontManager(r), page).Process([]byte("q 10 0 0 10 0 0 cm /Im1 Do Q"))
	if err != nil || len(images) != 1 {
		t.Errorf("got %d images, err %v; want 1, nil", len(images), err)
	}
}
