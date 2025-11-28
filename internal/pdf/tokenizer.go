package pdf

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// TokenType represents the type of a token
type TokenType int

const (
	TokenError TokenType = iota
	TokenEOF
	TokenKeyword    // obj, endobj, stream, endstream, xref, trailer, startxref, true, false, null, R
	TokenNumeric    // 123, -12.34
	TokenName       // /Name
	TokenString     // (String)
	TokenHexString  // <Hex>
	TokenArrayStart // [
	TokenArrayEnd   // ]
	TokenDictStart  // <<
	TokenDictEnd    // >>
)

// Token represents a lexical token
type Token struct {
	Type  TokenType
	Value string
}

// Tokenizer reads tokens from a reader
type Tokenizer struct {
	r          *bufio.Reader
	pushedBack []Token
}

// NewTokenizer creates a new tokenizer
func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		r: bufio.NewReader(r),
	}
}

// UnreadToken pushes a token back to the stream
func (t *Tokenizer) UnreadToken(tok Token) {
	t.pushedBack = append(t.pushedBack, tok)
}

// NextToken returns the next token from the stream
func (t *Tokenizer) NextToken() (Token, error) {
	if len(t.pushedBack) > 0 {
		tok := t.pushedBack[len(t.pushedBack)-1]
		t.pushedBack = t.pushedBack[:len(t.pushedBack)-1]
		return tok, nil
	}

	err := t.skipWhitespaceAndComments()
	if err != nil {
		return Token{Type: TokenEOF}, err
	}

	ch, err := t.readByte()
	if err != nil {
		return Token{Type: TokenEOF}, err
	}

	switch ch {
	case '[':
		return Token{Type: TokenArrayStart, Value: "["}, nil
	case ']':
		return Token{Type: TokenArrayEnd, Value: "]"}, nil
	case '(':
		return t.readString()
	case '<':
		// Could be HexString <...> or DictStart <<
		next, err := t.peekByte()
		if err == nil && next == '<' {
			t.readByte() // consume second <
			return Token{Type: TokenDictStart, Value: "<<"}, nil
		}
		return t.readHexString()
	case '>':
		// Could be DictEnd >> or just > (end of hex string - handled in readHexString)
		// But if we see it here, it might be >>
		next, err := t.peekByte()
		if err == nil && next == '>' {
			t.readByte() // consume second >
			return Token{Type: TokenDictEnd, Value: ">>"}, nil
		}
		return Token{Type: TokenError, Value: ">"}, fmt.Errorf("unexpected >")
	case '/':
		return t.readName()
	case '%':
		// Should have been handled by skipWhitespaceAndComments, but just in case
		t.unreadByte()
		t.skipWhitespaceAndComments()
		return t.NextToken()
	default:
		t.unreadByte()
		return t.readNumericOrKeyword()
	}
}

func (t *Tokenizer) skipWhitespaceAndComments() error {
	for {
		ch, err := t.readByte()
		if err != nil {
			return err
		}

		if isWhitespace(ch) {
			continue
		}

		if ch == '%' {
			// Comment, skip to EOL
			for {
				ch, err = t.readByte()
				if err != nil {
					return err
				}
				if ch == '\r' || ch == '\n' {
					break
				}
			}
			continue
		}

		t.unreadByte()
		return nil
	}
}

func (t *Tokenizer) readName() (Token, error) {
	var buf bytes.Buffer
	for {
		ch, err := t.readByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return Token{}, err
		}

		if isWhitespace(ch) || isDelimiter(ch) {
			t.unreadByte()
			break
		}
		buf.WriteByte(ch)
	}
	// TODO: Handle #xx escapes
	return Token{Type: TokenName, Value: buf.String()}, nil
}

func (t *Tokenizer) readString() (Token, error) {
	var buf bytes.Buffer
	parens := 1 // We already consumed the first '('
	escaped := false

	for {
		ch, err := t.readByte()
		if err != nil {
			return Token{}, err
		}

		if escaped {
			// TODO: Handle octal escapes \ddd
			buf.WriteByte(ch)
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			continue
		}

		if ch == '(' {
			parens++
		} else if ch == ')' {
			parens--
		}

		if parens == 0 {
			break
		}

		buf.WriteByte(ch)
	}
	return Token{Type: TokenString, Value: buf.String()}, nil
}

func (t *Tokenizer) readHexString() (Token, error) {
	var buf bytes.Buffer
	for {
		ch, err := t.readByte()
		if err != nil {
			return Token{}, err
		}

		if ch == '>' {
			break
		}
		if isWhitespace(ch) {
			continue
		}
		buf.WriteByte(ch)
	}
	return Token{Type: TokenHexString, Value: buf.String()}, nil
}

