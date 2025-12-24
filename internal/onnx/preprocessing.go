package onnx

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"

	"golang.org/x/image/draw"
)

// Preprocess prepares an image for YOLO inference.
// Returns the preprocessed tensor and transformation parameters for postprocessing.
//
// This is a pure function: same input always produces same output.
//
// The preprocessing pipeline:
// 1. Letterbox resize (maintain aspect ratio with padding)
// 2. Pad to exact targetSize x targetSize
// 3. Convert to BCHW float32 tensor
// 4. Normalize to [0, 1]
func Preprocess(img image.Image, targetSize, stride int) (tensor []float32, scale, padX, padY float64) {
	// Calculate letterbox dimensions
	bounds := img.Bounds()
	origW, origH := bounds.Dx(), bounds.Dy()

	// Calculate scale to fit within targetSize while maintaining aspect ratio
	scale = math.Min(float64(targetSize)/float64(origW), float64(targetSize)/float64(origH))

	// Calculate actual resized dimensions
	resizedW := int(float64(origW) * scale)
	resizedH := int(float64(origH) * scale)

	// Calculate padding to center the image in targetSize x targetSize canvas
	padX = float64(targetSize-resizedW) / 2
	padY = float64(targetSize-resizedH) / 2

	// Create padded canvas with gray background (114, 114, 114)
	// Output is always exactly targetSize x targetSize
	padded := PadImage(
		ResizeImage(img, resizedW, resizedH),
		targetSize,
		targetSize,
		int(padX),
		int(padY),
	)

	// Convert to tensor
	tensor = ImageToTensor(padded)

	return tensor, scale, padX, padY
}

// LetterboxScale calculates the scale factor to fit an image in target dimensions
// while maintaining aspect ratio.
//
// This is a pure function.
func LetterboxScale(origW, origH, targetW, targetH int) (scale float64, newW, newH int) {
	scaleW := float64(targetW) / float64(origW)
	scaleH := float64(targetH) / float64(origH)
	scale = math.Min(scaleW, scaleH)

	newW = int(float64(origW) * scale)
	newH = int(float64(origH) * scale)

	return scale, newW, newH
}

// alignToStride rounds up to the nearest multiple of stride.
func alignToStride(size, stride int) int {
	if size%stride == 0 {
		return size
	}
	return ((size / stride) + 1) * stride
}

// ResizeImage resizes an image to the specified dimensions using high-quality interpolation.
//
// This is a pure function.
func ResizeImage(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

// PadImage creates a new image with the given dimensions, centered content, and gray padding.
//
// This is a pure function.
func PadImage(img image.Image, width, height, padLeft, padTop int) image.Image {
	// Create canvas with gray background (114, 114, 114) - YOLO standard
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	gray := color.RGBA{R: 114, G: 114, B: 114, A: 255}

	// Fill with gray
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(x, y, gray)
		}
	}

	// Draw the image at the padded position
	bounds := img.Bounds()
	draw.Draw(dst, image.Rect(padLeft, padTop, padLeft+bounds.Dx(), padTop+bounds.Dy()),
		img, bounds.Min, draw.Over)

	return dst
}

// ImageToTensor converts an image to a BCHW float32 tensor normalized to [0, 1].
//
// This is a pure function.
//
// Output shape: [1, 3, height, width] flattened to []float32
// Channel order: RGB
func ImageToTensor(img image.Image) []float32 {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Allocate tensor: batch=1, channels=3, height, width
	tensor := make([]float32, 3*height*width)

	// Fill tensor in BCHW format (batch dimension is implicit as size 1)
	// Channel 0 = R, Channel 1 = G, Channel 2 = B
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()

			// RGBA returns 16-bit values, convert to 8-bit and normalize to [0, 1]
			idx := y*width + x
			tensor[0*height*width+idx] = float32(r>>8) / 255.0 // R channel
			tensor[1*height*width+idx] = float32(g>>8) / 255.0 // G channel
			tensor[2*height*width+idx] = float32(b>>8) / 255.0 // B channel
		}
	}

	return tensor
}

// TensorShape returns the shape of a BCHW tensor given image dimensions.
// Returns [1, 3, height, width].
func TensorShape(width, height int) []int64 {
	return []int64{1, 3, int64(height), int64(width)}
}

// DefaultImageDecoder decodes PNG image bytes to image.Image.
func DefaultImageDecoder(data []byte) (image.Image, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, ErrInvalidImage
	}
	return img, nil
}

// ImageTransform is a function type for image transformations.
// This supports functional composition of preprocessing steps.
type ImageTransform func(image.Image) image.Image

// Compose creates a new ImageTransform that applies transforms in order.
//
// This enables functional composition:
//
//	preprocess := Compose(
//	    Resize(1024, 1024),
//	    Letterbox(32),
//	)
func Compose(transforms ...ImageTransform) ImageTransform {
	return func(img image.Image) image.Image {
		for _, t := range transforms {
			img = t(img)
		}
		return img
	}
}

// ResizeTransform returns an ImageTransform that resizes to the given dimensions.
func ResizeTransform(width, height int) ImageTransform {
	return func(img image.Image) image.Image {
		return ResizeImage(img, width, height)
	}
}
