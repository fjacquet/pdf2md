package pdf

import (
	"bytes"
	"strings"
	"testing"
)

func TestTokenizer_NextToken(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			input: "  % Comment\n/Name (String) <ABCD> [ 123 45.67 ] << >> true false null R",
			expected: []Token{
				{Type: TokenName, Value: "Name"},
				{Type: TokenString, Value: "String"},
				{Type: TokenHexString, Value: "ABCD"},
				{Type: TokenArrayStart, Value: "["},
				{Type: TokenNumeric, Value: "123"},
				{Type: TokenNumeric, Value: "45.67"},
				{Type: TokenArrayEnd, Value: "]"},
				{Type: TokenDictStart, Value: "<<"},
				{Type: TokenDictEnd, Value: ">>"},
				{Type: TokenKeyword, Value: "true"},
				{Type: TokenKeyword, Value: "false"},
				{Type: TokenKeyword, Value: "null"},
				{Type: TokenKeyword, Value: "R"},
			},
		},
		{
			input: "/Name#20With#20Spaces",
			expected: []Token{
				{Type: TokenName, Value: "Name With Spaces"}, // Hex escapes decoded: #20 = space
			},
		},
	}

	for i, tt := range tests {
		tokenizer := NewTokenizer(strings.NewReader(tt.input))
		for j, expected := range tt.expected {
			tok, err := tokenizer.NextToken()
			if err != nil {
				t.Errorf("Test %d: Unexpected error at token %d: %v", i, j, err)
				break
			}
			if tok.Type != expected.Type {
				t.Errorf("Test %d: Token %d type mismatch: got %v, want %v", i, j, tok.Type, expected.Type)
			}
			if tok.Value != expected.Value {
				t.Errorf("Test %d: Token %d value mismatch: got %q, want %q", i, j, tok.Value, expected.Value)
			}
		}
		// Expect EOF
		tok, err := tokenizer.NextToken()
		if err == nil && tok.Type != TokenEOF {
			t.Errorf("Test %d: Expected EOF, got %v", i, tok)
		}
	}
}

func TestTokenizer_ReadString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`(Simple String)`, "Simple String"},
		{`(Nested (Parens))`, "Nested (Parens)"},
		{`(Escaped \n \r \t \b \f \( \) \\)`, "Escaped \n \r \t \b \f ( ) \\"},
		{`(Octal \101 \102 \103)`, "Octal A B C"},
		{`(Split \
Line)`, "Split Line"}, // Backslash at EOL ignored? Implementation says yes.
	}

	for i, tt := range tests {
		tokenizer := NewTokenizer(strings.NewReader(tt.input))
		tok, err := tokenizer.NextToken()
		if err != nil {
			t.Errorf("Test %d: Unexpected error: %v", i, err)
			continue
		}
		if tok.Type != TokenString {
			t.Errorf("Test %d: Expected String token, got %v", i, tok.Type)
		}
		if tok.Value != tt.expected {
			t.Errorf("Test %d: Value mismatch: got %q, want %q", i, tok.Value, tt.expected)
		}
	}
}

func TestTokenizer_ReadHexString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`<414243>`, "414243"},
		{`< 41 42 43 >`, "414243"},
		{`<41424>`, "41424"}, // Odd length is valid?
	}

	for i, tt := range tests {
		tokenizer := NewTokenizer(strings.NewReader(tt.input))
		tok, err := tokenizer.NextToken()
		if err != nil {
			t.Errorf("Test %d: Unexpected error: %v", i, err)
			continue
		}
		if tok.Type != TokenHexString {
			t.Errorf("Test %d: Expected HexString token, got %v", i, tok.Type)
		}
		if tok.Value != tt.expected {
			t.Errorf("Test %d: Value mismatch: got %q, want %q", i, tok.Value, tt.expected)
		}
	}
}

func TestTokenizer_ReadStream(t *testing.T) {
	// Test with length
	// Note: ReadStream expects to be called after "stream" keyword is consumed.
	// NextToken consumes "stream" and returns it, so the reader is positioned after "stream".
	// We simulate that by starting with the EOL after "stream".
	tokenizer := NewTokenizer(strings.NewReader("\nHello World\nendstream"))
	content, err := tokenizer.ReadStream(11)
	if err != nil {
		t.Fatalf("ReadStream failed: %v", err)
	}
	if string(content) != "Hello World" {
		t.Errorf("Expected 'Hello World', got %q", string(content))
	}

	// Test without length (scan for endstream)
	tokenizer = NewTokenizer(strings.NewReader("\nData without length\nendstream"))
	content, err = tokenizer.ReadStream(-1)
	if err != nil {
		t.Fatalf("ReadStream(-1) failed: %v", err)
	}
	// Note: The implementation might include the EOL before endstream or not.
	// Let's check what we get.
	// "Data without length" (19 chars)
	if !bytes.Contains(content, []byte("Data without length")) {
		t.Errorf("Expected content to contain 'Data without length', got %q", string(content))
	}
}
