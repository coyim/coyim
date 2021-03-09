package xmpp

import (
	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

type PingSuite struct{}

var _ = Suite(&PingSuite{})

func (s *PingSuite) Test_conn_SendPing(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		log: testLogger(),
		out: mockOut,
		jid: "juliet@example.com/chamber",
	}
	conn.inflights = make(map[data.Cookie]inflight)

	reply, cookie, err := conn.SendPing()
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Matches, "<iq xmlns='jabber:client' from='juliet@example.com/chamber' type='get' id='.+'><ping xmlns=\"urn:xmpp:ping\"></ping></iq>")
	c.Assert(reply, NotNil)
	c.Assert(cookie, NotNil)
}

func (s *PingSuite) Test_conn_SendPingReply(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		log: testLogger(),
		out: mockOut,
		jid: "juliet@example.com/chamber",
	}
	conn.inflights = make(map[data.Cookie]inflight)

	err := conn.SendPingReply("huff")
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Matches, "<iq xmlns='jabber:client' from='juliet@example.com/chamber' type='result' id='huff'></iq>")
}
