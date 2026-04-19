package onnx

import (
	"fmt"
	"os"

	"github.com/fjacquet/pdf2md/internal/modelcache"
)

const (
	// DefaultModelURL is the Hugging Face URL for the DocLayout-YOLO ONNX model.
	DefaultModelURL = "https://huggingface.co/wybxc/DocLayout-YOLO-DocStructBench-onnx/resolve/main/doclayout_yolo_docstructbench_imgsz1024.onnx"

	// DefaultModelHash is the expected SHA256 hash of the model file.
	DefaultModelHash = "fece9af02f618b603ff7921ccec6861d13e7e1f9830e091dfb7e8ad9311e5b21"

	// ModelFileName is the default filename for the downloaded model.
	ModelFileName = "doclayout-yolo.onnx"

	// ModelCacheDir is the subdirectory under user home for cached models.
	ModelCacheDir = "models"
)

// ResolveModelPath finds or downloads the DocLayout-YOLO ONNX model.
//
// Resolution order:
//  1. configPath (if non-empty)
//  2. PDF2MD_MODEL_PATH env var
//  3. ~/.pdf2md/models/doclayout-yolo.onnx — downloaded from Hugging Face if absent
func ResolveModelPath(configPath string) (string, error) {
	cacheDir, err := modelcache.Dir(ModelCacheDir)
	if err != nil {
		return "", err
	}
	path, err := modelcache.ResolveOrDownload(
		configPath, "PDF2MD_MODEL_PATH",
		cacheDir, ModelFileName, DefaultModelURL, DefaultModelHash,
	)
	if err != nil {
		return "", translateCacheErr(err)
	}
	return path, nil
}

// ResolveRuntimePath delegates to modelcache.ResolveRuntime.
func ResolveRuntimePath(configPath string) (string, error) {
	path, err := modelcache.ResolveRuntime(configPath)
	if err != nil {
		return "", translateCacheErr(err)
	}
	return path, nil
}

// DefaultModelLoader loads a model file from disk.
func DefaultModelLoader(path string) ([]byte, error) {
	data, err := os.ReadFile(path) //nolint:gosec // G304: model path resolved from user config or cache
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrModelNotFound, err)
	}
	return data, nil
}

// DefaultRuntimeInitializer validates that the ONNX Runtime library exists.
// The actual runtime initialization is handled by onnxruntime_go in detector.go.
func DefaultRuntimeInitializer(libraryPath string) error {
	if !modelcache.Exists(libraryPath) {
		return ErrRuntimeNotFound
	}
	return nil
}

// translateCacheErr maps modelcache sentinels to this package's public errors
// so callers keep getting ErrModelNotFound / ErrModelDownloadFailed / ErrRuntimeNotFound.
func translateCacheErr(err error) error {
	switch {
	case err == nil:
		return nil
	case errorsIs(err, modelcache.ErrRuntimeNotFound):
		return ErrRuntimeNotFound
	case errorsIs(err, modelcache.ErrNotFound):
		return fmt.Errorf("%w: %v", ErrModelNotFound, err)
	case errorsIs(err, modelcache.ErrDownloadFailed),
		errorsIs(err, modelcache.ErrHashMismatch):
		return fmt.Errorf("%w: %v", ErrModelDownloadFailed, err)
	default:
		return err
	}
}

// errorsIs is a tiny wrapper to avoid a direct "errors" import clash with
// this file's package-scoped errors.go. Using a local helper keeps the
// callsite readable.
func errorsIs(err, target error) bool {
	for e := err; e != nil; {
		if e == target {
			return true
		}
		type unwrapper interface{ Unwrap() error }
		u, ok := e.(unwrapper)
		if !ok {
			return false
		}
		e = u.Unwrap()
	}
	return false
}
