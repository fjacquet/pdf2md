package pdf

import (
	"fmt"

	"github.com/fjacquet/pdf2md/internal/types"
)

// ExtractLinks extracts links from a page dictionary
func (r *Reader) ExtractLinks(pageDict Dictionary) ([]types.Link, error) {
	annotsObj, ok := pageDict[Name("Annots")]
	if !ok {
		return nil, nil
	}

	// Resolve Annots array
	obj, err := r.Resolve(annotsObj)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve Annots: %w", err)
	}

	annotsArr, ok := obj.(Array)
	if !ok {
		return nil, nil // Not an array
	}

	var links []types.Link

	for _, annotRef := range annotsArr {
		// Resolve Annotation
		annotObj, err := r.Resolve(annotRef)
		if err != nil {
			continue
		}

		annotDict, ok := annotObj.(Dictionary)
		if !ok {
			continue
		}

		// Check Subtype
		subtype, ok := annotDict[Name("Subtype")].(Name)
		if !ok || subtype != "Link" {
			continue
		}

		// Get Action
		actionObj, ok := annotDict[Name("A")]
		if !ok {
			continue
		}

		actionDictObj, err := r.Resolve(actionObj)
		if err != nil {
			continue
		}

		actionDict, ok := actionDictObj.(Dictionary)
		if !ok {
			continue
		}

		// Check Action Type (S)
		s, ok := actionDict[Name("S")].(Name)
		if !ok || s != "URI" {
			continue
		}

		// Get URI
		uriObj, ok := actionDict[Name("URI")]
		if !ok {
			continue
		}

		// URI is usually a string, but could be indirect (unlikely for string)
		uriResolved, err := r.Resolve(uriObj)
		if err != nil {
			continue
		}

		var uriStr string
		if s, ok := uriResolved.(StringLiteral); ok {
			uriStr = string(s)
		} else if h, ok := uriResolved.(HexString); ok {
			uriStr = string(h) // HexString is just bytes
		} else {
			continue
		}

		// Get Rect
		rectObj, ok := annotDict[Name("Rect")]
		if !ok {
			continue
		}

		rectResolved, err := r.Resolve(rectObj)
		if err != nil {
			continue
		}

		rectArr, ok := rectResolved.(Array)
		if !ok || len(rectArr) != 4 {
			continue
		}

		rect := make([]float64, 4)
		for i := 0; i < 4; i++ {
			valObj, err := r.Resolve(rectArr[i])
			if err != nil {
				continue
			}

			if val, ok := valObj.(Integer); ok {
				rect[i] = float64(val)
			} else if val, ok := valObj.(Real); ok {
				rect[i] = float64(val)
			}
		}

		links = append(links, types.Link{
			URI:  uriStr,
			Rect: rect,
		})
	}

	return links, nil
}
