package xmpp

import (
	"encoding/xml"
	"io"

	. "gopkg.in/check.v1"
)

type ConnectionXmppSuite struct{}

var _ = Suite(&ConnectionXmppSuite{})

func (s *ConnectionXmppSuite) Test_Next_returnsErrorIfOneIsEncountered(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:foo xmlns:stream='http://etherx.jabber.org/streams' to='hello'></stream:foo>")}
	conn := Conn{
		in: xml.NewDecoder(mockIn),
	}

	_, err := conn.Next()
	c.Assert(err.Error(), Equals, "unexpected XMPP message http://etherx.jabber.org/streams <foo/>")
}

func (s *ConnectionXmppSuite) Test_Next_returnsErrorIfFailingToParseIQID(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='abczzzz'></client:iq>")}
	conn := Conn{
		in: xml.NewDecoder(mockIn),
	}

	_, err := conn.Next()
	c.Assert(err.Error(), Equals, "xmpp: failed to parse id from iq: strconv.ParseUint: parsing \"abczzzz\": invalid syntax")
}

func (s *ConnectionXmppSuite) Test_Next_returnsNothingIfThereIsNoInflightMatching(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000'></client:iq>")}
	conn := Conn{
		in: xml.NewDecoder(mockIn),
	}

	_, err := conn.Next()
	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXmppSuite) Test_makeInOut_returnsANewDecoderAndOriginalWriterWhenNoConfigIsGiven(c *C) {
	mockBoth := &mockConnIOReaderWriter{}
	_, rout := makeInOut(mockBoth, nil)
	c.Assert(rout, Equals, mockBoth)
}

func (s *ConnectionXmppSuite) Test_makeInOut_returnsANewDecoderAndWrappedWriterWhenConfigIsGiven(c *C) {
	mockBoth := &mockConnIOReaderWriter{}
	mockInLog := &mockConnIOReaderWriter{}
	config := &Config{InLog: mockInLog, OutLog: mockInLog}
	_, rout := makeInOut(mockBoth, config)
	c.Assert(rout, Not(Equals), mockBoth)
}

func (s *ConnectionXmppSuite) Test_Dial_returnsErrorFromGetFeatures(c *C) {
	rw := &mockConnIOReaderWriter{}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXmppSuite) Test_Dial_returnsErrorFromAuthenticateIfSkipTLS(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'><str:features></str:features>")}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "xmpp: PLAIN authentication is not an option")
}

func (s *ConnectionXmppSuite) Test_Dial_returnsErrorFromSecondFeatureCheck(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>")}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err, Equals, io.EOF)
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n")
}

func (s *ConnectionXmppSuite) Test_Dial_returnsErrorFromIQReturn(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>" +
			"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"</str:features>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "unmarshal <iq>: EOF")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_returnsWorkingConnIfEverythingPasses(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>" +
			"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"</str:features>" +
			"<client:iq xmlns:client='jabber:client'></client:iq>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err, IsNil)
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_failsIfTheServerDoesntSupportTLS(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: false}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "xmpp: server doesn't support TLS")
}

func (s *ConnectionXmppSuite) Test_Dial_failsIfReceivingEOFAfterStartingTLS(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'>" +
			"</starttls>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: false}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXmppSuite) Test_Dial_failsIfReceivingTheWrongNamespaceAfterStarttls(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'>" +
			"</starttls>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<str:proceed>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: false}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "xmpp: expected <proceed> after <starttls> but got <proceed> in http://etherx.jabber.org/streams")
}

func (s *ConnectionXmppSuite) Test_Dial_failsIfReceivingTheWrongTagName(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'>" +
			"</starttls>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<things xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: false}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "xmpp: expected <proceed> after <starttls> but got <things> in urn:ietf:params:xml:ns:xmpp-tls")
}

func (s *ConnectionXmppSuite) Test_Dial_failsWhenStartingAHandshake(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'>" +
			"</starttls>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<proceed xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: false}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXmppSuite) Test_Dial_failsIfDecodingFallbackFails(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true, CreateCallback: func(title, instructions string, fields []interface{}) error {
		return nil
	}}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "unmarshal <iq>: EOF")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_failsIfAccountCreationFails(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<iq xmlns='jabber:client' type='something'></iq>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true, CreateCallback: func(title, instructions string, fields []interface{}) error {
		return nil
	}}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "xmpp: account creation failed")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>",
	)
}
