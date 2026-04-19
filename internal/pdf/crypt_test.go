package pdf

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"testing"
)

// pkcs7Pad appends the PKCS#7 padding bytes to bring data up to a multiple
// of blockSize. Test helper only.
func pkcs7Pad(data []byte, blockSize int) []byte {
	pad := blockSize - (len(data) % blockSize)
	out := make([]byte, len(data)+pad)
	copy(out, data)
	for i := len(data); i < len(out); i++ {
		out[i] = byte(pad) //nolint:gosec // G115: pad ∈ [1, blockSize] ≤ 16, fits in a byte by construction
	}
	return out
}

// aesEncryptForTest encrypts plaintext with AES-128-CBC using the given IV.
// Test helper only; production code never encrypts, only decrypts.
func aesEncryptForTest(t *testing.T, key, iv, plaintext []byte) []byte {
	t.Helper()
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("aes.NewCipher: %v", err)
	}
	out := make([]byte, len(plaintext))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(out, plaintext)
	return out
}

// TestPadPassword verifies the fixed 32-byte padding algorithm.
func TestPadPassword(t *testing.T) {
	tests := []struct {
		name     string
		in       []byte
		wantLast byte // sanity: first byte of the padding constant = 0x28
	}{
		{"empty password", nil, 0x28},
		{"short password", []byte("hi"), 0x28},
		{"exactly 32 bytes", bytes.Repeat([]byte{0x41}, 32), 0x41},
		{"over 32 bytes (truncated)", bytes.Repeat([]byte{0x42}, 40), 0x42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := padPassword(tt.in)
			if len(got) != 32 {
				t.Fatalf("want 32 bytes, got %d", len(got))
			}
		})
	}
	// Specifically: empty password should begin with the padding constant.
	padded := padPassword(nil)
	if padded[0] != 0x28 || padded[31] != 0x7A {
		t.Errorf("empty-password padding wrong: got first=%#x last=%#x", padded[0], padded[31])
	}
}

// TestDeriveFileKey_Roundtrip verifies that the key derivation + /U
// computation are self-consistent: computing /U from a derived key and then
// verifying recovers the original result.
func TestDeriveFileKey_Roundtrip(t *testing.T) {
	// Arbitrary but fixed inputs emulating a V=2/R=3 encrypt dict.
	o := bytes.Repeat([]byte{0xAA}, 32)
	id := bytes.Repeat([]byte{0xCC}, 16)
	p := int32(-1852) // typical Acrobat permission value
	keyLen := 128

	key := deriveFileKey(nil, o, p, id, 3, keyLen, true)
	if len(key) != 16 {
		t.Fatalf("want 16-byte key, got %d", len(key))
	}

	computed := computeU(key, id, 3)
	if !verifyUserPassword(key, id, computed, 3) {
		t.Error("roundtrip: verifyUserPassword rejected the key it derived")
	}

	// Sanity: a flipped bit in the stored /U should fail verification.
	bad := append([]byte{}, computed...)
	bad[0] ^= 0x01
	if verifyUserPassword(key, id, bad, 3) {
		t.Error("verifyUserPassword accepted tampered /U")
	}
}

// TestObjectKey_Shape exercises Algorithm 1 (§7.6.3.1): per-object key length
// is min(FileKey.len + 5, 16) for RC4, and always 16 for AES.
func TestObjectKey_Shape(t *testing.T) {
	h := &SecurityHandler{FileKey: bytes.Repeat([]byte{0x01}, 16)}

	rc4Key := h.ObjectKey(7, 0, cipherRC4)
	if len(rc4Key) != 16 { // 16 + 5 = 21, capped at 16
		t.Errorf("RC4 object key: want 16 bytes, got %d", len(rc4Key))
	}

	aesKey := h.ObjectKey(7, 0, cipherAES128)
	if len(aesKey) != 16 {
		t.Errorf("AES object key: want 16 bytes, got %d", len(aesKey))
	}

	if bytes.Equal(rc4Key, aesKey) {
		t.Error("RC4 and AES object keys should differ (AES appends sAlT)")
	}
}

// TestObjectKey_ShortFileKey covers the RC4-40 case where FileKey is 5 bytes
// and the per-object key is 10 bytes (5 + 5).
func TestObjectKey_ShortFileKey(t *testing.T) {
	h := &SecurityHandler{FileKey: bytes.Repeat([]byte{0x01}, 5)}
	k := h.ObjectKey(1, 0, cipherRC4)
	if len(k) != 10 {
		t.Errorf("want 10-byte key for RC4-40, got %d", len(k))
	}
}

