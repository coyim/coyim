package xmpp

import (
	"errors"

	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

type PresenceXMPPSuite struct{}

var _ = Suite(&PresenceXMPPSuite{})

func (s *PresenceXMPPSuite) Test_SignalPresence_sendsPresenceInformation(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		log: testLogger(),
		out: mockOut,
	}

	err := conn.SignalPresence("fo'o")
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Equals, "<presence><show>fo&apos;o</show></presence>")
}

func (s *PresenceXMPPSuite) Test_SignalPresence_returnsWriterError(c *C) {
	mockOut := &mockConnIOReaderWriter{err: errors.New("foo bar")}
	conn := conn{
		log: testLogger(),
		out: mockOut,
	}

	err := conn.SignalPresence("fo'o")
	c.Assert(err.Error(), Equals, "foo bar")
}

func (s *PresenceXMPPSuite) Test_SendPresence_sendsPresenceWithTheIdGiven(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		log: testLogger(),
		out: mockOut,
	}

	expectedPresence := `<presence xmlns="jabber:client" id="123456&amp;" to="someone&lt;strange&gt;@foo.com" type="subsc&#39;ribe"></presence>`

	err := conn.SendPresence("someone<strange>@foo.com", "subsc'ribe", "123456&", "")
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Equals, expectedPresence)
}

func (s *PresenceXMPPSuite) Test_SendPresence_sendsPresenceWithRandomID(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		log:  testLogger(),
		out:  mockOut,
		rand: &mockConnIOReaderWriter{read: []byte("123555111654")},
	}

	expectedPresence := `<presence xmlns="jabber:client" id="3544672884359377457" to="someone&lt;strange&gt;@foo.com" type="subsc&#39;ribe"></presence>`
	err := conn.SendPresence("someone<strange>@foo.com", "subsc'ribe", "", "")
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Equals, expectedPresence)
}

func (s *PresenceXMPPSuite) Test_SendPresence_returnsWriterError(c *C) {
	mockOut := &mockConnIOReaderWriter{err: errors.New("bar foo")}
	conn := conn{
		log: testLogger(),
		out: mockOut,
	}

	err := conn.SendPresence("someone<strange>@foo.com", "subsc'ribe", "abc", "")
	c.Assert(err.Error(), Equals, "bar foo")
}

func (s *PresenceXMPPSuite) Test_SendPresence_addsStatusToSubscribeMessage(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		log: testLogger(),
		out: mockOut,
	}

	expectedPresence := `<presence xmlns="jabber:client" id="123" to="someone&lt;strange&gt;@foo.com" type="subscribe"><status>do you want &lt;to&gt;?</status></presence>`
	err := conn.SendPresence("someone<strange>@foo.com", "subscribe", "123", "do you want <to>?")
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Equals, expectedPresence)
}

func (s *PresenceXMPPSuite) Test_conn_SendMUCPresence_works(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	cn := conn{
		log: testLogger(),
		out: mockOut,
	}

	err := cn.SendMUCPresence("fo'o", &data.MUC{
		Password: "bla",
	})
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Matches, "<presence xmlns=\"jabber:client\" id=\".+\" to=\"fo&#39;o\"><x xmlns=\"http://jabber.org/protocol/muc\"><password>bla</password></x></presence>")
}
