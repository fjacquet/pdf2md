package ocr

import (
	"fmt"
	"image"
	"log/slog"
	"math"
	"runtime"

	"github.com/fjacquet/pdf2md/internal/types"
	ort "github.com/yalue/onnxruntime_go"
)

// PaddleRecognizer implements Recognizer using PaddleOCR PP-OCRv5 mobile
// det+rec models via ONNX Runtime. Shares the runtime library with
// internal/onnx.
type PaddleRecognizer struct {
	config    *Config
	deps      *Deps
	available bool
	chars     []string
	detSess   *ort.DynamicAdvancedSession
	recSess   *ort.DynamicAdvancedSession
	detInput  string
	detOutput string
	recInput  string
	recOutput string
}

// NewPaddleRecognizer constructs and initializes a PaddleRecognizer.
// If any model/runtime is missing, returns a non-available recognizer (no error):
// callers check IsAvailable() and skip the OCR step gracefully.
func NewPaddleRecognizer(config *Config, deps *Deps) (*PaddleRecognizer, error) {
	if config == nil {
		config = DefaultConfig()
	}
	if deps == nil {
		deps = DefaultDeps()
	}
	r := &PaddleRecognizer{config: config, deps: deps}
	if err := r.initialize(); err != nil {
		return nil, err
	}
	return r, nil
}

// RecognizePage implements Recognizer.
// pageWidth/pageHeight are PDF points; returned TextBlocks use PDF coordinates
// (origin bottom-left, Y increases upward).
func (r *PaddleRecognizer) RecognizePage(pageImage image.Image, pageWidth, pageHeight float64) ([]types.TextBlock, error) {
	if !r.available {
		return nil, ErrRecognizerNotAvailable
	}
	if pageImage == nil {
		return nil, ErrInvalidImage
	}

	boxes, err := r.detectBoxes(pageImage)
	if err != nil {
		return nil, err
	}
	if len(boxes) == 0 {
		return nil, nil
	}

	imgBounds := pageImage.Bounds()
	imgW, imgH := imgBounds.Dx(), imgBounds.Dy()

	blocks := make([]types.TextBlock, 0, len(boxes))
	for _, box := range boxes {
		text, conf, err := r.recognizeBox(pageImage, box)
		if err != nil || text == "" || conf < r.config.RecMinConfidence {
			continue
		}
		blocks = append(blocks, boxToTextBlock(box, text, imgW, imgH, pageWidth, pageHeight))
	}
	return blocks, nil
}

// Close implements Recognizer.
func (r *PaddleRecognizer) Close() error {
	r.available = false
	if r.detSess != nil {
		r.detSess.Destroy()
		r.detSess = nil
	}
	if r.recSess != nil {
		r.recSess.Destroy()
		r.recSess = nil
	}
	return nil
}

// IsAvailable implements Recognizer.
func (r *PaddleRecognizer) IsAvailable() bool { return r.available }

// initialize resolves paths, boots ONNX runtime, loads dict, and opens sessions.
// Missing resources are non-fatal: recognizer simply stays unavailable.
func (r *PaddleRecognizer) initialize() error {
	logger := slog.Default()

	runtimePath, err := r.deps.ResolveRuntime(r.config.RuntimePath)
	if err != nil {
		logger.Debug("ocr: ONNX runtime not found", "error", err)
		return nil
	}
	if !ort.IsInitialized() {
		ort.SetSharedLibraryPath(runtimePath)
		if err := ort.InitializeEnvironment(); err != nil {
			logger.Debug("ocr: ONNX init failed", "error", err)
			return nil
		}
	}

	detPath, err := r.deps.ResolveDetModel(r.config.DetModelPath)
	if err != nil {
		logger.Debug("ocr: det model not found", "error", err)
		return nil
	}
	recPath, err := r.deps.ResolveRecModel(r.config.RecModelPath)
	if err != nil {
		logger.Debug("ocr: rec model not found", "error", err)
		return nil
	}

	chars, err := r.deps.LoadDict(r.config.DictPath)
	if err != nil {
		logger.Debug("ocr: dict load failed", "error", err)
		return nil
	}
	// CTC blank is implicitly class 0; chars index i maps to class i+1.
	r.chars = chars

	detSess, detIn, detOut, err := openSession(detPath, r.sessionOptions(logger))
	if err != nil {
		logger.Debug("ocr: det session failed", "error", err)
		return nil
	}
	recSess, recIn, recOut, err := openSession(recPath, r.sessionOptions(logger))
	if err != nil {
		detSess.Destroy()
		logger.Debug("ocr: rec session failed", "error", err)
		return nil
	}

	r.detSess, r.detInput, r.detOutput = detSess, detIn, detOut
	r.recSess, r.recInput, r.recOutput = recSess, recIn, recOut
	r.available = true
	logger.Debug("ocr: PaddleRecognizer ready", "det", detPath, "rec", recPath)
	return nil
}

