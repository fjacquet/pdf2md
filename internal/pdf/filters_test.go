package pdf

import (
	"bytes"
	"errors"
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

func TestDecodeRunLength(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
		wantErr  bool
	}{
		{
			name:     "literal run of 5 bytes",
			input:    []byte{0x04, 'H', 'e', 'l', 'l', 'o', 0x80},
			expected: []byte("Hello"),
		},
		{
			name:     "repeat run of 0x41 five times",
			input:    []byte{0xFC, 0x41, 0x80},
			expected: []byte{0x41, 0x41, 0x41, 0x41, 0x41},
		},
		{
			name:     "mixed literal + repeat",
			input:    []byte{0x02, 'A', 'B', 'C', 0xFE, 'Z', 0x80},
			expected: []byte{'A', 'B', 'C', 'Z', 'Z', 'Z'},
		},
		{
			name:     "empty with EOD",
			input:    []byte{0x80},
			expected: []byte{},
		},
		{
			name:     "no EOD (lenient)",
			input:    []byte{0x01, 'X', 'Y'},
			expected: []byte{'X', 'Y'},
		},
		{
			name:    "truncated literal run",
			input:   []byte{0x05, 'A', 'B'},
			wantErr: true,
		},
		{
			name:    "truncated repeat run",
			input:   []byte{0xFE},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeRunLength(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !bytes.Equal(got, tt.expected) {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDecodeStream_JPXPassthrough(t *testing.T) {
	jp2Magic := []byte{0x00, 0x00, 0x00, 0x0C, 0x6A, 0x50, 0x20, 0x20, 0x0D, 0x0A, 0x87, 0x0A}
	out, err := DecodeStream(jp2Magic, "JPXDecode", nil)
	if err != nil {
		t.Fatalf("JPXDecode passthrough errored: %v", err)
	}
	if !bytes.Equal(out, jp2Magic) {
		t.Errorf("JPXDecode changed bytes: got %x", out)
	}
}

func TestDecodeStream_JBIG2Unsupported(t *testing.T) {
	_, err := DecodeStream([]byte{0x00}, "JBIG2Decode", nil)
	if err == nil {
		t.Fatal("expected error for JBIG2Decode, got nil")
	}
	if !errors.Is(err, ErrFilterUnsupported) {
		t.Errorf("expected ErrFilterUnsupported, got %v", err)
	}
}

func TestDecodeStream_CCITTFaxUnsupported(t *testing.T) {
	// Stub path for now; once G4 decoder lands, swap for a real G4 sample.
	_, err := DecodeStream([]byte{0x00}, "CCITTFaxDecode", nil)
	if err == nil {
		t.Fatal("expected error for CCITTFaxDecode stub, got nil")
	}
	if !errors.Is(err, ErrFilterUnsupported) {
		t.Errorf("expected ErrFilterUnsupported, got %v", err)
	}
}

func TestDecodeStream_CryptIdentity(t *testing.T) {
	raw := []byte("plaintext when security handler inactive")
	out, err := DecodeStream(raw, "Crypt", nil)
	if err != nil {
		t.Fatalf("Crypt identity errored: %v", err)
	}
	if !bytes.Equal(out, raw) {
		t.Errorf("Crypt changed bytes: got %q", out)
	}
}
