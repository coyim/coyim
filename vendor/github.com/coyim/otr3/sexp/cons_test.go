package sexp

import (
	"bufio"
	"bytes"
	"testing"
)

func Test_Cons_First_ReturnsTheFirstElement(t *testing.T) {
	c := Cons{Sstring("First"), Sstring("Second")}
	assertEquals(t, c.First(), Sstring("First"))
}

func Test_Cons_Second_ReturnsTheSecondElement(t *testing.T) {
	c := Cons{Sstring("First"), Sstring("Second")}
	assertEquals(t, c.Second(), Sstring("Second"))
}

func Test_Cons_Value_ReturnsTheFirstElement(t *testing.T) {
	c := Cons{Sstring("First"), Sstring("Second")}
	assertEquals(t, c.Value(), Sstring("First"))
}

func Test_Cons_String_ReturnsTheConsFormatted(t *testing.T) {
	c := Cons{Sstring("First"), Sstring("Second")}
	assertEquals(t, c.String(), "(\"First\" . \"Second\")")
}

func Test_ReadList_returnsNilIfAskedToReadNonList(t *testing.T) {
	res := ReadList(bufio.NewReader(bytes.NewReader([]byte("a"))))
	assertEquals(t, res, nil)
}

func Test_ReadList_returnsNilIfAskedToReadListThatDoesntEndWell(t *testing.T) {
	res := ReadList(bufio.NewReader(bytes.NewReader([]byte("("))))
	assertEquals(t, res, nil)
}
