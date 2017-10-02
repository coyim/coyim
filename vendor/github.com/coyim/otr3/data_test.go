package otr3

import (
	"math/big"
	"testing"
)

func Test_appendWordWillAppendTheWord(t *testing.T) {
	before := []byte{0x12, 0x14}
	result := appendWord(before, 0x334D1215)
	assertDeepEquals(t, result, []byte{0x12, 0x14, 0x33, 0x4D, 0x12, 0x15})
}

func Test_appendDataWillAppendBytes(t *testing.T) {
	before := []byte{0x13, 0x54}
	result := appendData(before, []byte{0x55, 0x12, 0x04, 0x8A, 0x00})
	assertDeepEquals(t, result, []byte{0x13, 0x54, 0x00, 0x00, 0x00, 0x05, 0x55, 0x12, 0x04, 0x8A, 0x00})
}

func Test_appendMPIWillAppendTheMPI(t *testing.T) {
	before := []byte{0x13, 0x54}
	result := appendMPI(before, new(big.Int).SetBytes([]byte{0x55, 0x12, 0x04, 0x8A, 0x00}))
	assertDeepEquals(t, result, []byte{0x13, 0x54, 0x00, 0x00, 0x00, 0x05, 0x55, 0x12, 0x04, 0x8A, 0x00})
}

func Test_appendMPIsWillAppendTheMPIs(t *testing.T) {
	before := []byte{0x13, 0x54}
	one := new(big.Int).SetBytes([]byte{0x55, 0x12, 0x04, 0x8A, 0x00})
	two := new(big.Int).SetBytes([]byte{0x01, 0x53, 0xCC})
	three := new(big.Int).SetBytes([]byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01})
	result := appendMPIs(before, one, two, three)
	assertDeepEquals(t, result, []byte{0x13, 0x54, 0x00, 0x00, 0x00, 0x05, 0x55, 0x12, 0x04, 0x8A, 0x00, 0x00, 0x00, 0x00, 0x03, 0x01, 0x53, 0xCC, 0x00, 0x00, 0x00, 0x0A, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01})
}

func Test_extractWord_extractsAllTheBytes(t *testing.T) {
	d := []byte{0x12, 0x14, 0x15, 0x17}
	_, result, _ := extractWord(d)
	assertDeepEquals(t, result, uint32(0x12141517))
}

func Test_extractWord_extractsWithError(t *testing.T) {
	d := []byte{0x12, 0x14, 0x15}
	_, result, ok := extractWord(d)
	assertDeepEquals(t, result, uint32(0))
	assertDeepEquals(t, ok, false)
}

func Test_extractShort_extractsAllTheBytes(t *testing.T) {
	d := []byte{0x12, 0x14}
	_, result, ok := extractShort(d)
	assertDeepEquals(t, result, uint16(0x1214))
	assertDeepEquals(t, ok, true)
}

func Test_extractShort_isNotOKIfThereIsNotEnoughData(t *testing.T) {
	d := []byte{0x12}
	_, result, ok := extractShort(d)
	assertDeepEquals(t, result, uint16(0))
	assertDeepEquals(t, ok, false)
}

func Test_extractData_extractsFromStartIndex(t *testing.T) {
	d := []byte{0x13, 0x54, 0x00, 0x00, 0x00, 0x05, 0x55, 0x12, 0x04, 0x8A, 0x00}
	index, result, ok := extractData(d[2:])
	assertDeepEquals(t, result, []byte{0x55, 0x12, 0x04, 0x8A, 0x00})
	assertDeepEquals(t, index, []byte{})
	assertDeepEquals(t, ok, true)
}

func Test_extractData_returnsNotOKIfThereIsntEnoughBytesForTheLength(t *testing.T) {
	d := []byte{0x13, 0x54, 0x00}
	_, _, ok := extractData(d)
	assertDeepEquals(t, ok, false)
}

func Test_extractData_returnsNotOKIfThereArentEnoughBytes(t *testing.T) {
	d := []byte{0x00, 0x00, 0x00, 0x02, 0x01}
	_, _, ok := extractData(d)
	assertDeepEquals(t, ok, false)
}

func Test_extractMPI_returnsNotOKIfThereIsNotEnoughBytesForLength(t *testing.T) {
	d := []byte{0x00, 0x00, 0x01}
	_, _, ok := extractMPI(d)
	assertDeepEquals(t, ok, false)
}

func Test_extractMPI_returnsNotOKIfThereIsNotEnoughBytesForTheMPI(t *testing.T) {
	d := []byte{0x00, 0x00, 0x00, 0x02, 0x01}
	_, _, ok := extractMPI(d)
	assertDeepEquals(t, ok, false)
}

func Test_extractMPIs_returnsNotOKIfThereIsNotEnoughBytesForLength(t *testing.T) {
	d := []byte{0x00, 0x00, 0x01}
	_, _, ok := extractMPIs(d)
	assertDeepEquals(t, ok, false)
}

func Test_extractMPIs_returnsNotOKIfOneOfTheMPIsInsideIsNotValid(t *testing.T) {
	d := []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x01}
	_, _, ok := extractMPIs(d)
	assertDeepEquals(t, ok, false)
}

func Test_extractMPIs_returnsNotOKIfThereAreNotEnoughMPIs(t *testing.T) {
	d := []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x01}
	_, _, ok := extractMPIs(d)
	assertDeepEquals(t, ok, false)
}

func Test_extractMPIs_returnsOKIfAnMPIIsReadCorrectly(t *testing.T) {
	d := []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x01}
	_, _, ok := extractMPIs(d)
	assertDeepEquals(t, ok, true)
}
