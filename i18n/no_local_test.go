package i18n

import (
	. "gopkg.in/check.v1"
)

func (s *I18NSuite) Test_NoLocal_Local_returnsTheString(c *C) {
	c.Assert(NoLocal.Local("foo"), Equals, "foo")
	c.Assert(NoLocal.Local("bar"), Equals, "bar")
}
