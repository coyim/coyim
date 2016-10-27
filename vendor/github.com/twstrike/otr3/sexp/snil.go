package sexp

// Snil is a representation of nil in an S-Expression.
type Snil struct{}

var snil = Snil{}

// First will return the nil value.
func (l Snil) First() Value {
	return l
}

// Second will return the nil value.
func (l Snil) Second() Value {
	return l
}

// Value will return the Golang nil.
func (l Snil) Value() interface{} {
	return nil
}

// String will return an empty list as a representation of nil.
func (l Snil) String() string {
	return "()"
}
