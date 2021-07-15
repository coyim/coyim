package text

import . "gopkg.in/check.v1"

type FormatSuite struct{}

var _ = Suite(&FormatSuite{})

func (s *FormatSuite) Test_ParseWithFormat_parsesASimpleTextWithoutFormatting(c *C) {
	res, ok := ParseWithFormat("hello world")
	c.Assert(ok, Equals, true)
	c.Assert(res, DeepEquals, FormattedText{&textFragment{"hello world"}})
}

func (s *FormatSuite) Test_ParseWithFormat_parsesATextWithOneFormat(c *C) {
	res, ok := ParseWithFormat("hello world, $nick{Luke}")
	c.Assert(ok, Equals, true)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello world, "},
		&fragmentWithFormat{
			"nick",
			&textFragment{"Luke"},
		},
	})
}

func (s *FormatSuite) Test_ParseWithFormat_parsesATextWithTwoFormats(c *C) {
	res, ok := ParseWithFormat("hello and welcome, $nick{Luke} - it's time to start - you are $role{an administrator}")
	c.Assert(ok, Equals, true)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello and welcome, "},
		&fragmentWithFormat{
			"nick",
			&textFragment{"Luke"},
		},
		&textFragment{" - it's time to start - you are "},
		&fragmentWithFormat{
			"role",
			&textFragment{"an administrator"},
		},
	})
}

func (s *FormatSuite) Test_ParseWithFormat_parsesATextWithDollarEscape(c *C) {
	res, ok := ParseWithFormat("$$")
	c.Assert(ok, Equals, true)
	c.Assert(res, DeepEquals, FormattedText{&textFragment{"$"}})

	res, ok = ParseWithFormat("hello $$")
	c.Assert(ok, Equals, true)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello "},
		&textFragment{"$"},
	})

	res, ok = ParseWithFormat("hello $$ somewhere")
	c.Assert(ok, Equals, true)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello "},
		&textFragment{"$"},
		&textFragment{" somewhere"},
	})

	res, ok = ParseWithFormat("hello $$followed{by} somewhere")
	c.Assert(ok, Equals, true)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello "},
		&textFragment{"$"},
		&textFragment{"followed{by} somewhere"},
	})
}

func (s *FormatSuite) Test_ParseWithFormat_parsesATextWithFailedEscape(c *C) {
	res, ok := ParseWithFormat("$")
	c.Assert(ok, Equals, false)
	c.Assert(res, DeepEquals, FormattedText{&textFragment{"$"}})

	res, ok = ParseWithFormat("hmm $")
	c.Assert(ok, Equals, false)
	c.Assert(res, DeepEquals, FormattedText{&textFragment{"hmm $"}})
}

func (s *FormatSuite) Test_ParseWithFormat_parsesATextWithEscapeOfEndingBracket(c *C) {
	res, ok := ParseWithFormat("hello $role{admin$}foo}")
	c.Assert(ok, Equals, true)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello "},
		&fragmentWithFormat{
			"role",
			&textFragment{"admin}foo"},
		},
	})

	res, ok = ParseWithFormat("hello $role{$}foo}")
	c.Assert(ok, Equals, true)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello "},
		&fragmentWithFormat{
			"role",
			&textFragment{"}foo"},
		},
	})

	res, ok = ParseWithFormat("hello $role{stf$}}")
	c.Assert(ok, Equals, true)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello "},
		&fragmentWithFormat{
			"role",
			&textFragment{"stf}"},
		},
	})
}

func (s *FormatSuite) Test_ParseWithFormat_parsesATextWithEscapeOfEscapeInsideFormat(c *C) {
	res, ok := ParseWithFormat("hello $role{admin$$}foo")
	c.Assert(ok, Equals, true)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello "},
		&fragmentWithFormat{
			"role",
			&textFragment{"admin$"},
		},
		&textFragment{"foo"},
	})
}

