package pdf

// Standard security handler for encrypted PDFs (ISO 32000-1 §7.6).
//
// Scope implemented in this pass:
//   - V=1 (RC4 40-bit, R=2)
//   - V=2 (RC4 variable length up to 128-bit, R=3)
//   - V=4 (RC4 or AES-128, R=4) — /CF with /V2 or /AESV2
//   - V=5 (AES-256, R=6) — /CF with /AESV3
//
// The blank password is tried first (the common "owner-password-only" case
// produced by Word and similar exporters), then the caller's user password.
//
// Out of scope (returns ErrEncrypted with a clear message):
//   - V=5 / R=5 (Adobe extension, superseded by R=6)
//   - /Adobe.PubSec public-key security handler
//   - Owner passwords

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5" //nolint:gosec // G501: MD5 required by PDF spec §7.6.3.3 for key derivation
	"crypto/rc4" //nolint:gosec // G401: RC4 required by PDF spec §7.6.3.4
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

// ErrEncrypted indicates a PDF is encrypted but we cannot (or will not)
// decrypt it. The wrapped message explains why.
var ErrEncrypted = errors.New("pdf: cannot decrypt document")

// passwordPadding is the 32-byte constant from ISO 32000-1 §7.6.3.3.
var passwordPadding = []byte{
	0x28, 0xBF, 0x4E, 0x5E, 0x4E, 0x75, 0x8A, 0x41,
	0x64, 0x00, 0x4E, 0x56, 0xFF, 0xFA, 0x01, 0x08,
	0x2E, 0x2E, 0x00, 0xB6, 0xD0, 0x68, 0x3E, 0x80,
	0x2F, 0x0C, 0xA9, 0xFE, 0x64, 0x53, 0x69, 0x7A,
}

// cipherMethod is the per-crypt-filter algorithm selector (resolved from
// the /CF dictionary when /V = 4).
type cipherMethod int

const (
	cipherNone   cipherMethod = iota // /Identity — never encrypted
	cipherRC4                        // /V2 (RC4 up to 128-bit)
	cipherAES128                     // /AESV2
	cipherAES256                     // /AESV3
)

// SecurityHandler holds everything needed to decrypt a document's streams
// and strings after the user password has been verified.
type SecurityHandler struct {
	V            int    // /V
	R            int    // /R
	KeyLength    int    // /Length in bits (default 40)
	FileKey      []byte // derived encryption key (length KeyLength/8)
	StreamCipher cipherMethod
	StringCipher cipherMethod
	EmbedCipher  cipherMethod // embedded-file crypt (not used today but parsed)
	EncryptMeta  bool
}

// NewSecurityHandler parses the /Encrypt dictionary and the trailer /ID,
// verifies the blank user password and then password (if non-empty), and
// returns a handler that can decrypt strings and streams. Returns ErrEncrypted
// if the scheme is unsupported or neither password opens the document.
func NewSecurityHandler(encrypt Dictionary, id, password []byte) (*SecurityHandler, error) {
	h := &SecurityHandler{EncryptMeta: true}

	if v, ok := intFromDict(encrypt, "V"); ok {
		h.V = v
	}
	if r, ok := intFromDict(encrypt, "R"); ok {
		h.R = r
	}
	if l, ok := intFromDict(encrypt, "Length"); ok {
		h.KeyLength = l
	} else {
		h.KeyLength = 40
	}
	if em, ok := encrypt[Name("EncryptMetadata")].(Boolean); ok {
		h.EncryptMeta = bool(em)
	}

	// Reject unsupported variants early.
	switch h.V {
	case 1, 2:
		h.StreamCipher = cipherRC4
		h.StringCipher = cipherRC4
	case 4, 5:
		stmF, _ := encrypt[Name("StmF")].(Name)
		strF, _ := encrypt[Name("StrF")].(Name)
		eff, _ := encrypt[Name("EFF")].(Name)
		cf, _ := encrypt[Name("CF")].(Dictionary)
		h.StreamCipher = resolveCipher(cf, stmF)
		h.StringCipher = resolveCipher(cf, strF)
		if eff != "" {
			h.EmbedCipher = resolveCipher(cf, eff)
		}
	default:
		return nil, fmt.Errorf("%w: /V=%d unsupported", ErrEncrypted, h.V)
	}

	if filter, _ := encrypt[Name("Filter")].(Name); filter != "" && filter != "Standard" {
		return nil, fmt.Errorf("%w: /Filter=%s (only Standard supported)", ErrEncrypted, filter)
	}

	if h.V == 5 {
		if h.R != 6 {
			return nil, fmt.Errorf("%w: V=5 R=%d unsupported (only R=6)", ErrEncrypted, h.R)
		}
		h.KeyLength = 256
		key, err := openAES256(encrypt, password)
		if err != nil {
			return nil, err
		}
		h.FileKey = key
		return h, nil
	}

	oVal, err := stringBytes(encrypt[Name("O")])
	if err != nil {
		return nil, fmt.Errorf("%w: missing /O: %v", ErrEncrypted, err)
	}
	uVal, err := stringBytes(encrypt[Name("U")])
	if err != nil {
		return nil, fmt.Errorf("%w: missing /U: %v", ErrEncrypted, err)
	}
	p, ok := intFromDict(encrypt, "P")
	if !ok {
		return nil, fmt.Errorf("%w: missing /P", ErrEncrypted)
	}

	for _, pw := range passwordCandidates(password) {
		key := deriveFileKey(pw, oVal, int32(p), id, h.R, h.KeyLength, h.EncryptMeta) //nolint:gosec // G115: /P is a 32-bit signed flags field per spec; narrowing is intentional
		if verifyUserPassword(key, id, uVal, h.R) {
			h.FileKey = key
			return h, nil
		}
	}
	return nil, errBadPassword(password)
}

