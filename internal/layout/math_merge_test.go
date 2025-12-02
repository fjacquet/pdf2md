package layout

import (
	"testing"
)

func TestMergeMath(t *testing.T) {
	// Test mergeWithStyle for math
	tests := []struct {
		a, b     string
		expected string
	}{
		{"$a$", "$b$", "$ab$"},
		{"$a$", " $b$", "$a b$"},
		{"$a$", "$+$", "$a+$"},
		{"$a+$", "$b$", "$a+b$"},
		{"$n$", "$2$", "$n2$"}, // Should probably be n^2 but we can't know without position. At least n2 is better than $n$$2$
	}

	for _, tt := range tests {
		got := mergeWithStyle(tt.a, tt.b)
		if got != tt.expected {
			t.Errorf("mergeWithStyle(%q, %q) = %q, want %q", tt.a, tt.b, got, tt.expected)
		}
	}
}
