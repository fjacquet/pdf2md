package extractor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fjacquet/pdf2md/internal/types"
)

// SaveImages saves the images of page (1-based) as images/<prefix>_p<page>_<n>.<format>.
// The PDF's own XObject name never reaches the file system: it is attacker
// controlled (#2F decodes to "/") and repeats across pages.
func SaveImages(images []types.Image, outputDir, prefix string, page int) ([]string, error) {
	if len(images) == 0 {
		return nil, nil
	}

	// Create images directory
	imgDir := filepath.Join(outputDir, "images")
	if err := os.MkdirAll(imgDir, 0o750); err != nil {
		return nil, fmt.Errorf("failed to create image directory: %w", err)
	}

	savedPaths := make([]string, len(images)) // lines up with images; "" when a write failed
	for i, img := range images {
		filename := fmt.Sprintf("%s_p%d_%d.%s", prefix, page, i+1, img.Format)

		path := filepath.Join(imgDir, filename)

		// Write data
		if err := os.WriteFile(path, img.Data, 0o600); err != nil {
			fmt.Printf("Warning: failed to save image %s: %v\n", filename, err)
			continue
		}

		// Store relative path for markdown
		savedPaths[i] = filepath.Join("images", filename)
	}

	return savedPaths, nil
}

// SaveGraphics saves the non-trivial graphics of page (1-based) as
// images/<prefix>_p<page>_<n>.svg. The returned paths line up with graphics;
// a skipped graphic gets "".
func SaveGraphics(graphics []types.VectorGraphic, outputDir, prefix string, page int) ([]string, error) {
	imagesDir := filepath.Join(outputDir, "images")
	if err := os.MkdirAll(imagesDir, 0o750); err != nil {
		return nil, err
	}

	savedPaths := make([]string, len(graphics))
	for i, vg := range graphics {
		// Skip trivial graphics (simple rectangles, lines, decorative elements)
		if IsTrivialGraphic(vg) {
			continue
		}

		filename := fmt.Sprintf("%s_p%d_%d.svg", prefix, page, i+1)

		path := filepath.Join(imagesDir, filename)

		// Convert to SVG
		svgContent := ToSVG(vg)

		if err := os.WriteFile(path, []byte(svgContent), 0o600); err != nil {
			return nil, err
		}

		// Return relative path for Markdown
		savedPaths[i] = filepath.Join("images", filename)
	}

	return savedPaths, nil
}
