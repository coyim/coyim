package xmpp

import (
	"encoding/xml"
	"errors"

	"github.com/coyim/coyim/sasl"
	"github.com/coyim/coyim/servers"
	"github.com/coyim/coyim/xmpp/data"
	xe "github.com/coyim/coyim/xmpp/errors"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

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
	c.Assert(e, Equals, xe.ErrAuthenticationFailed)
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

func (s *SaslXMPPSuite) Test_conn_AuthenticationFailure(c *C) {
	cn := &conn{}
	c.Assert(cn.AuthenticationFailure(), Equals, xe.ErrAuthenticationFailed)

	cn.features = data.StreamFeatures{
		Mechanisms: data.SaslMechanisms{
			Mechanism: []string{
				"X-GOOGLE-TOKEN",
			},
		},
	}
	c.Assert(cn.AuthenticationFailure(), Equals, xe.ErrGoogleAuthenticationFailed)
}

func (s *SaslXMPPSuite) Test_conn_authenticateWithPreferedMethod_handlesBrokenServer(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	brokenRand := &mockConnIOReaderWriter{
		err: errors.New("io rand error"),
	}

	cn := &conn{
		log:  l,
		rand: brokenRand,
	}
	cn.known = &servers.Server{BrokenSCRAM: true}
	cn.features = data.StreamFeatures{
		Mechanisms: data.SaslMechanisms{
			Mechanism: []string{
				"SCRAM-SHA-512-PLUS",
				"SCRAM-SHA-512",
				"SCRAM-SHA-256-PLUS",
				"SCRAM-SHA-256",
				"SCRAM-SHA-1-PLUS",
				"SCRAM-SHA-1",
				"DIGEST-MD5",
			},
		},
	}

	e := cn.authenticateWithPreferedMethod("foo", "bar")

	c.Assert(e, ErrorMatches, "EOF")
	c.Assert(len(hook.Entries), Equals, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "sasl: server supports mechanisms")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[1].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[1].Message, Equals, "sasl: authenticating via")
	c.Assert(hook.Entries[1].Data, HasLen, 1)
	c.Assert(hook.Entries[1].Data["mechanism"], Equals, "DIGEST-MD5")
}

func (s *SaslXMPPSuite) Test_conn_BindResource(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{
		read: []byte(`<iq xmlns="jabber:client">
<bind xmlns="urn:ietf:params:xml:ns:xmpp-bind">
  <jid>hello@bar.com/foo</jid>
  <resource>foobar</resource>
</bind>
</iq>`),
	}
	cn := &conn{
		resource: "hello",
		out:      mockOut,
		in:       xml.NewDecoder(mockIn),
	}

	e := cn.BindResource()
	c.Assert(e, IsNil)
	c.Assert(cn.resource, Equals, "foo")
	c.Assert(string(mockOut.write), Equals, "<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'><resource>hello</resource></bind></iq>")
}

type testSaslMechanism struct {
	setPropertyReturn error
	stepReturn        error
}

func (m *testSaslMechanism) NewClient() sasl.Session {
	return m
}

func (m *testSaslMechanism) SetProperty(sasl.Property, string) error {
	return m.setPropertyReturn
}

func (m *testSaslMechanism) Step(sasl.Token) (sasl.Token, error) {
	return sasl.Token{}, m.stepReturn
}

func (m *testSaslMechanism) NeedsMore() bool {
	return false
}

func (m *testSaslMechanism) SetChannelBinding([]byte) {
}

func (s *SaslXMPPSuite) Test_conn_authenticateWith_unknownMechanism(c *C) {
	cn := &conn{}

	e := cn.authenticateWith("something unknown", "hello", "foo")
	c.Assert(e, ErrorMatches, "the requested mechanism is not supported")
}

func (s *SaslXMPPSuite) Test_conn_authenticateWith_mechanismThatFailsOnFirstStep(c *C) {
	ms := &testSaslMechanism{
		stepReturn: errors.New("marker marker"),
	}
	sasl.RegisterMechanism("mechanism: Test_conn_authenticateWith_mechanismThatFailsOnFirstStep", ms)

	cn := &conn{}
	e := cn.authenticateWith("mechanism: Test_conn_authenticateWith_mechanismThatFailsOnFirstStep", "hello", "foo")
	c.Assert(e, ErrorMatches, "marker marker")
}

func (s *SaslXMPPSuite) Test_conn_challengeLoop_failsIfStepFails(c *C) {
	ms := &testSaslMechanism{
		stepReturn: errors.New("marker marker"),
	}

	mockIn := &mockConnIOReaderWriter{
		read: []byte(`
<challenge xmlns="urn:ietf:params:xml:ns:xmpp-sasl">Zm9vIGJhcg==</challenge>
`),
	}

	cn := &conn{
		log: testLogger(),
		in:  xml.NewDecoder(mockIn),
	}
	e := cn.challengeLoop(ms)
	c.Assert(e, ErrorMatches, "marker marker")
}
