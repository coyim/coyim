package xmpp

import (
	. "gopkg.in/check.v1"
)

type MUCSuite struct{}

var _ = Suite(&MUCSuite{})

func (s *MUCSuite) Test_CanJoinRoom(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		out:  mockOut,
		rand: &mockConnIOReaderWriter{read: []byte("123555111654")},
	}

	err := conn.enterRoom("coyim", "chat.coy.im", "i_am_coy")
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Equals, "<presence "+
		"id='3544672884359377457' "+
		"to='coyim@chat.coy.im/i_am_coy' "+
		"type=''>"+
		"<x xmlns='http://jabber.org/protocol/muc'/>"+
		"</presence>")
}
