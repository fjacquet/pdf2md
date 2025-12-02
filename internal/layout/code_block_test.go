package layout

import (
	"testing"
)

func TestIsCodeBlock_Math(t *testing.T) {
	analyzer := NewAnalyzer()

	tests := []struct {
		text     string
		expected bool
	}{
		{"$n$", false},
		{"(1)", false},
		{"[1]", false},
		{"{a, b}", false},
		{"f(x) = { x if x > 0; 0 otherwise }", false}, // Math with braces
		{"func main() {", true},                       // Actual code
		{"if (x > 0) {", true},                        // Actual code
	}

	for _, tt := range tests {
		got := analyzer.isCodeBlock(tt.text)
		if got != tt.expected {
			t.Errorf("isCodeBlock(%q) = %v, want %v", tt.text, got, tt.expected)
		}
	}
}
