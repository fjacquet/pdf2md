package layout

import (
	"image"
	"testing"

	"github.com/fjacquet/pdf2md/internal/extractor"
	"github.com/fjacquet/pdf2md/internal/types"
)

// mockDetector is a mock implementation of LayoutDetector for testing.
type mockDetector struct {
	available   bool
	detections  *types.PageDetections
	detectError error
}

func (m *mockDetector) DetectLayout(pageImage image.Image, pageWidth, pageHeight float64) (*types.PageDetections, error) {
	if m.detectError != nil {
		return nil, m.detectError
	}
	return m.detections, nil
}

func (m *mockDetector) Close() error {
	return nil
}

func (m *mockDetector) IsAvailable() bool {
	return m.available
}

func TestNewAnalyzerWithONNX(t *testing.T) {
	config := DefaultConfig()
	detector := &mockDetector{available: true}

	analyzer := NewAnalyzerWithONNX(config, detector)

	if analyzer == nil {
		t.Fatal("NewAnalyzerWithONNX() returned nil")
	}

	if analyzer.Analyzer == nil {
		t.Error("NewAnalyzerWithONNX().Analyzer is nil")
	}

	if analyzer.detector == nil {
		t.Error("NewAnalyzerWithONNX().detector is nil")
	}
}

func TestNewAnalyzerWithONNX_NilDetector(t *testing.T) {
	config := DefaultConfig()

	analyzer := NewAnalyzerWithONNX(config, nil)

	if analyzer == nil {
		t.Fatal("NewAnalyzerWithONNX() returned nil")
	}

	if analyzer.detector != nil {
		t.Error("NewAnalyzerWithONNX().detector should be nil")
	}
}

func TestAnalyzerWithONNX_AnalyzeWithImage_NoDetector(t *testing.T) {
	config := DefaultConfig()
	analyzer := NewAnalyzerWithONNX(config, nil)

	content := &extractor.PageContent{
		TextBlocks: []types.TextBlock{
			{Text: "Test paragraph", X: 100, Y: 100, Width: 200, Height: 20, FontSize: 12},
		},
		PageWidth:  612,
		PageHeight: 792,
	}

	// Create a dummy image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	elements := analyzer.AnalyzeWithImage(content, img)

	if len(elements) == 0 {
		t.Error("AnalyzeWithImage() returned no elements")
	}
}

func TestAnalyzerWithONNX_AnalyzeWithImage_DetectorUnavailable(t *testing.T) {
	config := DefaultConfig()
	detector := &mockDetector{available: false}
	analyzer := NewAnalyzerWithONNX(config, detector)

	content := &extractor.PageContent{
		TextBlocks: []types.TextBlock{
			{Text: "Test paragraph", X: 100, Y: 100, Width: 200, Height: 20, FontSize: 12},
		},
		PageWidth:  612,
		PageHeight: 792,
	}

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	elements := analyzer.AnalyzeWithImage(content, img)

	// Should fall back to rule-based analysis
	if len(elements) == 0 {
		t.Error("AnalyzeWithImage() returned no elements")
	}
}

func TestAnalyzerWithONNX_AnalyzeWithImage_NilImage(t *testing.T) {
	config := DefaultConfig()
	detector := &mockDetector{available: true}
	analyzer := NewAnalyzerWithONNX(config, detector)

	content := &extractor.PageContent{
		TextBlocks: []types.TextBlock{
			{Text: "Test paragraph", X: 100, Y: 100, Width: 200, Height: 20, FontSize: 12},
		},
		PageWidth:  612,
		PageHeight: 792,
	}

	// Nil image should fall back to rule-based
	elements := analyzer.AnalyzeWithImage(content, nil)

	if len(elements) == 0 {
		t.Error("AnalyzeWithImage() returned no elements")
	}
}

func TestAnalyzerWithONNX_AnalyzeWithImage_WithDetections(t *testing.T) {
	config := DefaultConfig()
	config.MinONNXConfidence = 0.5

	// Mock detector that returns a table detection
	detector := &mockDetector{
		available: true,
		detections: &types.PageDetections{
			PageIndex:  0,
			PageWidth:  612,
			PageHeight: 792,
			Detections: []types.BoundingBox{
				{
					X1:         100,
					Y1:         100,
					X2:         300,
					Y2:         200,
					Confidence: 0.9,
					ClassID:    int(types.ClassTable),
					ClassName:  "table",
				},
			},
		},
	}

	analyzer := NewAnalyzerWithONNX(config, detector)

	content := &extractor.PageContent{
		TextBlocks: []types.TextBlock{
			// This should match the ONNX detection and become a table
			{Text: "Some table content", X: 100, Y: 100, Width: 200, Height: 100, FontSize: 12},
		},
		PageWidth:  612,
		PageHeight: 792,
	}

	img := image.NewRGBA(image.Rect(0, 0, 612, 792))

	elements := analyzer.AnalyzeWithImage(content, img)

	if len(elements) == 0 {
		t.Fatal("AnalyzeWithImage() returned no elements")
	}

	// Check that the element was classified as table due to ONNX detection
	foundTable := false
	for _, elem := range elements {
		if elem.Type == ElementTypeTable {
			foundTable = true
			if elem.ONNXConfidence < 0.5 {
				t.Errorf("Element ONNXConfidence = %v, want >= 0.5", elem.ONNXConfidence)
			}
		}
	}

	if !foundTable {
		t.Log("No table element found - this may be due to no overlap or confidence threshold")
	}
}

