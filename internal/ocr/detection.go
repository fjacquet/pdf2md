package ocr

import (
	"image"
	"math"

	"golang.org/x/image/draw"
)

// ImageNet normalization constants used by PaddleOCR's detection models.
var (
	detMean = [3]float32{0.485, 0.456, 0.406}
	detStd  = [3]float32{0.229, 0.224, 0.225}
)

// DetectionBox is an axis-aligned bounding box in image-pixel coordinates
// (origin top-left, Y increases downward).
type DetectionBox struct {
	MinX, MinY int
	MaxX, MaxY int
	Score      float64
}

// Width returns the box width in pixels.
func (b DetectionBox) Width() int { return b.MaxX - b.MinX }

// Height returns the box height in pixels.
func (b DetectionBox) Height() int { return b.MaxY - b.MinY }

// PreprocessDet prepares an image for DB-style text detection.
//
// Pure function. Pipeline:
//  1. Resize so short side == shortSide, long side aligned up to stride.
//  2. Convert to CHW float32 tensor normalized with ImageNet mean/std.
//
// Returns tensor, resized width, resized height, and scale (resized / original).
// The caller uses scale to map detected boxes back to original-image pixels.
func PreprocessDet(img image.Image, shortSide, stride int) (tensor []float32, resizedW, resizedH int, scaleX, scaleY float64) {
	bounds := img.Bounds()
	origW, origH := bounds.Dx(), bounds.Dy()

	scale := float64(shortSide) / float64(minInt(origW, origH))
	resizedW = alignUp(int(math.Round(float64(origW)*scale)), stride)
	resizedH = alignUp(int(math.Round(float64(origH)*scale)), stride)

	scaleX = float64(resizedW) / float64(origW)
	scaleY = float64(resizedH) / float64(origH)

	resized := resizeRGBA(img, resizedW, resizedH)
	tensor = normalizeToTensor(resized, detMean, detStd)
	return tensor, resizedW, resizedH, scaleX, scaleY
}

// PostprocessDet converts a DB probability map to bounding boxes.
//
// Pure function. Pipeline:
//  1. Threshold the probability map at dbThreshold to get a binary mask.
//  2. Label connected components (4-connectivity).
//  3. Compute the bounding box and mean probability per component.
//  4. Discard boxes where mean prob < boxThreshold or short side < minSize.
//  5. Scale boxes back to original-image pixels via scaleX/scaleY.
//
// probMap is a flat [mapH][mapW] slice of probabilities in [0, 1].
// Connected components are simple and sufficient for near-horizontal PDF text;
// the Suzuki-Abe contour approach used by PaddleOCR is overkill here.
func PostprocessDet(probMap []float32, mapW, mapH int, scaleX, scaleY, dbThreshold, boxThreshold float64, minSize int) []DetectionBox {
	if len(probMap) != mapW*mapH {
		return nil
	}
	mask := thresholdMask(probMap, dbThreshold)
	labels, nLabels := labelComponents(mask, mapW, mapH)
	if nLabels == 0 {
		return nil
	}
	boxes := componentBoxes(probMap, labels, nLabels, mapW, mapH)
	return filterAndScaleBoxes(boxes, scaleX, scaleY, boxThreshold, minSize)
}

// thresholdMask returns a flat bool slice where prob >= threshold.
func thresholdMask(probMap []float32, threshold float64) []bool {
	mask := make([]bool, len(probMap))
	t := float32(threshold)
	for i, p := range probMap {
		mask[i] = p >= t
	}
	return mask
}

// labelComponents assigns a unique positive label to each 4-connected
// foreground component. Background is 0. Returns labels slice and count.
// Two-pass union-find keeps this O(N·α(N)).
func labelComponents(mask []bool, w, h int) ([]int, int) {
	labels := make([]int, len(mask))
	uf := newUnionFind()

	for y := range h {
		for x := range w {
			idx := y*w + x
			if !mask[idx] {
				continue
			}
			left, up := 0, 0
			if x > 0 {
				left = labels[idx-1]
			}
			if y > 0 {
				up = labels[idx-w]
			}
			switch {
			case left == 0 && up == 0:
				labels[idx] = uf.make()
			case left != 0 && up == 0:
				labels[idx] = left
			case left == 0 && up != 0:
				labels[idx] = up
			default:
				labels[idx] = left
				uf.union(left, up)
			}
		}
	}

	// Second pass: flatten to canonical labels 1..N.
	remap := make(map[int]int)
	next := 0
	for i, l := range labels {
		if l == 0 {
			continue
		}
		root := uf.find(l)
		canonical, ok := remap[root]
		if !ok {
			next++
			canonical = next
			remap[root] = canonical
		}
		labels[i] = canonical
	}
	return labels, next
}

