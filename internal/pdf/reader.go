package pdf

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// Reader reads a PDF file
type Reader struct {
	f         *os.File
	size      int64
	Trailer   Dictionary
	XrefTable map[int]int64
	Root      Dictionary
}

// NewReader creates a new PDF reader
func NewReader(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}

	r := &Reader{
		f:         f,
		size:      info.Size(),
		XrefTable: make(map[int]int64),
	}

	if err := r.readTrailer(); err != nil {
		r.Close()
		return nil, err
	}

	return r, nil
}

// Close closes the file
func (r *Reader) Close() error {
	return r.f.Close()
}

func (r *Reader) readTrailer() error {
	// Read last 1KB
	bufSize := int64(1024)
	if r.size < bufSize {
		bufSize = r.size
	}

	buf := make([]byte, bufSize)
	_, err := r.f.ReadAt(buf, r.size-bufSize)
	if err != nil {
		return err
	}

	// Find %%EOF
	eofIdx := bytes.LastIndex(buf, []byte("%%EOF"))
	if eofIdx == -1 {
		return fmt.Errorf("%%EOF marker not found")
	}

	// Find startxref before %%EOF
	// Scan backwards from eofIdx
	// We expect:
	// startxref
	// <offset>
	// %%EOF

	// Let's just search for "startxref" in the buffer
	startxrefIdx := bytes.LastIndex(buf[:eofIdx], []byte("startxref"))
	if startxrefIdx == -1 {
		return fmt.Errorf("startxref marker not found")
	}

	// Read offset
	// Skip "startxref" and whitespace
	offsetStr := string(bytes.TrimSpace(buf[startxrefIdx+9 : eofIdx]))
	var xrefOffset int64
	_, err = fmt.Sscanf(offsetStr, "%d", &xrefOffset)
	if err != nil {
		return fmt.Errorf("invalid startxref offset: %v", err)
	}

	// Parse Xref
	if err := r.readXref(xrefOffset); err != nil {
		return fmt.Errorf("failed to read xref: %v", err)
	}

	// Parse Trailer Dictionary
	// The trailer dictionary is usually before startxref.
	// But wait, `readXref` usually finds the trailer dictionary if it's a standard xref table.
	// Standard:
	// xref
	// ...
	// trailer
	// << ... >>
	// startxref

	// So we need to read xref first.

	return nil
}

func (r *Reader) readXref(offset int64) error {
	_, err := r.f.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	// Create tokenizer starting at offset
	// We need a way to reset tokenizer or create new one.
	// Since Tokenizer takes io.Reader, we can pass the file (which is at offset).
	// But we need to be careful about buffering.

	t := NewTokenizer(r.f)

	// Expect "xref"
	tok, err := t.NextToken()
	if err != nil {
		return err
	}
	if tok.Value != "xref" {
		return fmt.Errorf("expected 'xref', got '%s'", tok.Value)
	}

	// Read subsections
	for {
		// Expect start_obj_num count
		tok, err := t.NextToken()
		if err != nil {
			return err
		}

		// If we hit "trailer", we are done with xref
		if tok.Value == "trailer" {
			break
		}

		startObj, err := toInt(tok)
		if err != nil {
			return fmt.Errorf("expected object number, got %v", tok)
		}

		tok, err = t.NextToken()
		if err != nil {
			return err
		}
		count, err := toInt(tok)
		if err != nil {
			return fmt.Errorf("expected count, got %v", tok)
		}

		// Read 'count' entries
		// Each entry is 20 bytes.
		// Tokenizer might have buffered some bytes.
		// This is tricky. The xref table is line-based fixed width.
		// Tokenizer skips whitespace, so it might work if we just read tokens.
		// Entry: <offset> <gen> <n|f>

		for i := 0; i < count; i++ {
			tokOffset, err := t.NextToken()
			if err != nil {
				return err
			}
			offset, err := toInt64(tokOffset)
			if err != nil {
				return err
			}

			// Gen number
			_, err = t.NextToken()
			if err != nil {
				return err
			}

			tokFlag, err := t.NextToken()
			if err != nil {
				return err
			}

			if tokFlag.Value == "n" {
				r.XrefTable[startObj+i] = offset
			}
		}
	}

	// We hit "trailer". Parse dictionary.
	p := NewParser(t)
	// Expect "<<"
	tok, err = t.NextToken()
	if err != nil {
		return err
	}
	if tok.Type != TokenDictStart {
		return fmt.Errorf("expected trailer dictionary start '<<', got %v", tok)
	}

	// We consumed '<<', but ParseDictionary expects to consume it?
	// No, ParseDictionary calls NextToken.
	// Wait, my ParseDictionary implementation:
	/*
		func (p *Parser) parseDictionary() (Dictionary, error) {
			dict := make(Dictionary)
			for {
				token, err := p.tokenizer.NextToken()
	*/
	// It expects the CONTENTS. It does NOT expect '<<' to be passed to it.
	// But `ParseObject` handles `<<` and calls `parseDictionary`.
	// Here we are calling `parseDictionary` manually?
	// No, `ParseObject` logic is:
	/*
		case TokenDictStart:
			return p.parseDictionary()
	*/
	// So `parseDictionary` assumes `<<` is ALREADY consumed.
	// Correct.

	// But wait, I just consumed `<<` with `t.NextToken()`.
	// So I can call `p.parseDictionary()`.

	// But `Parser` doesn't expose `parseDictionary`.
	// I should expose it or use `ParseObject` but I already consumed `<<`.
	// I can `UnreadToken`!

	t.UnreadToken(tok)
	obj, err := p.ParseObject()
	if err != nil {
		return fmt.Errorf("failed to parse trailer dictionary: %v", err)
	}

	dict, ok := obj.(Dictionary)
	if !ok {
		return fmt.Errorf("trailer is not a dictionary")
	}
	r.Trailer = dict

	// Get Root
	if rootRef, ok := dict[Name("Root")].(IndirectRef); ok {
		rootObj, err := r.ReadObject(rootRef.ObjectNumber)
		if err != nil {
			return fmt.Errorf("failed to read Root object: %v", err)
		}
		if rootDict, ok := rootObj.(Dictionary); ok {
			r.Root = rootDict
		} else {
			return fmt.Errorf("Root object is not a dictionary")
		}
	}

	return nil
}

