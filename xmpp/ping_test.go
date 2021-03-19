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

func (s *PingSuite) Test_conn_ReceivePong(c *C) {
	t1 := getTimeNowWithPrecission()
	conn := conn{
		lastPongResponse: t1,
	}

	conn.ReceivePong()
	c.Assert(t1, Not(Equals), conn.lastPongResponse)
}

func (s *PingSuite) Test_ParsePong_failsIfStanzaIsNotIQ(c *C) {
	rp := data.Stanza{
		Value: "hello",
	}

	e := ParsePong(rp)
	c.Assert(e, ErrorMatches, "xmpp: ping request resulted in tag of type.*")
}

func (s *PingSuite) Test_ParsePong_failsIfIQTypeIsError(c *C) {
	rp := data.Stanza{
		Value: &data.ClientIQ{
			Type: "error",
			Error: data.StanzaError{
				Text: "oh no",
			},
		},
	}

	e := ParsePong(rp)
	c.Assert(e, ErrorMatches, "xmpp: ping request resulted in an error: oh no")
}

func (s *PingSuite) Test_ParsePong_failsIfIQTypeIsUnknown(c *C) {
	rp := data.Stanza{
		Value: &data.ClientIQ{
			Type: "bla",
		},
	}

	e := ParsePong(rp)
	c.Assert(e, ErrorMatches, "xmpp: ping request resulted in an unexpected type")
}

func (s *PingSuite) Test_ParsePong_worksWithCorrectType(c *C) {
	rp := data.Stanza{
		Value: &data.ClientIQ{
			Type: "result",
		},
	}

	e := ParsePong(rp)
	c.Assert(e, IsNil)
}
