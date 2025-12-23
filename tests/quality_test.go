package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestQuality_EndToEnd(t *testing.T) {
	// Paths
	// Assuming we run test from 'tests' dir or root.
	// If running from root with 'go test ./tests/...', os.Getwd() usually returns the package dir.
	// Let's verify where we are.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	var binaryPath, inputPath string
	if strings.HasSuffix(wd, "tests") {
		binaryPath = filepath.Join("..", "pdf2md")
		inputPath = filepath.Join("..", "testdata", "test.pdf")
	} else {
		// Assume root
		binaryPath = filepath.Join(".", "pdf2md")
		inputPath = filepath.Join(".", "testdata", "test.pdf")
	}

	outputPath := filepath.Join(wd, "output_quality.md")

	// Check binary existence
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatalf("Binary not found at %s. Please run 'go build ./cmd/pdf2md' first.", binaryPath)
	}

	// Check input existence
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		t.Fatalf("Input PDF not found at %s", inputPath)
	}

	// Remove previous output
	_ = os.Remove(outputPath)
	defer func() { _ = os.Remove(outputPath) }()

	// Run command: pdf2md <input> <output>
	cmd := exec.Command(binaryPath, inputPath, outputPath) //nolint:gosec // Test binary with controlled paths
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file not created at %s", outputPath)
	}

	// Read output content
	contentBytes, err := os.ReadFile(outputPath) //nolint:gosec // Test reads output file we just created
	if err != nil {
		t.Fatal(err)
	}
	content := string(contentBytes)

	// Quality Checks

	// 1. Headers
	if !strings.Contains(content, "# ") {
		t.Error("Output missing headers (# )")
	}

	// 2. Lists
	if !strings.Contains(content, "- ") && !strings.Contains(content, "* ") {
		t.Error("Output missing unordered lists (- or *)")
	}

	// 3. Links
	// Look for [text](url) pattern roughly
	if !strings.Contains(content, "](") {
		t.Error("Output missing links (](...) pattern)")
	}

	// 4. Code Blocks
	if !strings.Contains(content, "```") {
		t.Error("Output missing code blocks (```)")
	}

	// 5. Images
	if !strings.Contains(content, "![") {
		t.Error("Output missing images (![...)")
	}

	// 6. Check for Table of Contents removal
	if strings.Contains(content, "Table of Contents") {
		t.Log("Warning: 'Table of Contents' string found in output. Verify if it should have been removed.")
	}

	t.Logf("Quality test passed. Output size: %d bytes", len(content))
}