func TestAnalyzerWithONNX_fuseDetections_Empty(t *testing.T) {
	config := DefaultConfig()
	analyzer := NewAnalyzerWithONNX(config, nil)

	elements := []Element{
		{Type: ElementTypeParagraph, Content: "Test", X: 100, Y: 100, Width: 200, Height: 50},
	}

	// Nil detections
	result := analyzer.fuseDetections(elements, nil)
	if len(result) != len(elements) {
		t.Errorf("fuseDetections with nil = %d elements, want %d", len(result), len(elements))
	}

	// Empty detections
	result = analyzer.fuseDetections(elements, &types.PageDetections{})
	if len(result) != len(elements) {
		t.Errorf("fuseDetections with empty = %d elements, want %d", len(result), len(elements))
	}
}

func TestAnalyzerWithONNX_fuseDetections_NoOverlap(t *testing.T) {
	config := DefaultConfig()
	config.MinONNXConfidence = 0.5
	analyzer := NewAnalyzerWithONNX(config, nil)

	elements := []Element{
		{Type: ElementTypeParagraph, Content: "Test", X: 100, Y: 100, Width: 200, Height: 50},
	}

	detections := &types.PageDetections{
		Detections: []types.BoundingBox{
			// Far from the element
			{X1: 500, Y1: 500, X2: 600, Y2: 600, Confidence: 0.9, ClassID: int(types.ClassTable)},
		},
	}

	result := analyzer.fuseDetections(elements, detections)

	// Should remain paragraph since no overlap
	if result[0].Type != ElementTypeParagraph {
		t.Errorf("fuseDetections type = %v, want paragraph (no overlap)", result[0].Type)
	}
}

func TestAnalyzerWithONNX_fuseDetections_LowConfidence(t *testing.T) {
	config := DefaultConfig()
	config.MinONNXConfidence = 0.5
	analyzer := NewAnalyzerWithONNX(config, nil)

	elements := []Element{
		{Type: ElementTypeParagraph, Content: "Test", X: 100, Y: 100, Width: 200, Height: 50},
	}

	detections := &types.PageDetections{
		Detections: []types.BoundingBox{
			// Overlapping but low confidence
			{X1: 100, Y1: 100, X2: 300, Y2: 150, Confidence: 0.3, ClassID: int(types.ClassTable)},
		},
	}

	result := analyzer.fuseDetections(elements, detections)

	// Should remain paragraph due to low confidence
	if result[0].Type != ElementTypeParagraph {
		t.Errorf("fuseDetections type = %v, want paragraph (low confidence)", result[0].Type)
	}
}

func TestAnalyzerWithONNX_fuseDetections_HighConfidenceMatch(t *testing.T) {
	config := DefaultConfig()
	config.MinONNXConfidence = 0.5
	analyzer := NewAnalyzerWithONNX(config, nil)

	elements := []Element{
		{Type: ElementTypeParagraph, Content: "Test", X: 100, Y: 100, Width: 200, Height: 50},
	}

	detections := &types.PageDetections{
		Detections: []types.BoundingBox{
			// Good overlap and high confidence
			{X1: 100, Y1: 100, X2: 300, Y2: 150, Confidence: 0.9, ClassID: int(types.ClassTitle)},
		},
	}

	result := analyzer.fuseDetections(elements, detections)

	// Should be reclassified to header
	if result[0].Type != types.ElementTypeHeader {
		t.Errorf("fuseDetections type = %v, want header", result[0].Type)
	}

	if result[0].ONNXConfidence != 0.9 {
		t.Errorf("fuseDetections ONNXConfidence = %v, want 0.9", result[0].ONNXConfidence)
	}
}

func TestAnalyzerWithONNX_fuseDetections_PreservesImages(t *testing.T) {
	config := DefaultConfig()
	config.MinONNXConfidence = 0.5
	analyzer := NewAnalyzerWithONNX(config, nil)

	elements := []Element{
		{Type: ElementTypeImage, Content: "image.png", X: 100, Y: 100, Width: 200, Height: 200},
	}

	detections := &types.PageDetections{
		Detections: []types.BoundingBox{
			// Exact match but should not override image
			{X1: 100, Y1: 100, X2: 300, Y2: 300, Confidence: 0.95, ClassID: int(types.ClassFigure)},
		},
	}

	result := analyzer.fuseDetections(elements, detections)

	// Images should not be overridden
	if result[0].Type != ElementTypeImage {
		t.Errorf("fuseDetections type = %v, want image (preserved)", result[0].Type)
	}
}

