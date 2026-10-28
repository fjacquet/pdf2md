package pdf

import (
	"bytes"
	"encoding/hex"
	"log"
	"maps"
	"regexp"
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
	MissingWidth float64       // FontDescriptor /MissingWidth, 0 when absent
	Type3Scale   float64       // FontMatrix[0]×1000 for Type 3 fonts, 0 otherwise
	codespaces   []codespace   // from the embedded encoding CMap of a composite font
	std14        *std14Metrics // standard 14 metrics matching BaseFont, if any
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

// resolve follows an indirect reference; other objects are returned as is.
func (fm *FontManager) resolve(obj Object) Object {
	if ref, ok := obj.(IndirectRef); ok && fm.reader != nil {
		if o, err := fm.reader.ReadObject(ref.ObjectNumber); err == nil {
			return o
		}
	}
	return obj
}

// Child returns a font manager holding fm's fonts overlaid with the fonts of
// resources, for a form XObject that names its own fonts.
func (fm *FontManager) Child(resources Dictionary) *FontManager {
	child := &FontManager{Fonts: maps.Clone(fm.Fonts), reader: fm.reader}
	if err := child.LoadFonts(resources); err != nil {
		log.Printf("Warning: failed to load form fonts: %v", err)
	}
	return child
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
				log.Printf("Error parsing font %s: %v\n", name, err)
				continue
			}
			fm.Fonts[name] = font
		}
	}
	return nil
}

//go:generate go run gen_std14/main.go

// Code is one character code read from a string, with its length in bytes.
type Code struct{ Value, Len int }

// codespace is one begincodespacerange entry: codes of len(lo) bytes whose
// bytes lie between lo and hi, byte by byte.
type codespace struct{ lo, hi []byte }

type std14Metrics struct {
	byName map[string]int // glyph name -> width
	byCode map[int]int    // code in the font's built-in encoding -> width
}

var (
	codespaceRange = regexp.MustCompile(`(?s)begincodespacerange(.*?)endcodespacerange`)
	hexPair        = regexp.MustCompile(`<([0-9A-Fa-f]+)>\s*<([0-9A-Fa-f]+)>`)
)

// parseCodespace reads the begincodespacerange sections of a CMap.
func parseCodespace(cmap []byte) []codespace {
	var out []codespace
	for _, sec := range codespaceRange.FindAllSubmatch(cmap, -1) {
		for _, m := range hexPair.FindAllSubmatch(sec[1], -1) {
			lo, err1 := hex.DecodeString(string(m[1]))
			hi, err2 := hex.DecodeString(string(m[2]))
			if err1 == nil && err2 == nil && len(lo) == len(hi) && len(lo) > 0 {
				out = append(out, codespace{lo, hi})
			}
		}
	}
	return out
}

// std14Aliases maps common names of the standard fonts to their AFM names.
var std14Aliases = map[string]string{"Arial": "Helvetica", "TimesNewRoman": "Times-Roman", "CourierNew": "Courier"}

// lookupStd14 returns the standard 14 metrics for baseFont, without its
// subset prefix (ABCDEF+), or nil.
func lookupStd14(baseFont string) *std14Metrics {
	if i := strings.IndexByte(baseFont, '+'); i == 6 {
		baseFont = baseFont[7:]
	}
	if alias, ok := std14Aliases[baseFont]; ok {
		baseFont = alias
	}
	if m, ok := std14Fonts[baseFont]; ok {
		return &m
	}
	return nil
}

// Codes splits raw into character codes: one byte for simple fonts, the
// encoding CMap's codespace for composite fonts, else two bytes
// (Identity-H/V).
// shortcut: predefined CMaps other than Identity are read as 2-byte codes, upgrade when a corpus document uses one.
func (f *Font) Codes(raw []byte) []Code {
	var out []Code
	for i := 0; i < len(raw); {
		n := 1
		switch {
		case f == nil || f.Subtype != "Type0":
		case len(f.codespaces) > 0:
			n = f.codeLen(raw[i:])
		case i+1 < len(raw):
			n = 2
		}
		v := 0
		for _, b := range raw[i : i+n] {
			v = v<<8 | int(b)
		}
		out = append(out, Code{v, n})
		i += n
	}
	return out
}

// codeLen is the length of the codespace range matching the start of b, 1
// when none matches.
func (f *Font) codeLen(b []byte) int {
	for _, cs := range f.codespaces {
		if len(cs.lo) > len(b) {
			continue
		}
		match := true
		for k := range cs.lo {
			if b[k] < cs.lo[k] || b[k] > cs.hi[k] {
				match = false
				break
			}
		}
		if match {
			return len(cs.lo)
		}
	}
	return 1
}

