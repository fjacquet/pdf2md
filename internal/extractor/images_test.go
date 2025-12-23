package extractor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fjacquet/pdf2md/internal/types"
)

func TestSaveImages(t *testing.T) {
	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "pdf2md_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	images := []types.Image{
		{
			ID:     "img1",
			Data:   []byte("fake image data"),
			Format: "png",
		},
		{
			ID:     "", // Should auto-generate name
			Data:   []byte("more fake data"),
			Format: "jpg",
		},
	}

	paths, err := SaveImages(images, tmpDir, "test")
	if err != nil {
		t.Fatalf("SaveImages failed: %v", err)
	}

	if len(paths) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(paths))
	}

	// Verify files exist
	expectedFiles := []string{
		"test_img1.png",
		"test_img_2.jpg",
	}

	for i, filename := range expectedFiles {
		fullPath := filepath.Join(tmpDir, "images", filename)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file %s not found", fullPath)
		}

		// Check return path is relative
		expectedRel := filepath.Join("images", filename)
		if paths[i] != expectedRel {
			t.Errorf("Expected relative path %s, got %s", expectedRel, paths[i])
		}
	}
}

func TestSaveGraphics(t *testing.T) {
	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "pdf2md_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Use a complex graphic with curves (not filtered as trivial)
	graphics := []types.VectorGraphic{
		{
			ID: "graphic1",
			Operations: []types.PathOperation{
				{Type: types.PathOpMoveTo, Points: []types.Point{{X: 0, Y: 0}}},
				{Type: types.PathOpCurveTo, Points: []types.Point{{X: 5, Y: 10}, {X: 15, Y: 10}, {X: 20, Y: 0}}},
				{Type: types.PathOpLineTo, Points: []types.Point{{X: 10, Y: 20}}},
				{Type: types.PathOpClose, Points: nil},
			},
			Width:  100,
			Height: 100,
		},
	}

	paths, err := SaveGraphics(graphics, tmpDir, "test")
	if err != nil {
		t.Fatalf("SaveGraphics failed: %v", err)
	}

	if len(paths) != 1 {
		t.Errorf("Expected 1 path, got %d", len(paths))
	}

	// Verify file exists
	filename := "test_graphic1.svg"
	fullPath := filepath.Join(tmpDir, "images", filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Errorf("Expected file %s not found", fullPath)
	}

	// Check content (basic check)
	content, err := os.ReadFile(fullPath) //nolint:gosec // Test reads file we just created
	if err != nil {
		t.Fatal(err)
	}
	if len(content) == 0 {
		t.Error("SVG file is empty")
	}
}
