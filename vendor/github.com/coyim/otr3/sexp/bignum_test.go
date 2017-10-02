package sexp

import (
	"bufio"
	"bytes"
	"math/big"
	"testing"
)

func Test_BigNum_First_generatesAPanic(t *testing.T) {
	defer checkForPanic(t, "not valid to call First on a BigNum")
	NewBigNum("ABCD").First()
}

func Test_BigNum_Second_generatesAPanic(t *testing.T) {
	defer checkForPanic(t, "not valid to call Second on a BigNum")
	NewBigNum("ABCD").Second()
}

func Test_BigNum_String_generatesAValidRepresentation(t *testing.T) {
	res := NewBigNum("ABCD").String()
	assertEquals(t, res, "#ABCD#")
}

func Test_BigNum_Value_returnsTheBigNumInside(t *testing.T) {
	res := NewBigNum("ABCD").Value().(*big.Int)
	assertDeepEquals(t, res, new(big.Int).SetBytes([]byte{0xAB, 0xCD}))
}

func Test_ReadBigNum_returnsNilIfAskedToReadNonBigNum(t *testing.T) {
	res := ReadBigNum(bufio.NewReader(bytes.NewReader([]byte(""))))
	assertEquals(t, res, nil)
}

func Test_ReadBigNum_returnsNilIfAskedToReadBigNumThatDoesntEndCorrectly(t *testing.T) {
	res := ReadBigNum(bufio.NewReader(bytes.NewReader([]byte("#ABCD"))))
	assertEquals(t, res, nil)
}
