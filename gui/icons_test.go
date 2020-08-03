package gui

import . "gopkg.in/check.v1"

type IconsSuite struct{}

var _ = Suite(&IconsSuite{})

func (*IconsSuite) Test_icons_getPath(c *C) {
	c.Assert(coyimIcon.getPath(), Matches, ".*?icon_256x256\\.png$")
}