// sessionOptions builds per-session options. CoreML is opt-in because
// PaddleOCR ops historically have CoreML compatibility issues.
func (r *PaddleRecognizer) sessionOptions(logger *slog.Logger) *ort.SessionOptions {
	opts, err := ort.NewSessionOptions()
	if err != nil {
		return nil
	}
	if r.config.UseCoreML && runtime.GOOS == "darwin" {
		if err := opts.AppendExecutionProviderCoreMLV2(map[string]string{}); err != nil {
			logger.Debug("ocr: CoreML unavailable, using CPU", "error", err)
		} else {
			logger.Info("ocr: CoreML enabled")
		}
	}
	return opts
}

// openSession introspects the ONNX model to discover input/output names,
// then opens a dynamic session. Works across different PaddleOCR export
// scripts that use different tensor names.
func openSession(path string, opts *ort.SessionOptions) (*ort.DynamicAdvancedSession, string, string, error) {
	inputs, outputs, err := ort.GetInputOutputInfo(path)
	if err != nil {
		return nil, "", "", fmt.Errorf("introspect %s: %w", path, err)
	}
	if len(inputs) == 0 || len(outputs) == 0 {
		return nil, "", "", fmt.Errorf("model has no inputs/outputs: %s", path)
	}
	in, out := inputs[0].Name, outputs[0].Name
	sess, err := ort.NewDynamicAdvancedSession(path, []string{in}, []string{out}, opts)
	if err != nil {
		return nil, "", "", fmt.Errorf("session for %s: %w", path, err)
	}
	return sess, in, out, nil
}

// detectBoxes runs the DB detector and returns boxes in original-image pixels.
func (r *PaddleRecognizer) detectBoxes(pageImage image.Image) ([]DetectionBox, error) {
	tensor, resizedW, resizedH, scaleX, scaleY := PreprocessDet(
		pageImage, r.config.DetInputSize, r.config.DetStride,
	)
	shape := ort.NewShape(1, 3, int64(resizedH), int64(resizedW))
	probMap, outShape, err := runSession(r.detSess, tensor, shape)
	if err != nil {
		return nil, fmt.Errorf("det inference: %w", err)
	}
	// Expected output shape: [1, 1, H, W]. Locate H and W robustly.
	h, w := findProbMapDims(outShape, resizedH, resizedW)
	if h == 0 || w == 0 || len(probMap) < h*w {
		return nil, fmt.Errorf("unexpected det output shape: %v", outShape)
	}
	// The DB output is at the *input* resolution; boxes land in resized-image
	// space, so the inverse scale back to original pixels uses scaleX/scaleY.
	return PostprocessDet(
		probMap[:h*w], w, h,
		scaleX, scaleY,
		r.config.DetDBThreshold, r.config.DetBoxThreshold, r.config.DetMinBoxSize,
	), nil
}

// findProbMapDims extracts the (H, W) spatial dims from an output shape,
// falling back to the preprocess resizedH/resizedW if the shape is degenerate.
func findProbMapDims(shape ort.Shape, fallbackH, fallbackW int) (int, int) {
	// Walk from the tail — the last two dims of a [B, C, H, W] or [B, H, W]
	// tensor are always H and W.
	if len(shape) >= 2 {
		return int(shape[len(shape)-2]), int(shape[len(shape)-1])
	}
	return fallbackH, fallbackW
}

