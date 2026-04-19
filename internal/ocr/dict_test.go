package ocr

import (
	"slices"
	"strings"
	"testing"
)

func TestLoadDict_Embedded(t *testing.T) {
	chars, err := LoadDict("")
	if err != nil {
		t.Fatalf("LoadDict(\"\") returned error: %v", err)
	}
	if len(chars) < 100 {
		t.Errorf("embedded dict has %d chars, expected ≥ 100 (Latin coverage)", len(chars))
	}

	// Spot-check ASCII letters and a few Latin-specific chars needed for FR/DE.
	must := []string{"a", "Z", "0", "é", "ü", "ß"}
	for _, c := range must {
		if !slices.Contains(chars, c) {
			t.Errorf("embedded dict missing %q", c)
		}
	}
}

func TestParseDict_SkipsBlankLines(t *testing.T) {
	got := parseDict(strings.NewReader("a\n\nb\nc\n"))
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
