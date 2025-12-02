package pdf

import (
	"testing"
)

func TestObjects_String(t *testing.T) {
	tests := []struct {
		obj      Object
		expected string
	}{
		{Integer(123), "123"},
		{Real(12.34), "12.340000"},
		{Boolean(true), "true"},
		{Boolean(false), "false"},
		{Name("MyName"), "/MyName"},
		{StringLiteral("Hello"), "(Hello)"},
		{HexString("ABCD"), "<ABCD>"},
		{Null{}, "null"},
		{IndirectRef{1, 0}, "1 0 R"},
		{Array{Integer(1), Integer(2)}, "[1 2]"},
		{Dictionary{Name("A"): Integer(1)}, "<< /A 1 >>"},
	}

	for _, tt := range tests {
		if s := tt.obj.String(); s != tt.expected {
			t.Errorf("Expected %q, got %q", tt.expected, s)
		}
	}
}
