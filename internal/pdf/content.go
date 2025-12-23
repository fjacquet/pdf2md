package pdf

import (
	"bytes"
)

// ExtractContent extracts the raw content stream from a page
func (r *Reader) ExtractContent(page Dictionary) ([]byte, error) {
	contents, ok := page[Name("Contents")]
	if !ok {
		return nil, nil // Empty page
	}

	var streams []Stream

	// Contents can be a single reference, a single stream, or an array of references/streams
	switch v := contents.(type) {
	case IndirectRef:
		obj, err := r.ReadObject(v.ObjectNumber)
		if err != nil {
			return nil, err
		}
		if s, ok := obj.(Stream); ok {
			streams = append(streams, s)
		} else if arr, ok := obj.(Array); ok {
			// Array of references (indirectly)
			for _, item := range arr {
				if ref, ok := item.(IndirectRef); ok {
					obj, err := r.ReadObject(ref.ObjectNumber)
					if err != nil {
						return nil, err
					}
					if s, ok := obj.(Stream); ok {
						streams = append(streams, s)
					}
				}
			}
		}
	case Array:
		for _, item := range v {
			if ref, ok := item.(IndirectRef); ok {
				obj, err := r.ReadObject(ref.ObjectNumber)
				if err != nil {
					return nil, err
				}
				if s, ok := obj.(Stream); ok {
					streams = append(streams, s)
				}
			}
		}
	case Stream:
		streams = append(streams, v)
	}

	var fullContent bytes.Buffer
	for _, s := range streams {
		data, err := decodeStream(s)
		if err != nil {
			return nil, err
		}
		fullContent.Write(data)
		// Add space between streams to avoid merging tokens
		fullContent.WriteByte(' ')
	}

	return fullContent.Bytes(), nil
}

func decodeStream(s Stream) ([]byte, error) {
	return DecodeStreamFromDict(s.Data, s.Dictionary)
}
