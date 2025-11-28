package pdf

import (
	"fmt"
	"strings"
)

// Object represents any PDF object
type Object interface {
	String() string
}

// Null represents a null object
type Null struct{}

func (n Null) String() string { return "null" }

// Boolean represents a boolean object
type Boolean bool

func (b Boolean) String() string {
	if b {
		return "true"
	}
	return "false"
}

// Integer represents an integer object
type Integer int64

func (i Integer) String() string { return fmt.Sprintf("%d", i) }

// Real represents a real number object
type Real float64

func (r Real) String() string { return fmt.Sprintf("%f", r) }

// Name represents a name object (e.g. /Type)
type Name string

func (n Name) String() string { return "/" + string(n) }

// StringLiteral represents a string object (...)
type StringLiteral string

func (s StringLiteral) String() string { return "(" + string(s) + ")" }

// HexString represents a hex string object <...>
type HexString string

func (s HexString) String() string { return "<" + string(s) + ">" }

// Array represents an array object [...]
type Array []Object

func (a Array) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, obj := range a {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(obj.String())
	}
	sb.WriteString("]")
	return sb.String()
}

// Dictionary represents a dictionary object << ... >>
type Dictionary map[Name]Object

func (d Dictionary) String() string {
	var sb strings.Builder
	sb.WriteString("<<")
	for k, v := range d {
		sb.WriteString(" ")
		sb.WriteString(k.String())
		sb.WriteString(" ")
		sb.WriteString(v.String())
	}
	sb.WriteString(" >>")
	return sb.String()
}

// Stream represents a stream object
type Stream struct {
	Dictionary Dictionary
	Data       []byte // Raw data (possibly compressed)
}

func (s Stream) String() string {
	return fmt.Sprintf("%s\nstream\n...%d bytes...\nendstream", s.Dictionary.String(), len(s.Data))
}

// IndirectRef represents an indirect reference (ObjNum GenNum R)
type IndirectRef struct {
	ObjectNumber     int
	GenerationNumber int
}

func (r IndirectRef) String() string {
	return fmt.Sprintf("%d %d R", r.ObjectNumber, r.GenerationNumber)
}

// IndirectObject represents a definition (ObjNum GenNum obj ... endobj)
type IndirectObject struct {
	ObjectNumber     int
	GenerationNumber int
	Object           Object
}
