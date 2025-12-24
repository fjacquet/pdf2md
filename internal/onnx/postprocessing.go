package onnx

import (
	"sort"

	"github.com/fjacquet/pdf2md/internal/types"
)

// Postprocess converts raw ONNX output to PageDetections.
// Applies confidence filtering, NMS, and coordinate scaling.
//
// This is a pure function: same inputs always produce same outputs.
//
// Parameters:
//   - rawDetections: model output in format [N][6] = [x1, y1, x2, y2, conf, class]
//   - confThreshold: minimum confidence to keep a detection
//   - nmsThreshold: IoU threshold for Non-Maximum Suppression
//   - scale: preprocessing scale factor
//   - padX, padY: preprocessing padding offsets
//   - pageWidth, pageHeight: target page dimensions in PDF points
func Postprocess(
	rawDetections [][]float32,
	confThreshold, nmsThreshold float64,
	scale, padX, padY float64,
	pageWidth, pageHeight float64,
) []types.BoundingBox {
	// 1. Parse raw detections and filter by confidence
	boxes := ParseDetections(rawDetections, confThreshold)

	// 2. Apply Non-Maximum Suppression
	boxes = ApplyNMS(boxes, nmsThreshold)

	// 3. Scale coordinates from model space to page space
	boxes = ScaleBoxesToPage(boxes, scale, padX, padY, pageWidth, pageHeight)

	return boxes
}

// ParseDetections converts raw model output to BoundingBox slice.
// Filters detections below the confidence threshold.
//
// This is a pure function.
func ParseDetections(rawDetections [][]float32, confThreshold float64) []types.BoundingBox {
	result := make([]types.BoundingBox, 0, len(rawDetections))

	for _, det := range rawDetections {
		if len(det) < 6 {
			continue
		}

		conf := float64(det[4])
		if conf < confThreshold {
			continue
		}

		classID := int(det[5])
		result = append(result, types.BoundingBox{
			X1:         float64(det[0]),
			Y1:         float64(det[1]),
			X2:         float64(det[2]),
			Y2:         float64(det[3]),
			Confidence: conf,
			ClassID:    classID,
			ClassName:  types.DocLayoutClass(classID).String(),
		})
	}

	return result
}

// ApplyNMS performs Non-Maximum Suppression on detections.
// Removes overlapping boxes, keeping the one with highest confidence.
//
// This is a pure function.
//
// Algorithm:
// 1. Sort boxes by confidence (descending)
// 2. For each box, mark overlapping lower-confidence boxes for removal
// 3. Return unmarked boxes
func ApplyNMS(boxes []types.BoundingBox, iouThreshold float64) []types.BoundingBox {
	if len(boxes) == 0 {
		return boxes
	}

	// Sort by confidence (descending)
	sorted := make([]types.BoundingBox, len(boxes))
	copy(sorted, boxes)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Confidence > sorted[j].Confidence
	})

	// Track which boxes to keep
	keep := make([]bool, len(sorted))
	for i := range keep {
		keep[i] = true
	}

	// Apply NMS
	for i := 0; i < len(sorted); i++ {
		if !keep[i] {
			continue
		}

		for j := i + 1; j < len(sorted); j++ {
			if !keep[j] {
				continue
			}

			// Only suppress same class
			if sorted[i].ClassID != sorted[j].ClassID {
				continue
			}

			// Calculate IoU and suppress if above threshold
			if sorted[i].IoU(sorted[j]) > iouThreshold {
				keep[j] = false
			}
		}
	}

	// Collect kept boxes
	result := make([]types.BoundingBox, 0, len(sorted))
	for i, box := range sorted {
		if keep[i] {
			result = append(result, box)
		}
	}

	return result
}

// ScaleBoxesToPage transforms box coordinates from model space to page space.
//
// This is a pure function.
//
// The transformation reverses preprocessing:
// 1. Remove padding offset
// 2. Divide by scale factor
// 3. Optionally clamp to page bounds
func ScaleBoxesToPage(
	boxes []types.BoundingBox,
	scale, padX, padY float64,
	pageWidth, pageHeight float64,
) []types.BoundingBox {
	result := make([]types.BoundingBox, len(boxes))

	for i, box := range boxes {
		result[i] = ScaleBoxToPage(box, scale, padX, padY, pageWidth, pageHeight)
	}

	return result
}

// ScaleBoxToPage transforms a single box from model space to page space.
//
// This is a pure function.
func ScaleBoxToPage(
	box types.BoundingBox,
	scale, padX, padY float64,
	pageWidth, pageHeight float64,
) types.BoundingBox {
	// Remove padding and scale
	x1 := (box.X1 - padX) / scale
	y1 := (box.Y1 - padY) / scale
	x2 := (box.X2 - padX) / scale
	y2 := (box.Y2 - padY) / scale

	// Clamp to page bounds
	x1 = clamp(x1, 0, pageWidth)
	y1 = clamp(y1, 0, pageHeight)
	x2 = clamp(x2, 0, pageWidth)
	y2 = clamp(y2, 0, pageHeight)

	return types.BoundingBox{
		X1:         x1,
		Y1:         y1,
		X2:         x2,
		Y2:         y2,
		Confidence: box.Confidence,
		ClassID:    box.ClassID,
		ClassName:  box.ClassName,
	}
}

// clamp restricts a value to [min, max].
func clamp(v, minVal, maxVal float64) float64 {
	if v < minVal {
		return minVal
	}
	if v > maxVal {
		return maxVal
	}
	return v
}

// FilterByConfidence returns boxes with confidence >= threshold.
//
// This is a pure function.
func FilterByConfidence(boxes []types.BoundingBox, threshold float64) []types.BoundingBox {
	result := make([]types.BoundingBox, 0, len(boxes))
	for _, box := range boxes {
		if box.Confidence >= threshold {
			result = append(result, box)
		}
	}
	return result
}

// FilterByClass returns boxes with the specified class ID.
//
// This is a pure function.
func FilterByClass(boxes []types.BoundingBox, classID int) []types.BoundingBox {
	result := make([]types.BoundingBox, 0, len(boxes))
	for _, box := range boxes {
		if box.ClassID == classID {
			result = append(result, box)
		}
	}
	return result
}

// SortByConfidence returns boxes sorted by confidence (descending).
//
// This is a pure function (returns a new slice).
func SortByConfidence(boxes []types.BoundingBox) []types.BoundingBox {
	result := make([]types.BoundingBox, len(boxes))
	copy(result, boxes)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Confidence > result[j].Confidence
	})
	return result
}

// SortByPosition returns boxes sorted by position (top-to-bottom, left-to-right).
//
// This is a pure function (returns a new slice).
func SortByPosition(boxes []types.BoundingBox) []types.BoundingBox {
	result := make([]types.BoundingBox, len(boxes))
	copy(result, boxes)
	sort.Slice(result, func(i, j int) bool {
		// Primary sort by Y (top to bottom)
		if result[i].Y1 != result[j].Y1 {
			return result[i].Y1 < result[j].Y1
		}
		// Secondary sort by X (left to right)
		return result[i].X1 < result[j].X1
	})
	return result
}
