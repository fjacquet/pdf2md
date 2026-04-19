package layout

// Config holds configurable thresholds for layout analysis.
// These values can be adjusted to fine-tune the analysis for different
// document types and formatting styles.
type Config struct {
	// Column detection
	ColumnGapThreshold float64 // Minimum gap to consider as column boundary (default: 10.0)

	// Header detection
	HeaderSizeRatio float64 // Minimum ratio to body font to detect as header (default: 1.2)

	// Header level thresholds (ratio of font size to body font)
	H1Ratio float64 // Default: 2.0
	H2Ratio float64 // Default: 1.75
	H3Ratio float64 // Default: 1.5
	H4Ratio float64 // Default: 1.3
	H5Ratio float64 // Default: 1.15
	H6Ratio float64 // Default: 1.0

	// Table detection
	TableMinColumns int     // Minimum columns to detect as table (default: 2)
	TableMinSpaces  int     // Minimum consecutive spaces for wide gap (default: 4)
	WideGapMultiple float64 // Multiplier for wide gap detection (default: 3.0)

	// Image filtering
	MinImageSize float64 // Minimum width/height to include images (default: 50.0)

	// Text merging
	LineHeightMultiplier  float64 // Multiplier for line height threshold (default: 1.6)
	IndentThreshold       float64 // Threshold to detect indentation (default: 5.0)
	WordGapMultiplier     float64 // Multiplier for word gap detection (default: 0.2)
	SentenceGapMultiplier float64 // Stricter gap for sentence boundaries (default: 1.5)

	// Equation detection
	EquationScoreThreshold int // Minimum score to classify as equation (default: 2)

	// ONNX Layout Detection
	// ONNX detection provides ML-based layout analysis using DocLayout-YOLO model.
	// When enabled, it overrides rule-based classification for detected regions.
	EnableONNX        bool    // Enable ONNX layout detection (default: true)
	ONNXModelPath     string  // Path to ONNX model file (auto-detected if empty)
	ONNXRuntimePath   string  // Path to ONNX Runtime library (auto-detected if empty)
	ONNXConfThreshold float64 // Minimum confidence for ONNX detections (default: 0.25)
	ONNXNMSThreshold  float64 // IoU threshold for Non-Maximum Suppression (default: 0.45)
	MinONNXConfidence float64 // Minimum confidence to override rule-based classification (default: 0.5)
	ONNXInputSize     int     // Model input size in pixels (default: 1024)
	ONNXStride        int     // Model stride for padding alignment (default: 32)
	ONNXUseCoreML     bool    // Enable CoreML acceleration on macOS (default: true on darwin)
}

// DefaultConfig returns the default layout configuration.
// These defaults work well for most PDF documents.
// ONNX detection is enabled by default for improved layout analysis.
func DefaultConfig() *Config {
	return &Config{
		// Column detection
		ColumnGapThreshold: 10.0,

		// Header detection
		HeaderSizeRatio: 1.2,

		// Header level thresholds
		H1Ratio: 2.0,
		H2Ratio: 1.75,
		H3Ratio: 1.5,
		H4Ratio: 1.3,
		H5Ratio: 1.15,
		H6Ratio: 1.0,

		// Table detection
		TableMinColumns: 2,
		TableMinSpaces:  4,
		WideGapMultiple: 3.0,

		// Image filtering
		MinImageSize: 50.0,

		// Text merging
		LineHeightMultiplier:  1.6,
		IndentThreshold:       5.0,
		WordGapMultiplier:     0.2,
		SentenceGapMultiplier: 1.5,

		// Equation detection
		EquationScoreThreshold: 2,

		// ONNX Layout Detection
		EnableONNX:        true,  // Enabled by default for better detection quality
		ONNXConfThreshold: 0.25,  // Standard YOLO confidence threshold
		ONNXNMSThreshold:  0.45,  // Standard YOLO NMS threshold
		MinONNXConfidence: 0.5,   // Override rules only with high confidence detections
		ONNXInputSize:     1024,  // DocLayout-YOLO input size
		ONNXStride:        32,    // YOLO stride for padding
		ONNXUseCoreML:     false, // Disabled: DocLayout-YOLO uses ops incompatible with CoreML
	}
}

// ConfigForAcademicPapers returns configuration optimized for academic papers.
// Academic papers often have:
// - Smaller column gaps (two-column layouts)
// - More math/equations
// - Section numbering patterns
func ConfigForAcademicPapers() *Config {
	cfg := DefaultConfig()
	cfg.ColumnGapThreshold = 15.0   // Wider gap for two-column detection
	cfg.EquationScoreThreshold = 1  // Lower threshold for equation detection
	cfg.TableMinSpaces = 5          // Higher threshold to avoid equation misdetection
	cfg.LineHeightMultiplier = 1.4  // Tighter line spacing
	cfg.SentenceGapMultiplier = 1.3 // Tighter sentence merging
	return cfg
}

// ConfigForTechnicalDocs returns configuration optimized for technical documentation.
// Technical docs often have:
// - Code blocks with keywords
// - Wider spacing
// - Clear header hierarchy
func ConfigForTechnicalDocs() *Config {
	cfg := DefaultConfig()
	cfg.HeaderSizeRatio = 1.15     // Lower threshold to catch more headers
	cfg.TableMinColumns = 3        // Higher threshold for tables
	cfg.LineHeightMultiplier = 1.8 // More generous line merging
	return cfg
}

// ConfigForScannedDocuments returns configuration optimized for scanned PDFs.
// Scanned documents often have:
// - Inconsistent spacing
// - OCR artifacts
// - Noisier text positioning
func ConfigForScannedDocuments() *Config {
	cfg := DefaultConfig()
	cfg.ColumnGapThreshold = 20.0 // Wider gap tolerance for scan artifacts
	cfg.IndentThreshold = 10.0    // More forgiving indent detection
	cfg.WordGapMultiplier = 0.3   // More forgiving word gaps
	cfg.MinImageSize = 100.0      // Filter out more noise/artifacts
	return cfg
}