// TestDecryptStream_RC4Roundtrip encrypts a payload with the same per-object
// key derivation the handler uses, then asserts DecryptStream recovers it.
func TestDecryptStream_RC4Roundtrip(t *testing.T) {
	h := &SecurityHandler{
		FileKey:      bytes.Repeat([]byte{0x42}, 16),
		StreamCipher: cipherRC4,
		StringCipher: cipherRC4,
	}
	plain := []byte("BT /F1 12 Tf (Hello) Tj ET")

	objNum, genNum := 42, 0
	key := h.ObjectKey(objNum, genNum, cipherRC4)
	ciphertext := rc4Apply(key, plain)

	got, err := h.DecryptStream(ciphertext, objNum, genNum)
	if err != nil {
		t.Fatalf("DecryptStream: %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Errorf("roundtrip mismatch:\n got  %q\n want %q", got, plain)
	}
}

// TestDecryptStream_AES128Roundtrip does the same for AES-128-CBC with a
// random IV prepended and PKCS#7 padding.
func TestDecryptStream_AES128Roundtrip(t *testing.T) {
	h := &SecurityHandler{
		FileKey:      bytes.Repeat([]byte{0x01}, 16),
		StreamCipher: cipherAES128,
		StringCipher: cipherAES128,
	}
	plain := []byte("q 1 0 0 1 100 100 cm /Im1 Do Q")

	objNum, genNum := 7, 0
	key := h.ObjectKey(objNum, genNum, cipherAES128)

	// Encrypt with AES-128-CBC, fixed IV (test-only), PKCS#7 pad.
	iv := bytes.Repeat([]byte{0x33}, 16)
	padded := pkcs7Pad(plain, 16)
	ciphertext := aesEncryptForTest(t, key, iv, padded)
	blob := append(append([]byte{}, iv...), ciphertext...)

	got, err := h.DecryptStream(blob, objNum, genNum)
	if err != nil {
		t.Fatalf("DecryptStream: %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Errorf("AES roundtrip mismatch:\n got  %q\n want %q", got, plain)
	}
}

// TestNewSecurityHandler_UnsupportedV rejects V=5 with ErrEncrypted.
func TestNewSecurityHandler_UnsupportedV(t *testing.T) {
	dict := Dictionary{
		Name("V"):      Integer(5),
		Name("R"):      Integer(6),
		Name("Length"): Integer(256),
		Name("O"):      StringLiteral(bytes.Repeat([]byte{0xAA}, 48)),
		Name("U"):      StringLiteral(bytes.Repeat([]byte{0xBB}, 48)),
		Name("P"):      Integer(-1),
	}
	_, err := NewSecurityHandler(dict, bytes.Repeat([]byte{0xCC}, 16))
	if !errors.Is(err, ErrEncrypted) {
		t.Fatalf("want ErrEncrypted, got %v", err)
	}
}

// TestNewSecurityHandler_NonStandardFilter rejects public-key etc.
func TestNewSecurityHandler_NonStandardFilter(t *testing.T) {
	dict := Dictionary{
		Name("Filter"): Name("Adobe.PubSec"),
		Name("V"):      Integer(4),
		Name("R"):      Integer(4),
		Name("Length"): Integer(128),
		Name("O"):      StringLiteral(bytes.Repeat([]byte{0xAA}, 32)),
		Name("U"):      StringLiteral(bytes.Repeat([]byte{0xBB}, 32)),
		Name("P"):      Integer(-1),
	}
	_, err := NewSecurityHandler(dict, bytes.Repeat([]byte{0xCC}, 16))
	if !errors.Is(err, ErrEncrypted) {
		t.Fatalf("want ErrEncrypted, got %v", err)
	}
}

func TestStripPKCS7(t *testing.T) {
	tests := []struct {
		name    string
		in      []byte
		want    []byte
		wantErr bool
	}{
		{"one pad byte", append([]byte("abc"), []byte{0x01}...), []byte("abc"), false},
		{"full block pad", bytes.Repeat([]byte{0x10}, 16), []byte{}, false},
		{"bad pad value", append([]byte("ab"), []byte{0xFF, 0xFF}...), nil, true},
		{"mixed pad bytes", append([]byte("ab"), []byte{0x02, 0x03}...), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stripPKCS7(tt.in, 16)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
