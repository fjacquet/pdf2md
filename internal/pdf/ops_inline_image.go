package pdf

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fjacquet/pdf2md/internal/types"
)

// handleInlineImage handles BI (begin inline image) operator
// Inline images are embedded directly in content streams: BI <dict> ID <data> EI
func (in *Interpreter) handleInlineImage(tokenizer *Tokenizer) error {
	// Parse the inline image dictionary (key/value pairs until ID)
	dict := make(Dictionary)

	for {
		tok, err := tokenizer.NextToken()
		if err != nil {
			return fmt.Errorf("unexpected end of inline image dict: %v", err)
		}

		// ID marks end of dictionary and start of image data
		if tok.Type == TokenKeyword && tok.Value == "ID" {
			break
		}

		// Parse key (should be a Name, but inline images use abbreviations without /)
		key := tok.Value

		// Expand abbreviations to full names
		fullKey := expandInlineImageKey(key)

		// Parse value
		valTok, err := tokenizer.NextToken()
		if err != nil {
			return fmt.Errorf("unexpected end of inline image value: %v", err)
		}

		val := parseInlineImageValue(valTok)
		dict[Name(fullKey)] = val
	}

	// Skip one byte of whitespace after ID
	// The tokenizer should handle this, but we need to read the raw data

	// Calculate expected data length if possible
	width := getIntFromDict(dict, "Width", "W")
	height := getIntFromDict(dict, "Height", "H")
	bpc := getIntFromDict(dict, "BitsPerComponent", "BPC")
	if bpc == 0 {
		bpc = 8 // Default
	}

	// ColorSpace affects bytes per pixel
	cs := getNameFromDict(dict, "ColorSpace", "CS")
	components := colorSpaceComponents(cs)

	expectedLen := 0
	if width > 0 && height > 0 {
		// Calculate expected raw data length
		bitsPerRow := width * components * bpc
		bytesPerRow := (bitsPerRow + 7) / 8
		expectedLen = bytesPerRow * height
	}

	// Read image data until EI
	// EI must be preceded by whitespace
	data, err := readInlineImageData(tokenizer, expectedLen)
	if err != nil {
		return fmt.Errorf("error reading inline image data: %v", err)
	}

	// Decode data based on filter
	filter := getNameFromDict(dict, "Filter", "F")
	if filter != "" {
		var decodeParms Dictionary
		if dp, ok := dict[Name("DecodeParms")].(Dictionary); ok {
			decodeParms = dp
		} else if dp, ok := dict[Name("DP")].(Dictionary); ok {
			decodeParms = dp
		}
		decoded, err := DecodeStream(data, Name(filter), decodeParms)
		if err != nil {
			fmt.Printf("Warning: Failed to decode inline image: %v\n", err)
			// Continue with raw data
		} else {
			data = decoded
		}
	}

	// Get position from current CTM
	x, y := in.State.CTM.Transform(0, 0)
	ctmWidth := in.State.CTM[0]
	ctmHeight := in.State.CTM[3]

	// Generate unique ID
	imageID := fmt.Sprintf("inline_%d", len(in.Images))

	// Skip very small images (1x1, 2x2) which are typically spacers or placeholders
	if width <= 2 && height <= 2 {
		return nil
	}

	// Check if it's JPEG
	isJPEG := filter == "DCTDecode" || filter == "DCT"

	if isJPEG {
		// JPEG data is already properly encoded
		in.Images = append(in.Images, types.Image{
			ID:     imageID,
			Data:   data,
			Format: "jpeg",
			X:      x,
			Y:      y,
			Width:  ctmWidth,
			Height: ctmHeight,
		})
		return nil
	}

	// For non-JPEG, encode raw pixel data as PNG
	pngData, err := encodePNG(data, width, height, components, bpc)
	if err != nil {
		// If encoding fails, skip this image
		fmt.Printf("Warning: Failed to encode inline image as PNG: %v\n", err)
		return nil
	}

	in.Images = append(in.Images, types.Image{
		ID:     imageID,
		Data:   pngData,
		Format: "png",
		X:      x,
		Y:      y,
		Width:  ctmWidth,
		Height: ctmHeight,
	})

	return nil
}

