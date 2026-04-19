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
	BaseFont     string
	Subtype      string
	ToUnicode    map[int]string
	Encoding     map[int]string // Custom encoding map (code -> glyph name)
	Widths       []float64
	FirstChar    int
	LastChar     int
	CIDWidths    map[int]float64
	DefaultWidth float64
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
		Encoding:  make(map[int]string),
	}

	if baseFont, ok := dict[Name("BaseFont")].(Name); ok {
		font.BaseFont = string(baseFont)
	}
	if subtype, ok := dict[Name("Subtype")].(Name); ok {
		font.Subtype = string(subtype)
	}

	// Parse Widths
	if firstChar, ok := dict[Name("FirstChar")].(Integer); ok {
		font.FirstChar = int(firstChar)
	}
	if lastChar, ok := dict[Name("LastChar")].(Integer); ok {
		font.LastChar = int(lastChar)
	}
	if widthsObj, ok := dict[Name("Widths")]; ok {
		// Widths can be an indirect reference
		var widthsArray Array
		if ref, ok := widthsObj.(IndirectRef); ok {
			obj, err := fm.reader.ReadObject(ref.ObjectNumber)
			if err == nil {
				if arr, ok := obj.(Array); ok {
					widthsArray = arr
				}
			}
		} else if arr, ok := widthsObj.(Array); ok {
			widthsArray = arr
		}

		if widthsArray != nil {
			font.Widths = make([]float64, len(widthsArray))
			for i, w := range widthsArray {
				if val, ok := w.(Integer); ok {
					font.Widths[i] = float64(val)
				} else if val, ok := w.(Real); ok {
					font.Widths[i] = float64(val)
				}
			}
		}
	}

	// Parse DescendantFonts for Type0
	if descendantFonts, ok := dict[Name("DescendantFonts")].(Array); ok && len(descendantFonts) > 0 {
		// Usually only one descendant font
		var cidFontDict Dictionary
		if ref, ok := descendantFonts[0].(IndirectRef); ok {
			obj, err := fm.reader.ReadObject(ref.ObjectNumber)
			if err == nil {
				if d, ok := obj.(Dictionary); ok {
					cidFontDict = d
				}
			}
		} else if d, ok := descendantFonts[0].(Dictionary); ok {
			cidFontDict = d
		}

		if cidFontDict != nil {
			fm.parseCIDFont(font, cidFontDict)
		}
	}

	// Parse Encoding
	if encodingObj, ok := dict[Name("Encoding")]; ok {
		fm.parseEncoding(font, encodingObj)
	} else if font.Subtype == "Type1" || font.Subtype == "TrueType" {
		// Default encoding based on Subtype
		// For Type1, default is StandardEncoding if not specified (usually)
		// For TrueType, it's complicated.
		// Type3 fonts MUST have an Encoding entry per spec, but handle gracefully if missing.
		// Let's assume StandardEncoding for simple fonts if nothing else.
		for i, name := range StandardEncoding {
			if name != "" {
				font.Encoding[i] = name
			}
		}
	}

	// Handle Type3 specific parsing
	if font.Subtype == "Type3" {
		fm.parseType3Font(font, dict)
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

			font.parseToUnicodeCMap(data)
		}
	}

	return font, nil
}

func (fm *FontManager) parseEncoding(font *Font, encodingObj Object) {
	// Encoding can be Name or Dictionary
	var encDict Dictionary
	var baseEncoding []string

	if ref, ok := encodingObj.(IndirectRef); ok {
		obj, err := fm.reader.ReadObject(ref.ObjectNumber)
		if err == nil {
			encodingObj = obj
		}
	}

	if name, ok := encodingObj.(Name); ok {
		switch string(name) {
		case "WinAnsiEncoding":
			baseEncoding = WinAnsiEncoding
		case "MacRomanEncoding":
			baseEncoding = MacRomanEncoding
		case "StandardEncoding":
			baseEncoding = StandardEncoding
		}
	} else if d, ok := encodingObj.(Dictionary); ok {
		encDict = d
		if baseEnc, ok := d[Name("BaseEncoding")].(Name); ok {
			switch string(baseEnc) {
			case "WinAnsiEncoding":
				baseEncoding = WinAnsiEncoding
			case "MacRomanEncoding":
				baseEncoding = MacRomanEncoding
			case "StandardEncoding":
				baseEncoding = StandardEncoding
			}
		}
	}

	// Apply base encoding
	for i, name := range baseEncoding {
		if name != "" {
			font.Encoding[i] = name
		}
	}

	// Apply Differences
	if encDict != nil {
		if diffs, ok := encDict[Name("Differences")].(Array); ok {
			var currentCode int
			for _, item := range diffs {
				if code, ok := item.(Integer); ok {
					currentCode = int(code)
				} else if name, ok := item.(Name); ok {
					font.Encoding[currentCode] = string(name)
					currentCode++
				}
			}
		}
	}
}

