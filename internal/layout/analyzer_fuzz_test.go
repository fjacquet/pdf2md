package layout

import (
	"testing"

	"github.com/fjacquet/pdf2md/internal/extractor"
)

func FuzzAnalyze(f *testing.F) {
	f.Add("Some random text", 10.0, 100.0, 10.0, 200.0)
	f.Add("Header", 50.0, 500.0, 20.0, 400.0)

	f.Fuzz(func(t *testing.T, text string, x, y, w, h float64) {
		blocks := []extractor.TextBlock{
			{
				Text:     text,
				X:        x,
				Y:        y,
				Width:    w,
				Height:   h,
				FontSize: 12.0,
			},
		}

		analyzer := NewAnalyzer()
		// Should not panic
		analyzer.Analyze(blocks, nil, nil)
	})
}
