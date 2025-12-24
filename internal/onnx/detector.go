// Package onnx provides ONNX-based document layout detection using DocLayout-YOLO model.
package onnx

import (
	"fmt"
	"image"
	"log/slog"
	"runtime"
	"sync"

	"github.com/fjacquet/pdf2md/internal/types"
	ort "github.com/yalue/onnxruntime_go"
)

// LayoutDetector performs document layout detection using ML models.
type LayoutDetector interface {
	// DetectLayout analyzes a page image and returns detected regions.
	// pageWidth and pageHeight are in PDF points.
	DetectLayout(pageImage image.Image, pageWidth, pageHeight float64) (*types.PageDetections, error)

	// Close releases resources held by the detector.
	Close() error

	// IsAvailable returns true if the detector is ready to use.
	IsAvailable() bool
}

// Config holds ONNX detector configuration.
type Config struct {
	// ModelPath is the path to the ONNX model file.
	// If empty, the model will be auto-downloaded to ~/.pdf2md/models/
	ModelPath string

	// RuntimePath is the path to the ONNX Runtime shared library.
	// If empty, standard system paths will be searched.
	RuntimePath string

	// ConfThreshold is the minimum confidence for a detection to be included.
	// Default: 0.25
	ConfThreshold float64

	// NMSThreshold is the IoU threshold for Non-Maximum Suppression.
	// Higher values allow more overlapping boxes.
	// Default: 0.45
	NMSThreshold float64

	// InputSize is the model input size (width and height).
	// DocLayout-YOLO uses 1024x1024.
	// Default: 1024
	InputSize int

	// Stride is the model stride for padding alignment.
	// Default: 32
	Stride int

	// UseCoreML enables CoreML execution provider on macOS.
	// This leverages Apple Neural Engine / GPU for faster inference on M1/M2/M3 Macs.
	// Falls back to CPU if CoreML is unavailable.
	// Default: true (on macOS)
	UseCoreML bool
}

// DefaultConfig returns a Config with sensible defaults for DocLayout-YOLO.
func DefaultConfig() *Config {
	return &Config{
		ConfThreshold: 0.25,
		NMSThreshold:  0.45,
		InputSize:     1024,
		Stride:        32,
		UseCoreML:     runtime.GOOS == "darwin", // Enable CoreML by default on macOS
	}
}

// ModelLoader is a function type for loading model bytes from a path.
// This allows dependency injection for testing.
type ModelLoader func(path string) ([]byte, error)

// RuntimeInitializer is a function type for initializing the ONNX runtime.
// This allows dependency injection for testing.
type RuntimeInitializer func(libraryPath string) error

// DetectorDeps holds injectable dependencies for the Detector.
// Following functional programming principles, all I/O dependencies are injected.
type DetectorDeps struct {
	LoadModel      ModelLoader
	InitRuntime    RuntimeInitializer
	ResolveModel   func(configPath string) (string, error)
	ResolveRuntime func(configPath string) (string, error)
}

var (
	ortInitOnce sync.Once
	ortInitErr  error
)

// Detector implements LayoutDetector using ONNX Runtime with DocLayout-YOLO model.
type Detector struct {
	config    *Config
	deps      *DetectorDeps
	available bool
	modelPath string
	session   *ort.DynamicAdvancedSession
}

// NewDetector creates a new LayoutDetector with the given configuration and dependencies.
// If deps is nil, default implementations will be used.
func NewDetector(config *Config, deps *DetectorDeps) (*Detector, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if deps == nil {
		deps = defaultDeps()
	}

	d := &Detector{
		config:    config,
		deps:      deps,
		available: false,
	}

	// Initialize will be called to set up the ONNX runtime
	if err := d.initialize(); err != nil {
		return nil, err
	}

	return d, nil
}

