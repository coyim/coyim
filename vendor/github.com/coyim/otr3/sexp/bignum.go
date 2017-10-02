package sexp

import (
	"bufio"
	"fmt"
	"math/big"
)

// BigNum is an S-Expression big number
type BigNum struct {
	val *big.Int
}

// NewBigNum creates a new BigNum from the hex formatted string given
func NewBigNum(s string) BigNum {
	res, _ := new(big.Int).SetString(s, 16)
	return BigNum{res}
}

// First will cause an error when called on a BigNum
func (s BigNum) First() Value {
	panic("not valid to call First on a BigNum")
}

// Second will cause an error when called on a BigNum
func (s BigNum) Second() Value {
	panic("not valid to call Second on a BigNum")
}

// String returns the BigNum formatted for printing in an S-Expression
func (s BigNum) String() string {
	return "#" + fmt.Sprintf("%X", s.val) + "#"
}

// Value returns the *big.Int value inside of this BigNum
func (s BigNum) Value() interface{} {
	return s.val
}

// ReadBigNumStart will expect the start character for an S-Expression bignum and return false if not encountered
func ReadBigNumStart(r *bufio.Reader) bool {
	return expect(r, '#')
}

// ReadBigNumEnd will expect the end character for an S-Expression bignum and return false if not encountered
func ReadBigNumEnd(r *bufio.Reader) bool {
	return expect(r, '#')
}

// ReadBigNum will read a bignum from the given reader and return it
func ReadBigNum(r *bufio.Reader) Value {
	ReadWhitespace(r)
	if !ReadBigNumStart(r) {
		return nil
	}
	result := ReadDataUntil(r, untilFixed('#'))
	if !ReadBigNumEnd(r) {
		return nil
	}
	return NewBigNum(string(result))
}
