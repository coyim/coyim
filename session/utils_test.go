package session

import (
	. "gopkg.in/check.v1"
)

type UtilsSuite struct{}

var _ = Suite(&UtilsSuite{})

func (s *UtilsSuite) Test_either_works(c *C) {
	c.Assert(either("", ""), Equals, "")
	c.Assert(either("foo", ""), Equals, "foo")
	c.Assert(either("", "bar"), Equals, "bar")
	c.Assert(either("foo", "bar"), Equals, "foo")
}
