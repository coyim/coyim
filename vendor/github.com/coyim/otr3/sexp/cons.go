package sexp

import "bufio"

// Cons is an S-Expression cons cell
type Cons struct {
	first  Value
	second Value
}

// First returns the first value of the cons
func (l Cons) First() Value {
	return l.first
}

// Second returns the second value of the cons
func (l Cons) Second() Value {
	return l.second
}

// Value returns the first value inside this Cons
func (l Cons) Value() interface{} {
	return l.first
}

// String returns the Cons formatted for printing in an S-Expression
func (l Cons) String() string {
	return "(" + l.First().String() + " . " + l.Second().String() + ")"
}

// List creates a chain of Cons cells ended with a nil
func List(values ...Value) Value {
	var result Value = snil
	l := len(values)

	for i := l - 1; i >= 0; i-- {
		result = Cons{values[i], result}
	}

	return result
}

// ReadListStart will expect the start character for an S-Expression list and return false if not encountered
func ReadListStart(r *bufio.Reader) bool {
	return expect(r, '(')
}

// ReadListEnd will expect the end character for an S-Expression list and return false if not encountered
func ReadListEnd(r *bufio.Reader) bool {
	return expect(r, ')')
}

// ReadList will read a list and return it
func ReadList(r *bufio.Reader) Value {
	ReadWhitespace(r)
	if !ReadListStart(r) {
		return nil
	}
	result := ReadListItem(r)
	if !ReadListEnd(r) {
		return nil
	}
	return result
}

// ReadListItem recursively read a list item and the next potential item
func ReadListItem(r *bufio.Reader) Value {
	ReadWhitespace(r)
	val, end := ReadValue(r)
	if end {
		return Snil{}
	}
	return Cons{val, ReadListItem(r)}
}