// passwordCandidates lists the passwords to try, in order: blank, then the
// caller's.
func passwordCandidates(password []byte) [][]byte {
	if len(password) == 0 {
		return [][]byte{nil}
	}
	return [][]byte{nil, password}
}

func errBadPassword(password []byte) error {
	if len(password) == 0 {
		return fmt.Errorf("%w: user password required", ErrEncrypted)
	}
	return fmt.Errorf("%w: wrong password", ErrEncrypted)
}

// openAES256 implements Algorithm 2.A for the user password (ISO 32000-2
// §7.6.4.3.3): it checks the password against /U and unwraps the file key
// from /UE.
func openAES256(encrypt Dictionary, password []byte) ([]byte, error) {
	u, err := stringBytes(encrypt[Name("U")])
	if err != nil || len(u) < 48 {
		return nil, fmt.Errorf("%w: missing or short /U", ErrEncrypted)
	}
	ue, err := stringBytes(encrypt[Name("UE")])
	if err != nil || len(ue) < 32 {
		return nil, fmt.Errorf("%w: missing or short /UE", ErrEncrypted)
	}

	for _, pw := range passwordCandidates(password) {
		// shortcut: no SASLprep normalisation, only the 127-byte truncation; add it if non-ASCII passwords fail.
		pw = pw[:min(len(pw), 127)]
		if !bytes.Equal(hashR6(pw, u[32:40], nil), u[:32]) {
			continue
		}
		block, err := aes.NewCipher(hashR6(pw, u[40:48], nil))
		if err != nil {
			return nil, fmt.Errorf("%w: aes.NewCipher: %v", ErrEncrypted, err)
		}
		key := make([]byte, 32)
		cipher.NewCBCDecrypter(block, make([]byte, aes.BlockSize)).CryptBlocks(key, ue[:32])
		return key, nil
	}
	return nil, errBadPassword(password)
}

// hashR6 is Algorithm 2.B (ISO 32000-2 §7.6.4.3.4).
func hashR6(password, salt, udata []byte) []byte {
	sum := sha256.Sum256(bytes.Join([][]byte{password, salt, udata}, nil))
	k := sum[:]

	for round := 1; ; round++ {
		k1 := bytes.Repeat(bytes.Join([][]byte{password, k, udata}, nil), 64)
		block, _ := aes.NewCipher(k[:16]) // k is always >= 32 bytes
		e := make([]byte, len(k1))
		cipher.NewCBCEncrypter(block, k[16:32]).CryptBlocks(e, k1)

		// The first 16 bytes as a big-endian number mod 3 equal their byte sum mod 3.
		n := 0
		for _, b := range e[:16] {
			n += int(b)
		}
		switch n % 3 {
		case 0:
			s := sha256.Sum256(e)
			k = s[:]
		case 1:
			s := sha512.Sum384(e)
			k = s[:]
		default:
			s := sha512.Sum512(e)
			k = s[:]
		}

		if round >= 64 && int(e[len(e)-1]) <= round-32 {
			return k[:32]
		}
	}
}

