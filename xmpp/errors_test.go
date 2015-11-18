package xmpp

import (
	"encoding/xml"

	. "gopkg.in/check.v1"
)

type ErrorsSuite struct{}

var _ = Suite(&ErrorsSuite{})

func (s *ErrorsSuite) Test_StreamError_Marshalling(c *C) {
	streamErr := StreamError{
		DefinedCondition: ConnectionTimeout,
	}

	expected := `<error xmlns="http://etherx.jabber.org/streams"><connection-timeout xmlns="urn:ietf:params:xml:ns:xmpp-streams"></connection-timeout></error>`
	enc, err := xml.Marshal(streamErr)

	c.Check(err, IsNil)
	c.Check(string(enc), Equals, expected)
}
