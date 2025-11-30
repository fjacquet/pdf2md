package pdf

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// XrefEntry represents an entry in the cross-reference table
type XrefEntry struct {
	Type         int   // 0=free, 1=in-use, 2=compressed
	Offset       int64 // For Type 1
	Gen          int   // For Type 1
	StreamObjNum int   // For Type 2
	StreamIndex  int   // For Type 2
}

// Reader reads a PDF file
type Reader struct {
	f         *os.File
	size      int64
	Trailer   Dictionary
	XrefTable map[int]XrefEntry
	Root      Dictionary
}

// NewReader creates a new PDF reader
func NewReader(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	r := &Reader{
		f:         f,
		size:      info.Size(),
		XrefTable: make(map[int]XrefEntry),
	}

	if err := r.readTrailer(); err != nil {
		r.Close()
		return nil, fmt.Errorf("failed to read trailer: %w", err)
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
		return fmt.Errorf("failed to read file tail: %w", err)
	}

	// Find %%EOF
	eofIdx := bytes.LastIndex(buf, []byte("%%EOF"))
	if eofIdx == -1 {
		return fmt.Errorf("%%EOF marker not found")
	}

	// Find startxref before %%EOF
	startxrefIdx := bytes.LastIndex(buf[:eofIdx], []byte("startxref"))
	if startxrefIdx == -1 {
		return fmt.Errorf("startxref marker not found")
	}

	// Read offset
	offsetStr := string(bytes.TrimSpace(buf[startxrefIdx+9 : eofIdx]))
	var xrefOffset int64
	_, err = fmt.Sscanf(offsetStr, "%d", &xrefOffset)
	if err != nil {
		return fmt.Errorf("invalid startxref offset: %w", err)
	}

	// Parse Xref
	if err := r.readXref(xrefOffset); err != nil {
		return fmt.Errorf("failed to read xref: %w", err)
	}

	return nil
}

func (r *Reader) readXref(offset int64) error {
	_, err := r.f.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek to xref offset %d: %w", offset, err)
	}

	t := NewTokenizer(r.f)

	// Check for "xref" keyword or XRef stream object
	tok, err := t.NextToken()
	if err != nil {
		return fmt.Errorf("failed to read token at xref offset: %w", err)
	}

	if tok.Value == "xref" {
		// Classic XRef Table
		return r.parseXRefTable(t)
	}

	// Possible XRef Stream (starts with object number)
	if tok.Type == TokenNumeric {
		// We have ObjNum (tok)
		// Expect GenNum
		tok2, err := t.NextToken()
		if err != nil {
			return fmt.Errorf("failed to read gen num: %w", err)
		}
		if tok2.Type != TokenNumeric {
			return fmt.Errorf("expected gen num, got %v", tok2)
		}

		// Expect "obj"
		tok3, err := t.NextToken()
		if err != nil {
			return fmt.Errorf("failed to read 'obj' keyword: %w", err)
		}
		if tok3.Value != "obj" {
			return fmt.Errorf("expected 'obj', got %v", tok3)
		}

		p := NewParser(t)
		obj, err := p.ParseObject()
		if err != nil {
			return fmt.Errorf("failed to parse object at xref offset: %w", err)
		}

		stream, ok := obj.(Stream)
		if !ok {
			return fmt.Errorf("expected XRef stream, got %T", obj)
		}

		// Check Type
		if typeName, ok := stream.Dictionary[Name("Type")].(Name); !ok || typeName != "XRef" {
			return fmt.Errorf("expected XRef stream, got Type %v", stream.Dictionary[Name("Type")])
		}

		return r.parseXRefStream(&stream)
	}

	return fmt.Errorf("expected 'xref' or object number, got %v", tok)
}