func (t *Tokenizer) readNumericOrKeyword() (Token, error) {
	var buf bytes.Buffer
	for {
		ch, err := t.readByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return Token{}, err
		}

		if isWhitespace(ch) || isDelimiter(ch) {
			t.unreadByte()
			break
		}
		buf.WriteByte(ch)
	}

	s := buf.String()
	// Check if numeric
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return Token{Type: TokenNumeric, Value: s}, nil
	}

	return Token{Type: TokenKeyword, Value: s}, nil
}

func (t *Tokenizer) readByte() (byte, error) {
	return t.r.ReadByte()
}

func (t *Tokenizer) unreadByte() error {
	return t.r.UnreadByte()
}

func (t *Tokenizer) peekByte() (byte, error) {
	bytes, err := t.r.Peek(1)
	if err != nil {
		return 0, err
	}
	return bytes[0], nil
}

func (t *Tokenizer) ReadStream(length int64) ([]byte, error) {
	// Skip EOL after "stream"
	// The stream keyword is followed by either CRLF or LF.
	ch, err := t.readByte()
	if err != nil {
		return nil, err
	}
	if ch == '\r' {
		next, err := t.peekByte()
		if err == nil && next == '\n' {
			t.readByte()
		}
	} else if ch == '\n' {
		// OK
	} else {
		// Spec says: "The keyword stream that follows the stream dictionary should be followed by an end-of-line marker consisting of either a carriage return and a line feed or just a line feed, and not by a carriage return alone."
		// But in practice, whitespace might vary.
		// If it's not EOL, maybe we should unread?
		// But "stream" must be followed by EOL.
		return nil, fmt.Errorf("expected EOL after stream")
	}

	if length >= 0 {
		// Read exact length
		data := make([]byte, length)
		_, err := io.ReadFull(t.r, data)
		return data, err
	}

	// Scan for "endstream"
	var buf bytes.Buffer
	// This is slow and naive.
	// We read byte by byte and check for "endstream".
	// Note: "endstream" must be preceded by EOL?
	// Spec: "There should be an end-of-line marker after the data and before endstream."

	// Implementation note: This is tricky because "endstream" might appear in binary data (unlikely but possible).
	// But without length, we have no choice.

	// Optimization: Read chunks.

	endMarker := []byte("endstream")
	window := make([]byte, len(endMarker))

	// Pre-fill window
	_, err = io.ReadFull(t.r, window)
	if err != nil {
		return nil, err
	}

	for {
		if bytes.Equal(window, endMarker) {
			// Found it!
			// But wait, we consumed "endstream".
			// We need to return data BEFORE "endstream".
			// And we need to make sure "endstream" is available for NextToken?
			// NextToken expects to read "endstream".
			// So we should Unread it?
			// Or just consume it and return data?
			// Parser expects to read "endstream" via NextToken.
			// So we should push "endstream" back to tokenizer?
			// Or just return data and let Parser consume "endstream".
			// But we already consumed it from `t.r`.
			// We can use `UnreadToken` if we construct a token.
			// Or we can just return data and let Parser know we consumed endstream?
			// No, Parser calls NextToken.

			// Better: Stop BEFORE "endstream".
			// But we don't know until we read it.
			// We can use `pushedBack` tokens? No, that's for tokens.
			// We can't push back bytes to `bufio.Reader` easily (only 1 byte).

			// Hack: Return data, and push a "endstream" Token to `pushedBack`.
			t.UnreadToken(Token{Type: TokenKeyword, Value: "endstream"})

			// The data in `buf` is everything read so far minus the window.
			// But wait, we need to handle the EOL before endstream.
			// The EOL is part of the stream data? No.
			// "The sequence of bytes that make up the stream data... followed by an EOL... followed by endstream"
			// So we should strip the EOL from the end of data.

			return buf.Bytes(), nil // Need to strip EOL
		}

		// Append first byte of window to buf
		buf.WriteByte(window[0])

		// Shift window
		copy(window, window[1:])

		// Read next byte
		b, err := t.readByte()
		if err != nil {
			return nil, err
		}
		window[len(window)-1] = b
	}
}

func isWhitespace(ch byte) bool {
	return ch == 0 || ch == 9 || ch == 10 || ch == 12 || ch == 13 || ch == 32
}

func isDelimiter(ch byte) bool {
	return ch == '(' || ch == ')' || ch == '<' || ch == '>' ||
		ch == '[' || ch == ']' || ch == '{' || ch == '}' ||
		ch == '/' || ch == '%'
}