// parseType3Font handles Type3 font specific parsing.
// Type3 fonts are user-defined fonts where each glyph is defined by a content stream.
// They use an Encoding dictionary to map character codes to glyph names,
// and a CharProcs dictionary to map glyph names to content streams.
func (fm *FontManager) parseType3Font(font *Font, dict Dictionary) {
	// Type3 fonts MUST have:
	// - FontBBox: bounding box
	// - FontMatrix: transformation matrix (usually [0.001 0 0 0.001 0 0])
	// - CharProcs: dictionary mapping character names to content streams
	// - Encoding: maps character codes to glyph names
	// - FirstChar, LastChar, Widths: like simple fonts

	// Parse FontMatrix (used for proper sizing)
	if matrixArr, ok := dict[Name("FontMatrix")].(Array); ok && len(matrixArr) == 6 {
		// FontMatrix is typically [0.001 0 0 0.001 0 0] for 1000 unit em square
		// We don't need to store this for basic text extraction
		_ = matrixArr
	}

	// Parse CharProcs to know which glyphs are defined
	// We don't need to execute the content streams for text extraction,
	// we just need to know the glyph names exist
	if charProcsObj, ok := dict[Name("CharProcs")]; ok {
		var charProcsDict Dictionary
		if ref, ok := charProcsObj.(IndirectRef); ok {
			obj, _ := fm.reader.ReadObject(ref.ObjectNumber)
			if d, ok := obj.(Dictionary); ok {
				charProcsDict = d
			}
		} else if d, ok := charProcsObj.(Dictionary); ok {
			charProcsDict = d
		}

		// CharProcs dictionary has Name -> Stream mappings
		// For text extraction, having the Encoding is usually sufficient
		// But we can verify glyphs exist
		if charProcsDict != nil {
			// Log or validate CharProcs entries if needed
			// For now, we rely on Encoding for character mapping
			_ = charProcsDict
		}
	}

	// For Type3 fonts without ToUnicode, we try to build a mapping from:
	// Character code -> Encoding -> Glyph name -> GlyphToUnicode
	// This is already handled in DecodeString via font.Encoding

	// If Encoding doesn't provide useful names, try to extract from CharProcs keys
	// This is a fallback for non-standard Type3 fonts
	if len(font.Encoding) == 0 {
		// Try to build encoding from CharProcs if present
		if charProcsObj, ok := dict[Name("CharProcs")]; ok {
			var charProcsDict Dictionary
			if ref, ok := charProcsObj.(IndirectRef); ok {
				obj, _ := fm.reader.ReadObject(ref.ObjectNumber)
				if d, ok := obj.(Dictionary); ok {
					charProcsDict = d
				}
			} else if d, ok := charProcsObj.(Dictionary); ok {
				charProcsDict = d
			}

			if charProcsDict != nil && font.FirstChar >= 0 && font.LastChar >= font.FirstChar {
				// Map character codes to glyph names in order
				// This is a heuristic and may not be correct for all Type3 fonts
				code := font.FirstChar
				for name := range charProcsDict {
					if code <= font.LastChar {
						font.Encoding[code] = string(name)
						code++
					}
				}
			}
		}
	}
}

func (fm *FontManager) parseCIDFont(font *Font, dict Dictionary) {
	font.CIDWidths = make(map[int]float64)
	font.DefaultWidth = 1000 // Default default width

	if dw, ok := dict[Name("DW")]; ok {
		if val, ok := dw.(Integer); ok {
			font.DefaultWidth = float64(val)
		} else if val, ok := dw.(Real); ok {
			font.DefaultWidth = float64(val)
		}
	}

	if wObj, ok := dict[Name("W")]; ok {
		var wArray Array
		if ref, ok := wObj.(IndirectRef); ok {
			obj, err := fm.reader.ReadObject(ref.ObjectNumber)
			if err == nil {
				if arr, ok := obj.(Array); ok {
					wArray = arr
				}
			}
		} else if arr, ok := wObj.(Array); ok {
			wArray = arr
		}

		if wArray != nil {
			i := 0
			for i < len(wArray) {
				// Format: c [w1 w2 ...]
				// OR: c_first c_last w

				// First element is always c (start CID)
				var startCID int
				if val, ok := wArray[i].(Integer); ok {
					startCID = int(val)
				} else {
					// Error or unexpected type
					i++
					continue
				}
				i++

				if i >= len(wArray) {
					break
				}

				// Check next element type
				if arr, ok := wArray[i].(Array); ok {
					// c [w1 w2 ...]
					for j, w := range arr {
						width := 0.0
						if val, ok := w.(Integer); ok {
							width = float64(val)
						} else if val, ok := w.(Real); ok {
							width = float64(val)
						}
						font.CIDWidths[startCID+j] = width
					}
					i++
				} else {
					// c_first c_last w
					// We already have c_first (startCID)
					// Next should be c_last
					var endCID int
					if val, ok := wArray[i].(Integer); ok {
						endCID = int(val)
					} else {
						i++
						continue
					}
					i++

					if i >= len(wArray) {
						break
					}

					// Next should be width
					var width float64
					if val, ok := wArray[i].(Integer); ok {
						width = float64(val)
					} else if val, ok := wArray[i].(Real); ok {
						width = float64(val)
					}
					i++

					for c := startCID; c <= endCID; c++ {
						font.CIDWidths[c] = width
					}
				}
			}
		}
	}
}

