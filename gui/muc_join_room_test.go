package gui

import (
	. "gopkg.in/check.v1"
)

type JoinMUCRoomSuite struct{}

var _ = Suite(&JoinMUCRoomSuite{})

func (*SignalsSuite) Test_joinRoom_validRoomName(c *C) {
	v := &mucJoinRoomView{}

	c.Assert(v.isValidRoomName("john@server.com"), Equals, true)
	c.Assert(v.isValidRoomName(" john@server.com"), Equals, false)
	c.Assert(v.isValidRoomName("@server.com"), Equals, false)
	c.Assert(v.isValidRoomName("#@server.com"), Equals, true)
	c.Assert(v.isValidRoomName("#@server.com "), Equals, false)
	c.Assert(v.isValidRoomName("john@server"), Equals, true)
	c.Assert(v.isValidRoomName("john doe@server"), Equals, false)
	c.Assert(v.isValidRoomName("johnD0%$2doe@server"), Equals, true)
	c.Assert(v.isValidRoomName("john&doe@server"), Equals, false)
}
