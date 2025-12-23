package pdf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

// DecodeStream decodes a stream based on its filter and optional parameters
func DecodeStream(data []byte, filter Name, params Dictionary) ([]byte, error) {
	var decoded []byte
	var err error

	switch filter {
	case "FlateDecode":
		decoded, err = decodeFlate(data)
	case "DCTDecode":
		// DCTDecode is usually JPEG, which is already compressed.
		// We return it as is, and the consumer should handle it as JPEG data.
		return data, nil
	case "ASCIIHexDecode":
		decoded, err = decodeASCIIHex(data)
	case "ASCII85Decode":
		decoded, err = decodeASCII85(data)
	case "LZWDecode":
		decoded, err = decodeLZW(data, params)
	default:
		// Unknown filter or no filter
		return data, nil
	}

	if err != nil {
		return nil, err
	}

	// Apply Predictor if present
	if params != nil {
		if predictorObj, ok := params[Name("Predictor")]; ok {
			var predictor int
			if p, ok := predictorObj.(Integer); ok {
				predictor = int(p)
			} else if p, ok := predictorObj.(Real); ok {
				predictor = int(p)
			}

			if predictor >= 10 {
				columns := 1
				if colsObj, ok := params[Name("Columns")]; ok {
					if c, ok := colsObj.(Integer); ok {
						columns = int(c)
					}
				}

				// Colors (components per pixel) - default 1
				colors := 1
				if colorsObj, ok := params[Name("Colors")]; ok {
					if c, ok := colorsObj.(Integer); ok {
						colors = int(c)
					}
				}

				// BitsPerComponent - default 8
				bpc := 8
				if bpcObj, ok := params[Name("BitsPerComponent")]; ok {
					if c, ok := bpcObj.(Integer); ok {
						bpc = int(c)
					}
				}

				return applyPredictor(decoded, predictor, columns, colors, bpc)
			}
		}
	}

	return decoded, nil
}

func decodeFlate(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}

func applyPredictor(data []byte, predictor, columns, colors, bpc int) ([]byte, error) {
	if predictor < 10 {
		return data, nil // Only PNG predictors supported for now
	}

	// PNG Predictors (10-15)
	// Rows are (Columns * Colors * BitsPerComponent / 8) bytes wide + 1 byte for filter type
	// But usually bpc=8, so bytesPerPixel = Colors.

	bytesPerPixel := (colors*bpc + 7) / 8
	rowLen := (columns*colors*bpc + 7) / 8
	stride := rowLen + 1 // +1 for filter byte

	if len(data)%stride != 0 {
		// It might be that the last row is incomplete or something, but standard says it should match.
		// However, let's be robust.
	}

	rows := len(data) / stride
	if rows == 0 {
		return data, nil
	}

	out := make([]byte, rows*rowLen)

	// Previous row (initialized to zero)
	prev := make([]byte, rowLen)

	for i := 0; i < rows; i++ {
		start := i * stride
		if start+stride > len(data) {
			break
		}

		row := data[start : start+stride]
		filter := row[0]
		raw := row[1:] // The data bytes

		destStart := i * rowLen
		dest := out[destStart : destStart+rowLen]

		switch filter {
		case 0: // None
			copy(dest, raw)
		case 1: // Sub
			for j := 0; j < rowLen; j++ {
				left := byte(0)
				if j >= bytesPerPixel {
					left = dest[j-bytesPerPixel]
				}
				dest[j] = raw[j] + left
			}
		case 2: // Up
			for j := 0; j < rowLen; j++ {
				up := prev[j]
				dest[j] = raw[j] + up
			}
		case 3: // Average
			for j := 0; j < rowLen; j++ {
				left := byte(0)
				if j >= bytesPerPixel {
					left = dest[j-bytesPerPixel]
				}
				up := prev[j]
				dest[j] = raw[j] + byte((int(left)+int(up))/2)
			}
		case 4: // Paeth
			for j := 0; j < rowLen; j++ {
				left := byte(0)
				if j >= bytesPerPixel {
					left = dest[j-bytesPerPixel]
				}
				up := prev[j]
				upLeft := byte(0)
				if j >= bytesPerPixel {
					upLeft = prev[j-bytesPerPixel]
				}

				dest[j] = raw[j] + paethPredictor(left, up, upLeft)
			}
		default:
			return nil, fmt.Errorf("unknown PNG filter %d", filter)
		}

		// Update prev for next iteration
		copy(prev, dest)
	}

	return out, nil
}