// GetWidth returns the width of the character code
func (f *Font) GetWidth(code int) float64 {
	if f.CIDWidths != nil {
		if w, ok := f.CIDWidths[code]; ok {
			return w
		}
		return f.DefaultWidth
	}

	if len(f.Widths) > 0 {
		if code >= f.FirstChar && code <= f.LastChar {
			idx := code - f.FirstChar
			if idx >= 0 && idx < len(f.Widths) {
				return f.Widths[idx]
			}
		}
	}
	// Default width?
	// For standard 14 fonts we should have metrics, but for now return 0 or estimate
	// If we return 0, the interpreter might use the fallback estimate.
	return 0
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
				// Try 2 bytes first
				code = int(data[i])<<8 | int(data[i+1])
				if val, ok := f.ToUnicode[code]; ok {
					res.WriteString(val)
					i += 2
					continue
				}
				// Try 1 byte if 2-byte not found (some fonts mix encodings)
				code = int(data[i])
				if val, ok := f.ToUnicode[code]; ok {
					res.WriteString(val)
					i++
					continue
				}
				// Fallback: treat as ASCII if printable
				if data[i] >= 32 && data[i] < 127 {
					res.WriteByte(data[i])
					i++
				} else {
					i += 2 // Skip 2-byte code we can't decode
				}
				continue
			}
			// 1 byte
			code = int(data[i])
			if val, ok := f.ToUnicode[code]; ok {
				res.WriteString(val)
			} else if name, ok := f.Encoding[code]; ok {
				if uni, ok := GlyphToUnicode[name]; ok {
					res.WriteString(uni)
				} else {
					// Unknown glyph name
					res.WriteByte(data[i])
				}
			} else {
				// Fallback: assume ASCII/Latin1
				res.WriteByte(data[i])
			}
			i++
		}
		return res.String()
	}

	// No CMap, use Encoding
	var res strings.Builder
	data := []byte(s)
	for i := 0; i < len(data); i++ {
		code := int(data[i])
		if name, ok := f.Encoding[code]; ok {
			if uni, ok := GlyphToUnicode[name]; ok {
				res.WriteString(uni)
			} else {
				// Unknown glyph name
				res.WriteByte(data[i])
			}
		} else {
			// Fallback: assume ASCII/Latin1
			res.WriteByte(data[i])
		}
	}
	return res.String()
}

// CalculateWidth calculates the width of the string in text space (1000 units)
func (f *Font) CalculateWidth(s string) float64 {
	var width float64
	data := []byte(s)
	i := 0
	isComposite := f.Subtype == "Type0"

	for i < len(data) {
		var code int
		if isComposite && i+1 < len(data) {
			code = int(data[i])<<8 | int(data[i+1])
			i += 2
		} else {
			code = int(data[i])
			i++
		}
		width += f.GetWidth(code)
	}
	return width
}

func (f *Font) parseToUnicodeCMap(data []byte) {
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
							f.ToUnicode[c] = string(rune(base + (c - start))) //nolint:gosec // CMap codepoint within Unicode range
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

// DecodeString decodes a string using the specified font
func (fm *FontManager) DecodeString(fontName Name, s string) string {
	if font, ok := fm.Fonts[fontName]; ok {
		return font.DecodeString(s)
	}
	return s
}

// CalculateWidth calculates the width of the string using the specified font
func (fm *FontManager) CalculateWidth(fontName Name, s string) float64 {
	if font, ok := fm.Fonts[fontName]; ok {
		return font.CalculateWidth(s)
	}
	return 0
}
