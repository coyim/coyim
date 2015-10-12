package cli

import (
	"reflect"
	"testing"
)

const (
	opFind = iota
	opNext
)

var priorityListTests = []struct {
	op      int
	in, out string
}{
	{opFind, "a", "anchor"},
	{opFind, "a", "anchor"},
	{opNext, "", "anvil"},
	{opNext, "", "anchor"},
	{opNext, "", "anvil"},
	{opNext, "", "anchor"},
	{opFind, "a", "anchor"},
	{opNext, "", "anvil"},
	{opFind, "a", "anvil"},
	{opFind, "a", "anvil"},

	{opFind, "bo", "bob"},
	{opNext, "", "boom"},
	{opNext, "", "bop"},

	{opFind, "a", "anvil"},
	{opFind, "b", "bop"},

	{opFind, "c", "charlie"},
	{opNext, "", "charlie"},
	{opNext, "", "charlie"},
}

func TestPriorityList(t *testing.T) {
	var pl priorityList

	for _, word := range []string{"bop", "boom", "bob", "anvil", "anchor", "charlie"} {
		pl.Insert(word)
	}

	for i, step := range priorityListTests {
		var out string

		switch step.op {
		case opFind:
			out, _ = pl.Find(step.in)
		case opNext:
			out = pl.Next()
		default:
			panic("unknown op")
		}
		if string(out) != step.out {
			t.Fatalf("failed at step %d: got %s, want %s", i, string(out), step.out)
		}
	}
}

var testCommands = []uiCommand{
	{"a", aCommand{}, "Desc"},
	{"b", bCommand{}, "Desc"},
}

type aCommand struct {
}

type bCommand struct {
	A string
	B string "uid"
	C string
}

var parseForCompletionTests = []struct {
	in             string
	ok             bool
	before, prefix string
	isCommand      bool
}{
	{"", false, "", "", false},
	{"/", true, "/", "", true},
	{"/a", true, "/", "a", true},
	{"/ab", true, "/", "ab", true},
	{"/a b", false, "", "", false},
	{"/a b ", false, "", "", false},
	{"/a ba", false, "", "", false},
	{"/a ba c", false, "", "", false},
	{"/a ba c d", false, "", "", false},

	{"/b c", false, "", "", false},
	{"/b c ", false, "", "", false},
	{"/b c d", true, "/b c ", "d", false},
	{"/b c de", true, "/b c ", "de", false},
	{"/b c de ", false, "", "", false},
	{"/b c de f", false, "", "", false},
}

func TestParseCommandForCompletion(t *testing.T) {
	for i, test := range parseForCompletionTests {
		before, prefix, isCommand, ok := parseCommandForCompletion(testCommands, test.in)
		if ok != test.ok {
			t.Errorf("#%d: result mismatch (should be %t)", i, test.ok)
			continue
		}

		if !ok {
			continue
		}

		if before != test.before {
			t.Errorf("#%d: mismatch with 'before': got '%s', want '%s'", i, string(before), test.before)
		}
		if prefix != test.prefix {
			t.Errorf("#%d: mismatch with 'prefix': got '%s', want '%s'", i, string(before), test.before)
		}
		if isCommand != test.isCommand {
			t.Errorf("#%d: isCommand incorrect, wanted %t", i, test.isCommand)
		}
	}
}

var parseCommandTests = []struct {
	in  string
	ok  bool
	out interface{}
}{
	{"/", false, nil},
	{"/bob", false, nil},
	{"/a", true, aCommand{}},
	{"/a b", false, nil},
	{"/a b c", false, nil},
	{"/b a", false, nil},
	{"/b a b", false, nil},
	{"/b a b ", false, nil},
	{"/b a b c", true, bCommand{"a", "b", "c"}},
	{"/b a b\\  c", true, bCommand{"a", "b ", "c"}},
	{"/b a \"b b\" c", true, bCommand{"a", "b b", "c"}},
	{"/b a \"b \\\"b\" c", true, bCommand{"a", "b \"b", "c"}},
}

func TestParseCommand(t *testing.T) {
	for i, test := range parseCommandTests {
		v, err := parseCommand(testCommands, []byte(test.in))
		if (len(err) == 0) != test.ok {
			t.Errorf("#%d: bad parse result, expected %t", i, test.ok)
			continue
		}

		if !test.ok {
			continue
		}

		if !reflect.DeepEqual(v, test.out) {
			t.Errorf("#%d: bad result: got %v, want %v", i, v, test.out)
		}
	}
}