func paethPredictor(a, b, c byte) byte {
	p := int(a) + int(b) - int(c)
	pa := abs(p - int(a))
	pb := abs(p - int(b))
	pc := abs(p - int(c))

	if pa <= pb && pa <= pc {
		return a
	} else if pb <= pc {
		return b
	}
	return c
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// decodeASCIIHex decodes ASCIIHexDecode filter
// Each pair of hexadecimal digits encodes one byte. EOD marker is '>'.
func decodeASCIIHex(data []byte) ([]byte, error) {
	var out bytes.Buffer
	var nibble byte
	haveNibble := false

	for _, b := range data {
		// Skip whitespace
		if isASCIIWhitespace(b) {
			continue
		}
		// EOD marker
		if b == '>' {
			break
		}

		var val byte
		switch {
		case b >= '0' && b <= '9':
			val = b - '0'
		case b >= 'A' && b <= 'F':
			val = b - 'A' + 10
		case b >= 'a' && b <= 'f':
			val = b - 'a' + 10
		default:
			return nil, fmt.Errorf("invalid hex character: %c", b)
		}

		if haveNibble {
			out.WriteByte(nibble<<4 | val)
			haveNibble = false
		} else {
			nibble = val
			haveNibble = true
		}
	}

	// If odd number of hex digits, treat last nibble as if followed by 0
	if haveNibble {
		out.WriteByte(nibble << 4)
	}

	return out.Bytes(), nil
}

// decodeASCII85 decodes ASCII85Decode (btoa) filter
// Encodes 4 bytes as 5 ASCII characters (base 85). 'z' represents 4 zero bytes.
// EOD marker is '~>'.
func decodeASCII85(data []byte) ([]byte, error) {
	var out bytes.Buffer
	var tuple [5]byte
	tupleIndex := 0

	for i := 0; i < len(data); i++ {
		b := data[i]

		// Skip whitespace
		if isASCIIWhitespace(b) {
			continue
		}

		// Check for EOD marker
		if b == '~' {
			if i+1 < len(data) && data[i+1] == '>' {
				break
			}
		}

		// 'z' is special: represents 4 zero bytes
		if b == 'z' {
			if tupleIndex != 0 {
				return nil, fmt.Errorf("'z' inside group in ASCII85")
			}
			out.Write([]byte{0, 0, 0, 0})
			continue
		}

		// Valid ASCII85 characters are '!' (33) to 'u' (117)
		if b < '!' || b > 'u' {
			return nil, fmt.Errorf("invalid ASCII85 character: %c (%d)", b, b)
		}

		tuple[tupleIndex] = b - '!'
		tupleIndex++

		if tupleIndex == 5 {
			// Decode 5 characters to 4 bytes
			val := uint32(tuple[0])*85*85*85*85 +
				uint32(tuple[1])*85*85*85 +
				uint32(tuple[2])*85*85 +
				uint32(tuple[3])*85 +
				uint32(tuple[4])

			out.WriteByte(byte(val >> 24))
			out.WriteByte(byte(val >> 16))
			out.WriteByte(byte(val >> 8))
			out.WriteByte(byte(val))
			tupleIndex = 0
		}
	}

	// Handle remaining bytes (partial group)
	if tupleIndex > 0 {
		// Pad with 'u' (84) to make 5 characters
		for i := tupleIndex; i < 5; i++ {
			tuple[i] = 84 // 'u' - '!' = 84
		}

		val := uint32(tuple[0])*85*85*85*85 +
			uint32(tuple[1])*85*85*85 +
			uint32(tuple[2])*85*85 +
			uint32(tuple[3])*85 +
			uint32(tuple[4])

		// Output only the bytes we actually have
		// tupleIndex chars encode tupleIndex-1 bytes
		numBytes := tupleIndex - 1
		if numBytes >= 1 {
			out.WriteByte(byte(val >> 24))
		}
		if numBytes >= 2 {
			out.WriteByte(byte(val >> 16))
		}
		if numBytes >= 3 {
			out.WriteByte(byte(val >> 8))
		}
	}

	return out.Bytes(), nil
}

// decodeLZW decodes LZWDecode filter
// PDF uses early-change variant of LZW with 8-bit initial code size
func decodeLZW(data []byte, params Dictionary) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	// Get EarlyChange parameter (default 1 in PDF)
	earlyChange := 1
	if params != nil {
		if ec, ok := params[Name("EarlyChange")]; ok {
			if v, ok := ec.(Integer); ok {
				earlyChange = int(v)
			}
		}
	}

	decoder := newLZWDecoder(data, earlyChange)
	decoded, err := decoder.decode()
	if err != nil {
		return nil, err
	}

	// Apply predictor if present (same as FlateDecode)
	if params != nil {
		if predictorObj, ok := params[Name("Predictor")]; ok {
			var predictor int
			if p, ok := predictorObj.(Integer); ok {
				predictor = int(p)
			} else if p, ok := predictorObj.(Real); ok {
				predictor = int(p)
			}

			if predictor >= 10 {
				columns := 1
				if colsObj, ok := params[Name("Columns")]; ok {
					if c, ok := colsObj.(Integer); ok {
						columns = int(c)
					}
				}
				colors := 1
				if colorsObj, ok := params[Name("Colors")]; ok {
					if c, ok := colorsObj.(Integer); ok {
						colors = int(c)
					}
				}
				bpc := 8
				if bpcObj, ok := params[Name("BitsPerComponent")]; ok {
					if c, ok := bpcObj.(Integer); ok {
						bpc = int(c)
					}
				}
				return applyPredictor(decoded, predictor, columns, colors, bpc)
			}
		}
	}

	return decoded, nil
}

