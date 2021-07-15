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
