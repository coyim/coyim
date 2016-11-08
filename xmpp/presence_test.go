package xmpp

import (
	"errors"

	. "gopkg.in/check.v1"
)

type PresenceXMPPSuite struct{}

var _ = Suite(&PresenceXMPPSuite{})

func (s *PresenceXMPPSuite) Test_SignalPresence_sendsPresenceInformation(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		out: mockOut,
	}

	err := conn.SignalPresence("fo'o")
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Equals, "<presence><show>fo&apos;o</show></presence>")
}

func (s *PresenceXMPPSuite) Test_SignalPresence_returnsWriterError(c *C) {
	mockOut := &mockConnIOReaderWriter{err: errors.New("foo bar")}
	conn := conn{
		out: mockOut,
	}

	err := conn.SignalPresence("fo'o")
	c.Assert(err.Error(), Equals, "foo bar")
}

func (s *PresenceXMPPSuite) Test_SendPresence_sendsPresenceWithTheIdGiven(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		out: mockOut,
	}

	err := conn.SendPresence("someone<strange>@foo.com", "subsc'ribe", "123456&", "")
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Equals, "<presence id='123456&amp;' to='someone&lt;strange&gt;@foo.com' type='subsc&apos;ribe'/>")
}

func (s *PresenceXMPPSuite) Test_SendPresence_sendsPresenceWithRandomID(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		out:  mockOut,
		rand: &mockConnIOReaderWriter{read: []byte("123555111654")},
	}

	err := conn.SendPresence("someone<strange>@foo.com", "subsc'ribe", "", "")
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Equals, "<presence id='3544672884359377457' to='someone&lt;strange&gt;@foo.com' type='subsc&apos;ribe'/>")
}

func (s *PresenceXMPPSuite) Test_SendPresence_returnsWriterError(c *C) {
	mockOut := &mockConnIOReaderWriter{err: errors.New("bar foo")}
	conn := conn{
		out: mockOut,
	}

	err := conn.SendPresence("someone<strange>@foo.com", "subsc'ribe", "abc", "")
	c.Assert(err.Error(), Equals, "bar foo")
}

func (s *PresenceXMPPSuite) Test_SendPresence_addsStatusToSubscribeMessage(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		out: mockOut,
	}

	err := conn.SendPresence("someone<strange>@foo.com", "subscribe", "123", "do you want <to>?")
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Equals, "<presence id='123' to='someone&lt;strange&gt;@foo.com' type='subscribe'><status>do you want &lt;to&gt;?</status></presence>")
}
