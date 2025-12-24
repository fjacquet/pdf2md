package onnx

import "errors"

// Error definitions for the onnx package.
var (
	// ErrDetectorNotAvailable is returned when the detector is not initialized or available.
	ErrDetectorNotAvailable = errors.New("onnx detector not available")

	// ErrModelNotFound is returned when the ONNX model file cannot be found.
	ErrModelNotFound = errors.New("onnx model not found")

	// ErrRuntimeNotFound is returned when the ONNX Runtime library cannot be found.
	ErrRuntimeNotFound = errors.New("onnx runtime library not found")

	// ErrInvalidImage is returned when the input image cannot be decoded.
	ErrInvalidImage = errors.New("invalid image data")

	// ErrInferenceFailed is returned when ONNX inference fails.
	ErrInferenceFailed = errors.New("onnx inference failed")

	// ErrModelDownloadFailed is returned when the model download fails.
	ErrModelDownloadFailed = errors.New("model download failed")
)
