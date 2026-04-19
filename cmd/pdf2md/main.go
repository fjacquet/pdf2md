// Package main provides the pdf2md command-line tool for converting PDF to Markdown.
package main

import (
	"fmt"
	"image"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/fjacquet/pdf2md/internal/extractor"
	"github.com/fjacquet/pdf2md/internal/layout"
	"github.com/fjacquet/pdf2md/internal/markdown"
	"github.com/fjacquet/pdf2md/internal/ocr"
	"github.com/fjacquet/pdf2md/internal/onnx"
	"github.com/fjacquet/pdf2md/internal/render"
)

var CLI struct {
	Input         string  `arg:"" optional:"" help:"Input PDF file path." type:"path"`
	Output        string  `arg:"" optional:"" help:"Output Markdown file path." type:"path"`
	Debug         bool    `help:"Enable debug logging."`
	ExcludeTop    float64 `name:"exclude-top" help:"Height from top to exclude (e.g. for headers)." default:"0"`
	ExcludeBottom float64 `name:"exclude-bottom" help:"Height from bottom to exclude (e.g. for footers)." default:"0"`

	// ONNX layout detection options
	NoONNX      bool   `name:"no-onnx" help:"Disable ONNX-based layout detection (uses rule-based only)."`
	ModelPath   string `name:"model-path" help:"Path to ONNX model file (auto-downloaded if not specified)."`
	RuntimePath string `name:"runtime-path" help:"Path to ONNX Runtime library (auto-detected if not specified)."`

	// OCR options (used to fall back on scanned / image-only pages)
	NoOCR           bool   `name:"no-ocr" help:"Disable OCR fallback on image-only pages."`
	OCRDetModelPath string `name:"ocr-det-model-path" help:"Path to OCR text-detection ONNX model."`
	OCRRecModelPath string `name:"ocr-rec-model-path" help:"Path to OCR text-recognition ONNX model."`
	OCRDictPath     string `name:"ocr-dict-path" help:"Path to OCR character dictionary (one char per line)."`
}

