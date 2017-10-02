package sexp

import (
	"bufio"
	"bytes"
	"testing"
)

func inp(s string) *bufio.Reader {
	return bufio.NewReader(bytes.NewBuffer([]byte(s)))
}

func Test_parse_willParseASymbol(t *testing.T) {
	result := Read(inp("hello"))
	assertDeepEquals(t, result, Symbol("hello"))
}

func Test_parse_willParseAString(t *testing.T) {
	result := Read(inp("\"hello\""))
	assertDeepEquals(t, result, Sstring("hello"))
}

func Test_parse_willParseABigNum(t *testing.T) {
	result := Read(inp("#123FFCADDD#"))
	assertDeepEquals(t, result, NewBigNum("123FFCADDD"))
}

func Test_parse_willParseAnEmptyList(t *testing.T) {
	result := Read(inp("()"))
	assertDeepEquals(t, result, List())
}

func Test_parse_willParseAListWithAnAtom(t *testing.T) {
	result := Read(inp("(an-atom)"))
	assertDeepEquals(t, result, List(Symbol("an-atom")))
}

func Test_parse_willParseAListWithAString(t *testing.T) {
	result := Read(inp("(\"an-atom\")"))
	assertDeepEquals(t, result, List(Sstring("an-atom")))
}

func Test_parse_willParseAListWithTwoAtoms(t *testing.T) {
	result := Read(inp("(an-atom another-atom)"))
	assertDeepEquals(t, result, List(Symbol("an-atom"), Symbol("another-atom")))
}

func Test_parse_willParseAListWithNestedLists(t *testing.T) {
	result := Read(inp("(an-atom (another-atom))"))
	assertDeepEquals(t, result, List(Symbol("an-atom"), List(Symbol("another-atom"))))
}

func Test_parse_willParseAListWithSeveralLists(t *testing.T) {
	result := Read(inp("(an-atom (another-atom) (a-third))"))
	assertDeepEquals(t, result, List(Symbol("an-atom"), List(Symbol("another-atom")), List(Symbol("a-third"))))
}
