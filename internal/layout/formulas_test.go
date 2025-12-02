package layout

import "testing"

func TestIsMathFont(t *testing.T) {
	tests := []struct {
		fontName string
		expected bool
	}{
		{"CMMI10", true},
		{"CMSY10", true},
		{"Arial", false},
		{"Symbol", true},
		{"Times-Roman", false},
		{"STIXGeneral", true},
	}

	for _, tt := range tests {
		if got := isMathFont(tt.fontName); got != tt.expected {
			t.Errorf("isMathFont(%q) = %v, want %v", tt.fontName, got, tt.expected)
		}
	}
}

func TestIsMathSymbol(t *testing.T) {
	if isMathSymbol("foo") {
		t.Error("isMathSymbol should return false")
	}
}
