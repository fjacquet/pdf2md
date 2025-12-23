package pdf

import (
	"bytes"
	"testing"
)

func TestDecodeASCIIHex(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
		wantErr  bool
	}{
		{
			name:     "simple hex",
			input:    []byte("48656C6C6F>"),
			expected: []byte("Hello"),
		},
		{
			name:     "lowercase hex",
			input:    []byte("48656c6c6f>"),
			expected: []byte("Hello"),
		},
		{
			name:     "with whitespace",
			input:    []byte("48 65 6C 6C 6F>"),
			expected: []byte("Hello"),
		},
		{
			name:     "odd length (trailing nibble)",
			input:    []byte("48656C6C6F3>"),
			expected: []byte("Hello0"), // 0x30 = '0'
		},
		{
			name:     "empty",
			input:    []byte(">"),
			expected: []byte{},
		},
		{
			name:    "invalid character",
			input:   []byte("48GG>"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := decodeASCIIHex(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestDecodeASCII85(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
		wantErr  bool
	}{
		{
			name:     "simple encoding",
			input:    []byte("87cURD]j7BEbo80~>"),
			expected: []byte("Hello world!"), // Correct decoded value
		},
		{
			name:     "with z (four zeros)",
			input:    []byte("z~>"),
			expected: []byte{0, 0, 0, 0},
		},
		{
			name:     "empty",
			input:    []byte("~>"),
			expected: []byte{},
		},
		{
			name:     "partial group (2 chars)",
			input:    []byte("@/~>"), // Encodes to 0x61 (partial)
			expected: []byte{0x61},   // 'a'
		},
		{
			name:    "z inside group",
			input:   []byte("87z~>"),
			wantErr: true,
		},
		{
			name:     "Man string",
			input:    []byte("9jqo^~>"), // Known encoding of "Man "
			expected: []byte("Man "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := decodeASCII85(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("got %v (%q), want %v (%q)", result, result, tt.expected, tt.expected)
			}
		})
	}
}

func TestDecodeLZW(t *testing.T) {
	// LZW test with known encoded data
	// This is a simple test - LZW encoding of "ABABABA"
	// The encoding depends on the specific LZW variant used

	tests := []struct {
		name        string
		input       []byte
		params      Dictionary
		expectedLen int // Check length since exact encoding varies
		wantErr     bool
	}{
		{
			name:        "empty input",
			input:       []byte{},
			params:      nil,
			expectedLen: 0,
		},
		{
			name: "simple LZW stream",
			// Clear code (256) followed by 'A' (65), 'B' (66), EOD (257)
			// In 9-bit codes: 100000000 01000001 01000010 100000001
			// This creates a minimal valid LZW stream
			input:       []byte{0x80, 0x20, 0x82, 0x10, 0x10},
			params:      nil,
			expectedLen: 2, // "AB"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := decodeLZW(tt.input, tt.params)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil && tt.expectedLen > 0 {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if len(result) != tt.expectedLen {
				t.Errorf("got length %d, want %d", len(result), tt.expectedLen)
			}
		})
	}
}

func TestDecodeStream_Integration(t *testing.T) {
	// Test the DecodeStream dispatcher
	tests := []struct {
		name     string
		data     []byte
		filter   Name
		params   Dictionary
		expected []byte
	}{
		{
			name:     "ASCIIHexDecode",
			data:     []byte("48656C6C6F>"),
			filter:   "ASCIIHexDecode",
			params:   nil,
			expected: []byte("Hello"),
		},
		{
			name:     "DCTDecode passthrough",
			data:     []byte{0xFF, 0xD8, 0xFF}, // JPEG magic
			filter:   "DCTDecode",
			params:   nil,
			expected: []byte{0xFF, 0xD8, 0xFF},
		},
		{
			name:     "Unknown filter passthrough",
			data:     []byte("raw data"),
			filter:   "UnknownFilter",
			params:   nil,
			expected: []byte("raw data"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DecodeStream(tt.data, tt.filter, tt.params)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}