// resolveCipher looks up /CF[name]/CFM and maps it to our cipherMethod.
func resolveCipher(cf Dictionary, name Name) cipherMethod {
	if name == "" || name == "Identity" {
		return cipherNone
	}
	entry, _ := cf[name].(Dictionary)
	if entry == nil {
		return cipherRC4
	}
	switch cfm, _ := entry[Name("CFM")].(Name); cfm {
	case "V2":
		return cipherRC4
	case "AESV2":
		return cipherAES128
	case "AESV3":
		return cipherAES256
	case "None":
		return cipherNone
	default:
		return cipherRC4
	}
}

// deriveFileKey implements Algorithm 2 (§7.6.3.3).
func deriveFileKey(password, o []byte, p int32, id []byte, r, keyLen int, encryptMeta bool) []byte {
	// Step a: pad password to 32 bytes.
	pw := padPassword(password)

	//nolint:gosec // G401/G501: MD5 is the spec-mandated KDF
	h := md5.New()
	h.Write(pw)
	h.Write(o)
	var pBuf [4]byte
	binary.LittleEndian.PutUint32(pBuf[:], uint32(p)) //nolint:gosec // G115: /P is 32-bit signed flags; bit pattern preserved across signed↔unsigned conversion
	h.Write(pBuf[:])
	h.Write(id)
	if r >= 4 && !encryptMeta {
		h.Write([]byte{0xFF, 0xFF, 0xFF, 0xFF})
	}
	digest := h.Sum(nil)

	n := keyLen / 8
	if r >= 3 {
		for i := 0; i < 50; i++ {
			//nolint:gosec // G401/G501
			d := md5.Sum(digest[:n])
			digest = d[:]
		}
	}
	out := make([]byte, n)
	copy(out, digest[:n])
	return out
}

// padPassword returns a 32-byte slice: password truncated or padded with the
// fixed padding constant.
func padPassword(pw []byte) []byte {
	out := make([]byte, 32)
	copy(out, pw)
	if len(pw) < 32 {
		copy(out[len(pw):], passwordPadding[:32-len(pw)])
	}
	return out
}

// verifyUserPassword implements Algorithm 6 (§7.6.3.4): compute the expected
// /U from the derived key and compare against the stored value.
func verifyUserPassword(key, id, storedU []byte, r int) bool {
	computed := computeU(key, id, r)
	if r >= 3 {
		// Compare only first 16 bytes (spec: remaining 16 are arbitrary padding).
		if len(computed) < 16 || len(storedU) < 16 {
			return false
		}
		return bytes.Equal(computed[:16], storedU[:16])
	}
	return bytes.Equal(computed, storedU)
}

// computeU implements Algorithm 4 (R=2) or Algorithm 5 (R>=3).
func computeU(key, id []byte, r int) []byte {
	if r == 2 {
		// RC4(key, padding)
		return rc4Apply(key, passwordPadding)
	}
	// R >= 3
	//nolint:gosec // G401/G501
	h := md5.New()
	h.Write(passwordPadding)
	h.Write(id)
	buf := h.Sum(nil) // 16 bytes

	buf = rc4Apply(key, buf)
	for i := 1; i <= 19; i++ {
		xorKey := make([]byte, len(key))
		for j, b := range key {
			xorKey[j] = b ^ byte(i)
		}
		buf = rc4Apply(xorKey, buf)
	}
	// Pad to 32 bytes (content of last 16 is irrelevant per spec).
	out := make([]byte, 32)
	copy(out, buf)
	return out
}

// ObjectKey derives the per-object key used for strings and streams inside
// object (objNum, genNum). method selects the KDF tweak: AES appends "sAlT".
func (h *SecurityHandler) ObjectKey(objNum, genNum int, method cipherMethod) []byte {
	//nolint:gosec // G401/G501: MD5 required by spec §7.6.3.1 algorithm 1
	hash := md5.New()
	hash.Write(h.FileKey)
	var buf [5]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(objNum))  //nolint:gosec // G115: PDF object numbers fit well within 32-bit signed range
	binary.LittleEndian.PutUint16(buf[3:], uint16(genNum)) //nolint:gosec // G115: generation numbers are capped at 65535 by spec
	hash.Write(buf[:5])
	if method == cipherAES128 {
		hash.Write([]byte{'s', 'A', 'l', 'T'})
	}
	digest := hash.Sum(nil)

	n := len(h.FileKey) + 5
	if n > 16 {
		n = 16
	}
	out := make([]byte, n)
	copy(out, digest[:n])
	return out
}