func main() {
	ctx := kong.Parse(&CLI)

	// If no input provided, print help and exit
	if CLI.Input == "" {
		_ = ctx.PrintUsage(false)
		return
	}

	// Configure logger based on Debug flag
	logLevel := slog.LevelInfo
	if CLI.Debug {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	// Ensure pdfium pool is closed on exit
	defer render.ClosePdfiumPool()

	inputPath := CLI.Input
	outputPath := CLI.Output

	// Build conversion options
	opts := convertOptions{
		excludeTop:      CLI.ExcludeTop,
		excludeBottom:   CLI.ExcludeBottom,
		enableONNX:      !CLI.NoONNX,
		modelPath:       CLI.ModelPath,
		runtimePath:     CLI.RuntimePath,
		enableOCR:       !CLI.NoOCR,
		ocrDetModelPath: CLI.OCRDetModelPath,
		ocrRecModelPath: CLI.OCRRecModelPath,
		ocrDictPath:     CLI.OCRDictPath,
	}

	// Convert PDF to Markdown
	markdown, err := convertPDFToMarkdown(inputPath, outputPath, opts)
	if err != nil {
		ctx.FatalIfErrorf(err)
	}

	// Output result
	if outputPath != "" {
		err = os.WriteFile(outputPath, []byte(markdown), 0o600)
		if err != nil {
			ctx.FatalIfErrorf(fmt.Errorf("error writing output: %w", err))
		}
		fmt.Printf("Successfully converted %s to %s\n", inputPath, outputPath)
	} else {
		fmt.Println(markdown)
	}
}

// convertOptions holds options for PDF to Markdown conversion.
type convertOptions struct {
	excludeTop      float64
	excludeBottom   float64
	enableONNX      bool
	modelPath       string
	runtimePath     string
	enableOCR       bool
	ocrDetModelPath string
	ocrRecModelPath string
	ocrDictPath     string
}

func convertPDFToMarkdown(inputPath string, outputPath string, opts convertOptions) (string, error) {
	logger := slog.Default()

	// Open PDF file (this part is kept for consistency)
	file, err := os.Open(inputPath) //nolint:gosec // G304: Path comes from CLI argument, intentional user input
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// Step 1: Initialize Extractor (using custom parser)
	ext, err := extractor.NewCustomExtractor(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create extractor: %w", err)
	}
	defer func() { _ = ext.Close() }()

	// Step 2: Initialize PDF renderer — needed for both ONNX layout and OCR.
	var renderer render.PageRenderer
	if opts.enableONNX || opts.enableOCR {
		renderer, err = render.NewPdfiumRenderer(inputPath, render.RenderConfigForONNX())
		if err != nil {
			logger.Warn("Failed to initialize PDF renderer, ONNX/OCR disabled", "error", err)
			renderer = nil
		} else {
			defer func() { _ = renderer.Close() }()
		}
	}

	// Step 3: Initialize layout configuration with ONNX settings
	config := layout.DefaultConfig()
	config.EnableONNX = opts.enableONNX && renderer != nil
	config.ONNXModelPath = opts.modelPath
	config.ONNXRuntimePath = opts.runtimePath

	// Step 4: Initialize ONNX detector if enabled
	var detector layout.LayoutDetector
	if config.EnableONNX {
		logger.Debug("ONNX enabled, initializing detector")
		onnxConfig := layout.CreateONNXConfig(config)
		det, err := onnx.NewDetector(onnxConfig, nil)
		switch {
		case err != nil:
			logger.Warn("Failed to initialize ONNX detector, using rule-based only", "error", err)
		case det.IsAvailable():
			detector = det
			defer func() { _ = det.Close() }()
			logger.Info("ONNX layout detection enabled")
		default:
			logger.Debug("ONNX detector created but not available (runtime/model not found)")
		}
	} else {
		logger.Debug("ONNX disabled", "enableONNX", opts.enableONNX, "rendererNil", renderer == nil)
	}

	// Step 4b: Initialize OCR recognizer (fallback for image-only pages).
	var recognizer ocr.Recognizer
	if opts.enableOCR && renderer != nil {
		ocrConfig := ocr.DefaultConfig()
		ocrConfig.DetModelPath = opts.ocrDetModelPath
		ocrConfig.RecModelPath = opts.ocrRecModelPath
		ocrConfig.DictPath = opts.ocrDictPath
		ocrConfig.RuntimePath = opts.runtimePath
		rec, err := ocr.NewPaddleRecognizer(ocrConfig, nil)
		switch {
		case err != nil:
			logger.Warn("Failed to initialize OCR recognizer", "error", err)
		case rec.IsAvailable():
			recognizer = rec
			defer func() { _ = rec.Close() }()
			logger.Info("OCR fallback enabled")
		default:
			logger.Debug("OCR recognizer not available (models or runtime missing)")
		}
	}

	// Initialize analyzer and builder
	analyzer := layout.NewAnalyzerWithONNX(config, detector)
	builder := markdown.NewBuilder()

	// Process each page
	var allElements []layout.Element
	pageCount, err := ext.GetPageCount()
	if err != nil {
		return "", fmt.Errorf("failed to get page count: %w", err)
	}

	for i := range pageCount {
		// Extract content
		content, err := ext.ExtractTextBlocks(i + 1)
		if err != nil {
			fmt.Printf("Error extracting page %d: %v\n", i+1, err)
			continue
		}

		// Save images
		outputDir := filepath.Dir(outputPath)
		if outputPath == "" {
			outputDir = "."
		}
		prefix := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))

		savedPaths, err := extractor.SaveImages(content.Images, outputDir, prefix)
		if err != nil {
			logger.Error("Failed to save images", "error", err)
		} else {
			// Update Image IDs to be the relative paths for Markdown
			for j, path := range savedPaths {
				content.Images[j].ID = path
			}
		}

		// Save graphics
		savedGraphicsPaths, err := extractor.SaveGraphics(content.Graphics, outputDir, prefix)
		if err != nil {
			logger.Error("Failed to save graphics", "error", err)
		} else {
			// Update Graphic IDs to be the relative paths for Markdown
			for j, path := range savedGraphicsPaths {
				content.Graphics[j].ID = path
			}
		}

		// Set exclusion zones
		analyzer.Exclusion = layout.ExclusionZone{
			Top:    opts.excludeTop,
			Bottom: opts.excludeBottom,
		}

		// Render page image — used by both ONNX detection and OCR fallback.
		var pageImage image.Image
		needsImage := (renderer != nil) && (detector != nil || recognizer != nil)
		if needsImage {
			pageImage, err = renderer.RenderPage(i, 0) // 0-indexed, default DPI
			if err != nil {
				logger.Warn("Failed to render page", "page", i+1, "error", err)
			}
		}

		// OCR fallback: when the page has no extractable text (scanned/image-only),
		// enrich content.TextBlocks with recognized text before layout analysis.
		if recognizer != nil && len(content.TextBlocks) == 0 && pageImage != nil {
			ocrBlocks, ocrErr := recognizer.RecognizePage(pageImage, content.PageWidth, content.PageHeight)
			if ocrErr != nil {
				logger.Warn("OCR failed", "page", i+1, "error", ocrErr)
			} else if len(ocrBlocks) > 0 {
				logger.Debug("OCR enriched page", "page", i+1, "blocks", len(ocrBlocks))
				content.TextBlocks = ocrBlocks
			}
		}

		// Analyze layout with ONNX detection (falls back to rule-based if no image)
		elements := analyzer.AnalyzeWithImage(content, pageImage)
		// Merge consecutive elements on same line
		elements = analyzer.MergeElements(elements)

		allElements = append(allElements, elements...)
	}

	// Post-processing: Remove Table of Contents (global)
	allElements = analyzer.RemoveTableOfContentsRange(allElements)

	// Post-processing: Filter repeated page headers (like "CHAPTER 3. INSTALLATION...")
	allElements = analyzer.FilterRepeatedPageHeaders(allElements)

	// Post-processing: Merge list item continuations across pages
	allElements = analyzer.MergeListContinuations(allElements)

	// Post-processing: Merge paragraphs across pages
	// This handles cases where a sentence starts on one page and ends on the next.
	allElements = analyzer.MergeParagraphLines(allElements)

	// Post-processing: Normalize header levels globally
	// This ensures consistent header hierarchy across all pages
	allElements = analyzer.NormalizeHeaderLevels(allElements)

	// Build Markdown
	md := builder.Build(allElements)

	return md, nil
}
