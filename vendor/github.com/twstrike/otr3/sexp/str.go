package sexp

import "bufio"

// Sstring represents an S-Expression symbol.
type Sstring string

// First will fail if called on an Sstring
func (s Sstring) First() Value {
	panic("not valid to call First on an SString")
}

// Second will fail if called on an Sstring
func (s Sstring) Second() Value {
	panic("not valid to call Second on an SString")
}

// String returns the string quoted as a string in an S-Expression
func (s Sstring) String() string {
	return "\"" + string(s) + "\""
}

// Value returns the string as a string
func (s Sstring) Value() interface{} {
	return string(s)
}

// ReadStringStart will read the string start character and return false if it is not encountered
func ReadStringStart(r *bufio.Reader) bool {
	return expect(r, '"')
}

// ReadStringEnd will read the string end character and return false if it is not encountered
func ReadStringEnd(r *bufio.Reader) bool {
	return expect(r, '"')
}

// ReadString will read a string from the reader
func ReadString(r *bufio.Reader) Value {
	ReadWhitespace(r)
	if !ReadStringStart(r) {
		return nil
	}
	result := ReadDataUntil(r, untilFixed('"'))
	if !ReadStringEnd(r) {
		return nil
	}
	return Sstring(result)
}
