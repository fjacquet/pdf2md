package pdf

import (
	"reflect"
	"strings"
	"testing"
)

func TestParser_ParseObject(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Object
	}{
		{"Integer", "123", Integer(123)},
		{"Real", "12.34", Real(12.34)},
		{"Name", "/MyName", Name("MyName")},
		{"String", "(Hello)", StringLiteral("Hello")},
		{"HexString", "<4142>", HexString("4142")},
		{"Boolean True", "true", Boolean(true)},
		{"Boolean False", "false", Boolean(false)},
		{"Null", "null", Null{}},
		{"IndirectRef", "10 0 R", IndirectRef{ObjectNumber: 10, GenerationNumber: 0}},
		{"Array", "[ 1 2 3 ]", Array{Integer(1), Integer(2), Integer(3)}},
		{"Dictionary", "<< /A 1 /B (Two) >>", Dictionary{Name("A"): Integer(1), Name("B"): StringLiteral("Two")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer(strings.NewReader(tt.input))
			parser := NewParser(tokenizer)
			obj, err := parser.ParseObject()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(obj, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, obj)
			}
		})
	}
}

func TestParser_ParseStream(t *testing.T) {
	input := "<< /Length 11 >> stream\nHello World\nendstream"
	tokenizer := NewTokenizer(strings.NewReader(input))
	parser := NewParser(tokenizer)
	obj, err := parser.ParseObject()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	stream, ok := obj.(Stream)
	if !ok {
		t.Fatalf("Expected Stream, got %T", obj)
	}

	if len(stream.Dictionary) != 1 {
		t.Errorf("Expected 1 dict entry, got %d", len(stream.Dictionary))
	}

	if string(stream.Data) != "Hello World" {
		t.Errorf("Expected 'Hello World', got %q", string(stream.Data))
	}
}

func TestParser_ParseStream_UnknownLength(t *testing.T) {
	input := "<< >> stream\nData without length\nendstream"
	tokenizer := NewTokenizer(strings.NewReader(input))
	parser := NewParser(tokenizer)
	obj, err := parser.ParseObject()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	stream, ok := obj.(Stream)
	if !ok {
		t.Fatalf("Expected Stream, got %T", obj)
	}

	// Note: implementation might include EOL or not depending on tokenizer behavior
	if !strings.Contains(string(stream.Data), "Data without length") {
		t.Errorf("Expected content to contain 'Data without length', got %q", string(stream.Data))
	}
}

func TestParser_ParseObjectFromToken(t *testing.T) {
	// Test the unexported method directly
	tokenizer := NewTokenizer(strings.NewReader(""))
	parser := NewParser(tokenizer)

	tests := []struct {
		token    Token
		expected Object
	}{
		{Token{Type: TokenNumeric, Value: "42"}, Integer(42)},
		{Token{Type: TokenNumeric, Value: "3.14"}, Real(3.14)},
		{Token{Type: TokenName, Value: "Foo"}, Name("Foo")},
		{Token{Type: TokenString, Value: "Bar"}, StringLiteral("Bar")},
		{Token{Type: TokenHexString, Value: "CAFE"}, HexString("CAFE")},
		{Token{Type: TokenKeyword, Value: "true"}, Boolean(true)},
		{Token{Type: TokenKeyword, Value: "false"}, Boolean(false)},
		{Token{Type: TokenKeyword, Value: "null"}, Null{}},
	}

	for _, tt := range tests {
		obj, err := parser.parseObjectFromToken(tt.token)
		if err != nil {
			t.Errorf("Unexpected error for %v: %v", tt.token, err)
			continue
		}
		if !reflect.DeepEqual(obj, tt.expected) {
			t.Errorf("Expected %v, got %v", tt.expected, obj)
		}
	}
}
