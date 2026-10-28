package main

import (
	"math"
	"slices"
	"testing"
)

func TestFindRegressions(t *testing.T) {
	base := map[string]float64{"a.pdf": 0.80, "b.pdf": 0.70, "c.pdf": 0.60}
	scores := map[string]float64{"a.pdf": 0.79, "b.pdf": 0.60, "d.pdf": 0.10}

	got := findRegressions(scores, base, 0.02)
	want := []string{
		"b.pdf: 0.6000 < baseline 0.7000",
		"c.pdf: missing (baseline 0.6000)",
	}
	if !slices.Equal(got, want) {
		t.Errorf("findRegressions() = %v, want %v", got, want)
	}
}

func TestExpandCmd(t *testing.T) {
	got := expandCmd("pdftotext {in} {out}", "a b.pdf", "o.md")
	want := []string{"pdftotext", "a b.pdf", "o.md"}
	if !slices.Equal(got, want) {
		t.Errorf("expandCmd() = %q, want %q", got, want)
	}
}

func TestWordF1(t *testing.T) {
	tests := []struct {
		gen, gt string
		want    float64
	}{
		{"the cat the", "the cat", 0.8}, // multiset ∩ = 2, sizes 3 and 2
		{"", "x", 0},
		{"A-b", "a b", 1}, // tokens are \w+, lowercased
	}
	for _, tt := range tests {
		if got := wordF1(tt.gen, tt.gt); math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("wordF1(%q, %q) = %v, want %v", tt.gen, tt.gt, got, tt.want)
		}
	}
}

func TestHeadingF1(t *testing.T) {
	tests := []struct {
		gen, gt   string
		withLevel bool
		want      float64
	}{
		{"# A\n## B", "# A\n## B", true, 1},
		{"# A\n# B", "# A\n## B", true, 0.5}, // one of two pairs matches
		{"# A\n# B", "# A\n## B", false, 1},
		{"## [Docs](x)", "## Docs", true, 1},
		{"", "", true, 1},
		{"text", "# A", true, 0},
	}
	for _, tt := range tests {
		if got := headingF1(tt.gen, tt.gt, tt.withLevel); math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("headingF1(%q, %q, %v) = %v, want %v", tt.gen, tt.gt, tt.withLevel, got, tt.want)
		}
	}
}
