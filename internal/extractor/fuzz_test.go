package extractor

import (
	"os"
	"path/filepath"
	"testing"
)

// FuzzExtract feeds whole (mutated) PDF files through the parser and the
// content-stream interpreter: any panic or hang is a bug.
func FuzzExtract(f *testing.F) {
	seeds, _ := filepath.Glob("testdata/*/*.pdf")
	more, _ := filepath.Glob("testdata/*.pdf")
	for _, p := range append(seeds, more...) {
		if b, err := os.ReadFile(p); err == nil { //nolint:gosec // G304: p comes from a Glob over testdata
			f.Add(b)
		}
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		path := filepath.Join(t.TempDir(), "f.pdf")
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatal(err)
		}
		ext, err := NewCustomExtractor(path)
		if err != nil {
			return
		}
		defer func() { _ = ext.Close() }()
		n, err := ext.GetPageCount()
		if err != nil {
			return
		}
		for i := 1; i <= min(n, 8); i++ {
			_, _ = ext.ExtractTextBlocks(i)
		}
	})
}