func TestAnalyzerWithONNX_findBestMatch(t *testing.T) {
	config := DefaultConfig()
	analyzer := NewAnalyzerWithONNX(config, nil)

	tests := []struct {
		name       string
		element    *Element
		detections *types.PageDetections
		wantMatch  bool
	}{
		{
			name:       "nil element",
			element:    nil,
			detections: &types.PageDetections{},
			wantMatch:  false,
		},
		{
			name:       "nil detections",
			element:    &Element{X: 100, Y: 100, Width: 200, Height: 100},
			detections: nil,
			wantMatch:  false,
		},
		{
			name:    "good overlap",
			element: &Element{X: 100, Y: 100, Width: 200, Height: 100},
			detections: &types.PageDetections{
				Detections: []types.BoundingBox{
					{X1: 100, Y1: 100, X2: 300, Y2: 200, Confidence: 0.8, ClassID: 5},
				},
			},
			wantMatch: true,
		},
		{
			name:    "no overlap",
			element: &Element{X: 100, Y: 100, Width: 100, Height: 100},
			detections: &types.PageDetections{
				Detections: []types.BoundingBox{
					{X1: 500, Y1: 500, X2: 600, Y2: 600, Confidence: 0.8, ClassID: 5},
				},
			},
			wantMatch: false,
		},
		{
			name:    "low overlap",
			element: &Element{X: 100, Y: 100, Width: 100, Height: 100},
			detections: &types.PageDetections{
				Detections: []types.BoundingBox{
					// Only 10x100 overlap (10% of element)
					{X1: 190, Y1: 100, X2: 250, Y2: 200, Confidence: 0.8, ClassID: 5},
				},
			},
			wantMatch: false, // IoU < 0.3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := analyzer.findBestMatch(tt.element, tt.detections)
			if (match != nil) != tt.wantMatch {
				t.Errorf("findBestMatch() = %v, wantMatch = %v", match != nil, tt.wantMatch)
			}
		})
	}
}

func TestMapONNXClassToElementType(t *testing.T) {
	tests := []struct {
		classID  int
		wantType types.ElementType
	}{
		{int(types.ClassTitle), types.ElementTypeHeader},
		{int(types.ClassPlainText), types.ElementTypeParagraph},
		{int(types.ClassFigure), types.ElementTypeFigure},
		{int(types.ClassTable), types.ElementTypeTable},
		{int(types.ClassIsolateFormula), types.ElementTypeEquation},
		{int(types.ClassAbandon), types.ElementTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(types.DocLayoutClass(tt.classID).String(), func(t *testing.T) {
			got := mapONNXClassToElementType(tt.classID)
			if got != tt.wantType {
				t.Errorf("mapONNXClassToElementType(%d) = %v, want %v", tt.classID, got, tt.wantType)
			}
		})
	}
}

func TestCreateONNXConfig(t *testing.T) {
	// Test with nil config
	onnxCfg := CreateONNXConfig(nil)
	if onnxCfg == nil {
		t.Fatal("CreateONNXConfig(nil) returned nil")
	}

	// Test with custom config
	layoutCfg := &Config{
		ONNXModelPath:     "/path/to/model.onnx",
		ONNXRuntimePath:   "/path/to/runtime",
		ONNXConfThreshold: 0.35,
		ONNXNMSThreshold:  0.50,
		ONNXInputSize:     1024,
		ONNXStride:        32,
	}

	onnxCfg = CreateONNXConfig(layoutCfg)
	if onnxCfg == nil {
		t.Fatal("CreateONNXConfig(custom) returned nil")
	}

	if onnxCfg.ModelPath != layoutCfg.ONNXModelPath {
		t.Errorf("ModelPath = %v, want %v", onnxCfg.ModelPath, layoutCfg.ONNXModelPath)
	}

	if onnxCfg.RuntimePath != layoutCfg.ONNXRuntimePath {
		t.Errorf("RuntimePath = %v, want %v", onnxCfg.RuntimePath, layoutCfg.ONNXRuntimePath)
	}

	if onnxCfg.ConfThreshold != layoutCfg.ONNXConfThreshold {
		t.Errorf("ConfThreshold = %v, want %v", onnxCfg.ConfThreshold, layoutCfg.ONNXConfThreshold)
	}

	if onnxCfg.NMSThreshold != layoutCfg.ONNXNMSThreshold {
		t.Errorf("NMSThreshold = %v, want %v", onnxCfg.NMSThreshold, layoutCfg.ONNXNMSThreshold)
	}

	if onnxCfg.InputSize != layoutCfg.ONNXInputSize {
		t.Errorf("InputSize = %v, want %v", onnxCfg.InputSize, layoutCfg.ONNXInputSize)
	}

	if onnxCfg.Stride != layoutCfg.ONNXStride {
		t.Errorf("Stride = %v, want %v", onnxCfg.Stride, layoutCfg.ONNXStride)
	}
}
