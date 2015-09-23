package ui

import (
	"bytes"
	"testing"
)

var escapingTests = []string{
	"",
	"foo",
	"foo\\",
	"foo\\x01",
	"العربية",
}

func TestEscaping(t *testing.T) {
	for _, test := range escapingTests {
		escaped := EscapeNonASCII(test)
		unescaped, err := UnescapeNonASCII(escaped)
		if err != nil {
			t.Errorf("Error unescaping '%s' (from '%s')", escaped, test)
			continue
		}
		if unescaped != test {
			t.Errorf("Unescaping didn't return the original value: '%s' -> '%s' -> '%s'", test, escaped, unescaped)
		}
	}
}

func TestHTMLStripping(t *testing.T) {
	raw := []byte("<hr>This is some <font color='green'>html</font><br />.")
	exp := []byte("This is some html.")
	res := StripHTML(raw)
	if !bytes.Equal(res, exp) {
		t.Errorf("HTML wasn't properly stripped: '%s' -> '%s' but expected '%s'", raw, res, exp)
	}

}
