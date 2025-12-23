package pdf

import (
	"testing"
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

	// Do Unknown
	in.Stack = []Object{Name("Unknown")}
	// Should return error
	if err := in.handleXObject("Do"); err == nil {
		t.Error("Expected error for Unknown XObject, got nil")
	}
	// Images count should remain 1
	if len(in.Images) != 1 {
		t.Errorf("Expected 1 image, got %d", len(in.Images))
	}

	// Missing Resources
	inNoRes := NewInterpreter(nil, nil)
	inNoRes.Stack = []Object{Name("Im1")}
	if err := inNoRes.handleXObject("Do"); err == nil {
		t.Error("Expected error for missing resources, got nil")
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