// expandInlineImageKey expands abbreviated keys to full names
func expandInlineImageKey(abbrev string) string {
	switch abbrev {
	case "BPC":
		return "BitsPerComponent"
	case "CS":
		return "ColorSpace"
	case "D":
		return "Decode"
	case "DP":
		return "DecodeParms"
	case "F":
		return "Filter"
	case "H":
		return "Height"
	case "IM":
		return "ImageMask"
	case "I":
		return "Interpolate"
	case "W":
		return "Width"
	default:
		return abbrev
	}
}

// expandInlineImageValue expands abbreviated filter values
func expandInlineImageValue(abbrev string) string {
	switch abbrev {
	case "AHx":
		return "ASCIIHexDecode"
	case "A85":
		return "ASCII85Decode"
	case "LZW":
		return "LZWDecode"
	case "Fl":
		return "FlateDecode"
	case "RL":
		return "RunLengthDecode"
	case "CCF":
		return "CCITTFaxDecode"
	case "DCT":
		return "DCTDecode"
	case "G":
		return "DeviceGray"
	case "RGB":
		return "DeviceRGB"
	case "CMYK":
		return "DeviceCMYK"
	case "I":
		return "Indexed"
	default:
		return abbrev
	}
}

// parseInlineImageValue parses a token into an appropriate Object
func parseInlineImageValue(tok Token) Object {
	switch tok.Type {
	case TokenNumeric:
		// Check if integer or real
		if strings.Contains(tok.Value, ".") {
			var f float64
			_, _ = fmt.Sscanf(tok.Value, "%f", &f)
			return Real(f)
		}
		var i int
		_, _ = fmt.Sscanf(tok.Value, "%d", &i)
		return Integer(i)
	case TokenName:
		// Expand abbreviations for filter/colorspace values
		return Name(expandInlineImageValue(tok.Value))
	case TokenKeyword:
		// Keywords like true/false
		switch tok.Value {
		case "true":
			return Boolean(true)
		case "false":
			return Boolean(false)
		default:
			// Abbreviated name without /
			return Name(expandInlineImageValue(tok.Value))
		}
	case TokenArrayStart:
		// Would need to parse full array, but inline images rarely use arrays
		return nil
	default:
		return Name(tok.Value)
	}
}

// getIntFromDict gets an integer from dict with primary and alternate keys
func getIntFromDict(dict Dictionary, key1, key2 string) int {
	if v, ok := dict[Name(key1)]; ok {
		if i, ok := v.(Integer); ok {
			return int(i)
		}
	}
	if v, ok := dict[Name(key2)]; ok {
		if i, ok := v.(Integer); ok {
			return int(i)
		}
	}
	return 0
}

// getNameFromDict gets a name from dict with primary and alternate keys
func getNameFromDict(dict Dictionary, key1, key2 string) string {
	if v, ok := dict[Name(key1)]; ok {
		if n, ok := v.(Name); ok {
			return string(n)
		}
	}
	if v, ok := dict[Name(key2)]; ok {
		if n, ok := v.(Name); ok {
			return string(n)
		}
	}
	return ""
}

// colorSpaceComponents returns the number of components for a color space
func colorSpaceComponents(cs string) int {
	switch cs {
	case "DeviceGray", "G":
		return 1
	case "DeviceRGB", "RGB":
		return 3
	case "DeviceCMYK", "CMYK":
		return 4
	case "Indexed", "I":
		return 1 // Indexed returns single index
	default:
		return 3 // Default to RGB
	}
}

// readInlineImageData reads the raw image data until EI marker
func readInlineImageData(tokenizer *Tokenizer, _ int) ([]byte, error) {
	// Get the underlying reader to read raw bytes
	// The ID keyword has been consumed, and there should be one whitespace char
	// followed by the image data, then whitespace + EI

	var buf bytes.Buffer

	// We need to read bytes until we find the pattern: whitespace + "EI" + (whitespace or EOF)
	// This is tricky because EI could appear in the image data

	// Strategy: Read one byte at a time, looking for potential EI marker
	// EI must be preceded by whitespace in valid PDF

	// Get access to raw bytes - this requires tokenizer support
	// For now, we'll use a simpler approach: read until we find EI preceded by whitespace

	rawData, err := tokenizer.ReadUntilEI()
	if err != nil {
		return nil, err
	}

	buf.Write(rawData)
	return buf.Bytes(), nil
}
