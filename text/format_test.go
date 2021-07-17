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

func (s *FormatSuite) Test_ParseWithFormat_failsOnFormatNameAtEnd(c *C) {
	res, ok := ParseWithFormat("hello $rol")
	c.Assert(ok, Equals, false)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello $rol"},
	})
}

func (s *FormatSuite) Test_ParseWithFormat_failsOnNoEndingBrace(c *C) {
	res, ok := ParseWithFormat("hello $role{bla")
	c.Assert(ok, Equals, false)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello $role{bla"},
	})
}

func (s *FormatSuite) Test_ParseWithFormat_incorrectEscapeInsideOfBrackets(c *C) {
	res, ok := ParseWithFormat("hello $role{bl$a} foo")
	c.Assert(ok, Equals, false)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello $role{bl$a} foo"},
	})
}

func (s *FormatSuite) Test_Join_simpleText(c *C) {
	res, _ := ParseWithFormat("hello world")
	txt, formats := res.Join()

	c.Assert(txt, Equals, "hello world")
	c.Assert(formats, HasLen, 0)
}

func (s *FormatSuite) Test_Join_moreThanOneTextFragment(c *C) {
	res, _ := ParseWithFormat("hello world")
	res2 := append(res, res...)
	txt, formats := res2.Join()

	c.Assert(txt, Equals, "hello worldhello world")
	c.Assert(formats, HasLen, 0)
}

func (s *FormatSuite) Test_Join_withASimpleFormat(c *C) {
	res, _ := ParseWithFormat("hello world, $nick{Luke} - what's up?")

	txt, formats := res.Join()

	c.Assert(txt, Equals, "hello world, Luke - what's up?")

	c.Assert(formats, HasLen, 1)
	c.Assert(formats[0], Equals, Formatting{
		13, 4, "nick",
	})
}

func (s *FormatSuite) Test_Join_withASimpleFormatAtStart(c *C) {
	res, _ := ParseWithFormat("$nick{Luke} - what's up?")

	txt, formats := res.Join()

	c.Assert(txt, Equals, "Luke - what's up?")

	c.Assert(formats, HasLen, 1)
	c.Assert(formats[0], Equals, Formatting{
		0, 4, "nick",
	})
}

func (s *FormatSuite) Test_Join_withASimpleFormatAtEnd(c *C) {
	res, _ := ParseWithFormat("hello world, $nick{Luke}")

	txt, formats := res.Join()

	c.Assert(txt, Equals, "hello world, Luke")

	c.Assert(formats, HasLen, 1)
	c.Assert(formats[0], Equals, Formatting{
		13, 4, "nick",
	})
}

func (s *FormatSuite) Test_Join_withMoreThanOneFormatAndEscapes(c *C) {
	res, _ := ParseWithFormat("hello and welcome $$42$$, $role{foo{$}bar$$} - you are $nick{someone}")

	txt, formats := res.Join()

	c.Assert(txt, Equals, "hello and welcome $42$, foo{}bar$ - you are someone")

	c.Assert(formats, HasLen, 2)
	c.Assert(formats[0], Equals, Formatting{
		24, 9, "role",
	})
	c.Assert(formats[1], Equals, Formatting{
		44, 7, "nick",
	})
}
