package main

import (
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
