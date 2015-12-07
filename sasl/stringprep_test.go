package sasl

import (
	"bufio"
	"strings"

	"golang.org/x/text/transform"

	. "gopkg.in/check.v1"
)

func (s *SASLSuite) TestScramNormalizesPassword(c *C) {
	// From: libidn-1.9/tests/tst_stringprep.c
	// See RFC 4013, section 3
	testCases := []struct {
		raw        string
		normalized string
	}{
		{"I\xC2\xADX", "IX"},
		{"user", "user"},
		{"USER", "USER"},
		{"\xC2\xAA", "a"},
		{"x\xC2\xADy", "xy"},
		{"\xE2\x85\xA3", "IV"},
		{"\xE2\x85\xA8", "IX"},
		//They should error because they have forbidden chars
		//{"\x07", ""},      //should error
		//{"\xD8\xA71", ""}, //shold error
	}

	for _, test := range testCases {
		t := transform.NewReader(strings.NewReader(test.raw), Stringprep)
		r := bufio.NewReader(t)
		normalized, _, err := r.ReadLine()

		c.Check(err, IsNil)
		c.Check(normalized, DeepEquals, []byte(test.normalized))
	}
}
