package pdf

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"
)

// Font represents a PDF font
type Font struct {
	BaseFont  string
	Subtype   string
	ToUnicode map[int]string
	Encoding  string // For simple fonts
}

// FontManager manages fonts for a page
type FontManager struct {
	Fonts  map[Name]*Font
	reader *Reader
}

// NewFontManager creates a new font manager
func NewFontManager(r *Reader) *FontManager {
	return &FontManager{
		Fonts:  make(map[Name]*Font),
		reader: r,
	}
}

// LoadFonts loads fonts from the page resources
func (fm *FontManager) LoadFonts(resources Dictionary) error {
	fontDictObj, ok := resources[Name("Font")]
	if !ok {
		return nil // No fonts
	}

	var fontDict Dictionary
	if ref, ok := fontDictObj.(IndirectRef); ok {
		obj, err := fm.reader.ReadObject(ref.ObjectNumber)
		if err != nil {
			return err
		}
		if d, ok := obj.(Dictionary); ok {
			fontDict = d
		}
	} else if d, ok := fontDictObj.(Dictionary); ok {
		fontDict = d
	}

	if fontDict == nil {
		return nil
	}

	for name, obj := range fontDict {
		var fObj Object
		if ref, ok := obj.(IndirectRef); ok {
			o, err := fm.reader.ReadObject(ref.ObjectNumber)
			if err != nil {
				return err
			}
			fObj = o
		} else {
			fObj = obj
		}

		if fDict, ok := fObj.(Dictionary); ok {
			font, err := fm.parseFont(fDict)
			if err != nil {
				// Log error but continue?
				fmt.Printf("Error parsing font %s: %v\n", name, err)
				continue
			}
			fm.Fonts[name] = font
		}
	}
	return nil
}

func (fm *FontManager) parseFont(dict Dictionary) (*Font, error) {
	font := &Font{
		ToUnicode: make(map[int]string),
	}

	if baseFont, ok := dict[Name("BaseFont")].(Name); ok {
		font.BaseFont = string(baseFont)
	}
	if subtype, ok := dict[Name("Subtype")].(Name); ok {
		font.Subtype = string(subtype)
	}

	// Check ToUnicode CMap
	if toUnicode, ok := dict[Name("ToUnicode")]; ok {
		var stream Stream
		if ref, ok := toUnicode.(IndirectRef); ok {
			obj, err := fm.reader.ReadObject(ref.ObjectNumber)
			if err != nil {
				return nil, err
			}
			if s, ok := obj.(Stream); ok {
				stream = s
			}
		} else if s, ok := toUnicode.(Stream); ok {
			stream = s
		}

		if stream.Data != nil {
			// Decompress if needed
			data, err := decodeStream(stream)
			if err != nil {
				return nil, err
			}

			err = font.parseToUnicodeCMap(data)
			if err != nil {
				return nil, err
			}
		}
	}

	return font, nil
}

// DecodeString decodes a PDF string using the font's encoding/CMap
func (f *Font) DecodeString(s string) string {
	if len(f.ToUnicode) > 0 {
		// Use CMap
		// Problem: Is the string 1-byte or 2-byte codes?
		// Usually composite fonts use 2 bytes (CID), simple fonts 1 byte.
		// We need to know if it's a composite font.
		// Subtype Type0 is composite.

		isComposite := f.Subtype == "Type0"

		var res strings.Builder
		data := []byte(s)
		i := 0
		for i < len(data) {
			var code int
			if isComposite && i+1 < len(data) {
				// Try 2 bytes
				code = int(data[i])<<8 | int(data[i+1])
				if val, ok := f.ToUnicode[code]; ok {
					res.WriteString(val)
					i += 2
					continue
				}
				// Fallback: if not found, maybe it's ASCII?
				// Or maybe we should output the CID?
				i += 2
			} else {
				// 1 byte
				code = int(data[i])
				if val, ok := f.ToUnicode[code]; ok {
					res.WriteString(val)
				} else {
					// Fallback: assume ASCII/Latin1
					res.WriteByte(data[i])
				}
				i++
			}
		}
		return res.String()
	}

	// No CMap, assume simple encoding (WinAnsi/MacRoman)
	// TODO: Implement Encoding lookup
	return s
}