func (r *Reader) parseXRefTable(t *Tokenizer) error {
	// Read subsections
	for {
		// Expect start_obj_num count
		tok, err := t.NextToken()
		if err != nil {
			return fmt.Errorf("failed to read xref subsection start: %w", err)
		}

		// If we hit "trailer", we are done with xref
		if tok.Value == "trailer" {
			break
		}

		startObj, err := toInt(tok)
		if err != nil {
			return fmt.Errorf("expected object number, got %v: %w", tok, err)
		}

		tok, err = t.NextToken()
		if err != nil {
			return fmt.Errorf("failed to read xref subsection count: %w", err)
		}
		count, err := toInt(tok)
		if err != nil {
			return fmt.Errorf("expected count, got %v: %w", tok, err)
		}

		for i := 0; i < count; i++ {
			tokOffset, err := t.NextToken()
			if err != nil {
				return fmt.Errorf("failed to read xref entry offset: %w", err)
			}
			offset, err := toInt64(tokOffset)
			if err != nil {
				return fmt.Errorf("invalid xref entry offset: %w", err)
			}

			// Gen number
			tokGen, err := t.NextToken()
			if err != nil {
				return fmt.Errorf("failed to read xref entry gen: %w", err)
			}
			gen, _ := toInt(tokGen)

			tokFlag, err := t.NextToken()
			if err != nil {
				return fmt.Errorf("failed to read xref entry flag: %w", err)
			}

			if tokFlag.Value == "n" {
				r.XrefTable[startObj+i] = XrefEntry{
					Type:   1,
					Offset: offset,
					Gen:    gen,
				}
			}
		}
	}

	// We hit "trailer". Parse dictionary.
	p := NewParser(t)
	// Expect "<<"
	tok, err := t.NextToken()
	if err != nil {
		return fmt.Errorf("failed to read trailer dictionary start: %w", err)
	}
	if tok.Type != TokenDictStart {
		return fmt.Errorf("expected trailer dictionary start '<<', got %v", tok)
	}

	t.UnreadToken(tok)
	obj, err := p.ParseObject()
	if err != nil {
		return fmt.Errorf("failed to parse trailer dictionary: %w", err)
	}

	dict, ok := obj.(Dictionary)
	if !ok {
		return fmt.Errorf("trailer is not a dictionary")
	}
	r.Trailer = dict

	return r.extractRoot()
}

func (r *Reader) parseXRefStream(stream *Stream) error {
	// 1. Get W array
	wArr, ok := stream.Dictionary[Name("W")].(Array)
	if !ok || len(wArr) != 3 {
		return fmt.Errorf("invalid W array in XRef stream")
	}
	w := make([]int, 3)
	for i := 0; i < 3; i++ {
		if val, ok := wArr[i].(Integer); ok {
			w[i] = int(val)
		} else if val, ok := wArr[i].(Real); ok {
			w[i] = int(val)
		} else {
			return fmt.Errorf("invalid W value")
		}
	}

	// 2. Get Index array (optional, default [0 Size])
	var index []int
	if idxArr, ok := stream.Dictionary[Name("Index")].(Array); ok {
		for _, val := range idxArr {
			if i, ok := val.(Integer); ok {
				index = append(index, int(i))
			}
		}
	} else {
		sizeObj, ok := stream.Dictionary[Name("Size")]
		if !ok {
			return fmt.Errorf("missing Size in XRef stream")
		}
		size := 0
		if s, ok := sizeObj.(Integer); ok {
			size = int(s)
		}
		index = []int{0, size}
	}

	// 3. Decode data
	data := stream.Data
	if filter, ok := stream.Dictionary[Name("Filter")].(Name); ok {
		decoded, err := DecodeStream(data, filter)
		if err != nil {
			return fmt.Errorf("failed to decode xref stream: %w", err)
		}
		data = decoded
	}

	// 4. Iterate
	entryLen := w[0] + w[1] + w[2]
	if entryLen == 0 {
		return fmt.Errorf("invalid W array (sum is 0)")
	}

	buf := bytes.NewReader(data)

	for i := 0; i < len(index); i += 2 {
		startObj := index[i]
		count := index[i+1]

		for j := 0; j < count; j++ {
			b := make([]byte, entryLen)
			_, err := io.ReadFull(buf, b)
			if err != nil {
				if err == io.EOF {
					break
				}
				return fmt.Errorf("failed to read xref entry: %w", err)
			}

			// Extract fields
			f1 := readField(b[0:w[0]])
			f2 := readField(b[w[0] : w[0]+w[1]])
			f3 := readField(b[w[0]+w[1] : w[0]+w[1]+w[2]])

			objNum := startObj + j

			switch f1 {
			case 0: // Free entry
				// Ignore
			case 1: // In-use entry (offset)
				r.XrefTable[objNum] = XrefEntry{
					Type:   1,
					Offset: f2,
					Gen:    int(f3),
				}
			case 2: // Compressed object (objNum of stream, index in stream)
				r.XrefTable[objNum] = XrefEntry{
					Type:         2,
					StreamObjNum: int(f2),
					StreamIndex:  int(f3),
				}
			}
		}
	}

	// 5. Handle Trailer keys
	r.Trailer = stream.Dictionary

	return r.extractRoot()
}

func (r *Reader) extractRoot() error {
	if rootRef, ok := r.Trailer[Name("Root")].(IndirectRef); ok {
		if _, ok := r.XrefTable[rootRef.ObjectNumber]; ok {
			rootObj, err := r.ReadObject(rootRef.ObjectNumber)
			if err != nil {
				// Don't fail if root object cannot be read, just return nil
				// But maybe we should log it?
				return nil
			}
			if rootDict, ok := rootObj.(Dictionary); ok {
				r.Root = rootDict
			}
		}
	}
	return nil
}

