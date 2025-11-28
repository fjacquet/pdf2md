package pdf

import (
	"fmt"
	"strconv"
)

// Parser parses PDF objects from a tokenizer
type Parser struct {
	tokenizer *Tokenizer
}

// NewParser creates a new parser
func NewParser(t *Tokenizer) *Parser {
	return &Parser{tokenizer: t}
}

// ParseObject parses the next object from the stream
func (p *Parser) ParseObject() (Object, error) {
	token, err := p.tokenizer.NextToken()
	if err != nil {
		return nil, err
	}

	switch token.Type {
	case TokenNumeric:
		// Check if it's an indirect reference (Obj Gen R)
		// 1. Parse first number (Obj)
		objNum, err := strconv.ParseInt(token.Value, 10, 64)
		if err != nil {
			// Real number
			f, _ := strconv.ParseFloat(token.Value, 64)
			return Real(f), nil
		}

		// 2. Peek next token
		t2, err := p.tokenizer.NextToken()
		if err != nil {
			// EOF or error, just return integer
			return Integer(objNum), nil
		}

		if t2.Type == TokenNumeric {
			// 3. Parse second number (Gen)
			genNum, err := strconv.ParseInt(t2.Value, 10, 64)
			if err == nil {
				// 4. Peek third token
				t3, err := p.tokenizer.NextToken()
				if err == nil {
					if t3.Type == TokenKeyword && t3.Value == "R" {
						// It IS an indirect reference!
						return IndirectRef{ObjectNumber: int(objNum), GenerationNumber: int(genNum)}, nil
					}
					// Not 'R', push back t3
					p.tokenizer.UnreadToken(t3)
				}
			}
		}

		// Not an indirect ref, push back t2
		p.tokenizer.UnreadToken(t2)

		return Integer(objNum), nil

	case TokenName:
		return Name(token.Value), nil

	case TokenString:
		return StringLiteral(token.Value), nil

	case TokenHexString:
		return HexString(token.Value), nil

	case TokenKeyword:
		switch token.Value {
		case "true":
			return Boolean(true), nil
		case "false":
			return Boolean(false), nil
		case "null":
			return Null{}, nil
		case "R":
			// Should not happen if we handle IndirectRef correctly
			return nil, fmt.Errorf("unexpected keyword 'R'")
		default:
			// Could be "obj", "endobj", "stream", etc.
			// Treat as keyword? Or error?
			// For ParseObject, we expect a value.
			return nil, fmt.Errorf("unexpected keyword '%s'", token.Value)
		}

	case TokenArrayStart:
		return p.parseArray()

	case TokenDictStart:
		return p.parseDictionary()

	default:
		return nil, fmt.Errorf("unexpected token type %v (%s)", token.Type, token.Value)
	}
}

func (p *Parser) parseArray() (Array, error) {
	var arr Array
	for {
		// Peek next token to check for end of array
		token, err := p.tokenizer.NextToken()
		if err != nil {
			return nil, err
		}

		if token.Type == TokenArrayEnd {
			break
		}

		// Not end, push back and parse object
		p.tokenizer.UnreadToken(token)

		obj, err := p.ParseObject()
		if err != nil {
			return nil, err
		}
		arr = append(arr, obj)
	}
	return arr, nil
}

func (p *Parser) parseDictionary() (Object, error) {
	dict := make(Dictionary)
	for {
		token, err := p.tokenizer.NextToken()
		if err != nil {
			return nil, err
		}

		if token.Type == TokenDictEnd {
			break
		}

		if token.Type != TokenName {
			return nil, fmt.Errorf("expected Name in dictionary, got %v", token)
		}
		key := Name(token.Value)

		val, err := p.ParseObject()
		if err != nil {
			return nil, err
		}

		dict[key] = val
	}

	// Check if followed by "stream"
	// We need to peek next token.
	// Tokenizer doesn't have Peek, but we can read and unread.

	tok, err := p.tokenizer.NextToken()
	if err != nil {
		// EOF or error, just return dict
		return dict, nil
	}

	if tok.Type == TokenKeyword && tok.Value == "stream" {
		// It is a stream!
		// We need to read stream data.
		// The length should be in dict["Length"].
		lengthObj, ok := dict[Name("Length")]
		var length int64
		if ok {
			if i, ok := lengthObj.(Integer); ok {
				length = int64(i)
			} else if _, ok := lengthObj.(IndirectRef); ok {
				// Length is indirect. We can't resolve it here easily without Reader access.
				// This is a limitation of separating Parser and Reader.
				// However, we can scan for "endstream".
				length = -1
			}
		} else {
			length = -1
		}

		data, err := p.tokenizer.ReadStream(length)
		if err != nil {
			return nil, err
		}

		// Expect endstream
		// ReadStream should consume up to endstream?
		// If length is known, we read length bytes. Then expect "endstream".
		// If length is -1, ReadStream scans for "endstream".

		// After ReadStream, we should be at "endstream".
		// Let's verify.
		tok, err = p.tokenizer.NextToken()
		if err != nil {
			return nil, err
		}
		if tok.Value != "endstream" {
			return nil, fmt.Errorf("expected 'endstream', got %v", tok)
		}

		return Stream{Dictionary: dict, Data: data}, nil
	}

	// Not a stream, push back
	p.tokenizer.UnreadToken(tok)

	return dict, nil
}

func (p *Parser) parseObjectFromToken(token Token) (Object, error) {
	// Helper to parse object given the first token
	switch token.Type {
	case TokenNumeric:
		if i, err := strconv.ParseInt(token.Value, 10, 64); err == nil {
			return Integer(i), nil
		}
		f, _ := strconv.ParseFloat(token.Value, 64)
		return Real(f), nil
	case TokenName:
		return Name(token.Value), nil
	case TokenString:
		return StringLiteral(token.Value), nil
	case TokenHexString:
		return HexString(token.Value), nil
	case TokenKeyword:
		switch token.Value {
		case "true":
			return Boolean(true), nil
		case "false":
			return Boolean(false), nil
		case "null":
			return Null{}, nil
		}
	case TokenArrayStart:
		return p.parseArray()
	case TokenDictStart:
		return p.parseDictionary()
	}
	return nil, fmt.Errorf("unexpected token %v", token)
}
