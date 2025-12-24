package layout

import (
	"image"
	"log"

	"github.com/fjacquet/pdf2md/internal/extractor"
	"github.com/fjacquet/pdf2md/internal/onnx"
	"github.com/fjacquet/pdf2md/internal/types"
)

// LayoutDetector defines the interface for ONNX-based layout detection.
// This interface is defined here (consumer side) per Go idioms.
type LayoutDetector interface {
	DetectLayout(pageImage image.Image, pageWidth, pageHeight float64) (*types.PageDetections, error)
	Close() error
	IsAvailable() bool
}

// AnalyzerWithONNX extends Analyzer with ONNX detection capabilities.
type AnalyzerWithONNX struct {
	*Analyzer
	detector LayoutDetector
}

// NewAnalyzerWithONNX creates an analyzer with ONNX support.
// If detector is nil, falls back to rule-based analysis only.
func NewAnalyzerWithONNX(config *Config, detector LayoutDetector) *AnalyzerWithONNX {
	return &AnalyzerWithONNX{
		Analyzer: NewAnalyzerWithConfig(config),
		detector: detector,
	}
}

// AnalyzeWithImage analyzes page content with optional ONNX detection.
// If pageImage is provided and ONNX detector is available, uses ML detection.
// Otherwise falls back to rule-based analysis.
func (a *AnalyzerWithONNX) AnalyzeWithImage(content *extractor.PageContent, pageImage image.Image) []Element {
	// First, run the standard rule-based analysis
	elements := a.Analyzer.Analyze(content)

	// If no detector or no image, return rule-based results
	if a.detector == nil || !a.detector.IsAvailable() || pageImage == nil {
		return elements
	}

	// Run ONNX detection
	detections, err := a.detector.DetectLayout(pageImage, content.PageWidth, content.PageHeight)
	if err != nil {
		log.Printf("ONNX detection failed, using rule-based fallback: %v", err)
		return elements
	}

	// Fuse ONNX detections with rule-based elements
	return a.fuseDetections(elements, detections)
}

// fuseDetections merges ONNX detections with rule-based elements.
// ONNX detections with high confidence override rule-based classification.
func (a *AnalyzerWithONNX) fuseDetections(elements []Element, detections *types.PageDetections) []Element {
	if detections == nil || len(detections.Detections) == 0 {
		return elements
	}

	minConfidence := 0.5 // Default minimum confidence for override
	if a.Config != nil && a.Config.MinONNXConfidence > 0 {
		minConfidence = a.Config.MinONNXConfidence
	}

	// Create a copy to avoid modifying the original
	result := make([]Element, len(elements))
	copy(result, elements)

	// For each element, find the best matching ONNX detection
	for i := range result {
		if result[i].Type == ElementTypeImage {
			continue // Don't override images
		}

		match := a.findBestMatch(&result[i], detections)
		if match != nil && match.Confidence >= minConfidence {
			newType := mapONNXClassToElementType(match.ClassID)
			if newType != types.ElementTypeUnknown && newType != result[i].Type {
				result[i].Type = newType
				result[i].ONNXClassID = match.ClassID
				result[i].ONNXConfidence = match.Confidence
			}
		}
	}

	return result
}

// findBestMatch finds the ONNX detection that best overlaps with the element.
func (a *AnalyzerWithONNX) findBestMatch(elem *Element, detections *types.PageDetections) *types.BoundingBox {
	if elem == nil || detections == nil {
		return nil
	}

	// Convert element to bounding box
	elemBox := types.BoundingBox{
		X1: elem.X,
		Y1: elem.Y,
		X2: elem.X + elem.Width,
		Y2: elem.Y + elem.Height,
	}

	var bestMatch *types.BoundingBox
	var bestIoU float64

	for i := range detections.Detections {
		det := &detections.Detections[i]
		iou := elemBox.IoU(*det)
		if iou > bestIoU && iou > 0.3 { // Minimum 30% overlap
			bestIoU = iou
			bestMatch = det
		}
	}

	return bestMatch
}

// mapONNXClassToElementType converts ONNX class ID to ElementType.
func mapONNXClassToElementType(classID int) types.ElementType {
	class := types.DocLayoutClass(classID)
	return class.ToElementType()
}

// CreateONNXConfig creates an onnx.Config from layout Config.
func CreateONNXConfig(cfg *Config) *onnx.Config {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	return &onnx.Config{
		ModelPath:     cfg.ONNXModelPath,
		RuntimePath:   cfg.ONNXRuntimePath,
		ConfThreshold: cfg.ONNXConfThreshold,
		NMSThreshold:  cfg.ONNXNMSThreshold,
		InputSize:     cfg.ONNXInputSize,
		Stride:        cfg.ONNXStride,
		UseCoreML:     cfg.ONNXUseCoreML,
	}
}
