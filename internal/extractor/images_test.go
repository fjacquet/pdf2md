package extractor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	paths, err := SaveImages(images, tmpDir, "test", 1)
	if err != nil {
		t.Fatalf("SaveImages failed: %v", err)
	}

	if len(paths) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(paths))
	}

	// Verify files exist
	expectedFiles := []string{
		"test_p1_1.png",
		"test_p1_2.jpg",
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

	paths, err := SaveGraphics(graphics, tmpDir, "test", 1)
	if err != nil {
		t.Fatalf("SaveGraphics failed: %v", err)
	}

	if len(paths) != 1 {
		t.Errorf("Expected 1 path, got %d", len(paths))
	}

	// Verify file exists
	filename := "test_p1_1.svg"
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

func curve(id string) types.VectorGraphic {
	return types.VectorGraphic{
		ID: id,
		Operations: []types.PathOperation{
			{Type: types.PathOpMoveTo, Points: []types.Point{{X: 0, Y: 0}}},
			{Type: types.PathOpCurveTo, Points: []types.Point{{X: 5, Y: 10}, {X: 15, Y: 10}, {X: 20, Y: 0}}},
			{Type: types.PathOpLineTo, Points: []types.Point{{X: 10, Y: 20}}},
			{Type: types.PathOpClose},
		},
		Width: 100, Height: 100,
	}
}

// A PDF-controlled image name (#2F decodes to "/") must not reach the file system.
func TestSaveImages_NameFromPDFIgnored(t *testing.T) {
	root := t.TempDir()
	out := filepath.Join(root, "a", "b")
	paths, err := SaveImages([]types.Image{{ID: "../../../../evil", Data: []byte("x"), Format: "png"}}, out, "test", 1)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join("images", "test_p1_1.png"); len(paths) != 1 || paths[0] != want {
		t.Errorf("paths = %v, want [%s]", paths, want)
	}
	_ = filepath.WalkDir(root, func(p string, d os.DirEntry, _ error) error {
		if !d.IsDir() && !strings.HasPrefix(p, filepath.Join(out, "images")) {
			t.Errorf("file written outside images/: %s", p)
		}
		return nil
	})
}

// pdfTeX names every page's first image Im1, and graphics restart at Path1.
func TestSave_PagesDoNotCollide(t *testing.T) {
	dir := t.TempDir()
	for page, data := range []string{"red", "blue"} {
		if _, err := SaveImages([]types.Image{{ID: "Im1", Data: []byte(data), Format: "png"}}, dir, "doc", page+1); err != nil {
			t.Fatal(err)
		}
		if _, err := SaveGraphics([]types.VectorGraphic{curve("Path1")}, dir, "doc", page+1); err != nil {
			t.Fatal(err)
		}
	}
	for page, want := range []string{"red", "blue"} {
		got, err := os.ReadFile(filepath.Join(dir, "images", fmt.Sprintf("doc_p%d_1.png", page+1))) //nolint:gosec // test file
		if err != nil || string(got) != want {
			t.Errorf("page %d image = %q, %v; want %q", page+1, got, err, want)
		}
		if _, err := os.Stat(filepath.Join(dir, "images", fmt.Sprintf("doc_p%d_1.svg", page+1))); err != nil {
			t.Errorf("page %d graphic: %v", page+1, err)
		}
	}
}

// Paths line up with the input: a skipped trivial graphic leaves "" in its slot,
// so callers can assign paths by index.
func TestSaveGraphics_PathsAlignWithInput(t *testing.T) {
	rect := types.VectorGraphic{ID: "Path1"} // no operations: trivial, skipped
	paths, err := SaveGraphics([]types.VectorGraphic{rect, curve("Path2")}, t.TempDir(), "doc", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 || paths[0] != "" || paths[1] != filepath.Join("images", "doc_p1_2.svg") {
		t.Errorf("paths = %q, want [\"\" images/doc_p1_2.svg]", paths)
	}
}
