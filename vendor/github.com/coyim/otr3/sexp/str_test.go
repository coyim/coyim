package sexp

import (
	"bufio"
	"bytes"
	"testing"
)

func Test_Sstring_First_generatesAPanic(t *testing.T) {
	defer checkForPanic(t, "not valid to call First on an SString")
	Sstring("ABCD").First()
}

func Test_Sstring_Second_generatesAPanic(t *testing.T) {
	defer checkForPanic(t, "not valid to call Second on an SString")
	Sstring("ABCD").Second()
}

func Test_Sstring_Value_returnsTheStringRepresentation(t *testing.T) {
	res := Sstring("ABB").Value()
	assertDeepEquals(t, res, "ABB")
}

func Test_Sstring_String_returnsTheStringRepresentation(t *testing.T) {
	res := Sstring("ABB").String()
	assertDeepEquals(t, res, "\"ABB\"")
}

func Test_ReadString_returnsNilIfAskedToReadNonString(t *testing.T) {
	res := ReadString(bufio.NewReader(bytes.NewReader([]byte("a"))))
	assertEquals(t, res, nil)
}

func Test_ReadString_returnsNilIfAskedToReadNonFinishedString(t *testing.T) {
	res := ReadString(bufio.NewReader(bytes.NewReader([]byte("\"a"))))
	assertEquals(t, res, nil)
}
