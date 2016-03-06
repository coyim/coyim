package sexp

import "bufio"

// Symbol represents an S-Expression symbol.
type Symbol string

// First will fail if called on a symbol
func (s Symbol) First() Value {
	panic("not valid to call First on a Symbol")
}

// Second will fail if called on a symbol
func (s Symbol) Second() Value {
	panic("not valid to call Second on a Symbol")
}

// String returns the symbol as a string
func (s Symbol) String() string {
	return string(s)
}

// Value returns the symbol as a string
func (s Symbol) Value() interface{} {
	return string(s)
}

// ReadSymbol will read a symbol from the reader
func ReadSymbol(r *bufio.Reader) Value {
	ReadWhitespace(r)
	result := ReadDataUntil(r, isNotSymbolCharacter)
	return Symbol(result)
}
