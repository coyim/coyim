package xmpp

import (
	"encoding/xml"

	. "gopkg.in/check.v1"
)

type SaslXmppSuite struct{}

var _ = Suite(&SaslXmppSuite{})

func (s *SaslXmppSuite) Test_authenticate_failsIfPlainIsNotAnOption(c *C) {
	conn := Conn{}

	e := conn.authenticate("", "")

	c.Assert(e.Error(), Equals, "xmpp: PLAIN authentication is not an option")
}

func (s *SaslXmppSuite) Test_authenticate_authenticatesWithUsernameAndPassword(c *C) {
	out := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>")}
	conn := Conn{
		rawOut: out,
		in:     xml.NewDecoder(mockIn),
		features: streamFeatures{
			Mechanisms: saslMechanisms{
				Mechanism: []string{"FOO", "PLAIN"},
			},
		},
	}

	e := conn.authenticate("foo", "bar")
	c.Assert(e, IsNil)
	c.Assert(string(out.write), Equals, "<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AGZvbwBiYXI=</auth>\n")
}

func (s *SaslXmppSuite) Test_authenticate_handlesFailure(c *C) {
	out := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<sasl:failure xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'><foobar></foobar></sasl:failure>")}
	conn := Conn{
		rawOut: out,
		in:     xml.NewDecoder(mockIn),
		features: streamFeatures{
			Mechanisms: saslMechanisms{
				Mechanism: []string{"FOO", "PLAIN"},
			},
		},
	}

	e := conn.authenticate("foo", "bar")
	c.Assert(e.Error(), Equals, "xmpp: authentication failure: foobar")
}

func (s *SaslXmppSuite) Test_authenticate_handlesWrongResponses(c *C) {
	out := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<sasl:something xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:something>")}
	conn := Conn{
		rawOut: out,
		in:     xml.NewDecoder(mockIn),
		features: streamFeatures{
			Mechanisms: saslMechanisms{
				Mechanism: []string{"FOO", "PLAIN"},
			},
		},
	}

	e := conn.authenticate("foo", "bar")
	c.Assert(e.Error(), Equals, "expected <success> or <failure>, got <> in ")
}
