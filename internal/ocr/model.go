package ocr

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fjacquet/pdf2md/internal/modelcache"
)

// Cache layout under ~/.pdf2md/:
//
//	models/ocr/det.onnx       - text detection model (DB head)
//	models/ocr/rec.onnx       - text recognition model (CRNN/SVTR head)
//	models/ocr/dict.txt       - character dictionary (one char per line)
//
// Recommended source: PaddleOCR Latin PP-OCRv3 or PP-OCRv4, exported to ONNX
// via paddle2onnx. See docs/ocr-setup.md for a one-time download script.
const (
	CacheSubdir      = "models/ocr"
	DetModelFilename = "det.onnx"
	RecModelFilename = "rec.onnx"
	DictFilename     = "dict.txt"

	EnvDet  = "PDF2MD_OCR_DET_PATH"
	EnvRec  = "PDF2MD_OCR_REC_PATH"
	EnvDict = "PDF2MD_OCR_DICT_PATH"
)

// ResolveDetModelPath locates the text-detection ONNX model.
//
// Resolution order:
//  1. configPath (if non-empty)
//  2. PDF2MD_OCR_DET_PATH env var
//  3. ~/.pdf2md/models/ocr/det.onnx
//
// Returns ErrModelNotFound with an actionable message if none resolve.
// Unlike internal/onnx, no auto-download is attempted: stable community
// mirrors for PaddleOCR ONNX are scarce, so the user installs once.
func ResolveDetModelPath(configPath string) (string, error) {
	return resolveOCRArtifact(configPath, EnvDet, DetModelFilename, "text detection")
}

// ResolveRecModelPath locates the text-recognition ONNX model.
// See ResolveDetModelPath for the resolution order.
func ResolveRecModelPath(configPath string) (string, error) {
	return resolveOCRArtifact(configPath, EnvRec, RecModelFilename, "text recognition")
}

// ResolveRuntimePath delegates to modelcache.ResolveRuntime — the same ONNX
// Runtime is shared with internal/onnx.
func ResolveRuntimePath(configPath string) (string, error) {
	path, err := modelcache.ResolveRuntime(configPath)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrRuntimeNotFound, err)
	}
	return path, nil
}

// resolveOCRArtifact is the common resolver for det/rec/dict files.
// Kept private; callers use the named helpers above.
func resolveOCRArtifact(configPath, envVar, filename, label string) (string, error) {
	if configPath != "" {
		if modelcache.Exists(configPath) {
			return configPath, nil
		}
		return "", fmt.Errorf("%w: %s", ErrModelNotFound, configPath)
	}
	if envPath := os.Getenv(envVar); envPath != "" {
		if modelcache.Exists(envPath) {
			return envPath, nil
		}
		return "", fmt.Errorf("%w: %s=%s", ErrModelNotFound, envVar, envPath)
	}
	cacheDir, err := modelcache.Dir(CacheSubdir)
	if err != nil {
		return "", err
	}
	cached := joinPath(cacheDir, filename)
	if modelcache.Exists(cached) {
		return cached, nil
	}
	return "", fmt.Errorf(
		"%w: %s model not found at %s (set --ocr-%s-model-path, the %s env var, or see docs/ocr-setup.md)",
		ErrModelNotFound, label, cached, label, envVar,
	)
}

// LoadDict reads the character dictionary (one char per line).
// If path is empty, the embedded Latin dictionary is returned.
func LoadDict(path string) ([]string, error) {
	if path == "" {
		if envPath := os.Getenv(EnvDict); envPath != "" {
			path = envPath
		}
	}
	if path == "" {
		return parseDict(strings.NewReader(embeddedLatinDict)), nil
	}

	f, err := os.Open(path) //nolint:gosec // G304: path resolved from user config or cache
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrModelNotFound, err)
	}
	defer func() { _ = f.Close() }()

	chars := parseDict(f)
	if len(chars) == 0 {
		return nil, ErrDictEmpty
	}
	return chars, nil
}

// parseDict reads one character per line, ignoring blank lines.
// Returns an empty slice if the reader yields nothing parseable.
func parseDict(r interface{ Read([]byte) (int, error) }) []string {
	sc := bufio.NewScanner(r)
	var out []string
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}

// joinPath joins a directory and filename using filepath semantics without
// importing path/filepath here (keeps this file compact and platform-correct
// because modelcache already uses filepath internally).
func joinPath(dir, name string) string { return dir + string(os.PathSeparator) + name }