func (f *Font) parseToUnicodeCMap(data []byte) error {
	// Simple parser for CMap
	// We look for:
	// <count> beginbfchar
	// <code> <unicode>
	// ...
	// endbfchar
	//
	// <count> beginbfrange
	// <start> <end> <unicode>
	// ...
	// endbfrange

	// Tokenize the CMap
	t := NewTokenizer(bytes.NewReader(data))

	// fmt.Printf("Parsing CMap for %s\n", f.BaseFont)

	for {
		tok, err := t.NextToken()
		if err != nil {
			break
		}

		if tok.Value == "beginbfchar" {
			// We need to know count. It was the previous token.
			// But we didn't save it.
			// Actually, we can just loop until endbfchar.
			for {
				codeTok, err := t.NextToken()
				if err != nil {
					break
				}
				if codeTok.Value == "endbfchar" {
					break
				}

				uniTok, err := t.NextToken()
				if err != nil {
					break
				}

				code := parseHexStringOrInt(codeTok)
				uni := parseHexStringToString(uniTok)

				f.ToUnicode[code] = uni
			}
		} else if tok.Value == "beginbfrange" {
			for {
				startTok, err := t.NextToken()
				if err != nil {
					break
				}
				if startTok.Value == "endbfrange" {
					break
				}

				endTok, err := t.NextToken()
				if err != nil {
					break
				}

				uniTok, err := t.NextToken()
				if err != nil {
					break
				}

				start := parseHexStringOrInt(startTok)
				end := parseHexStringOrInt(endTok)

				// uniTok can be <hex> (start unicode) or [ <hex> <hex> ... ] (array)
				if uniTok.Type == TokenHexString {
					// Range mapping: start->uni, start+1->uni+1 ...
					// We need to increment the last byte of unicode string?
					// Or treat unicode as integer?
					// Unicode is usually UTF-16BE in hex string.
					uniBase := parseHexStringToUTF16(uniTok.Value)
					if len(uniBase) == 1 {
						base := int(uniBase[0])
						for c := start; c <= end; c++ {
							f.ToUnicode[c] = string(rune(base + (c - start)))
						}
					}
				} else if uniTok.Type == TokenArrayStart {
					// Array of destination codes
					// [ <uni1> <uni2> ... ]
					// We need to parse the array.
					// Since our Tokenizer returns [ as token, we need to read until ]
					c := start
					for {
						valTok, err := t.NextToken()
						if err != nil {
							break
						}
						if valTok.Type == TokenArrayEnd {
							break
						}

						uni := parseHexStringToString(valTok)
						f.ToUnicode[c] = uni
						c++
					}
				}
			}
		}
	}
	return nil
}

func parseHexStringOrInt(t Token) int {
	if t.Type == TokenHexString {
		// <001A> -> 0x1A
		b, _ := hex.DecodeString(t.Value)
		if len(b) == 1 {
			return int(b[0])
		} else if len(b) == 2 {
			return int(b[0])<<8 | int(b[1])
		}
	}
	// Fallback
	i, _ := strconv.Atoi(t.Value)
	return i
}

func parseHexStringToString(t Token) string {
	if t.Type == TokenHexString {
		// UTF-16BE hex string
		// <0041> -> "A"
		runes := parseHexStringToUTF16(t.Value)
		return string(runes)
	}
	return t.Value
}

func parseHexStringToUTF16(s string) []rune {
	b, _ := hex.DecodeString(s)
	u16s := make([]uint16, len(b)/2)
	for i := 0; i < len(u16s); i++ {
		u16s[i] = uint16(b[2*i])<<8 | uint16(b[2*i+1])
	}
	return utf16.Decode(u16s)
}
