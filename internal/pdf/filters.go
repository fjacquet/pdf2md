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