// Width is the advance of c in thousandths of em of text space (ISO 32000
// §9.6–9.7): composite fonts use /W then /DW; simple fonts use /Widths,
// then /MissingWidth, then the standard 14 metrics, then 500. A nil font
// (unknown Tf name) has width 500.
// shortcut: vertical metrics (Identity-V, /W2) ignored, upgrade when a vertical CJK document is in the corpus.
func (f *Font) Width(c Code) float64 {
	if f == nil {
		return 500
	}
	if f.Subtype == "Type0" || f.CIDWidths != nil {
		if w, ok := f.CIDWidths[c.Value]; ok {
			return w
		}
		if f.DefaultWidth > 0 {
			return f.DefaultWidth
		}
		return 1000
	}
	if i := c.Value - f.FirstChar; i >= 0 && i < len(f.Widths) {
		if f.Type3Scale != 0 {
			return f.Widths[i] * f.Type3Scale
		}
		return f.Widths[i]
	}
	if f.MissingWidth > 0 {
		return f.MissingWidth
	}
	if f.std14 != nil {
		if w, ok := f.std14.byName[f.Encoding[c.Value]]; ok {
			return float64(w)
		}
		if w, ok := f.std14.byCode[c.Value]; ok {
			return float64(w)
		}
	}
	return 500
}

// SpaceWidth is the width of the font's space, in thousandths of em: the
// code mapped to " " by ToUnicode or named "space" by the encoding, else the
// standard 14 space, else 250. A width that is not in (0, 1000] also gives
// 250, so a zero-width space glyph cannot make every gap a word gap.
func (f *Font) SpaceWidth() float64 {
	w := 0.0
	switch {
	case f == nil:
	case f.spaceCode() != nil && f.declares(*f.spaceCode()):
		w = f.Width(*f.spaceCode())
	case f.std14 != nil:
		w = float64(f.std14.byName["space"])
	}
	if w <= 0 || w > 1000 {
		return 250
	}
	return w
}

// declares reports whether the font gives c its own width (/W, /Widths or
// standard 14 metrics), rather than a default such as /DW or /MissingWidth.
func (f *Font) declares(c Code) bool {
	if f.Subtype == "Type0" || f.CIDWidths != nil {
		_, ok := f.CIDWidths[c.Value]
		return ok
	}
	i := c.Value - f.FirstChar
	return (i >= 0 && i < len(f.Widths)) || f.std14 != nil
}

// spaceCode is the lowest code mapped to " " by ToUnicode, else the lowest
// code the encoding names "space" (WinAnsi also names 160 so), or nil.
func (f *Font) spaceCode() *Code {
	n := 1
	if f.Subtype == "Type0" {
		n = 2
	}
	best := -1
	for code, s := range f.ToUnicode {
		if s == " " && (best < 0 || code < best) {
			best = code
		}
	}
	if best < 0 {
		for code, name := range f.Encoding {
			if name == "space" && (best < 0 || code < best) {
				best = code
			}
		}
	}
	if best < 0 {
		return nil
	}
	return &Code{best, n}
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
	font.std14 = lookupStd14(font.BaseFont)
	if fd, ok := fm.resolve(dict[Name("FontDescriptor")]).(Dictionary); ok {
		font.MissingWidth = toFloat(fd[Name("MissingWidth")])
	}
	if m, ok := dict[Name("FontMatrix")].(Array); ok && len(m) == 6 && font.Subtype == "Type3" {
		font.Type3Scale = toFloat(m[0]) * 1000
	}
	if enc, ok := fm.resolve(dict[Name("Encoding")]).(Stream); ok && font.Subtype == "Type0" {
		if data, err := decodeStream(enc); err == nil {
			font.codespaces = parseCodespace(data)
		}
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
	var res strings.Builder
	data := []byte(s)
	if len(f.ToUnicode) == 0 {
		// No CMap, use Encoding byte by byte
		for _, b := range data {
			f.writeSimple(&res, b)
		}
		return res.String()
	}
	for _, c := range f.Codes(data) {
		switch {
		case f.ToUnicode[c.Value] != "":
			res.WriteString(f.ToUnicode[c.Value])
		case c.Len == 1:
			f.writeSimple(&res, byte(c.Value)) //nolint:gosec // G115: a 1-byte code is < 256
		default:
			// Unmapped multi-byte code: keep the bytes that are printable ASCII.
			for k := c.Len - 1; k >= 0; k-- {
				if b := byte(c.Value >> (8 * k)); b >= 32 && b < 127 { //nolint:gosec // G115: keeps one byte of the code on purpose
					res.WriteByte(b)
				}
			}
		}
	}
	return res.String()
}

// writeSimple writes a 1-byte code through the encoding's glyph name, else
// as Latin-1.
func (f *Font) writeSimple(res *strings.Builder, b byte) {
	if name, ok := f.Encoding[int(b)]; ok {
		if uni, ok := GlyphToUnicode[name]; ok {
			res.WriteString(uni)
			return
		}
	}
	res.WriteByte(b)
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
