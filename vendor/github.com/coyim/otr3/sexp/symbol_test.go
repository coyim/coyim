package sexp

import "testing"

func Test_Symbol_First_generatesAPanic(t *testing.T) {
	defer checkForPanic(t, "not valid to call First on a Symbol")
	Symbol("ABCD").First()
}

func Test_Symbol_Second_generatesAPanic(t *testing.T) {
	defer checkForPanic(t, "not valid to call Second on a Symbol")
	Symbol("ABCD").Second()
}

func Test_Symbol_Value_returnsTheStringRepresentation(t *testing.T) {
	res := Symbol("ABB").Value()
	assertDeepEquals(t, res, "ABB")
}

func Test_Symbol_String_returnsTheStringRepresentation(t *testing.T) {
	res := Symbol("ABB").String()
	assertDeepEquals(t, res, "ABB")
}
