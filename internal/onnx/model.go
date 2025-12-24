package onnx

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const (
	// DefaultModelURL is the Hugging Face URL for the DocLayout-YOLO ONNX model.
	DefaultModelURL = "https://huggingface.co/wybxc/DocLayout-YOLO-DocStructBench-onnx/resolve/main/doclayout_yolo_docstructbench_imgsz1024.onnx"

	// DefaultModelHash is the expected SHA256 hash of the model file.
	// This ensures integrity after download.
	DefaultModelHash = "fece9af02f618b603ff7921ccec6861d13e7e1f9830e091dfb7e8ad9311e5b21"

	// ModelFileName is the default filename for the downloaded model.
	ModelFileName = "doclayout-yolo.onnx"

	// ModelCacheDir is the subdirectory under user home for cached models.
	ModelCacheDir = ".pdf2md/models"
)

// ResolveModelPath finds or downloads the ONNX model.
//
// Resolution order:
// 1. If configPath is provided and exists, use it
// 2. Check PDF2MD_MODEL_PATH environment variable
// 3. Check default cache location (~/.pdf2md/models/doclayout-yolo.onnx)
// 4. Download from Hugging Face
func ResolveModelPath(configPath string) (string, error) {
	// 1. Check explicit config path
	if configPath != "" {
		if fileExists(configPath) {
			return configPath, nil
		}
		return "", fmt.Errorf("%w: %s", ErrModelNotFound, configPath)
	}

	// 2. Check environment variable
	if envPath := os.Getenv("PDF2MD_MODEL_PATH"); envPath != "" {
		if fileExists(envPath) {
			return envPath, nil
		}
		return "", fmt.Errorf("%w: PDF2MD_MODEL_PATH=%s", ErrModelNotFound, envPath)
	}

	// 3. Check default cache location
	cacheDir, err := getModelCacheDir()
	if err != nil {
		return "", err
	}

	modelPath := filepath.Join(cacheDir, ModelFileName)
	if fileExists(modelPath) {
		return modelPath, nil
	}

	// 4. Download model
	if err := downloadModel(DefaultModelURL, modelPath, DefaultModelHash); err != nil {
		return "", fmt.Errorf("%w: %v", ErrModelDownloadFailed, err)
	}

	return modelPath, nil
}

// ResolveRuntimePath finds the ONNX Runtime shared library.
//
// Resolution order:
// 1. If configPath is provided and exists, use it
// 2. Check ORT_LIB_PATH environment variable
// 3. Check default cache location (~/.pdf2md/lib/)
// 4. Search standard system paths
func ResolveRuntimePath(configPath string) (string, error) {
	// 1. Check explicit config path
	if configPath != "" {
		if fileExists(configPath) {
			return configPath, nil
		}
		return "", fmt.Errorf("%w: %s", ErrRuntimeNotFound, configPath)
	}

	// 2. Check environment variable
	if envPath := os.Getenv("ORT_LIB_PATH"); envPath != "" {
		if fileExists(envPath) {
			return envPath, nil
		}
		// Try as directory
		libPath := filepath.Join(envPath, getRuntimeLibName())
		if fileExists(libPath) {
			return libPath, nil
		}
		return "", fmt.Errorf("%w: ORT_LIB_PATH=%s", ErrRuntimeNotFound, envPath)
	}

	// 3. Check default cache location
	homeDir, err := os.UserHomeDir()
	if err == nil {
		libPath := filepath.Join(homeDir, ".pdf2md", "lib", getRuntimeLibName())
		if fileExists(libPath) {
			return libPath, nil
		}
	}

	// 4. Search standard system paths
	standardPaths := getStandardLibraryPaths()
	for _, dir := range standardPaths {
		libPath := filepath.Join(dir, getRuntimeLibName())
		if fileExists(libPath) {
			return libPath, nil
		}
	}

	return "", ErrRuntimeNotFound
}

// getRuntimeLibName returns the platform-specific ONNX Runtime library name.
func getRuntimeLibName() string {
	switch runtime.GOOS {
	case "darwin":
		return "libonnxruntime.dylib"
	case "windows":
		return "onnxruntime.dll"
	default: // Linux and others
		return "libonnxruntime.so"
	}
}

// getStandardLibraryPaths returns standard library search paths for the current OS.
func getStandardLibraryPaths() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{
			"/usr/local/lib",
			"/opt/homebrew/lib",
			"/opt/local/lib",
		}
	case "windows":
		return []string{
			`C:\Program Files\ONNX Runtime\lib`,
			`C:\onnxruntime\lib`,
		}
	default: // Linux
		return []string{
			"/usr/lib",
			"/usr/local/lib",
			"/usr/lib/x86_64-linux-gnu",
			"/usr/lib/aarch64-linux-gnu",
		}
	}
}

// getModelCacheDir returns the model cache directory, creating it if necessary.
func getModelCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	cacheDir := filepath.Join(homeDir, ModelCacheDir)
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cacheDir, nil
}

// downloadModel downloads a model from URL to destPath.
// If expectedHash is non-empty, verifies the download.
func downloadModel(url, destPath, expectedHash string) error {
	// Create parent directory
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return err
	}

	// Download to temp file
	tmpPath := destPath + ".tmp"
	defer os.Remove(tmpPath) // Clean up on error

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer f.Close()

	// Hash while downloading
	hash := sha256.New()
	writer := io.MultiWriter(f, hash)

	if _, err := io.Copy(writer, resp.Body); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// Close file before moving
	f.Close()

	// Verify hash if provided
	if expectedHash != "" {
		actualHash := fmt.Sprintf("%x", hash.Sum(nil))
		if actualHash != expectedHash {
			return fmt.Errorf("hash mismatch: expected %s, got %s", expectedHash, actualHash)
		}
	}

	// Move to final location
	if err := os.Rename(tmpPath, destPath); err != nil {
		return fmt.Errorf("failed to move downloaded file: %w", err)
	}

	return nil
}

// fileExists checks if a file exists and is not a directory.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DefaultModelLoader loads a model file from disk.
func DefaultModelLoader(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrModelNotFound, err)
	}
	return data, nil
}

// DefaultRuntimeInitializer validates that the ONNX Runtime library exists.
// The actual runtime initialization is handled by onnxruntime_go in detector.go.
func DefaultRuntimeInitializer(libraryPath string) error {
	if !fileExists(libraryPath) {
		return ErrRuntimeNotFound
	}
	return nil
}