func readField(b []byte) int64 {
	var val int64 = 0
	for _, x := range b {
		val = (val << 8) | int64(x)
	}
	return val
}

// ReadObject reads an indirect object by number
func (r *Reader) ReadObject(objNum int) (Object, error) {
	entry, ok := r.XrefTable[objNum]
	if !ok {
		return nil, fmt.Errorf("object %d not found in xref", objNum)
	}

	if entry.Type == 2 {
		return r.readCompressedObject(entry.StreamObjNum, entry.StreamIndex)
	}

	_, err := r.f.Seek(entry.Offset, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("failed to seek to object offset %d: %w", entry.Offset, err)
	}

	t := NewTokenizer(r.f)

	// Expect: ObjNum GenNum obj
	tok, err := t.NextToken()
	if err != nil {
		return nil, fmt.Errorf("failed to read object header: %w", err)
	}
	if id, err := toInt(tok); err != nil || id != objNum {
		return nil, fmt.Errorf("expected object id %d, got %v", objNum, tok)
	}

	tok, err = t.NextToken()
	if err != nil {
		return nil, fmt.Errorf("failed to read object gen: %w", err)
	}
	// Gen num - ignore for now

	tok, err = t.NextToken()
	if err != nil {
		return nil, fmt.Errorf("failed to read 'obj' keyword: %w", err)
	}
	if tok.Value != "obj" {
		return nil, fmt.Errorf("expected 'obj', got %v", tok)
	}

	// Parse content
	p := NewParser(t)
	obj, err := p.ParseObject()
	if err != nil {
		return nil, fmt.Errorf("failed to parse object content: %w", err)
	}

	// Expect: endobj
	tok, err = t.NextToken()
	if err != nil {
		return nil, fmt.Errorf("failed to read 'endobj': %w", err)
	}
	if tok.Value != "endobj" {
		return nil, fmt.Errorf("expected 'endobj', got %v", tok)
	}

	return obj, nil
}

func (r *Reader) readCompressedObject(streamObjNum, index int) (Object, error) {
	// Read the stream object
	obj, err := r.ReadObject(streamObjNum)
	if err != nil {
		return nil, fmt.Errorf("failed to read object stream %d: %w", streamObjNum, err)
	}

	stream, ok := obj.(Stream)
	if !ok {
		return nil, fmt.Errorf("object %d is not a stream", streamObjNum)
	}

	return r.parseObjStm(stream, index)
}

func (r *Reader) parseObjStm(stream Stream, targetIndex int) (Object, error) {
	// Decode data
	data, err := decodeStream(stream)
	if err != nil {
		return nil, fmt.Errorf("failed to decode object stream: %w", err)
	}

	// Parse N and First from dictionary
	nObj, ok := stream.Dictionary[Name("N")]
	if !ok {
		return nil, fmt.Errorf("missing N in ObjStm")
	}
	n, ok := nObj.(Integer)
	if !ok {
		return nil, fmt.Errorf("N is not an integer")
	}

	firstObj, ok := stream.Dictionary[Name("First")]
	if !ok {
		return nil, fmt.Errorf("missing First in ObjStm")
	}
	first, ok := firstObj.(Integer)
	if !ok {
		return nil, fmt.Errorf("First is not an integer")
	}

	// Parse header (N pairs of integers)
	// We can use a tokenizer on the data
	t := NewTokenizer(bytes.NewReader(data))

	var offset int
	found := false

	for i := 0; i < int(n); i++ {
		// ObjNum
		_, err := t.NextToken()
		if err != nil {
			return nil, fmt.Errorf("failed to read objstm header: %w", err)
		}
		// Offset
		offTok, err := t.NextToken()
		if err != nil {
			return nil, fmt.Errorf("failed to read objstm offset: %w", err)
		}
		off, err := toInt(offTok)
		if err != nil {
			return nil, fmt.Errorf("invalid objstm offset: %w", err)
		}

		if i == targetIndex {
			offset = off
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("index %d not found in ObjStm", targetIndex)
	}

	// The object is at data[first + offset]
	start := int(first) + offset
	if start >= len(data) {
		return nil, fmt.Errorf("object offset out of bounds")
	}

	// Parse object
	// Create a new parser for the object data
	p := NewParser(NewTokenizer(bytes.NewReader(data[start:])))
	return p.ParseObject()
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
		return 0, fmt.Errorf("failed to read /Pages object: %w", err)
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
		return nil, fmt.Errorf("failed to read page tree node: %w", err)
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
