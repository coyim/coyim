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
		{"user\u200B", "user "},
		{"user\u2002", "user "},
		{"\xC2\xAA", "a"},
		{"x\xC2\xADy", "xy"},
		{"\xE2\x85\xA3", "IV"},
		{"\xE2\x85\xA8", "IX"},
		{"\u034F\u1806\u180Bb\u180C\u180Dy\u200Ct\u200D\u2060\uFE00e\uFE01\uFE02\uFE03\uFE04\uFE05\uFE06\uFE07\uFE08\uFE09\uFE0A\uFE0B\uFE0C\uFE0D\uFE0E\uFE0F\uFEFF", "byte"},
		//They should error because they have forbidden chars
		//{"\x07", ""},      //should error
		//{"\xD8\xA71", ""}, //shold error
	}

	for _, test := range testCases {
		t := transform.NewReader(strings.NewReader(test.raw), Stringprep)
		r := bufio.NewReader(t)
		normalized, _, err := r.ReadLine()

		c.Check(err, IsNil)
		c.Check(string(normalized), DeepEquals, test.normalized)
	}
}

func identity(r rune) rune {
	return r
}

func (s *SASLSuite) Test_replaceTransformed_Transform_doesntDealWithTooShortDestination(c *C) {
	_, _, err := replaceTransformer(identity).Transform([]byte{}, []byte("user"), true)
	c.Assert(err, Equals, transform.ErrShortDst)
}

func (s *SASLSuite) Test_replaceTransformed_Transform_doesntDealWithIncompleteRune(c *C) {
	var dst [4]byte
	_, _, err := replaceTransformer(identity).Transform(dst[:], []byte("\xE2\x85"), false)
	c.Assert(err, Equals, transform.ErrShortSrc)
}