// recognizeBox crops, runs rec inference, and CTC-decodes.
func (r *PaddleRecognizer) recognizeBox(pageImage image.Image, box DetectionBox) (string, float64, error) {
	tensor, tensorW := PreprocessRec(
		pageImage, box, r.config.RecInputHeight, r.config.RecInputMaxWidth,
	)
	if len(tensor) == 0 {
		return "", 0, nil
	}
	shape := ort.NewShape(1, 3, int64(r.config.RecInputHeight), int64(tensorW))
	logits, outShape, err := runSession(r.recSess, tensor, shape)
	if err != nil {
		return "", 0, fmt.Errorf("rec inference: %w", err)
	}
	// Recognition output: [1, T, C].
	t, c := findRecDims(outShape)
	if t == 0 || c == 0 || len(logits) < t*c {
		return "", 0, nil
	}
	text, conf := CTCDecode(logits[:t*c], t, c, r.chars)
	return text, conf, nil
}

// findRecDims pulls (T, C) from a [1, T, C] or [T, C] shape.
func findRecDims(shape ort.Shape) (int, int) {
	switch len(shape) {
	case 3:
		return int(shape[1]), int(shape[2])
	case 2:
		return int(shape[0]), int(shape[1])
	default:
		return 0, 0
	}
}

// runSession is a small helper that wraps tensor lifecycle and extracts
// float32 output + shape.
func runSession(sess *ort.DynamicAdvancedSession, tensor []float32, shape ort.Shape) ([]float32, ort.Shape, error) {
	if sess == nil {
		return nil, nil, ErrRecognizerNotAvailable
	}
	input, err := ort.NewTensor(shape, tensor)
	if err != nil {
		return nil, nil, fmt.Errorf("tensor: %w", err)
	}
	defer input.Destroy()

	outputs := []ort.Value{nil}
	if err := sess.Run([]ort.Value{input}, outputs); err != nil {
		return nil, nil, fmt.Errorf("run: %w", err)
	}
	if outputs[0] == nil {
		return nil, nil, fmt.Errorf("nil output")
	}
	defer outputs[0].Destroy()

	out, ok := outputs[0].(*ort.Tensor[float32])
	if !ok {
		return nil, nil, fmt.Errorf("unexpected output type")
	}
	return out.GetData(), out.GetShape(), nil
}

// boxToTextBlock converts an image-space detection box into a PDF-space TextBlock.
//
// Pure function. Image coords have origin top-left and Y-down; PDF coords
// have origin bottom-left and Y-up. FontSize is estimated from box height
// so downstream layout heuristics (header detection, etc.) behave as if the
// OCR-derived blocks came from the extractor.
func boxToTextBlock(box DetectionBox, text string, imgW, imgH int, pageWidth, pageHeight float64) types.TextBlock {
	sx := pageWidth / float64(imgW)
	sy := pageHeight / float64(imgH)
	x := float64(box.MinX) * sx
	w := float64(box.Width()) * sx
	h := float64(box.Height()) * sy
	// Flip Y: image MinY (top) maps to PDF Y at top of box; TextBlock Y is baseline.
	yTop := pageHeight - float64(box.MinY)*sy
	y := yTop - h
	return types.TextBlock{
		Text:     text,
		X:        x,
		Y:        y,
		Width:    w,
		Height:   h,
		FontSize: estimateFontSize(h),
		FontName: "OCR",
	}
}

// estimateFontSize maps box height (in PDF points) to an approximate font size.
// PP-OCRv5 boxes are tight around cap-height; PDF font size ≈ box height / 0.7
// (cap-height is ~70% of em-height for common fonts). Clamp to a reasonable
// range to keep outliers from tripping the layout analyzer.
func estimateFontSize(boxHeightPts float64) float64 {
	size := boxHeightPts / 0.7
	return math.Max(6, math.Min(72, size))
}
