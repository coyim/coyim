package sexp

import (
	"bufio"
	"io"
)

// Value is an S-Expression value
type Value interface {
	First() Value
	Second() Value
	Value() interface{}
	String() string
}

func peek(r *bufio.Reader) (c byte, e error) {
	c, e = r.ReadByte()
	if e != io.EOF {
		r.UnreadByte()
	}
	return
}

// Read will read an S-Expression from the given reader
func Read(r *bufio.Reader) Value {
	res, _ := ReadValue(r)
	return res
}

// ReadWhitespace will read from the reader until no whitespace is encountered
func ReadWhitespace(r *bufio.Reader) {
	c, e := peek(r)
	for e != io.EOF && isWhitespace(c) {
		r.ReadByte()
		c, e = peek(r)
	}
}

// ReadValue will read and return an S-Expression value from the reader
func ReadValue(r *bufio.Reader) (Value, bool) {
	ReadWhitespace(r)
	c, err := peek(r)
	if err != nil {
		return nil, true
	}
	switch c {
	case '(':
		return ReadList(r), false
	case ')':
		return nil, true
	case '"':
		return ReadString(r), false
	case '#':
		return ReadBigNum(r), false
	default:
		return ReadSymbol(r), false
	}
}

func isWhitespace(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}

func isNotSymbolCharacter(c byte) bool {
	if isWhitespace(c) {
		return true
	}
	switch c {
	case '(', ')':
		return true
	default:
		return false
	}
}

func expect(r *bufio.Reader, c byte) bool {
	ReadWhitespace(r)
	res, err := r.ReadByte()
	if res != c {
		r.UnreadByte()
	}

	return res == c && err != io.EOF
}

func untilFixed(b byte) func(byte) bool {
	return func(until byte) bool {
		return until == b
	}
}

// ReadDataUntil will read and collect bytes from the reader until it encounters EOF or the given function returns true.
func ReadDataUntil(r *bufio.Reader, until func(byte) bool) []byte {
	result := make([]byte, 0, 10)
	c, err := peek(r)
	for err != io.EOF && !until(c) {
		r.ReadByte()
		result = append(result, c)
		c, err = peek(r)
	}
	return result
}
