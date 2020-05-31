package xmpp

import (
	"encoding/xml"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/errors"

	. "gopkg.in/check.v1"
)

type SaslXMPPSuite struct{}

var _ = Suite(&SaslXMPPSuite{})

func (s *SaslXMPPSuite) Test_authenticate_failsIfPlainIsNotAnOption(c *C) {
	conn := conn{log: testLogger()}

	err := conn.Authenticate("", "")
	c.Assert(err, Equals, errUnsupportedSASLMechanism)
}

func (s *SaslXMPPSuite) Test_authenticate_authenticatesWithUsernameAndPassword(c *C) {
	out := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>")}
	conn := conn{
		log:    testLogger(),
		rawOut: out,
		in:     xml.NewDecoder(mockIn),
		features: data.StreamFeatures{
			Mechanisms: data.SaslMechanisms{
				Mechanism: []string{"FOO", "PLAIN"},
			},
		},
	}

	e := conn.Authenticate("foo", "bar")
	c.Assert(e, IsNil)
	c.Assert(string(out.write), Equals, "<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AGZvbwBiYXI=</auth>\n")
}

func (s *SaslXMPPSuite) Test_authenticate_handlesFailure(c *C) {
	out := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<sasl:failure xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'><foobar></foobar></sasl:failure>")}
	conn := conn{
		log:    testLogger(),
		rawOut: out,
		in:     xml.NewDecoder(mockIn),
		features: data.StreamFeatures{
			Mechanisms: data.SaslMechanisms{
				Mechanism: []string{"FOO", "PLAIN"},
			},
		},
	}

	e := conn.Authenticate("foo", "bar")
	c.Assert(e.Error(), Equals, "xmpp: authentication failure: foobar")
}

func (s *SaslXMPPSuite) Test_authenticate_handlesWrongResponses(c *C) {
	out := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<sasl:something xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:something>")}
	conn := conn{
		log:    testLogger(),
		rawOut: out,
		in:     xml.NewDecoder(mockIn),
		features: data.StreamFeatures{
			Mechanisms: data.SaslMechanisms{
				Mechanism: []string{"FOO", "PLAIN"},
			},
		},
	}

	e := conn.Authenticate("foo", "bar")
	c.Assert(e, Equals, errors.ErrAuthenticationFailed)
}

func (s *SaslXMPPSuite) Test_digestMD5_authenticatesWithUsernameAndPassword(c *C) {
	out := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte(
		"<challenge xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>cmVhbG09ImNveS5pbSIsbm9uY2U9Ik9BNk1HOXRFUUdtMmhoIixxb3A9ImF1dGgiLGNoYXJzZXQ9dXRmLTgsYWxnb3JpdGhtPW1kNS1zZXNz</challenge>\n" +
			"<challenge xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>cnNwYXV0aD1lYTQwZjYwMzM1YzQyN2I1NTI3Yjg0ZGJhYmNkZmZmZA==</challenge>\n" +
			"<success xmlns='urn:ietf:params:xml:ns:xmpp-sasl'/>\n",
	)}

	mockRand := &mockConnIOReaderWriter{read: []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
	}}

	conn := conn{
		log:    testLogger(),
		rawOut: out,
		in:     xml.NewDecoder(mockIn),
		rand:   mockRand,
		features: data.StreamFeatures{
			Mechanisms: data.SaslMechanisms{
				Mechanism: []string{"DIGEST-MD5"},
			},
		},
	}

	expectedOut := "<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='DIGEST-MD5'></auth>\n" +
		"<response xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>Y2hhcnNldD11dGYtOCx1c2VybmFtZT0iZm9vIixyZWFsbT0iY295LmltIixub25jZT0iT0E2TUc5dEVRR20yaGgiLG5jPTAwMDAwMDAxLGNub25jZT0iMDEwMjAzMDQwNTA2MDciLGRpZ2VzdC11cmk9InhtcHAvY295LmltIixyZXNwb25zZT00ZGVlODYyNjkxOTZiNmUxNGI5Zjc2OWZhYmQ5OTdiZCxxb3A9YXV0aA==</response>\n" +
		"<response xmlns='urn:ietf:params:xml:ns:xmpp-sasl'></response>\n"

	e := conn.Authenticate("foo", "bar")
	c.Assert(e, IsNil)
	c.Assert(string(out.write), Equals, expectedOut)
}