// DetectLayout implements LayoutDetector.
func (d *Detector) DetectLayout(pageImage image.Image, pageWidth, pageHeight float64) (*types.PageDetections, error) {
	if !d.available {
		return nil, ErrDetectorNotAvailable
	}

	if pageImage == nil {
		return nil, ErrDetectorNotAvailable
	}

	// 1. Preprocess image for YOLO
	tensor, scale, padX, padY := Preprocess(pageImage, d.config.InputSize, d.config.Stride)

	// 2. Run inference
	rawDetections, err := d.runInference(tensor)
	if err != nil {
		return nil, err
	}

	// 3. Postprocess: filter by confidence, apply NMS, scale to page coordinates
	detections := Postprocess(
		rawDetections,
		d.config.ConfThreshold,
		d.config.NMSThreshold,
		scale,
		padX,
		padY,
		pageWidth,
		pageHeight,
	)

	return &types.PageDetections{
		PageWidth:  pageWidth,
		PageHeight: pageHeight,
		Detections: detections,
	}, nil
}

// Close implements LayoutDetector.
func (d *Detector) Close() error {
	d.available = false
	if d.session != nil {
		d.session.Destroy()
		d.session = nil
	}
	return nil
}

// IsAvailable implements LayoutDetector.
func (d *Detector) IsAvailable() bool {
	return d.available
}

// initialize sets up the ONNX runtime and loads the model.
func (d *Detector) initialize() error {
	logger := slog.Default()

	// 1. Resolve runtime library path
	runtimePath, err := d.deps.ResolveRuntime(d.config.RuntimePath)
	if err != nil {
		logger.Debug("ONNX runtime not found", "error", err)
		return nil // Non-fatal: detector simply won't be available
	}
	logger.Debug("ONNX runtime found", "path", runtimePath)

	// 2. Initialize ONNX Runtime (once globally)
	ortInitOnce.Do(func() {
		ort.SetSharedLibraryPath(runtimePath)
		ortInitErr = ort.InitializeEnvironment()
	})
	if ortInitErr != nil {
		logger.Debug("ONNX environment init failed", "error", ortInitErr)
		return nil // Non-fatal
	}
	logger.Debug("ONNX environment initialized")

	// 3. Resolve model path (may trigger download)
	modelPath, err := d.deps.ResolveModel(d.config.ModelPath)
	if err != nil {
		logger.Debug("ONNX model not found", "error", err)
		return nil // Non-fatal
	}
	d.modelPath = modelPath
	logger.Debug("ONNX model found", "path", modelPath)

	// 4. Create session options (with hardware acceleration if available)
	sessionOptions, err := d.createSessionOptions(logger)
	if err != nil {
		logger.Debug("Failed to create session options", "error", err)
		return nil // Non-fatal
	}
	defer sessionOptions.Destroy()

	// 5. Create ONNX session
	// DocLayout-YOLO input: "images" [1, 3, 1024, 1024]
	// DocLayout-YOLO output: "output0" [1, 15, 8400] (YOLOv8 format)
	session, err := ort.NewDynamicAdvancedSession(
		modelPath,
		[]string{"images"},
		[]string{"output0"},
		sessionOptions,
	)
	if err != nil {
		logger.Debug("ONNX session creation failed", "error", err)
		return nil // Non-fatal: session creation failed
	}

	d.session = session
	d.available = true
	logger.Debug("ONNX session created successfully")
	return nil
}

// createSessionOptions creates ONNX session options with hardware acceleration.
// On macOS with UseCoreML=true, enables CoreML execution provider for Apple Neural Engine/GPU.
// Falls back to CPU execution if hardware acceleration is unavailable.
func (d *Detector) createSessionOptions(logger *slog.Logger) (*ort.SessionOptions, error) {
	options, err := ort.NewSessionOptions()
	if err != nil {
		return nil, fmt.Errorf("failed to create session options: %w", err)
	}

	logger.Debug("Creating session options", "useCoreML", d.config.UseCoreML, "os", runtime.GOOS)

	// Try to enable CoreML on macOS
	if d.config.UseCoreML && runtime.GOOS == "darwin" {
		// CoreML options for optimal performance
		// See: https://onnxruntime.ai/docs/execution-providers/CoreML-ExecutionProvider.html
		coreMLOptions := map[string]string{
			"CoreMLFlags": "0", // Default flags (can use COREML_FLAG_USE_CPU_AND_GPU = 1)
		}

		if err := options.AppendExecutionProviderCoreMLV2(coreMLOptions); err != nil {
			// CoreML not available, fall back to CPU (this is not fatal)
			logger.Debug("CoreML execution provider not available, using CPU", "error", err)
		} else {
			logger.Info("CoreML execution provider enabled (Apple Neural Engine/GPU)")
		}
	}

	return options, nil
}

