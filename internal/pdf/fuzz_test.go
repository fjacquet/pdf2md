package pdf

import (
	"bytes"
	"testing"
)

func FuzzParseObject(f *testing.F) {
	for _, s := range []string{
		"<< /Type /Page /Kids [1 0 R 2 0 R] /Count 2 >>",
		"(a \\(nested\\) string \\101) <48656C6C6F> /Na#6De 3.14 -7 true null",
		"[[[[[[[[[[]]]]]]]]]]",
		"<< /Length 5 >> stream\nhello\nendstream",
	} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(_ *testing.T, data []byte) {
		p := NewParser(NewTokenizer(bytes.NewReader(data)))
		for range 64 {
			if _, err := p.ParseObject(); err != nil {
				return
			}
		}
	})
}

var fuzzFilters = []Name{
	"FlateDecode", "ASCIIHexDecode", "ASCII85Decode", "LZWDecode",
	"RunLengthDecode", "CCITTFaxDecode", "DCTDecode",
}

func FuzzDecodeStream(f *testing.F) {
	f.Add([]byte("48656C6C6F>"), uint8(1), int64(0), int64(0))
	f.Add([]byte("87cURD]i,\"Ebo80~>"), uint8(2), int64(0), int64(0))
	f.Add([]byte{0x80, 0x0b, 0x60, 0x50, 0x22, 0x0c, 0x0c, 0x85, 0x01}, uint8(3), int64(0), int64(0))
	f.Add([]byte{2, 'a', 'b', 'c', 0xfe, 'x', 0x80}, uint8(4), int64(0), int64(0))
	f.Add([]byte{0x78, 0x9c, 0x03, 0x00, 0x00, 0x00, 0x00, 0x01}, uint8(0), int64(12), int64(4))
	f.Fuzz(func(_ *testing.T, data []byte, filter uint8, predictor, columns int64) {
		params := Dictionary{}
		if predictor != 0 {
			params["Predictor"] = Integer(predictor)
			params["Columns"] = Integer(columns)
		}
		_, _ = DecodeStream(data, fuzzFilters[int(filter)%len(fuzzFilters)], params)
	})
}
