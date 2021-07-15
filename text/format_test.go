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
			&compositeFragment{
				[]string{
					"admin",
					"}",
					"foo",
				},
			},
		},
	})

	res, ok = ParseWithFormat("hello $role{$}foo}")
	c.Assert(ok, Equals, true)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello "},
		&fragmentWithFormat{
			"role",
			&compositeFragment{
				[]string{
					"}",
					"foo",
				},
			},
		},
	})

	res, ok = ParseWithFormat("hello $role{stf$}}")
	c.Assert(ok, Equals, true)
	c.Assert(res, DeepEquals, FormattedText{
		&textFragment{"hello "},
		&fragmentWithFormat{
			"role",
			&compositeFragment{
				[]string{
					"stf",
					"}",
				},
			},
		},
	})
}

// Incorrect escape inside of brackets
// Failure of parsing the format
//  - no format name / invalid format name
//  - no starting brace
//  - no ending brace
