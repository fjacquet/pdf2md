// Package ocr provides ONNX-based optical character recognition for scanned
// or image-only PDF pages. It is used as a fallback when the PDF content
// stream contains no extractable text.
//
// The pipeline mirrors internal/onnx/: pure preprocessing/postprocessing
// helpers, stateful Recognizer struct that owns ONNX sessions, dependency
// injection via a Deps struct for testability.
package ocr

import (
	"errors"
	"image"

	"github.com/fjacquet/pdf2md/internal/types"
)

// Recognizer performs text detection and recognition on a rendered page image.
// Implementations own the ONNX sessions and character dictionary.
type Recognizer interface {
	// RecognizePage detects and recognizes text in the given image.
	// pageWidth and pageHeight are in PDF points; returned TextBlocks use
	// the PDF coordinate system (origin bottom-left, Y increases upward).
	RecognizePage(pageImage image.Image, pageWidth, pageHeight float64) ([]types.TextBlock, error)

	// Close releases ONNX sessions and other resources.
	Close() error

	// IsAvailable reports whether the recognizer is ready to use.
	// Returns false when model files or runtime could not be resolved.
	IsAvailable() bool
}

// Config holds recognizer configuration. All fields have sensible defaults
// from DefaultConfig; override only what the caller needs.
type Config struct {
	// DetModelPath is the path to the text-detection ONNX model
	// (PP-OCRv5 mobile det). Auto-downloaded if empty.
	DetModelPath string

	// RecModelPath is the path to the text-recognition ONNX model
	// (PP-OCRv5 mobile rec). Auto-downloaded if empty.
	RecModelPath string

	// DictPath is the path to the character dictionary (one char per line).
	// If empty, the embedded Latin dictionary is used.
	DictPath string

	// RuntimePath is the path to the ONNX Runtime shared library.
	// Reuses the same library as the layout detector; if empty, standard
	// system paths are searched.
	RuntimePath string

	// DetInputSize is the short-side target for detection input in pixels.
	// DB-style models are fully convolutional; we resize so the short side
	// is DetInputSize and the long side is aligned to DetStride.
	// Default: 960.
	DetInputSize int

	// DetStride rounds the detection input dimensions. Default: 32.
	DetStride int

	// DetDBThreshold is the binarization threshold on the probability map.
	// Default: 0.3.
	DetDBThreshold float64

	// DetBoxThreshold is the minimum mean probability inside a box for it
	// to be kept. Default: 0.5.
	DetBoxThreshold float64

	// DetMinBoxSize discards detected regions smaller than this many pixels
	// on the short side. Default: 3.
	DetMinBoxSize int

	// RecInputHeight is the fixed input height for the recognition model.
	// Crops are resized to this height keeping aspect ratio, then padded to
	// RecInputMaxWidth. Default: 48 (PP-OCRv5).
	RecInputHeight int

	// RecInputMaxWidth is the maximum input width for the recognition model.
	// Longer crops are split; shorter are right-padded. Default: 320.
	RecInputMaxWidth int

	// RecMinConfidence drops recognized text below this mean confidence.
	// Default: 0.5.
	RecMinConfidence float64

	// UseCoreML enables the CoreML execution provider on macOS.
	// Default: false (PaddleOCR ops have historical CoreML compatibility issues).
	UseCoreML bool
}

// DefaultConfig returns a Config with production-ready defaults for
// PP-OCRv5 mobile (det + rec).
func DefaultConfig() *Config {
	return &Config{
		DetInputSize:     960,
		DetStride:        32,
		DetDBThreshold:   0.3,
		DetBoxThreshold:  0.5,
		DetMinBoxSize:    3,
		RecInputHeight:   48,
		RecInputMaxWidth: 320,
		RecMinConfidence: 0.5,
		UseCoreML:        false,
	}
}

// PathResolver resolves (and potentially downloads) a model or dictionary
// path from a user-supplied hint (which may be empty, a filename, or a
// full path). Returned path must exist on disk.
type PathResolver func(configPath string) (string, error)

// DictLoader reads a character dictionary from the given path or returns
// the embedded default when path is empty. One character per line.
type DictLoader func(path string) ([]string, error)

// Deps holds injectable I/O dependencies. Pure logic does not go here —
// only things that touch the filesystem, network, or ONNX runtime.
// Following the same pattern as internal/onnx.DetectorDeps.
type Deps struct {
	ResolveDetModel PathResolver
	ResolveRecModel PathResolver
	ResolveRuntime  PathResolver
	LoadDict        DictLoader
}

// DefaultDeps returns production implementations that auto-download models
// and load the embedded Latin dictionary.
func DefaultDeps() *Deps {
	return &Deps{
		ResolveDetModel: ResolveDetModelPath,
		ResolveRecModel: ResolveRecModelPath,
		ResolveRuntime:  ResolveRuntimePath,
		LoadDict:        LoadDict,
	}
}

// Sentinel errors.
var (
	ErrRecognizerNotAvailable = errors.New("ocr: recognizer not available")
	ErrModelNotFound          = errors.New("ocr: model file not found")
	ErrModelDownloadFailed    = errors.New("ocr: model download failed")
	ErrRuntimeNotFound        = errors.New("ocr: ONNX runtime not found")
	ErrInvalidImage           = errors.New("ocr: invalid page image")
	ErrDictEmpty              = errors.New("ocr: dictionary is empty")
)