// runInference executes the ONNX model on the preprocessed tensor.
// Returns raw detections in format [N, 6] = [x1, y1, x2, y2, conf, class].
func (d *Detector) runInference(tensor []float32) ([][]float32, error) {
	if d.session == nil {
		return nil, ErrDetectorNotAvailable
	}

	// Create input tensor [1, 3, 1024, 1024]
	inputSize := d.config.InputSize
	inputShape := ort.NewShape(1, 3, int64(inputSize), int64(inputSize))
	inputTensor, err := ort.NewTensor(inputShape, tensor)
	if err != nil {
		return nil, fmt.Errorf("failed to create input tensor: %w", err)
	}
	defer inputTensor.Destroy()

	// Run inference with auto-allocated output
	outputs := []ort.Value{nil}
	err = d.session.Run([]ort.Value{inputTensor}, outputs)
	if err != nil {
		return nil, fmt.Errorf("inference failed: %w", err)
	}

	// Get output tensor
	if outputs[0] == nil {
		return nil, fmt.Errorf("no output from model")
	}
	defer outputs[0].Destroy()

	// Cast to float32 tensor and get data
	outputTensor, ok := outputs[0].(*ort.Tensor[float32])
	if !ok {
		return nil, fmt.Errorf("unexpected output tensor type")
	}

	outputData := outputTensor.GetData()
	outputShape := outputTensor.GetShape()

	// Parse YOLOv8 output format: [1, 15, 8400]
	// 15 = 4 (bbox) + 11 (class scores for DocLayout)
	// 8400 = number of predictions
	return d.parseYOLOv8Output(outputData, outputShape)
}

// parseYOLOv8Output converts YOLOv8 output to detection format.
// Input shape: [1, 15, 8400] where 15 = 4 (x,y,w,h) + 11 classes
// Output: slice of [x1, y1, x2, y2, confidence, class_id]
func (d *Detector) parseYOLOv8Output(data []float32, shape ort.Shape) ([][]float32, error) {
	if len(shape) != 3 {
		return nil, fmt.Errorf("unexpected output shape: %v", shape)
	}

	numFeatures := int(shape[1]) // 15 (4 bbox + 11 classes)
	numPreds := int(shape[2])    // 8400

	if numFeatures < 5 {
		return nil, fmt.Errorf("insufficient features: %d", numFeatures)
	}

	numClasses := numFeatures - 4
	var detections [][]float32

	for i := 0; i < numPreds; i++ {
		// Extract bbox (center x, center y, width, height)
		cx := data[0*numPreds+i]
		cy := data[1*numPreds+i]
		w := data[2*numPreds+i]
		h := data[3*numPreds+i]

		// Find best class and its score
		bestClass := 0
		bestScore := float32(0)
		for c := 0; c < numClasses; c++ {
			score := data[(4+c)*numPreds+i]
			if score > bestScore {
				bestScore = score
				bestClass = c
			}
		}

		// Convert center format to corner format
		x1 := cx - w/2
		y1 := cy - h/2
		x2 := cx + w/2
		y2 := cy + h/2

		detections = append(detections, []float32{
			x1, y1, x2, y2, bestScore, float32(bestClass),
		})
	}

	return detections, nil
}

// defaultDeps returns default implementations for DetectorDeps.
func defaultDeps() *DetectorDeps {
	return &DetectorDeps{
		LoadModel:      DefaultModelLoader,
		InitRuntime:    DefaultRuntimeInitializer,
		ResolveModel:   ResolveModelPath,
		ResolveRuntime: ResolveRuntimePath,
	}
}
