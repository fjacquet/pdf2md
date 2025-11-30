package extractor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fjacquet/pdf2md/internal/types"
)

// SaveImages saves extracted images to the specified directory
func SaveImages(images []types.Image, outputDir, prefix string) ([]string, error) {
	if len(images) == 0 {
		return nil, nil
	}

	// Create images directory
	imgDir := filepath.Join(outputDir, "images")
	if err := os.MkdirAll(imgDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create image directory: %w", err)
	}

	var savedPaths []string
	for i, img := range images {
		// Generate filename
		// Use ID if available, otherwise index
		filename := fmt.Sprintf("%s_img_%d.%s", prefix, i+1, img.Format)
		if img.ID != "" {
			// Sanitize ID?
			filename = fmt.Sprintf("%s_%s.%s", prefix, img.ID, img.Format)
		}

		path := filepath.Join(imgDir, filename)

		// Write data
		if err := os.WriteFile(path, img.Data, 0644); err != nil {
			fmt.Printf("Warning: failed to save image %s: %v\n", filename, err)
			continue
		}

		// Store relative path for markdown
		relPath := filepath.Join("images", filename)
		savedPaths = append(savedPaths, relPath)

		// Update Image struct with path?
		// We can't modify the slice in place easily if we want to return paths.
		// But the caller might need to know which image maps to which path.
	}

	return savedPaths, nil
}

// SaveGraphics saves vector graphics to the output directory
func SaveGraphics(graphics []types.VectorGraphic, outputDir, prefix string) ([]string, error) {
	imagesDir := filepath.Join(outputDir, "images")
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		return nil, err
	}

	var savedPaths []string
	for i, vg := range graphics {
		// Generate filename
		filename := vg.ID
		if filename == "" {
			filename = fmt.Sprintf("graphic_%d", i+1)
		}
		filename = fmt.Sprintf("%s_%s.svg", prefix, filename)

		path := filepath.Join(imagesDir, filename)

		// Convert to SVG
		svgContent := ToSVG(vg)

		if err := os.WriteFile(path, []byte(svgContent), 0644); err != nil {
			return nil, err
		}

		// Return relative path for Markdown
		relPath := filepath.Join("images", filename)
		savedPaths = append(savedPaths, relPath)
	}

	return savedPaths, nil
}