func (s *FormatSuite) Test_ParseWithFormat_failsOnMissingFormattingName(c *C) {
	res, ok := ParseWithFormat("hello $*** what's up?")
	c.Assert(ok, Equals, false)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello $*** what's up?"},
	})

	res, ok = ParseWithFormat("hello ${hmm} what's up?")
	c.Assert(ok, Equals, false)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello ${hmm} what's up?"},
	})
}

func (s *FormatSuite) Test_ParseWithFormat_failsOnMissingFormattingText(c *C) {
	res, ok := ParseWithFormat("hello $role what's up?")
	c.Assert(ok, Equals, false)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello $role what's up?"},
	})
}

func (s *FormatSuite) Test_Join_simpleText(c *C) {
	res, _ := ParseWithFormat("hello world")
	txt, st, l, f := res.Join()

	c.Assert(txt, Equals, "hello world")
	c.Assert(st, HasLen, 0)
	c.Assert(l, HasLen, 0)
	c.Assert(f, HasLen, 0)
}

func (s *FormatSuite) Test_Join_moreThanOneTextFragment(c *C) {
	res, _ := ParseWithFormat("hello world")
	res2 := append(res, res...)
	txt, st, l, f := res2.Join()

	c.Assert(txt, Equals, "hello worldhello world")
	c.Assert(st, HasLen, 0)
	c.Assert(l, HasLen, 0)
	c.Assert(f, HasLen, 0)
}

func (s *FormatSuite) Test_Join_withASimpleFormat(c *C) {
	res, _ := ParseWithFormat("hello world, $nick{Luke} - what's up?")

	txt, st, l, f := res.Join()

	c.Assert(txt, Equals, "hello world, Luke - what's up?")

	c.Assert(st, HasLen, 1)
	c.Assert(st[0], Equals, 13)

	c.Assert(l, HasLen, 1)
	c.Assert(l[0], Equals, 4)

	c.Assert(f, HasLen, 1)
	c.Assert(f[0], Equals, "nick")
}

func (s *FormatSuite) Test_Join_withASimpleFormatAtStart(c *C) {
	res, _ := ParseWithFormat("$nick{Luke} - what's up?")

	txt, st, l, f := res.Join()

	c.Assert(txt, Equals, "Luke - what's up?")

	c.Assert(st, HasLen, 1)
	c.Assert(st[0], Equals, 0)

	c.Assert(l, HasLen, 1)
	c.Assert(l[0], Equals, 4)

	c.Assert(f, HasLen, 1)
	c.Assert(f[0], Equals, "nick")
}

func (s *FormatSuite) Test_Join_withASimpleFormatAtEnd(c *C) {
	res, _ := ParseWithFormat("hello world, $nick{Luke}")

	txt, st, l, f := res.Join()

	c.Assert(txt, Equals, "hello world, Luke")

	c.Assert(st, HasLen, 1)
	c.Assert(st[0], Equals, 13)

	c.Assert(l, HasLen, 1)
	c.Assert(l[0], Equals, 4)

	c.Assert(f, HasLen, 1)
	c.Assert(f[0], Equals, "nick")
}

func (s *FormatSuite) Test_Join_withMoreThanOneFormatAndEscapes(c *C) {
	res, _ := ParseWithFormat("hello and welcome $$42$$, $role{foo{$}bar$$} - you are $nick{someone}")

	txt, st, l, f := res.Join()

	c.Assert(txt, Equals, "hello and welcome $42$, foo{}bar$ - you are someone")

	c.Assert(st, HasLen, 2)
	c.Assert(st[0], Equals, 24)
	c.Assert(st[1], Equals, 44)

	c.Assert(l, HasLen, 2)
	c.Assert(l[0], Equals, 9)
	c.Assert(l[1], Equals, 7)

	c.Assert(f, HasLen, 2)
	c.Assert(f[0], Equals, "role")
	c.Assert(f[1], Equals, "nick")
}

// Incorrect escape inside of brackets
// Failure of parsing the format
//  - no ending brace