// DecryptString decrypts the bytes of a string belonging to object (objNum, genNum).
func (h *SecurityHandler) DecryptString(data []byte, objNum, genNum int) ([]byte, error) {
	return h.decrypt(data, objNum, genNum, h.StringCipher)
}

// DecryptStream decrypts the raw bytes of a stream belonging to object (objNum, genNum).
// The caller should run DecodeStream afterwards.
func (h *SecurityHandler) DecryptStream(data []byte, objNum, genNum int) ([]byte, error) {
	return h.decrypt(data, objNum, genNum, h.StreamCipher)
}

func (h *SecurityHandler) decrypt(data []byte, objNum, genNum int, method cipherMethod) ([]byte, error) {
	switch method {
	case cipherNone:
		return data, nil
	case cipherRC4:
		key := h.ObjectKey(objNum, genNum, method)
		return rc4Apply(key, data), nil
	case cipherAES128:
		key := h.ObjectKey(objNum, genNum, method)
		return aesDecrypt(key, data)
	case cipherAES256:
		// AESV3 uses the file key directly, with no per-object derivation.
		return aesDecrypt(h.FileKey, data)
	default:
		return nil, fmt.Errorf("%w: unknown crypt method", ErrEncrypted)
	}
}

// rc4Apply runs RC4 over data with the given key. RC4 is symmetric.
func rc4Apply(key, data []byte) []byte {
	//nolint:gosec // G401: RC4 required by PDF spec §7.6.3.4
	c, _ := rc4.NewCipher(key) // only errors on empty key
	out := make([]byte, len(data))
	c.XORKeyStream(out, data)
	return out
}

// aesDecrypt decrypts AES-128-CBC ciphertext with a 16-byte IV prefix and
// strips PKCS#7 padding.
func aesDecrypt(key, data []byte) ([]byte, error) {
	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("%w: AES ciphertext shorter than IV", ErrEncrypted)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%w: aes.NewCipher: %v", ErrEncrypted, err)
	}
	iv := data[:aes.BlockSize]
	ct := data[aes.BlockSize:]
	if len(ct)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("%w: AES ciphertext not block-aligned", ErrEncrypted)
	}
	out := make([]byte, len(ct))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(out, ct)
	return stripPKCS7(out, aes.BlockSize)
}

func stripPKCS7(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}
	pad := int(data[len(data)-1])
	if pad == 0 || pad > blockSize || pad > len(data) {
		return nil, fmt.Errorf("%w: invalid PKCS#7 padding", ErrEncrypted)
	}
	for _, b := range data[len(data)-pad:] {
		if int(b) != pad {
			return nil, fmt.Errorf("%w: invalid PKCS#7 padding bytes", ErrEncrypted)
		}
	}
	return data[:len(data)-pad], nil
}

// stringBytes extracts raw bytes from a StringLiteral or HexString.
// A HexString holds its hex digits as written, so it is decoded here.
func stringBytes(o Object) ([]byte, error) {
	switch s := o.(type) {
	case StringLiteral:
		return []byte(s), nil
	case HexString:
		return hexBytes(s)
	case nil:
		return nil, fmt.Errorf("nil")
	default:
		return nil, fmt.Errorf("not a string: %T", o)
	}
}

// intFromDict pulls an int out of a dictionary, accepting both Integer and Real.
func intFromDict(d Dictionary, key string) (int, bool) {
	switch v := d[Name(key)].(type) {
	case Integer:
		return int(v), true
	case Real:
		return int(v), true
	}
	return 0, false
}

// hexBytes decodes the digits of a HexString; an odd digit count is padded
// with a trailing 0 (ISO 32000-1 §7.3.4.3).
func hexBytes(s HexString) ([]byte, error) {
	digits := string(s)
	if len(digits)%2 == 1 {
		digits += "0"
	}
	b, err := hex.DecodeString(digits)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid hex string: %v", ErrEncrypted, err)
	}
	return b, nil
}