// componentBoxes aggregates bounding boxes and mean scores per label.
func componentBoxes(probMap []float32, labels []int, nLabels, w, h int) []DetectionBox {
	boxes := make([]DetectionBox, nLabels)
	sums := make([]float64, nLabels)
	counts := make([]int, nLabels)

	for i := range boxes {
		boxes[i] = DetectionBox{MinX: w, MinY: h, MaxX: 0, MaxY: 0}
	}
	for y := range h {
		for x := range w {
			l := labels[y*w+x]
			if l == 0 {
				continue
			}
			b := &boxes[l-1]
			if x < b.MinX {
				b.MinX = x
			}
			if x > b.MaxX {
				b.MaxX = x
			}
			if y < b.MinY {
				b.MinY = y
			}
			if y > b.MaxY {
				b.MaxY = y
			}
			sums[l-1] += float64(probMap[y*w+x])
			counts[l-1]++
		}
	}
	for i := range boxes {
		if counts[i] > 0 {
			boxes[i].Score = sums[i] / float64(counts[i])
			boxes[i].MaxX++ // convert inclusive-max to exclusive-max for Width/Height
			boxes[i].MaxY++
		}
	}
	return boxes
}

// filterAndScaleBoxes drops low-confidence/small boxes and scales back to
// original-image pixels.
func filterAndScaleBoxes(boxes []DetectionBox, scaleX, scaleY, boxThreshold float64, minSize int) []DetectionBox {
	out := make([]DetectionBox, 0, len(boxes))
	for _, b := range boxes {
		if b.Score < boxThreshold {
			continue
		}
		if minInt(b.Width(), b.Height()) < minSize {
			continue
		}
		out = append(out, DetectionBox{
			MinX:  int(math.Round(float64(b.MinX) / scaleX)),
			MinY:  int(math.Round(float64(b.MinY) / scaleY)),
			MaxX:  int(math.Round(float64(b.MaxX) / scaleX)),
			MaxY:  int(math.Round(float64(b.MaxY) / scaleY)),
			Score: b.Score,
		})
	}
	return out
}

// resizeRGBA resizes using Catmull-Rom (high quality, matches onnx package).
func resizeRGBA(img image.Image, width, height int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

// normalizeToTensor converts an image to a CHW float32 tensor normalized
// with ImageNet mean/std. Channel order is RGB.
//
// Output shape: [1, 3, height, width] flattened.
func normalizeToTensor(img *image.RGBA, mean, std [3]float32) []float32 {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	plane := h * w
	tensor := make([]float32, 3*plane)

	stride := img.Stride
	for y := range h {
		row := img.Pix[y*stride : y*stride+w*4]
		for x := range w {
			r := float32(row[x*4]) / 255.0
			g := float32(row[x*4+1]) / 255.0
			b := float32(row[x*4+2]) / 255.0
			idx := y*w + x
			tensor[idx] = (r - mean[0]) / std[0]
			tensor[plane+idx] = (g - mean[1]) / std[1]
			tensor[2*plane+idx] = (b - mean[2]) / std[2]
		}
	}
	return tensor
}

// alignUp rounds size up to the nearest multiple of stride.
func alignUp(size, stride int) int {
	if size%stride == 0 {
		return size
	}
	return ((size / stride) + 1) * stride
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// unionFind is a tiny disjoint-set for connected-component labelling.
type unionFind struct{ parent []int }

func newUnionFind() *unionFind { return &unionFind{parent: []int{0}} }

func (u *unionFind) make() int {
	id := len(u.parent)
	u.parent = append(u.parent, id)
	return id
}

func (u *unionFind) find(x int) int {
	for u.parent[x] != x {
		u.parent[x] = u.parent[u.parent[x]]
		x = u.parent[x]
	}
	return x
}

func (u *unionFind) union(a, b int) {
	ra, rb := u.find(a), u.find(b)
	if ra == rb {
		return
	}
	if ra < rb {
		u.parent[rb] = ra
	} else {
		u.parent[ra] = rb
	}
}
