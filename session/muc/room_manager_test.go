package muc

import (
	"github.com/coyim/coyim/xmpp/jid"
	. "gopkg.in/check.v1"
)

type MucRoomManagerSuite struct{}

var _ = Suite(&MucRoomManagerSuite{})

func (*MucRoomManagerSuite) Test_NewRoomManager(c *C) {
	rr := NewRoomManager()
	c.Assert(rr.rooms, Not(IsNil))
}

func (*MucRoomManagerSuite) Test_RoomManager_GetRoom(c *C) {
	rr := NewRoomManager()
	room := &Room{}
	rr.rooms["foo@bar.com"] = room

	nr, ok := rr.GetRoom(jid.ParseBare("foo@bar.com"))

	c.Assert(ok, Equals, true)
	c.Assert(nr, Equals, room)
}

func (*MucRoomManagerSuite) Test_RoomManager_AddRoom(c *C) {
	rr := NewRoomManager()

	ok := rr.AddRoom(&Room{ID: jid.ParseBare("foo@bar.com")})
	c.Assert(ok, Equals, true)

	ok = rr.AddRoom(&Room{ID: jid.ParseBare("foo@bar.com")})
	c.Assert(ok, Equals, false)
}

func (*MucRoomManagerSuite) Test_RoomManager_LeaveRoom(c *C) {
	rr := NewRoomManager()
	_ = rr.AddRoom(&Room{ID: jid.ParseBare("foo@bar.com")})

	ok := rr.LeaveRoom(jid.ParseBare("somewhere@bar.com"))
	c.Assert(ok, Equals, false)

	ok = rr.LeaveRoom(jid.ParseBare("foo@bar.com"))
	c.Assert(ok, Equals, true)
}
