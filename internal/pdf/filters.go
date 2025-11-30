package pdf

import (
	"bytes"
	"compress/zlib"
	"io"
)

// DecodeStream decodes a stream based on its filter
func DecodeStream(data []byte, filter Name) ([]byte, error) {
	switch filter {
	case "FlateDecode":
		return decodeFlate(data)
	case "DCTDecode":
		// DCTDecode is usually JPEG, which is already compressed.
		// We return it as is, and the consumer should handle it as JPEG data.
		return data, nil
	default:
		// Unknown filter or no filter
		return data, nil
	}
}

func decodeFlate(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}
