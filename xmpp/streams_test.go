package xmpp

import (
	"encoding/xml"
	"errors"

	. "gopkg.in/check.v1"
)

type StreamsXmppSuite struct{}

var _ = Suite(&StreamsXmppSuite{})

func (s *StreamsXmppSuite) Test_sendInitialStreamHeader_returnsErrorIfSomethingGoesWrongWithFmtPrintf(c *C) {
	conn := Conn{
		out:          &mockConnIOReaderWriter{err: errors.New("Hello")},
		originDomain: "foo.com",
	}

	err := conn.sendInitialStreamHeader()
	c.Assert(err, Not(IsNil))
}

func (s *StreamsXmppSuite) Test_sendInitialStreamHeader_returnsErrorIfSomethingGoesWrongWithReadingAStream(c *C) {
	mockIn := &mockConnIOReaderWriter{err: errors.New("Hello")}
	conn := Conn{
		out:          &mockConnIOReaderWriter{},
		in:           xml.NewDecoder(mockIn),
		originDomain: "foo.com",
	}

	err := conn.sendInitialStreamHeader()
	c.Assert(err, Not(IsNil))
}

func (s *StreamsXmppSuite) Test_sendInitialStreamHeader_sendsInitialStreamHeaderToOutput(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{err: errors.New("Hello")}
	conn := Conn{
		out:          mockOut,
		in:           xml.NewDecoder(mockIn),
		originDomain: "somewhere.org",
	}
	conn.sendInitialStreamHeader()
	c.Assert(string(mockOut.write), Equals, "<?xml version='1.0'?><stream:stream to='somewhere.org' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n")
}

func (s *StreamsXmppSuite) Test_sendInitialStreamHeader_expectsResponseStreamHeaderInReturn(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><stream:stream xmlns:stream='http://etherx.jabber.org/streams' version='1.0'></stream:stream>")}
	conn := Conn{
		out:          mockOut,
		in:           xml.NewDecoder(mockIn),
		originDomain: "somewhereElse.org",
	}

	err := conn.sendInitialStreamHeader()
	c.Assert(err.Error(), Equals, "xmpp: error to unmarshal <features>: EOF")
}

func (s *StreamsXmppSuite) Test_sendInitialStreamHeader_failsIfReturnedStreamIsNotCorrectNamespace(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><str:stream xmlns:str='http://etherx.jabber.org/streams2' version='1.0'>")}
	conn := Conn{
		out:          mockOut,
		in:           xml.NewDecoder(mockIn),
		originDomain: "somewhereElse.org",
	}

	err := conn.sendInitialStreamHeader()
	c.Assert(err.Error(), Equals, "xmpp: expected <stream> but got <stream> in http://etherx.jabber.org/streams2")
}

func (s *StreamsXmppSuite) Test_sendInitialStreamHeader_failsIfReturnedElementIsNotStream(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><str:feature xmlns:str='http://etherx.jabber.org/streams' version='1.0'>")}
	conn := Conn{
		out:          mockOut,
		in:           xml.NewDecoder(mockIn),
		originDomain: "somewhereElse.org",
	}

	err := conn.sendInitialStreamHeader()
	c.Assert(err.Error(), Equals, "xmpp: expected <stream> but got <feature> in http://etherx.jabber.org/streams")
}

func (s *StreamsXmppSuite) Test_sendInitialStreamHeader_expectsFeaturesInReturn(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'><str:features></str:features>")}
	conn := Conn{
		out:          mockOut,
		in:           xml.NewDecoder(mockIn),
		originDomain: "somewhereElse.org",
	}

	err := conn.sendInitialStreamHeader()
	c.Assert(err, IsNil)
	expected := streamFeatures{}
	expected.XMLName = xml.Name{Space: "http://etherx.jabber.org/streams", Local: "features"}
	c.Assert(conn.features, DeepEquals, expected)
}

func (s *StreamsXmppSuite) Test_sendInitialStreamHeader_receiveResponseStreamHeaderInReturn(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte(`
	<?xml version='1.0'?>
	<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>
		<str:features>
			<bind xmlns='urn:ietf:params:xml:ns:xmpp-bind' />
			<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'>
				<required/>
			</starttls>
			<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>
				<mechanism>PLAIN</mechanism>
				<mechanism>X-OAUTH2</mechanism>
				<mechanism>X-GOOGLE-TOKEN</mechanism>
			</mechanisms>
			<register xmlns='http://jabber.org/features/iq-register'/>
		</str:features>
	`)}
	conn := Conn{
		out:          mockOut,
		in:           xml.NewDecoder(mockIn),
		originDomain: "somewhereElse.org",
	}

	err := conn.sendInitialStreamHeader()
	c.Assert(err, IsNil)
	expected := streamFeatures{
		XMLName: xml.Name{Space: "http://etherx.jabber.org/streams", Local: "features"},
		Bind: bindBind{
			XMLName: xml.Name{Space: "urn:ietf:params:xml:ns:xmpp-bind", Local: "bind"}, Resource: "", Jid: "",
		},
		StartTLS: tlsStartTLS{
			XMLName:  xml.Name{Space: "urn:ietf:params:xml:ns:xmpp-tls", Local: "starttls"},
			Required: xml.Name{Space: "urn:ietf:params:xml:ns:xmpp-tls", Local: "required"},
		},
		Mechanisms: saslMechanisms{
			XMLName:   xml.Name{Space: "urn:ietf:params:xml:ns:xmpp-sasl", Local: "mechanisms"},
			Mechanism: []string{"PLAIN", "X-OAUTH2", "X-GOOGLE-TOKEN"},
		},
		InBandRegistration: &inBandRegistration{xml.Name{Space: "http://jabber.org/features/iq-register", Local: "register"}},
	}

	c.Assert(conn.features, DeepEquals, expected)
}
