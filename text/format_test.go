package text

import . "gopkg.in/check.v1"

type FormatSuite struct{}

var _ = Suite(&FormatSuite{})

func (s *FormatSuite) Test_ParseWithFormat_parsesASimpleTextCorrectly(c *C) {
	res, ok := ParseWithFormat("hello world")
	c.Assert(ok, Equals, true)
	c.Assert(res, DeepEquals, FormattedText{&textFragment{"hello world"}})
}
