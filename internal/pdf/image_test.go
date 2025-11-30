package pdf

import (
	"testing"
)

func TestImageExtraction(t *testing.T) {
	// Create a mock interpreter with resources
	fm := &FontManager{}

	// Create a mock image stream
	imageStream := &Stream{
		Dictionary: Dictionary{
			Name("Type"):             Name("XObject"),
			Name("Subtype"):          Name("Image"),
			Name("Width"):            Integer(100),
			Name("Height"):           Integer(100),
			Name("ColorSpace"):       Name("DeviceRGB"),
			Name("BitsPerComponent"): Integer(8),
			// Name("Filter"):   Name("FlateDecode"), // Removed to avoid zlib error
		},
		Data: []byte{1, 2, 3, 4}, // Mock data
	}

	// Create resources dictionary with XObject
	resources := Dictionary{
		Name("XObject"): Dictionary{
			Name("Im1"): imageStream, // Direct object for simplicity in this test
		},
	}

	interpreter := NewInterpreter(fm, resources)

	// Mock content stream: q 100 0 0 100 50 50 cm /Im1 Do Q
	// q: Save state
	// 100 0 0 100 50 50 cm: Scale 100x100, translate to 50,50
	// /Im1 Do: Draw image Im1
	// Q: Restore state
	content := "q 100 0 0 100 50 50 cm /Im1 Do Q"

	// Process
	_, images, _, err := interpreter.Process([]byte(content))
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	if len(images) != 1 {
		t.Fatalf("Expected 1 image, got %d", len(images))
	}

	img := images[0]
	if img.ID != "Im1" {
		t.Errorf("Expected ID Im1, got %s", img.ID)
	}
	if img.Width != 100 {
		t.Errorf("Expected Width 100, got %f", img.Width)
	}
	if img.Height != 100 {
		t.Errorf("Expected Height 100, got %f", img.Height)
	}
	if img.X != 50 {
		t.Errorf("Expected X 50, got %f", img.X)
	}
	if img.Y != 50 {
		t.Errorf("Expected Y 50, got %f", img.Y)
	}
}
