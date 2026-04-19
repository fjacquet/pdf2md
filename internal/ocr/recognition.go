package ocr

import (
	"image"

	"golang.org/x/image/draw"
)

// PaddleOCR recognition normalization: (pixel/255 - 0.5) / 0.5 == pixel/127.5 - 1.0.
// Produces values in [-1, 1] with channel order RGB.

// PreprocessRec crops a bounding box from the page image, resizes it to the
// model's input height keeping aspect ratio (truncated to maxWidth), and
// produces a CHW float32 tensor in [-1, 1].
//
// Pure function. Returns the tensor and the actual tensor width (≤ maxWidth),
// which the caller needs to pass as the dynamic W dimension in the ONNX
// input shape [1, 3, H, W].
//
// A crop wider than maxWidth is truncated rather than split: PaddleOCR's
// recognition head handles up to ~320-px wide strips in one pass, and PDF
// text blocks rarely exceed this at 300-DPI rendering. If users report
// truncation in practice, the right fix is raising RecInputMaxWidth, not
// adding split-and-stitch complexity here.
func PreprocessRec(pageImg image.Image, box DetectionBox, height, maxWidth int) ([]float32, int) {
	crop := cropImage(pageImg, box)
	if crop == nil {
		return nil, 0
	}
	bounds := crop.Bounds()
	cw, ch := bounds.Dx(), bounds.Dy()
	if cw == 0 || ch == 0 {
		return nil, 0
	}

	ratio := float64(cw) / float64(ch)
	targetW := max(min(int(float64(height)*ratio), maxWidth), 1)

	resized := resizeRGBA(crop, targetW, height)
	tensor := recTensor(resized, maxWidth)
	return tensor, maxWidth
}

// cropImage returns the sub-image inside box, clamped to the page bounds.
// Pure function. Returns nil if the clamped box is empty.
func cropImage(img image.Image, box DetectionBox) image.Image {
	bounds := img.Bounds()
	x0 := clamp(box.MinX, bounds.Min.X, bounds.Max.X)
	y0 := clamp(box.MinY, bounds.Min.Y, bounds.Max.Y)
	x1 := clamp(box.MaxX, bounds.Min.X, bounds.Max.X)
	y1 := clamp(box.MaxY, bounds.Min.Y, bounds.Max.Y)
	if x1-x0 <= 0 || y1-y0 <= 0 {
		return nil
	}
	// SubImage is the fast path when available; otherwise draw into a new RGBA.
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	if si, ok := img.(subImager); ok {
		return si.SubImage(image.Rect(x0, y0, x1, y1))
	}
	dst := image.NewRGBA(image.Rect(0, 0, x1-x0, y1-y0))
	draw.Draw(dst, dst.Bounds(), img, image.Point{X: x0, Y: y0}, draw.Src)
	return dst
}

// recTensor normalizes an RGBA image into a [1, 3, H, maxWidth] tensor.
// Crops narrower than maxWidth are right-padded with zeros (matches
// PaddleOCR's inference behavior).
//
// Pure function.
func recTensor(img *image.RGBA, maxWidth int) []float32 {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w > maxWidth {
		w = maxWidth
	}
	plane := h * maxWidth
	tensor := make([]float32, 3*plane)

	stride := img.Stride
	for y := range h {
		row := img.Pix[y*stride : y*stride+w*4]
		for x := range w {
			r := float32(row[x*4])/127.5 - 1.0
			g := float32(row[x*4+1])/127.5 - 1.0
			b := float32(row[x*4+2])/127.5 - 1.0
			idx := y*maxWidth + x
			tensor[idx] = r
			tensor[plane+idx] = g
			tensor[2*plane+idx] = b
		}
	}
	return tensor
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
