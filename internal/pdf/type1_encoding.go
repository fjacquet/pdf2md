package pdf

import (
	"bytes"
	"regexp"
	"strconv"
)

var dupPut = regexp.MustCompile(`dup\s+(\d+)\s*/([^\s/\[\]{}()<>]+)\s+put`)

// type1BuiltinEncoding reads the built-in encoding of a Type 1 font program:
// the "dup <code> /<name> put" entries in its cleartext part, before eexec.
// It returns nil when the program uses StandardEncoding or has no entries.
func type1BuiltinEncoding(fontFile []byte) map[int]string {
	if i := bytes.Index(fontFile, []byte("eexec")); i >= 0 {
		fontFile = fontFile[:i]
	}
	var enc map[int]string
	for _, m := range dupPut.FindAllSubmatch(fontFile, -1) {
		code, err := strconv.Atoi(string(m[1]))
		if err != nil || code > 255 {
			continue
		}
		if enc == nil {
			enc = map[int]string{}
		}
		enc[code] = string(m[2])
	}
	return enc
}
