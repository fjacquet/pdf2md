// Package modelcache provides a single source of truth for downloading,
// caching, and resolving ONNX models and the ONNX Runtime shared library.
// Both internal/onnx and internal/ocr depend on it.
//
// All functions are pure except for the obvious I/O (filesystem, network).
package modelcache

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Sentinel errors shared by all callers.
var (
	ErrNotFound        = errors.New("modelcache: file not found")
	ErrDownloadFailed  = errors.New("modelcache: download failed")
	ErrHashMismatch    = errors.New("modelcache: hash mismatch after download")
	ErrRuntimeNotFound = errors.New("modelcache: ONNX runtime library not found")
)

// Root is the cache root relative to the user's home directory.
const Root = ".pdf2md"

// Dir returns (and creates) ~/.pdf2md/<subdir>/.
// Pass "models" for ONNX layout models, "models/ocr" for OCR models, etc.
func Dir(subdir string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("modelcache: home dir: %w", err)
	}
	dir := filepath.Join(home, Root, subdir)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return "", fmt.Errorf("modelcache: create %s: %w", dir, err)
	}
	return dir, nil
}

// Exists reports whether path points to a regular file.
func Exists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// Download fetches url to destPath, verifying SHA-256 if expectedHash is
// non-empty. The download is atomic (writes to destPath.tmp then renames).
func Download(url, destPath, expectedHash string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0o750); err != nil {
		return fmt.Errorf("%w: mkdir: %v", ErrDownloadFailed, err)
	}

	tmpPath := destPath + ".tmp"
	defer func() { _ = os.Remove(tmpPath) }()

	client := &http.Client{Timeout: 10 * time.Minute}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("%w: request: %v", ErrDownloadFailed, err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: HTTP %d", ErrDownloadFailed, resp.StatusCode)
	}

	f, err := os.Create(tmpPath) //nolint:gosec // G304: destPath is under the user's home cache dir
	if err != nil {
		return fmt.Errorf("%w: create: %v", ErrDownloadFailed, err)
	}

	hash := sha256.New()
	if _, err := io.Copy(io.MultiWriter(f, hash), resp.Body); err != nil {
		_ = f.Close()
		return fmt.Errorf("%w: copy: %v", ErrDownloadFailed, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("%w: close: %v", ErrDownloadFailed, err)
	}

	if expectedHash != "" {
		actual := fmt.Sprintf("%x", hash.Sum(nil))
		if actual != expectedHash {
			return fmt.Errorf("%w: expected %s, got %s", ErrHashMismatch, expectedHash, actual)
		}
	}

	if err := os.Rename(tmpPath, destPath); err != nil {
		return fmt.Errorf("%w: rename: %v", ErrDownloadFailed, err)
	}
	return nil
}

// ResolveOrDownload returns the first existing path among:
//  1. configPath, if non-empty (error if it doesn't exist)
//  2. the value of envVar, if set (error if it doesn't exist)
//  3. cacheDir/filename — downloaded from url (with optional sha256) if absent
func ResolveOrDownload(configPath, envVar, cacheDir, filename, url, sha256hash string) (string, error) {
	if configPath != "" {
		if Exists(configPath) {
			return configPath, nil
		}
		return "", fmt.Errorf("%w: %s", ErrNotFound, configPath)
	}

	if envVar != "" {
		if envPath := os.Getenv(envVar); envPath != "" {
			if Exists(envPath) {
				return envPath, nil
			}
			return "", fmt.Errorf("%w: %s=%s", ErrNotFound, envVar, envPath)
		}
	}

	cached := filepath.Join(cacheDir, filename)
	if Exists(cached) {
		return cached, nil
	}

	if err := Download(url, cached, sha256hash); err != nil {
		return "", err
	}
	return cached, nil
}

// RuntimeLibName returns the platform-specific ONNX Runtime library filename.
func RuntimeLibName() string {
	switch runtime.GOOS {
	case "darwin":
		return "libonnxruntime.dylib"
	case "windows":
		return "onnxruntime.dll"
	default:
		return "libonnxruntime.so"
	}
}

// StandardLibraryPaths returns standard directories to search for the
// ONNX Runtime shared library on the current OS.
func StandardLibraryPaths() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{"/usr/local/lib", "/opt/homebrew/lib", "/opt/local/lib"}
	case "windows":
		return []string{`C:\Program Files\ONNX Runtime\lib`, `C:\onnxruntime\lib`}
	default:
		return []string{
			"/usr/lib",
			"/usr/local/lib",
			"/usr/lib/x86_64-linux-gnu",
			"/usr/lib/aarch64-linux-gnu",
		}
	}
}

// ResolveRuntime locates the ONNX Runtime shared library.
//
// Resolution order:
//  1. configPath (if non-empty)
//  2. ORT_LIB_PATH env var (file or directory)
//  3. ~/.pdf2md/lib/<libname>
//  4. OS standard directories
func ResolveRuntime(configPath string) (string, error) {
	if configPath != "" {
		if Exists(configPath) {
			return configPath, nil
		}
		return "", fmt.Errorf("%w: %s", ErrRuntimeNotFound, configPath)
	}

	libName := RuntimeLibName()

	if envPath := os.Getenv("ORT_LIB_PATH"); envPath != "" {
		if Exists(envPath) {
			return envPath, nil
		}
		if p := filepath.Join(envPath, libName); Exists(p) {
			return p, nil
		}
		return "", fmt.Errorf("%w: ORT_LIB_PATH=%s", ErrRuntimeNotFound, envPath)
	}

	if home, err := os.UserHomeDir(); err == nil {
		if p := filepath.Join(home, Root, "lib", libName); Exists(p) {
			return p, nil
		}
	}

	for _, dir := range StandardLibraryPaths() {
		if p := filepath.Join(dir, libName); Exists(p) {
			return p, nil
		}
	}

	return "", ErrRuntimeNotFound
}
