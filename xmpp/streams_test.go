package xmpp

import (
	"encoding/xml"
	"errors"

	"github.com/twstrike/coyim/xmpp/data"

	. "gopkg.in/check.v1"
)

type StreamsXMPPSuite struct{}

var _ = Suite(&StreamsXMPPSuite{})

func (s *StreamsXMPPSuite) Test_sendInitialStreamHeader_returnsErrorIfSomethingGoesWrongWithFmtPrintf(c *C) {
	conn := conn{
		out:          &mockConnIOReaderWriter{err: errors.New("Hello")},
		originDomain: "foo.com",
	}

	err := conn.SendInitialStreamHeader()
	c.Assert(err, Not(IsNil))
}

func (s *StreamsXMPPSuite) Test_sendInitialStreamHeader_returnsErrorIfSomethingGoesWrongWithReadingAStream(c *C) {
	mockIn := &mockConnIOReaderWriter{err: errors.New("Hello")}
	conn := conn{
		out:          &mockConnIOReaderWriter{},
		in:           xml.NewDecoder(mockIn),
		originDomain: "foo.com",
	}

	err := conn.SendInitialStreamHeader()
	c.Assert(err, Not(IsNil))
}

func (s *StreamsXMPPSuite) Test_sendInitialStreamHeader_sendsInitialStreamHeaderToOutput(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{err: errors.New("Hello")}
	conn := conn{
		out:          mockOut,
		in:           xml.NewDecoder(mockIn),
		originDomain: "somewhere.org",
	}
	conn.SendInitialStreamHeader()
	c.Assert(string(mockOut.write), Equals, "<?xml version='1.0'?><stream:stream to='somewhere.org' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n")
}

func (s *StreamsXMPPSuite) Test_sendInitialStreamHeader_expectsResponseStreamHeaderInReturn(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><stream:stream xmlns:stream='http://etherx.jabber.org/streams' version='1.0'></stream:stream>")}
	conn := conn{
		out:          mockOut,
		in:           xml.NewDecoder(mockIn),
		originDomain: "somewhereElse.org",
	}

	err := conn.SendInitialStreamHeader()
	c.Assert(err.Error(), Equals, "xmpp: error to unmarshal <features>: EOF")
}

func (s *StreamsXMPPSuite) Test_sendInitialStreamHeader_failsIfReturnedStreamIsNotCorrectNamespace(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><str:stream xmlns:str='http://etherx.jabber.org/streams2' version='1.0'>")}
	conn := conn{
		out:          mockOut,
		in:           xml.NewDecoder(mockIn),
		originDomain: "somewhereElse.org",
	}

	err := conn.SendInitialStreamHeader()
	c.Assert(err.Error(), Equals, "xmpp: expected <stream> but got <stream> in http://etherx.jabber.org/streams2")
}

func (s *StreamsXMPPSuite) Test_sendInitialStreamHeader_failsIfReturnedElementIsNotStream(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><str:feature xmlns:str='http://etherx.jabber.org/streams' version='1.0'>")}
	conn := conn{
		out:          mockOut,
		in:           xml.NewDecoder(mockIn),
		originDomain: "somewhereElse.org",
	}

	err := conn.SendInitialStreamHeader()
	c.Assert(err.Error(), Equals, "xmpp: expected <stream> but got <feature> in http://etherx.jabber.org/streams")
}

func (s *StreamsXMPPSuite) Test_sendInitialStreamHeader_expectsFeaturesInReturn(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'><str:features></str:features>")}
	conn := conn{
		out:          mockOut,
		in:           xml.NewDecoder(mockIn),
		originDomain: "somewhereElse.org",
	}

	err := conn.SendInitialStreamHeader()
	c.Assert(err, IsNil)
	expected := data.StreamFeatures{}
	expected.XMLName = xml.Name{Space: "http://etherx.jabber.org/streams", Local: "features"}
	c.Assert(conn.features, DeepEquals, expected)
}

func (s *StreamsXMPPSuite) Test_sendInitialStreamHeader_receiveResponseStreamHeaderInReturn(c *C) {
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
	conn := conn{
		out:          mockOut,
		in:           xml.NewDecoder(mockIn),
		originDomain: "somewhereElse.org",
	}

	err := conn.SendInitialStreamHeader()
	c.Assert(err, IsNil)
	expected := data.StreamFeatures{
		XMLName: xml.Name{Space: "http://etherx.jabber.org/streams", Local: "features"},
		Bind: data.BindBind{
			XMLName: xml.Name{Space: "urn:ietf:params:xml:ns:xmpp-bind", Local: "bind"}, Resource: "", Jid: "",
		},
		StartTLS: data.StartTLS{
			XMLName:  xml.Name{Space: "urn:ietf:params:xml:ns:xmpp-tls", Local: "starttls"},
			Required: xml.Name{Space: "urn:ietf:params:xml:ns:xmpp-tls", Local: "required"},
		},
		Mechanisms: data.SaslMechanisms{
			XMLName:   xml.Name{Space: "urn:ietf:params:xml:ns:xmpp-sasl", Local: "mechanisms"},
			Mechanism: []string{"PLAIN", "X-OAUTH2", "X-GOOGLE-TOKEN"},
		},
		InBandRegistration: &data.InBandRegistration{xml.Name{Space: "http://jabber.org/features/iq-register", Local: "register"}},
	}

	c.Assert(conn.features, DeepEquals, expected)
}