func (s *SaslXMPPSuite) Test_digestMD5_serverFailsToVerifyChallenge(c *C) {
	out := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte(
		"<challenge xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>cmVhbG09ImNveS5pbSIsbm9uY2U9Ik9BNk1HOXRFUUdtMmhoIixxb3A9ImF1dGgiLGNoYXJzZXQ9dXRmLTgsYWxnb3JpdGhtPW1kNS1zZXNz</challenge>\n" +
			"<response xmlns='urn:ietf:params:xml:ns:xmpp-sasl'/>",
	)}

	mockRand := &mockConnIOReaderWriter{read: []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
	}}

	conn := conn{
		log:    testLogger(),
		rawOut: out,
		in:     xml.NewDecoder(mockIn),
		rand:   mockRand,
		features: data.StreamFeatures{
			Mechanisms: data.SaslMechanisms{
				Mechanism: []string{"DIGEST-MD5"},
			},
		},
	}

	expectedOut := "<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='DIGEST-MD5'></auth>\n" +
		"<response xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>Y2hhcnNldD11dGYtOCx1c2VybmFtZT0iZm9vIixyZWFsbT0iY295LmltIixub25jZT0iT0E2TUc5dEVRR20yaGgiLG5jPTAwMDAwMDAxLGNub25jZT0iMDEwMjAzMDQwNTA2MDciLGRpZ2VzdC11cmk9InhtcHAvY295LmltIixyZXNwb25zZT00ZGVlODYyNjkxOTZiNmUxNGI5Zjc2OWZhYmQ5OTdiZCxxb3A9YXV0aA==</response>\n"

	e := conn.Authenticate("foo", "bar")
	c.Assert(e.Error(), Equals, "xmpp: unexpected <response> in urn:ietf:params:xml:ns:xmpp-sasl")
	c.Assert(string(out.write), Equals, expectedOut)
}

func (s *SaslXMPPSuite) Test_scramSHA1Auth_authenticatesWithUsernameAndPassword(c *C) {
	out := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte(
		"<challenge xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>cj1iNWNmZjYxOTAwMTNlNmttdWE1REVtUEFaak9NcHE0VEhXSlE9PSxzPURrRVdNMjBxRTE5c3V2ckhoUHI3SEE9PSxpPTQwOTY=</challenge>\n" +
			"<success xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>dj1rNW41OTVxVzUwVHlFMnErSjBjVWY5eVQ4djQ9</success>\n",
	)}

	mockRand := &mockConnIOReaderWriter{read: []byte{
		0xb5, 0xcf, 0xf6, 0x19, 0x00, 0x13, 0xe6,
	}}

	conn := conn{
		log:    testLogger(),
		rawOut: out,
		in:     xml.NewDecoder(mockIn),
		rand:   mockRand,
		features: data.StreamFeatures{
			Mechanisms: data.SaslMechanisms{
				Mechanism: []string{"SCRAM-SHA-1"},
			},
		},
	}

	expectedOut := "<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='SCRAM-SHA-1'>eSwsbj11c2VyLHI9YjVjZmY2MTkwMDEzZTY=</auth>\n" +
		"<response xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>Yz1lU3dzLHI9YjVjZmY2MTkwMDEzZTZrbXVhNURFbVBBWmpPTXBxNFRIV0pRPT0scD1SZnFnNDlqYkJmMWJHQ2t3RlRiby9EdkhtVUk9</response>\n"

	e := conn.Authenticate("user", "pencil")
	c.Assert(e, IsNil)
	c.Assert(string(out.write), Equals, expectedOut)
}
