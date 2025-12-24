package onnx

import (
	"image"
	"image/color"
	"testing"
)

// createTestImage creates a solid color test image of specified size.
func createTestImage(width, height int, c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func TestLetterboxScale(t *testing.T) {
	tests := []struct {
		name                  string
		origW, origH          int
		targetW, targetH      int
		wantScale             float64
		wantNewW, wantNewH    int
	}{
		{
			name:      "square_to_square",
			origW:     100, origH: 100,
			targetW:   200, targetH: 200,
			wantScale: 2.0, wantNewW: 200, wantNewH: 200,
		},
		{
			name:      "wide_to_square",
			origW:     200, origH: 100,
			targetW:   200, targetH: 200,
			wantScale: 1.0, wantNewW: 200, wantNewH: 100,
		},
		{
			name:      "tall_to_square",
			origW:     100, origH: 200,
			targetW:   200, targetH: 200,
			wantScale: 1.0, wantNewW: 100, wantNewH: 200,
		},
		{
			name:      "downscale",
			origW:     1000, origH: 800,
			targetW:   500, targetH: 500,
			wantScale: 0.5, wantNewW: 500, wantNewH: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scale, newW, newH := LetterboxScale(tt.origW, tt.origH, tt.targetW, tt.targetH)
			if scale != tt.wantScale {
				t.Errorf("scale = %v, want %v", scale, tt.wantScale)
			}
			if newW != tt.wantNewW {
				t.Errorf("newW = %v, want %v", newW, tt.wantNewW)
			}
			if newH != tt.wantNewH {
				t.Errorf("newH = %v, want %v", newH, tt.wantNewH)
			}
		})
	}
}

func TestResizeImage(t *testing.T) {
	img := createTestImage(100, 100, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	resized := ResizeImage(img, 50, 50)
	bounds := resized.Bounds()

	if bounds.Dx() != 50 || bounds.Dy() != 50 {
		t.Errorf("ResizeImage size = %dx%d, want 50x50", bounds.Dx(), bounds.Dy())
	}
}

func TestPadImage(t *testing.T) {
	img := createTestImage(50, 50, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	padded := PadImage(img, 100, 100, 25, 25)
	bounds := padded.Bounds()

	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("PadImage size = %dx%d, want 100x100", bounds.Dx(), bounds.Dy())
	}

	// Check padding color (should be gray 114)
	r, g, b, _ := padded.At(0, 0).RGBA()
	// RGBA() returns 16-bit values, shift to compare 8-bit
	if r>>8 != 114 || g>>8 != 114 || b>>8 != 114 {
		t.Errorf("PadImage corner color = (%d, %d, %d), want gray (114)", r>>8, g>>8, b>>8)
	}

	// Check center color (should be red)
	r, g, b, _ = padded.At(50, 50).RGBA()
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 {
		t.Errorf("PadImage center color = (%d, %d, %d), want red (255, 0, 0)", r>>8, g>>8, b>>8)
	}
}

func TestImageToTensor(t *testing.T) {
	// Create a 2x2 test image with known colors
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})   // Red
	img.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})   // Green
	img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})   // Blue
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255}) // White

	tensor := ImageToTensor(img)

	// Expected tensor size: 3 channels * 2 * 2 = 12
	if len(tensor) != 12 {
		t.Fatalf("tensor length = %d, want 12", len(tensor))
	}

	// Check tensor values (BCHW format)
	// Channel 0 (R): [1.0, 0, 0, 1.0]
	// Channel 1 (G): [0, 1.0, 0, 1.0]
	// Channel 2 (B): [0, 0, 1.0, 1.0]

	// Red channel
	if tensor[0] != 1.0 { // pixel (0,0) R
		t.Errorf("tensor[0] = %v, want 1.0 (red pixel R)", tensor[0])
	}
	if tensor[1] != 0.0 { // pixel (1,0) R
		t.Errorf("tensor[1] = %v, want 0.0 (green pixel R)", tensor[1])
	}

	// Green channel (starts at index 4)
	if tensor[4] != 0.0 { // pixel (0,0) G
		t.Errorf("tensor[4] = %v, want 0.0 (red pixel G)", tensor[4])
	}
	if tensor[5] != 1.0 { // pixel (1,0) G
		t.Errorf("tensor[5] = %v, want 1.0 (green pixel G)", tensor[5])
	}

	// Blue channel (starts at index 8)
	if tensor[10] != 1.0 { // pixel (0,1) B
		t.Errorf("tensor[10] = %v, want 1.0 (blue pixel B)", tensor[10])
	}
}

func TestTensorShape(t *testing.T) {
	shape := TensorShape(1024, 1024)

	if len(shape) != 4 {
		t.Fatalf("shape length = %d, want 4", len(shape))
	}

	if shape[0] != 1 || shape[1] != 3 || shape[2] != 1024 || shape[3] != 1024 {
		t.Errorf("shape = %v, want [1, 3, 1024, 1024]", shape)
	}
}

func TestPreprocess(t *testing.T) {
	// Create a test image
	img := createTestImage(800, 600, color.RGBA{R: 128, G: 128, B: 128, A: 255})

	tensor, scale, padX, padY := Preprocess(img, 1024, 32)

	// Check tensor is not empty
	if len(tensor) == 0 {
		t.Error("tensor is empty")
	}

	// Check scale is reasonable
	if scale <= 0 || scale > 2 {
		t.Errorf("scale = %v, want > 0 and <= 2", scale)
	}

	// Check padding is non-negative
	if padX < 0 || padY < 0 {
		t.Errorf("padding = (%v, %v), want non-negative", padX, padY)
	}
}

func TestCompose(t *testing.T) {
	img := createTestImage(100, 100, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	// Compose two resize operations
	transform := Compose(
		ResizeTransform(50, 50),
		ResizeTransform(25, 25),
	)

	result := transform(img)
	bounds := result.Bounds()

	if bounds.Dx() != 25 || bounds.Dy() != 25 {
		t.Errorf("Compose result size = %dx%d, want 25x25", bounds.Dx(), bounds.Dy())
	}
}
