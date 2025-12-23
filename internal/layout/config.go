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
}

// DefaultConfig returns the default layout configuration.
// These defaults work well for most PDF documents.
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