// lzwDecoder handles LZW decompression
type lzwDecoder struct {
	data        []byte
	pos         int
	bitPos      int
	codeSize    int
	nextCode    int
	earlyChange int
	table       [][]byte
}

const (
	lzwClearCode = 256
	lzwEODCode   = 257
	lzwFirstCode = 258
	lzwMaxCode   = 4095
)

func newLZWDecoder(data []byte, earlyChange int) *lzwDecoder {
	d := &lzwDecoder{
		data:        data,
		earlyChange: earlyChange,
	}
	d.reset()
	return d
}

func (d *lzwDecoder) reset() {
	d.codeSize = 9
	d.nextCode = lzwFirstCode
	// Initialize table with single-byte sequences
	d.table = make([][]byte, lzwMaxCode+1)
	for i := 0; i < 256; i++ {
		d.table[i] = []byte{byte(i)}
	}
}

func (d *lzwDecoder) readCode() (int, error) {
	code := 0
	for i := 0; i < d.codeSize; i++ {
		if d.pos >= len(d.data) {
			return 0, io.EOF
		}
		bit := (d.data[d.pos] >> (7 - d.bitPos)) & 1
		code = (code << 1) | int(bit)
		d.bitPos++
		if d.bitPos == 8 {
			d.bitPos = 0
			d.pos++
		}
	}
	return code, nil
}

func (d *lzwDecoder) decode() ([]byte, error) {
	var out bytes.Buffer
	var prevSeq []byte

	for {
		code, err := d.readCode()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if code == lzwEODCode {
			break
		}

		if code == lzwClearCode {
			d.reset()
			prevSeq = nil
			continue
		}

		var seq []byte
		if code < d.nextCode {
			// Code is in the table
			seq = d.table[code]
		} else if code == d.nextCode && prevSeq != nil {
			// Special case: code not yet in table
			seq = append([]byte{}, prevSeq...)
			seq = append(seq, prevSeq[0])
		} else {
			return nil, fmt.Errorf("invalid LZW code: %d (nextCode=%d)", code, d.nextCode)
		}

		out.Write(seq)

		// Add new entry to table
		if prevSeq != nil && d.nextCode <= lzwMaxCode {
			newSeq := append([]byte{}, prevSeq...)
			newSeq = append(newSeq, seq[0])
			d.table[d.nextCode] = newSeq
			d.nextCode++

			// Increase code size if needed
			// PDF uses earlyChange=1 by default (increase before code is used)
			threshold := 1 << d.codeSize
			if d.earlyChange == 1 {
				threshold--
			}
			if d.nextCode > threshold && d.codeSize < 12 {
				d.codeSize++
			}
		}

		prevSeq = seq
	}

	return out.Bytes(), nil
}

func isASCIIWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\f' || b == 0
}

// DecodeStreamFromDict decodes stream data using filter(s) specified in the dictionary.
// Handles both single filters and arrays of filters.
func DecodeStreamFromDict(data []byte, dict Dictionary) ([]byte, error) {
	filter := dict[Name("Filter")]
	if filter == nil {
		return data, nil
	}

	// Get DecodeParms (can be single dictionary or array of dictionaries)
	decodeParmsObj := dict[Name("DecodeParms")]
	var decodeParmsArr []Dictionary

	if arr, ok := decodeParmsObj.(Array); ok {
		for _, item := range arr {
			if d, ok := item.(Dictionary); ok {
				decodeParmsArr = append(decodeParmsArr, d)
			} else {
				decodeParmsArr = append(decodeParmsArr, nil) // null entry
			}
		}
	} else if d, ok := decodeParmsObj.(Dictionary); ok {
		decodeParmsArr = []Dictionary{d}
	}

	// Handle single filter (Name)
	if name, ok := filter.(Name); ok {
		var parms Dictionary
		if len(decodeParmsArr) > 0 {
			parms = decodeParmsArr[0]
		}
		return DecodeStream(data, name, parms)
	}

	// Handle Array of filters - apply in order
	if arr, ok := filter.(Array); ok {
		var err error
		for i, f := range arr {
			filterName, ok := f.(Name)
			if !ok {
				return nil, fmt.Errorf("filter array element %d is not a Name", i)
			}

			var parms Dictionary
			if i < len(decodeParmsArr) {
				parms = decodeParmsArr[i]
			}

			data, err = DecodeStream(data, filterName, parms)
			if err != nil {
				return nil, fmt.Errorf("filter %s (index %d) failed: %w", filterName, i, err)
			}
		}
		return data, nil
	}

	return data, nil // No filter or unknown filter type - return raw
}