// ReadObject reads an indirect object by number
func (r *Reader) ReadObject(objNum int) (Object, error) {
	offset, ok := r.XrefTable[objNum]
	if !ok {
		return nil, fmt.Errorf("object %d not found in xref", objNum)
	}

	_, err := r.f.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	t := NewTokenizer(r.f)

	// Expect: ObjNum GenNum obj
	tok, err := t.NextToken()
	if err != nil {
		return nil, err
	}
	if id, err := toInt(tok); err != nil || id != objNum {
		return nil, fmt.Errorf("expected object id %d, got %v", objNum, tok)
	}

	tok, err = t.NextToken()
	if err != nil {
		return nil, err
	}
	// Gen num - ignore for now

	tok, err = t.NextToken()
	if err != nil {
		return nil, err
	}
	if tok.Value != "obj" {
		return nil, fmt.Errorf("expected 'obj', got %v", tok)
	}

	// Parse content
	p := NewParser(t)
	obj, err := p.ParseObject()
	if err != nil {
		return nil, err
	}

	// Expect: endobj
	tok, err = t.NextToken()
	if err != nil {
		return nil, err
	}
	if tok.Value != "endobj" {
		return nil, fmt.Errorf("expected 'endobj', got %v", tok)
	}

	return obj, nil
}

// GetPageCount returns the total number of pages
func (r *Reader) GetPageCount() (int, error) {
	if r.Root == nil {
		return 0, fmt.Errorf("root not found")
	}

	pagesRef, ok := r.Root[Name("Pages")].(IndirectRef)
	if !ok {
		return 0, fmt.Errorf("/Pages not found or not a reference")
	}

	pagesObj, err := r.ReadObject(pagesRef.ObjectNumber)
	if err != nil {
		return 0, err
	}

	pagesDict, ok := pagesObj.(Dictionary)
	if !ok {
		return 0, fmt.Errorf("/Pages is not a dictionary")
	}

	countObj, ok := pagesDict[Name("Count")]
	if !ok {
		return 0, fmt.Errorf("/Count not found in /Pages")
	}

	if count, ok := countObj.(Integer); ok {
		return int(count), nil
	}

	return 0, fmt.Errorf("/Count is not an integer")
}

// GetPage returns the dictionary for the specified page (1-based index)
func (r *Reader) GetPage(pageIndex int) (Dictionary, error) {
	if r.Root == nil {
		return nil, fmt.Errorf("root not found")
	}

	pagesRef, ok := r.Root[Name("Pages")].(IndirectRef)
	if !ok {
		return nil, fmt.Errorf("/Pages not found")
	}

	// Start traversal
	return r.traversePageTree(pagesRef, &pageIndex)
}

func (r *Reader) traversePageTree(ref IndirectRef, pageIndex *int) (Dictionary, error) {
	obj, err := r.ReadObject(ref.ObjectNumber)
	if err != nil {
		return nil, err
	}

	dict, ok := obj.(Dictionary)
	if !ok {
		return nil, fmt.Errorf("page tree node is not a dictionary")
	}

	typeVal, ok := dict[Name("Type")].(Name)
	if !ok {
		return nil, fmt.Errorf("missing /Type in page tree node")
	}

	if typeVal == "Page" {
		// Found a leaf node
		(*pageIndex)--
		if *pageIndex == 0 {
			return dict, nil
		}
		return nil, nil // Not this page
	}

	if typeVal == "Pages" {
		// Intermediate node
		kids, ok := dict[Name("Kids")].(Array)
		if !ok {
			return nil, fmt.Errorf("missing /Kids in pages node")
		}

		// Optimization: Check /Count to skip subtrees
		// If we are looking for page 100, and this node has 50 pages, we can skip it.
		// But for now, let's do simple traversal.

		for _, kid := range kids {
			kidRef, ok := kid.(IndirectRef)
			if !ok {
				return nil, fmt.Errorf("kid is not a reference")
			}

			page, err := r.traversePageTree(kidRef, pageIndex)
			if err != nil {
				return nil, err
			}
			if page != nil {
				return page, nil
			}
		}
		return nil, nil
	}

	return nil, fmt.Errorf("unknown node type: %s", typeVal)
}

func toInt(t Token) (int, error) {
	if t.Type != TokenNumeric {
		return 0, fmt.Errorf("not a number")
	}
	// Handle "0000000000"
	var val int
	_, err := fmt.Sscanf(t.Value, "%d", &val)
	return val, err
}

func toInt64(t Token) (int64, error) {
	if t.Type != TokenNumeric {
		return 0, fmt.Errorf("not a number")
	}
	var val int64
	_, err := fmt.Sscanf(t.Value, "%d", &val)
	return val, err
}
