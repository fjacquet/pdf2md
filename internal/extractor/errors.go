package extractor

import "errors"

var (
	// ErrNotInitialized is returned when extractor is used before Open()
	ErrNotInitialized = errors.New("extractor not initialized")

	// ErrInvalidPage is returned when page number is out of range
	ErrInvalidPage = errors.New("invalid page number")

	// ErrNotReaderAt is returned when the reader doesn't implement io.ReaderAt
	ErrNotReaderAt = errors.New("reader must implement io.ReaderAt")
)
