package layout

import (
	"testing"
)

func TestIsCodeBlock_Shell(t *testing.T) {
	analyzer := NewAnalyzer()

	tests := []struct {
		text     string
		expected bool
	}{
		{"$ npm install", true},
		{"$ echo hello", true},
		{"$ \\alpha $", false}, // Math formula
		{"$ x = 1 $", false},   // Math formula
		{"# Comment", true},
		{"#1", false}, // Numbered list or reference?
	}

	for _, tt := range tests {
		got := analyzer.isCodeBlock(tt.text)
		if got != tt.expected {
			t.Errorf("isCodeBlock(%q) = %v, want %v", tt.text, got, tt.expected)
		}
	}
}
