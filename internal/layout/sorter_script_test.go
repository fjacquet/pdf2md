package layout

import (
	"strings"
	"testing"

	"github.com/fjacquet/pdf2md/internal/extractor"
)

// "ΔT_AGN = 10^9 K.": the superscript sits 3.9 pt above the baseline and
// the subscript 1.7 pt below it; both belong to the line, in X order.
func TestSortY_SubscriptAndSuperscriptStayOnTheirLine(t *testing.T) {
	b := func(x, y, fs float64, s string) extractor.TextBlock {
		return extractor.TextBlock{Text: s, X: x, Y: y, FontSize: fs, Height: fs}
	}
	got := sortY([]extractor.TextBlock{
		b(72, 204.2, 10.9, "change of Δ"),
		b(431.4, 204.2, 10.9, "T"),
		b(437.8, 202.5, 8, "AGN"),
		b(460.9, 204.2, 10.9, "= 10"),
		b(483.8, 208.1, 8, "9"),
		b(492.5, 204.2, 10.9, "K."),
	})
	var texts []string
	for _, g := range got {
		texts = append(texts, g.Text)
	}
	if want := "change of Δ|T|AGN|= 10|9|K."; strings.Join(texts, "|") != want {
		t.Errorf("order = %q, want %q", strings.Join(texts, "|"), want)
	}
}
